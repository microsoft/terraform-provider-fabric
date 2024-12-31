resource "fabric_connection" "example" {
  display_name      = "example parent"
  connectivity_type = "ShareableCloud"
  privacy_level     = "Organizational"
  connection_details = {
    type            = "FTP"
    creation_method = "FTP.Contents"
    parameters = {
      "server" = "ftp.example.com"
    }
  }
  credential_details = {
    connection_encryption = "NotEncrypted"
    credential_type       = "Basic"
    single_sign_on_type   = "None"
    skip_test_connection  = false
    basic_credentials = {
      username = "user"
      password = "password"
    }
  }
}
