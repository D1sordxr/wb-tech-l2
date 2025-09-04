package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"wb-tech-l2/13/go-cut/internal/config"
	"wb-tech-l2/13/go-cut/internal/cut"
)

var (
	appName  = `go-cut`
	shortMsg = `Cut utility for extracting sections from each line of files`
	longMsg  = `
Go-cut is a clone of the Unix cut utility.
It extracts parts from each line of input and outputs them.

Examples:
  echo "a:b:c:d" | go-cut -d ":" -f 1,3
  echo "a:b:c:d" | go-cut -d ":" -f 2-4
  echo "a:b:c:d" | go-cut -d ":" -f 1,3-4 -s
`
	appConfig = &config.Cut{}
	rootCmd   = &cobra.Command{
		Use:   appName,
		Short: shortMsg,
		Long:  longMsg,
		Run:   runApp,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "shows app usage")

	rootCmd.Flags().StringVarP(&appConfig.Fields, "fields", "f", "", "select only these fields")
	rootCmd.Flags().StringVarP(&appConfig.Delimiter, "delimiter", "d", "\t", "use DELIM instead of TAB for field delimiter")
	rootCmd.Flags().BoolVarP(&appConfig.SeparatedOnly, "separated", "s", false, "do not print lines not containing delimiters")
}

func exitWithErrorMessage(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	os.Exit(1)
}

func runApp(_ *cobra.Command, args []string) {
	var input io.Reader = os.Stdin
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			exitWithErrorMessage(fmt.Sprintf("cannot open file: %s", err))
		}
		defer func() { _ = file.Close() }()
		input = file
	}

	if err := cut.Process(input, os.Stdout, cut.Opts{
		Fields:        appConfig.Fields,
		Delimiter:     appConfig.Delimiter,
		SeparatedOnly: appConfig.SeparatedOnly,
	}); err != nil {
		exitWithErrorMessage(fmt.Sprintf("processing error: %s", err))
	}
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
