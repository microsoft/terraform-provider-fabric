---
page_title: "Configuring a Fabric Capacity"
subcategory: "Guides"
description: |-
  This guide outlines the steps to provision Microsoft Fabric Capacity.
---

# Configuring a Fabric Capacity in Azure and granting access

This guide outlines the steps to provision Microsoft Fabric Capacity in Azure and grant access to the Capacity for Service Principal (SPN), Managed Identity (MSI) or User to use with the Fabric Terraform Provider.

---

~> Be aware that Fabric capacities will start incurring costs upon creation. It is your responsibility to properly size and pause Fabric Capacities as needed, and to stay informed about costs. You are also responsible for any excessive expenses. Please [plan your capacity size](https://learn.microsoft.com/fabric/enterprise/plan-capacity) to fit your requirements.

---

## Provisioning Fabric Capacity

### Using Azure Portal

To add new capacity, visit the [Azure Portal](https://portal.azure.com/#browse/Microsoft.Fabric%2Fcapacities) and follow the wizard's steps.

### Using Terraform

Terraform can also help you create and manage a Fabric Capacity, but not through this specific provider, as it does not support Azure resources. This provider is designed exclusively for managing resources within Microsoft Fabric, which is separate from Azure. For Azure resources, consider using the [AzAPI provider](https://registry.terraform.io/providers/Azure/azapi) provider instead.

To assist you in getting started, check out some examples available at <https://aka.ms/FabricTF/quickstart>. This repository contains Fabric Capacity configuration for Terraform using AzAPI.

## Granting access to Fabric Capacity

Fabric Capacity can be accessed by Service Principals (SPN), Managed Identities (MSI) and Users. The following sections explain how to grant each entity access. This step is necessary for the Fabric Terraform Provider to use `fabric_capacity` / `fabric_capacities` Data-Sources or to assign capacity to a Workspace using the `fabric_workspace` Resource.

The steps outlined below illustrate the manual method for granting access to the Capacity via Azure Portal. Automation of this process is also possible using various tools such as [Azure REST API](https://learn.microsoft.com/rest/api/microsoftfabric/fabric-capacities/update), [Azure CLI](https://learn.microsoft.com/cli/azure/), [Azure PowerShell](https://learn.microsoft.com/powershell/azure/), or Terraform (with [AzAPI](https://registry.terraform.io/providers/Azure/azapi) and [AzureAD](https://registry.terraform.io/providers/hashicorp/azuread) providers).

-> If you choose to automate this section, ensure you use the principal ID (Enterprise Application Object ID) for Service Principals (SPN) and Managed Identities (MSI). For Users, please use User Principal Name (UPN).

### Service Principals (SPN) and Managed Identities (MSI)

1. Sign in to the [Azure Portal](https://portal.azure.com/).
1. Browse to the Fabric Capacity you want to grant access to.
1. Go to **Settings** > **Capacity administrators**.
1. Click **Add** and select the **Enterprise applications**
1. Search for the Service Principal or Managed Identity you want to grant access to and click **Select**.
1. Click **Save**.

### Users

1. Sign in to the [Azure Portal](https://portal.azure.com/).
1. Browse to the Fabric Capacity you want to grant access to.
1. Go to **Settings** > **Capacity administrators**.
1. Click **Add** and select the **Users**
1. Search for the User you want to grant access to and click **Select**.
1. Click **Save**.

## Configure Microsoft Fabric to allow Service Principals or Managed Identities

Please follow [Configure Microsoft Fabric to allow Service Principals (SPN) and Managed Identities (MSI)](./auth_app_reg_spn.md#configure-microsoft-fabric-to-allow-service-principals-spn-and-managed-identities-msi) guide.
