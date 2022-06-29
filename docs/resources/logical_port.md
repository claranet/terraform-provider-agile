---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "agile_logical_port Resource - terraform-provider-agile"
subcategory: ""
description: |-
  Manages Logical Ports.
---

# agile_logical_port (Resource)

Manages Logical Ports.

## Example Usage

```terraform
resource "agile_logical_port" "example" {
  name            = "example"
  description     = "This Logical Port is created by terraform"
  tenant_id       = "11ade37a-79d0-482f-a7a0-6ad070e1d05d"
  fabric_id       = "f1429224-1860-4bdb-8cc8-98ccc0f5563a"
  logic_switch_id = "6c0a96d3-0789-47e6-9dbc-66ac5ba2e519"
  access_info {
    mode = "Uni"
    type = "Dot1q"
    vlan = 1218
    qinq {
      inner_vid_begin = 10
      inner_vid_end   = 10
      outer_vid_begin = 10
      outer_vid_end   = 10
      rewrite_action  = "PopDouble"
    }
    location {
      device_group_id = "e13784fb-499f-4c30-8f9c-e49e6c98fdbb"
      device_id       = "9e3a5bee-3d95-3bf7-90f5-09bd2177324b"
      port_id         = "589c87dd-7222-3c09-87b7-d09a236af285"
    }
    location {
      device_group_id = "e13784fb-499f-4c30-8f9c-e49e6c98fdbb"
      device_id       = "b4f6d9ed-0f1d-3f7a-82f1-a4a7ea4f84d4"
      port_id         = "4c142b5e-1858-33b2-a03e-71dcc3b37360"
    }
    subinterface_number = 18
  }
  additional {
    producer = "Terraform"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_info` (Block Set, Min: 1, Max: 1) Access info Settings. (see [below for nested schema](#nestedblock--access_info))
- `logic_switch_id` (String) Logical switch to which a logical port belongs.
- `name` (String) Logical port name.

### Optional

- `additional` (Block Set, Max: 1) Additional Settings. (see [below for nested schema](#nestedblock--additional))
- `description` (String) Logical port description.
- `fabric_id` (String) Fabric to which the logical port belongs.
- `tenant_id` (String) Tenant to which a logical port belongs. This parameter is automatically obtained by the controller.

### Read-Only

- `id` (String) Logical port ID.

<a id="nestedblock--access_info"></a>
### Nested Schema for `access_info`

Required:

- `location` (Block Set, Min: 1, Max: 4000) Location Settings. (see [below for nested schema](#nestedblock--access_info--location))
- `mode` (String) Port mode, which can be UNI or NNI.
- `type` (String) Logical port type, which can be DOT1Q, DEFAULT, UNTAG, or QINQ.

Optional:

- `qinq` (Block Set, Max: 1) Qinq Settings. (see [below for nested schema](#nestedblock--access_info--qinq))
- `subinterface_number` (Number) Number of an access sub-interface.
- `vlan` (Number) Access VLAN ID. This parameter is mandatory when the access type is dot1q.

<a id="nestedblock--access_info--location"></a>
### Nested Schema for `access_info.location`

Required:

- `device_id` (String) Specified physical device.
- `port_id` (String) Specified physical port.

Optional:

- `device_group_id` (String) Device group ID of a physical device.

Read-Only:

- `device_ip` (String) Device management IP address.
- `port_name` (String) Port name.


<a id="nestedblock--access_info--qinq"></a>
### Nested Schema for `access_info.qinq`

Required:

- `inner_vid_begin` (Number) Start VLAN ID of the inner VLAN tag for QinQ.
- `outer_vid_begin` (Number) Start VLAN ID of the outer VLAN tag for QinQ.
- `rewrite_action` (String) Rewrite action of QinQ, which can be POPDOUBLE or PASSTHROUGH.

Optional:

- `inner_vid_end` (Number) End VLAN ID of the inner VLAN tag for QinQ.
- `outer_vid_end` (Number) End VLAN ID of the outer VLAN tag for QinQ.



<a id="nestedblock--additional"></a>
### Nested Schema for `additional`

Optional:

- `producer` (String) This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used. Defaults to `default`.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import agile_logical_port.myport 78750ef6-d054-4181-9143-4640dff220e1
```