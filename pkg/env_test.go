// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkg_test

import (
	"github.com/vmware-tanzu/build-image-action/pkg"
	"reflect"
	"testing"
)

func TestEnvParsing(t *testing.T) {
	type test struct {
		input string
		want  map[string]string
	}

	tests := []test{
		{input: "", want: map[string]string{}},
		{input: "BP_JAVA_VERSION=17", want: map[string]string{"BP_JAVA_VERSION": "17"}},
		{input: "      BP_JAVA_VERSION=17", want: map[string]string{"BP_JAVA_VERSION": "17"}},
		{input: "BP_JAVA_VERSION=17     ", want: map[string]string{"BP_JAVA_VERSION": "17"}},
		{input: "BP_JAVA_VERSION=17 BP_OTHER_VALUE=something", want: map[string]string{"BP_JAVA_VERSION": "17", "BP_OTHER_VALUE": "something"}},
		{input: "BP_JAVA_VERSION=17\tBP_OTHER_VALUE=something", want: map[string]string{"BP_JAVA_VERSION": "17", "BP_OTHER_VALUE": "something"}},
		{input: "BP_JAVA_VERSION=17\nBP_OTHER_VALUE=something", want: map[string]string{"BP_JAVA_VERSION": "17", "BP_OTHER_VALUE": "something"}},
		{input: "BP_OTHER_PROPERTY=a=b", want: map[string]string{"BP_OTHER_PROPERTY": "a=b"}},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := pkg.ParseEnvVars(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
