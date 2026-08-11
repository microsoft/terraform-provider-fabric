# Virtual Network Gateway
resource "fabric_gateway" "example" {
  type                            = "VirtualNetwork"
  display_name                    = "example"
  inactivity_minutes_before_sleep = 30
  number_of_member_gateways       = 1
  virtual_network_azure_resource = {
    resource_group_name  = "example resource group"
    virtual_network_name = "example virtual network"
    subnet_name          = "example subnet"
    subscription_id      = "00000000-0000-0000-0000-000000000000"
  }
  capacity_id = "11111111-1111-1111-1111-111111111111"
}

# Virtual Network Gateway with an autoscaling member gateway range
resource "fabric_gateway" "example_autoscale" {
  type                            = "VirtualNetwork"
  display_name                    = "example autoscale"
  inactivity_minutes_before_sleep = 30

  # mutually exclusive with number_of_member_gateways; both bounds must be set together
  min_member_gateway_count = 1
  max_member_gateway_count = 5

  virtual_network_azure_resource = {
    resource_group_name  = "example resource group"
    virtual_network_name = "example virtual network"
    subnet_name          = "example subnet"
    subscription_id      = "00000000-0000-0000-0000-000000000000"
  }
  capacity_id = "11111111-1111-1111-1111-111111111111"
}

# Streaming Virtual Network Gateway
resource "fabric_gateway" "example_streaming" {
  type         = "StreamingVirtualNetwork"
  display_name = "example streaming"
  virtual_network_azure_resource = {
    resource_group_name  = "example resource group"
    virtual_network_name = "example virtual network"
    subnet_name          = "example subnet"
    subscription_id      = "00000000-0000-0000-0000-000000000000"
  }
}
