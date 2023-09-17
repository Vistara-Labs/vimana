package utils

import (
	"fmt"
	"io"
	"net/http"
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
	fmt.Println("Command execution completed", cmd)
	return nil
}

func downloadFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return "", err
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}
