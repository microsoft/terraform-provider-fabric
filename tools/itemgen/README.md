# Fabric Item Generator

This tool helps automate the process of onboarding new items to the Fabric Terraform provider. It generates all the necessary files with the correct structure and content based on templates.

## Usage

To generate a new item, run the generator with the following flags:

```bash
go run tools/itemgen/main.go \
  -item-name="<item-name>" \
  -items-name="<items-name>" \
  -item-type=<item-category> \
  -rename-allowed=<true|false> \
  -is-preview=<true|false> \
  -is-spn-supported=<true|false>
```

### Parameters

- `-item-name`: The display name of the item (e.g., "Data Pipeline")
- `-items-name`: The display name of the item in plural form (e.g., "Data Pipelines")
- `-item-type`: The category of item to generate (see Item Types section below)
- `-rename-allowed`: Whether the item can be renamed (default: true)
- `-is-preview`: Whether the item is in preview (default: false)
- `-is-spn-supported`: Whether the item supports SPN (default: false)

### Item Types

The generator supports 6 different item types:

1. **Definition Only** (`definition`)
   - Has definition but no properties or config
   - Example: `datapipeline`
   - Uses: `NewDataSourceFabricItemDefinition`, `NewDataSourceFabricItems`, `NewResourceFabricItemDefinition`

1. **Properties Only** (`properties`)
   - Has properties but no definition or config
   - Example: `environment`
   - Uses: `NewDataSourceFabricItemProperties`, `NewDataSourceFabricItemsProperties`, `NewResourceFabricItemProperties`

1. **Definition and Properties** (`definition-properties`)
   - Has both definition and properties but no config
   - Example: `mirroreddatabase`
   - Uses: `NewDataSourceFabricItemDefinitionProperties`, `NewDataSourceFabricItemsProperties`, `NewResourceFabricItemDefinitionProperties`

1. **Config and Properties** (`config-properties`)
   - Has config and properties but no definition
   - Example: `lakehouse`
   - Uses: `NewDataSourceFabricItemProperties`, `NewDataSourceFabricItemsProperties`, `NewResourceFabricItemConfigProperties`

1. **Config, Definition, and Properties** (`config-definition-properties`)
   - Has config, definition, and properties
   - Example: `eventhouse`
   - Uses: `NewDataSourceFabricItemDefinitionProperties`, `NewDataSourceFabricItemsProperties`, `NewResourceFabricItemConfigDefinitionProperties`

### Examples

1. Generate an item with definition (like datapipeline):

   ```bash
   go run tools/itemgen/main.go \
   -item-name="Data Pipeline" \
   -items-name="Data Pipelines" \
   -item-type=definition \
   -rename-allowed=true \
   -is-preview=false \
   -is-spn-supported=true
   ```

1. Generate an item with properties (like environment):

   ```bash
   go run tools/itemgen/main.go \
   -item-name="Environment" \
   -items-name="Environments" \
   -item-type=properties \
   -rename-allowed=true \
   -is-preview=false \
   -is-spn-supported=true
   ```

## Generated Files

The generator will create different files based on the item type:

### Common Files (All Types)

1. `base.go` - Base functionality and type definitions
1. `base_test.go` - Base test setup
1. `data_<type>.go` - Single item data source
1. `data_<type>_test.go` - Single item data source tests
1. `data_<types>.go` - Multiple items data source
1. `data_<types>_test.go` - Multiple items data source tests
1. `resource_<type>.go` - Resource definition
1. `resource_<type>_test.go` - Resource tests

### Additional Files for Items with Properties

For items with properties (including definition-properties, config-properties, and config-definition-properties), the following additional files are generated:

1. `schema_data_<type>.go` - Data source schema
1. `schema_resource_<type>.go` - Resource schema
1. `models.go` - Data models

## File Naming Convention

The generator automatically handles file naming based on the provided item name:

- Single item files use the lowercase, underscore-separated version of the item name
- Plural item files use the lowercase, underscore-separated version of the plural item name
- Example: "Data Pipeline" becomes "data_pipeline" and "data_pipelines"

## Customization

After generating the files, you may need to:

1. Review and update the generated files.
1. Complete all `TBD` placeholders.
1. Add the resource/data-source/s to the provider configuration.
1. Update the well-known script.
1. Run the tests to verify the implementation.

## Templates

The templates are located in `tools/itemgen/templates/<item-type>/`. They follow the structure of the example items and can be modified to add new features or change the generated code structure.

## Safety Features

The generator includes several safety features:

- Checks for existing files before overwriting
- Prompts for confirmation when overwriting existing files
- Provides clear error messages and guidance
- Generates consistent file naming and structure
