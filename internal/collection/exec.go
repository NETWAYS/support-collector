package collection

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

const DefaultTimeout = 5 * time.Second

func LoadCommandOutput(command string, arguments ...string) (output []byte, err error) {
	return LoadCommandOutputWithTimeout(DefaultTimeout, command, arguments...)
}

func LoadCommandOutputWithTimeout(timeout time.Duration, command string, arguments ...string) (
	output []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, arguments...)

	output, err = cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}

		// append error to output
		if err != nil {
			output = append(output, []byte("\nCommand Error: "+err.Error())...)
		}

		err = fmt.Errorf("command not successful: '%s': %w", cmd.String(), err)

		return
	}

	return
}
