resource "agile_tenant" "example" {
  name        = "example"
  description = "This tenant is created by terraform"
  producer    = "Terraform"
  quota {
    logic_vas_num    = 10
    logic_router_num = 10
    logic_switch_num = 10
  }

  multicast_quota {
    acl_num      = 10
    acl_rule_num = 10
  }

  res_pool {
    fabric_ids           = ["cf121667-d982-4049-8fd2-fc31857b6613"]
    external_gateway_ids = ["924e8a8f-9b64-4a0b-84fc-54c5a3b0efc6"]
    vmm_ids              = ["d71c0b73-b80b-42ed-aa11-d51ac207c7ba"]
    dhcp_group_ids       = ["8726d0c1-5c87-4572-9429-13afddcaeafb"]
  }
}