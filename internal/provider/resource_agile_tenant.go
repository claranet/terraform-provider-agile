package provider

import (
	"context"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"terraform-provider-agile/tools"
)

func resourceAgileTenant() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages Agile Tenant.",
		CreateContext: resourceAgileTenantCreate,
		ReadContext:   resourceAgileTenantRead,
		UpdateContext: resourceAgileTenantUpdate,
		DeleteContext: resourceAgileTenantDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileTenantImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Tenant ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Tenant name.",
				Required:    true,
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 128),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Tenant description.",
				Optional:    true,
				Default:     "",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"producer": {
				Type:        schema.TypeString,
				Description: "Producer. This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used.",
				Optional:    true,
				Default:     "default",
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(1, 36),
				),
			},
			"multicast_capability": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the multicast capability for the tenant. By default, the multicast capability of a tenant is disabled.",
				Computed:    true,
			},
			"multicast_quota": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_num": {
							Type:        schema.TypeInt,
							Description: "Maximum number of ACL rules that can be created for the tenant when the multicast capability is enabled. The value is an integer in the range from 0 to 1000",
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 1000),
							),
						},
						"acl_rule_num": {
							Type:        schema.TypeInt,
							Description: "Maximum number of ACL rules that can be created for the tenant when the multicast capability is enabled. The value is an integer in the range from 0 to 3000",
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 3000),
							),
						},
					},
				},
			},
			"quota": {
				Type:       schema.TypeSet,
				Required:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logic_vas_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of logical VASs that can be created for the tenant. The value is an integer in the range from 0 to 2000.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 2000),
							),
						},
						"logic_router_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of logical routers that can be created for the tenant. The value is an integer in the range from 0 to 1000.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 1000),
							),
						},
						"logic_switch_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximum number of logical switches that can be created for the tenant. The value is an integer in the range from 0 to 6000.",
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(0, 6000),
							),
						},
					},
				},
			},
			"res_pool": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_gateway_ids": {
							Type:        schema.TypeSet,
							Description: "External gateway that can be used by the tenant. UUID Version 4 Format.",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
						"fabric_ids": {
							Type:        schema.TypeSet,
							Description: "Fabrics network that can be used by the tenant. UUID Version 4 Format.",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
						"vmm_ids": {
							Type:        schema.TypeSet,
							Description: "VMMs that can be used by the tenant. This parameter is required only in network virtualization - computing scenarios. UUID Version 4 Format.",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
						"dhcp_group_ids": {
							Type:        schema.TypeSet,
							Description: "DHCP groups that can be used by the tenant to dynamically allocate IP addresses to VMs in VPC services. UUID Version 4 Format.",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ComputedIf("multicast_capability", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("multicast_quota")
			}),
			// Agile Controller Bug: Fields can't back to null
			customdiff.ForceNewIfChange("res_pool.0.external_gateway_ids", func(ctx context.Context, old, new, meta interface{}) bool {
				if old != nil && (new == nil || len(new.(*schema.Set).List()) == 0) {
					return true
				}
				return false
			}),
			customdiff.ForceNewIfChange("res_pool.0.fabric_ids", func(ctx context.Context, old, new, meta interface{}) bool {
				if old != nil && (new == nil || len(new.(*schema.Set).List()) == 0) {
					return true
				}
				return false
			}),
			customdiff.ForceNewIfChange("res_pool.0.vmm_ids", func(ctx context.Context, old, new, meta interface{}) bool {
				if old != nil && (new == nil || len(new.(*schema.Set).List()) == 0) {
					return true
				}
				return false
			}),
			customdiff.ForceNewIfChange("res_pool.0.dhcp_group_ids", func(ctx context.Context, old, new, meta interface{}) bool {
				if old != nil && (new == nil || len(new.(*schema.Set).List()) == 0) {
					return true
				}
				return false
			}),
		),
	}
}

func resourceAgileTenantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tenant: Beginning Creation")
	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	tenant, errTenant := NewTenantAttributes(d)

	if errTenant != nil {
		return errTenant
	}

	err := agileClient.CreateTenant(agile.String(id.String()), agile.String(name), tenant)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAgileTenantRead(ctx, d, meta)
}

func resourceAgileTenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	agileClient := meta.(*agile.Client)

	id := d.Id()
	tenant, err := agileClient.GetTenant(id)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setTenantAttributes(tenant, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileTenantUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Tenant: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	tenant, errTenant := NewTenantAttributes(d)

	if errTenant != nil {
		return errTenant
	}

	_, err := agileClient.UpdateTenant(agile.String(d.Id()), agile.String(name), tenant)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileTenantRead(ctx, d, meta)
}

func resourceAgileTenantDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	err := agileClient.DeleteTenant(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func resourceAgileTenantImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()

	tenant, err := agileClient.GetTenant(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setTenantAttributes(tenant, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func NewTenantAttributes(d *schema.ResourceData) (*models.TenantAttributes, diag.Diagnostics) {
	tenantAttr := models.TenantAttributes{
		Producer: agile.String(d.Get("producer").(string)),
	}

	if d.HasChange("description") {
		old, new := d.GetChange("description")
		if old != nil && new == nil {
			tenantAttr.Description = agile.String("")
		} else {
			tenantAttr.Description = agile.String(d.Get("description").(string))
		}
	} else {
		if _, ok := d.GetOk("description"); ok {
			tenantAttr.Description = agile.String(d.Get("description").(string))
		}
	}

	if _, ok := d.GetOk("description"); ok {
		tenantAttr.Description = agile.String(d.Get("description").(string))
	}

	quota := d.Get("quota").(*schema.Set).List()[0].(map[string]interface{})
	tenantAttr.Quota = &models.TenantQuota{
		LogicSwitchNum: agile.Int32(int32(quota["logic_switch_num"].(int))),
		LogicVasNum:    agile.Int32(int32(quota["logic_vas_num"].(int))),
		LogicRouterNum: agile.Int32(int32(quota["logic_router_num"].(int))),
	}
	if _, ok := d.GetOk("multicast_quota"); ok {
		multicastQuota := d.Get("multicast_quota").([]interface{})
		multicastQuotaItem := multicastQuota[0].(map[string]interface{})
		multicastQuotaAttr := models.TenantMulticastQuota{}
		if val, ok := multicastQuotaItem["acl_num"]; ok {
			multicastQuotaAttr.AclNum = agile.Int32(int32(val.(int)))
		}
		if val, ok := multicastQuotaItem["acl_rule_num"]; ok {
			multicastQuotaAttr.AclRuleNum = agile.Int32(int32(val.(int)))
		}
		tenantAttr.MulticastQuota = &multicastQuotaAttr
		tenantAttr.MulticastCapability = agile.Bool(true)
	} else {
		tenantAttr.MulticastCapability = agile.Bool(false)
	}

	if _, ok := d.GetOk("res_pool"); ok {
		resPool := d.Get("res_pool").([]interface{})
		resPoolItem := resPool[0].(map[string]interface{})
		resPoolAttr := models.TenantResPool{}

		if resPoolItem["external_gateway_ids"] != nil {
			resPoolAttr.ExternalGatewayIds = tools.ExtractSliceOfStrings(resPoolItem["external_gateway_ids"].(*schema.Set).List())
		}

		if resPoolItem["fabric_ids"] != nil {
			resPoolAttr.FabricIds = tools.ExtractSliceOfStrings(resPoolItem["fabric_ids"].(*schema.Set).List())

		}

		if resPoolItem["vmm_ids"] != nil {
			resPoolAttr.VmmIds = tools.ExtractSliceOfStrings(resPoolItem["vmm_ids"].(*schema.Set).List())
		}

		if resPoolItem["dhcp_group_ids"] != nil {
			resPoolAttr.DhcpGroupIds = tools.ExtractSliceOfStrings(resPoolItem["dhcp_group_ids"].(*schema.Set).List())
		}
		tenantAttr.ResPool = &resPoolAttr
	}

	return &tenantAttr, nil

}

func setTenantAttributes(tenant *models.Tenant, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("name", *tenant.Name)
	d.Set("multicast_capability", *tenant.MulticastCapability)

	if tenant.Producer != nil {
		d.Set("producer", *tenant.Producer)
	}

	d.Set("quota", []interface{}{
		map[string]interface{}{
			"logic_vas_num":    *tenant.Quota.LogicVasNum,
			"logic_router_num": *tenant.Quota.LogicRouterNum,
			"logic_switch_num": *tenant.Quota.LogicSwitchNum,
		},
	})
	if tenant.MulticastQuota != nil {
		aclNum := *tenant.MulticastQuota.AclNum
		aclRuleNum := *tenant.MulticastQuota.AclRuleNum
		if aclNum != 0 || aclRuleNum != 0 {
			d.Set("multicast_quota", []interface{}{
				map[string]interface{}{
					"acl_num":      aclNum,
					"acl_rule_num": aclRuleNum,
				},
			})
		}
	}

	if tenant.ResPool != nil {
		if len(tenant.ResPool.ExternalGatewayIds) != 0 || len(tenant.ResPool.FabricIds) != 0 || len(tenant.ResPool.DhcpGroupIds) != 0 || len(tenant.ResPool.VmmIds) != 0 {
			d.Set("res_pool", []interface{}{
				map[string]interface{}{
					"external_gateway_ids": tools.CreateSliceOfStrings(tenant.ResPool.ExternalGatewayIds),
					"fabric_ids":           tools.CreateSliceOfStrings(tenant.ResPool.FabricIds),
					"dhcp_group_ids":       tools.CreateSliceOfStrings(tenant.ResPool.DhcpGroupIds),
					"vmm_ids":              tools.CreateSliceOfStrings(tenant.ResPool.VmmIds),
				},
			})
		}
	}

	return d, nil
}
