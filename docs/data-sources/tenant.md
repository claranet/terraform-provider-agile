---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "agile_tenant Data Source - terraform-provider-agile"
subcategory: ""
description: |-
  Data source can be used to retrieve Tenant by name.
---

# agile_tenant (Data Source)

Data source can be used to retrieve Tenant by name.

## Example Usage

```terraform
data "agile_tenant" "example" {
  name = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Tenant name.

### Read-Only

- **description** (String) Tenant description.
- **id** (String) Tenant ID.
- **multicast_capability** (Boolean) Whether the multicast capability is supported.
- **producer** (String) Producer.

