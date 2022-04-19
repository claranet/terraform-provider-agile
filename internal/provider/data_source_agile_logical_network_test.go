package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceAgileLogicalNetwork(t *testing.T) {
	dataSourceName := "data.agile_logical_network.this"
	resourceName := "agile_logical_network.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileLogicalNetwork,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "multicast_capability", resourceName, "multicast_capability"),
				),
			},
		},
	})
}

const testAccDataSourceAgileLogicalNetwork = `
resource "agile_logical_network" "this" {
	name        = "tf_acc_tests_logicalNetwork"
	description = "Logical Network created via Terraform Tests"
	type = "Transit"
	multicast_capability = "false"
	tenant_id = "7e0ba3e8-280d-420c-951a-b2fe79b4b68a"
	fabrics_id = [ "f1429224-1860-4bdb-8cc8-98ccc0f5563a" ]
}

data "agile_logical_network" "this" {
 name = agile_logical_network.this.name
}
`
