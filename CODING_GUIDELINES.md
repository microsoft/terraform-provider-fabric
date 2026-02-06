# 🧑‍💻 Coding Guidelines

- Terraform HCL Attribute naming uses lowered snake cases based on SDK parameters name, i.e. `CapacityID` -> `capacity_id`, `DisplayName` -> `display_name`, etc

- Resource and DataSources are stored under `internal/services` folder. There should be a package name that encompass a fabric item (such as workspace).

    i.e. `workspace` use own folder called `workspace` and package name is `workspace` as well

- Basic file structure per package is

  - `data_<service>.go` or `data_<service>_<operation>.go` for data-sources
    - i.e. `data_workspace.go`
  - `resource_<service>.go` or `resource_<service>_<operation>.go` for resources
    - i.e. `resource_workspace.go`
    - i.e. `resource_workspace_role_assignment.go`
  - `models.go` contains reusable model across all data-sources/resources with DTO conversion helpers if needed see [Workspace models.go](./internal/services/workspace/models.go)

- Service usage in provider uses the constructor methods `NewDataSource<Service>`, `NewResource<Service>` prefixed with service package name from imports, i.e. `workspace.NewDataSourceWorkspace`, `notebook.NewResourceNotebook`

- Common error titles are part of constants – it enforces error naming consistency.This is used for Error Summaries, not error details.  Errors are in `./internal/common/errors.go`

- Fabric SDK imports use aliases to avoid naming collisions with similar packages. Alias starts with fab + package name, i.e. alias for SDK package `core` is `fabcore`. The same rule applies to SDK fakes, like `fabfake`

- If tests are not dependent on each other utilize `ParallelTest` otherwise just `Test`

- Unit tests must start with `TestUnit_` prefix following with the DataSource/Resource name and a description that indicates test case, i.e. `TestUnit_WorkspaceDataSource_CRUD`. The same rule applies to the acceptance test, but the prefix is `TestAcc_`

- Unit Tests utilize Fakes for mocking responses
- Acceptance Tests execute the real API call to Fabric

- run Unit tests command: `task testunit`
- Execute a single Unit test or group. For example:

  - `task testunit -- WorkspaceResource`
  - `task testunit -- WorkspaceResource_CRUD`

    The designator after `testunit` is a part of the test name, specifically following `TestUnit_`. For example if a Test is called `TestUnit_WorkspaceResource_CRUD` the designator is `WorkspaceResource_CRUD` for this specific test and group designator is `WorkspaceResource`

- Run Acceptance tests command: `task testacc`
- Execute a single Acceptance test or group. For example:

  - `task testacc -- WorkspaceResource`
  - `task testacc -- WorkspaceResource_CRUD`

    The designator after `testunit` is a part of the test name, specifically following `TestAcc_`. For example if a Test is called `TestAcc_WorkspaceResource_CRUD` the designator is `WorkspaceResource_CRUD` for this specific test and group designator is `WorkspaceResource`

- All provider services (data-sources/resources) should use use a `black-box` test pattern: <https://pkg.go.dev/testing> and should be placed right next to the file being tested.

    Tests are prefixed with service name, following with TF type and test suffix, i.e. `data_workspace_test.go`, `resource_workspace_test.go`

- Make sure that Microsoft links do not contain any language or region identifiers such as `en-us`. This includes links in markdown files and also in .go files with `MarkdownDescription` properties.

- Each .go file must start with header:

```go
// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0
```
