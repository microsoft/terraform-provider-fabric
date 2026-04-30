---
applyTo: "internal/**/*.go"
---

# Coding Standards — Go Files

## Copyright Header

Every `.go` file: `// Copyright Microsoft Corporation 2026` + `// SPDX-License-Identifier: MPL-2.0`

## SDK Import Aliases

`fab` + package name: `fabcore`, `fablakehouse`, `fabfake`

## HCL Naming

SDK PascalCase → TF snake_case: `CapacityID` → `capacity_id`, `SQLEndpointProperties` → `sql_endpoint_properties`

## Microsoft Link Rule

Never include `en-us` locale: `https://learn.microsoft.com/fabric/...`

## Error Constants

Use `common.Err*` from `internal/common/errors.go` for error summaries.

## Always Use `MarkdownDescription`

Never use `Description` in schema attributes.

## Fabric Item File Structure

| File                                 | Purpose                                                              |
| ------------------------------------ | -------------------------------------------------------------------- |
| `base.go`                            | Constants: `FabricItemType`, `ItemTypeInfo`, `itemDefinitionFormats` |
| `resource_<type>.go`                 | Resource constructor and closures                                    |
| `data_<type>.go` / `data_<types>.go` | Singular / plural data source                                        |
| `models.go`                          | Model structs with `set()` methods                                   |
| `schema_resource_<type>.go`          | Resource schema attributes                                           |
| `schema_data_<type>.go`              | Data source schema attributes                                        |

Constructor naming: `<pkg>.NewResource<Type>`, `<pkg>.NewDataSource<Type>`, `<pkg>.NewDataSource<Types>`
