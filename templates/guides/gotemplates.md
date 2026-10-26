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

Provider 1.x preserves the template behavior of Sprout v1.0.3 when upgrading to
Sprout v1.1.2. Existing supported templates do not need to change for this update.
This is compatibility with the provider's existing Sprout functions, not with all
Sprig functions. Functions from unsupported registries remain unavailable.

The compatibility layer is planned for removal in provider 2.0. Update your
templates together with that major upgrade, using the migration examples below.

| Functions | Provider 1.x behavior | Planned provider 2.x behavior |
| --------- | --------------------- | ---------------------------- |
| `get`, `set`, `unset`, `hasKey`, `pick`, `omit` | Accept both map-first and map-last calls | Map-last calls only |
| `append`, `prepend`, `slice`, `without` | Accept both orders when types are unambiguous; a list-like first argument uses the legacy list-first interpretation | List-last calls only |
| `regexFindAll`, `regexSplit` | Pattern, input string, limit | Pattern, limit, input string |
| `regexReplaceAll`, `regexReplaceAllLiteral` | Pattern, input string, replacement | Pattern, replacement, input string |
| `dig` | Split keys on every dot; backslashes are literal and empty path segments are allowed | Interpret backslash escapes and reject empty dotted path segments |
| `substr` | End index `0` selects an empty substring | End index `0` selects through the end of the string |
| `ellipsisBoth` | Preserve safe byte slicing when no truncation is needed; otherwise use Sprout v1.1.2 behavior | Use rune offsets consistently |
| `toCamelCase`, `toPascalCase`, `toSnakeCase`, `toKebabCase`, `toDotCase`, `toPathCase`, `toConstantCase` | Preserve Sprout v1.0.3 capitalization and word boundaries | Use the new Sprout case-conversion rules |
| `mustRegexFind`, `mustRegexFindAll`, `mustRegexMatch`, `mustRegexSplit`, `mustRegexReplaceAll`, `mustRegexReplaceAllLiteral` | Keep existing aliases | Use the corresponding function without the `must` prefix |

There is no per-template switch between the old and new behavior. For map and
list functions, provider 1.x accepts both argument orders when the types make
the target clear. The regex functions below are pinned to the 1.x argument
order. Where the same call can mean two things, provider 1.x preserves the
established meaning; the new interpretation is unavailable until provider 2.0.

In particular, keep these limitations in mind:

- **List functions:** `append`, `prepend`, `slice`, and `without`, including
  their `must*` aliases, treat a list-like first argument as the legacy target.
  If the value to append or prepend is itself a list, use the legacy target-first
  form. For example, `append (list "a") (list "b")` appends the second list as
  one value to the first, producing `["a",["b"]]`. Do not use
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
- **`dig`:** provider 1.x splits keys on every dot; backslashes do not escape
  dots, and empty path segments are allowed. The new escape syntax is not
  available. Use `get` to look up a literal key containing a dot. `dig` also
  has no positional default argument; pipe its result to `default` instead.
- **`substr`:** an end index of `0` (or a negative end equal to the string's
  byte length, which resolves to index `0`) returns an empty string in 1.x.
  The new meaning, “through the end of the string,” is not selectable until 2.0.
- **Case conversions:** the seven functions listed in the table use the old
  output for every call in 1.x; there is no alternate syntax for the new rules.
  The compatibility adapter does not override `toTitleCase`.
- **`ellipsisBoth`:** the old byte-offset behavior is retained in its
  non-truncating branch only when the byte offset is within the string. Because
  it is byte-based, an in-range offset can split a multi-byte rune. After the
  minimum-width checks, positive offsets beyond the string use Sprout v1.1.2's
  bounds handling rather than restoring the old out-of-range slice panic.
  Negative offsets remain invalid and can produce an execution error. The new
  consistent rune-offset behavior is not selectable in 1.x.

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

For a fallback from `dig`, pipe its result to `default`:

```gotemplate
{{ $dict | dig "name" | default "unknown" }}
```

Regex argument-order changes must accompany the provider 2.0 upgrade:

```gotemplate
{{/* Provider 1.x */}}
{{ regexFindAll "[a-z]+" "ab12cd" -1 }}
{{ regexReplaceAll "[0-9]+" "item123" "X" }}

{{/* Planned provider 2.x */}}
{{ "ab12cd" | regexFindAll "[a-z]+" -1 }}
{{ "item123" | regexReplaceAll "[0-9]+" "X" }}
```

Check generated names and definition content when migrating case conversions.
For example, `toCamelCase "firstName"` produces `firstname` in provider 1.x but
`firstName` with Sprout's new rules. `toPascalCase "hELLO world"` changes from
`HELLOWorld` to `HelloWorld`. These differences can change content hashes.

Library fixes for crashes, concurrent random shuffling, and division by zero are
retained. Compatibility does not restore crashes or invalid arithmetic results.
Invalid template calls continue to produce provider error diagnostics rather than
silently returning a successful payload.
