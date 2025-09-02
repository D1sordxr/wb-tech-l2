package cmd

import (
	"fmt"
	"os"
	"wb-tech-l2/12/go-grep/internal/config"
	"wb-tech-l2/12/go-grep/internal/reader"

	"github.com/spf13/cobra"
)

var (
	appName  = `go-grep`
	shortMsg = `Grep utility for text files and input`
	longMsg  = `
A pattern search utility that supports various grep options similar to Unix grep command.
Supports regular expressions, context lines, ignore case, inverse match and more.
`
	appConfig = new(config.Grep)
	rootCmd   = &cobra.Command{
		Use:   appName,
		Short: shortMsg,
		Long:  longMsg,
		Run:   runApp,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("help", "", false, "shows app usage")

	// Optional flags for file path and pattern
	rootCmd.Flags().StringVarP(&appConfig.Pattern, "pattern", "e", "", "pattern to search for (required)")
	rootCmd.Flags().StringVarP(&appConfig.FilePath, "file", "f", "", "read from file (default: stdin)")

	// Context flags
	rootCmd.Flags().IntVarP(&appConfig.AfterContext, "after-context", "A", 0, "print N lines after match")
	rootCmd.Flags().IntVarP(&appConfig.BeforeContext, "before-context", "B", 0, "print N lines before match")
	rootCmd.Flags().IntVarP(&appConfig.Context, "context", "C", 0, "print N lines around match")

	// Bool flags
	rootCmd.Flags().BoolVarP(&appConfig.CountOnly, "count", "c", false, "print only count of matching lines")
	rootCmd.Flags().BoolVarP(&appConfig.IgnoreCase, "ignore-case", "i", false, "ignore case distinctions")
	rootCmd.Flags().BoolVarP(&appConfig.InvertMatch, "invert-match", "v", false, "select non-matching lines")
	rootCmd.Flags().BoolVarP(&appConfig.FixedString, "fixed-strings", "F", false, "interpret pattern as fixed string")
	rootCmd.Flags().BoolVarP(&appConfig.LineNumber, "line-number", "n", false, "print line number with output")

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Help about any command",
		Run: func(cmd *cobra.Command, _ []string) {
			_ = cmd.Usage()
		},
	})
}

func exitWithErrorMessage(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", message)
	os.Exit(1)
}

func mustSetupPattern(args []string) {
	if appConfig.Pattern == "" && len(args) > 0 {
		appConfig.Pattern = args[0]
	}

	if appConfig.Pattern == "" {
		exitWithErrorMessage("pattern is required")
	}
}

func runApp(_ *cobra.Command, args []string) {
	mustSetupPattern(args)

	_, err := reader.NewService(appConfig.FilePath).ReadLines()
	if err != nil {
		exitWithErrorMessage(err.Error())
	}

	// TODO main service
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
