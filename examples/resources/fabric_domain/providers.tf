terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "1.0.0" # Check for the latest version on the Terraform Registry
    }
  }
}

provider "fabric" {}
