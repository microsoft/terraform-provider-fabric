# MCP Instructions: Properties Item Post-Generation Steps

The Fabric Item Generator has successfully created a new properties-based item. Follow these steps to complete the implementation.

## Step 1: Obtain API Definition

### Option A: Use the MCP Tool (Recommended)

**Use the `get_fabric_swagger` MCP tool to get the swagger definition:**

- The MCP server provides a built-in tool that returns the complete Swagger definition for Fabric artifacts
- This includes all API endpoints, request/response schemas, and data models
- The tool returns a comprehensive OpenAPI/Swagger specification in JSON format

**Usage:**

```text
Use the get_fabric_swagger MCP tool to get the complete API specification
```

## Step 2: Create Properties Contract

### File: `internal/services/<item_type>/models.go`

Once you have the definition JSON:

1. **Analyze the properties object**: Look for the item type's properties definition in the JSON
2. **Create the properties contract**: Under the item type object, create a properties contract that matches the API specification
3. **Map data types**: Ensure proper mapping between API types and Go types

Example structure:

```go
type <ItemType>Properties struct {
    // Properties based on API definition
}
```

## Step 3: Implement Schema Attributes

### Files

- `schema_data_<item_type>.go`
- `schema_resource_<item_type>.go`

Using the properties contract from Step 4:

1. **Complete `getPropertiesAttributes` function**: Fill in the schema attributes based on the API definition
2. **Map each property**: For each property in the API definition, create corresponding Terraform schema attributes
3. **Set proper validation**: Add validation rules, required fields, and default values as specified in the API

Example pattern:

```go
func getPropertiesAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        "property_name": schema.StringAttribute{
            Description: "Description from API spec",
            Required:    true, // or Optional based on API
            // Additional validation as needed
        },
        // Additional properties...
    }
}
```

## Step 4: Add Properties Validations

### Files to Update

- `*_test.go` files in the service directory (resource and data source tests)

After implementing the schema attributes:

1. **Review test files**: Look for TBD comments that indicate where properties validation should be added
2. **Add specific property validations**: Replace the TBD comments with actual property value validations
3. **Validate computed properties**: Ensure computed properties from the API are properly validated in tests

Example pattern for replacing TBD comments:

```go
// Replace this TBD comment:
// TBD: Add properties validation - validate specific property values based on API definition

// With specific property validations like:
resource.TestCheckResourceAttr(testResourceItemFQN, "properties.specific_property", "expected_value"),
resource.TestCheckResourceAttrSet(testResourceItemFQN, "properties.computed_property"),
resource.TestCheckResourceAttrWith(testResourceItemFQN, "properties.validated_property", func(value string) error {
    // Custom validation logic
    return nil
}),
```

**Important Note**: When implementing the fake test files and removing TBD comments, you should complete ALL TBDs in a single generation/update rather than going through them one by one across multiple generations.

## Step 5: Complete Fake Test Implementation

### Files to Update

- `internal/testhelp/fakes/fabric_<item_type>.go`

**`NewRandom<ItemType>()` function**: Contains TODO for properties implementation

### Required Actions

**Complete the Properties in `NewRandom<ItemType>()` function:**

1. **Locate the TODO comment**: Find the comment `// TODO: Add Properties field with appropriate test data based on API definition`
2. **Implement the Properties field**: Based on your API definition analysis, add the Properties field with realistic test data
3. **Follow the pattern**: Use the same pattern as other fake implementations (e.g., `fabric_kqldatabase.go`)

Example implementation:

```go
func NewRandom<ItemType>() fab<package>.<ItemType> {
 return fab<package>.<ItemType>{
  ID:          to.Ptr(testhelp.RandomUUID()),
  DisplayName: to.Ptr(testhelp.RandomName()),
  Description: to.Ptr(testhelp.RandomName()),
  WorkspaceID: to.Ptr(testhelp.RandomUUID()),
  Type:        to.Ptr(fab<package>.ItemType<ItemType>),
  Properties: &fab<package>.Properties{
   // Fill with appropriate test data based on API definition
   PropertyName1: to.Ptr("test_value"),
   PropertyName2: to.Ptr(testhelp.RandomUUID()),
   // ... additional properties
  },
 }
}
```

## Step 6: Validation and Testing

After completing the above steps:

1. **Run tests**: use the task file command `task testunit -- <itemType>`

## Common Patterns to Follow

### Property Type Mapping

- **String fields**: Use `schema.StringAttribute`
- **Boolean fields**: Use `schema.BoolAttribute`
- **Integer fields**: Use `schema.Int64Attribute`
- **Object fields**: Use `schema.SingleNestedAttribute`
- **Array fields**: Use `schema.ListAttribute` or `schema.SetAttribute`

### Validation Rules

- **Required fields**: Set `Required: true`
- **Optional fields**: Set `Optional: true`
- **Computed fields**: Set `Computed: true`
- **Default values**: Use `Default: <default_value>`

## Troubleshooting

If you encounter issues:

1. **Schema mismatch**: Double-check the API definition against your schema
2. **Type errors**: Ensure proper type mapping between API and Terraform
3. **Test failures**: Review the fake server configuration and test data
4. **Provider registration**: Verify the resource and data source names are correct

## Final Checklist

- [ ] Provider registration updated
- [ ] API definition obtained and analyzed
- [ ] Properties contract created
- [ ] Schema attributes implemented
- [ ] Properties validations added to test files
- [ ] Fake test Properties field implemented with realistic test data
- [ ] ALL TBD comments replaced with actual implementations
- [ ] Tests passing
- [ ] Documentation updated

Your properties-based item is now ready for integration into the Fabric Terraform provider.
