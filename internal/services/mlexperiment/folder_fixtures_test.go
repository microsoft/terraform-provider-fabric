// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package mlexperiment_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabfake "github.com/microsoft/fabric-sdk-go/fabric/fake"

	"github.com/microsoft/terraform-provider-fabric/internal/testhelp"
	"github.com/microsoft/terraform-provider-fabric/internal/testhelp/fakes"
)

func TestUnit_FolderFixtures_ExperimentLifecycle(t *testing.T) {
	t.Setenv("FABRIC_TESTACC_WELLKNOWN", `{"FolderRS1":{"id":"stale-first"},"FolderRS2":{"id":"stale-second"}}`)

	server := fakes.NewFakeServer()
	// Generic item fakes cannot store folder entities in their item conversion map.
	server.ServerFactory.Core.FoldersServer = fakes.NewFakeServer().ServerFactory.Core.FoldersServer
	testCase := testhelp.NewTestUnitCase(t, &testResourceItemFQN, server.ServerFactory, nil,
		experimentCRUDSteps(t, testhelp.RandomUUID()))
	testCase.ProtoV6ProviderFactories = testhelp.GetTestUnitProtoV6ProviderFactories(&fabric.ClientOptions{
		Transport: &folderFixtureTransport{transport: fabfake.NewServerFactoryTransport(server.ServerFactory)},
	}, nil)
	resource.Test(t, testCase)
}

// The in-memory fake stores need serialized access, even when Terraform runs resources concurrently.
type folderFixtureTransport struct {
	mutex     sync.Mutex
	transport interface {
		Do(request *http.Request) (*http.Response, error)
	}
}

func (transport *folderFixtureTransport) Do(request *http.Request) (*http.Response, error) {
	transport.mutex.Lock()
	defer transport.mutex.Unlock()

	return transport.transport.Do(request)
}

func TestUnit_FolderFixtures_ExperimentCRUD(t *testing.T) {
	t.Setenv("FABRIC_TESTACC_WELLKNOWN", `{"FolderRS1":{"id":"stale-first"},"FolderRS2":{"id":"stale-second"}}`)
	const workspaceID = "00000000-0000-0000-0000-000000000001"
	seen := map[string]bool{}

	for run := range 2 {
		steps := experimentCRUDSteps(t, workspaceID)
		if len(steps) != 2 {
			t.Fatalf("expected create and move steps, got %d", len(steps))
		}
		var previous map[string]bool
		var previousTarget string

		for step, config := range steps {
			file, diags := hclsyntax.ParseConfig([]byte(config.Config), "fixture.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatal(diags)
			}

			folders := map[string]bool{}
			var target string

			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				if block.Type != "resource" {
					continue
				}

				if block.Labels[0] == "fabric_folder" {
					folders[block.Labels[1]] = true

					value, diagnostics := block.Body.Attributes["workspace_id"].Expr.Value(nil)
					if diagnostics.HasErrors() || value.AsString() != workspaceID {
						t.Fatal("owned folder must use the requested workspace")
					}
				}

				if block.Labels[0] == "fabric_ml_experiment" {
					refs := block.Body.Attributes["folder_id"].Expr.Variables()
					if len(refs) != 1 || len(refs[0]) != 3 || refs[0].RootName() != "fabric_folder" {
						t.Fatal("item must reference an owned folder, not a static fixture ID")
					}

					target = refs[0][1].(hcl.TraverseAttr).Name
					if refs[0][2].(hcl.TraverseAttr).Name != "id" {
						t.Fatal("item must depend on the folder ID")
					}
				}
			}

			if len(folders) != 2 || !folders[target] {
				t.Fatal("both owned folders must remain in every CRUD step")
			}

			if step == 0 {
				for name := range folders {
					if seen[name] {
						t.Fatalf("run %d reused another run's folder", run)
					}

					seen[name] = true
				}
			} else {
				for name := range previous {
					if !folders[name] {
						t.Fatal("moving the item must not destroy its previous folder")
					}
				}

				if target == previousTarget {
					t.Fatal("second step must move to the other owned folder")
				}
			}

			previous, previousTarget = folders, target
		}
	}
}
