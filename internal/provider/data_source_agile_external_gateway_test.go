package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAgileExternalGateway(t *testing.T) {
	dataSourceName := "data.agile_external_gateway.this"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgileExternalGateway,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "id", "b620662c-4c9f-46d4-9798-6728d4ef7131"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "TESTE"),
					resource.TestCheckResourceAttr(dataSourceName, "description", ""),
					resource.TestCheckResourceAttr(dataSourceName, "gateway_type", "Public"),
					resource.TestCheckResourceAttr(dataSourceName, "vrf_name", "teste"),
				),
			},
		},
	})
}

const testAccDataSourceAgileExternalGateway = `
data "agile_external_gateway" "this" {
 name = "TESTE"
}
`
