package run

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/vmware-tanzu/build-image-action/hack/log"
	"os/exec"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
)

type cmd struct {
	argv         []string
	envOverrides map[string]string
	stdout       *bytes.Buffer
	stderr       *bytes.Buffer
}

func Cmd(argv ...string) *cmd {
	return &cmd{
		argv:         argv,
		envOverrides: map[string]string{},
		stdout:       new(bytes.Buffer),
		stderr:       new(bytes.Buffer),
	}
}

func (c *cmd) Env(key, value string) *cmd {
	c.envOverrides[key] = value
	return c
}

func (c *cmd) Run(ctx context.Context) error {
	logger := log.L(ctx).
		WithName("run").
		WithValues("argv", prettyArgv(c.argv))

	ctx = log.ToContext(ctx, logger)

	if err := c.run(ctx); err != nil {
		if c.stdout.Len() > 0 {
			err = multierror.Append(err, fmt.Errorf(
				"stdout: %s", c.stdout.String(),
			))
		}

		if c.stderr.Len() > 0 {
			err = multierror.Append(err, fmt.Errorf(
				"stderr: %s", c.stderr.String(),
			))
		}

		return err
	}

	return nil
}

func (c *cmd) RunWithOutput(ctx context.Context) (string, string, error) {
	if err := c.Run(ctx); err != nil {
		return "", "", err
	}

	return c.stdout.String(), c.stderr.String(), nil
}

func (c *cmd) run(ctx context.Context) error {
	logger := log.L(ctx)

	command := exec.CommandContext(ctx, c.argv[0], c.argv[1:]...)
	command.Stdout = c.stdout
	command.Stderr = c.stderr
	for key, value := range c.envOverrides {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value))
	}

	errC := make(chan error)

	go func() {
		errC <- command.Run()
		close(errC)
	}()

	debugTicker := time.NewTicker(5 * time.Second)
	startTime := time.Now()

	for {
		select {
		case <-debugTicker.C:
			logger.V(1).Info("still running",
				"elapsed", duration.ShortHumanDuration(
					time.Since(startTime),
				))
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errC:
			if err != nil {
				return fmt.Errorf("run: %w", err)
			}

			return nil
		}
	}
}

// wrapSlice wraps a slice of strings of `n` elements to one of size `maxItems`
// where, in case of wrapping, the last element is replaced by ellipsis
// (`...`).
//
// e.g.:
//
//	wrapSlice([]string{"a","b","c"}, 2) == []string{"a","..."}
func prettyArgv(slice []string) string {
	const maxItems = 8

	n := min(len(slice), maxItems)
	res := make([]string, n)
	copy(res, slice[0:n])

	if n == maxItems {
		res[n-1] = "..."
	}

	return strings.Join(res, " ")
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
