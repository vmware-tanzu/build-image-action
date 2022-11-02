// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log"
	"os"
	"strings"
)

func ParseEnvVars(in string) map[string]string {
	m := make(map[string]string)
	in = strings.TrimSpace(in)

	for _, field := range strings.Fields(in) {
		const numberOfFields = 2
		split := strings.SplitN(field, "=", numberOfFields)
		m[split[0]] = split[1]
	}

	return m
}

func MustGetEnv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf("Environment Var %s must be set", name)
	}
	return val
}
