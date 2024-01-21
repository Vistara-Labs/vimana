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

func ExecBashCmd(cmd *exec.Cmd, node_info string, options ...CommandOption) error {
	for _, option := range options {
		option(cmd)
	}
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	pid := cmd.Process.Pid
	PIDFile := GetPIDFileName(node_info)
	savePID(pid, PIDFile)

	// use goroutine waiting, manage process
	// this is important, otherwise the process becomes in S mode
	go func() {
		err = cmd.Wait()
		fmt.Printf("Command finished with error: %v", err)
	}()
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
