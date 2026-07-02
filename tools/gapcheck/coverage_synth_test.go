// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"

	"github.com/microsoft/terraform-provider-fabric/tools/internal/toolutil"
)

// fakeSDKCore is a minimal stand-in for a fabric-sdk-go client package. Its
// import path (github.com/microsoft/fabric-sdk-go/fabric/core) is what makes
// collectCoveredCalls treat calls into it as SDK calls.
const fakeSDKCore = `package core

type FooClient struct{}

func (c *FooClient) GetFoo() string   { return "" }
func (c *FooClient) ListFoos() string { return "" }

func NewFooClient() *FooClient { return &FooClient{} }

// FooResponse is an SDK-path type whose methods are NOT client methods.
type FooResponse struct{}

func (r *FooResponse) Value() string { return "" }

// Helper is a package-level SDK function (no receiver).
func Helper() string { return "" }
`

// fakeService calls the fake SDK through every shape collectCoveredCalls must
// handle. Each call is annotated with the expected outcome.
const fakeService = `package svc

import fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

// wrapper embeds the SDK client so its methods are promoted.
type wrapper struct {
	*fabcore.FooClient
}

// localGetter is a service-local interface (declared outside the SDK path).
type localGetter interface {
	GetFoo() string
}

func bareVariable() {
	c := fabcore.NewFooClient()
	_ = c.GetFoo() // RECORDED: core/FooClient/GetFoo
}

func valueReceiverAndList() {
	var c fabcore.FooClient
	_ = c.ListFoos() // RECORDED: core/FooClient/ListFoos
}

func promoted(w wrapper) {
	_ = w.GetFoo() // RECORDED via embedding: core/FooClient/GetFoo
}

func methodValue() {
	c := fabcore.NewFooClient()
	f := c.GetFoo
	_ = f() // NOT recorded: call target is an identifier, not a selector
}

func viaLocalInterface(g localGetter) {
	_ = g.GetFoo() // NOT recorded: resolves to svc.localGetter, off the SDK path
}

func nonClientReceiver(r *fabcore.FooResponse) {
	_ = r.Value() // NOT recorded: receiver type does not end in "Client"
}

func packageLevelFunc() {
	_ = fabcore.Helper() // NOT recorded: no receiver
}
`

// writeSyntheticModules lays out two local modules under dir: a fake fabric-sdk
// module and a caller module that requires it via a local replace. Returning
// the caller module dir lets packages.Load type-check the realistic shape where
// the caller's import path does NOT start with the SDK module path.
func writeSyntheticModules(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	sdkDir := filepath.Join(dir, "sdk")
	coreDir := filepath.Join(sdkDir, "fabric", "core")
	svcDir := filepath.Join(dir, "svc")

	for _, d := range []string{coreDir, svcDir} {
		err := os.MkdirAll(d, 0o750)
		if err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	files := map[string]string{
		filepath.Join(sdkDir, "go.mod"):   "module " + toolutil.SDKModulePath + "\n\ngo 1.23\n",
		filepath.Join(coreDir, "core.go"): fakeSDKCore,
		filepath.Join(svcDir, "go.mod"):   "module example.com/svc\n\ngo 1.23\n\nrequire " + toolutil.SDKModulePath + " v0.0.0\n\nreplace " + toolutil.SDKModulePath + " => ../sdk\n",
		filepath.Join(svcDir, "svc.go"):   fakeService,
	}

	for path, content := range files {
		err := os.WriteFile(path, []byte(content), 0o600)
		if err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	return svcDir
}

// Test_collectCoveredCalls_shapes deterministically verifies which SDK call
// shapes collectCoveredCalls resolves, using a tiny synthetic module pair. It
// type-checks in well under a second — unlike the whole-SDK integration test —
// and pins the resolution behaviour (embedded/value/bare-variable) plus the
// documented limitations (method values, service-local interfaces, non-client
// receivers, package-level funcs).
func Test_collectCoveredCalls_shapes(t *testing.T) {
	t.Parallel()

	svcDir := writeSyntheticModules(t)

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports |
			packages.NeedDeps,
		Dir: svcDir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatalf("packages.Load: %v", err)
	}

	covered := map[string]struct{}{}

	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			for _, e := range pkg.Errors {
				t.Errorf("package %s: %v", pkg.PkgPath, e)
			}

			t.FailNow()
		}

		for _, file := range pkg.Syntax {
			collectCoveredCalls(pkg, file, covered)
		}
	}

	wantPresent := []string{
		"core/FooClient/GetFoo",   // bare variable + promoted (embedding)
		"core/FooClient/ListFoos", // value receiver
	}
	for _, key := range wantPresent {
		if _, ok := covered[key]; !ok {
			t.Errorf("expected %q to be recorded as covered; got keys %v", key, keysOf(covered))
		}
	}

	// Shapes that must NOT be recorded. Value() is the only other selector call
	// into the SDK path, so any leakage would show up as an extra key.
	if _, ok := covered["core/FooResponse/Value"]; ok {
		t.Error("non-client receiver FooResponse.Value must not be recorded")
	}

	// The only recordable keys are the two client methods; method values,
	// local-interface dispatch, and package-level funcs contribute nothing.
	if len(covered) != len(wantPresent) {
		t.Errorf("covered has unexpected extra keys: got %v, want exactly %v", keysOf(covered), wantPresent)
	}
}

func keysOf(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
