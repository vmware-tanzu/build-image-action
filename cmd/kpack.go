package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vmware-tanzu/build-image-action/pkg/kpack"
)

var TIMEOUT = 3600

func init() {
	rootCmd.AddCommand(NewKpackCmd())
}

func NewKpackCmd() *cobra.Command {
	c := kpack.Config{}
	var kpackCmd = &cobra.Command{
		Use:   "kpack",
		Short: "Create a kpack build on cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("::debug:: kpack build")

			fmt.Println("::debug:: tag", c.Tag)
			fmt.Println("::debug:: namespace", c.Namespace)
			fmt.Println("::debug:: gitServer", c.GitServer)
			fmt.Println("::debug:: gitRepo", c.GitRepo)
			fmt.Println("::debug:: gitSha", c.GitSha)
			fmt.Println("::debug:: env", c.Env)
			fmt.Println("::debug:: serviceAccountName", c.ServiceAccountName)
			fmt.Println("::debug:: clusterBuilder", c.ClusterBuilderName)
			fmt.Println("::debug:: timeout", c.Timeout)
			fmt.Println("::debug:: cleanup", c.Cleanup)

			err := c.Build()
			if err != nil {
				fmt.Printf("::error:: error performing build %+v\n", err)
				return err
			}
			return nil
		},
	}

	kpackCmd.Flags().StringVarP(&c.CaCert, "ca-cert", "c", "", "ca cert to access cluster")
	kpackCmd.Flags().StringVarP(&c.Token, "token", "t", "", "token to access cluster")
	kpackCmd.Flags().StringVarP(&c.Server, "server", "s", "", "server address of cluster")
	kpackCmd.MarkFlagsRequiredTogether("ca-cert", "token", "server")

	kpackCmd.Flags().StringVarP(&c.Namespace, "namespace", "n", "", "kubernetes namespace to create the build")
	_ = kpackCmd.MarkFlagRequired("namespace")
	kpackCmd.Flags().StringVarP(&c.GitServer, "github-server-url", "u", "", "github server url for the source location")
	_ = kpackCmd.MarkFlagRequired("github-server-url")
	kpackCmd.Flags().StringVarP(&c.GitRepo, "github-repository", "r", "", "github repository for the source location")
	_ = kpackCmd.MarkFlagRequired("github-repository")
	kpackCmd.Flags().StringVar(&c.GitSha, "github-sha", "", "sha of source to build")
	_ = kpackCmd.MarkFlagRequired("github-sha")
	kpackCmd.Flags().StringVar(&c.Tag, "tag", "", "docker tag to build")
	_ = kpackCmd.MarkFlagRequired("tag")
	kpackCmd.Flags().StringVarP(&c.Env, "env-vars", "e", "", "list of build time environment variables")
	kpackCmd.Flags().StringVarP(&c.ServiceAccountName, "service-account-name", "a", "default", "service account name that will be used for credential lookup")
	kpackCmd.Flags().StringVarP(&c.ClusterBuilderName, "cluster-builder", "b", "default", "cluster builder to use for the build")
	kpackCmd.Flags().Int64Var(&c.Timeout, "timeout", int64(TIMEOUT), "max active time that the pod can run for in seconds")

	kpackCmd.Flags().StringVarP(&c.ActionOutput, "github-action-output", "o", "", "location to store output of the build")
	_ = kpackCmd.MarkFlagRequired("github-action-output")

	kpackCmd.Flags().BoolVar(&c.Cleanup, "cleanup", true, "delete build resource")

	return kpackCmd
}
