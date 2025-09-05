package handler

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var validCommands = map[string]func(args []string, w io.Writer) error{
	"echo": echo,
	"cd":   cd,
	"pwd":  pwd,
	"kill": kill,
	"ps":   ps,
}

func HandleCommand(command string, args []string, w io.Writer) error {
	if c, ok := validCommands[command]; ok {
		return c(args, w)
	}
	return fmt.Errorf("command not found: %s", command)
}

func echo(args []string, w io.Writer) error {
	_, err := w.Write([]byte(strings.Join(args, " ") + "\n"))
	return err
}

func cd(args []string, _ io.Writer) error {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(home)
	}
	return os.Chdir(args[0])
}

func pwd(_ []string, w io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(dir + "\n"))
	return err
}

func kill(args []string, _ io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: missing operand")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid PID")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Signal(syscall.SIGTERM)
}

func ps(_ []string, w io.Writer) error {
	processes, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}

	for _, process := range processes {
		if _, err = strconv.Atoi(process.Name()); err == nil {
			_, err = w.Write([]byte(process.Name() + "\n"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
