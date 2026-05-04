# Terraform Provider for Microsoft Fabric

This repository is **terraform-provider-fabric** — the official [Terraform](https://www.terraform.io/) provider for [Microsoft Fabric](https://learn.microsoft.com/fabric/).

## Technology Stack

- **Language:** Go (see `go.mod` for version)
- **Framework:** [HashiCorp Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) (NOT the older Terraform SDK)
- **Fabric SDK:** [Microsoft Fabric Go SDK](https://github.com/microsoft/fabric-sdk-go) (`github.com/microsoft/fabric-sdk-go`)
- **Task Runner:** [Task](https://taskfile.dev/) — all build, test, lint, and doc commands are defined in `Taskfile.yml`
- **Linters:** `golangci-lint`, `markdownlint`, `tflint`
- **Doc Generator:** `tfplugindocs`

## Environment Setup

- **Preferred:** DevContainer (see `.devcontainer/`) — pre-configured with all tools
- **Local:** Go, Task, PowerShell Core (`pwsh`), Terraform CLI. See `DEVELOPER.md` for full requirements.
- **Install Task:** `winget install Task.Task` (Windows) or see `Taskfile.yml` header for Linux/macOS
- **Install tools:** `task tools` — installs `golangci-lint`, `tfplugindocs`, `gotestsum`, and other dev tools
- **Provider dev overrides:** `task dev-overrides` — sets up `dev.tfrc` for local testing

## Project Structure

| Directory                  | Purpose                                                                                 |
| -------------------------- | --------------------------------------------------------------------------------------- |
| `internal/services/`       | All resource and data source implementations, one package per Fabric item or service    |
| `internal/pkg/fabricitem/` | Generic abstraction layer for Fabric Item resources (~80% of resources use this)        |
| `internal/pkg/utils/`      | Shared utility functions (`IsErrNotFound`, `GetDiagsFromError`, enum converters)        |
| `internal/provider/`       | Provider definition, configuration, and registration of all resources/data sources      |
| `internal/common/`         | Shared error constants (`common.Err*`), warnings, and base models                       |
| `internal/framework/`      | Custom types (`customtypes.UUID`, `customtypes.URL`), validators, and plan modifiers    |
| `internal/auth/`           | Authentication: credential types, OIDC, token acquisition                               |
| `internal/testhelp/`       | Test utilities, well-known fixtures, and `fakes/` directory for unit test fake handlers |
| `tools/itemgen/`           | Code generator to scaffold new Fabric Item service packages                             |
| `tools/scripts/`           | PowerShell scripts for well-known resource setup and maintenance                        |
| `examples/`                | Terraform HCL examples for each resource and data source                                |
| `docs/`                    | Auto-generated provider documentation (do NOT edit manually)                            |
| `templates/`               | Go templates for documentation generation (`tfplugindocs`)                              |

## Resource Categories

1. **Fabric Items (~80%)** — Use the generic `internal/pkg/fabricitem/` abstraction. Scaffold with `go run tools/itemgen/main.go`. Canonical reference: `internal/services/lakehouse/`
2. **Non-Item resources (~20%)** — Custom CRUD implementations using `superschema`. Examples: Workspace, Gateway, Connection, Domain, Shortcut. References: `internal/services/workspace/`, `internal/services/gateway/`, `internal/services/connection/`
3. **Sub-resources** — Scoped under a parent resource (role assignments, workspace settings, spark settings, tags, folders, schedulers). References: `internal/services/workspacera/`, `internal/services/sparkcustompool/`, `internal/services/tags/`

## Common Commands (Task Runner)

| Command                         | Purpose                                       |
| ------------------------------- | --------------------------------------------- |
| `task build`                    | Build development provider binary             |
| `task tools`                    | Install dev tools (linters, docgen)           |
| `task testunit -- <Name> <Pkg>` | Run unit tests matching name in package       |
| `task testacc -- <Name> <Pkg>`  | Run acceptance tests matching name in package |
| `task lint`                     | Run all linters (Go, TF, Markdown)            |
| `task docs`                     | Auto-generate provider documentation          |
| `task deps:up`                  | Update Go dependencies                        |

**NEVER run `go test` directly.** Always use `task testunit` or `task testacc`. The Taskfile sets required environment variables (e.g., `FABRIC_PREVIEW=true`) that are missing from raw `go test`, causing preview resource tests to fail.

Test name designator is the portion after `TestUnit_` or `TestAcc_`. The optional second argument scopes to a specific package for faster execution. For example, `task testunit -- WorkspaceResource ./internal/services/workspace/` runs all tests matching `TestUnit_WorkspaceResource*` in that package only. Without the package path, it defaults to `./...` (all packages).

## Provider Registration

New resources and data sources must be registered in `internal/provider/provider.go`:

- Import the service package (alphabetical order)
- Add constructor to `Resources()` and/or `DataSources()`
- If the constructor requires `ctx` (schema uses `supertypes`), wrap it: `func() resource.Resource { return pkg.NewResourceType(ctx) }`

## Documentation

- Schema descriptions via `MarkdownDescription` in schema attributes (never `Description`)
- Examples in `examples/` directory following `tfplugindocs` conventions
- Auto-generate docs with `task docs` — this reads schema + examples and writes to `docs/`
- Doc templates in `templates/` — for guides, index, and custom content

## Coding Conventions

- **SDK import aliases:** `fab` + package name (e.g. `fabcore`, `fablakehouse`, `fabfake`)
- **HCL naming:** SDK PascalCase → Terraform snake_case (`CapacityID` → `capacity_id`)
- **Error constants:** Use `common.Err*` from `internal/common/errors.go`
- **Black-box tests:** Test packages use `_test` suffix (e.g. `package lakehouse_test`)
- **Parallelism:** Use `resource.ParallelTest` unless tests have ordered dependencies
- **Coverage target:** >80% for new contributions
