# Example 1 - Notebook without definition
resource "fabric_notebook" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  folder_id    = "11111111-1111-1111-1111-111111111111"
}

# Example 2 - Notebook with definition bootstrapping only
resource "fabric_notebook" "example_definition_bootstrap" {
  display_name              = "example"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "ipynb"
  definition = {
    "notebook-content.ipynb" = {
      source = "${local.path}/notebook.ipynb.tmpl"
    }
  }
}

# Example 3 - Notebook with definition update when source or tokens changed
resource "fabric_notebook" "example_definition_update" {
  display_name = "example"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "ipynb"
  definition = {
    "notebook-content.ipynb" = {
      source = "${local.path}/notebook.ipynb.tmpl"
      tokens = {
        "MESSAGE" = "World"
        "MyValue" = "Lorem Ipsum"
      }
    }
  }
}
