# Example 1 - Notebook without definition
resource "fabric_notebook" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
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

# Example 4 - Notebook with custom tokens delimiter
resource "fabric_notebook" "example_custom_delimiter" {
  display_name = "example"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "ipynb"
  definition = {
    "notebook-content.ipynb" = {
      source           = "${local.path}/notebook.ipynb.tmpl"
      tokens_delimiter = "##"
      tokens = {
        "MESSAGE" = "World"
        "MyValue" = "Lorem Ipsum"
      }
    }
  }
}

# Example 5 - Notebook with parameters processing mode
resource "fabric_notebook" "example_parameters" {
  display_name = "example"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "ipynb"
  definition = {
    "notebook-content.ipynb" = {
      source          = "${local.path}/notebook.ipynb.tmpl"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "JsonPathReplace"
          find  = "$.metadata.kernelspec.display_name"
          value = "Python 3.10"
        },
        {
          type  = "TextReplace"
          find  = "OldValue"
          value = "NewValue"
        }
      ]
    }
  }
}
