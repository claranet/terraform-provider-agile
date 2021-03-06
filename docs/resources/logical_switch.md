---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "agile_logical_switch Resource - terraform-provider-agile"
subcategory: ""
description: |-
  Manages Logical Switches.
---

# agile_logical_switch (Resource)

Manages Logical Switches.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `logic_network_id` (String) Logical network where a logical switch is located.

### Optional

- `additional` (Block List, Max: 1) Additional Settings. (see [below for nested schema](#nestedblock--additional))
- `bd` (Number) BD ID of a logical switch.
- `description` (String) Logical switch description.
- `mac_address` (String) MAC address of a logical switch.
- `name` (String) Logical switch name.
- `storm_suppress` (Block List, Max: 1) Storm Suppress Settings. (see [below for nested schema](#nestedblock--storm_suppress))
- `tenant_id` (String) Tenant ID. In the northbound direction, the value can be either specified or not. The controller can automatically obtain the tenant ID from a logical network.
- `vni` (Number) Logical switch VNI.

### Read-Only

- `id` (String) Logical switch ID.

<a id="nestedblock--additional"></a>
### Nested Schema for `additional`

Optional:

- `producer` (String) This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used. Defaults to `default`.


<a id="nestedblock--storm_suppress"></a>
### Nested Schema for `storm_suppress`

Optional:

- `broadcast_cbs` (Number) CBS of broadcast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.
- `broadcast_cbs_unit` (String) CBS unit of broadcast packets. The value can be bytes, Kbytes, or Mbytes.
- `broadcast_cir` (Number) Broadcast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.
- `broadcast_cir_unit` (String) CIR unit of broadcast packets. The value can be Gbit/s, Mbit/s, or kbit/s.
- `broadcast_enable` (Boolean) Whether to enable the broadcast function Defaults to `false`.
- `multicast_cbs` (Number) CBS of multicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.
- `multicast_cbs_unit` (String) CBS unit of multicast packets. The value can be bytes, Kbytes, or Mbytes.
- `multicast_cir` (Number) Multicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.
- `multicast_cir_unit` (String) CIR unit of multicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.
- `multicast_enable` (Boolean) Whether to enable the multicast function. Defaults to `false`.
- `unicast_cbs` (Number) CBS of unicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.
- `unicast_cbs_unit` (String) CBS unit of unicast packets. The value can be byte, Kbytes, or Mbytes.
- `unicast_cir` (Number) Unicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.
- `unicast_cir_unit` (String) CIR unit of unicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.
- `unicast_enable` (Boolean) Whether to enable the unicast function.  Defaults to `false`.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import agile_logical_switch.myswitch 78750ef6-d054-4181-9143-4640dff220e1
```
