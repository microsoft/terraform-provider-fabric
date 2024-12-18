// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const wellKnownEnvKey = "FABRIC_TESTACC_WELLKNOWN"

var wellKnownFilePath = getFixtureFilePath(".wellknown.json")

// wellKnownData     map[string]any

func IsWellKnownDataAvailable() bool {
	wellKnownJSON := os.Getenv(wellKnownEnvKey)

	// if the environment variable is set, we don't need to check the file
	if wellKnownJSON != "" {
		return true
	}

	// check if the file exists
	_, err := os.Stat(wellKnownFilePath)

	return !os.IsNotExist(err)
}

func getFixtureFilePath(sourcePath string) string {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled

	return filepath.Join(filepath.Dir(filename), "fixtures", sourcePath)
}

func WellKnown() map[string]any {
	if !IsWellKnownDataAvailable() {
		panicOnError(fmt.Errorf("well-known resources file %s does not exist", wellKnownFilePath))
	}

	var wk map[string]any

	if wellKnownJSON, ok := os.LookupEnv(wellKnownEnvKey); ok {
		err := json.Unmarshal([]byte(wellKnownJSON), &wk)
		panicOnError(err)

		return wk
	}

	// read the file into a string
	wellKnownJSONBytes, err := os.ReadFile(wellKnownFilePath)
	panicOnError(err)

	// parse the json string
	err = json.Unmarshal(wellKnownJSONBytes, &wk)
	panicOnError(err)

	return wk
}

func panicOnError(err error) {
	if err != nil {
		panic(err) // lintignore:R009
	}
}
