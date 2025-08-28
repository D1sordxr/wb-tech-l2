package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// To use this application you need to write string to unpack.
// Example: go run ./9/main.go --data=qwe\45

func main() {
	var data string
	flag.StringVar(&data, "data", "", "string for unpack")
	flag.Parse()

	result, err := UnpackString(data)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(result)
}

var (
	ErrEmptyStr   = errors.New("empty string")
	ErrInvalidStr = errors.New("invalid string")
)

func UnpackString(s string) (string, error) {
	if s == "" {
		return "", ErrEmptyStr
	}
	if len(s) > 0 && unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidStr
	}

	var (
		result     strings.Builder
		prevChar   rune
		escapeMode bool
	)

	for i := 0; i < len(s); i++ {
		char := rune(s[i])

		if char == '\\' && !escapeMode {
			escapeMode = true
			continue
		}

		if escapeMode {
			result.WriteRune(prevChar)
			prevChar = char
			escapeMode = false
			continue
		}

		if unicode.IsDigit(char) {
			if prevChar == 0 {
				return "", ErrInvalidStr
			}

			start := i
			for i < len(s) && unicode.IsDigit(rune(s[i])) {
				i++
			}
			numStr := s[start:i]
			i--

			count, err := strconv.Atoi(numStr)
			if err != nil || count <= 0 {
				return "", fmt.Errorf("invalid number: %w", err)
			}

			result.WriteString(strings.Repeat(string(prevChar), count))
			prevChar = 0
			continue
		}

		if prevChar != 0 {
			result.WriteRune(prevChar)
		}
		prevChar = char
	}

	if escapeMode {
		if prevChar != 0 {
			result.WriteRune(prevChar)
		}
		result.WriteRune('\\')
	} else if prevChar != 0 {
		result.WriteRune(prevChar)
	}

	return result.String(), nil
}
