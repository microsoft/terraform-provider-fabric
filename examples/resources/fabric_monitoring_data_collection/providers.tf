terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source = "microsoft/fabric"
    }
  }
}

provider "fabric" {
  preview = true
}
