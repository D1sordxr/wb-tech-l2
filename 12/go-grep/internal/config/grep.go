package config

type Grep struct {
	Pattern string

	// Reader settings
	FilePath string // -f

	// Context settings
	AfterContext  int // -A
	BeforeContext int // -B
	Context       int // -C

	// Bool flags
	CountOnly   bool // -c
	IgnoreCase  bool // -i
	InvertMatch bool // -v
	FixedString bool // -F
	LineNumber  bool // -n
}
