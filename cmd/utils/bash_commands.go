package utils

import (
	"fmt"
	"os"
	"os/exec"
)

type CommandOption func(cmd *exec.Cmd)

func WithOutputToStdout() CommandOption {
	return func(cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
	}
}

func WithErrorsToStderr() CommandOption {
	return func(cmd *exec.Cmd) {
		cmd.Stderr = os.Stderr
	}
}

func ExecBashCmd(cmd *exec.Cmd, options ...CommandOption) error {
	for _, option := range options {
		option(cmd)
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}
	return nil
}
