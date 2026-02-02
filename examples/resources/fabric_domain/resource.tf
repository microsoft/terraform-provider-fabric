resource "fabric_domain" "parent" {
  display_name = "example parent"
}

resource "fabric_domain" "child" {
  display_name     = "example child"
  description      = "This is an example child domain"
  parent_domain_id = fabric_domain.parent.id
}

# Fabric Domain operations require admin API access and may fail if the Service Principal has the Tenant.ReadWrite.All permission assigned
