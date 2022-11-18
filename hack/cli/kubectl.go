package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vmware-tanzu/build-image-action/hack/log"
	"github.com/vmware-tanzu/build-image-action/hack/run"
)

type kubectl struct {
	cmd          []string
	name         string
	namespace    string
	files        []string
	target       interface{}
	outputFormat []string

	argv []string
}

func Kubectl() *kubectl {
	return &kubectl{}
}

func (k *kubectl) Delete(resource string) *kubectl {
	k.cmd = []string{"kubectl", "delete", resource}
	return k
}

func (k *kubectl) Apply() *kubectl {
	k.cmd = []string{"kubectl", "apply"}
	return k
}

func (k *kubectl) Create(resource ...string) *kubectl {
	k.cmd = append([]string{"kubectl", "create"}, resource...)
	return k
}

func (k *kubectl) Get(resource string) *kubectl {
	k.cmd = []string{"kubectl", "get", resource}
	return k
}

func (k *kubectl) Namespace(namespace string) *kubectl {
	k.namespace = namespace
	return k
}

func (k *kubectl) Name(name string) *kubectl {
	k.name = name
	return k
}

func (k *kubectl) F(fpath string) *kubectl {
	k.files = append(k.files, fpath)
	return k
}

func (k *kubectl) Flags(flags ...string) *kubectl {
	k.cmd = append(k.cmd, flags...)
	return k
}

func (k *kubectl) Output(o string) *kubectl {
	k.outputFormat = []string{"-o", o}
	return k
}

func (k *kubectl) Into(target interface{}) *kubectl {
	k.target = target
	if len(k.outputFormat) == 0 {
		k = k.Output("json")
	}

	return k
}

func (k *kubectl) build() error {
	if len(k.cmd) == 0 {
		return fmt.Errorf("command not set")
	}

	k.argv = k.cmd

	if k.namespace != "" {
		k.argv = append(k.argv, "-n", k.namespace)
	}

	if k.name != "" {
		k.argv = append(k.argv, k.name)
	}

	for _, fpath := range k.files {
		k.argv = append(k.argv, "-f", fpath)
	}

	if k.target != nil && len(k.outputFormat) == 0 {
		return fmt.Errorf("target set but no output format")
	}

	if len(k.outputFormat) != 0 {
		k.argv = append(k.argv, k.outputFormat...)
	}

	return nil
}

func (k *kubectl) RunWithJsonOutput(ctx context.Context) (string, error) {
	k.outputFormat = []string{"-o", "json"}

	logger := log.L(ctx).WithName("kubectl")
	ctx = log.ToContext(ctx, logger)

	if err := k.build(); err != nil {
		return "", fmt.Errorf("build: %w", err)
	}

	stdout, _, err := run.Cmd(k.argv...).RunWithOutput(ctx)
	if err != nil {
		return "", fmt.Errorf("kubectl: %w", err)
	}

	return stdout, nil
}

func (k *kubectl) Run(ctx context.Context) error {
	logger := log.L(ctx).WithName("kubectl")
	ctx = log.ToContext(ctx, logger)

	if err := k.build(); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	stdout, _, err := run.Cmd(k.argv...).RunWithOutput(ctx)
	if err != nil {
		return fmt.Errorf("kubectl: %w", err)
	}

	if k.target != nil {
		if err := json.Unmarshal([]byte(stdout), k.target); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}

	return nil
}
