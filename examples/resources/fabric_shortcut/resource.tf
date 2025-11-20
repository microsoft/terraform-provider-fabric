# Example of using the fabric_shortcut resource
resource "fabric_shortcut" "onelake" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    onelake = {
      workspace_id = "00000000-0000-0000-0000-000000000000"
      item_id      = "00000000-0000-0000-0000-000000000000"
      find         = "MyTargetPath"
    }
  }
}

resource "fabric_shortcut" "adls_gen2" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    adls_gen2 = {
      location      = "MyLocation"
      subfind       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "amazon_s3" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    amazon_s3 = {
      location      = "MyLocation"
      subfind       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "google_cloud_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    google_cloud_storage = {
      location      = "MyLocation"
      subfind       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "s3_compatible" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    s3_compatible = {
      location      = "MyLocation"
      bucket        = "MyBucket"
      subfind       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "dataverse" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    dataverse = {
      table_name         = "MyTableName"
      deltalake_folder   = "MyDeltaLakeFolder"
      environment_domain = "MyEnvironmentDomainURI"
      bucket             = "MyBucket"
      subfind            = "MySubpath"
      connection_id      = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_shortcut" "azure_blob_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  find         = "MyShortcutPath"
  target = {
    azure_blob_storage = {
      location      = "MyLocation"
      subfind       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}
