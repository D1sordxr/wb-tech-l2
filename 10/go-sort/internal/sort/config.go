package sort

type Config struct {
	FileName             string
	Column               int
	Numeric              bool
	Reverse              bool
	Unique               bool
	Month                bool
	IgnoreTrailingBlanks bool
	CheckSorted          bool
	HumanNumeric         bool
}

func (c *Config) GetFileName() string  { return c.FileName }
func (c *Config) GetColumn() int       { return c.Column }
func (c *Config) IsNumeric() bool      { return c.Numeric }
func (c *Config) IsReverse() bool      { return c.Reverse }
func (c *Config) IsUnique() bool       { return c.Unique }
func (c *Config) IsMonth() bool        { return c.Month }
func (c *Config) IgnoreBlanks() bool   { return c.IgnoreTrailingBlanks }
func (c *Config) CheckIsSorted() bool  { return c.CheckSorted }
func (c *Config) IsHumanNumeric() bool { return c.HumanNumeric }
