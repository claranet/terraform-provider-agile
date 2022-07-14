terraform {
  required_providers {
    agile = {
      source = "claranet/agile"
    }
  }
}

provider "agile" {}

resource "agile_logical_switch" "example" {
  name             = "example"
  description      = "This Logical Switch is created by terraform"
  logic_network_id = "5308df55-1709-404f-b4f8-4d8947d8f0c4"
  tenant_id        = "cd27d9cf-9be0-4852-a560-2d6e05fd3c1e"
  mac_address      = "00:00:5E:00:01:02"

  storm_suppress {
    broadcast_enable   = true
    multicast_enable   = true
#    unicast_enable     = true
    broadcast_cbs      = 10000
    broadcast_cbs_unit = "byte"
    broadcast_cir      = 100
    broadcast_cir_unit = "kbps"
    multicast_cbs      = 10000
    multicast_cbs_unit = "byte"
    multicast_cir      = 100
    multicast_cir_unit = "kbps"
#    unicast_cbs        = 10000
#    unicast_cbs_unit   = "byte"
#    unicast_cir        = 100
#    unicast_cir_unit   = "kbps"
  }

  additional {
    producer = "Terraform"
  }

}

output "id" {
  value = agile_logical_switch.example.id
}