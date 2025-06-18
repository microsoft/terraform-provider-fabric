# Copyright (c) Microsoft Corporation
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_version = ">= 1.8, < 2.0"
  required_providers {
    fabric = {
      source  = "microsoft/fabric"
      version = "1.2.0"
    }
  }
}

provider "fabric" {
  preview = true
}


resource "fabric_ml_experiment" "example" {
  display_name = "example-ml-model"
  description  = "This is an example ML experiment."
  workspace_id = "b1ed84ce-1876-4be9-a6f9-15ba5557c6c4"
}

resource "fabric_ml_experiment" "example2" {
  display_name = "example-ml-model2"
  description  = "This is an example ML experiment 2."
  workspace_id = "b1ed84ce-1876-4be9-a6f9-15ba5557c6c4"
}
