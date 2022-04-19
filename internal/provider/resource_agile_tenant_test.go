package provider

import (
	"fmt"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jinzhu/copier"
	"strings"
	"terraform-provider-agile/tools"
	"testing"
)

func TestAccAgileTenant_Basic(t *testing.T) {
	name := "tf_acc_tests_tenant"
	tenantAttr := models.TenantAttributes{
		Quota: &models.TenantQuota{
			LogicVasNum:    agile.Int32(10),
			LogicRouterNum: agile.Int32(15),
			LogicSwitchNum: agile.Int32(20),
		},
	}

	resourceName := "agile_tenant.this"
	var tenant models.Tenant

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcTenantConfig_Basic(name, &tenantAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenant),
					testAccCheckAgileTenantAttributes(name, &tenant, &tenantAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "producer", "default"),
					resource.TestCheckResourceAttr(resourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_vas_num", fmt.Sprint(*tenantAttr.Quota.LogicVasNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_switch_num", fmt.Sprint(*tenantAttr.Quota.LogicSwitchNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_router_num", fmt.Sprint(*tenantAttr.Quota.LogicRouterNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "false"),
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

func TestAccAgileTenant_Complete(t *testing.T) {
	name := "tf_acc_tests_tenant"

	fabrics := GetFabrics()
	externalGateways := GetExternalGateways()

	tenantAttr := models.TenantAttributes{
		Description:         agile.String("Tenant Created via Terraform Agile Provider Acceptance tests"),
		Producer:            agile.String("terraform"),
		MulticastCapability: agile.Bool(true),
		Quota: &models.TenantQuota{
			LogicVasNum:    agile.Int32(10),
			LogicRouterNum: agile.Int32(15),
			LogicSwitchNum: agile.Int32(20),
		},
		MulticastQuota: &models.TenantMulticastQuota{
			AclNum:     agile.Int32(10),
			AclRuleNum: agile.Int32(15),
		},
		ResPool: &models.TenantResPool{
			FabricIds:          []*string{fabrics[3].Id},
			ExternalGatewayIds: []*string{externalGateways[0].Id, externalGateways[5].Id},
		},
	}

	resourceName := "agile_tenant.this"
	var tenant models.Tenant

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcTenantConfig_Complete(name, &tenantAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenant),
					testAccCheckAgileTenantAttributes(name, &tenant, &tenantAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *tenantAttr.Description),
					resource.TestCheckResourceAttr(resourceName, "producer", *tenantAttr.Producer),
					resource.TestCheckResourceAttr(resourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_vas_num", fmt.Sprint(*tenantAttr.Quota.LogicVasNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_switch_num", fmt.Sprint(*tenantAttr.Quota.LogicSwitchNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_router_num", fmt.Sprint(*tenantAttr.Quota.LogicRouterNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.0.acl_num", fmt.Sprint(*tenantAttr.MulticastQuota.AclNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.0.acl_rule_num", fmt.Sprint(*tenantAttr.MulticastQuota.AclRuleNum)),
					resource.TestCheckResourceAttr(resourceName, "res_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.fabric_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.fabric_ids.0", *tenantAttr.ResPool.FabricIds[0]),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.external_gateway_ids.#", fmt.Sprint(len(tenantAttr.ResPool.ExternalGatewayIds))),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.external_gateway_ids.0", *tenantAttr.ResPool.ExternalGatewayIds[0]),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.external_gateway_ids.1", *tenantAttr.ResPool.ExternalGatewayIds[1]),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "true"),
				),
			},
		},
	})
}

func TestAccAgileTenant_Update(t *testing.T) {
	name := "tf_acc_tests_tenant"

	fabrics := GetFabrics()
	externalGateways := GetExternalGateways()

	tenantAttr := models.TenantAttributes{
		Description:         agile.String("Tenant Created via Terraform Agile Provider Acceptance tests"),
		Producer:            agile.String("terraform"),
		MulticastCapability: agile.Bool(true),
		Quota: &models.TenantQuota{
			LogicVasNum:    agile.Int32(10),
			LogicRouterNum: agile.Int32(15),
			LogicSwitchNum: agile.Int32(20),
		},
		MulticastQuota: &models.TenantMulticastQuota{
			AclNum:     agile.Int32(10),
			AclRuleNum: agile.Int32(15),
		},
		ResPool: &models.TenantResPool{
			FabricIds:          []*string{fabrics[3].Id},
			ExternalGatewayIds: []*string{externalGateways[0].Id, externalGateways[5].Id},
		},
	}

	var tenantAttrUpdate models.TenantAttributes
	copier.Copy(&tenantAttrUpdate, &tenantAttr)
	tenantAttrUpdate.Description = agile.String("Tenant Updated via Terraform Agile Provider Acceptance tests")
	tenantAttrUpdate.Quota.LogicRouterNum = agile.Int32(10)
	tenantAttrUpdate.MulticastQuota.AclRuleNum = agile.Int32(7)
	tenantAttrUpdate.ResPool.FabricIds = []*string{fabrics[3].Id}

	resourceName := "agile_tenant.this"
	var tenant models.Tenant
	var tenantUpdated models.Tenant

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcTenantConfig_Complete(name, &tenantAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenant),
					testAccCheckAgileTenantAttributes(name, &tenant, &tenantAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "true"),
				),
			},
			{
				Config: testAccCheckAcTenantConfig_Complete(name, &tenantAttrUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenantUpdated),
					testAccCheckAgileTenantAttributes(name, &tenantUpdated, &tenantAttrUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *tenantAttrUpdate.Description),
					resource.TestCheckResourceAttr(resourceName, "producer", *tenantAttrUpdate.Producer),
					resource.TestCheckResourceAttr(resourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_vas_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicVasNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_switch_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicSwitchNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_router_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicRouterNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.0.acl_num", fmt.Sprint(*tenantAttrUpdate.MulticastQuota.AclNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.0.acl_rule_num", fmt.Sprint(*tenantAttrUpdate.MulticastQuota.AclRuleNum)),
					resource.TestCheckResourceAttr(resourceName, "res_pool.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.fabric_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.fabric_ids.0", *tenantAttrUpdate.ResPool.FabricIds[0]),
					resource.TestCheckResourceAttr(resourceName, "res_pool.0.external_gateway_ids.#", fmt.Sprint(len(tenantAttr.ResPool.ExternalGatewayIds))),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "true"),
				),
			},
			{
				Config: testAccCheckAcTenantConfig_WithoutMulticast(name, &tenantAttrUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenantUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", *tenantAttrUpdate.Description),
					resource.TestCheckResourceAttr(resourceName, "producer", *tenantAttrUpdate.Producer),
					resource.TestCheckResourceAttr(resourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_vas_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicVasNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_switch_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicSwitchNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_router_num", fmt.Sprint(*tenantAttrUpdate.Quota.LogicRouterNum)),
					resource.TestCheckResourceAttr(resourceName, "multicast_quota.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "multicast_capability", "false"),
				),
			},
		},
	})
}

func TestAccAgileTenant_UpdateForceNew(t *testing.T) {
	name := "tf_acc_tests_tenant"
	tenantAttr := models.TenantAttributes{
		Quota: &models.TenantQuota{
			LogicVasNum:    agile.Int32(10),
			LogicRouterNum: agile.Int32(15),
			LogicSwitchNum: agile.Int32(20),
		},
	}

	resourceName := "agile_tenant.this"
	var tenant models.Tenant
	var tenantUpdated models.Tenant

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckAgileTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAcTenantConfig_Basic(name, &tenantAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenant),
					testAccCheckAgileTenantAttributes(name, &tenant, &tenantAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccCheckAcTenantConfig_Basic("tf_acc_tests_tenant_updated", &tenantAttr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAgileTenantExists(resourceName, &tenantUpdated),
					testAccCheckAgileTenantAttributes("tf_acc_tests_tenant_updated", &tenantUpdated, &tenantAttr),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "tf_acc_tests_tenant_updated"),
					resource.TestCheckResourceAttr(resourceName, "quota.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_vas_num", fmt.Sprint(*tenantAttr.Quota.LogicVasNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_switch_num", fmt.Sprint(*tenantAttr.Quota.LogicSwitchNum)),
					resource.TestCheckResourceAttr(resourceName, "quota.0.logic_router_num", fmt.Sprint(*tenantAttr.Quota.LogicRouterNum)),
				),
			},
		},
	})
}

func testAccCheckAcTenantConfig_Basic(name string, tenant *models.TenantAttributes) string {
	return fmt.Sprintf(`
	resource "agile_tenant" "this" {
		name        = "%s"
		quota {
			logic_vas_num = "%d"
			logic_router_num = "%d"
			logic_switch_num = "%d"
		}
	} 
	`, name, *tenant.Quota.LogicVasNum, *tenant.Quota.LogicRouterNum, *tenant.Quota.LogicSwitchNum)
}

func testAccCheckAcTenantConfig_Complete(name string, tenant *models.TenantAttributes) string {
	externalGateways := tools.CreateSliceOfStrings(tenant.ResPool.ExternalGatewayIds)
	externalGatewaysResult := "\"" + strings.Join(externalGateways, "\",\"") + "\""
	return fmt.Sprintf(`
	resource "agile_tenant" "this" {
		name        = "%s"
		description = "%s"
        producer    = "%s"
		quota {
			logic_vas_num = "%d"
			logic_router_num = "%d"
			logic_switch_num = "%d"
		}
		multicast_quota {
			acl_num = %d
			acl_rule_num = %d
		}
		res_pool {
    		fabric_ids = [ "%s" ]
            external_gateway_ids = [ %s ]
  		}
	} 
	`, name,
		*tenant.Description,
		*tenant.Producer,
		*tenant.Quota.LogicVasNum,
		*tenant.Quota.LogicRouterNum,
		*tenant.Quota.LogicSwitchNum,
		*tenant.MulticastQuota.AclNum,
		*tenant.MulticastQuota.AclRuleNum,
		*tenant.ResPool.FabricIds[0],
		externalGatewaysResult,
	)
}

func testAccCheckAcTenantConfig_WithoutMulticast(name string, tenant *models.TenantAttributes) string {
	return fmt.Sprintf(`
	resource "agile_tenant" "this" {
		name        = "%s"
		description = "%s"
        producer    = "%s"
		quota {
			logic_vas_num = "%d"
			logic_router_num = "%d"
			logic_switch_num = "%d"
		}
		res_pool {
    		fabric_ids = [ "%s" ]
  		}
	} 
	`, name,
		*tenant.Description,
		*tenant.Producer,
		*tenant.Quota.LogicVasNum,
		*tenant.Quota.LogicRouterNum,
		*tenant.Quota.LogicSwitchNum,
		*tenant.ResPool.FabricIds[0],
	)
}

func GetFabrics() []*models.Fabric {
	agileClient := testAccProvider.Meta().(*agile.Client)
	fabrics, err := agileClient.ListFabrics(nil)
	if err != nil {
		panic("error query Fabrics")
	}

	return fabrics
}

func GetExternalGateways() []*models.ExternalGateway {
	agileClient := testAccProvider.Meta().(*agile.Client)
	externalGateways, err := agileClient.ListExternalGateways(nil)
	if err != nil {
		panic("error query External Gateways")
	}
	return externalGateways
}

func testAccCheckAgileTenantExists(name string, tenant *models.Tenant) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("tenant %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no tenant id was set")
		}

		agileClient := testAccProvider.Meta().(*agile.Client)

		tenantFound, err := agileClient.GetTenant(rs.Primary.ID)
		if err != nil {
			return err
		}

		if *tenantFound.Id != rs.Primary.ID {
			return fmt.Errorf("tenant %s not found", rs.Primary.ID)
		}

		*tenant = *tenantFound
		return nil
	}
}

func testAccCheckAgileTenantAttributes(name string, tenant *models.Tenant, attributes *models.TenantAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if name != *tenant.Name {
			return fmt.Errorf("bad tenant Name %s", *tenant.Name)
		}

		if attributes.Description != nil && *tenant.Description != *attributes.Description {
			return fmt.Errorf("bad tenant description %s", *tenant.Description)
		}

		if attributes.Producer != nil && (*tenant.Producer != *attributes.Producer) {
			return fmt.Errorf("bad tenant producer %s", *tenant.Producer)
		}

		if *tenant.Quota.LogicVasNum != *attributes.Quota.LogicVasNum {
			return fmt.Errorf("bad tenant quota logical vas num %d", *tenant.Quota.LogicVasNum)
		}

		if *tenant.Quota.LogicRouterNum != *attributes.Quota.LogicRouterNum {
			return fmt.Errorf("bad tenant quota logical router num %d", *tenant.Quota.LogicRouterNum)
		}

		if *tenant.Quota.LogicSwitchNum != *attributes.Quota.LogicSwitchNum {
			return fmt.Errorf("bad tenant quota logical switch num %d", *tenant.Quota.LogicSwitchNum)
		}

		if attributes.MulticastQuota != nil {
			if *tenant.MulticastQuota.AclNum != *attributes.MulticastQuota.AclNum {
				//fmt.Println(*tenant.MulticastQuota.AclNum)
				//fmt.Println(*attributes.MulticastQuota.AclNum)
				return fmt.Errorf("bad tenant multicast quota acl num %d", *tenant.MulticastQuota.AclNum)
			}

			if *tenant.MulticastQuota.AclRuleNum != *attributes.MulticastQuota.AclRuleNum {
				return fmt.Errorf("bad tenant multicast quota acl rule num %d", *tenant.MulticastQuota.AclRuleNum)
			}
		}

		if attributes.ResPool != nil {
			for i, fabric := range attributes.ResPool.FabricIds {
				if *tenant.ResPool.FabricIds[i] != *fabric {
					return fmt.Errorf("bad tenant res pool fabric id %s", *tenant.ResPool.FabricIds[i])
				}
			}

			for i, externalGateway := range attributes.ResPool.ExternalGatewayIds {
				if *tenant.ResPool.ExternalGatewayIds[i] != *externalGateway {
					return fmt.Errorf("bad tenant res pool external gateway %s", *tenant.ResPool.ExternalGatewayIds[i])
				}
			}

		}
		return nil
	}
}

func testAccCheckAgileTenantDestroy(s *terraform.State) error {
	agileClient := testAccProvider.Meta().(*agile.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "agile_tenant" {
			tenant, err := agileClient.GetTenant(rs.Primary.ID)

			if tenant != nil {
				return fmt.Errorf("tenant %s still exists", *tenant.Name)
			}

			if err == nil {
				return fmt.Errorf("tenant %s still exists", *tenant.Name)
			}

		} else {
			continue
		}
	}

	return nil
}
