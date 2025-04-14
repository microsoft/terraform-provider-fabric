# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- `task build` - Build development provider binary
- `task lint` - Run all linters (files, Go, Terraform, Markdown)
- `task testunit` - Run unit tests
- `task testunit -- [TestName]` - Run specific unit test (e.g., `task testunit -- WorkspaceResource_CRUD`)
- `task testacc` - Run acceptance tests 
- `task testacc -- [TestName]` - Run specific acceptance test

## Code Style
- Files must start with copyright header (`// Copyright (c) Microsoft Corporation\n// SPDX-License-Identifier: MPL-2.0`)
- Follow Go black-box testing patterns with tests in `*_test.go` files
- Use `TestUnit_` prefix for unit tests and `TestAcc_` for acceptance tests
- Use snake_case for Terraform HCL attributes (e.g., `CapacityID` → `capacity_id`)
- Fabric SDK imports use aliases to avoid collisions (e.g., `fabcore` for SDK package `core`)
- Service constructors follow pattern: `New[DataSource|Resource][ServiceName]`
- Package directory structure follows: `internal/services/[service]` with consistent file naming
- Error titles come from constants in `internal/common/errors.go`
- Microsoft links should not contain language identifiers (like `en-us`)