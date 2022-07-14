package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceAgileLogicalSwitch(t *testing.T) {
	dataSourceName := "data.agile_logical_switch.this"
	resourceName := "agile_logical_switch.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileLogicalSwitch,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logic_network_id", resourceName, "logic_network_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bd", resourceName, "bd"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vni", resourceName, "vni"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_id", resourceName, "tenant_id"),
				),
			},
		},
	})
}

const testAccDataSourceAgileLogicalSwitch = `
resource "agile_logical_switch" "this" {
	name             = "tf_acc_tests_logical_switch"
	description      = "Logical Switch created via Terraform Tests"
    tenant_id         = "cd27d9cf-9be0-4852-a560-2d6e05fd3c1e"
	logic_network_id = "5308df55-1709-404f-b4f8-4d8947d8f0c4" 
	additional {
		producer = " terraform"
	}
}

data "agile_logical_switch" "this" {
 name = agile_logical_switch.this.name
}
`
