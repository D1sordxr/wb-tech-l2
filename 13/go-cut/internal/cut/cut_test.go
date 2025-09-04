package cut

import (
	"bytes"
	"strings"
	"testing"
)

func TestCutProcessor(t *testing.T) {
	tests := []struct {
		name           string
		fields         string
		delimiter      string
		separatedOnly  bool
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "basic field selection",
			fields:         "1,3",
			delimiter:      ":",
			input:          "a:b:c:d\ne:f:g:h\n",
			expectedOutput: "a:c\ne:g\n",
		},
		{
			name:           "range selection",
			fields:         "2-4",
			delimiter:      ":",
			input:          "a:b:c:d:e\n",
			expectedOutput: "b:c:d\n",
		},
		{
			name:           "separated only flag",
			fields:         "1",
			delimiter:      ":",
			separatedOnly:  true,
			input:          "a:b:c\nno_delimiter\nx:y:z\n",
			expectedOutput: "a\nx\n",
		},
		{
			name:           "without separated only flag",
			fields:         "1",
			delimiter:      ":",
			separatedOnly:  false,
			input:          "a:b:c\nno_delimiter\nx:y:z\n",
			expectedOutput: "a\nno_delimiter\nx\n",
		},
		{
			name:           "tab delimiter by default",
			fields:         "2",
			delimiter:      "",
			input:          "a\tb\tc\nd\te\tf\n",
			expectedOutput: "b\ne\n",
		},
		{
			name:           "no fields specified",
			fields:         "",
			delimiter:      ":",
			input:          "a:b:c\n",
			expectedOutput: "a:b:c\n",
		},
		{
			name:        "invalid fields format",
			fields:      "1-2-3",
			expectError: true,
		},
		{
			name:        "invalid field number",
			fields:      "abc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			input := strings.NewReader(tt.input)

			opts := Opts{
				Fields:        tt.fields,
				Delimiter:     tt.delimiter,
				SeparatedOnly: tt.separatedOnly,
			}

			err := Process(input, &output, opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if output.String() != tt.expectedOutput {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expectedOutput, output.String())
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		fields         string
		delimiter      string
		input          string
		expectedOutput string
	}{
		{
			name:           "empty input",
			fields:         "1",
			delimiter:      ":",
			input:          "",
			expectedOutput: "",
		},
		{
			name:           "single field with delimiter",
			fields:         "1",
			delimiter:      ":",
			input:          "hello\n",
			expectedOutput: "hello\n",
		},
		{
			name:           "multiple ranges",
			fields:         "1,3-4",
			delimiter:      ":",
			input:          "a:b:c:d:e\n",
			expectedOutput: "a:c:d\n",
		},
		{
			name:           "overlapping ranges",
			fields:         "1-3,2-4",
			delimiter:      ":",
			input:          "a:b:c:d:e\n",
			expectedOutput: "a:b:c:d\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			input := strings.NewReader(tt.input)

			opts := Opts{
				Fields:    tt.fields,
				Delimiter: tt.delimiter,
			}

			err := Process(input, &output, opts)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if output.String() != tt.expectedOutput {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expectedOutput, output.String())
			}
		})
	}
}

func TestParseFieldsErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		fields    string
		expectErr bool
	}{
		{"invalid range format", "1-2-3", true},
		{"invalid start range", "a-3", true},
		{"invalid end range", "1-b", true},
		{"reverse range", "3-1", true},
		{"negative field", "-1", true},
		{"negative range start", "-1-3", true},
		{"negative range end", "1--3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseFields(tt.fields)
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for %q, but got none", tt.fields)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error for %q: %v", tt.fields, err)
			}
		})
	}
}
