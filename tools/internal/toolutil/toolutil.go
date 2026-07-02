// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

// Package toolutil provides shared helpers for internal CLI tools
// (previewcheck, gapcheck, etc.).
package toolutil

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	// SDKModulePath is the Go module path of the Fabric SDK.
	SDKModulePath = "github.com/microsoft/fabric-sdk-go"

	// DefaultServicesDir is the default path to the services directory.
	DefaultServicesDir = "internal/services"
)

// ModuleRoot walks up from the working directory to find the directory holding
// go.mod.
func ModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s upward", dir)
		}

		dir = parent
	}
}

// SortedKeys returns the sorted keys of a string set.
func SortedKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	return keys
}

// Outf writes formatted output to stdout (for CLI tool results).
func Outf(format string, a ...any) { fmt.Printf(format, a...) } //nolint:forbidigo

// Errf writes formatted output to stderr (for CLI tool diagnostics).
func Errf(format string, a ...any) { fmt.Fprintf(os.Stderr, format, a...) }

// LoadServicePackages type-checks every package under servicesDir (relative to
// root) and returns them. It loads full type information and dependencies so
// callers can resolve identifiers to their declared objects (e.g. fabric-sdk-go
// client methods reached through local variables, parameters, or embedding).
func LoadServicePackages(root, servicesDir string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports |
			packages.NeedDeps,
		Dir:   root,
		Tests: false,
	}

	pattern := "./" + filepath.ToSlash(servicesDir) + "/..."

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, fmt.Errorf("loading packages under %q: %w", servicesDir, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found under %q", servicesDir)
	}

	return pkgs, nil
}

// ParseClientMethod parses a Go source line declaring a method, of the form
//
//	func (recv *Type) Method(...)
//	func (recv Type) Method(...)
//
// It returns the receiver's bare type name (with any leading '*' stripped) and
// the method name, or empty strings if the line is not a method declaration or
// the method is unexported. It is a lightweight line scanner for
// machine-generated SDK client files, not a full Go parser.
func ParseClientMethod(line string) (recvType, method string) { //nolint:nonamedreturns // named for clarity: receiver type + method name
	const prefix = "func ("
	if !strings.HasPrefix(line, prefix) {
		return "", ""
	}

	// Isolate the receiver clause between the opening '(' and its closing ')'.
	closeIdx := strings.IndexByte(line, ')')
	if closeIdx <= len(prefix) {
		return "", ""
	}

	// The receiver type is the last field, e.g. "client *WorkspacesClient".
	recv := strings.Fields(line[len(prefix):closeIdx])
	if len(recv) == 0 {
		return "", ""
	}

	recvType = strings.TrimPrefix(recv[len(recv)-1], "*")
	if recvType == "" {
		return "", ""
	}

	// The method name follows the receiver clause, up to its parameter list.
	rest := strings.TrimSpace(line[closeIdx+1:])

	paren := strings.IndexByte(rest, '(')
	if paren <= 0 {
		return "", ""
	}

	method = rest[:paren]
	if method[0] < 'A' || method[0] > 'Z' {
		return "", "" // unexported, so not a public SDK method
	}

	return recvType, method
}
