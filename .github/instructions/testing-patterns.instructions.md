---
applyTo: "internal/**/*_test.go"
---

# Testing Patterns

## Test Types

| Suffix                                     | Purpose                      | Key Pattern                                                                |
| ------------------------------------------ | ---------------------------- | -------------------------------------------------------------------------- |
| `_Attributes`                              | Schema constraint validation | Missing required args, invalid UUIDs, unexpected/conflicting attributes    |
| `_CRUD`                                    | Full lifecycle               | Create+Read, Update+Read, Delete (automatic). Setup fake entities first.   |
| `_ImportState`                             | Import validation            | Invalid format, invalid segments, successful `workspaceID/entityID` import |
| DataSource                                 | Read scenarios               | Read-by-id, read-by-name, not-found for each                               |
| `_CRUD_Configuration` / `_CRUD_Definition` | Config/definition lifecycle  | For items with configuration or definitions                                |

## Naming Convention

- **Unit:** `TestUnit_<TypeName>Resource_CRUD`, `TestUnit_<TypeName>Resource_Attributes`, `TestUnit_<TypeName>DataSource`
- **Acceptance:** `TestAcc_<TypeName>Resource_CRUD`, `TestAcc_<TypeName>DataSource`

Group = designator after prefix. `TestUnit_WorkspaceResource_CRUD` → group `WorkspaceResource`.

## Black-Box Testing & File Setup

Test packages **must** use `_test` suffix (e.g. `package lakehouse_test`).

Each test file references shared variables from `base_test.go`:

```go
var testResourceItemFQN, testResourceItemHeader = testhelp.TFResource(common.ProviderTypeName, itemTypeInfo.Type, "test")
var testDataSourceItemFQN, testDataSourceItemHeader = testhelp.TFDataSource(common.ProviderTypeName, itemTypeInfo.Type, "test")
```

## Test Helpers

- `at.CompileConfig(header, map[string]any{...})` — compile HCL from header + attributes
- `at.JoinConfigs(...)` — join multiple config blocks
- `resource.ComposeAggregateTestCheckFunc(...)` — combine attribute checks
- `resource.TestCheckResourceAttr(...)` / `TestCheckResourceAttrPtr(...)` / `TestCheckResourceAttrSet(...)` / `TestCheckNoResourceAttr(...)`

## Parallelism & Coverage

Use `resource.ParallelTest(t, ...)` unless tests have ordered dependencies. Target **>80%** coverage.

## Well-Known Resources & Fake Handlers

- **Acceptance tests:** `testhelp.WellKnown()["WorkspaceRS"]` (resource tests) / `["WorkspaceDS"]` (data source tests)
- **Unit tests:** `fakes.FakeServer.Upsert(entity)` + `testhelp.NewTestUnitCase(t, &fqn, fakes.FakeServer.ServerFactory, nil, steps)`

## Run Commands

**CRITICAL: NEVER run `go test` directly.** Always use the Task runner — it sets required environment variables (e.g., `FABRIC_PREVIEW=true`) that are missing from raw `go test`. Tests for preview resources WILL FAIL without these variables.

- `task testunit -- <Group> <Pkg>` — run unit tests scoped to a package (e.g., `task testunit -- LakehouseResource ./internal/services/lakehouse/`)
- `task testacc -- <Group> <Pkg>` — run acceptance tests scoped to a package (e.g., `task testacc -- LakehouseResource_CRUD ./internal/services/lakehouse/`)

Always provide the package path as the second argument for faster execution. Without it, the runner scans all packages (`./...`).

## Non-Item Testing Specifics

### Fake Pattern Decision Tree

Non-Item resources use **either** centralized fakes or inline fakes depending on their API shape:

| Condition                                                                              | Pattern                                                          | Location                          | Example                                   |
| -------------------------------------------------------------------------------------- | ---------------------------------------------------------------- | --------------------------------- | ----------------------------------------- |
| API uses standard `entityID` or `parentID/entityID` paths                              | **Centralized** — `simpleIDOperations` or `parentIDOperations`   | `internal/testhelp/fakes/`        | workspace, domain, folder                 |
| Resource is polymorphic with multiple subtypes sharing a client                        | **Centralized** — shared configure + per-subtype wrappers        | `internal/testhelp/fakes/`        | gateway, connection                       |
| API has 3+ path parameters or non-standard composite IDs that don't fit `typedHandler` | **Inline** — package-level map store + direct handler assignment | `fake_test.go` in service package | shortcut (`workspaceID/itemID/path/name`) |

**Default:** Use centralized fakes unless the SDK method signatures don't fit the `typedHandler` interface (i.e., more than 2 path segments or non-standard ID composition).

### Inline Fake Pattern (`fake_test.go`)

For resources that require inline fakes, each service defines its own fake functions directly in the test package:

```go
// fake_test.go — package <service>_test
func fake<Type>(exampleResp fabcore.<Type>) func(ctx context.Context, ...) (resp azfake.Responder[...], errResp azfake.ErrorResponder) {
    return func(_ context.Context, ...) (resp azfake.Responder[...], errResp azfake.ErrorResponder) {
        resp = azfake.Responder[...]{}
        resp.SetResponse(http.StatusOK, ..., nil)
        return resp, errResp
    }
}

func NewRandom<Type>() fabcore.<Type> {
    return fabcore.<Type>{
        ID:   new(testhelp.RandomUUID()),
        // ... populate all fields with random test data
    }
}
```

### Inline Fake Wiring

For resources using inline fakes (per the decision tree above), wire them directly to the SDK fake server. For resources needing state across operations, use a map store:

```go
var fakeStore = map[string]fabcore.<Type>{}
```

For simpler resources (like role assignments), pass the fake response directly to the handler function.

### Non-Item Import State Formats

Non-Item resources may use composite IDs beyond the standard `workspaceID/entityID`:

| Resource           | Import ID Format                          |
| ------------------ | ----------------------------------------- |
| Shortcut           | `workspaceID/itemID/path/name`            |
| Role Assignments   | `parentID/roleAssignmentID`               |
| Standard Non-Items | `workspaceID/entityID` or just `entityID` |

### Polymorphic Test Data

Non-Items with polymorphic types (gateway, connection, role assignments) should generate test entities covering different type variants:

```go
func NewRandom<Types>() fabcore.<Types> {
    return fabcore.<Types>{
        Value: []fabcore.<Type>{
            { /* variant A — e.g. GroupPrincipal */ },
            { /* variant B — e.g. UserPrincipal */ },
            { /* variant C — e.g. ServicePrincipal */ },
        },
    }
}
```

Reference: `internal/services/shortcut/fake_test.go`, `internal/services/workspacera/fake_test.go`

### Fake Responder Types

Non-Item fakes use three responder types depending on the SDK operation:

| Responder Type           | SDK Operation      | Setup Method          | Example                     |
| ------------------------ | ------------------ | --------------------- | --------------------------- |
| `azfake.Responder`       | Sync (Get/Create)  | `SetResponse`         | `connectionra/fake_test.go` |
| `azfake.PagerResponder`  | List (paged)       | `AddPage`             | `tags/fake_test.go`         |
| `azfake.PollerResponder` | Long-running (LRO) | `SetTerminalResponse` | `workspacegit/fake_test.go` |

```go
// Sync — standard Get/Create/Update/Delete
resp = azfake.Responder[fabcore.ClientGetResponse]{}
resp.SetResponse(http.StatusOK, fabcore.ClientGetResponse{Entity: entity}, nil)

// Pager — list operations returning pages
resp = azfake.PagerResponder[fabcore.ClientListResponse]{}
resp.AddPage(http.StatusOK, fabcore.ClientListResponse{Entities: entities}, nil)

// Poller — long-running operations (initialize, commit, etc.)
resp = azfake.PollerResponder[fabcore.ClientInitializeResponse]{}
resp.SetTerminalResponse(http.StatusOK, fabcore.ClientInitializeResponse{Result: result}, nil)
```

### Stateful Fakes for CRUD Tests

When a resource's unit CRUD test needs to verify updates, use a state struct so the Get fake returns the mutated entity after Update:

```go
type entityState struct {
    currentEntity fabcore.Entity
}

func fakeStatefulGet(state *entityState) func(...) (resp, errResp) {
    return func(...) (resp, errResp) {
        resp.SetResponse(http.StatusOK, Response{Entity: state.currentEntity}, nil)
        return resp, errResp
    }
}

func fakeStatefulUpdate(updatedEntity fabcore.Entity, state *entityState) func(...) (resp, errResp) {
    return func(...) (resp, errResp) {
        state.currentEntity = updatedEntity
        resp.SetResponse(http.StatusOK, Response{Entity: state.currentEntity}, nil)
        return resp, errResp
    }
}
```

Reference: `internal/services/connectionra/fake_test.go`

### Error Simulation in Fakes

Use `fabfake.SetResponseError` to simulate SDK error responses (e.g. 404 Not Found):

```go
import fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

// Inside a fake handler — return not-found error
if _, ok := store[id]; !ok {
    errResp.SetError(fabfake.SetResponseError(http.StatusNotFound, "ItemNotFound", "Item not found"))
    resp.SetResponse(http.StatusNotFound, Response{}, nil)
    return resp, errResp
}
```

Reference: `internal/services/tags/fake_test.go`, `internal/services/externaldatashare/fake_test.go`
