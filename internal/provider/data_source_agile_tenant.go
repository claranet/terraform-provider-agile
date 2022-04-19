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
)

func dataSourceAgileTenant() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve Tenant by name.",
		ReadContext: dataSourceAgileTenantRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Tenant name.",
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
				Description: "Tenant ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Tenant description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"producer": {
				Description: "Producer.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"multicast_capability": {
				Description: "Whether the multicast capability is supported.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceAgileTenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	tenants, err := agileClient.ListTenants(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var tenant models.Tenant
	underscore.Chain(tenants).Find(func(t models.Tenant, _ int) bool {
		return *t.Name == name
	}).Value(&tenant)

	if tenant.Id == nil || *tenant.Id == "" {
		return diag.Errorf("No Tenant with name %s found", name)
	}

	d.SetId(*tenant.Id)
	d.Set("name", tenant.Name)
	d.Set("description", *tenant.Description)
	d.Set("producer", *tenant.Producer)
	d.Set("multicast_capability", *tenant.MulticastCapability)
	return nil
}
