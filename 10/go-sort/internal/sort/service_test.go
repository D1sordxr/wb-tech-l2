package sort

import (
	"testing"
)

func TestService_extractField(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		column   int
		expected string
	}{
		{
			name:     "negative column returns whole line",
			line:     "hello\tworld",
			column:   -1,
			expected: "hello\tworld",
		},
		{
			name:     "column 0 returns first field",
			line:     "apple\tbanana\tcherry",
			column:   0,
			expected: "apple",
		},
		{
			name:     "column 1 returns second field",
			line:     "apple\tbanana\tcherry",
			column:   1,
			expected: "banana",
		},
		{
			name:     "column out of bounds returns empty string",
			line:     "apple\tbanana",
			column:   5,
			expected: "",
		},
		{
			name:     "empty line with any column",
			line:     "",
			column:   0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				config: &Config{Column: tt.column},
				parser: new(parser),
			}

			result := svc.extractField(tt.line)
			if result != tt.expected {
				t.Errorf("extractField(%q, %d) = %q, expected %q", tt.line, tt.column, result, tt.expected)
			}
		})
	}
}

func TestService_comparePrepared(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		a        sortableLine
		b        sortableLine
		expected bool
	}{
		{
			name:     "string comparison",
			config:   &Config{},
			a:        sortableLine{stringValue: "apple"},
			b:        sortableLine{stringValue: "banana"},
			expected: true,
		},
		{
			name:     "numeric comparison",
			config:   &Config{Numeric: true},
			a:        sortableLine{numberValue: 10},
			b:        sortableLine{numberValue: 20},
			expected: true,
		},
		{
			name:     "month comparison",
			config:   &Config{Month: true},
			a:        sortableLine{monthIndex: 1},
			b:        sortableLine{monthIndex: 2},
			expected: true,
		},
		{
			name:     "fallback to string when numeric parsing fails for one",
			config:   &Config{Numeric: true},
			a:        sortableLine{stringValue: "apple", numberValue: 10},
			b:        sortableLine{stringValue: "banana", numberValue: 0},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{config: tt.config}
			result := svc.comparePrepared(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("comparePrepared() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestService_uniquePrepared(t *testing.T) {
	tests := []struct {
		name     string
		input    []sortableLine
		expected []sortableLine
	}{
		{
			name:     "empty slice",
			input:    []sortableLine{},
			expected: []sortableLine{},
		},
		{
			name: "no duplicates",
			input: []sortableLine{
				{original: "apple"},
				{original: "banana"},
				{original: "cherry"},
			},
			expected: []sortableLine{
				{original: "apple"},
				{original: "banana"},
				{original: "cherry"},
			},
		},
		{
			name: "with duplicates",
			input: []sortableLine{
				{original: "apple"},
				{original: "apple"},
				{original: "banana"},
				{original: "banana"},
				{original: "cherry"},
			},
			expected: []sortableLine{
				{original: "apple"},
				{original: "banana"},
				{original: "cherry"},
			},
		},
		{
			name: "consecutive duplicates only",
			input: []sortableLine{
				{original: "apple"},
				{original: "banana"},
				{original: "apple"},
				{original: "cherry"},
			},
			expected: []sortableLine{
				{original: "apple"},
				{original: "banana"},
				{original: "apple"},
				{original: "cherry"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			result := svc.uniquePrepared(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("uniquePrepared() length = %d, expected %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i].original != tt.expected[i].original {
					t.Errorf("uniquePrepared()[%d] = %q, expected %q", i, result[i].original, tt.expected[i].original)
				}
			}
		})
	}
}

func TestService_IsSorted(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		lines    []string
		expected bool
	}{
		{
			name:     "empty lines are sorted",
			config:   &Config{},
			lines:    []string{},
			expected: true,
		},
		{
			name:     "single line is sorted",
			config:   &Config{},
			lines:    []string{"apple"},
			expected: true,
		},
		{
			name:     "sorted strings",
			config:   &Config{},
			lines:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "unsorted strings",
			config:   &Config{},
			lines:    []string{"banana", "apple", "cherry"},
			expected: false,
		},
		{
			name:     "sorted numeric by first column",
			config:   &Config{Numeric: true, Column: 0},
			lines:    []string{"10\tapple", "20\tbanana", "30\tcherry"},
			expected: true,
		},
		{
			name:     "reverse sorted with reverse flag",
			config:   &Config{Reverse: true},
			lines:    []string{"cherry", "banana", "apple"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				config: tt.config,
				lines:  tt.lines,
				parser: new(parser),
			}

			result := svc.IsSorted()
			if result != tt.expected {
				t.Errorf("IsSorted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestService_Sort(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		input    []string
		expected []string
	}{
		{
			name:     "basic string sort",
			config:   &Config{},
			input:    []string{"banana", "apple", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "numeric sort",
			config:   &Config{Numeric: true},
			input:    []string{"100", "20", "3"},
			expected: []string{"3", "20", "100"},
		},
		{
			name:     "unique filter",
			config:   &Config{Unique: true},
			input:    []string{"apple", "apple", "banana", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "sort by column",
			config:   &Config{Column: 1},
			input:    []string{"1\tbanana", "2\tapple", "3\tcherry"},
			expected: []string{"2\tapple", "1\tbanana", "3\tcherry"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				config: tt.config,
				lines:  tt.input,
				parser: new(parser),
			}

			result, err := svc.Sort()
			if err != nil {
				t.Errorf("Sort() unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Sort() length = %d, expected %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Sort()[%d] = %q, expected %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestService_Sort_CheckIfSorted(t *testing.T) {
	svc := &Service{
		config: &Config{CheckSorted: true},
		lines:  []string{"apple", "banana", "cherry"},
		parser: new(parser),
	}

	result, err := svc.Sort()
	if err != nil {
		t.Errorf("Sort() with sorted input should return nil error, got: %v", err)
	}
	if result != nil {
		t.Errorf("Sort() with CheckIsSorted should return nil result, got: %v", result)
	}

	svc.lines = []string{"banana", "apple", "cherry"}
	_, err = svc.Sort()
	if err == nil {
		t.Error("Sort() with unsorted input should return error")
	}
}
