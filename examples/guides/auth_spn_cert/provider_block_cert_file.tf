variable "client_certificate" {
  description = "The path to the client certificate file."
  type        = string
}
variable "client_certificate_password" {
  description = "The password for the client certificate."
  type        = string
  sensitive   = true
}

# We strongly recommend using the required_providers block to set the Fabric Provider source and version being used
terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "0.0.0" # Check for the latest version on the Terraform Registry
    }
  }
}

# Configure the Microsoft Fabric Provider
provider "fabric" {
  tenant_id                    = "00000000-0000-0000-0000-000000000000"
  client_id                    = "00000000-0000-0000-0000-000000000000"
  client_certificate_file_path = var.client_certificate
  client_certificate_password  = var.client_certificate_password
}
