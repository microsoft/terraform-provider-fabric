# Configure a connection using Workspace Identity authentication
resource "fabric_connection" "workspace_identity_connection" {
  workspace_id = "12345678-1234-1234-1234-123456789012"
  display_name = "WorkspaceIdentityConnection"
  description  = "Connection using Workspace Identity authentication"

  # Connection configuration
  properties {
    connectivity_type = "ShareableCloud"
    privacy_level     = "Organizational"
    
    connection_details {
      type            = "AzureBlobStorage"
      creation_method = "AzureBlobStorage"
      
      parameters {
        name      = "url"
        data_type = "Text"
        value     = "https://yourstorage.blob.core.windows.net"
      }
    }
    
    credential_details {
      single_sign_on_type    = "None"
      connection_encryption  = "Encrypted"
      skip_test_connection   = false
      
      credentials {
        credential_type = "WorkspaceIdentity"
        # No additional fields required for WorkspaceIdentity
      }
    }
  }
}