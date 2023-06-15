package provider

import (
	"fmt"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jinzhu/copier"
	"testing"
)

func TestAccAgileEndPort_Complete(t *testing.T) {
	name := "tf_acc_tests_endPort"
	endPortAttr := models.EndPortAttributes{
		Description:    agile.String("End Port created via Terraform Tests"),
		LogicPortId:    agile.String("8682b032-6bb3-46f8-b7b5-2c8bcfceefff"),
		LogicNetworkId: agile.String("5308df55-1709-404f-b4f8-4d8947d8f0c4"),
		Location:       agile.String("10"),
		VmName:         agile.String("10"),
		Ipv4:           []*string{agile.String("192.168.1.1")},
		Ipv6:           []*string{agile.String("FE80::A1")},
	}

	resourceName := "agile_end_port.this"
	var endPort models.EndPort

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileEndPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcEndPortConfigComplete(name, &endPortAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileEndPortExists(resourceName, &endPort),
					testAccCheckAgileEndPortAttributes(name, &endPort, &endPortAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *endPortAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "logic_port_id", *endPortAttr.LogicPortId),
					resource.TestCheckResourceAttr(resourceName, "logic_network_id", *endPortAttr.LogicNetworkId),
					resource.TestCheckResourceAttr(resourceName, "location", *endPortAttr.Location),
					resource.TestCheckResourceAttr(resourceName, "vm_name", *endPortAttr.VmName),
					resource.TestCheckResourceAttr(resourceName, "ipv4", *endPortAttr.Ipv4[0]),
					resource.TestCheckResourceAttr(resourceName, "ipv6", *endPortAttr.Ipv6[0]),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAgileEndPort_Update(t *testing.T) {
	name := "tf_acc_tests_endPort"
	endPortAttr := models.EndPortAttributes{
		Description:    agile.String("End Port created via Terraform Tests"),
		LogicPortId:    agile.String("8682b032-6bb3-46f8-b7b5-2c8bcfceefff"),
		LogicNetworkId: agile.String("5308df55-1709-404f-b4f8-4d8947d8f0c4"),
		Location:       agile.String("10"),
		VmName:         agile.String("10"),
		Ipv4:           []*string{agile.String("192.168.1.1")},
		Ipv6:           []*string{agile.String("FE80::A1")},
	}

	var endPortAttrUpdate models.EndPortAttributes
	copier.Copy(&endPortAttrUpdate, &endPortAttr)
	endPortAttrUpdate.Description = agile.String("End Port Updated via Terraform Agile Provider Acceptance tests")
	endPortAttr.Location = agile.String("11")
	endPortAttr.VmName = agile.String("11")

	resourceName := "agile_end_port.this"
	var endPort models.EndPort
	var endPortUpdated models.EndPort

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileEndPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcEndPortConfigComplete(name, &endPortAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileEndPortExists(resourceName, &endPort),
					testAccCheckAgileEndPortAttributes(name, &endPort, &endPortAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *endPort.Description),
					resource.TestCheckResourceAttr(resourceName, "location", *endPortAttr.Location),
					resource.TestCheckResourceAttr(resourceName, "vm_name", *endPortAttr.VmName),
				),
			},
			{
				Config: testAccCheckAcEndPortConfigComplete(name, &endPortAttrUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileEndPortExists(resourceName, &endPortUpdated),
					testAccCheckAgileEndPortAttributes(name, &endPortUpdated, &endPortAttrUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *endPortAttrUpdate.Description),
					resource.TestCheckResourceAttr(resourceName, "logic_port_id", *endPortAttrUpdate.LogicPortId),
					resource.TestCheckResourceAttr(resourceName, "logic_network_id", *endPortAttrUpdate.LogicNetworkId),
					resource.TestCheckResourceAttr(resourceName, "location", *endPortAttrUpdate.Location),
					resource.TestCheckResourceAttr(resourceName, "vm_name", *endPortAttrUpdate.VmName),
					resource.TestCheckResourceAttr(resourceName, "ipv4", *endPortAttrUpdate.Ipv4[0]),
					resource.TestCheckResourceAttr(resourceName, "ipv6", *endPortAttrUpdate.Ipv6[0]),
				),
			},
		},
	})
}

func testAccCheckAgileEndPortAttributes(name string, endPort *models.EndPort, attributes *models.EndPortAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *endPort.Name {
			return fmt.Errorf("bad end port Name %s", *endPort.Name)
		}

		if attributes.Description != nil && *endPort.Description != *attributes.Description {
			return fmt.Errorf("bad end port description %s", *endPort.Description)
		}

		if attributes.LogicPortId != nil && *endPort.LogicPortId != *attributes.LogicPortId {
			return fmt.Errorf("bad end port logic port ID %s", *endPort.LogicPortId)
		}

		if attributes.LogicNetworkId != nil && *endPort.LogicNetworkId != *attributes.LogicNetworkId {
			return fmt.Errorf("bad end port Logic Network ID %s", *endPort.LogicNetworkId)
		}

		if attributes.Location != nil && *endPort.Location != *attributes.Location {
			return fmt.Errorf("bad end port Logic location %s", *endPort.Location)
		}

		if attributes.VmName != nil && *endPort.VmName != *attributes.VmName {
			return fmt.Errorf("bad end port Logic Vm Name %s", *endPort.VmName)
		}

		if attributes.Ipv4[0] != nil && *endPort.Ipv4[0] != *attributes.Ipv4[0] {
			return fmt.Errorf("bad end port ipv4 %s", *endPort.Ipv4[0])
		}

		if attributes.Ipv6[0] != nil && *endPort.Ipv6[0] != *attributes.Ipv6[0] {
			return fmt.Errorf("bad end port ipv6 %s", *endPort.Ipv6[0])
		}

		return nil
	}
}

func testAccCheckAgileEndPortExists(name string, endort *models.EndPort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("end port %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no end port id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		endPortFound, err := agileClient.GetEndPort(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *endPortFound.Id != rs.Primary.ID {
			return fmt.Errorf("end port %s not found", rs.Primary.ID)
		}

		*endort = *endPortFound
		return nil
	}
}

func testAccCheckAgileEndPortDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_end_port" {
			endPort, err := agileClient.GetEndPort(rs.Primary.ID)

			if err == nil {
				return fmt.Errorf("end port %s still exists", *endPort.Name)
			}

			if endPort != nil {
				return fmt.Errorf("end port %s still exists", *endPort.Name)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAcEndPortConfigComplete(name string, endPort *models.EndPortAttributes) string {
	return fmt.Sprintf(`
	resource "agile_end_port" "this" {
		name             = "%s"
        description      = "%s"
        logic_port_id    = "%s"
        logic_network_id = "%s"
        location         = "%s"
		vm_name          = "%s"
		ipv4             = "%s"
		ipv6             = "%s"
	} 
	`, name, *endPort.Description, *endPort.LogicPortId, *endPort.LogicNetworkId, *endPort.Location,
		*endPort.VmName, *endPort.Ipv4[0], *endPort.Ipv6[0],
	)
}
