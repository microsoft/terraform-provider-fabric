// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package testhelp_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

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
	type workflow struct {
		Concurrency concurrencyConfig         `yaml:"concurrency"`
		Jobs        map[string]map[string]any `yaml:"jobs"`
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

	if tests.Jobs["test"] == nil {
		t.Fatal("acceptance test matrix job is missing")
	}

	if _, ok := tests.Jobs["test"]["concurrency"]; ok {
		t.Fatal("a per-matrix-job lock would discard pending test matrix legs")
	}
}
