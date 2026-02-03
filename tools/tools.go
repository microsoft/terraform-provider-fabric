// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

//go:generate go install github.com/go-task/task/v3/cmd/task@latest
//go:generate task tools
