package sort

import "flag"

type Config struct {
	column    int
	isNumeric bool
	isReverse bool
	isUnique  bool
}

func (c *Config) Setup() {
	flag.IntVar(&c.column, "k", 0, "column number")
	flag.BoolVar(&c.isNumeric, "n", c.isNumeric, "isNumeric")
	flag.BoolVar(&c.isReverse, "r", c.isReverse, "isReverse")
	flag.BoolVar(&c.isUnique, "u", c.isUnique, "isUnique")

	flag.Parse()
}
