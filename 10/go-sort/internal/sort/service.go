package sort

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Service struct {
	config *Config
	reader io.ReadCloser
	parser *parser
	lines  []string
}

func NewService(config *Config) *Service {
	svc := &Service{
		config: config,
		lines:  make([]string, 0, 32),
		parser: new(parser),
	}

	if svc.config.GetFileName() == "" {
		svc.reader = os.Stdin
		fmt.Println("Reading text from STDIN. Enter text (press Ctrl+D to finish):")
	} else {
		file, err := os.Open(svc.config.GetFileName())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			os.Exit(1)
		}
		svc.reader = file
	}

	if svc.config.Column > 0 { // column idx
		svc.config.Column--
	}

	return svc
}

func (s *Service) MustReadLines() {
	defer func() { _ = s.reader.Close() }()

	scanner := bufio.NewScanner(s.reader)

	for scanner.Scan() {
		line := scanner.Text()

		if s.config.IgnoreBlanks() {
			line = strings.TrimRight(line, " \t")
		}

		s.lines = append(s.lines, line)
	}

	if s.config.CheckIsSorted() {
		if s.IsSorted() {
			fmt.Println("Data is sorted")
			os.Exit(0)
		} else {
			fmt.Println("Data is not sorted")
			os.Exit(1)
		}
	}
}

func (s *Service) MustWriteLines() {
	sortedLines, err := s.Sort()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	for _, line := range sortedLines {
		fmt.Println(line)
	}
}

type sortableLine struct {
	original    string
	stringValue string
	numberValue float64
	monthIndex  int
}

func (s *Service) Sort() ([]string, error) {
	prepared := make([]sortableLine, len(s.lines))
	for i, line := range s.lines {
		key := s.extractField(line)
		sl := sortableLine{
			original:    line,
			stringValue: key,
		}

		if s.config.IsNumeric() {
			if num, err := s.parser.ParseFloat(key); err == nil {
				sl.numberValue = num
			}
		}

		if s.config.IsHumanNumeric() {
			if num, err := s.parser.ParseHumanNumber(key); err == nil {
				sl.numberValue = num
			}
		}

		if s.config.IsMonth() {
			if month, err := s.parser.ParseMonth(key); err == nil {
				sl.monthIndex = month
			}
		}

		prepared[i] = sl
	}

	less := func(i, j int) bool {
		a, b := prepared[i], prepared[j]
		result := s.comparePrepared(a, b)

		if s.config.IsReverse() {
			return !result
		}
		return result
	}

	if s.config.CheckIsSorted() {
		if sort.SliceIsSorted(prepared, less) {
			return nil, nil
		}
		return nil, fmt.Errorf("input is not sorted")
	}

	sort.Slice(prepared, less)

	if s.config.IsUnique() {
		prepared = s.uniquePrepared(prepared)
	}

	result := make([]string, len(prepared))
	for i, sl := range prepared {
		result[i] = sl.original
	}

	return result, nil
}

func (s *Service) comparePrepared(a, b sortableLine) bool {
	if s.config.IsMonth() && a.monthIndex != 0 && b.monthIndex != 0 {
		return a.monthIndex < b.monthIndex
	}

	if (s.config.IsNumeric() || s.config.IsHumanNumeric()) && a.numberValue != 0 && b.numberValue != 0 {
		return a.numberValue < b.numberValue
	}

	return a.stringValue < b.stringValue
}

func (s *Service) IsSorted() bool {
	prepared := make([]sortableLine, len(s.lines))
	for i, line := range s.lines {
		key := s.extractField(line)
		sl := sortableLine{
			original:    line,
			stringValue: key,
		}
		prepared[i] = sl
	}

	less := func(i, j int) bool {
		a, b := prepared[i], prepared[j]
		result := s.comparePrepared(a, b)

		if s.config.IsReverse() {
			return !result
		}
		return result
	}

	return sort.SliceIsSorted(prepared, less)
}

func (s *Service) extractField(line string) string {
	if s.config.GetColumn() < 0 {
		return line
	}

	fields := strings.Split(line, "\t")
	if s.config.GetColumn() < len(fields) {
		return fields[s.config.GetColumn()]
	}

	return ""
}

func (s *Service) uniquePrepared(prepared []sortableLine) []sortableLine {
	if len(prepared) == 0 {
		return prepared
	}

	unique := []sortableLine{prepared[0]}
	for i := 1; i < len(prepared); i++ {
		if prepared[i].original != prepared[i-1].original {
			unique = append(unique, prepared[i])
		}
	}
	return unique
}
