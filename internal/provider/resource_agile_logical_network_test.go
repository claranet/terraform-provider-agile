package provider

import (
	"fmt"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAgileLogicalNetwork_Complete(t *testing.T) {
	name := "tf_acc_tests_logicalNetwork"

	logicalNetworkAttr := models.LogicalNetworkAttributes{
		Description:         agile.String("Logical Network created via Terraform Tests"),
		TenantId:            agile.String("7e0ba3e8-280d-420c-951a-b2fe79b4b68a"),
		FabricId:            []*string{agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a")},
		MulticastCapability: agile.Bool(false),
		Type:                agile.String("Instance"),
		Additional: &models.LogicalNetworkAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	resourceName := "agile_logical_network.this"
	var logicalNetwork models.LogicalNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalNetworkConfig_Complete(name, &logicalNetworkAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalNetworkExists(resourceName, &logicalNetwork),
					testAccCheckAgileLogicalNetworkAttributes(name, &logicalNetwork, &logicalNetworkAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalNetworkAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", *logicalNetworkAttr.TenantId),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", fmt.Sprint(*logicalNetworkAttr.MulticastCapability)),
					resource.TestCheckResourceAttr(resourceName, "fabrics_id.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "fabrics_id.0", *logicalNetworkAttr.FabricId[0]),
					resource.TestCheckResourceAttr(resourceName, "additional.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "additional.0.producer", *logicalNetworkAttr.Additional.Producer),
					resource.TestCheckResourceAttr(resourceName, "type", *logicalNetworkAttr.Type),
					resource.TestCheckResourceAttr(resourceName, "is_vpc_deployed", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"additional"},
			},
		},
	})
}

func TestAccAgileLogicalNetwork_Update(t *testing.T) {
	name := "tf_acc_tests_logicalNetwork"

	logicalNetworkAttr := models.LogicalNetworkAttributes{
		Description:         agile.String("Logical Network created via Terraform Tests"),
		TenantId:            agile.String("7e0ba3e8-280d-420c-951a-b2fe79b4b68a"),
		FabricId:            []*string{agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a")},
		MulticastCapability: agile.Bool(false),
		Type:                agile.String("Instance"),
		Additional: &models.LogicalNetworkAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	logicalNetworkUpdate := logicalNetworkAttr
	logicalNetworkUpdate.Description = agile.String("Logical Network Updated via Terraform Agile Provider Acceptance tests")
	logicalNetworkUpdate.Type = agile.String("Transit")

	resourceName := "agile_logical_network.this"
	var logicalNetwork models.LogicalNetwork
	var logicalNetworkUpdated models.LogicalNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalNetworkConfig_Complete(name, &logicalNetworkAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalNetworkExists(resourceName, &logicalNetwork),
					testAccCheckAgileLogicalNetworkAttributes(name, &logicalNetwork, &logicalNetworkAttr),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_vpc_deployed", "true"),
				),
			},
			{
				Config: testAccCheckAcLogicalNetworkConfig_Complete(name, &logicalNetworkUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalNetworkExists(resourceName, &logicalNetworkUpdated),
					testAccCheckAgileLogicalNetworkAttributes(name, &logicalNetworkUpdated, &logicalNetworkUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalNetworkUpdate.Description),
					resource.TestCheckResourceAttr(resourceName, "type", *logicalNetworkUpdate.Type),
					resource.TestCheckResourceAttr(resourceName, "is_vpc_deployed", "true"),
				),
			},
		},
	})
}

func testAccCheckAcLogicalNetworkConfig_Complete(name string, logicalNetwork *models.LogicalNetworkAttributes) string {
	return fmt.Sprintf(`
	resource "agile_logical_network" "this" {
	  name        = "%s"
      description = "%s"
	  tenant_id   = "%s"
	  fabrics_id  = [ "%s" ]
	  type = "%s"
      multicast_capability = "%t"
	  additional {
		producer = "%s"
	  }
	}
	`, name, *logicalNetwork.Description, *logicalNetwork.TenantId, *logicalNetwork.FabricId[0], *logicalNetwork.Type, *logicalNetwork.MulticastCapability, *logicalNetwork.Additional.Producer)
}

func testAccCheckAgileLogicalNetworkExists(name string, logicalNetwork *models.LogicalNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("logical network %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no logical network id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		logicalNetworkFound, err := agileClient.GetLogicalNetwork(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *logicalNetworkFound.Id != rs.Primary.ID {
			return fmt.Errorf("logical network %s not found", rs.Primary.ID)
		}

		*logicalNetwork = *logicalNetworkFound
		return nil
	}
}

func testAccCheckAgileLogicalNetworkAttributes(name string, logicalNetwork *models.LogicalNetwork, attributes *models.LogicalNetworkAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *logicalNetwork.Name {
			return fmt.Errorf("bad logical network Name %s", *logicalNetwork.Name)
		}

		if attributes.Description != nil && *logicalNetwork.Description != *attributes.Description {
			return fmt.Errorf("bad logical network description %s", *logicalNetwork.Description)
		}

		if attributes.Additional != nil {
			if *logicalNetwork.Additional.Producer != *attributes.Additional.Producer {
				return fmt.Errorf("bad logical network producer %s", *logicalNetwork.Additional.Producer)
			}
		}

		if attributes.MulticastCapability != nil && *logicalNetwork.MulticastCapability != *attributes.MulticastCapability {
			return fmt.Errorf("bad logical network multicast %t", *logicalNetwork.MulticastCapability)
		}

		if attributes.Type != nil && *logicalNetwork.Type != *attributes.Type {
			return fmt.Errorf("bad logical network type %s", *logicalNetwork.Type)
		}

		if attributes.TenantId != nil && *logicalNetwork.TenantId != *attributes.TenantId {
			return fmt.Errorf("bad logical network tenant id %s", *logicalNetwork.TenantId)
		}

		if attributes.FabricId[0] != nil && *logicalNetwork.FabricId[0] != *attributes.FabricId[0] {
			return fmt.Errorf("bad logical network fabric %s", *logicalNetwork.FabricId[0])
		}

		return nil
	}
}

func testAccCheckAgileLogicalNetworkDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_logical_network" {
			logicalNetwork, err := agileClient.GetLogicalNetwork(rs.Primary.ID)

			if logicalNetwork != nil {
				return fmt.Errorf("logical network %s still exists", *logicalNetwork.Name)
			}

			if err == nil {
				return fmt.Errorf("logical network %s still exists", *logicalNetwork.Name)
			}

		} else {
			continue
		}
	}

	return nil
}
