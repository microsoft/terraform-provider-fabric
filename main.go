// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

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
