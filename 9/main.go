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

// To use this application you need to write string to UnpackString.
// Example: go run ./9/main.go --data=qwe\45

func main() {
	var data string
	flag.StringVar(&data, "data", "", "string for UnpackString")
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
	switch {
	case s == "":
		return "", ErrEmptyStr
	case len(s) > 0 && unicode.IsDigit(rune(s[0])):
		return "", ErrInvalidStr
	}

	type entry struct {
		char  rune
		count int
	}

	defaultEntry := func(char rune) entry {
		return entry{char, 1}
	}

	safeIsUnpackChar := func(idx int) int { // todo: nums of indexes to skip
		var num string
		if idx+1 < len(s) && unicode.IsDigit(rune(s[idx+1])) {
			for j := idx + 1; j < len(s); j++ {
				if unicode.IsDigit(rune(s[j])) {
					num += string(rune(s[j]))
				} else {
					break
				}
			}
		}
		if num == "" {
			num = "1"
		}
		res, _ := strconv.Atoi(num)
		return res
	}

	var escape bool

	entries := make([]entry, 0, len(s))
	for i := 0; i < len(s); i++ {
		if escape {
			count := safeIsUnpackChar(i)
			entries = append(entries, entry{rune(s[i]), count})
			escape = false
			continue
		}

		if s[i] == '\\' {
			escape = true
			continue
		}

		if unicode.IsDigit(rune(s[i])) {
			continue
		}

		if i < len(s)-1 {
			if unicode.IsDigit(rune(s[i+1])) {
				count := safeIsUnpackChar(i)
				entries = append(entries, entry{rune(s[i]), count})
			} else {
				entries = append(entries, defaultEntry(rune(s[i])))
			}
		} else {
			entries = append(entries, defaultEntry(rune(s[i])))
		}
	}

	// var result string
	// for _, e := range entries { result += strings.Repeat(string(e.char), e.count) }

	totalLength := 0
	for _, e := range entries {
		totalLength += e.count
	}

	var res strings.Builder
	res.Grow(totalLength)
	for _, e := range entries {
		for range e.count {
			res.WriteRune(e.char)
		}
	}

	return res.String(), nil
}

func optimizedForLongStrings(s string) (string, error) {
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
