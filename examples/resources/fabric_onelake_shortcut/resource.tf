# Example of using the fabric_onelake_shortcut resource
resource "fabric_onelake_shortcut" "onelake" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    onelake = {
      workspace_id = "00000000-0000-0000-0000-000000000000"
      item_id      = "00000000-0000-0000-0000-000000000000"
      path         = "/MyTargetPath"
    }
  }
}

resource "fabric_onelake_shortcut" "adlsgen2" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    adlsGen2 = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_onelake_shortcut" "amazon_s3" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    amazons3 = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_onelake_shortcut" "google_cloud_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    google_cloud_storage = {
      location      = "MyLocation"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_onelake_shortcut" "s3_compatible" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    s3_compatible = {
      location      = "MyLocation"
      bucket        = "MyBucket"
      subpath       = "MySubpath"
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_onelake_shortcut" "dataverse" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    dataverse = {
      table_name         = "MyTableName"
      deltaLake_folder   = "MyDeltaLakeFolder"
      environment_domain = "MyEnvironmentDomainURI"
      bucket             = "MyBucket"
      subpath            = "MySubpath"
      connection_id      = "00000000-0000-0000-0000-000000000000"
    }
  }
}

resource "fabric_onelake_shortcut" "external_data_share_target" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "00000000-0000-0000-0000-000000000000"
  name         = "MyShortcutName"
  path         = "/MyShortcutPath"
  target = {
    external_data_share_target = {
      connection_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}
