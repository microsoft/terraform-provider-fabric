# Sprout upgrade: PR onboarding

## What this PR is trying to do

The provider's Go-template renderer previously used **native Sprout v1.0.3**.
This PR upgrades it to v1.1.2 as part of a grouped dependency update. Several
template functions changed signatures or results, so an unadapted upgrade can
change generated Fabric item definitions or their hashes, or make plans fail.

The design is a **temporary, selective provider 1.x compatibility layer**:
preserve important existing calls and selected results while using native v1.1.2
where changes are accepted. It is not complete v1.0.3 emulation and is not
Sprig compatibility. The intended removal point is provider 2.0, with a
customer migration guide before removal. There is no per-template mode switch.

The governing tradeoff is:

- **More compatibility:** fewer existing template changes, but more local
  parsing or conversion logic that could diverge from v1.0.3 and must be
  maintained, reviewed, and removed later.
- **Thinner wrappers:** more behavior comes directly from upstream and less
  local logic must be maintained, but customers can see changed output or
  errors in 1.x.

A corrected result is still a compatibility change if a customer relied on
the old result. A wrapper is not necessarily correct just because it is small.

## Where to look

| File | Role |
| --- | --- |
| `internal/pkg/transforms/payload.go` | Selects supported Sprout registries, builds the function map and installs the compatibility layer. |
| `internal/pkg/transforms/template_compat_v1.go` | Central registration of argument adapters, aliases and selected result compatibility. |
| `internal/pkg/transforms/template_compat_v1_strings.go` | Local old-style camel and Pascal conversion. |
| `internal/testhelp/template_compat_v1.go` and `internal/testhelp/fixtures/template_compat_v1_*.json` | Shared renderer and Notebook cases and old-output fixtures. |
| `internal/pkg/transforms/template_compat_v1_test.go` | Rendering, payload-hash and error regression checks. |
| `internal/services/notebook/template_compat_v1_test.go` | Resource lifecycle checks using the shared cases. |
| `templates/guides/gotemplates.md` | Source of customer-facing behavior and migration documentation; `docs/guides/gotemplates.md` is generated. |
| `internal/pkg/transforms/testdata/compat_discussion/dig.go` | **Discussion only:** possible old `dig` implementation, not registered or used by the provider. |
| `internal/pkg/transforms/testdata/compat_discussion/separated_case.go` | **Discussion only:** possible old separated-case algorithm, not registered or used by the provider. |

The discussion package is inside `testdata`, so normal package discovery and
provider builds do not include it. Its tests can be run by explicitly selecting
that package. Adopting either implementation requires an explicit registration,
integration tests and documentation updates; simply adding the files does not
change template behavior.

## How to read the function table

**Compatibility** describes the current provider 1.x behavior relative to
native Sprout v1.0.3. **Risk** combines possible customer impact with the
complexity of maintaining the current code; it is qualitative, not a frequency
estimate. “Selected” means tested or intentionally preserved behavior, not a
proof that all inputs or diagnostics are identical.

| Function | Current compatibility | Main risk and available tradeoff |
| --- | --- | --- |
| `get` | Both map-first and map-last forms are adapted. | **Medium:** removing the adapter breaks old calls; validation and error text can differ from v1.0.3. Keep for call compatibility or remove for a thinner breaking 1.x update. |
| `hasKey` | Both map-first and map-last forms are adapted. | **Medium:** same call-order and validation tradeoff as `get`. |
| `unset` | Both map-first and map-last forms are adapted. | **Medium:** old calls need the adapter; preserve mutation/error behavior in tests. |
| `set` | Both map-first and map-last forms are adapted; nil maps are rejected before delegation. | **Medium:** legacy ordering and map-valued arguments need care. Removing it breaks old calls. |
| `pick`, `omit` | Old map-first and new map-last calls are normalized. | **Medium:** variadic validation remains upstream; test empty/invalid arguments and both orders rather than assuming exact diagnostics. |
| `append`, `prepend` | Both orders work when identifiable; a list-like or `nil` first argument selects the old list-first interpretation. `mustAppend` and `mustPrepend` point to the adapters. | **High ambiguity:** two list operands cannot reveal intent. One swap avoids the v1.0.3 recursive stack overflow, so those formerly failing cases are deliberately not reproduced. Removing the adapter breaks old calls. |
| `slice`, `mustSlice` | With one argument, return it unchanged, including scalar, `nil`, map or an empty list. With indices, normalize the old list-first form. | **Medium:** this restores old pass-through but accepts non-lists. Removing pass-through makes the layer thinner and validates inputs, but changes old results including `[]` versus `null`. |
| `without`, `mustWithout` | Normalize legacy list-first and identifiable list-last forms; `nil` first uses the old interpretation. | **Medium to high:** ambiguity remains, and upstream uses comparable values as map keys; nested list/map values are not generally valid inputs. Removing the adapter breaks old calls. |
| `regexFindAll`, `regexSplit` | Retain v1.0.3 order: pattern, input, limit. Their old `must*` aliases are retained. | **Medium:** removing adapters breaks old calls. New-order calls are not interchangeable during 1.x. |
| `regexReplaceAll`, `regexReplaceAllLiteral` | Retain v1.0.3 order: pattern, input, replacement, including `must*` aliases. | **High customer impact:** all arguments are strings, so new-order calls can silently produce the wrong result. Keep old meaning in 1.x or require a coordinated migration. |
| `regexFind`, `regexMatch` and their old `must*` aliases | Native functions retained; missing old aliases restored. | **Low maintenance:** removing the alias breaks templates using it. |
| `regexFindGroups`, `regexFindAllGroups`, `regexFindNamed`, `regexFindAllNamed`, `regexQuoteMeta` | Native functions; no argument reorder applied. | **Low adapter risk:** do not assume all regex functions had the same old order. |
| `substr` | Preserve an end index resolving to zero as an empty result; otherwise delegate upstream. | **Medium:** removing the adapter can silently change generated content. It is a small targeted compatibility rule, not a claim of exhaustive equivalence. |
| `toCamelCase`, `toPascalCase` | Local implementation preserves selected v1.0.3 capitalization and generated names. | **Medium to high maintenance:** algorithm and Unicode edge behavior must be compared to old output. Removing it reduces code but changes names such as `firstname` or `HTMLParser`. |
| `toSnakeCase`, `toKebabCase`, `toDotCase`, `toPathCase`, `toConstantCase` | Native v1.1.2; output differences are accepted in 1.x. | **Medium to high customer impact:** names and hashes may change. The discussion-only alternative can be adopted if stronger preservation justifies owning more algorithmic code. |
| `dig` | Native v1.1.2 parser, with new escaping and empty-segment validation; no Sprig positional default. | **Medium to high customer impact:** existing backslash or empty-segment paths can fail or select differently. The discussion-only old parser is an option if compatibility is prioritized. |
| `ellipsisBoth` | Native v1.1.2 rune handling, including the Unicode fix. | **Medium:** affected Unicode output and hashes can change. Reproducing old byte behavior improves exact output compatibility but can produce invalid UTF-8. |
| Other registered functions | Native v1.1.2 unless explicitly overridden. | **Unknown per function:** upstream changes can remain. Evaluate differences before claiming complete compatibility. |

Aliases are resolved before the compatibility overrides, so affected old
aliases must be explicitly pointed to the same adapter. Unsupported registries
remain excluded. Sprigin is not a substitute: it targets Sprig semantics,
including different `dig`, error and other behaviors.

## Concrete tradeoffs and examples

### Map, list and regex calls

```gotemplate
{{ append (list "a") "b" }}                 {{/* old list-first: retained */}}
{{ list "a" | append "b" }}                 {{/* new order: retained */}}
{{ list "a" | append (list "b") | toJson }} {{/* ["b",["a"]] in 1.x */}}

{{ regexReplaceAll "a" "abc" "x" }}        {{/* xbc in 1.x */}}
```

In the nested-list pipeline the first argument is treated as the destination,
because argument values cannot distinguish old and new intent. The regex
replacement order cannot be inferred from types because all three values are
strings. Removing these adapters is possible, but requires users to migrate
calls with the provider upgrade and can silently alter replacement output.

### `slice` without indices

| Expression | Current provider 1.x | Native v1.1.2 |
| --- | --- | --- |
| `slice (list) \| toJson` | `[]` | `null` |
| `slice "abc"` | `abc` | Error |
| `slice nil \| toJson` | `null` | Error |

The scalar/`nil` pass-through deliberately matches v1.0.3 even though those
inputs are not lists. The same applies to `mustSlice`. Removing the special
case would make validation stricter but would be a breaking behavior choice.

### `dig`

```gotemplate
{{ dict "a" (dict "" (dict "b" "value")) | dig "a..b" }}
```

This returned `value` in v1.0.3; the current native parser returns an error.
Old backslashes were literal; new `\\.` and `\\\\` template-string forms
express a literal dot and backslash to the new parser, respectively. Adopting
the discussion `dig` implementation would restore old splitting and literal
backslashes, but would make the new escape interpretation unavailable in 1.x.

### Case conversion

| Expression | v1.0.3 | Current provider 1.x |
| --- | --- | --- |
| `toCamelCase "firstName"` | `firstname` | `firstname` |
| `toPascalCase "HTMLParser"` | `HTMLParser` | `HTMLParser` |
| `toSnakeCase "0"` | `_0` | `0` |
| `toPathCase "0"` | `/0` | `0` |

Camel/Pascal are currently wrapped; the other five converters use upstream.
The discussion implementation could add old separated behavior, including
legacy byte-based lookahead and leading separators. Conversely, removing
camel/Pascal wrappers would make those functions native too and may change
generated names or hashes. Neither choice guarantees all inputs without
adequate old-version comparison.

## What could be added or removed

| Choice | Add/remove | Benefit | Cost |
| --- | --- | --- | --- |
| More `dig` compatibility | Register and integrate the discussion-only legacy parser. | Preserves old dotted and backslash path meaning. | More code and tests; loses new escaping in 1.x. |
| More case compatibility | Register separated converters using the discussion-only algorithm. | Preserves more old generated names and hashes. | More algorithmic code and old quirks; verify Unicode and boundaries. |
| Stricter `slice` validation | Remove one-argument pass-through and decide separately whether to preserve valid empty lists. | Simpler behavior for invalid inputs. | Breaks old scalar/`nil` calls and potentially `[]` serialization. |
| Native camel/Pascal | Remove their two overrides and isolated helper. | Less local algorithmic code. | Changes existing generated output. |
| Native `substr` | Remove the zero-end adapter. | A smaller shim. | Changes zero-end results. |
| Native map/list or regex order | Remove individual adapters and redirect aliases accordingly. | Less adapter code and earlier alignment with v1.1.2. | Existing calls fail or silently change; major-version migration is safer. |
| Full old-library behavior | Keep v1.0.3 or assess a separate frozen implementation. | More direct old implementation. | Defers upstream fixes or adds dependency/fork overhead; can retain process-fatal defects. |

Do not restore the old recursive list overflow or byte-splitting Unicode behavior
solely to claim equivalence. Decide explicitly whether other old quirks are
worth preserving. Customer usage and the severity of generated-output changes
should guide the decision.

## Evidence and limits

- Old-output fixtures and selected tests were captured or checked against
  v1.0.3. They demonstrate covered cases, **not all possible inputs**.
- The renderer checks results, errors and payload hashes; the Notebook test
  exercises a representative item lifecycle. Live Fabric acceptance requires
  configured fixtures and credentials.
- Discussion implementations are not part of the production function map.
  Their focused tests and a one-off comparison against v1.0.3 (25,000
  separated-case outputs and 11 `dig` cases) cover examples, not every input.
- Error behavior, nil and empty values, alias rebinding, ambiguous arguments,
  Unicode, case conversion and restrictions on available registries deserve
  particular attention during review.
- Regenerate customer-facing docs from `templates/` when behavior changes.
  Keep the 2.0 migration guide aligned with whichever 1.x choices are selected.
