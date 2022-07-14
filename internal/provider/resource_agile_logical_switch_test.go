package provider

import (
	"fmt"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccAgileLogicalSwitch_Complete(t *testing.T) {
	name := "tf_acc_tests_logicalSwitch"

	logicalSwitchAttr := models.LogicalSwitchAttributes{
		Description:    agile.String("Logical Switch created via Terraform Tests"),
		LogicNetworkId: agile.String("5308df55-1709-404f-b4f8-4d8947d8f0c4"),
		Bd:             agile.Int32(4067),
		Vni:            agile.Int32(5096),
		MacAddress:     agile.String("00:00:5E:00:01:02"),
		TenantId:       agile.String("cd27d9cf-9be0-4852-a560-2d6e05fd3c1e"),
		StormSuppress: &models.LogicalSwitchStormSuppress{
			BroadcastEnable:  agile.Bool(true),
			MulticastEnable:  agile.Bool(true),
			UnicastEnable:    agile.Bool(true),
			BroadcastCbs:     agile.String("10000"),
			BroadcastCbsUnit: agile.String("byte"),
			BroadcastCir:     agile.Int64(100),
			BroadcastCirUnit: agile.String("kbps"),
			MulticastCbs:     agile.String("10000"),
			MulticastCbsUnit: agile.String("byte"),
			MulticastCir:     agile.Int64(100),
			MulticastCirUnit: agile.String("kbps"),
			UnicastCbs:       agile.String("10000"),
			UnicastCbsUnit:   agile.String("byte"),
			UnicastCir:       agile.Int64(100),
			UnicastCirUnit:   agile.String("kbps"),
		},
		Additional: &models.LogicalSwitchAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	resourceName := "agile_logical_switch.this"
	var logicalSwitch models.LogicalSwitch

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalSwitchDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalSwitchConfig_Complete(name, &logicalSwitchAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalSwitchExists(resourceName, &logicalSwitch),
					testAccCheckAgileLogicalSwitchAttributes(name, &logicalSwitch, &logicalSwitchAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalSwitchAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", *logicalSwitchAttr.TenantId),
					resource.TestCheckResourceAttr(resourceName, "vni", fmt.Sprint(*logicalSwitchAttr.Vni)),
					resource.TestCheckResourceAttr(resourceName, "logic_network_id", *logicalSwitchAttr.LogicNetworkId),
					resource.TestCheckResourceAttr(resourceName, "bd", fmt.Sprint(*logicalSwitchAttr.Bd)),
					resource.TestCheckResourceAttr(resourceName, "mac_address", *logicalSwitchAttr.MacAddress),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.broadcast_enable", fmt.Sprint(*logicalSwitchAttr.StormSuppress.BroadcastEnable)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.multicast_enable", fmt.Sprint(*logicalSwitchAttr.StormSuppress.MulticastEnable)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.unicast_enable", fmt.Sprint(*logicalSwitchAttr.StormSuppress.UnicastEnable)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.broadcast_cbs", fmt.Sprint(*logicalSwitchAttr.StormSuppress.BroadcastCbs)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.broadcast_cbs_unit", *logicalSwitchAttr.StormSuppress.BroadcastCbsUnit),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.broadcast_cir", fmt.Sprint(*logicalSwitchAttr.StormSuppress.BroadcastCir)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.broadcast_cir_unit", *logicalSwitchAttr.StormSuppress.BroadcastCirUnit),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.unicast_cbs", fmt.Sprint(*logicalSwitchAttr.StormSuppress.UnicastCbs)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.unicast_cbs_unit", *logicalSwitchAttr.StormSuppress.UnicastCbsUnit),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.unicast_cir", fmt.Sprint(*logicalSwitchAttr.StormSuppress.UnicastCir)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.unicast_cir_unit", *logicalSwitchAttr.StormSuppress.UnicastCirUnit),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.multicast_cbs", fmt.Sprint(*logicalSwitchAttr.StormSuppress.MulticastCbs)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.multicast_cbs_unit", *logicalSwitchAttr.StormSuppress.MulticastCbsUnit),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.multicast_cir", fmt.Sprint(*logicalSwitchAttr.StormSuppress.MulticastCir)),
					resource.TestCheckResourceAttr(resourceName, "storm_suppress.0.multicast_cir_unit", *logicalSwitchAttr.StormSuppress.MulticastCirUnit),
					resource.TestCheckResourceAttr(resourceName, "additional.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "additional.0.producer", *logicalSwitchAttr.Additional.Producer),
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

func TestAccAgileLogicalSwitch_Update(t *testing.T) {
	name := "tf_acc_tests_logicalSwitch"

	logicalSwitchAttr := models.LogicalSwitchAttributes{
		Description:    agile.String("Logical Switch created via Terraform Tests"),
		LogicNetworkId: agile.String("5308df55-1709-404f-b4f8-4d8947d8f0c4"),
		Bd:             agile.Int32(4067),
		Vni:            agile.Int32(5096),
		MacAddress:     agile.String("00:00:5E:00:01:02"),
		TenantId:       agile.String("cd27d9cf-9be0-4852-a560-2d6e05fd3c1e"),
		StormSuppress: &models.LogicalSwitchStormSuppress{
			BroadcastEnable:  agile.Bool(true),
			MulticastEnable:  agile.Bool(true),
			UnicastEnable:    agile.Bool(true),
			BroadcastCbs:     agile.String("10000"),
			BroadcastCbsUnit: agile.String("byte"),
			BroadcastCir:     agile.Int64(100),
			BroadcastCirUnit: agile.String("kbps"),
			MulticastCbs:     agile.String("10000"),
			MulticastCbsUnit: agile.String("byte"),
			MulticastCir:     agile.Int64(100),
			MulticastCirUnit: agile.String("kbps"),
			UnicastCbs:       agile.String("10000"),
			UnicastCbsUnit:   agile.String("byte"),
			UnicastCir:       agile.Int64(100),
			UnicastCirUnit:   agile.String("kbps"),
		},
		Additional: &models.LogicalSwitchAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	logicalSwitchUpdate := logicalSwitchAttr
	logicalSwitchUpdate.Description = agile.String("Logical Switch Updated via Terraform Agile Provider Acceptance tests")

	resourceName := "agile_logical_switch.this"
	var logicalSwitch models.LogicalSwitch
	var logicalSwitchUpdated models.LogicalSwitch

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalSwitchConfig_Complete(name, &logicalSwitchAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalSwitchExists(resourceName, &logicalSwitch),
					testAccCheckAgileLogicalSwitchAttributes(name, &logicalSwitch, &logicalSwitchAttr),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccCheckAcLogicalSwitchConfig_Complete(name, &logicalSwitchUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalSwitchExists(resourceName, &logicalSwitchUpdated),
					testAccCheckAgileLogicalSwitchAttributes(name, &logicalSwitchUpdated, &logicalSwitchUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalSwitchUpdate.Description),
				),
			},
		},
	})
}

func testAccCheckAcLogicalSwitchConfig_Complete(name string, logicalSwitch *models.LogicalSwitchAttributes) string {
	return fmt.Sprintf(`
	resource "agile_logical_switch" "this" {
	  name             = "%s"
      description      = "%s"
	  logic_network_id = "%s"
      mac_address      = "%s"
      tenant_id        = "%s"
      bd               = %d
	  vni              = %d
      storm_suppress {
		broadcast_enable = %t
        multicast_enable = %t
        unicast_enable = %t
		broadcast_cbs = %s
        broadcast_cbs_unit = "%s"
		broadcast_cir = %d
 		broadcast_cir_unit = "%s"
        unicast_cbs = %s
        unicast_cbs_unit = "%s"
		unicast_cir = %d
 		unicast_cir_unit = "%s"
		multicast_cbs = %s
        multicast_cbs_unit = "%s"
		multicast_cir = %d
 		multicast_cir_unit = "%s"
      }

	 additional {
		producer = "%s"
	 }

	}
	`, name, *logicalSwitch.Description, *logicalSwitch.LogicNetworkId, *logicalSwitch.MacAddress,
		*logicalSwitch.TenantId, *logicalSwitch.Bd, *logicalSwitch.Vni, *logicalSwitch.StormSuppress.BroadcastEnable, *logicalSwitch.StormSuppress.MulticastEnable, *logicalSwitch.StormSuppress.UnicastEnable,
		*logicalSwitch.StormSuppress.BroadcastCbs, *logicalSwitch.StormSuppress.BroadcastCbsUnit, *logicalSwitch.StormSuppress.BroadcastCir, *logicalSwitch.StormSuppress.BroadcastCirUnit,
		*logicalSwitch.StormSuppress.UnicastCbs, *logicalSwitch.StormSuppress.UnicastCbsUnit, *logicalSwitch.StormSuppress.UnicastCir, *logicalSwitch.StormSuppress.UnicastCirUnit,
		*logicalSwitch.StormSuppress.MulticastCbs, *logicalSwitch.StormSuppress.MulticastCbsUnit, *logicalSwitch.StormSuppress.MulticastCir, *logicalSwitch.StormSuppress.MulticastCirUnit,
		*logicalSwitch.Additional.Producer,
	)
}

func testAccCheckAgileLogicalSwitchExists(name string, logicalSwitch *models.LogicalSwitch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("logical switch %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no logical switch id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		logicalSwitchFound, err := agileClient.GetLogicalSwitch(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *logicalSwitchFound.Id != rs.Primary.ID {
			return fmt.Errorf("logical switch %s not found", rs.Primary.ID)
		}

		*logicalSwitch = *logicalSwitchFound
		return nil
	}
}

func testAccCheckAgileLogicalSwitchAttributes(name string, logicalSwitch *models.LogicalSwitch, attributes *models.LogicalSwitchAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *logicalSwitch.Name {
			return fmt.Errorf("bad logical switch Name %s", *logicalSwitch.Name)
		}

		if attributes.Description != nil && *logicalSwitch.Description != *attributes.Description {
			return fmt.Errorf("bad logical switch description %s", *logicalSwitch.Description)
		}

		if attributes.Additional != nil {
			if *logicalSwitch.Additional.Producer != *attributes.Additional.Producer {
				return fmt.Errorf("bad logical switch producer %s", *logicalSwitch.Additional.Producer)
			}
		}

		if attributes.Vni != nil && *logicalSwitch.Vni != *attributes.Vni {
			return fmt.Errorf("bad logical switch vni %d", *logicalSwitch.Vni)
		}

		if attributes.TenantId != nil && *logicalSwitch.TenantId != *attributes.TenantId {
			return fmt.Errorf("bad logical switch tenant id %s", *logicalSwitch.TenantId)
		}

		if attributes.Bd != nil && *logicalSwitch.Bd != *attributes.Bd {
			return fmt.Errorf("bad logical switch type %d", *logicalSwitch.Bd)
		}

		if attributes.MacAddress != nil && *logicalSwitch.MacAddress != *attributes.MacAddress {
			return fmt.Errorf("bad logical switch mac address %s", *logicalSwitch.MacAddress)
		}

		if attributes.LogicNetworkId != nil && *logicalSwitch.LogicNetworkId != *attributes.LogicNetworkId {
			return fmt.Errorf("bad logical switch network id %s", *logicalSwitch.LogicNetworkId)
		}

		if attributes.StormSuppress != nil {
			if *logicalSwitch.StormSuppress.BroadcastEnable != *attributes.StormSuppress.BroadcastEnable {
				return fmt.Errorf("bad logical switch storm supress broadcast enable %t", *logicalSwitch.StormSuppress.BroadcastEnable)
			}

			if *logicalSwitch.StormSuppress.BroadcastCbs != *attributes.StormSuppress.BroadcastCbs {
				return fmt.Errorf("bad logical switch storm supress broadcast cbs %s", *logicalSwitch.StormSuppress.BroadcastCbs)
			}

			if *logicalSwitch.StormSuppress.BroadcastCbsUnit != *attributes.StormSuppress.BroadcastCbsUnit {
				return fmt.Errorf("bad logical switch storm supress broadcast cbs unit %s", *logicalSwitch.StormSuppress.BroadcastCbsUnit)
			}

			if *logicalSwitch.StormSuppress.BroadcastCir != *attributes.StormSuppress.BroadcastCir {
				return fmt.Errorf("bad logical switch storm supress broadcast cir %d", *logicalSwitch.StormSuppress.BroadcastCir)
			}

			if *logicalSwitch.StormSuppress.BroadcastCirUnit != *attributes.StormSuppress.BroadcastCirUnit {
				return fmt.Errorf("bad logical switch storm supress broadcast cir unit %s", *logicalSwitch.StormSuppress.BroadcastCirUnit)
			}

			if *logicalSwitch.StormSuppress.UnicastEnable != *attributes.StormSuppress.UnicastEnable {
				return fmt.Errorf("bad logical switch storm supress unicast enable %t", *logicalSwitch.StormSuppress.UnicastEnable)
			}

			if *logicalSwitch.StormSuppress.UnicastCbs != *attributes.StormSuppress.UnicastCbs {
				return fmt.Errorf("bad logical switch storm supress unicast cbs %s", *logicalSwitch.StormSuppress.UnicastCbs)
			}

			if *logicalSwitch.StormSuppress.UnicastCbsUnit != *attributes.StormSuppress.UnicastCbsUnit {
				return fmt.Errorf("bad logical switch storm supress unicast cbs unit %s", *logicalSwitch.StormSuppress.UnicastCbsUnit)
			}

			if *logicalSwitch.StormSuppress.UnicastCir != *attributes.StormSuppress.UnicastCir {
				return fmt.Errorf("bad logical switch storm supress unicast cir %d", *logicalSwitch.StormSuppress.UnicastCir)
			}

			if *logicalSwitch.StormSuppress.UnicastCirUnit != *attributes.StormSuppress.UnicastCirUnit {
				return fmt.Errorf("bad logical switch storm supress unicast cir unit %s", *logicalSwitch.StormSuppress.UnicastCirUnit)
			}

			if *logicalSwitch.StormSuppress.MulticastEnable != *attributes.StormSuppress.MulticastEnable {
				return fmt.Errorf("bad logical switch storm supress unicast enable %t", *logicalSwitch.StormSuppress.MulticastEnable)
			}

			if *logicalSwitch.StormSuppress.MulticastCbs != *attributes.StormSuppress.MulticastCbs {
				return fmt.Errorf("bad logical switch storm supress multicast cbs %s", *logicalSwitch.StormSuppress.MulticastCbs)
			}

			if *logicalSwitch.StormSuppress.MulticastCbsUnit != *attributes.StormSuppress.MulticastCbsUnit {
				return fmt.Errorf("bad logical switch storm supress multicast cbs unit %s", *logicalSwitch.StormSuppress.MulticastCbsUnit)
			}

			if *logicalSwitch.StormSuppress.MulticastCir != *attributes.StormSuppress.MulticastCir {
				return fmt.Errorf("bad logical switch storm supress multicast cir %d", *logicalSwitch.StormSuppress.MulticastCir)
			}

			if *logicalSwitch.StormSuppress.MulticastCirUnit != *attributes.StormSuppress.MulticastCirUnit {
				return fmt.Errorf("bad logical switch storm supress multicast cir unit %s", *logicalSwitch.StormSuppress.MulticastCirUnit)
			}
		}

		return nil
	}
}

func testAccCheckAgileLogicalSwitchDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_logical_switch" {
			logicalSwitch, err := agileClient.GetLogicalSwitch(rs.Primary.ID)

			if logicalSwitch != nil {
				return fmt.Errorf("logical switch %s still exists", *logicalSwitch.Name)
			}

			if err == nil {
				return fmt.Errorf("logical switch %s still exists", *logicalSwitch.Name)
			}

		} else {
			continue
		}
	}

	return nil
}
