# TF-MCP: Terraform Fabric MCP Server

This is a Model Context Protocol (MCP) server that provides workflow guidance for Fabric Terraform provider development.

## Features

- **get_fabric_item_properties_workflow**: Returns step-by-step instructions for completing properties-based Fabric item implementations
- **get_fabric_swagger**: Returns the Swagger definition of the Fabric artifact with complete API specification

## Setup for Use in Other VS Code Projects

### Method 1: Using the MCP Configuration

1. **Copy the MCP configuration** to your project's MCP settings:

    ```json
    {
      "servers": {
        "tf-mcp": {
          "command": "go",
          "args": [
            "run",
            "main.go"
          ],
          "cwd": "${workspaceFolder}/tools/itemgen/tf-mcp"
        }
      }
    }
    ```

2. **Add to your VS Code settings.json** or MCP configuration file in your target project.

### Method 2: Build and Run Executable

1. **Build the executable** in the tf-mcp directory:

    ```bash
    cd tools/itemgen/tf-mcp
    go build -o tf-mcp.exe .
    ```

2. **Use the executable in MCP configuration**:

    ```json
    {
      "servers": {
        "tf-mcp": {
          "command": "./tf-mcp.exe",
          "args": [],
          "cwd": "${workspaceFolder}/tools/itemgen/tf-mcp"
        }
      }
    }
    ```

### Method 3: Direct Go Execution

1. **Run directly** from the tf-mcp directory:

    ```bash
    cd tools/itemgen/tf-mcp
    go run main.go
    ```

## MCP Client Configuration

To use this server in another VS Code project with GitHub Copilot, add this to your MCP configuration:

### VS Code Settings (settings.json)

```json
{
  "servers": {
    "tf-mcp": {
      "command": "go",
      "args": [
        "run",
        "main.go"
      ],
      "cwd": "${workspaceFolder}/tools/itemgen/tf-mcp"
    }
  }
}
```

### Alternative: Using Built Executable

```json
{
  "servers": {
    "tf-mcp": {
      "command": "./tf-mcp.exe",
      "args": [],
      "cwd": "${workspaceFolder}/tools/itemgen/tf-mcp"
    }
  }
}
```

## Available Tools

### get_fabric_item_properties_workflow

Returns the complete workflow instructions for properties-based Fabric items, including:

- Provider registration steps
- Fake server configuration
- API definition requirements
- Properties contract creation
- Schema implementation
- Testing and validation

### get_fabric_swagger

Returns the Swagger definition of the Fabric artifact, providing:

- Complete API specification for Fabric items
- Endpoint definitions for CRUD operations
- Request and response schemas
- Data models and definitions
- Parameter specifications
- HTTP status codes and error responses

The swagger definition includes:

- `/workspaces/{workspaceId}/items` - List and create items
- `/workspaces/{workspaceId}/items/{itemId}` - Get, update, and delete specific items
- Comprehensive schema definitions for Item, ItemsResponse, CreateItemRequest, and UpdateItemRequest
- Support for various Fabric item types (Lakehouse, Warehouse, Notebook, Report, etc.)

## Usage Examples

Once configured, you can use the tools in your Copilot conversations:

### Get Properties Workflow

```text
@tf-mcp Please provide the workflow for completing a properties-based Fabric item
```

The server will return the complete step-by-step instructions from the `properties.md` template.

### Get Fabric Swagger Definition

```text
@tf-mcp Please provide the Swagger definition for Fabric artifacts
```

The server will return the complete Swagger/OpenAPI specification with all endpoint definitions, schemas, and data models for Fabric items.

## Development

To modify or extend this server:

1. Edit `main.go` to add new tools or modify existing ones
2. Update `go.mod` if adding dependencies
3. Test locally with `go run main.go`
4. Build with `go build -o tf-mcp.exe .`

## File Structure

The server expects the following files in the same directory:

- `properties.md` - Contains the properties workflow instructions
- `definition.json` - Contains the Swagger/OpenAPI specification

## Troubleshooting

- Ensure Go 1.24.4+ is installed
- Verify all dependencies are available: `go mod tidy`
- Check that the path to tf-mcp directory is correct in your configuration
- Ensure the `properties.md` and `definition.json` files exist in the tf-mcp directory
- For build issues, try `go clean` and rebuild
