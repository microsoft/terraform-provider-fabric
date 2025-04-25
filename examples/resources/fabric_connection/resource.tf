# Configure a SQL connection using a virtual network gateway with Basic authentication
resource "fabric_connection" "sql_connection_basic" {
  display_name = "ContosoVirtualNetworkGatewayConnection"

  # Connection configuration
  connectivity_type = "VirtualNetworkGateway"
  gateway_id        = "93491300-cfbd-402f-bf17-9ace59a92354"
  privacy_level     = "Organizational"

  connection_details {
    type            = "SQL"
    creation_method = "SQL"

    parameters {
      name      = "server"
      data_type = "Text"
      value     = "contoso.database.windows.net"
    }

    parameters {
      name      = "database"
      data_type = "Text"
      value     = "sales"
    }
  }

  credential_details {
    single_sign_on_type   = "None"
    connection_encryption = "Encrypted"
    skip_test_connection  = false

    credentials {
      credential_type = "Basic"
      username        = "admin"
      password        = "your-password-here" # Sensitive value, better to use variables
    }
  }
}

# Configure an Azure SQL connection using a virtual network gateway with Service Principal authentication
resource "fabric_connection" "sql_connection_service_principal" {
  display_name = "AzureSQLVirtualNetworkGatewayConnectionWithSPN"

  # Connection configuration
  connectivity_type = "VirtualNetworkGateway"
  gateway_id        = "93491300-cfbd-402f-bf17-9ace59a92354"
  privacy_level     = "Organizational"

  connection_details {
    type            = "AzureSQL"
    creation_method = "AzureSQL"

    parameters {
      name      = "server"
      data_type = "Text"
      value     = "azuresql.database.windows.net"
    }

    parameters {
      name      = "database"
      data_type = "Text"
      value     = "analytics"
    }
  }

  credential_details {
    single_sign_on_type   = "MicrosoftEntraID"
    connection_encryption = "Encrypted"
    skip_test_connection  = false

    credentials {
      credential_type    = "ServicePrincipal"
      application_id     = "87654321-4321-4321-4321-210987654321"
      application_secret = "your-service-principal-secret" # Sensitive value, better to use variables
      tenant_id          = "12345678-1234-1234-1234-123456789012"
    }
  }
}
