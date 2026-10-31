# Developer Guide <!-- omit in toc -->

- [Microsoft Fabric prerequisites](#microsoft-fabric-prerequisites)
- [Development Environment](#development-environment)
  - [DevContainer development](#devcontainer-development)
    - [Prerequisites](#prerequisites)
    - [Opening the DevContainer](#opening-the-devcontainer)
  - [Local development](#local-development)
    - [Requirements](#requirements)
    - [Linux](#linux)
    - [Windows](#windows)
    - [Linting](#linting)
      - [Markdown](#markdown)
    - [Provider Dev Overrides](#provider-dev-overrides)
- [Building](#building)
- [Debugging](#debugging)
  - [Debugging with VS Code](#debugging-with-vs-code)
- [Testing](#testing)
  - [Well-Known resources](#well-known-resources)
  - [Unit Tests](#unit-tests)
  - [Acceptance Tests](#acceptance-tests)
- [Dependencies](#dependencies)
  - [Adding](#adding)
  - [Updating](#updating)
- [Documentation](#documentation)
- [Release](#release)

---

The Terraform Provider for Microsoft Fabric extends Terraform's capabilities to allow Terraform to manage Microsoft Fabric infrastructure and services. The provider is built on the modern [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and NOT on the the older Terraform SDK. Ensure that you are referencing the correct [Plugin Framework documentation](https://developer.hashicorp.com/terraform/plugin/framework) when developing for this provider.

If you want to contribute to the provider, refer to the [Contributing Guide](./CONTRIBUTING.md) which can help you learn about the different types of contributions you can make to the repo. The following documentation will help developers get setup and prepared to make code contributions to the repo.

## Microsoft Fabric prerequisites

- Microsoft Entra ID Tenant with Global Admin access
- Microsoft Fabric with Admin access

You can request a free development Entra ID tenant via [Microsoft 365 Developer Program](https://developer.microsoft.com/microsoft-365/dev-program).
Using created tenant above, request a free Fabric trial license using [Get started with Microsoft Fabric](https://www.microsoft.com/microsoft-fabric/getting-started)
page.

## Development Environment

Depends on your preferences you can use pre-configured development environment with DevContainer or your machine directly.

### DevContainer development

DevContainer is the **preferred** way for contribution to the project.

The [DevContainer](https://containers.dev/) feature in Visual Studio Code creates a consistent and isolated development environment. A DevContainer is a Docker container that has all the tools and dependencies needed to work with the codebase. You can open any folder inside the container and use VS Code's full feature set, including IntelliSense, code navigation, debugging, and extensions.

#### Prerequisites

To use the DevContainer in this repo, you need to have the following prerequisites:

- [Docker](https://www.docker.com/products/docker-desktop/)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed in VS Code.

#### Opening the DevContainer

Once you have the prerequisites, you can follow these steps to open the repo in a DevContainer:

1. Clone or fork this repo to your local machine.
1. Open VS Code and press F1 to open the command palette. Type "Remote-Containers: Open Folder in Container..." and select it.
1. Browse to the folder where you cloned or forked the repo and click "Open".
1. VS Code will reload and start building the DevContainer image. This may take a few minutes depending on your network speed and the size of the image.
1. When the DevContainer is ready, you will see "Dev Container: Terraform Provider" in the lower left corner of the VS Code status bar. You can also open a new terminal (Ctrl+Shift+`) and see that you are inside the container.
1. You can now edit, run, debug, and test the code as if you were on your local machine. Any changes you make will be reflected in the container and in your local file system.

> [!NOTE]
> To work with the repository you will need to verify or configure your GIT credentials, you can do it as follows in the dev Container terminal:

- Verify Git user name and email:

```shell
git config --list
```

You should see your username and email listed, if they do not appear or you want to change them you must establish them following the step below, (to quit the "git config" mode type "q").

- Change or set your Git username and email in the DevContainer:

```shell
git config --global user.name "Your Name"
git config --global user.email "your.email@address"
```

> [!NOTE]
> If you logging to docker container's shell outside the VS Code, in order to work with git repository, run the following commands:

```shell
export SSH_AUTH_SOCK=$(ls -t /tmp/vscode-ssh-auth* | head -1)
export REMOTE_CONTAINERS_IPC=$(ls -t /tmp/vscode-remote-containers-ipc* | head -1)
```

For more information about DevContainers, you can check out the [DevContainer documentation](https://code.visualstudio.com/docs/devcontainers/containers) and [sharing Git credentials with your container](https://code.visualstudio.com/remote/advancedcontainers/sharing-git-credentials).

### Local development

DevContainer is the **preferred** way for contribution to the project. It contains all necessary tools and configuration that is ready to start corking on the code without thinking on development environment setup.

Local development is still possible on Windows, Linux and macOS, but requires additional step to setup development environment.

> [!NOTE]
> Treat all instructions, commands or scripts in `Local development` section as examples. Depending on your local environment and configuration, the final commands or script may vary.

#### Requirements

- [Git](https://git-scm.com/downloads) `>= 2.47.1`
- [Go](https://go.dev/doc/install) `>= 1.24.4`
  - We recommend you to use Go version manager [go-nv/goenv](https://github.com/go-nv/goenv/blob/master/INSTALL.md)
    - `goenv install 1.24.4`
- [Terraform](https://developer.hashicorp.com/terraform/downloads) `>= 1.12.2`
  - We recommend you to use Terraform version manager [tfutils/tfenv](https://github.com/tfutils/tfenv/blob/master/README.md)
    - `tfenv install 1.12.2`, `tfenv use 1.12.2`
- [Task](https://taskfile.dev/installation) `>= 3.40.1`

#### Linux

- [Git](https://git-scm.com/downloads)
- [Go](https://go.dev/doc/install)
- [Terraform](https://developer.hashicorp.com/terraform/install)
- [Task](https://taskfile.dev/installation)

Below you can find examples of tools setup for Ubuntu/Debian.

Install Git

```shell
sudo apt update && sudo apt install git
```

Install Go

```shell
goVersion=$(curl https://go.dev/dl/?mode=json | jq -r '.[0].version')
curl -LO https://go.dev/dl/$goVersion.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local/ -xzf $goVersion.linux-amd64.tar.gz
rm -f $goVersion.linux-amd64.tar.gz
echo 'export PATH="$PATH:/usr/local/go/bin"' | sudo tee /etc/profile.d/go-lang.sh >/dev/null
echo 'export GOROOT=/usr/local/go' >>~/.bashrc
echo 'export GOPATH=$HOME/go' >>~/.bashrc
echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >>~/.bashrc
source ~/.bashrc
```

Install Terraform

```shell
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install terraform
```

Install Task

```shell
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
```

#### Windows

- [Git](https://gitforwindows.org)
- [Go](https://go.dev/doc/install)
- [Terraform](https://developer.hashicorp.com/terraform/install#windows)
- [Task](https://taskfile.dev/installation)

> For _Git for Windows_, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".
> _Git for Windows_ must be installed per steps above

Install via [winget](https://learn.microsoft.com/windows/package-manager/winget/#install-winget)

```powershell
winget install GoLang.Go Hashicorp.Terraform Task.Task
```

or via [Chocolatey](https://chocolatey.org/install)

```powershell
choco install golang terraform go-task -y
refreshenv
```

#### Linting

##### Markdown

To lint Markdown files use `markdownlint` tools. For the installation you can run:

```shell
task install:markdownlint
```

You can integrate `markdownlint` with vscode by using extension: [DavidAnson.vscode-markdownlint](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint)

To execute linter, run:

```shell
task lint:md
```

> [!NOTE]
> Markdown linter runs automatically just after `task docs` as well to lint autogenerated files.

#### Provider Dev Overrides

> [!NOTE]
> You can omit below configuration if you work under DevContainer. It uses pre-configured `terraformrc`.

With Terraform v0.14 and later, [development overrides for provider developers](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) can be leveraged in order to use the provider built from source.

> [!WARNING]
> If you already use `terraformrc`, be careful not to overwrite your existing settings and only add the necessary configuration for this project to your existing ones.

To do this, populate a Terraform CLI configuration file `~/.terraformrc` for all platforms **other** than Windows:

```hcl
provider_installation {
  dev_overrides {
    "microsoft/fabric" = "/home/{REPLACE WITH YOUR SOURCE CODE PATH}/terraform-provider-fabric/bin/{REPLACE WITH darwin FOR MACOS OR linux FOR LINUX}-amd64"
  }

  # Install all other providers directly from their origin provider
  # registries as normal. If you omit this, no other providers will be available.
  direct {}
}
```

For Windows use `terraform.rc` in the `%APPDATA%` directory:

```hcl
provider_installation {
  dev_overrides {
    "microsoft/fabric" = "C:\\Users\\{YOUR USERNAME}\\{REPLACE WITH YOUR SOURCE CODE PATH}\\terraform-provider-fabric\\bin\\windows-amd64"
  }

  # Install all other providers directly from their origin provider
  # registries as normal. If you omit this, no other providers will be available.
  direct {}
}
```

## Building

To compile the provider, run `task deps` and `task build`. This will build the provider and put the provider binary in the `bin` directory.

```shell
task deps
task build
```

## Debugging

This provider support [terraform plugin debugging](https://developer.hashicorp.com/terraform/plugin/debugging) pattern.

### Debugging with VS Code

1. Open VS Code with the source code root folder
1. Launch ["Run and Debug"](https://code.visualstudio.com/docs/editor/debugging) from VS Code.

  When the debugger start you should see the Terraform debugging information similar to below:

  ```text
  Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

   Command Prompt: set "TF_REATTACH_PROVIDERS={"registry.terraform.io/microsoft/fabric":{"Protocol":"grpc","ProtocolVersion":6,"Pid":69004,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:56897"}}}"
   PowerShell: $env:TF_REATTACH_PROVIDERS='{"registry.terraform.io/microsoft/fabric":{"Protocol":"grpc","ProtocolVersion":6,"Pid":69004,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:56897"}}}'
  ```

1. Copy `TF_REATTACH_PROVIDERS` value from the Debug Console
1. Set `TF_REATTACH_PROVIDERS` environment variable in the terminal with the value copied from the above step

  ```shell
  # Linux / macOS
  export TF_REATTACH_PROVIDERS={...}
  ```

  ```powershell
  # Windows
  $env:TF_REATTACH_PROVIDERS='{...}'
  ```

1. Add breakpoints
1. `cd` to a parent folder where `main.tf` exists
1. Run `terraform` commands

## Testing

### Well-Known resources

To run tests, especially "Acceptance Tests" some resources on the Fabric, Entra and Azure DevOps side have to be pre-created first. To setup them, set input environment variables:

```text
# Required
FABRIC_TESTACC_WELLKNOWN_ENTRA_TENANT_ID="<ENTRA TENANT ID>"
FABRIC_TESTACC_WELLKNOWN_AZURE_SUBSCRIPTION_ID="<AZURE SUBSCRIPTION ID>"
FABRIC_TESTACC_WELLKNOWN_FABRIC_CAPACITY_NAME="<FABRIC CAPACITY NAME>"
FABRIC_TESTACC_WELLKNOWN_AZDO_ORGANIZATION_NAME="<AZURE DEVOPS ORGANIZATION NAME>"
FABRIC_TESTACC_WELLKNOWN_NAME_PREFIX="<RESOURCES PREFIX>"
FABRIC_TESTACC_WELLKNOWN_AZURE_RESOURCE_GROUP_NAME="<AZURE RESOURCE GROUP NAME>"
FABRIC_TESTACC_WELLKNOWN_AZURE_LOCATION="<AZURE LOCATION>"

# Optional
FABRIC_TESTACC_WELLKNOWN_NAME_SUFFIX=""
FABRIC_TESTACC_WELLKNOWN_NAME_BASE=""
```

You can set those variables into `./wellknown.env` files as well.

Then run:

```shell
task testacc:setup
```

Setup writes `internal/testhelp/fixtures/.wellknown.json`. A populated
`FABRIC_TESTACC_WELLKNOWN` environment variable takes precedence over that file.
GitHub workflows use the `FABRIC_TESTACC_WELLKNOWN` Actions variable; a local
setup run does not publish or refresh that variable.

### Unit Tests

> [!NOTE]
> Unit tests won't create the actual resources since they will be run against a fake server.

To run all unit tests

```shell
task testunit
```

To run single unit test

```shell
task testunit -- <test_name>
```

> [!NOTE]
> The tests require permissions on the folders, these permissions are assigned when creating your container. If you have permission problems when running the unit tests, you can rebuild your development container or run the following commands again to assign the permissions to the necessary folders.

```shell
sudo chown -R vscode /workspace
sudo chown -R vscode /go/pkg
```

### Acceptance Tests

> [!NOTE]
> Acceptance tests will create the actual resources since they will be run against real APIs.

The MLModel and MLExperiment CRUD tests create their own folders in `WorkspaceRS`,
use Terraform references to order item creation and moves, and destroy their owned
resources afterwards. They do not require pre-created `FolderRS1`/`FolderRS2`.

With the executed workflows retaining `fabric-testacc-environment` and
`cancel-in-progress: false`, Test and Cleanup share a **repository-wide,
workflow-level** slot. It lasts until the whole workflow exits, including
docs-only/build jobs, environment-approval waits and the test matrix.
GitHub's default single pending slot can be replaced by a newer run; dispatch
ordering is not guaranteed. Both workflows schedule at 02:00 UTC, so pending
Cleanup can be superseded, visible as a cancelled Cleanup run. Re-run a superseded
required check after the group is **idle: no in-progress run and no pending run**.
An environment-approval wait is in-progress even when no test job is executing.
A Cleanup dispatch is a newer group run and cancels/replaces a pending required
Test run; owner coordination does not replace the idle check.

`pull_request` uses the PR merge-ref workflow. An incoming PR workflow can rename
or remove the group, or set `cancel-in-progress: true` and cancel an active run;
the guarantees and recovery instructions here do not cover those edits. Upstream fork-PR runs participate
only while their executed file retains the group and false setting. Review workflow
changes before granting any required outside-collaborator **workflow approval**.
Separate **job environment approval** gates jobs and environment secrets; it does
not protect workflow-level concurrency. Fork-owned workflows, local commands,
other repositories and manual cleanup remain outside this lock.

The in-workflow **pre-clean** is best-effort and removes recoverable items only.
**Terraform destroy** runs inside each test job; standalone **Cleanup** can delete
active items and non-preserved folders. Cancellation, timeout or runner loss can skip destroy,
leaving active objects the next pre-clean will not reclaim. After interruption or
pending-Cleanup replacement, coordinate recovery with the environment owner.
No cancellation or rejection atomically clears the group or rolls back resource
changes; queued runs can be promoted and new runs can arrive.

1. Inspect the queue and cancel unwanted pending runs **before slow workspace
   assessment**, including Cleanup that must not start. Recheck before every
   controllable or expected exit, not only rejection/cancellation: natural completion, timeout, runner loss or
   rejection of Cleanup's only job can also promote pending work immediately.
   Promotion is neither approval nor idle; promoted jobs may wait for a new
   environment review or start at once. Keep observing until the group is idle. If
   unwanted Cleanup has already started, cancel it and assess changes; deletion
   is not undone. This ordering reduces risk but is not an atomic guarantee.
2. Before approving any waiting Cleanup job, review the **workflow and referenced
   script at that run's actual executing revision/SHA**, including the PR merge
   revision where applicable, not a head-branch or newer branch tip. Check
   target/preservation behavior, `development`, the exact script/mode invocation,
   the group and false setting. Standalone `cleanup-script.yml` must use
   `Remove-WorkspaceRSItems.ps1 -Force` without `-RecoverableItemsOnly`; Test
   pre-clean must use `-Force -RecoverableItemsOnly`. Any other invocation is
   changed execution, not approved. For a new `workflow_dispatch`, review the selected
   ref's resolved revision; do not dispatch unreviewed/moving work on the assumption
   that an approval prompt will stop it. Unknown or changed execution is not
   approved. Distinguish standalone `cleanup-script.yml` (`-Force`) from Test's
   best-effort pre-clean (`-Force -RecoverableItemsOnly`); approval starts the
   intended job, and standalone Cleanup then hard-deletes.
3. Rejection prevents that job starting but siblings may continue; for standalone
   Cleanup, rejecting its only job can end the run and promote pending work.
   Cancellation ends the selected run and frees its running slot, which can
   immediately promote pending work; rejection of Cleanup's only job can do the
   same, while Test siblings may continue after one rejection. Promotion is not
   approval: a later environment gate may or may not wait and is not a barrier to
   rely on. Neither action atomically makes the group idle or rolls back changes.
   Job timeouts do not bound approval waits.
   Resolve queued/waiting work with the owner and recheck until the group is idle.
4. Do not dispatch unneeded Cleanup; cancel any unwanted pending or approval-waiting
   Cleanup. Dispatch approved Cleanup only in an observed idle group, continuing
   to monitor arrivals and queued work through run exit and inventory assessment.
   These manual observations do not enforce a post-run inventory barrier: if
   unsafe queued work cannot be held back, stop normal recovery and coordinate
   an owner-reviewed plan rather than claim an atomic lock.
   After standalone `cleanup-script.yml` exits, a `success` conclusion or
   "deleted successfully" summary is **inconclusive**: the script has non-fatal
   fetch/delete errors. Require the `Found`/`No items found`, `Found`/`No folders
   found`, and `Found`/`No recoverable items found` results for all three fetches;
   any `Failed to fetch` message is unsuccessful regardless of WARN/ERROR level
   or later success text. When a deletion summary exists, all three failed counts
   must be zero. A zero-delete early return is acceptable only with complete fetch
   evidence and independent inventory of active items, folders and recoverable
   items matching the expected retained state from the pre-run assessment.
   After every run, independently inventory all three resource categories against
   the expected deleted/preserved state recorded before execution; do not infer it
   from counts or success text. Script totals and `only preserved items
   found` are not that inventory. The expected `-Force` banner and rate-limit retry
   warnings alone are neither success nor failure. Missing required evidence,
   fetch/delete failures, inventory mismatch, timeout or cancellation means
   unsuccessful Cleanup; stop for an owner-reviewed partial-state plan,
   without automatic Cleanup retry or Test start. This standalone gate does not
   turn Test's best-effort recoverable pre-clean into a new hard gate.
   Recheck idle immediately before a required Test re-run and monitor until it
   becomes the intended in-progress run, without concurrent Cleanup. Re-run checks
   affected by supersession, failure, cancellation or rejection only after the
   owner confirms the workspace and an idle group, in a separate window from Cleanup.
Do not run manual or standalone cleanup alongside acceptance tests.

Existing cleanup-job limits are 30 minutes and test-matrix limits are 100 minutes.
Job limits do not bound concurrency-queue or environment-approval waits; jobs
without explicit limits retain platform defaults. Matrix legs still run in parallel
in `WorkspaceRS`, and non-ML acceptance tests still share `FolderRS1`/`FolderRS2`.
Additional job bounds and workspace/stable-identity lifecycle isolation are
follow-ups, not guarantees supplied by this change.

To run all acceptance tests

```shell
task testacc
```

To run single acceptance test

```shell
task testacc -- <test_name>
```

## Dependencies

### Adding

This provider uses [Go modules](https://go.dev/wiki/Modules). Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to the provider:

```shell
go get github.com/author/dependency
task deps
```

### Updating

In the terminal type and run:

```sh
task deps
```

## Documentation

User documentation markdown files in [./docs](./docs/) are auto-generated by the [terraform plugin docs tool](https://github.com/hashicorp/terraform-plugin-docs).

**DO NOT** manually edit the markdown files in [./docs](./docs/). If you need to edit documentation edit the following sources:

- schema information in the provider, resource, and data-source golang files that are in [./internal/services](./internal/services)
- [template files](./templates/)

```sh
task docs
```

User documentation is temporarily served on GitHub Pages which requires the [ghpages.yml GitHub workflow](./.github/workflows/ghpages.yml) to transform [./docs](./docs/) markdown files into a static website. Once this provider is published to the Terraform registry, documentation will be hosted on the registry instead.

## Release

Our releases use [semantic versioning](https://semver.org/).

Given a version number MAJOR.MINOR.PATCH, increment the:

- MAJOR version when you make incompatible changes
- MINOR version when you add functionality in a backward compatible manner
- PATCH version when you make backward compatible bug fixes

Use the `beta` suffix with identifier to the MAJOR.MINOR.PATCH format for preview release such as `v0.1.0-beta.1`.
