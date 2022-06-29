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

func TestAccAgileLogicalPort_Complete(t *testing.T) {
	name := "tf_acc_tests_logicalPort"
	logicalPortAttr := models.LogicalPortAttributes{
		Description:   agile.String("Logical Port created via Terraform Tests"),
		TenantId:      agile.String("11ade37a-79d0-482f-a7a0-6ad070e1d05d"),
		FabricId:      agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a"),
		LogicSwitchId: agile.String("6c0a96d3-0789-47e6-9dbc-66ac5ba2e519"),
		AccessInfo: &models.LogicalPortAccessInfo{
			Mode: agile.String("Uni"),
			Type: agile.String("Dot1q"),
			Vlan: agile.Int32(1218),
			Qinq: &models.LogicalPortAccessInfoQinq{
				InnerVidBegin: agile.Int32(10),
				InnerVidEnd:   agile.Int32(10),
				OuterVidBegin: agile.Int32(10),
				OuterVidEnd:   agile.Int32(10),
				RewriteAction: agile.String("PopDouble"),
			},
			Location: []*models.LogicalPortAccessInfoLocation{
				{
					DeviceGroupId: agile.String("e13784fb-499f-4c30-8f9c-e49e6c98fdbb"),
					DeviceId:      agile.String("9e3a5bee-3d95-3bf7-90f5-09bd2177324b"),
					PortId:        agile.String("589c87dd-7222-3c09-87b7-d09a236af285"),
				},
				{
					DeviceGroupId: agile.String("e13784fb-499f-4c30-8f9c-e49e6c98fdbb"),
					DeviceId:      agile.String("b4f6d9ed-0f1d-3f7a-82f1-a4a7ea4f84d4"),
					PortId:        agile.String("4c142b5e-1858-33b2-a03e-71dcc3b37360"),
				},
			},
			SubinterfaceNumber: agile.Int32(16),
		},
		Additional: &models.LogicalPortAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	resourceName := "agile_logical_port.this"
	var logicalPort models.LogicalPort

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalPortConfigComplete(name, &logicalPortAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalPortExists(resourceName, &logicalPort),
					testAccCheckAgileLogicalPortAttributes(name, &logicalPort, &logicalPortAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalPortAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", *logicalPortAttr.TenantId),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", *logicalPortAttr.FabricId),
					resource.TestCheckResourceAttr(resourceName, "logic_switch_id", *logicalPortAttr.LogicSwitchId),
					resource.TestCheckResourceAttr(resourceName, "access_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.mode", *logicalPortAttr.AccessInfo.Mode),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.type", *logicalPortAttr.AccessInfo.Type),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.vlan", fmt.Sprint(*logicalPortAttr.AccessInfo.Vlan)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.inner_vid_begin", fmt.Sprint(*logicalPortAttr.AccessInfo.Qinq.InnerVidBegin)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.inner_vid_end", fmt.Sprint(*logicalPortAttr.AccessInfo.Qinq.InnerVidEnd)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.outer_vid_begin", fmt.Sprint(*logicalPortAttr.AccessInfo.Qinq.OuterVidBegin)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.outer_vid_end", fmt.Sprint(*logicalPortAttr.AccessInfo.Qinq.OuterVidEnd)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.rewrite_action", *logicalPortAttr.AccessInfo.Qinq.RewriteAction),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.device_group_id", *logicalPortAttr.AccessInfo.Location[0].DeviceGroupId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.device_id", *logicalPortAttr.AccessInfo.Location[0].DeviceId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.port_id", *logicalPortAttr.AccessInfo.Location[0].PortId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.device_group_id", *logicalPortAttr.AccessInfo.Location[1].DeviceGroupId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.device_id", *logicalPortAttr.AccessInfo.Location[1].DeviceId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.port_id", *logicalPortAttr.AccessInfo.Location[1].PortId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.subinterface_number", fmt.Sprint(*logicalPortAttr.AccessInfo.SubinterfaceNumber)),
					resource.TestCheckResourceAttr(resourceName, "additional.0.producer", *logicalPortAttr.Additional.Producer),
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

func TestAccAgileLogicalPort_Update(t *testing.T) {
	name := "tf_acc_tests_logicalPort"

	logicalPortAttr := models.LogicalPortAttributes{
		Description:   agile.String("Logical Port created via Terraform Tests"),
		TenantId:      agile.String("11ade37a-79d0-482f-a7a0-6ad070e1d05d"),
		FabricId:      agile.String("f1429224-1860-4bdb-8cc8-98ccc0f5563a"),
		LogicSwitchId: agile.String("6c0a96d3-0789-47e6-9dbc-66ac5ba2e519"),
		AccessInfo: &models.LogicalPortAccessInfo{
			Mode: agile.String("Uni"),
			Type: agile.String("Dot1q"),
			Vlan: agile.Int32(1218),
			Qinq: &models.LogicalPortAccessInfoQinq{
				InnerVidBegin: agile.Int32(10),
				InnerVidEnd:   agile.Int32(10),
				OuterVidBegin: agile.Int32(10),
				OuterVidEnd:   agile.Int32(10),
				RewriteAction: agile.String("PopDouble"),
			},
			Location: []*models.LogicalPortAccessInfoLocation{
				{
					DeviceGroupId: agile.String("e13784fb-499f-4c30-8f9c-e49e6c98fdbb"),
					DeviceId:      agile.String("9e3a5bee-3d95-3bf7-90f5-09bd2177324b"),
					PortId:        agile.String("589c87dd-7222-3c09-87b7-d09a236af285"),
				},
				{
					DeviceGroupId: agile.String("e13784fb-499f-4c30-8f9c-e49e6c98fdbb"),
					DeviceId:      agile.String("b4f6d9ed-0f1d-3f7a-82f1-a4a7ea4f84d4"),
					PortId:        agile.String("4c142b5e-1858-33b2-a03e-71dcc3b37360"),
				},
			},
			SubinterfaceNumber: agile.Int32(16),
		},
		Additional: &models.LogicalPortAdditional{
			Producer: agile.String("Terraform"),
		},
	}

	var logicalPortAttrUpdate models.LogicalPortAttributes
	copier.Copy(&logicalPortAttrUpdate, &logicalPortAttr)
	logicalPortAttrUpdate.Description = agile.String("Logical Port Updated via Terraform Agile Provider Acceptance tests")

	resourceName := "agile_logical_port.this"
	var logicalPort models.LogicalPort
	var logicalPortUpdated models.LogicalPort

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileLogicalPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcLogicalPortConfigComplete(name, &logicalPortAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalPortExists(resourceName, &logicalPort),
					testAccCheckAgileLogicalPortAttributes(name, &logicalPort, &logicalPortAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalPortAttr.Description),
				),
			},
			{
				Config: testAccCheckAcLogicalPortConfigComplete(name, &logicalPortAttrUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileLogicalPortExists(resourceName, &logicalPortUpdated),
					testAccCheckAgileLogicalPortAttributes(name, &logicalPortUpdated, &logicalPortAttrUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *logicalPortAttrUpdate.Description),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", *logicalPortAttrUpdate.TenantId),
					resource.TestCheckResourceAttr(resourceName, "fabric_id", *logicalPortAttrUpdate.FabricId),
					resource.TestCheckResourceAttr(resourceName, "logic_switch_id", *logicalPortAttrUpdate.LogicSwitchId),
					resource.TestCheckResourceAttr(resourceName, "access_info.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.mode", *logicalPortAttrUpdate.AccessInfo.Mode),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.type", *logicalPortAttrUpdate.AccessInfo.Type),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.vlan", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.Vlan)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.inner_vid_begin", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.Qinq.InnerVidBegin)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.inner_vid_end", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.Qinq.InnerVidEnd)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.outer_vid_begin", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.Qinq.OuterVidBegin)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.outer_vid_end", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.Qinq.OuterVidEnd)),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.qinq.0.rewrite_action", *logicalPortAttrUpdate.AccessInfo.Qinq.RewriteAction),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.device_group_id", *logicalPortAttrUpdate.AccessInfo.Location[0].DeviceGroupId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.device_id", *logicalPortAttrUpdate.AccessInfo.Location[0].DeviceId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.0.port_id", *logicalPortAttrUpdate.AccessInfo.Location[0].PortId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.device_group_id", *logicalPortAttrUpdate.AccessInfo.Location[1].DeviceGroupId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.device_id", *logicalPortAttrUpdate.AccessInfo.Location[1].DeviceId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.location.1.port_id", *logicalPortAttrUpdate.AccessInfo.Location[1].PortId),
					resource.TestCheckResourceAttr(resourceName, "access_info.0.subinterface_number", fmt.Sprint(*logicalPortAttrUpdate.AccessInfo.SubinterfaceNumber)),
					resource.TestCheckResourceAttr(resourceName, "additional.0.producer", *logicalPortAttrUpdate.Additional.Producer),
				),
			},
		},
	})
}

func testAccCheckAgileLogicalPortAttributes(name string, logicalPort *models.LogicalPort, attributes *models.LogicalPortAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *logicalPort.Name {
			return fmt.Errorf("bad logical port Name %s", *logicalPort.Name)
		}

		if attributes.Description != nil && *logicalPort.Description != *attributes.Description {
			return fmt.Errorf("bad logical port description %s", *logicalPort.Description)
		}

		if attributes.TenantId != nil && *logicalPort.TenantId != *attributes.TenantId {
			return fmt.Errorf("bad logical port Tenant ID %s", *logicalPort.TenantId)
		}

		if attributes.FabricId != nil && *logicalPort.FabricId != *attributes.FabricId {
			return fmt.Errorf("bad logical port Logic Fabric ID %s", *logicalPort.FabricId)
		}

		if attributes.LogicSwitchId != nil && *logicalPort.LogicSwitchId != *attributes.LogicSwitchId {
			return fmt.Errorf("bad logical port Logic Switch ID %s", *logicalPort.LogicSwitchId)
		}

		if attributes.Additional != nil {
			if *logicalPort.Additional.Producer != *attributes.Additional.Producer {
				return fmt.Errorf("bad logical port producer %s", *logicalPort.Additional.Producer)
			}
		}

		if attributes.AccessInfo.Mode != nil && *logicalPort.AccessInfo.Mode != *attributes.AccessInfo.Mode {
			return fmt.Errorf("bad logical port Access Info Mode %s", *logicalPort.AccessInfo.Mode)
		}

		if attributes.AccessInfo.Type != nil && *logicalPort.AccessInfo.Type != *attributes.AccessInfo.Type {
			return fmt.Errorf("bad logical port Access Info type %s", *logicalPort.AccessInfo.Type)
		}

		if attributes.AccessInfo.SubinterfaceNumber != nil && *logicalPort.AccessInfo.SubinterfaceNumber != *attributes.AccessInfo.SubinterfaceNumber {
			return fmt.Errorf("bad logical port Access Info Subinterface Number %d", *logicalPort.AccessInfo.SubinterfaceNumber)
		}

		if attributes.AccessInfo.Qinq.InnerVidBegin != nil && *logicalPort.AccessInfo.Qinq.InnerVidBegin != *attributes.AccessInfo.Qinq.InnerVidBegin {
			return fmt.Errorf("bad logical port Access Info Qinq Inner Vid begin %d", *logicalPort.AccessInfo.Qinq.InnerVidBegin)
		}

		if attributes.AccessInfo.Qinq.InnerVidEnd != nil && *logicalPort.AccessInfo.Qinq.InnerVidEnd != *attributes.AccessInfo.Qinq.InnerVidEnd {
			return fmt.Errorf("bad logical port Access Info Qinq Inner Vid end %d", *logicalPort.AccessInfo.Qinq.InnerVidEnd)
		}

		if attributes.AccessInfo.Qinq.OuterVidBegin != nil && *logicalPort.AccessInfo.Qinq.OuterVidBegin != *attributes.AccessInfo.Qinq.OuterVidBegin {
			return fmt.Errorf("bad logical port Access Info Qinq Outer Vid begin %d", *logicalPort.AccessInfo.Qinq.InnerVidBegin)
		}

		if attributes.AccessInfo.Qinq.OuterVidEnd != nil && *logicalPort.AccessInfo.Qinq.OuterVidEnd != *attributes.AccessInfo.Qinq.OuterVidEnd {
			return fmt.Errorf("bad logical port Access Info Qinq Outer Vid end %d", *logicalPort.AccessInfo.Qinq.OuterVidEnd)
		}

		if attributes.AccessInfo.Qinq.RewriteAction != nil && *logicalPort.AccessInfo.Qinq.RewriteAction != *attributes.AccessInfo.Qinq.RewriteAction {
			return fmt.Errorf("bad logical port Access Info Qinq Rewrite Action %s", *logicalPort.AccessInfo.Qinq.RewriteAction)
		}

		for i, location := range attributes.AccessInfo.Location {

			if location.DeviceGroupId != nil && *logicalPort.AccessInfo.Location[i].DeviceGroupId != *location.DeviceGroupId {
				return fmt.Errorf("bad logical port access info device group id %s", *logicalPort.AccessInfo.Location[0].DeviceGroupId)
			}

			if location.PortId != nil && *logicalPort.AccessInfo.Location[i].PortId != *location.PortId {
				return fmt.Errorf("bad logical port access info port id %s", *logicalPort.AccessInfo.Location[0].PortId)
			}

			if location.DeviceId != nil && *logicalPort.AccessInfo.Location[i].DeviceId != *location.DeviceId {
				return fmt.Errorf("bad logical port access info device id %s", *logicalPort.AccessInfo.Location[0].DeviceId)
			}
		}
		return nil
	}
}

func testAccCheckAgileLogicalPortExists(name string, logicalPort *models.LogicalPort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("logical port %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no logical port id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		logicalPortFound, err := agileClient.GetLogicalPort(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *logicalPortFound.Id != rs.Primary.ID {
			return fmt.Errorf("logical port %s not found", rs.Primary.ID)
		}

		*logicalPort = *logicalPortFound
		return nil
	}
}

func testAccCheckAgileLogicalPortDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_logical_port" {
			logicalPort, err := agileClient.GetLogicalPort(rs.Primary.ID)

			if logicalPort != nil {
				return fmt.Errorf("logical port %s still exists", *logicalPort.Name)
			}

			if err == nil {
				return fmt.Errorf("logical port %s still exists", *logicalPort.Name)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAcLogicalPortConfigComplete(name string, logicalPort *models.LogicalPortAttributes) string {
	return fmt.Sprintf(`
	resource "agile_logical_port" "this" {
		name        = "%s"
        description = "%s"
        tenant_id = "%s"
        fabric_id = "%s"
        logic_switch_id = "%s"
		access_info {
			mode = "%s"
			type = "%s"
			vlan = "%d"
			qinq {
			  inner_vid_begin = "%d"
			  inner_vid_end = "%d"
			  outer_vid_begin = "%d"
			  outer_vid_end = "%d"
			  rewrite_action = "%s"
			}
			location {
				device_group_id = "%s"
				device_id = "%s"
				port_id = "%s"
			}
			
			location {
				device_group_id = "%s"
				device_id = "%s"
				port_id = "%s"
			}

			subinterface_number = "%d"
		}
		additional {
			producer = "%s"
		}
	} 
	`, name, *logicalPort.Description, *logicalPort.TenantId, *logicalPort.FabricId, *logicalPort.LogicSwitchId,
		*logicalPort.AccessInfo.Mode, *logicalPort.AccessInfo.Type, *logicalPort.AccessInfo.Vlan, *logicalPort.AccessInfo.Qinq.InnerVidBegin,
		*logicalPort.AccessInfo.Qinq.InnerVidEnd, *logicalPort.AccessInfo.Qinq.OuterVidBegin, *logicalPort.AccessInfo.Qinq.OuterVidEnd,
		*logicalPort.AccessInfo.Qinq.RewriteAction, *logicalPort.AccessInfo.Location[0].DeviceGroupId, *logicalPort.AccessInfo.Location[0].DeviceId,
		*logicalPort.AccessInfo.Location[0].PortId, *logicalPort.AccessInfo.Location[1].DeviceGroupId, *logicalPort.AccessInfo.Location[1].DeviceId,
		*logicalPort.AccessInfo.Location[1].PortId, *logicalPort.AccessInfo.SubinterfaceNumber, *logicalPort.Additional.Producer,
	)
}
