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
			name:     "basic unpack",
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
