# Example 1 - ShareableCloud Connection
resource "fabric_connection" "example_cloud" {
  display_name      = "example"
  connectivity_type = "ShareableCloud"
  privacy_level     = "Organizational"
  connection_details = {
    type            = "FTP"
    creation_method = "FTP.Contents"
    parameters = [
      {
        name  = "server"
        value = "ftp.example.com"
      }
    ]
  }
  credential_details = {
    connection_encryption = "NotEncrypted"
    credential_type       = "Basic"
    single_sign_on_type   = "None"
    skip_test_connection  = false
    basic_credentials = {
      username            = "user"
      password_wo         = "...secret_password..."
      password_wo_version = 1
    }
  }
}

# Example 2 - VirtualNetworkGateway Connection
resource "fabric_connection" "example_virtual_network_gateway" {
  gateway_id        = "00000000-0000-0000-0000-000000000000"
  display_name      = "example"
  connectivity_type = "VirtualNetworkGateway"
  privacy_level     = "Organizational"
  connection_details = {
    type            = "FTP"
    creation_method = "FTP.Contents"
    parameters = [
      {
        name  = "server"
        value = "ftp.example.com"
      }
    ]
  }
  credential_details = {
    connection_encryption = "NotEncrypted"
    credential_type       = "Basic"
    single_sign_on_type   = "None"
    skip_test_connection  = false
    basic_credentials = {
      username            = "user"
      password_wo         = "...secret_password..."
      password_wo_version = 1
    }
  }
}

# OAuth2 credential_type is not supported in the provider
