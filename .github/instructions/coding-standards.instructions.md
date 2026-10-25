---
applyTo: "internal/**/*.go"
---

# Coding Standards — Go Files

## Copyright Header

Every `.go` file: `// Copyright Microsoft Corporation 2026` + `// SPDX-License-Identifier: MPL-2.0`

## SDK Import Aliases

`fab` + package name: `fabcore`, `fablakehouse`, `fabfake`

## Microsoft Link Rule

Never include `en-us` locale: `https://learn.microsoft.com/fabric/...`

## Always Use `MarkdownDescription`

Never use `Description` in schema attributes — lint fails.

## Error Constants

Use `common.Err*` from `internal/common/errors.go`: `ErrCreate`, `ErrRead`, `ErrUpdate`, `ErrDelete` for CRUD operation summaries; `ErrConfigRead` for provider config errors.

## Constructor Naming

`<pkg>.NewResource<Type>`, `<pkg>.NewDataSource<Type>`, `<pkg>.NewDataSource<Types>`
