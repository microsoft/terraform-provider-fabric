resource "fabric_connection" "example" {
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
