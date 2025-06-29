# TF-MCP: Terraform Fabric MCP Server

This is a Model Context Protocol (MCP) server that provides workflow guidance for Fabric Terraform provider development.

## Features

- **getfabricitempropertiesworkflow**: Returns step-by-step instructions for completing properties-based Fabric item implementations
- **getFabricSwagger**: Returns the Swagger definition of the Fabric artifact with complete API specification

## Setup for Use in Other VS Code Projects

### Method 1: Using the MCP Configuration

1. **Copy the MCP configuration** to your project's MCP settings:

```json
{
  "mcpServers": {
    "tf-mcp": {
      "command": "python",
      "args": ["-m", "server"],
      "cwd": "q:\\Repos\\TF-MCP",
      "env": {}
    }
  }
}
```

2. **Add to your VS Code settings.json** or MCP configuration file in your target project.

### Method 2: Install as a Package

1. **Install the package** in your target project:

```bash
pip install -e q:\Repos\TF-MCP
```

2. **Run the server**:

```bash
tf-mcp-server
```

### Method 3: Direct Python Execution

1. **Run directly** from the TF-MCP directory:

```bash
cd q:\Repos\TF-MCP
python server.py
```

## MCP Client Configuration

To use this server in another VS Code project with GitHub Copilot, add this to your MCP configuration:

### VS Code Settings (settings.json)

```json
{
  "mcp.servers": {
    "tf-mcp": {
      "command": "python",
      "args": ["q:\\Repos\\TF-MCP\\server.py"],
      "env": {}
    }
  }
}
```

### Alternative: Using uv (if available)

```json
{
  "mcp.servers": {
    "tf-mcp": {
      "command": "uv",
      "args": ["--directory", "q:\\Repos\\TF-MCP", "run", "python", "server.py"],
      "env": {}
    }
  }
}
```

## Available Resources

### getfabricitempropertiesworkflow

Returns the complete workflow instructions for properties-based Fabric items, including:

- Provider registration steps
- Fake server configuration
- API definition requirements
- Properties contract creation
- Schema implementation
- Testing and validation

### getFabricSwagger

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

Once configured, you can use the resources in your Copilot conversations:

### Get Properties Workflow

```
@tf-mcp Please provide the workflow for completing a properties-based Fabric item
```

The server will return the complete step-by-step instructions from the properties.md template.

### Get Fabric Swagger Definition

```
@tf-mcp Please provide the Swagger definition for Fabric artifacts
```

The server will return the complete Swagger/OpenAPI specification with all endpoint definitions, schemas, and data models for Fabric items.

## Development

To modify or extend this server:

1. Edit `server.py` to add new resources or tools
2. Update `pyproject.toml` if adding dependencies
3. Test locally with `python server.py`

## Troubleshooting

- Ensure Python 3.13+ is installed
- Verify all dependencies are installed: `pip install -r requirements.txt` or use the pyproject.toml
- Check that the path to TF-MCP directory is correct in your configuration
- Ensure the `itemgen/templates/mds/properties.md` file exists
