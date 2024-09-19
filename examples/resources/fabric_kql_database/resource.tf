# Create a ReadWrite KQL database example
resource "fabric_kql_database" "example1" {
  display_name = "example1"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type = "ReadWrite"
    eventhouse_id = "11111111-1111-1111-1111-111111111111"
  }
}

# Create a Shortcut KQL database to source Azure Data Explorer cluster example
resource "fabric_kql_database" "example2" {
  display_name = "example2"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type        = "Shortcut"
    eventhouse_id        = "11111111-1111-1111-1111-111111111111"
    source_cluster_uri   = "https://clustername.westus.kusto.windows.net"
    source_database_name = "MyDatabase"
  }
}

# Create a Shortcut KQL database to source Azure Data Explorer cluster with invitation token example
resource "fabric_kql_database" "example3" {
  display_name = "example3"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type    = "Shortcut"
    eventhouse_id    = "11111111-1111-1111-1111-111111111111"
    invitation_token = "eyJ0...InvitationToken...iJKV"
  }
}

# Create a Shortcut KQL database to source KQL database example
resource "fabric_kql_database" "example4" {
  display_name = "example4"
  workspace_id = "00000000-0000-0000-0000-000000000000"

  configuration = {
    database_type        = "Shortcut"
    eventhouse_id        = "11111111-1111-1111-1111-111111111111"
    source_database_name = "MyDatabase"
  }
}
