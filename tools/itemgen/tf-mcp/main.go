package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server with tool capabilities
	s := server.NewMCPServer(
		"TF-MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register tools using the available API (separate tool definition from handler)
	s.AddTool(getFabricItemPropertiesWorkflowTool(), handleGetFabricItemPropertiesWorkflow)
	s.AddTool(getFabricSwaggerTool(), handleGetFabricSwagger)

	// Start the server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// getFabricItemPropertiesWorkflowTool returns the properties item post-generation workflow instructions.
func getFabricItemPropertiesWorkflowTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_fabric_item_properties_workflow",
		Description: "Returns the properties item post-generation workflow instructions. This provides step-by-step guidance for completing a properties-based Fabric item implementation.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}
}

// handleGetFabricItemPropertiesWorkflow handles the get_fabric_item_properties_workflow tool call.
func handleGetFabricItemPropertiesWorkflow(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Construct the path to the properties.md file
	propertiesMdPath := filepath.Join("properties.md")

	// Try to read the file
	content, err := os.ReadFile(propertiesMdPath)
	if err != nil {
		if os.IsNotExist(err) {
			errorMessage := `# Error: Properties workflow file not found

The properties.md workflow file could not be found at the expected location.
Please ensure the file exists at: tools/itemgen/tf-mcp/properties.md

This file should contain the step-by-step instructions for completing a properties-based Fabric item implementation.`

			return &mcp.CallToolResult{}, errors.New(errorMessage)
		}

		return &mcp.CallToolResult{}, fmt.Errorf("could not read properties workflow file: %v", err)
	}

	// Return the content as a successful result - let the error message go through the error return for now
	// We'll adjust this once we understand the Content structure better
	return &mcp.CallToolResult{}, fmt.Errorf("content: %s", string(content))
}

// getFabricSwaggerTool returns the Swagger definition of the Fabric artifact.
func getFabricSwaggerTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_fabric_swagger",
		Description: "Returns the Swagger definition of the Fabric artifact. This provides the complete API specification for Fabric items including endpoints, request/response schemas, and data models.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}
}

// handleGetFabricSwagger handles the getFabricSwagger tool call.
func handleGetFabricSwagger(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Construct the path to the definition.json file
	definitionJsonPath := filepath.Join("definition.json")

	// Try to read the file
	content, err := os.ReadFile(definitionJsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			errorMessage := `# Error: Swagger definition file not found

The definition.json file could not be found at the expected location.
Please ensure the file exists at: tools/itemgen/tf-mcp/definition.json

This file should contain the Swagger API specification for Fabric artifacts.`

			return &mcp.CallToolResult{}, errors.New(errorMessage)
		}

		return &mcp.CallToolResult{}, fmt.Errorf("could not read swagger definition file: %v", err)
	}

	// Parse JSON to validate it and format it
	var swaggerData any
	if err := json.Unmarshal(content, &swaggerData); err != nil {
		return &mcp.CallToolResult{}, fmt.Errorf("invalid JSON in swagger definition file: %v", err)
	}

	// Format the JSON with indentation
	formattedJson, err := json.MarshalIndent(swaggerData, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{}, fmt.Errorf("could not format swagger definition: %v", err)
	}

	// Return the formatted JSON through error for now until we figure out Content structure
	return &mcp.CallToolResult{}, fmt.Errorf("swagger: %s", string(formattedJson))
}
