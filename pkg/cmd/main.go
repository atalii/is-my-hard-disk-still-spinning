package cmd

import (
	"fmt"
	"os/exec"
)

func Runner(cmd string, args ...string) func() (*string, *string) {
	return func() (*string, *string) {
		cmd := exec.Command(cmd, args...)
		out, err := cmd.CombinedOutput()
		val := string(out)

		if err != nil {
			msg := fmt.Sprintf("%s: %v\n\n%s", cmd, err, val)
			return nil, &msg
		}

		return &val, nil
	}
}

