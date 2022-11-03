package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/build-image-action/pkg"
	"github.com/vmware-tanzu/build-image-action/pkg/version"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

//go:generate go run -modfile ../hack/tools/go.mod github.com/maxbrunsfeld/counterfeiter/v6 -generate

const sleepTimeBetweenChecks = 3

type Config struct {
	Namespace          string
	CaCert             string
	Server             string
	Token              string
	GitServer          string
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

func GetClusterBuilderStatus(ctx context.Context, client client.Client, name string) (string, string, error) {
	builder := &v1alpha2.ClusterBuilder{}
	err := client.Get(ctx, types.NamespacedName{Name: name}, builder)
	if err != nil {
		return "", "", err
	}

	return builder.Status.LatestImage, builder.Status.Stack.RunImage, nil
}

func getBuildStatus(ctx context.Context, client client.Client, namespace string, name string) (string, string, string, error) {
	build := &v1alpha2.Build{}
	err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, build)
	if err != nil {
		return "", "", "", err
	}

	var statusMessage string
	for _, condition := range build.Status.Conditions {
		if condition.Type == "Succeeded" {
			if condition.Status == "False" {
				statusMessage = condition.Message
			}
		}
	}

	return build.Status.PodName, build.Status.LatestImage, statusMessage, nil
}

func WaitForBuildToStart(ctx context.Context, client client.Client, namespace string, name string) error {
	for {
		var podName string
		var statusMessage string
		podName, _, statusMessage, err := getBuildStatus(ctx, client, namespace, name)
		if err != nil {
			return err
		}

		if statusMessage != "" {
			return errors.New(statusMessage)
		}

		if podName != "" {
			fmt.Printf("::debug:: build has started\n")
			fmt.Printf("::debug:: Building... podName=%s, starting streaming\n", podName)
			//StreamPodLogs(ctx, client, namespace, podName)
			break
		}

		time.Sleep(sleepTimeBetweenChecks * time.Second)
	}

	return nil
}

func (c *Config) Build(ctx context.Context, client client.Client) error {
	builderLatestImage, builderStackRunImage, err := GetClusterBuilderStatus(ctx, client, "my-builder")
	if err != nil {
		fmt.Printf("::debug:: failed to get cluster builder: %+v", err)
		return err
	}

	buildName, err := CreateBuild(ctx, client, &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kpack.io/v1alpha2",
			"kind":       "Build",
			"metadata": map[string]interface{}{
				"generateName": strings.ReplaceAll(c.GitRepo, "/", "-") + "-",
				"namespace":    c.Namespace,
				"annotations": map[string]interface{}{
					"app.kubernetes.io/managed-by": "vmware-tanzu/build-image-action " + version.Version,
				},
			},
			"spec": map[string]interface{}{
				"builder": map[string]interface{}{
					"image": builderLatestImage,
				},
				"runImage": map[string]interface{}{
					"image": builderStackRunImage,
				},
				"serviceAccountName": c.ServiceAccountName,
				"source": map[string]interface{}{
					"git": map[string]interface{}{
						"url":      fmt.Sprintf("%s/%s", c.GitServer, c.GitRepo),
						"revision": c.GitSha,
					},
				},
				"tags": []string{
					c.Tag,
				},
				"env": pkg.KeyValueArray(pkg.ParseEnvVars(c.Env)),
			},
		},
	})
	if err != nil {
		fmt.Printf("::debug:: failed to create build: %+v", err)
		return err
	}

	err = WaitForBuildToStart(ctx, client, c.Namespace, buildName)
	if err != nil {
		fmt.Printf("::debug:: build did not start: %+v", err)
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

	config.GitServer = pkg.MustGetEnv("GITHUB_SERVER_URL")
	config.GitRepo = pkg.MustGetEnv("GITHUB_REPOSITORY")
	config.GitSha = pkg.MustGetEnv("GITHUB_SHA")

	config.Tag = pkg.MustGetEnv("TAG")
	config.Env = os.Getenv("ENV_VARS")
	config.ServiceAccountName = os.Getenv("SERVICE_ACCOUNT_NAME")
	config.GithubOutput = pkg.MustGetEnv("GITHUB_OUTPUT")

	return kpackCmd
}
