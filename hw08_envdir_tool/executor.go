package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	for key, val := range env {
		if val.NeedRemove {
			os.Unsetenv(key)
			continue
		}

		os.Setenv(key, val.Value)
	}

	execute := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	execute.Stdout = os.Stdout

	if err := execute.Run(); err != nil {
		return 1
	}
	return 0
}
