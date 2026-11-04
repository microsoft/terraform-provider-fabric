# Example 1 - User Data Function without definition
resource "fabric_user_data_function" "example" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"
}

# Example 2 - Item with definition bootstrapping only
resource "fabric_user_data_function" "example_definition_bootstrap" {
  display_name              = "example2"
  description               = "example with definition bootstrapping"
  workspace_id              = "00000000-0000-0000-0000-000000000000"
  definition_update_enabled = false
  format                    = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json"
    }
    "function_app.py" = {
      source = "${local.path}/function_app.py"
    }
    "resources/functions.json" = {
      source = "${local.path}/resources/functions.json"
    }
  }
}

# Example 3 - Item with definition update when source or tokens changed
resource "fabric_user_data_function" "example_definition_update" {
  display_name = "example3"
  description  = "example with definition update when source or tokens changed"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json"
    }
    "function_app.py" = {
      source = "${local.path}/function_app.py.tmpl"
      tokens = {
        "FunctionName" = "hello_fabric"
        "Greeting"     = "Hello"
      }
    }
    "resources/functions.json" = {
      source = "${local.path}/resources/functions.json"
    }
  }
}

# Example 4 - Item with custom tokens delimiter
resource "fabric_user_data_function" "example_custom_delimiter" {
  display_name = "example4"
  description  = "example with custom tokens delimiter"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json"
    }
    "function_app.py" = {
      source           = "${local.path}/function_app-custom-delimiter.py.tmpl"
      tokens_delimiter = "<<>>"
      tokens = {
        "FunctionName" = "hello_fabric"
        "Greeting"     = "Hello"
      }
    }
    "resources/functions.json" = {
      source = "${local.path}/resources/functions.json"
    }
  }
}

# Example 5 - Item with parameters processing mode
resource "fabric_user_data_function" "example_parameters" {
  display_name = "example5"
  description  = "example with parameters processing mode"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  format       = "Default"
  definition = {
    "definition.json" = {
      source = "${local.path}/definition.json"
    }
    "function_app.py" = {
      source          = "${local.path}/function_app.py"
      processing_mode = "parameters"
      parameters = [
        {
          type  = "TextReplace"
          find  = "hello_fabric"
          value = "greet_fabric"
        },
        {
          type  = "TextReplace"
          find  = "Hello"
          value = "Welcome"
        }
      ]
    }
    "resources/functions.json" = {
      source = "${local.path}/resources/functions.json"
    }
  }
}
