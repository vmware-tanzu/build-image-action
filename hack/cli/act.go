package cli

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/build-image-action/hack/run"
)

type act struct {
	argv []string
}

func Act() *act {
	return &act{}
}

func (a *act) Run(ctx context.Context) error {
	argv := []string{
		"../../emmjohnson/act/dist/local/act", //TODO
		"-j e2e",
		"--container-architecture linux/amd64",
		"-s NAMESPACE=dev",
		"-s SERVER=TODO",
		"-s TOKEN=TODO",
		"-s CA_CERT=TODO",
		"--env GITHUB_OUTPUT=output-act.txt",
		"--env GITHUB_SERVER_URL=https://github.com",
		"--env GITHUB_REPOSITORY=spring-projects/spring-petclinic",
		"--env GITHUB_SHA=9ecdc1111e3da388a750ace41a125287d9620534",
	}

	run.Cmd()
	err := run.Cmd(argv...).Run(ctx)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
