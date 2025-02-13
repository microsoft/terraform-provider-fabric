---
page_title: "Getting started with the Terraform Provider for Microsoft Fabric"
subcategory: "Guides"
description: |-

---

# Getting started

[HashiCorp Terraform](https://www.terraform.io/) is a popular open-source tool for creating safe and predictable cloud infrastructure across several cloud providers.
You can use the Microsoft Fabric Terraform Provider to manage your Microsoft Fabric workspaces using a flexible, powerful tool.
The goal of the Microsoft Fabric Terraform Provider is to support automation of the most complicated aspects of deploying and managing Microsoft Fabric.
Microsoft Fabric customers are using the Microsoft Fabric Terraform Provider to deploy and manage clusters and jobs and to configure data access.

In this section, you install and configure requirements to use Terraform and the Microsoft Fabric Terraform Provider on your local development machine.
You then configure Terraform authentication. Following this section, this article provides a sample configuration that you can experiment with to provision a Microsoft Fabric Notebook and Lakehouse.

## Requirements

1. You must have the Terraform CLI. See [Download Terraform](https://www.terraform.io/downloads.html) on the Terraform website.

1. You must have a Terraform project. In your terminal, create an empty directory and then switch to it. (Each separate set of Terraform configuration files must be in its own directory, which is called a Terraform project.) For example:

    ```bash
    mkdir terraform_demo && cd terraform_demo
    ```

    Include Terraform configurations for your project in one or more configuration files in your Terraform project. For information about the configuration file syntax, see [Terraform Language Documentation](https://developer.hashicorp.com/terraform/language) on the Terraform website.

1. You must configure authentication for your Terraform project. See [Authentication](./auth_app_reg_spn.md) in the Microsoft Fabric Terraform Provider documentation.
1. You must have a [Fabric Capacity](https://learn.microsoft.com/fabric/enterprise/licenses#capacity) provisioned in Azure. See [Configuring a Fabric Capacity](./fabric_capacity_setup.md) in the Microsoft Fabric Terraform Provider documentation.

-> Please keep the capacity name handy, as we will use it below to fetch the capacity id.

## Sample configuration

This section provides a sample configuration that you can experiment with to provision a Microsoft Fabric Notebook and a Lakehouse. It assumes that you have already set up the requirements, as well as created a Terraform project and configured the project with Terraform authentication as described in the previous section.

1. Create a new file named `provider.tf` in your Terraform project directory.
1. Add the following code to `provider.tf` to define a dependency on the Microsoft Fabric Terraform Provider:

    ```terraform
    # We strongly recommend using the required_providers block to set the Fabric Provider source and version being used
    terraform {
      required_version = ">= 1.8, < 2.0"
      required_providers {
        fabric = {
          source  = "microsoft/fabric"
          version = "0.1.0-beta.9"
        }
      }
    }

    # Configure the Microsoft Fabric Terraform Provider
    provider "fabric" {
      # Configuration options
    }
    ```

1. Create another file named `variables.tf`, and add the following code. This file represents input variables that can be used to configure a Notebook and Lakehouse.

    ```terraform
    variable "workspace_display_name" {
      description = "A name for the getting started workspace."
      type        = string
    }

    variable "notebook_display_name" {
      description = "A name for the subdirectory to store the notebook."
      type        = string
    }

    variable "notebook_definition_update_enabled" {
      description = "Whether to update the notebook definition."
      type        = bool
      default     = true
    }

    variable "notebook_definition_path" {
      description = "The path to the notebook's definition file."
      type        = string
    }

    variable "capacity_name" {
      description = "The name of the capacity to use."
      type = string
    }
    ```

1. Create a file named `workspace.tf` and add the following hcl code to represent a Microsoft Fabric workspace. We will also add a data source to fetch the Microsoft Fabric Capacity id by name (see requirements section).

    ```terraform
    data "fabric_capacity" "capacity" {
      display_name = var.capacity_name
    }

    resource "fabric_workspace" "example_workspace" {
      display_name = var.workspace_display_name
      description = "Getting started workspace"
      capacity_id = data.fabric_capacity.capacity.id
    }
    ```

1. Create a file named notebook.ipynb in the same folder and copy the content of [this example notebook](https://github.com/Azure-Samples/modern-data-warehouse-dataops/blob/main/single_tech_samples/fabric/fabric_ci_cd/src/notebooks/nb-city-safety.ipynb).
1. Create a file named `notebook.tf` and add the following hcl code to represent a Notebook. This Notebook references the workspace created in step 4, specifically using the workspace id.

    ```terraform
    resource "fabric_notebook" "example_notebook" {
      workspace_id = fabric_workspace.example_workspace.id
      display_name = var.notebook_display_name
      definition_update_enabled = var.notebook_definition_update_enabled
      definition = {
        "notebook-content.ipynb" = {
          source = var.notebook_definition_path
        }
      }
    }
    ```

1. Create another file named `terraform.tfvars`, and add the following code. This file specifies the Notebook's properties. Learn more about [tfvars file](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files).

    ```terraform
    workspace_display_name = "example workspace"
    notebook_display_name = "example notebook"
    notebook_definition_update_enabled = true
    notebook_definition_path = "notebook.ipynb"
    capacity_name = "<use capacity name configured above>"
    ```

1. Create another file named `outputs.tf`, this is where we will define [Terraform output values](https://developer.hashicorp.com/terraform/language/values/outputs). Add the following code:

    ```terraform
    output "capacity_id" {
      value = data.fabric_capacity.capacity.id
    }

    output "notebook_id" {
      value = fabric_notebook.example_notebook.id
    }
    ```

1. Run `terraform init`. If there are any errors, fix them, and then run the command again.
1. Run `terraform plan -out=plan.tfplan`. If there are any errors, fix them, and then run the command again. In this example, we are capturing the output to a plan file named `plan.tfplan`.
1. Run `terraform apply plan.tfplan`. This command applies the changes required to reach the desired state of the configuration. If there are any errors, fix them, and then run the command again.
1. Verify that the Workspace and Notebook were created in Microsoft Fabric. In the output of the `terraform apply` command, find the Notebook id and capacity id.
1. when you are done with this sample, delete the Notebook, and workspace from Microsoft Fabric by running `terraform destroy`.

## Troubleshooting

See [Troubleshooting](./troubleshooting.md) guide.
