resource "agile_logical_network" "example" {
  name                 = "tf_test_logical_network"
  description          = "This Logical Network is created by terraform"
  tenant_id            = "65825aae-8804-4356-b29d-76d9daca9ad8"
  fabrics_id           = ["dd61883a-440b-4336-84aa-8e43b9f33b6a"]
  multicast_capability = false
  type                 = "Transit"
  additional {
    producer = "Terraform"
  }
}