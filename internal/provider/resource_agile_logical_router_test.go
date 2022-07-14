package provider

import (
	"fmt"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAgileLogicalRouter_Complete(t *testing.T) {
	name := "tf_acc_tests_logicalRouter"

	logicalRouterAttr := models.LogicalRouterAttributes{
		Description:    agile.String("Logical Router created via Terraform Tests"),
		LogicNetworkId: agile.String("acfd8aaf-c6dc-499d-8020-bebd85b1f0e6"),
		Type:           agile.String("Normal"),
		VrfName:        agile.String("Management_67766776"),
		Vni:            agile.Int32(4226),
		RouterLocations: []*models.LogicalRouterLocations{
			{
				FabricId:   agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a"),
				FabricRole: agile.String("master"),
			},
		},
		//Additional: &models.LogicalRouterAdditional{
		//	Producer: agile.String("Terraform"),
		//},
	}

	resourceName := "agile_logical_router.this"
	var logicalRouter models.LogicalRouter

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalRouterConfig_Complete(name, &logicalRouterAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalRouterExists(resourceName, &logicalRouter),
					testAccCheckAgileLogicalRouterAttributes(name, &logicalRouter, &logicalRouterAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalRouterAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "vrf_name", *logicalRouterAttr.VrfName),
					resource.TestCheckResourceAttr(resourceName, "vni", fmt.Sprint(*logicalRouterAttr.Vni)),
					resource.TestCheckResourceAttr(resourceName, "logic_network_id", *logicalRouterAttr.LogicNetworkId),
					resource.TestCheckResourceAttr(resourceName, "router_locations.#", fmt.Sprint(len(logicalRouterAttr.RouterLocations))),
					resource.TestCheckResourceAttr(resourceName, "router_locations.0.fabric_role", *logicalRouterAttr.RouterLocations[0].FabricRole),
					resource.TestCheckResourceAttr(resourceName, "router_locations.0.fabric_id", *logicalRouterAttr.RouterLocations[0].FabricId),
					//resource.TestCheckResourceAttr(resourceName, "additional.#", "1"),
					//resource.TestCheckResourceAttr(resourceName, "additional.0.producer", *logicalRouterAttr.Additional.Producer),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAgileLogicalRouter_Update(t *testing.T) {
	name := "tf_acc_tests_logicalRouter"

	logicalRouterAttr := models.LogicalRouterAttributes{
		Description:    agile.String("Logical Router created via Terraform Tests"),
		LogicNetworkId: agile.String("acfd8aaf-c6dc-499d-8020-bebd85b1f0e6"),
		Type:           agile.String("Normal"),
		VrfName:        agile.String("Management_67766776"),
		Vni:            agile.Int32(4226),
		RouterLocations: []*models.LogicalRouterLocations{
			{
				FabricId:   agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a"),
				FabricRole: agile.String("master"),
			},
		},
		//Additional: &models.LogicalRouterAdditional{
		//	Producer: agile.String("Terraform"),
		//},
	}

	logicalRouterUpdate := logicalRouterAttr
	logicalRouterUpdate.Description = agile.String("Logical Network Updated via Terraform Agile Provider Acceptance tests")

	resourceName := "agile_logical_router.this"
	var logicalRouter models.LogicalRouter
	var logicalRouterUpdated models.LogicalRouter

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalRouterConfig_Complete(name, &logicalRouterAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalRouterExists(resourceName, &logicalRouter),
					testAccCheckAgileLogicalRouterAttributes(name, &logicalRouter, &logicalRouterAttr),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccCheckAcLogicalRouterConfig_Complete(name, &logicalRouterUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalRouterExists(resourceName, &logicalRouterUpdated),
					testAccCheckAgileLogicalRouterAttributes(name, &logicalRouterUpdated, &logicalRouterUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalRouterUpdate.Description),
				),
			},
		},
	})
}

func testAccCheckAcLogicalRouterConfig_Complete(name string, logicalRouter *models.LogicalRouterAttributes) string {
	return fmt.Sprintf(`
	resource "agile_logical_router" "this" {
	  name        = "%s"
      description = "%s"
	  logic_network_id   = "%s"
	  type  = "%s"
      vrf_name = "%s"
      vni = "%d"
      router_locations {
		fabric_id = "%s"
		fabric_role = "%s"
      }
	}
	`, name, *logicalRouter.Description, *logicalRouter.LogicNetworkId, *logicalRouter.Type, *logicalRouter.VrfName, *logicalRouter.Vni,
		*logicalRouter.RouterLocations[0].FabricId, *logicalRouter.RouterLocations[0].FabricRole)
}

func testAccCheckAgileLogicalRouterExists(name string, logicalRouter *models.LogicalRouter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("logical router %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no logical router id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		logicalRouterFound, err := agileClient.GetLogicalRouter(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *logicalRouterFound.Id != rs.Primary.ID {
			return fmt.Errorf("logical router %s not found", rs.Primary.ID)
		}

		*logicalRouter = *logicalRouterFound
		return nil
	}
}

func testAccCheckAgileLogicalRouterAttributes(name string, logicalRouter *models.LogicalRouter, attributes *models.LogicalRouterAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *logicalRouter.Name {
			return fmt.Errorf("bad logical router Name %s", *logicalRouter.Name)
		}

		if attributes.Description != nil && *logicalRouter.Description != *attributes.Description {
			return fmt.Errorf("bad logical router description %s", *logicalRouter.Description)
		}

		//if attributes.Additional != nil {
		//	if *logicalRouter.Additional.Producer != *attributes.Additional.Producer {
		//		return fmt.Errorf("bad logical router producer %s", *logicalRouter.Additional.Producer)
		//	}
		//}

		if attributes.Vni != nil && *logicalRouter.Vni != *attributes.Vni {
			return fmt.Errorf("bad logical router vni %d", *logicalRouter.Vni)
		}

		if attributes.VrfName != nil && *logicalRouter.VrfName != *attributes.VrfName {
			return fmt.Errorf("bad logical router vrf name %s", *logicalRouter.Type)
		}

		if attributes.Type != nil && *logicalRouter.Type != *attributes.Type {
			return fmt.Errorf("bad logical router type %s", *logicalRouter.Type)
		}

		if attributes.LogicNetworkId != nil && *logicalRouter.LogicNetworkId != *attributes.LogicNetworkId {
			return fmt.Errorf("bad logical router network id %s", *logicalRouter.LogicNetworkId)
		}

		for i, location := range attributes.RouterLocations {

			if location.FabricName != nil && *logicalRouter.RouterLocations[i].FabricName != *location.FabricName {
				return fmt.Errorf("bad logical router router location fabric name %s", *logicalRouter.RouterLocations[i].FabricName)
			}

			if location.FabricRole != nil && *logicalRouter.RouterLocations[i].FabricRole != *location.FabricRole {
				return fmt.Errorf("bad logical router router location fabric role %s", *logicalRouter.RouterLocations[i].FabricRole)
			}

			// TODO Falta Router Locations
		}

		return nil
	}
}

func testAccCheckAgileLogicalRouterDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_logical_router" {
			logicalRouter, err := agileClient.GetLogicalRouter(rs.Primary.ID)

			if logicalRouter != nil {
				return fmt.Errorf("logical router %s still exists", *logicalRouter.Name)
			}

			if err == nil {
				return fmt.Errorf("logical router %s still exists", *logicalRouter.Name)
			}

		} else {
			continue
		}
	}

	return nil
}
