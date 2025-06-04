# We strongly recommend using the required_providers block to set the Fabric Provider source and version being used
terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "1.1.0"
    }
  }
}

# Configure the Microsoft Fabric Terraform Provider
provider "fabric" {
  # Configuration options
}

ephemeral "fabric_eventstream_source_connection" "example" {
  workspace_id   = "3ee46a65-b214-4db1-87d6-cf7602d4bd39"
  eventstream_id = "ae6a8b1f-05c6-45bf-9b05-8afa20ac29fd"
  source_id      = "7f77b4a1-9989-4bae-8c9c-a5dd6ef62689"
}

output "example" {
  value = ephemeral.fabric_eventstream_source_connection.example
  # ephemeral = true
  sensitive = true
}
