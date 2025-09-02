package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Service struct {
	filePath string
}

func NewService(filePath string) *Service {
	return &Service{filePath: filePath}
}

func (s *Service) ReadLines() ([]string, error) {
	const op = "reader.Service.ReadLines"

	var reader io.ReadCloser
	if s.filePath == "" {
		reader = os.Stdin
		fmt.Println("Reading text from STDIN. Enter text (press Ctrl+D to finish):")
	} else {
		file, err := os.Open(s.filePath)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer func() { _ = file.Close() }()
		reader = file
	}

	lines := make([]string, 0, 32)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}
