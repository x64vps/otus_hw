package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var ErrEmptyCommand = errors.New("empty command")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (int, error) {
	if len(cmd) == 0 {
		return 0, ErrEmptyCommand
	}

	runner := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	runner.Stdin = os.Stdin
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr

	for _, e := range os.Environ() {
		s := strings.Split(e, "=")

		v, exists := env[s[0]]
		if exists && v.NeedRemove {
			continue
		}

		runner.Env = append(runner.Env, e)
	}

	for k, v := range env {
		if v.NeedRemove {
			continue
		}

		runner.Env = append(runner.Env, k+"="+v.Value)
	}

	if err := runner.Start(); err != nil {
		return 0, err
	}

	if err := runner.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok { //nolint:errorlint
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		} else {
			return 0, err
		}
	}

	return 0, nil
}
