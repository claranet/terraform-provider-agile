package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAgileTenant(t *testing.T) {
	dataSourceName := "data.agile_tenant.this"
	resourceName := "agile_tenant.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileTenant,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "producer", resourceName, "producer"),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "false"),
				),
			},
		},
	})
}

const testAccDataSourceAgileTenant = `

resource "agile_tenant" "this" {
	name  = "tf_acc_tests_tenant"
	producer = "terraform_producer"
    description = "terraform acceptance tests"
	quota {
		logic_vas_num = "10"
		logic_router_num = "15"
		logic_switch_num = "18"
	}
}


data "agile_tenant" "this" {
 name = agile_tenant.this.name
}
`
