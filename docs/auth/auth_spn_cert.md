---
page_title: "Authenticating using a Service Principal and Client Certificate"
subcategory: "Authentication"
description: |-

---

# Authenticating using a Service Principal and Client Certificate

---

!> We recommend using either a [Service Principal with OpenID Connect (OIDC)](./auth_spn_oidc.md) or [Managed Service Identity (MSI)](./auth_msi.md) authentication when running Terraform non-interactively (such as when running Terraform in a CI server), and authenticating using the Azure CLI when running Terraform locally.

---

## Setting up Entra Application and Service Principal

Follow [Creating an App Registration for the Service Principal context (SPN)](./auth_app_reg_spn.md) guide.

## Generating Client Certificate

Firstly we need to create a certificate which can be used for authentication. To do that we're going to generate a private key and self-signed certificate using OpenSSL or LibreSSL (this can also be achieved using PowerShell, however that's outside the scope of this document).

```shell
# Generates RSA 4096-bit private key with AES-256 encryption
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -aes256 -pass pass:'YourPrivateKeyPassword' -out private.key

# Generates a self-signed certificate valid for 1 year
openssl req -subj '/CN=myclientcert/O=Contoso Inc./ST=WA/C=US' -x509 -sha256 -days 365 -passin pass:'YourPrivateKeyPassword' -key private.key -out client.crt

# Generates a PKCS12 bundle from a private key and a certificate
openssl pkcs12 -export -passin pass:'YourPrivateKeyPassword' -password pass:'YourBundlePassword' -inkey private.key -in client.crt -out bundle.pfx
```

## Adding a Client Certificate to Entra App

### Using Azure Portal

1. In the [Microsoft Entra admin center](https://entra.microsoft.com), in **App registrations**, select your application.
1. Select **Certificates & secrets** > `Certificates` > `Upload certificate`.
1. Select the file you want to upload. It must be one of the following file types: *.cer, .pem, .crt*.
1. Select `Add`.

### Using Azure CLI

Run the following command to upload the certificate:

```shell
# See https://learn.microsoft.com/cli/azure/ad/app/credential#az-ad-app-credential-reset for more details.
az ad app credential reset --id "00000000-0000-0000-0000-000000000000" --append --cert "@~/client.crt"
```

Replace `00000000-0000-0000-0000-000000000000` with the Application (client) ID of the service principal.

### Using Azure PowerShell

Run the following command to upload the certificate:

```powershell
$certValue = [System.Convert]::ToBase64String([System.IO.File]::ReadAllBytes('client.crt'))
Get-AzADApplication -ApplicationId '00000000-0000-0000-0000-000000000000' | New-AzADAppCredential -CertValue $certValue
```

Replace `00000000-0000-0000-0000-000000000000` with the Application (client) ID of the service principal.

## Configuring Terraform to use the Client Certificate

Now that we have our Client Certificate uploaded to Entra App and ready to use, it's possible to configure Terraform in a few different ways.

The provider can be configured to read the certificate bundle from the `.pfx` file in your filesystem, or alternatively you can pass a base64-encoded copy of the certificate bundle directly to the provider.

### Environment Variables

Our recommended approach is storing the credentials as Environment Variables, for example:

#### Reading the certificate bundle from the filesystem (env vars)

```shell
# sh
export FABRIC_TENANT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_CERTIFICATE_FILE_PATH="/path/to/my/client/bundle.pfx"
export FABRIC_CLIENT_CERTIFICATE_PASSWORD="YourBundlePassword"
```

```powershell
# PowerShell
$env:FABRIC_TENANT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_CERTIFICATE_FILE_PATH = 'C:\Users\myusername\Documents\my\client\bundle.pfx'
$env:FABRIC_CLIENT_CERTIFICATE_PASSWORD = 'YourBundlePassword'
```

#### Passing the encoded certificate bundle directly (env vars)

```shell
# sh
export FABRIC_TENANT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export FABRIC_CLIENT_CERTIFICATE="$(base64 -w0 /path/to/my/client/bundle.pfx)"
export FABRIC_CLIENT_CERTIFICATE_PASSWORD="YourBundlePassword"
```

```powershell
# PowerShell
$env:FABRIC_TENANT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_ID = '00000000-0000-0000-0000-000000000000'
$env:FABRIC_CLIENT_CERTIFICATE = [System.Convert]::ToBase64String([System.IO.File]::ReadAllBytes('C:\Users\myusername\Documents\my\client\bundle.pfx'))
$env:FABRIC_CLIENT_CERTIFICATE_PASSWORD = 'YourBundlePassword'
```

The following Terraform and Provider blocks can be specified, where `0.0.0-preview` is the version of the Fabric Provider that you'd like to use:

```terraform
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
provider "fabric" {}
```

### Provider Block

It's also possible to configure these variables either directly or from variables in your provider block.

#### Reading the certificate bundle from the filesystem (provider block)

The following Terraform and Provider blocks can be specified, where `0.0.0-preview` is the version of the Fabric Provider that you'd like to use:

```terraform
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
```

#### Passing the encoded certificate bundle directly (provider block)

The following Terraform and Provider blocks can be specified, where `0.0.0-preview` is the version of the Fabric Provider that you'd like to use:

```terraform
variable "client_certificate" {
  description = "The client certificate."
  type        = string
  sensitive   = true
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
  tenant_id                   = "00000000-0000-0000-0000-000000000000"
  client_id                   = "00000000-0000-0000-0000-000000000000"
  client_certificate          = var.client_certificate
  client_certificate_password = var.client_certificate_password
}
```
