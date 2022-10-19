// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"context"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"os"
	"regexp"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/stern/stern/stern"
	"k8s.io/apimachinery/pkg/fields"
)

type SternTailer struct {
}

func (s *SternTailer) Tail(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, podName string) error {
	t := "{{color .PodColor \"[\"}}{{color .PodColor .ContainerName}}{{color .PodColor \"]\"}} {{.Message}}\n"

	functions := map[string]interface{}{
		"color": func(color color.Color, text string) string {
			return color.SprintFunc()(text)
		},
	}
	parsedTemplate, err := template.New("log").Funcs(functions).Parse(t)
	if err != nil {
		panic(err)
	}

	configStern := stern.Config{
		Namespaces:     []string{namespace},
		Location:       time.Local,
		LabelSelector:  labels.Everything(),
		ContainerQuery: regexp.MustCompile(".*"),
		ContainerStates: []stern.ContainerState{
			stern.RUNNING,
		},
		InitContainers: true,
		Since:          1 * time.Second,
		PodQuery:       regexp.MustCompile(podName),
		FieldSelector:  fields.Everything(),
		Template:       parsedTemplate,
		Out:            os.Stdout,
		ErrOut:         os.Stderr,
		Follow:         true,
	}

	return Run(ctx, clientSet, &configStern)
}
