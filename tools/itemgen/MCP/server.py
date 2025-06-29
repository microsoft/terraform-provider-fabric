import os
import json
from mcp.server.fastmcp import FastMCP

# Initialize FastMCP server
mcp = FastMCP("TF-MCP")


@mcp.tool()
def get_fabric_item_properties_workflow() -> str:
    """
    Returns the properties item post-generation workflow instructions.
    This provides step-by-step guidance for completing a properties-based Fabric item implementation.
    """
    properties_md_path = os.path.join(
        "tools", "itemgen", "MCP", "itemgen.md")

    try:
        with open(properties_md_path, 'r', encoding='utf-8') as file:
            content = file.read()
        return content
    except FileNotFoundError:
        return """# Error: Properties workflow file not found

The properties.md workflow file could not be found at the expected location.
Please ensure the file exists at: itemgen/templates/mds/properties.md

This file should contain the step-by-step instructions for completing a properties-based Fabric item implementation."""
    except Exception as e:
        return f"""# Error: Could not read properties workflow file

An error occurred while reading the properties workflow file:
{str(e)}

Please check the file permissions and try again."""


@mcp.tool()
def getFabricSwagger() -> str:
    """
    Returns the Swagger definition of the Fabric artifact.
    This provides the complete API specification for Fabric items including endpoints,
    request/response schemas, and data models.
    """
    definition_json_path = os.path.join(
        "tools", "itemgen", "MCP", "definition.json")

    try:
        with open(definition_json_path, 'r', encoding='utf-8') as file:
            swagger_data = json.load(file)
        return json.dumps(swagger_data, indent=2)
    except FileNotFoundError:
        return """# Error: Swagger definition file not found

The definition.json file could not be found at the expected location.
Please ensure the file exists at: tools/itemgen/MCP/definition.json

This file should contain the Swagger API specification for Fabric artifacts."""
    except json.JSONDecodeError as e:
        return f"""# Error: Invalid JSON in swagger definition file

The definition.json file contains invalid JSON:
{str(e)}

Please check the JSON syntax and try again."""
    except Exception as e:
        return f"""# Error: Could not read swagger definition file

An error occurred while reading the swagger definition file:
{str(e)}

Please check the file permissions and try again."""


if __name__ == "__main__":
    mcp.run()


def main():
    """Entry point for the MCP server"""
    mcp.run()
