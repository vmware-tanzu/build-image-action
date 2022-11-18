package tests

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/build-image-action/hack/cli"
	"github.com/vmware-tanzu/build-image-action/hack/log"
	"github.com/vmware-tanzu/build-image-action/hack/run"
	"testing"
)

func TestAct(t *testing.T) {
	ctx := log.ToContext(context.Background(), log.New().
		WithName("act"),
	)

	if err := run.ConcurrentlyWithContext(ctx,
		cli.Act().Run,
	); err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
}
