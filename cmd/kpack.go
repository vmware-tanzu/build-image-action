package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/build-image-action/pkg/kpack"
	"log"
	"os"
)

func init() {
	rootCmd.AddCommand(kpackCmd)
}

var kpackCmd = &cobra.Command{
	Use:   "kpack",
	Short: "Create a kpack build on cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("::debug:: kpack build")

		c := kpack.Config{
			CaCert:             os.Getenv("CA_CERT"),
			Token:              os.Getenv("TOKEN"),
			Server:             os.Getenv("SERVER"),
			Namespace:          MustGetEnv("NAMESPACE"),
			GitServer:          MustGetEnv("GITHUB_SERVER_URL"),
			GitRepo:            MustGetEnv("GITHUB_REPOSITORY"),
			GitSha:             MustGetEnv("GITHUB_SHA"),
			Tag:                MustGetEnv("TAG"),
			Env:                os.Getenv("ENV_VARS"),
			ServiceAccountName: os.Getenv("SERVICE_ACCOUNT_NAME"),
			ActionOutput:       MustGetEnv("GITHUB_OUTPUT"),
		}

		fmt.Println("::debug:: tag", c.Tag)
		fmt.Println("::debug:: namespace", c.Namespace)
		fmt.Println("::debug:: gitServer", c.GitServer)
		fmt.Println("::debug:: gitRepo", c.GitRepo)
		fmt.Println("::debug:: gitSha", c.GitSha)
		fmt.Println("::debug:: env", c.Env)
		fmt.Println("::debug:: serviceAccountName", c.ServiceAccountName)

		c.Build()
	},
}

func MustGetEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Environment Var %s must be set", name)
	}
	return val
}
