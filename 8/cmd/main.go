package main

import (
	"fmt"
	"os"

	"wb-tech-l2/8/internal/ntp"
)

func main() {
	timeService := new(ntp.Service)
	time, err := timeService.GetTime()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(time)
	os.Exit(0)
}
