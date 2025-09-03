package grep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"wb-tech-l2/12/go-grep/internal/config"
)

type Service struct {
	cfg *config.Grep
}

func NewService(cfg *config.Grep) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Process(r io.Reader, w io.Writer) error {
	lines, err := readLines(r)
	if err != nil {
		return fmt.Errorf("failed to read lines: %w", err)
	}

	matcher, err := s.buildMatcher()
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	matched, count := s.matchLines(lines, matcher)

	if s.cfg.CountOnly {
		return s.printCount(w, count)
	}

	toPrint := s.applyContext(matched)
	return s.printLines(w, lines, matched, toPrint)
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (s *Service) buildMatcher() (func(string) bool, error) {
	if s.cfg.FixedString {
		if s.cfg.IgnoreCase {
			pattern := strings.ToLower(s.cfg.Pattern)
			return func(line string) bool {
				return strings.Contains(strings.ToLower(line), pattern)
			}, nil
		}
		return func(line string) bool {
			return strings.Contains(line, s.cfg.Pattern)
		}, nil
	}

	pattern := s.cfg.Pattern
	if s.cfg.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return func(line string) bool {
		return re.MatchString(line)
	}, nil
}

func (s *Service) matchLines(lines []string, matcher func(string) bool) ([]bool, int) {
	matched := make([]bool, len(lines))
	count := 0

	for i, line := range lines {
		matches := matcher(line)
		if s.cfg.InvertMatch {
			matches = !matches
		}
		if matches {
			matched[i] = true
			count++
		}
	}

	return matched, count
}

func (s *Service) applyContext(matched []bool) []bool {
	n := len(matched)
	toPrint := make([]bool, n)

	after := s.cfg.AfterContext
	before := s.cfg.BeforeContext

	if s.cfg.Context > 0 {
		after = s.cfg.Context
		before = s.cfg.Context
	}

	for i := 0; i < n; i++ {
		if !matched[i] {
			continue
		}

		toPrint[i] = true

		for k := i - before; k < i; k++ {
			if k >= 0 {
				toPrint[k] = true
			}
		}

		for k := i + 1; k <= i+after && k < n; k++ {
			toPrint[k] = true
		}
	}

	return toPrint
}

func (s *Service) printCount(w io.Writer, count int) error {
	_, err := fmt.Fprintf(w, "%d\n", count)
	return err
}

func (s *Service) printLines(w io.Writer, lines []string, matched, toPrint []bool) error {
	for i := 0; i < len(lines); i++ {
		if !toPrint[i] {
			continue
		}

		var output string
		if s.cfg.LineNumber {
			prefix := fmt.Sprintf("%d:", i+1)
			if !matched[i] {
				prefix = fmt.Sprintf("%d-", i+1)
			}
			output = prefix + lines[i] + "\n"
		} else {
			output = lines[i] + "\n"
		}

		if _, err := w.Write([]byte(output)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ProcessLines(lines []string) []string {
	buffer := strings.NewReader(strings.Join(lines, "\n"))
	output := &strings.Builder{}

	if err := s.Process(buffer, output); err != nil {
		return []string{fmt.Sprintf("Error: %v", err)}
	}

	return strings.Split(strings.TrimSpace(output.String()), "\n")
}

type options struct {
	config.Grep
}

type Svc func(reader io.Reader, writer io.Writer, cfg ...*options) error

func NewGrep() Svc {
	return func(reader io.Reader, writer io.Writer, cfg ...*options) error {
		return errors.New("not implemented")
	}
}
