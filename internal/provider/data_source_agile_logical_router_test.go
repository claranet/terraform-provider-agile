package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceAgileLogicalRouter(t *testing.T) {
	dataSourceName := "data.agile_logical_router.this"
	resourceName := "agile_logical_router.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileLogicalRouter,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logic_network_id", resourceName, "logic_network_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttr(resourceName, "router_locations.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "router_locations.0.fabric_role", resourceName, "router_locations.0.fabric_role"),
					resource.TestCheckResourceAttrPair(dataSourceName, "router_locations.0.fabric_id", resourceName, "router_locations.0.fabric_id"),
				),
			},
		},
	})
}

const testAccDataSourceAgileLogicalRouter = `
resource "agile_logical_router" "this" {
	name             = "tf_acc_tests_logicalRouter"
	description      = "Logical Router created via Terraform Tests"
	type             = "Normal"
    logic_network_id = "acfd8aaf-c6dc-499d-8020-bebd85b1f0e6"
    vrf_name         = "tf_acc_tests_logicalRouter_4226"
    vni              = 4226
	router_locations {
		fabric_id   = "f1429224-1860-4bdb-8cc8-98ccc0f5563a"
		fabric_role = "master"
	}
}

data "agile_logical_router" "this" {
 name = agile_logical_router.this.name
}
`
