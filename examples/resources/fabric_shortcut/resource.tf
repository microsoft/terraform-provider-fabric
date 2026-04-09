# Example of using the fabric_shortcut resource
resource "fabric_shortcut" "onelake" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    onelake = {
      workspace_id = "00000000-0000-0000-0000-000000000000"
      item_id      = "22222222-2222-2222-2222-222222222222"
      path         = "Tables/myTablesFolder/someTableSubFolder"
    }
  }
}

resource "fabric_shortcut" "adls_gen2" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    adls_gen2 = {
      location      = "https://[account-name].dfs.core.windows.net"
      subpath       = "[container]/[subfolder]"
      connection_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}

resource "fabric_shortcut" "amazon_s3" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    amazon_s3 = {
      location      = "https://[bucket-name].s3.[region-code].amazonaws.com"
      subpath       = "MySubpath"
      connection_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}

resource "fabric_shortcut" "google_cloud_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    google_cloud_storage = {
      location      = "https://[bucket-name].storage.googleapis.com"
      subpath       = "/folder"
      connection_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}

resource "fabric_shortcut" "s3_compatible" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    s3_compatible = {
      location      = "https://s3endpoint.contoso.com"
      bucket        = "MyBucket"
      subpath       = "/folder"
      connection_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}

resource "fabric_shortcut" "dataverse" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    dataverse = {
      table_name         = "MyTableName"
      deltalake_folder   = "MyDeltaLakeFolder"
      environment_domain = "https://[orgname].crm[xx].dynamics.com"
      connection_id      = "22222222-2222-2222-2222-222222222222"
    }
  }
}

resource "fabric_shortcut" "azure_blob_storage" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id      = "11111111-1111-1111-1111-111111111111"
  name         = "MyShortcutName"
  path         = "MyShortcutPath"
  target = {
    azure_blob_storage = {
      location      = "https://[account-name].blob.core.windows.net"
      subpath       = "/mycontainer/mysubfolder"
      connection_id = "22222222-2222-2222-2222-222222222222"
    }
  }
}
