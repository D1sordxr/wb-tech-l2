package cut

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Opts struct {
	Fields        string
	Delimiter     string
	SeparatedOnly bool
}

type validOpts struct {
	fields        map[int]struct{}
	delimiter     string
	separatedOnly bool
}

func Process(r io.Reader, w io.Writer, rawOpts ...Opts) error {
	opts := validOpts{
		fields:    make(map[int]struct{}),
		delimiter: "\t",
	}

	if len(rawOpts) > 0 {
		var err error
		opts, err = parseOpts(rawOpts[0])
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		containsDelimiter := strings.Contains(line, opts.delimiter)

		if opts.separatedOnly && !containsDelimiter {
			continue
		}

		if !containsDelimiter {
			_, _ = fmt.Fprintln(w, line)
			continue
		}

		fields := strings.Split(line, opts.delimiter)

		if len(opts.fields) == 0 {
			_, _ = fmt.Fprintln(w, line)
			continue
		}

		var outputFields []string
		for i := 0; i < len(fields); i++ {
			fieldIndex := i + 1
			if _, ok := opts.fields[fieldIndex]; ok {
				outputFields = append(outputFields, fields[i])
			}
		}

		if len(outputFields) > 0 {
			_, _ = fmt.Fprintln(w, strings.Join(outputFields, opts.delimiter))
		}
	}

	return scanner.Err()
}

func parseOpts(raw Opts) (validOpts, error) {
	result := validOpts{
		delimiter:     "\t",
		separatedOnly: raw.SeparatedOnly,
	}

	if raw.Delimiter != "" {
		result.delimiter = raw.Delimiter
	}

	if raw.Fields != "" {
		fields, err := parseFields(raw.Fields)
		if err != nil {
			return validOpts{}, err
		}
		result.fields = fields
	}

	return result, nil
}

func parseFields(fieldsStr string) (map[int]struct{}, error) {
	result := make(map[int]struct{})

	if fieldsStr == "" {
		return result, nil
	}

	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start of range: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end of range: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("range start cannot be greater than end: %s", part)
			}

			for i := start; i <= end; i++ {
				if i > 0 {
					result[i] = struct{}{}
				}
			}
		} else {
			fieldNum, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			if fieldNum > 0 {
				result[fieldNum] = struct{}{}
			}
		}
	}

	return result, nil
}
