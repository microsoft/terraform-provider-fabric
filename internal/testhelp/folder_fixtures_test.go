// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package testhelp_test

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUnit_FolderFixtures_WorkflowIsolation(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the workflow regression test")
	}

	root := filepath.Join(filepath.Dir(filename), "..", "..")
	type concurrencyConfig struct {
		Group            string `yaml:"group"`
		CancelInProgress bool   `yaml:"cancel-in-progress"` //nolint:tagliatelle // GitHub Actions uses this field name.
	}
	type workflowStep struct {
		ID   string            `yaml:"id"`
		Uses string            `yaml:"uses"`
		With map[string]string `yaml:"with"`
	}
	type workflowJob struct {
		Concurrency yaml.Node         `yaml:"concurrency"`
		If          string            `yaml:"if"`
		Needs       yaml.Node         `yaml:"needs"`
		Outputs     map[string]string `yaml:"outputs"`
		Steps       []workflowStep    `yaml:"steps"`
	}
	type workflow struct {
		Concurrency concurrencyConfig      `yaml:"concurrency"`
		Jobs        map[string]workflowJob `yaml:"jobs"`
	}
	read := func(name string) workflow {
		t.Helper()

		data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", name))
		if err != nil {
			t.Fatal(err)
		}
		var result workflow

		err = yaml.Unmarshal(data, &result)
		if err != nil {
			t.Fatal(err)
		}

		return result
	}
	tests := read("test.yml")

	cleanup := read("cleanup-script.yml")
	if tests.Concurrency.Group == "" || strings.Contains(tests.Concurrency.Group, "${{") ||
		tests.Concurrency.Group != cleanup.Concurrency.Group {
		t.Fatal("acceptance tests and cleanup must hold the same workflow-level environment lock")
	}

	if tests.Concurrency.CancelInProgress || cleanup.Concurrency.CancelInProgress {
		t.Fatal("a queued run must not cancel active resource creation or teardown")
	}

	changes, ok := tests.Jobs["changes"]
	if !ok {
		t.Fatal("changes job is missing")
	}

	// Verify the real YAML wiring, not Picomatch or GitHub expression execution.
	filterID := ""
	var filters map[string][]string

	for _, step := range changes.Steps {
		if !strings.HasPrefix(step.Uses, "dorny/paths-filter@") {
			continue
		}

		require.Empty(t, filterID, "changes must have only one paths-filter step")
		require.NotEmpty(t, step.ID, "paths-filter step must have an ID")

		filterID = step.ID
		require.NoError(t, yaml.Unmarshal([]byte(step.With["filters"]), &filters))
	}

	if filterID == "" || changes.Outputs["src"] != "${{ steps."+filterID+".outputs.src }}" {
		t.Fatal("changes.src must expose the paths-filter src output")
	}

	requiredPatterns := []string{
		"**.go", "go.mod", "go.sum",
		".github/workflows/test.yml", ".github/workflows/cleanup-script.yml",
	}
	for _, pattern := range requiredPatterns {
		if !slices.Contains(filters["src"], pattern) {
			t.Fatalf("src filter must include %q", pattern)
		}
	}

	if len(filters["src"]) != len(requiredPatterns) {
		t.Fatal("src must retain the existing selection scope without documentation paths")
	}

	for _, name := range []string{"cleanup-workspace", "test-auth-spn", "test"} {
		job, exists := tests.Jobs[name]
		if !exists {
			t.Fatalf("development job %q is missing", name)
		}

		if job.Concurrency.Kind != 0 {
			t.Fatalf("development job %q must use the workflow lock", name)
		}

		if job.If != "needs.changes.outputs.src == 'true'" {
			t.Fatalf("development job %q must use changes.src", name)
		}

		var needs []string
		if job.Needs.Kind == yaml.ScalarNode {
			needs = []string{job.Needs.Value}
		} else {
			require.NoError(t, job.Needs.Decode(&needs))
		}

		if !slices.Contains(needs, "changes") ||
			(name != "cleanup-workspace" && !slices.Contains(needs, "cleanup-workspace")) {
			t.Fatalf("development job %q must wait for changes and its cleanup prerequisite", name)
		}
	}

	cleanupJob, ok := cleanup.Jobs["cleanup-workspace"]
	if !ok || cleanupJob.Concurrency.Kind != 0 || cleanupJob.If != "" || cleanupJob.Needs.Kind != 0 {
		t.Fatal("standalone cleanup must remain runnable under the shared workflow lock")
	}
}
