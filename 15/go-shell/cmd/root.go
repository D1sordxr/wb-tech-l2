package cmd

import (
	"fmt"
	"os"
	"wb-tech-l2/15/go-shell/internal/handler"

	"github.com/spf13/cobra"
)

var (
	appName  = "go-shell"
	shortMsg = "Mini Unix shell implementation in Go"
	longMsg  = `
Go-shell is a mini Unix shell implementation that supports basic commands.

Built-in commands:
  echo    - display a line of text
  cd      - change the working directory
  pwd     - print working directory
  kill    - send a signal to a process
  ps      - report a snapshot of current processes

Examples:
  go-shell echo "Hello World"
  go-shell cd /tmp
  go-shell pwd
  go-shell kill 1234
  go-shell ps
`
)

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "shows app usage")
}

func setupCommandArgs(args []string) []string {
	commandArgs := make([]string, len(args)-1)
	copy(commandArgs, args[1:])
	return commandArgs
}

func runApp(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	if err := handler.HandleCommand(
		args[0],
		setupCommandArgs(args),
		os.Stdout,
	); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: shortMsg,
	Long:  longMsg,
	Run:   runApp,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
