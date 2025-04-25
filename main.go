// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/microsoft/terraform-provider-fabric/internal/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary.
var version = "dev"

// goreleaser can pass other information to the main package, such as the specific commit
// https://goreleaser.com/cookbooks/using-main.version/

func main() {
	var debug bool
	var printVersion bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&printVersion, "version", false, "print the version")
	flag.Parse()

	if printVersion {
		log.Printf("Version: %s\n", version)

		// Create a new provider instance to get access to its resources
		log.Println("Registered resources:")
		p := provider.New(version)
		ctx := context.Background()
		resources := p.Resources(ctx)

		// Extract resource names from the resource functions
		var resourceTypes []string

		for _, resourceFunc := range resources {
			// Get the function name which typically includes the resource name
			funcName := getFunctionName(resourceFunc)

			// Extract resource name using the last component
			parts := strings.Split(funcName, ".")
			if len(parts) > 0 {
				lastPart := parts[len(parts)-1]

				// Handle NewResourceXXX format
				if strings.HasPrefix(lastPart, "NewResource") {
					resourceName := strings.TrimPrefix(lastPart, "NewResource")
					resourceName = toSnakeCase(resourceName)
					resourceTypes = append(resourceTypes, "fabric_"+resourceName)
				} else if strings.HasPrefix(lastPart, "NewDataSource") {
					// Handle datasource functions if present
					resourceName := strings.TrimPrefix(lastPart, "NewDataSource")
					resourceName = toSnakeCase(resourceName)
					resourceTypes = append(resourceTypes, "fabric_"+resourceName)
				}
			}
		}

		// Sort the resource types for consistent output
		sort.Strings(resourceTypes)

		// Print all resource types
		for _, resourceType := range resourceTypes {
			log.Printf("  - %s\n", resourceType)
		}

		return
	}

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/microsoft/fabric",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.NewFunc(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// getFunctionName returns the name of a function
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// toSnakeCase converts a camel case string to snake case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
