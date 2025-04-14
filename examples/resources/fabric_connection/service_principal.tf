# Configure a connection using Service Principal authentication
resource "fabric_connection" "spn_connection" {
  workspace_id = "12345678-1234-1234-1234-123456789012"
  display_name = "ServicePrincipalConnection"
  description  = "Connection using Service Principal authentication"

  # Connection configuration
  properties {
    connectivity_type = "ShareableCloud"
    privacy_level     = "Organizational"
    
    connection_details {
      type            = "AzureDataLakeStorage"
      creation_method = "AzureDataLakeStorage"
      
      parameters {
        name      = "url"
        data_type = "Text"
        value     = "https://yourdatalake.dfs.core.windows.net"
      }
    }
    
    credential_details {
      single_sign_on_type    = "None"
      connection_encryption  = "Encrypted"
      skip_test_connection   = false
      
      credentials {
        credential_type    = "ServicePrincipal"
        application_id     = "87654321-4321-4321-4321-210987654321"
        application_secret = "your-app-secret-here" # Sensitive value, better to use variables
        tenant_id          = "00000000-0000-0000-0000-000000000000"
      }
    }
  }
}