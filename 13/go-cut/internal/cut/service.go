package cut

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"wb-tech-l2/13/go-cut/internal/config"
)

type Processor struct {
	fields        map[int]struct{}
	delimiter     string
	separatedOnly bool
}

func NewProcessor(cfg *config.Cut) (*Processor, error) {
	fieldSet, err := parseFields(cfg.Fields)
	if err != nil {
		return nil, err
	}

	delimiter := cfg.Delimiter
	if delimiter == "" {
		delimiter = "\t"
	}

	return &Processor{
		fields:        fieldSet,
		delimiter:     delimiter,
		separatedOnly: cfg.SeparatedOnly,
	}, nil
}

func (p *Processor) Process(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		containsDelimiter := strings.Contains(line, p.delimiter)

		if p.separatedOnly && !containsDelimiter {
			continue
		}

		if !containsDelimiter {
			_, _ = fmt.Fprintln(output, line)
			continue
		}

		fields := strings.Split(line, p.delimiter)

		var outputFields []string
		for i := 0; i < len(fields); i++ {
			fieldIndex := i + 1
			if _, ok := p.fields[fieldIndex]; ok {
				outputFields = append(outputFields, fields[i])
			}
		}

		if len(outputFields) > 0 {
			_, _ = fmt.Fprintln(output, strings.Join(outputFields, p.delimiter))
		}
	}

	return scanner.Err()
}
