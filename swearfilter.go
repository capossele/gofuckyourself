package swearfilter

import (
	"regexp"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//SwearFilter contains settings for the swear filter
type SwearFilter struct {
	//Options to tell the swear filter how to operate
	DisableNormalize                bool //Disables normalization of alphabetic characters if set to true (ex: à -> a)
	DisableSpacedTab                bool //Disables converting tabs to singular spaces (ex: [tab][tab] -> [space][space])
	DisableMultiWhitespaceStripping bool //Disables stripping down multiple whitespaces (ex: hello[space][space]world -> hello[space]world)
	DisableZeroWidthStripping       bool //Disables stripping zero-width spaces
	EnableSpacedBypass              bool //Disables testing for spaced bypasses (if hell is in filter, look for occurrences of h and detect only alphabetic characters that follow; ex: h[space]e[space]l[space]l[space] -> hell)
	DisableSimpleRegex              bool //Disables using strings.HasPrefix if the string starts with ^ and strings.HasSuffix if it ends with $. Only strings.Contains will be used
	EnableFullRegex                 bool //Enables treating each word in the wordlist as a regex

	//A list of words to check against the filters
	BadWords       map[string]struct{}
	BadWordRegexps map[string]*regexp.Regexp
	mutex          sync.RWMutex
}

//NewSwearFilter returns an initialized SwearFilter struct to check messages against
func NewSwearFilter(enableSpacedBypass bool, enableFullRegex bool, uhohwords ...string) (filter *SwearFilter) {
	filter = &SwearFilter{
		EnableSpacedBypass: enableSpacedBypass,
		EnableFullRegex:    enableFullRegex,
		BadWords:           make(map[string]struct{}),
	}

	filter.Add(uhohwords...)
	return
}

//Check will return any words that trip an enabled swear filter, an error if any, or nothing if you've removed all the words for some reason
func (filter *SwearFilter) Check(msg string) (trippedWords []string, err error) {
	filter.mutex.RLock()
	defer filter.mutex.RUnlock()

	if filter.EnableFullRegex {
		if filter.BadWordRegexps == nil || len(filter.BadWordRegexps) == 0 {
			return nil, nil
		}
	} else {
		if filter.BadWords == nil || len(filter.BadWords) == 0 {
			return nil, nil
		}
	}

	message := strings.ToLower(msg)

	//Normalize the text
	if !filter.DisableNormalize {
		bytes := make([]byte, len(message))
		normalize := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
			return unicode.Is(unicode.Mn, r)
		}), norm.NFC)
		_, _, err = normalize.Transform(bytes, []byte(message), true)
		if err != nil {
			return nil, err
		}
		message = string(bytes)
	}

	//Turn tabs into spaces
	if !filter.DisableSpacedTab {
		message = strings.Replace(message, "\t", " ", -1)
	}

	//Get rid of zero-width spaces
	if !filter.DisableZeroWidthStripping {
		message = strings.Replace(message, "\u200b", "", -1)
	}

	//Convert multiple re-occurring whitespaces into a single space
	if !filter.DisableMultiWhitespaceStripping {
		regexLeadCloseWhitepace := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
		message = regexLeadCloseWhitepace.ReplaceAllString(message, "")
		regexInsideWhitespace := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
		message = regexInsideWhitespace.ReplaceAllString(message, "")
	}

	trippedWords = make([]string, 0)
	checkSpace := false

	if filter.EnableFullRegex {
		for swear := range filter.BadWordRegexps {
			if swear == " " {
				checkSpace = true
				continue
			}

			if filter.scan(message, swear) {
				trippedWords = append(trippedWords, swear)
			} else if filter.EnableSpacedBypass {
				nospaceMessage := strings.Replace(message, " ", "", -1)
				if filter.scan(nospaceMessage, swear) {
					trippedWords = append(trippedWords, swear)
				}
			}
		}
	} else {
		for swear := range filter.BadWords {
			if swear == " " {
				checkSpace = true
				continue
			}

			if filter.scan(message, swear) {
				trippedWords = append(trippedWords, swear)
			} else if filter.EnableSpacedBypass {
				nospaceMessage := strings.Replace(message, " ", "", -1)
				if filter.scan(nospaceMessage, swear) {
					trippedWords = append(trippedWords, swear)
				}
			}
		}
	}

	if checkSpace && message == "" {
		trippedWords = append(trippedWords, " ")
	}

	return
}

func (filter *SwearFilter) scan(message string, swear string) bool {
	if filter.EnableFullRegex {
		return filter.BadWordRegexps[swear].MatchString(message)
	} else {
		if filter.DisableSimpleRegex {
			return strings.Contains(message, swear)
		} else {
			hasSimplePrefix := false
			if string(swear[0]) == "^" {
				hasSimplePrefix = true
				return strings.HasPrefix(message, swear[1:])
			}

			hasSimpleSuffix := false
			strLen := len(swear)
			if string(swear[strLen-1]) == "$" {
				hasSimpleSuffix = true
				return strings.HasSuffix(message, swear[:strLen-1])
			}

			// fallback to substring matching
			if !hasSimplePrefix && !hasSimpleSuffix {
				return strings.Contains(message, swear)
			}
		}
	}

	return false
}

//Add appends the given word to the uhohwords list
func (filter *SwearFilter) Add(badWords ...string) {
	filter.mutex.Lock()
	defer filter.mutex.Unlock()

	if filter.EnableFullRegex {
		if filter.BadWordRegexps == nil {
			filter.BadWordRegexps = make(map[string]*regexp.Regexp)
		}

		for _, word := range badWords {
			filter.BadWordRegexps[word] = regexp.MustCompile(word)
		}
	} else {
		if filter.BadWords == nil {
			filter.BadWords = make(map[string]struct{})
		}

		for _, word := range badWords {
			filter.BadWords[word] = struct{}{}
		}
	}
}

//Delete deletes the given word from the uhohwords list
func (filter *SwearFilter) Delete(badWords ...string) {
	filter.mutex.Lock()
	defer filter.mutex.Unlock()

	if filter.EnableFullRegex {
		for _, word := range badWords {
			delete(filter.BadWordRegexps, word)
		}
	} else {
		for _, word := range badWords {
			delete(filter.BadWords, word)
		}
	}
}

//Load return the uhohwords list
func (filter *SwearFilter) Load() (activeWords []string) {
	filter.mutex.RLock()
	defer filter.mutex.RUnlock()

	if filter.EnableFullRegex {
		if filter.BadWordRegexps == nil {
			return nil
		}

		for word := range filter.BadWordRegexps {
			activeWords = append(activeWords, word)
		}
	} else {
		if filter.BadWords == nil {
			return nil
		}

		for word := range filter.BadWords {
			activeWords = append(activeWords, word)
		}
	}
	return
}
