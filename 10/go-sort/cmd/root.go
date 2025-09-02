package cmd

import (
	"fmt"
	"os"

	"wb-tech-l2/10/go-sort/internal/sort"

	"github.com/spf13/cobra"
)

var (
	appConfig = new(sort.Config)
	appName   = `go-sort`
	shortMsg  = `Sort utility for text files and input`
	longMsg   = `
A sorting utility that supports various sorting options similar to Unix sort command. 
Supports sorting by columns, numeric sorting, reverse order, and more.
`
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: shortMsg,
	Long:  longMsg,
	Run:   runApp,
}

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, "shows app usage")

	rootCmd.Flags().StringVarP(&appConfig.FileName, "file", "f", "", "read from file")
	rootCmd.Flags().IntVarP(&appConfig.Column, "key", "k", 1, "sort by column number")
	rootCmd.Flags().BoolVarP(&appConfig.Numeric, "numeric", "n", false, "sort numerically")
	rootCmd.Flags().BoolVarP(&appConfig.Reverse, "reverse", "r", false, "reverse sort order")
	rootCmd.Flags().BoolVarP(&appConfig.Unique, "unique", "u", false, "output only unique lines")
	rootCmd.Flags().BoolVarP(&appConfig.Month, "month", "M", false, "sort by month names")
	rootCmd.Flags().BoolVarP(&appConfig.IgnoreTrailingBlanks, "ignore-blanks", "b", false, "ignore trailing blanks")
	rootCmd.Flags().BoolVarP(&appConfig.CheckSorted, "check", "c", false, "check if data is sorted")
	rootCmd.Flags().BoolVarP(&appConfig.HumanNumeric, "human-numeric", "h", false, "sort by human-readable numbers")

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "help about any command",
		Long:  longMsg,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Usage()
		},
	})
}

func runApp(_ *cobra.Command, _ []string) {
	sortSvc := sort.NewService(appConfig)
	sortSvc.MustReadLines()
	sortSvc.MustWriteLines()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
