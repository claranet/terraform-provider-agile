resource "agile_logical_router" "example" {
  name             = "example"
  description      = "This Logical Router is created by terraform"
  logic_network_id = "acfd8aaf-c6dc-499d-8020-bebd85b1f0e6"
  type             = "Normal"
  vrf_name         = "Management_67766776"
  vni              = 4226
  router_locations {
    fabric_id   = "f1429224-1860-4bdb-8cc8-98ccc0f5563a"
    fabric_role = "master"
  }
}