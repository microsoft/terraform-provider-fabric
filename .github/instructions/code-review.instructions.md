---
applyTo: "internal/**/*.go"
---

# Go Code Review Checklist — terraform-provider-fabric

## Purpose & Scope

Guidance for Copilot code review on Go files under `internal/`. Flag deviations
from the conventions below. These rules complement (not repeat) the authoring
instructions and are scoped to review only.

## File Hygiene

- Every `.go` file must start with the two-line header:
  `// Copyright Microsoft Corporation 2026` then
  `// SPDX-License-Identifier: MPL-2.0`. Flag files missing either line.
- Flag Microsoft docs links containing an `en-us` (or any) locale segment.

```go
// Avoid
// https://learn.microsoft.com/en-us/fabric/...
// Prefer
// https://learn.microsoft.com/fabric/...
```

## Schema

- Flag `Description:` in schema attributes — this repo requires
  `MarkdownDescription:` and `Description` fails lint.

```go
// Avoid
Description: "The workspace ID."
// Prefer
MarkdownDescription: "The workspace ID."
```

## Descriptions & Docs Text

Description defects propagate into generated `docs/`. Flag:

- Copy/paste descriptions referencing the wrong subject/attribute (e.g. an
  "inbound" description on an `outbound` attribute, or another resource's name).
- Attribute names in prose using camelCase instead of the snake_case HCL name
  (`passwordReference` → `password_reference`).
- Missing space when concatenating enum lists, e.g. `"Possible values:" + ...`
  rendering as `Possible values:`Active``. Add a trailing space/newline.
- `en-us` (or any locale) in `learn.microsoft.com` URLs in `DocsURL` too.
- Non-tenant-specific/placeholder-unfriendly example URLs presented as required.

## Schema Validators

- Write-only secret attributes must pair with conditional validation. When a
  `*_wo` value is optional (to allow a `*_reference`), flag a missing
  `AlsoRequires` linking it to its `*_wo_version` (and vice versa).

```go
stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("password_wo_version"))
```

- Definition path-key validators: flag exact `OneOf`-style matching when the
  format declares wildcard parts (e.g. `EntityTypes/*`, `Files/.../*`); use the
  pattern-based validator instead.
- Flag `mapvalidator.SizeAtMost(1)` (or too-small limits) on `definition` when
  the format legitimately has multiple parts.
- Flag `customtypes.UUIDType{}` in schema whose model field is `types.String`
  (or any schema custom type / model type mismatch) — causes decode errors.

## State & Diagnostics

- Append and check operation diagnostics **before** writing state. Flag
  `resp.State.Set(...)` (or `RemoveResource`) called before returning on a prior
  error (e.g. tag sync, `SyncTags`) — a failed apply must not persist new state.
- Flag ignored diagnostics from `model.set(...)` / property setters.
- `Delete` must call `resp.State.RemoveResource(ctx)` on success, including
  early-return paths.

## Preview Mode Parity

- If the resource's `Configure`/CRUD calls `fabricitem.IsPreviewMode(...)`, the
  matching data source must too. Flag data sources missing the preview check.
- Flag hard-coded `false` for a preview/GA flag when the type carries
  `TypeInfo.IsPreview` — pass the real value.

## Errors

- Flag raw/inline error summary strings in CRUD paths. Use the shared
  `common.Error*` constants (`ErrorCreateHeader`, `ErrorReadHeader`,
  `ErrorUpdateHeader`, `ErrorDeleteHeader`, `ErrorListHeader`, and their
  `*Details` counterparts).

```go
// Avoid
resp.Diagnostics.AddError("create failed", err.Error())
// Prefer
resp.Diagnostics.AddError(common.ErrorCreateHeader, ...)
```

## Naming

- SDK import aliases must be `fab` + package name (e.g. `fabcore`,
  `fablakehouse`, `fabfake`). Flag other aliases.
- HCL attribute names must be snake_case mapping from SDK PascalCase
  (`CapacityID` → `capacity_id`). Flag camelCase/PascalCase HCL names.
- Constructors must follow `New<Kind><Type>` (e.g. `NewResourceLakehouse`,
  `NewDataSourceLakehouse`).

## Provider Registration

- A new resource/data source must be registered in
  `internal/provider/provider.go` (import in alphabetical order and added to
  `Resources()` / `DataSources()`). Flag new types that are not registered.
- New/changed resources should ship a matching example under `examples/` and
  regenerated docs (`task docs`). If this change adds or renames a resource,
  flag when a corresponding `examples/` entry is absent from the PR.

## Tests

- Test files must use the black-box `<pkg>_test` package (e.g.
  `package lakehouse_test`). Flag white-box test packages.
- Prefer `resource.ParallelTest` unless tests have ordered dependencies. Flag
  `ParallelTest` when the test mutates the shared `fakes.FakeServer` (via
  `Upsert`, etc.) — the fake server is unsynchronized and this races.
- Tests must assert the **new** behavior, not just generic attributes. Flag new
  definition formats / job-type combos / attributes added without a case that
  exercises and checks them.
- Test function names follow `TestUnit_<Type>_*` / `TestAcc_<Type>_*`
  (e.g. `TestAcc_WorkspaceNetworkCommunicationPolicyResource_CRUD`).
- Flag any suggestion or doc that runs `go test` directly; this repo uses
  `task testunit` / `task testacc` (required env like `FABRIC_PREVIEW=true`).
- New code should target >80% coverage; flag new exported behavior lacking tests.

## General

- Flag use of deprecated SDK calls or deprecated libraries.
- Flag missing error handling on SDK/network calls and unhandled `err` values.
- Flag dereferencing SDK pointers (e.g. `Etag`) or indexing slices
  (`Value[0]`) without a nil / length check.
- Flag error messages naming the wrong operation (e.g. "delete" text in
  `Create`) or the wrong entity (e.g. "Workspace" in a Tag data source).
- Flag validators/helpers that are implemented but whose invocation is
  commented out, and dead/unused code (fails golangci-lint).
- Flag `Update` request builders that omit a newly added updatable attribute.
- Flag secrets, tokens, or connection strings committed in code or test fixtures.
