package provider

import (
	"context"
	underscore "github.com/ahl5esoft/golang-underscore"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"terraform-provider-agile/tools"
)

func dataSourceAgileLogicalNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve Logical Network by name.",
		ReadContext: dataSourceAgileLogicalNetworkRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Logical network name.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 255),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"id": {
				Description: "Logical network ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Logical network description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": {
				Description: "Tenant to which a logical network (VPC) belongs.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"fabrics_id": {
				Description: "ID of the fabrics associated with the logical network.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"type": {
				Description: "Logical network type, which can be Instance or Transit.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"multicast_capability": {
				Description: "Whether the multicast capability is supported.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"additional": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Additional Settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"producer": {
							Type:        schema.TypeString,
							Description: "This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used.",
							Computed:    true,
						},
					},
				},
			},
			"is_vpc_deployed": {
				Type:        schema.TypeBool,
				Description: "Indicates if VPC is deployed",
				Computed:    true,
			},
		},
	}
}

func dataSourceAgileLogicalNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	logicalNetworks, err := agileClient.ListLogicalNetworks(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var logicalNetwork models.LogicalNetwork
	underscore.Chain(logicalNetworks).Find(func(l models.LogicalNetwork, _ int) bool {
		return *l.Name == name
	}).Value(&logicalNetwork)

	if logicalNetwork.Id == nil || *logicalNetwork.Id == "" {
		return diag.Errorf("No Logical Network with name %s found", name)
	}

	d.SetId(*logicalNetwork.Id)
	d.Set("name", logicalNetwork.Name)
	d.Set("description", *logicalNetwork.Description)
	d.Set("tenant_id", *logicalNetwork.TenantId)
	d.Set("fabrics_id", tools.CreateSliceOfStrings(logicalNetwork.FabricId))
	d.Set("multicast_capability", *logicalNetwork.MulticastCapability)
	d.Set("type", *logicalNetwork.Type)
	d.Set("is_vpc_deployed", *logicalNetwork.IsVpcDeployed)
	if logicalNetwork.Additional != nil {
		d.Set("additional", []interface{}{
			map[string]string{
				"producer": *logicalNetwork.Additional.Producer,
			},
		})
	}
	return nil
}
