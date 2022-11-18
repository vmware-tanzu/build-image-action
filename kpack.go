//go:build mage
// +build mage

package main

import (
	"context"
	"github.com/magefile/mage/mg"
	"github.com/vmware-tanzu/build-image-action/hack"
	"github.com/vmware-tanzu/build-image-action/hack/cli"
	"github.com/vmware-tanzu/build-image-action/hack/log"
)

func Kpack(ctx context.Context) error {
	mg.CtxDeps(ctx,
		KpackController,
		hack.EnsureNamespaceExists("dev"),
		hack.CreateSecret("dev"),
		hack.ApplyExtraObjects(),
		hack.ApplyRbac(),
	)

	return nil
}

func KpackController(ctx context.Context) error {
	ctx = log.ToContext(ctx, log.L(ctx).WithName("kpack"))

	log.L(ctx).Info("installing kpack")
	return installKpack(ctx)
}

func installKpack(ctx context.Context) error {
	return cli.Kubectl().
		Apply().
		F("https://github.com/pivotal/kpack/releases/download/v0.7.2/release-0.7.2.yaml").
		Run(ctx)
}
