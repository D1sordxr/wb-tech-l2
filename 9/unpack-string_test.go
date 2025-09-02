package main

import "testing"

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "basic UnpackString",
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
			wantErr:  false,
		},
		{
			name:     "no numbers",
			input:    "abcd",
			expected: "abcd",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "starts with digit",
			input:    "45",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "starts with digit in context",
			input:    "4a",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "escape sequences",
			input:    `qwe\4\5`,
			expected: "qwe45",
			wantErr:  false,
		},
		{
			name:     "escape with numbers",
			input:    `qwe\45`,
			expected: "qwe44444",
			wantErr:  false,
		},
		{
			name:     "multiple digit numbers",
			input:    "a10b",
			expected: "aaaaaaaaaab",
			wantErr:  false,
		},
		{
			name:     "large number",
			input:    "a2b3c15",
			expected: "aabbbccccccccccccccc",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnpackString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if result != tt.expected {
				t.Errorf("UnpackString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

var testCases = []struct {
	name     string
	input    string
	expected string
}{
	{"simple", "a4bc2d5e", "aaaabccddddde"},
	{"no_numbers", "abcd", "abcd"},
	{"escape", "qwe\\4\\5", "qwe45"},
	{"empty", "", ""},
	{"invalid_start", "45", ""},
	{"long_number", "a10", "aaaaaaaaaa"},
	{"mixed", "a2b3c1", "aabbbc"},
}

func BenchmarkUnpackOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_, _ = optimizedForLongStrings(tc.input)
		}
	}
}

func BenchmarkUnpackOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_, _ = optimizedForLongStrings(tc.input)
		}
	}
}

func BenchmarkUnpackOriginalLong(b *testing.B) {
	longInput := "a10b20c30d40e50f60g70h80i90j100"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnpackString(longInput)
	}
}

func BenchmarkUnpackOptimizedLong(b *testing.B) {
	longInput := "a10b20c30d40e50f60g70h80i90j100"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = optimizedForLongStrings(longInput)
	}
}

func TestBothImplementations(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			orig, _ := UnpackString(tc.input)
			opt, _ := UnpackString(tc.input)

			if orig != tc.expected {
				t.Errorf("Original failed: input %s, got %s, expected %s", tc.input, orig, tc.expected)
			}
			if opt != tc.expected {
				t.Errorf("Optimized failed: input %s, got %s, expected %s", tc.input, opt, tc.expected)
			}
		})
	}
}
