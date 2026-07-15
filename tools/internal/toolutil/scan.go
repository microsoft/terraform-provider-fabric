// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package toolutil

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

// SDKCall identifies one distinct fabric-sdk-go function invoked by a service
// package, along with the position of its declaration in the SDK source so the
// caller can inspect the doc comment that documents the API.
type SDKCall struct {
	ID  string         // "pkgName.FuncName", e.g. "core.CreateShortcut"
	Pos token.Position // position of the SDK function declaration
}

// CollectSDKCalls walks every call expression in the package and returns the
// distinct fabric-sdk-go functions it invokes. The package must be loaded with
// full type information and dependencies (see LoadServicePackages) so selector
// expressions can be resolved to their declared objects (e.g. SDK client methods
// reached through local variables, parameters, or embedding).
func CollectSDKCalls(pkg *packages.Package) []SDKCall {
	seen := map[string]struct{}{}

	var calls []SDKCall

	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Only x.Y(...) style calls; SDK funcs are always reached via a selector.
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Resolve the selector to its real declared object via the type checker.
			fn, ok := pkg.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
			if !ok || fn.Pkg() == nil {
				return true
			}

			if !strings.HasPrefix(fn.Pkg().Path(), SDKModulePath) {
				return true
			}

			// Identify the callee as "pkgName.FuncName" and skip duplicates.
			id := fn.Pkg().Name() + "." + fn.Name()
			if _, dup := seen[id]; dup {
				return true
			}

			seen[id] = struct{}{}
			calls = append(calls, SDKCall{ID: id, Pos: pkg.Fset.Position(fn.Pos())})

			return true
		})
	}

	return calls
}

// ClientMethodDoc pairs an SDK client method name with the doc comment block
// that immediately precedes its declaration.
type ClientMethodDoc struct {
	Name string
	Doc  string
}

// DocCache lazily reads SDK source files and caches their lines, the doc comment
// above declarations, and per-package CRUD method scans.
type DocCache struct {
	files map[string][]string
	docs  map[string]string
	crud  map[string][]ClientMethodDoc
}

// NewDocCache returns an empty DocCache ready for use.
func NewDocCache() *DocCache {
	return &DocCache{
		files: map[string][]string{},
		docs:  map[string]string{},
		crud:  map[string][]ClientMethodDoc{},
	}
}

// ReadFile returns the lines of path, reading and caching them on first access.
func (d *DocCache) ReadFile(path string) ([]string, error) {
	if lines, ok := d.files[path]; ok {
		return lines, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	d.files[path] = lines

	return lines, nil
}

// DocCommentAt returns the contiguous //-comment block directly above the
// declaration at pos (joined by newlines), or "" if pos is unknown or the file
// cannot be read. Results are cached per position.
func (d *DocCache) DocCommentAt(pos token.Position) string {
	if pos.Filename == "" || pos.Line <= 0 {
		return ""
	}

	key := fmt.Sprintf("%s:%d", pos.Filename, pos.Line)
	if v, ok := d.docs[key]; ok {
		return v
	}

	doc := ""

	lines, err := d.ReadFile(pos.Filename)
	if err == nil {
		doc = DocCommentAbove(lines, pos.Line)
	}

	d.docs[key] = doc

	return doc
}

// ScanItemCRUD scans the *_client.go files in an SDK package directory and
// returns the item's exported CRUD methods together with their doc comments,
// sorted by method name. Results are cached per (directory, item).
func (d *DocCache) ScanItemCRUD(pkgDir, item string) []ClientMethodDoc {
	cacheKey := pkgDir + "\x00" + item
	if v, ok := d.crud[cacheKey]; ok {
		return v
	}

	var methods []ClientMethodDoc

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		d.crud[cacheKey] = nil

		return nil
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), "_client.go") {
			continue
		}

		lines, rErr := d.ReadFile(filepath.Join(pkgDir, e.Name()))
		if rErr != nil {
			continue
		}

		for i, line := range lines {
			_, name := ParseClientMethod(line)
			if name == "" || !IsItemCRUD(name, item) {
				continue
			}

			// Line numbers passed to DocCommentAbove are 1-indexed.
			methods = append(methods, ClientMethodDoc{Name: name, Doc: DocCommentAbove(lines, i+1)})
		}
	}

	slices.SortFunc(methods, func(a, b ClientMethodDoc) int { return strings.Compare(a.Name, b.Name) })

	d.crud[cacheKey] = methods

	return methods
}

// DocCommentAbove returns the contiguous //-style comment block that sits
// directly above the 1-indexed declLine, joined by newlines in source order.
func DocCommentAbove(lines []string, declLine int) string {
	// declLine is 1-indexed; the line directly above it is at slice index declLine-2.
	end := declLine - 2
	if end < 0 || end >= len(lines) || !strings.HasPrefix(strings.TrimSpace(lines[end]), "//") {
		return ""
	}

	// Walk up to the top of the contiguous //-comment block.
	start := end
	for start-1 >= 0 && strings.HasPrefix(strings.TrimSpace(lines[start-1]), "//") {
		start--
	}

	doc := make([]string, 0, end-start+1)
	for i := start; i <= end; i++ {
		doc = append(doc, strings.TrimSpace(lines[i]))
	}

	return strings.Join(doc, "\n")
}

// ExtractItemTypeInfoBool returns the value of the named boolean field inside the
// package's ItemTypeInfo composite literal, and whether the field was found with
// a literal true/false value.
func ExtractItemTypeInfoBool(pkg *packages.Package, field string) (bool, bool) { //nolint:revive // value + found
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}

			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || !HasName(vs.Names, "ItemTypeInfo") {
					continue
				}

				for _, val := range vs.Values {
					lit, ok := val.(*ast.CompositeLit)
					if !ok {
						continue
					}

					if v, ok := extractBoolField(lit, field); ok {
						return v, true
					}
				}
			}
		}
	}

	return false, false
}

func extractBoolField(lit *ast.CompositeLit, field string) (bool, bool) { //nolint:revive // value + found
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok || key.Name != field {
			continue
		}

		ident, ok := kv.Value.(*ast.Ident)
		if !ok || (ident.Name != "true" && ident.Name != "false") {
			continue
		}

		return ident.Name == "true", true
	}

	return false, false
}

// FindFabricItemType locates the FabricItemType const (= fabcore.ItemType<Name>)
// in the package and returns the SDK constant name plus a file in the fabcore
// package (used to derive the SDK fabric/ directory).
func FindFabricItemType(pkg *packages.Package) (string, string) { //nolint:revive // SDK const name + fabcore file path
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}

			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || !HasName(vs.Names, "FabricItemType") {
					continue
				}

				for _, val := range vs.Values {
					sel, ok := val.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					obj := pkg.TypesInfo.ObjectOf(sel.Sel)
					if obj == nil || obj.Pkg() == nil {
						continue
					}

					return obj.Name(), pkg.Fset.Position(obj.Pos()).Filename
				}
			}
		}
	}

	return "", ""
}

// FabricRoot returns the fabric-sdk-go "fabric" directory given a file inside the
// fabric/core package (e.g. .../fabric-sdk-go@v/fabric/core/x.go -> .../fabric).
func FabricRoot(coreFile string) string {
	dir := filepath.Dir(coreFile)
	if filepath.Base(dir) != "core" {
		return ""
	}

	return filepath.Dir(dir)
}

// SDKPackageDir resolves the package's FabricItemType constant to its dedicated
// fabric-sdk-go package directory, returning the directory and the item name.
// Overrides maps an item name to a non-default SDK package directory name when
// the default lowercase derivation does not match.
func SDKPackageDir(pkg *packages.Package, overrides map[string]string) (string, string, bool) { //nolint:revive // dir + item name + ok
	constName, anyFile := FindFabricItemType(pkg)
	if constName == "" || anyFile == "" {
		return "", "", false
	}

	// ItemType constants are named ItemType<Name>; the SDK package is fabric/<name>.
	item := strings.TrimPrefix(constName, "ItemType")
	if item == constName || item == "" {
		return "", "", false
	}

	fabricDir := FabricRoot(anyFile)
	if fabricDir == "" {
		return "", "", false
	}

	pkgName := strings.ToLower(item)
	if override, has := overrides[item]; has {
		pkgName = override
	}

	return filepath.Join(fabricDir, pkgName), item, true
}

// IsItemCRUD reports whether method is a CRUD operation on the item itself,
// e.g. GetNotebook, BeginCreateNotebook, ListNotebooks, UpdateNotebook.
func IsItemCRUD(method, item string) bool {
	if !strings.Contains(strings.ToLower(method), strings.ToLower(item)) {
		return false
	}

	for _, verb := range []string{"Get", "List", "Create", "Update", "Delete"} {
		if strings.Contains(method, verb) {
			return true
		}
	}

	return false
}

// HasName reports whether target appears among names.
func HasName(names []*ast.Ident, target string) bool {
	for _, n := range names {
		if n.Name == target {
			return true
		}
	}

	return false
}

// PkgName returns the package's name, falling back to the import-path base.
func PkgName(pkg *packages.Package) string {
	if pkg.Name != "" {
		return pkg.Name
	}

	return filepath.Base(pkg.PkgPath)
}
