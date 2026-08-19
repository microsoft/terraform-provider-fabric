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

Use the `common.Error*` constants from `internal/common/errors.go`: `ErrorCreateHeader`, `ErrorReadHeader`, `ErrorUpdateHeader`, `ErrorDeleteHeader`, `ErrorListHeader` (and their `*Details` counterparts) for CRUD operation summaries; `ErrorInvalidConfig` for provider config errors.

## Constructor Naming

`<pkg>.NewResource<Type>`, `<pkg>.NewDataSource<Type>`, `<pkg>.NewDataSource<Types>`
