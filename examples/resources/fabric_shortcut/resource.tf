# Example of using the fabric_shortcut resource
resource "fabric_shortcut" "onelake" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    onelake = {
      workspace_id = "00000000-0000-0000-0000-000000000000"
      item_id      = "00000000-0000-0000-0000-000000000000"
      path         = "MyTargetPath"
    }
  }
}

resource "fabric_shortcut" "adls_gen2" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    adls_gen2 = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "amazon_s3" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    amazon_s3 = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "google_cloud_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    google_cloud_storage = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "s3_compatible" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    s3_compatible = {
      location      = "MyLocation"
      bucket        = "MyBucket"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "dataverse" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    dataverse = {
      table_name         = "MyTableName"
      deltalake_folder   = "MyDeltaLakeFolder"
      environment_domain = "MyEnvironmentDomainURI"
      bucket             = "MyBucket"
      subpath            = "MySubpath"
      connection_id      = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "azure_blob_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    azure_blob_storage = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}
