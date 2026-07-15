// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package toolutil

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// Exclusion suppresses a check for a single service package, together with the
// reason the mismatch is intentional.
type Exclusion struct {
	Service string `json:"service" yaml:"service"`
	Reason  string `json:"reason"  yaml:"reason"`
}

// ExclusionsFile is the top-level structure of an exclusions YAML file.
type ExclusionsFile struct {
	Exclusions []Exclusion `yaml:"exclusions"`
}

// LoadExclusions reads an exclusions YAML file and returns a map of service
// package name to reason. Every entry must set both service and reason.
func LoadExclusions(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading exclusions file %s: %w", path, err)
	}

	var ef ExclusionsFile

	err = yaml.Unmarshal(data, &ef)
	if err != nil {
		return nil, fmt.Errorf("parsing exclusions file %s: %w", path, err)
	}

	result := make(map[string]string, len(ef.Exclusions))

	for _, e := range ef.Exclusions {
		if e.Service == "" || e.Reason == "" {
			return nil, fmt.Errorf("exclusion entry missing required field (service, reason): %+v", e)
		}

		result[e.Service] = e.Reason
	}

	return result, nil
}
