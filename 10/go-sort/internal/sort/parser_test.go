package sort

import "testing"

func TestParseFloat(t *testing.T) {
	p := new(parser)

	tests := []struct {
		v     string
		exp   float64
		isErr bool
	}{
		{"absolute cinema", 0, true},
		{"1337", 1337, false},
		{"  42.52	 ", 42.52, false},
	}

	for _, tt := range tests {
		got, err := p.ParseFloat(tt.v)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseFloat(%q) error = %v, wantErr %v", tt.v, err, tt.isErr)
		}
		if !tt.isErr && got != tt.exp {
			t.Errorf("ParseFloat(%q) = %v, exp %v", tt.v, got, tt.exp)
		}
	}
}

func TestParseMonth(t *testing.T) {
	p := new(parser)

	tests := []struct {
		v     string
		exp   int
		isErr bool
	}{
		{"Jan", 1, false},
		{"auG", 8, false},
		{"  Feb			  ", 2, false},
		{"month", 0, true},
	}

	for _, tt := range tests {
		got, err := p.ParseMonth(tt.v)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseMonth(%q) error = %v, wantErr %v", tt.v, err, tt.isErr)
		}
		if !tt.isErr && got != tt.exp {
			t.Errorf("ParseMonth(%q) = %v, exp %v", tt.v, got, tt.exp)
		}
	}
}

func TestParseHumanNumber(t *testing.T) {
	p := new(parser)

	tests := []struct {
		v     string
		exp   float64
		isErr bool
	}{
		{"777", 777, false},
		{"1K", 1024, false},
		{"2M", 2 * 1024 * 1024, false},
		{"100G", 100 * 1024 * 1024 * 1024, false},
		{"error", 0, true},
	}

	for _, tt := range tests {
		got, err := p.ParseHumanNumber(tt.v)
		if (err != nil) != tt.isErr {
			t.Errorf("ParseHumanNumber(%q) error = %v, wantErr %v", tt.v, err, tt.isErr)
		}
		if !tt.isErr && got != tt.exp {
			t.Errorf("ParseHumanNumber(%q) = %v, exp %v", tt.v, got, tt.exp)
		}
	}
}
