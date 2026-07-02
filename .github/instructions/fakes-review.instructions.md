---
applyTo: "internal/testhelp/fakes/**/*.go,internal/services/**/fake_test.go"
---

# Fake Handler Review Checklist — terraform-provider-fabric

## Purpose & Scope

Guidance for Copilot code review of unit-test fake handlers: the shared
handlers under `internal/testhelp/fakes/` and the inline `fake_test.go` files
under `internal/services/**`. These build SDK response objects for the fake
server. Flag the recurring defects below.

## Pointer Helpers (most common defect)

- `new(...)` takes a **type**, not a value. Flag `new(someValue)`,
  `new("literal")`, `new(testhelp.RandomUUID())`, or `new(fabpkg.EnumConst)` —
  these do not compile.
- Use the pointer helper for values. Files under `internal/testhelp/fakes/`
  import `github.com/Azure/azure-sdk-for-go/sdk/azcore/to` unaliased and call
  `to.Ptr(...)`; some service packages alias it as `azto.Ptr(...)`. Prefer the
  form already used by the file under review.

```go
// Avoid — does not compile
Etag: new("fake-etag"),
State: new(fabcore.ConnectionAccessActionTypeDeny),

// Prefer (in fakes/**; match the file's existing convention)
Etag: to.Ptr("fake-etag"),
State: to.Ptr(fabcore.ConnectionAccessActionTypeDeny),
```

## Response Completeness & Parity

- If the resource/data source model's `set(...)` reads and dereferences a field
  (e.g. `Etag`), the fake response must populate it. Flag omitted fields that
  the model dereferences — tests will panic or diverge from real responses.
- Keep fake responses representative of the real API shape (same required
  fields, nested objects, enum values).

## Robustness

- Flag slice indexing without a length guard (`Value[0]`, `entity.Value[0].ID`)
  and pointer dereferences without a nil check — brittle handlers cause
  hard-to-debug panics.

## Naming

- Constructors: `NewRandom<Type>` (e.g. `NewRandomConnection`). Flag awkward or
  incorrect pluralization (e.g. `DataAccessesSecurity`) and typos.

## Fixture References

Fixture files themselves live under `internal/testhelp/fixtures/**` (outside
this file's scope). When a fake or test **references** a fixture:

- Flag template tokens the test/fake passes (e.g. a `SCOPE` value) that the
  referenced template never uses, and vice versa.
- Flag references to fixture paths that do not exist in the repo.

## Test Interaction

- Fakes mutate a shared, unsynchronized `FakeServer`. Flag `resource.ParallelTest`
  in tests that `Upsert` into the shared fake server — this races.
- Flag typos in identifiers used across the fake and its tests
  (e.g. `workspaeID` → `workspaceID`) and mismatched IDs between the inserted
  fake entity and the config that reads it back.
