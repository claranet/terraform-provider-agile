package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAgileFabric(t *testing.T) {
	dataSourceName := "data.agile_fabric.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileFabric,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "id", "804c7c74-5586-48bf-9cea-96a6d4d3f3a5"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "tesye"),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceName, "network_type", "Distributed"),
					resource.TestCheckResourceAttr(dataSourceName, "physical_network_mode", "Vxlan"),
					resource.TestCheckResourceAttr(dataSourceName, "multicast_capability", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "micro_segment", "false"),
				),
			},
		},
	})
}

const testAccDataSourceAgileFabric = `
data "agile_fabric" "this" {
 name = "tesye"
}
`
