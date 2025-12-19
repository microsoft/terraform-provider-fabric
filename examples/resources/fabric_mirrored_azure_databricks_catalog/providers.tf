terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "0.0.0" # Check for the latest version on the Terraform Registry
    }
  }
}

provider "fabric" {}

locals {
  path = abspath(join("/", [path.root, "..", "..", "..", "internal", "testhelp", "fixtures", "mirrored_azure_databricks_catalog"]))
}
