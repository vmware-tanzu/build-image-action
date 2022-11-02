package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/build-image-action/pkg"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run -modfile ../hack/tools/go.mod github.com/maxbrunsfeld/counterfeiter/v6 -generate

type Config struct {
	Namespace          string
	CaCert             string
	Server             string
	Token              string
	GitRepo            string
	GitSha             string
	Tag                string
	Env                string
	ServiceAccountName string
	GithubOutput       string
}

//counterfeiter:generate sigs.k8s.io/controller-runtime/pkg/client.Client
func CreateBuild(ctx context.Context, client client.Client, build *unstructured.Unstructured) (string, error) {
	fmt.Printf("::debug:: creating resource %+v\n", build)

	err := client.Create(ctx, build)
	if err != nil {
		return "", err
	}

	return build.GetName(), nil
}

func (c *Config) Build(ctx context.Context, client client.Client) error {
	_, err := CreateBuild(ctx, client, &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kpack.io/v1alpha2",
			"kind":       "Build",
			"metadata": map[string]interface{}{
				"namespace": c.Namespace,
			},
		},
	})

	if err != nil {
		fmt.Printf("::debug:: failed to create build: %+v", err)
		return err
	}

	return nil
}

//func init() {
//	rootCmd.AddCommand(NewCmdBuild())
//}

func NewCmdBuild() *cobra.Command {

	var config *Config

	var kpackCmd = &cobra.Command{
		Use:   "kpack",
		Short: "Create kpack build",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("::debug:: tag", config.Tag)
			fmt.Println("::debug:: namespace", config.Namespace)
			//fmt.Println("::debug:: gitRepo", gitRepo)
			//fmt.Println("::debug:: gitSha", gitSha)
			//fmt.Println("::debug:: env", env)
			//fmt.Println("::debug:: serviceAccountName", serviceAccountName)

			//config.Build()
		},
	}

	config.Namespace = pkg.MustGetEnv("NAMESPACE")
	config.CaCert = os.Getenv("CA_CERT")
	config.Server = os.Getenv("SERVER")
	config.Token = os.Getenv("TOKEN")

	config.GitRepo = fmt.Sprintf("%s/%s", pkg.MustGetEnv("GITHUB_SERVER_URL"), pkg.MustGetEnv("GITHUB_REPOSITORY"))
	config.GitSha = pkg.MustGetEnv("GITHUB_SHA")
	config.Tag = pkg.MustGetEnv("TAG")
	config.Env = os.Getenv("ENV_VARS")
	config.ServiceAccountName = os.Getenv("SERVICE_ACCOUNT_NAME")
	config.GithubOutput = pkg.MustGetEnv("GITHUB_OUTPUT")

	return kpackCmd
}
