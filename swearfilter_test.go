package swearfilter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	filter := NewSwearFilter(true, "fuck", "hell")
	if filter.DisableNormalize {
		t.Errorf("Filter option DisableNormalize was incorrect, got: %t, want: %t", filter.DisableNormalize, false)
	}
	if filter.DisableSpacedTab {
		t.Errorf("Filter option DisableSpacedTab was incorrect, got: %t, want: %t", filter.DisableSpacedTab, false)
	}
	if filter.DisableMultiWhitespaceStripping {
		t.Errorf("Filter option DisableMultiWhitespaceStripping was incorrect, got: %t, want: %t", filter.DisableMultiWhitespaceStripping, false)
	}
	if filter.DisableZeroWidthStripping {
		t.Errorf("Filter option DisableZeroWidthStripping was incorrect, got: %t, want: %t", filter.DisableZeroWidthStripping, false)
	}
	if !filter.EnableSpacedBypass {
		t.Errorf("Filter option EnableSpacedBypass was incorrect, got: %t, want: %t", filter.EnableSpacedBypass, true)
	}
	if len(filter.BadWords) != 2 {
		t.Errorf("Filter option BadWords was incorrect, got length: %d, want length: %d", len(filter.BadWords), 2)
	}
	if filter.DisableSimpleRegex {
		t.Errorf("Filter option DisableSimpleRegex was incorrect, got: %t, want: %t", filter.DisableSimpleRegex, false)
	}
	if filter.EnableFullRegex {
		t.Errorf("Filter option EnableFullRegex was incorrect, got: %t, want: %t", filter.EnableFullRegex, false)
	}
}

func TestCheck(t *testing.T) {
	filter := NewSwearFilter(true, "fuck")
	messages := []string{"fucking", "fûçk", "asdf", "what the f u c k dude"}

	for i := 0; i < len(messages); i++ {
		trippers, err := filter.Check(messages[i])
		if err != nil {
			t.Errorf("Check failed due to external dependency: %v", err)
		}
		switch i {
		case 0, 1, 3:
			if len(trippers) != 1 {
				t.Errorf("Check did not act as expected, got trippers length: %d, want trippers length: %d", len(trippers), 1)
			}
			if trippers[0] != "fuck" {
				t.Errorf("Check did not act as expected, got first tripper: %s, want first tripper: %s", trippers[0], "fuck")
			}
		case 2:
			if len(trippers) != 0 {
				t.Errorf("Check did not act as expected, got trippers length: %d, want trippers length: %d", len(trippers), 0)
			}
		default:
			t.Errorf("Check test invalid, got test messages length: %d, want test messages length: %d", len(messages), 4)
		}
	}
}

func TestCheck2(t *testing.T) {
	filter := NewSwearFilter(true, "fuck")
	messages := []string{"FuCking", "fûçk", "asdf", "what the f u c k dude"}

	for i := 0; i < len(messages); i++ {
		trippers, err := filter.Check(messages[i])
		require.NoError(t, err)
		t.Log(trippers)
	}
}

func TestCheckSimpleRegex(t *testing.T) {
	tests := map[string]bool{
		"fuckx":   true,
		"xfuck":   false,
		"xthing":  true,
		"thingx":  true,
		"xthingx": true,
		"xthat":   true,
		"thatx":   false,
		"ok":      false,
	}

	for _, spacedBypass := range []bool{true, false} {
		filter := NewSwearFilter(true, "^fuck", "thing", "that$")
		filter.EnableSpacedBypass = spacedBypass

		for test, expected := range tests {
			trippers, err := filter.Check(test)
			t.Logf("test=%s, expected=%v, trippers=%v\n", test, expected, trippers)
			require.NoError(t, err)
			require.Equal(t, expected, len(trippers) > 0)
		}
	}
}

func TestCheckSimpleRegexDisabled(t *testing.T) {
	tests := map[string]bool{
		"fuckx":   false,
		"xfuck":   false,
		"xthing":  true,
		"thingx":  true,
		"xthingx": true,
		"xthat":   false,
		"thatx":   false,
		"ok":      false,
	}

	for _, spacedBypass := range []bool{true, false} {
		filter := NewSwearFilter(true, "^fuck", "thing", "that$")
		filter.DisableSimpleRegex = true
		filter.EnableSpacedBypass = spacedBypass

		for test, expected := range tests {
			trippers, err := filter.Check(test)
			t.Logf("test=%s, expected=%v, trippers=%v\n", test, expected, trippers)
			require.NoError(t, err)
			require.Equal(t, expected, len(trippers) > 0)
		}
	}
}
