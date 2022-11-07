// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package kpack

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	"github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"github.com/vmware-tanzu/build-image-action/pkg"
	"github.com/vmware-tanzu/build-image-action/pkg/logs"
	"github.com/vmware-tanzu/build-image-action/pkg/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const sleepTimeBetweenChecks = 3

var (
	v1alpha2Builds         = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "builds"}
	v1alpha2ClusterBuilder = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "clusterbuilders"}
)

type Config struct {
	CaCert    string
	Token     string
	Server    string
	Namespace string

	GitServer string
	GitRepo   string
	GitSha    string

	Tag                string
	Env                string
	ServiceAccountName string
	ClusterBuilderName string

	ActionOutput string
}

func (c *Config) Build() {
	decodedCaCert, err := base64.StdEncoding.DecodeString(c.CaCert)
	if err != nil {
		panic(err)
	}

	var config *rest.Config

	if c.CaCert == "" && c.Server == "" && c.Token == "" {
		// assume we are currently running inside the cluster we want to create the image resource in
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	} else {
		config = &rest.Config{
			TLSClientConfig: rest.TLSClientConfig{
				CAData: decodedCaCert,
			},
			Host:        c.Server,
			BearerToken: c.Token,
		}
	}

	ctx := context.Background()

	restMapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
	restMapper.Add(schema.GroupVersionKind{Group: "kpack.io", Version: "v1alpha2", Kind: "ClusterBuilder"}, meta.RESTScopeRoot)
	restMapper.Add(schema.GroupVersionKind{Group: "kpack.io", Version: "v1alpha2", Kind: "Build"}, meta.RESTScopeNamespace)

	client, err := client.New(config, client.Options{Mapper: restMapper})
	if err != nil {
		panic(err)
	}

	v1alpha2.AddToScheme(scheme.Scheme)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	clusterBuilder, runImage, err := GetClusterBuilder(ctx, client, c.ClusterBuilderName)
	if err != nil {
		panic(err)
	}

	build := &v1alpha2.Build{
		TypeMeta: metav1.TypeMeta{
			Kind:       "build",
			APIVersion: "kpack.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: strings.ReplaceAll(c.GitRepo, "/", "-") + "-",
			Namespace:    c.Namespace,
			Annotations: map[string]string{
				"app.kubernetes.io/managed-by": "vmware-tanzu/build-image-action " + version.Version,
			},
		},
		Spec: v1alpha2.BuildSpec{
			Tags: []string{c.Tag},
			Builder: v1alpha1.BuildBuilderSpec{
				Image: clusterBuilder,
			},
			ServiceAccountName: c.ServiceAccountName,
			Source: v1alpha1.SourceConfig{
				Git: &v1alpha1.Git{
					URL:      fmt.Sprintf("%s/%s", c.GitServer, c.GitRepo),
					Revision: c.GitSha,
				},
				Blob:     nil,
				Registry: nil,
				SubPath:  "",
			},
			RunImage: v1alpha2.BuildSpecImage{
				Image: runImage,
			},
			Env: KeyValueArray(pkg.ParseEnvVars(c.Env)),
		},
	}

	name, err := CreateBuild(ctx, client, build)
	if err != nil {
		panic(err)
	}

	for {
		var podName string
		var statusMessage string
		podName, _, statusMessage, err = GetBuildStatus(ctx, client, c.Namespace, name)
		if err != nil {
			panic(err)
		}

		if statusMessage != "" {
			panic(statusMessage)
		}

		if podName != "" {
			fmt.Printf("::debug:: build has started\n")
			fmt.Printf("::debug:: Building... podName=%s, starting streaming\n", podName)
			StreamPodLogs(ctx, clientset, c.Namespace, podName)
			break
		}

		time.Sleep(sleepTimeBetweenChecks * time.Second)
	}

	for {
		fmt.Printf("::debug:: checking if build is complete...\n")
		var latestImage string
		var statusMessage string
		_, latestImage, statusMessage, err = GetBuildStatus(ctx, client, c.Namespace, name)
		if err != nil {
			panic(err)
		}

		if statusMessage != "" {
			panic(statusMessage)
		}

		if latestImage != "" {
			fmt.Printf("::debug:: build is complete\n")

			err = Append(c.ActionOutput, fmt.Sprintf("name=%s\n", latestImage))
			if err != nil {
				panic(err)
			}
			break
		}

		time.Sleep(sleepTimeBetweenChecks * time.Second)
	}
}

func GetClusterBuilderStatus(ctx context.Context, client client.Client, name string) (string, string, error) {
	builder := &v1alpha2.ClusterBuilder{}
	err := client.Get(ctx, types.NamespacedName{Name: name}, builder)
	if err != nil {
		return "", "", err
	}

	return builder.Status.LatestImage, builder.Status.Stack.RunImage, nil
}

func CreateBuild(ctx context.Context, client client.Client, build *v1alpha2.Build) (string, error) {
	fmt.Printf("::debug:: creating resource %+v\n", build)

	err := client.Create(ctx, build)
	if err != nil {
		return "", err
	}

	return build.GetName(), nil
}

func GetBuildStatus(ctx context.Context, client client.Client, namespace string, name string) (string, string, string, error) {
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

func KeyValueArray(vars map[string]string) []corev1.EnvVar {
	var values []corev1.EnvVar
	for k, v := range vars {
		values = append(values, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	fmt.Printf("::debug:: parsed environment variables to %s\n", values)
	return values
}

func StreamPodLogs(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string) {
	go func() {
		st := logs.SternTailer{}
		err := st.Tail(ctx, clientSet, namespace, podName)
		if err != nil {
			log.Fatalf("issue streaming logs: %s", err)
		}
	}()
}

func Append(file string, name string) error {
	const filePermissions = 0644
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(name); err != nil {
		return err
	}
	return nil
}
