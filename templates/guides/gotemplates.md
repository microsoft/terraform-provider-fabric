---
page_title: "Advanced use cases with Go templating"
subcategory: "Guides"
description: |-

---

# Parametrization with Templates and Tokens

Various Microsoft Fabric items (such as Notebook, Report, Semantic Model, Spark Job Definition) use "Definition" for configuration. This object outlines the structure and format of the item. Each Fabric item's definition includes essential system files that specify its characteristics, with different formats and required files for each type.

Definition parts often contain elements that need customization to suit your needs. To make them dynamic, use Go templates. A template file with tokens (key-value pairs) can create dynamic content based on variables and functions in Go.

[Go template](https://pkg.go.dev/text/template) is a templating engine integrated into the Go programming language (a key component of Terraform Provider). It lets you create templates with placeholders and populate them with dynamic data. When used in conjunction with Go templating, Configuration as Code through Terraform Provider for Microsoft Fabric can address complex scenarios for managing and deploying Fabric item definitions.

~> If you opt to use Go templating, we cannot promise assistance for any problems or complications that might occur. We recommend proceeding carefully and ensuring that you have the required expertise to handle any potential issues.

## Template file

In this example, the JSON template contains the payload that will be uploaded to the Fabric API endpoint for the Report. It allows you to reference all defined parameters of the configuration via `{{ .PARAMETER_NAME }}` syntax. For example, `definition.json` uses `{{ .SemanticModelID }}` parameter what you can specified in the `tokens` section in you terraform file:

```json
{
 "version": "4.0",
 "datasetReference": {
  "byPath": null,
  "byConnection": {
   "connectionString": null,
   "pbiServiceModelId": null,
   "pbiModelVirtualServerName": "sobe_wowvirtualserver",
   "pbiModelDatabaseName": "{{ .SemanticModelID }}",
   "name": "EntityDataSource",
   "connectionType": "pbiServiceXmlaStyleLive"
  }
 }
}
```

## Available Go template functions

The Provider includes multiple helper functions offered by the [sprout](https://docs.atom.codes/sprout) library, which you can use within the template file.

-> The Terraform Provider for Microsoft Fabric leverages Go templates, which enable the creation of more intricate templates. However, we strongly advise keeping your templates straightforward. Simply reference variables using `{{ .PARAMETER_NAME }}`.

- [Conversion](https://docs.atom.codes/sprout/registries/conversion) - The Conversion includes a collection of functions designed to convert one data type to another directly within your templates. This allows for seamless type transformations.
- [Checksum](https://docs.atom.codes/sprout/registries/checksum) - The Checksum offers functions to generate and verify checksums, ensuring data integrity. It supports various algorithms for reliable error detection and data validation.
- [Encoding](https://docs.atom.codes/sprout/registries/encoding) - The Encoding offers methods for encoding and decoding data in different formats, allowing for flexible data representation and storage within your templates.
- [Maps](https://docs.atom.codes/sprout/registries/maps) - The Maps offers tools for creating, manipulating, and interacting with map data structures, facilitating efficient data organization and retrieval.
- [Numeric](https://docs.atom.codes/sprout/registries/numeric) - The Numeric includes a range of utilities for performing numerical operations and calculations, making it easier to handle numbers and perform math functions in your templates.
- [Random](https://docs.atom.codes/sprout/registries/random) - The Random provides functions to generate random numbers, strings, and other data types, useful for scenarios requiring randomness or unique identifiers.
- [Regexp](https://docs.atom.codes/sprout/registries/regexp) - The Regexp includes functions for pattern matching and string manipulation using regular expressions, providing powerful text processing capabilities.
- [Semver](https://docs.atom.codes/sprout/registries/semver) - The Semver is designed to handle semantic versioning, offering functions to compare and manage version numbers consistently across your projects.
- [Slices](https://docs.atom.codes/sprout/registries/slices) - The Slices provides utilities for working with slice data structures, including functions for filtering, sorting, and transforming slices in a flexible manner.
- [Std](https://docs.atom.codes/sprout/registries/std) - The Std provides a set of standard functions for common tasks, included by default, making it easy to perform basic operations without additional setup
- [Strings](https://docs.atom.codes/sprout/registries/strings) - The Strings offers a comprehensive set of functions for manipulating strings, including formatting, splitting, joining, and other common string operations.
- [Time](https://docs.atom.codes/sprout/registries/time) - The Time provides tools to manage and manipulate dates, times, and time-related calculations, making it easy to handle time-based data in your projects.
- [Uniqueid](https://docs.atom.codes/sprout/registries/uniqueid) - The Uniqueid offers functions to generate unique identifiers, such as UUIDs, which are essential for creating distinct and traceable entities in your applications.

Unsupported function groups in the provider what are available in sprout: [Crypto](https://docs.atom.codes/sprout/registries/crypto), [Env](https://docs.atom.codes/sprout/registries/env), [Filesystem](https://docs.atom.codes/sprout/registries/filesystem), [Reflect](https://docs.atom.codes/sprout/registries/reflect).

## Template compatibility in provider 1.x

Provider 1.x preserves legacy call forms and selected results from Sprout v1.0.3
when upgrading to Sprout v1.1.2, except for the intentional library changes
described below. Affected templates can render different content or fail
validation. This is compatibility with the provider's existing Sprout functions,
not with all Sprig functions. Functions from unsupported registries remain unavailable.

The compatibility layer is planned for removal in provider 2.0. Update calls
listed as changing in 2.x together with that major upgrade, using the migration
examples below.

| Functions | Provider 1.x behavior | Planned provider 2.x behavior |
| --------- | --------------------- | ---------------------------- |
| `get`, `set`, `unset`, `hasKey`, `pick`, `omit` | Accept both map-first and map-last calls | Map-last calls only |
| `append`, `prepend`, `slice`, `without` | Accept both orders when types are unambiguous; a list-like or `nil` first argument uses the legacy list-first interpretation | List-last calls only |
| `regexFindAll`, `regexSplit` | Pattern, input string, limit | Pattern, limit, input string |
| `regexReplaceAll`, `regexReplaceAllLiteral` | Pattern, input string, replacement | Pattern, replacement, input string |
| `dig` | Interpret backslash escapes and reject empty dotted path segments (intentional parsing change) | Same native path parsing |
| `substr` | End index `0` selects an empty substring | End index `0` selects through the end of the string |
| `ellipsisBoth` | Use rune offsets, including when no truncation is needed (intentional Unicode bug fix) | Same rune-based behavior |
| `toCamelCase` | Preserve Sprout v1.0.3 capitalization | Use the new Sprout camel-case rules |
| `toPascalCase` | Preserve Sprout v1.0.3 capitalization, including acronyms | Use the new Sprout Pascal-case rules |
| `toSnakeCase`, `toKebabCase`, `toDotCase`, `toPathCase`, `toConstantCase` | Use native Sprout v1.1.2 conversion rules (intentional output changes) | Same native conversion rules |
| `mustRegexFind`, `mustRegexFindAll`, `mustRegexMatch`, `mustRegexSplit`, `mustRegexReplaceAll`, `mustRegexReplaceAllLiteral` | Keep existing aliases | Use the corresponding function without the `must` prefix |

There is no per-template switch between the old and new behavior. For map and
list functions, provider 1.x accepts both argument orders when the types make
the target clear. The regex functions below are pinned to the 1.x argument
order. Where the same call can mean two things, provider 1.x preserves the
established meaning; the new interpretation is unavailable until provider 2.0.

In particular, keep these limitations in mind:

- **List functions:** `append`, `prepend`, `slice`, and `without`, including
  their `must*` aliases, treat a list-like or `nil` first argument as the legacy
  target. If a value to append or prepend is a list, or a value to append,
  prepend, or omit is `nil`, use the legacy target-first form. For example,
  `append (list "a") (list "b")` appends the second list as one value to the
  first, producing `["a",["b"]]`, and `append (list "a") nil` produces
  `["a",null]`. Do not use
  `append (list "b") (list "a")` expecting the new value-first interpretation;
  provider 1.x treats the first list as the target. A pipeline does not
  disambiguate them either:
  `list "a" | append (list "b")` produces `["b",["a"]]`.
- **`regexFindAll` and `regexSplit`:** use the 1.x order—pattern, input, limit.
  The planned 2.x order—pattern, limit, input—fails with an argument-type
  error in 1.x. This also applies to their `must*` aliases.
- **`regexReplaceAll` and `regexReplaceAllLiteral`:** use the 1.x
  pattern, input, replacement order. All three arguments are strings, so a
  2.x-order call can be accepted but interpreted using the 1.x order, producing
  an unintended result. For example, `regexReplaceAll "a" "x" "abc"` renders
  `"x"` in 1.x, not `"xbc"`. The same limitation applies to the `must*` aliases.
- **`substr`:** an end index of `0` (or a negative end equal to the string's
  byte length, which resolves to index `0`) returns an empty string in 1.x.
  The new meaning, “through the end of the string,” is not selectable until 2.0.
- **`slice` without indices:** a non-nil empty list remains an empty list in 1.x.
  Native Sprout returns `nil`, which encodes as `null` rather than `[]` in JSON.
- **`toCamelCase` and `toPascalCase`:** keep their old conversion rules to
  preserve generated names. New conversion rules are not selectable until 2.0.
  Other case converters, including `toTitleCase`, are not overridden.

Existing aliases such as `mustAppend`, `mustSlice`, and `mustRegexReplaceAll` keep
the same behavior as their corresponding functions. Other regex functions retain
their argument order, including `regexFindAllGroups` and `regexFindAllNamed`.

For example, map and list calls can already use the pipeline forms in provider 1.x:

```gotemplate
{{/* Existing calls, supported in provider 1.x */}}
{{ get $dict "name" }}
{{ append $list "value" }}

{{/* Direct map-last and list-last forms also work for unambiguous arguments */}}
{{ get "name" $dict }}
{{ append "value" $list }}

{{/* Pipeline forms, supported in 1.x and planned for 2.x */}}
{{ $dict | get "name" }}
{{ $list | append "value" }}
```

The `dig` signature is unchanged: keys followed by the dictionary, with no
positional default argument. For a fallback, pipe its result to `default`:

```gotemplate
{{ $dict | dig "name" | default "unknown" }}
```

`default` also replaces empty values such as `""`, `0`, and `false`, not just
missing keys.

Regex argument-order changes must accompany the provider 2.0 upgrade:

```gotemplate
{{/* Provider 1.x */}}
{{ regexFindAll "[a-z]+" "ab12cd" -1 }}
{{ regexReplaceAll "[0-9]+" "item123" "X" }}

{{/* Planned provider 2.x */}}
{{ "ab12cd" | regexFindAll "[a-z]+" -1 }}
{{ "item123" | regexReplaceAll "[0-9]+" "X" }}
```

When upgrading to provider 2.0, review generated names and definition content.
For example, `toCamelCase "firstName"` produces `firstname` in 1.x but `firstName`
with Sprout's new rules; `toPascalCase "HTMLParser"` changes from `HTMLParser` to
`HtmlParser`. These changes can affect content hashes. Also review empty lists
passed to `slice` without indices: `slice (list) | toJson` changes from `[]` to
`null` with native Sprout.

### Intentional library changes in provider 1.x

`slice` and its `mustSlice` alias now validate single-argument calls using
Sprout's native behavior. A list without indices remains supported, but a scalar
or `nil` now produces an error instead of being returned unchanged. A valid empty
list without indices retains its original result. The legacy list-first argument
order remains supported for calls with indices.

`toSnakeCase`, `toKebabCase`, `toDotCase`, `toPathCase`, and
`toConstantCase` use native Sprout v1.1.2 rather than preserving the old conversion
algorithm. Improved word-boundary handling and removal of unwanted leading-digit
separators are active in this 1.x update:

| Expression | Before this update | After this update |
| ---------- | ------------------ | ----------------- |
| `toSnakeCase "0"` | `_0` | `0` |
| `toKebabCase "0"` | `-0` | `0` |
| `toDotCase "0"` | `.0` | `0` |
| `toPathCase "0"` | `/0` | `0` |
| `toConstantCase "0"` | `_0` | `0` |

These output changes can affect generated names, definition content, and hashes.
Review affected templates now rather than waiting for 2.0. `toCamelCase` and
`toPascalCase` retain their old conversion behavior during 1.x as a deliberate
naming-compatibility choice, not a signature adapter.

`dig` uses Sprout v1.1.2 directly rather than preserving the old path parser.
Unescaped dots still separate nested keys, but `\.` selects a literal dot and
`\\` selects a literal backslash. Empty dotted path segments (`"a..b"`, `".a"`,
and `"a."`) and invalid or trailing backslash escapes now produce errors.
An empty-string key remains accessible as a separate argument, for example
`dig "a" "" "b" $dict`.

Template string literals need their own backslash escaping:

```gotemplate
{{/* Previously: look up a key containing one literal backslash */}}
{{ dig "a\\b" $dict }}

{{/* Now: escape the backslash for dig as well as the template string */}}
{{ dig "a\\\\b" $dict }}

{{/* Now also supported: look up the literal key "a.b" */}}
{{ dig "a\\.b" $dict }}
```

These native parsing rules are active in this 1.x update, not deferred to 2.0.
Existing backslash-containing paths can select a different key or fail, and
empty dotted segments must be passed as separate keys. Review affected templates
and their generated content and hashes.

`ellipsisBoth` uses Sprout v1.1.2 directly: offsets count Unicode code points
(runes), not bytes, in both truncating and non-truncating paths. The old
non-truncating path could split a multi-byte character; that bug is not preserved.
For example, `ellipsisBoth 2 10 "éabc"` previously returned `"abc"` and now returns
`"bc"`. Existing Unicode templates using the affected path can therefore render
different content and hashes during this 1.x update; review their generated
definitions. The minimum-width rules and truncation algorithm are unchanged.
This fix is already active in 1.x; it does not wait for provider 2.0.

Library fixes for crashes, concurrent random shuffling, and division by zero are
retained. Compatibility does not restore crashes or invalid arithmetic results.
Invalid template calls continue to produce provider error diagnostics rather than
silently returning a successful payload.
