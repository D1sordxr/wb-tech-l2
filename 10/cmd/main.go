package main

import (
	"fmt"
	"wb-tech-l2/10/internal/sort"
)

func main() {
	cfg := new(sort.Config)
	cfg.Setup()

	fmt.Println(cfg)
}
