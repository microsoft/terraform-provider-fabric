resource "fabric_workspace_encryption" "example" {
  workspace_id   = "00000000-0000-0000-0000-000000000000"
  key_identifier = "https://example-vault.vault.azure.net/keys/example-key/"
}
