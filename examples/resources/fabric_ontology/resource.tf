# Ontology generation 2 (new experience) only. REST API v1 is still used.
# A create without a definition defaults to generation 2.
resource "fabric_ontology" "example" {
  display_name = "OntologyExample"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Bootstrap a generation 2 ontology from its minimal TMDL definition.
# The service synthesizes model.tmdl and namespaces/default.tmdl when omitted.
resource "fabric_ontology" "example_definition_bootstrap" {
  display_name              = "OntologyBootstrap"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "TMDL"
  definition = {
    ".platform" = {
      source = "${path.module}/definition/platform.json"
    }
    "database.tmdl" = {
      source = "${path.module}/definition/database.tmdl"
    }
  }
}

# Manage the definition, including changes to source files and template tokens.
# TMDL files are passed through as text; no JSON conversion is required.
resource "fabric_ontology" "example_definition_update" {
  display_name = "OntologyManaged"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "TMDL"
  definition = {
    ".platform" = {
      source = "${path.module}/definition/platform.json"
    }
    "database.tmdl" = {
      source = "${path.module}/definition/database.tmdl"
    }
    "model.tmdl" = {
      source = "${path.module}/definition/model.tmdl.tmpl"
      tokens = {
        Namespace = "operations"
      }
    }
    "namespaces/default.tmdl" = {
      source = "${path.module}/definition/namespaces/default.tmdl"
    }
    "namespaces/operations.tmdl" = {
      source = "${path.module}/definition/namespaces/namespace.tmdl.tmpl"
      tokens = {
        Namespace = "operations"
      }
    }
  }
}
