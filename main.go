// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vmware-tanzu/build-image-action/cmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"os"
)

func main() {
	cmd.Execute()
}

const sleepTimeBetweenChecks = 3

var (
	v1alpha2Builds         = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "builds"}
	v1alpha2ClusterBuilder = schema.GroupVersionResource{Group: "kpack.io", Version: "v1alpha2", Resource: "clusterbuilders"}
)

func GetBuild(ctx context.Context, client dynamic.Interface, namespace string, build string) (string, string, string, error) {
	got, err := client.Resource(v1alpha2Builds).Namespace(namespace).Get(ctx, build, metav1.GetOptions{})
	if err != nil {
		return "", "", "", err
	}

	podName, _, err := unstructured.NestedString(got.Object, "status", "podName")
	if err != nil {
		return "", "", "", err
	}

	latestImage, _, err := unstructured.NestedString(got.Object, "status", "latestImage")
	if err != nil {
		return "", "", "", err
	}

	conditions, _, err := unstructured.NestedSlice(got.Object, "status", "conditions")
	if err != nil {
		return "", "", "", err
	}

	var statusMessage string
	for _, condition := range conditions {
		conditionObj, ok := condition.(map[string]interface{})
		if !ok {
			return "", "", "", errors.New("unable to cast condition to map[string]interface{}")
		}
		if conditionObj["type"] == "Succeeded" {
			if conditionObj["status"] == "False" {
				statusMessage, ok = conditionObj["message"].(string)
				if !ok {
					return "", "", "", errors.New("unable to cast condition message to string")
				}
			}
		}
	}

	return podName, latestImage, statusMessage, nil
}

//func main() {
//	caCert := os.Getenv("CA_CERT")
//	server := os.Getenv("SERVER")
//	namespace := MustGetEnv("NAMESPACE")
//	token := os.Getenv("TOKEN")
//
//	gitRepo := fmt.Sprintf("%s/%s", MustGetEnv("GITHUB_SERVER_URL"), MustGetEnv("GITHUB_REPOSITORY"))
//	gitSha := MustGetEnv("GITHUB_SHA")
//	tag := MustGetEnv("TAG")
//	env := os.Getenv("ENV_VARS")
//	serviceAccountName := os.Getenv("SERVICE_ACCOUNT_NAME")
//	githubOutput := MustGetEnv("GITHUB_OUTPUT")
//
//	fmt.Println("::debug:: tag", tag)
//	fmt.Println("::debug:: namespace", namespace)
//	fmt.Println("::debug:: gitRepo", gitRepo)
//	fmt.Println("::debug:: gitSha", gitSha)
//	fmt.Println("::debug:: env", env)
//	fmt.Println("::debug:: serviceAccountName", serviceAccountName)
//
//	decodedCaCert, err := base64.StdEncoding.DecodeString(caCert)
//	if err != nil {
//		panic(err)
//	}
//
//	var config *rest.Config
//
//	if caCert == "" && server == "" && token == "" {
//		// assume we are currently running inside the cluster we want to create the image resource in
//		config, err = rest.InClusterConfig()
//		if err != nil {
//			panic(err)
//		}
//	} else {
//		config = &rest.Config{
//			TLSClientConfig: rest.TLSClientConfig{
//				CAData: decodedCaCert,
//			},
//			Host:        server,
//			BearerToken: token,
//		}
//	}
//
//	ctx := context.Background()
//
//	dynamicClient, err := dynamic.NewForConfig(config)
//	if err != nil {
//		panic(err)
//	}
//
//	client, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		panic(err)
//	}
//
//	clusterBuilder, runImage, err := GetClusterBuilder(ctx, dynamicClient, "default")
//	if err != nil {
//		panic(err)
//	}
//
//	build := &unstructured.Unstructured{
//		Object: map[string]interface{}{
//			"apiVersion": "kpack.io/v1alpha2",
//			"kind":       "Build",
//			"metadata": map[string]interface{}{
//				"generateName": strings.ReplaceAll(MustGetEnv("GITHUB_REPOSITORY"), "/", "-") + "-",
//				"namespace":    namespace,
//				"annotations": map[string]interface{}{
//					"app.kubernetes.io/managed-by": "vmware-tanzu/build-image-action " + version.Version,
//				},
//			},
//			"spec": map[string]interface{}{
//				"builder": map[string]interface{}{
//					"image": clusterBuilder,
//				},
//				"runImage": map[string]interface{}{
//					"image": runImage,
//				},
//				"serviceAccountName": serviceAccountName,
//				"source": map[string]interface{}{
//					"git": map[string]interface{}{
//						"url":      gitRepo,
//						"revision": gitSha,
//					},
//				},
//				"tags": []string{
//					tag,
//				},
//				"env": KeyValueArray(pkg.ParseEnvVars(env)),
//			},
//		},
//	}
//
//	name, err := CreateBuild(ctx, dynamicClient, namespace, build)
//	if err != nil {
//		panic(err)
//	}
//
//	for {
//		var podName string
//		var statusMessage string
//		podName, _, statusMessage, err = GetBuild(ctx, dynamicClient, namespace, name)
//		if err != nil {
//			panic(err)
//		}
//
//		if statusMessage != "" {
//			panic(statusMessage)
//		}
//
//		if podName != "" {
//			fmt.Printf("::debug:: build has started\n")
//			fmt.Printf("::debug:: Building... podName=%s, starting streaming\n", podName)
//			StreamPodLogs(ctx, client, namespace, podName)
//			break
//		}
//
//		time.Sleep(sleepTimeBetweenChecks * time.Second)
//	}
//
//	for {
//		fmt.Printf("::debug:: checking if build is complete...\n")
//		var latestImage string
//		var statusMessage string
//		_, latestImage, statusMessage, err = GetBuild(ctx, dynamicClient, namespace, name)
//		if err != nil {
//			panic(err)
//		}
//
//		if statusMessage != "" {
//			panic(statusMessage)
//		}
//
//		if latestImage != "" {
//			fmt.Printf("::debug:: build is complete\n")
//
//			err = Append(githubOutput, fmt.Sprintf("name=%s\n", latestImage))
//			if err != nil {
//				panic(err)
//			}
//			break
//		}
//
//		time.Sleep(sleepTimeBetweenChecks * time.Second)
//	}
//}

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
