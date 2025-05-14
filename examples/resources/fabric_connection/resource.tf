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

# Example 2 - OnPremisesGateway Connection
# resource "fabric_connection" "example_on_premises_gateway" {
#   display_name      = "example"
#   connectivity_type = "OnPremisesGateway"
#   privacy_level     = "Organizational"
#   connection_details = {
#     type            = "FTP"
#     creation_method = "FTP.Contents"
#     parameters = [
#       {
#         name  = "server"
#         value = "ftp.example.com"
#       }
#     ]
#   }
#   credential_details = {
#     connection_encryption = "NotEncrypted"
#     credential_type       = "Basic"
#     single_sign_on_type   = "None"
#     skip_test_connection  = false
#     basic_credentials = {
#       username            = "user"
#       password_wo         = "...secret_password..."
#       password_wo_version = 1
#     }
#   }
# }
