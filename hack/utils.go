package hack

import (
	"context"
	"github.com/vmware-tanzu/build-image-action/hack/cli"
	"github.com/vmware-tanzu/build-image-action/hack/log"
)

func EnsureNamespaceExists(namespace string) func(context.Context) error {
	return func(ctx context.Context) error {
		ctx = log.ToContext(ctx, log.L(ctx).
			WithName("ensure-namespace-exists").
			WithValues("namespace", namespace))
		log.L(ctx).Info("ensuring")

		if namespaceExists(ctx, namespace) {
			log.L(ctx).V(1).Info("namespace already exists")
			return nil
		}

		return cli.Kubectl().
			Create("namespace").
			Name(namespace).
			Run(ctx)
	}
}

func namespaceExists(ctx context.Context, ns string) bool {
	return cli.Kubectl().Get("namespace/"+ns).Run(ctx) == nil
}

func secretExists(ctx context.Context, ns string, name string) bool {
	return cli.Kubectl().Namespace(ns).Get("secret/"+name).Run(ctx) == nil
}

func ApplyExtraObjects() func(context.Context) error {
	return func(ctx context.Context) error {
		ctx = log.ToContext(ctx, log.L(ctx).
			WithName("extra-objects"))

		log.L(ctx).Info("applying")

		return cli.Kubectl().
			Apply().
			F("./hack/kpack.yaml").
			Run(ctx)
	}
}

func ApplyRbac() func(context.Context) error {
	return func(ctx context.Context) error {
		ctx = log.ToContext(ctx, log.L(ctx).
			WithName("rbac"))

		log.L(ctx).Info("applying")

		return cli.Kubectl().
			Apply().
			F("./config/rbac.yaml").
			Run(ctx)
	}
}

func CreateSecret(namespace string) func(context.Context) error {
	return func(ctx context.Context) error {
		ctx = log.ToContext(ctx, log.L(ctx).
			WithName("create-secret").
			WithValues("namespace", namespace))
		log.L(ctx).Info("creating")

		if secretExists(ctx, namespace, "kpack-registry-credentials") {
			log.L(ctx).V(1).Info("secret already exists")
			return nil
		}

		return cli.Kubectl().
			Create("secret", "docker-registry").
			Name("kpack-registry-credentials").
			Flags("--docker-username=_json_key",
				"--docker-password=TODO",
				"--docker-server=gcr.io").
			Namespace(namespace).
			Run(ctx)
	}
}
