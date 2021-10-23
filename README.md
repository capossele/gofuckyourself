# gofuckyourself

[![GoDoc](https://godoc.org/github.com/JoshuaDoes/gofuckyourself?status.svg)](https://godoc.org/github.com/JoshuaDoes/gofuckyourself)
[![Go Report Card](https://goreportcard.com/badge/github.com/JoshuaDoes/gofuckyourself)](https://goreportcard.com/report/github.com/JoshuaDoes/gofuckyourself)
[![cover.run](https://cover.run/go/github.com/JoshuaDoes/gofuckyourself.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2FJoshuaDoes%2Fgofuckyourself)

A sanitization-based swear filter for Go.

# Installing
Just `import github.com/capossele/swearfilter` if using go modules.

# Example
```Go
package main

import (
	"fmt"

	"github.com/capossele/swearfilter"
)

var message = "This is a fûçking message with shitty swear words asswipe."
var swears = []string{"fuck", "shit", "^ass"}

func main() {
	filter := swearfilter.New(false, false, false, false, false, swears...)
	swearFound, swearsFound, err := filter.Check(message)
	fmt.Println("Swear found: ", swearFound)
	fmt.Println("Swears tripped: ", swearsFound)
	fmt.Println("Error: ", err)
}
```
### Output
```
> go run main.go
Swear found:  true
Swears tripped:  [fuck shit ^ass]
Error:  <nil>
```

## Options
By default substring testing is performed, e.g. so `abc` will match any of `1abc`, `1abc2` and `abc2`.

To help keep word lists concise but performance good, simple (simulated) regex matching is supported. The only control characters supported are `^` and `$`. These will perform prefix/suffix string match tests respectively. E.g. so `^ass` will match `asses` but not `pass`. These simple regexes aren't compiled to regexes internally so may be faster (but they haven't been benchmarked).

Full regex support can be enabled by passing the relevant parameter when calling `swearfilter.New`. In this case, each swear word will be compiled to a regex and tested with regex matching.

## License
The source code for gofuckyourself is released under the MIT License. See LICENSE for more details.

## Donations
All donations are appreciated and help me stay awake at night to work on this more. Even if it's not much, it helps a lot in the long run!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)