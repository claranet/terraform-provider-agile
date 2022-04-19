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

func dataSourceAgileExternalGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve External Gateways by name.",
		ReadContext: dataSourceAgileExternalGatewayRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "External gateway name.",
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
				Description: "External gateway ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "External gateway description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"gateway_type": {
				Description: "External gateway type, which can be Public or Private.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_telco_gateway": {
				Description: "Indicates if is a Telco cloud gateway.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"vrf_name": {
				Description: "VRF Name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAgileExternalGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	externalGateways, err := agileClient.ListExternalGateways(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var externalGateway models.ExternalGateway
	underscore.Chain(externalGateways).Find(func(e models.ExternalGateway, _ int) bool {
		return *e.Name == name
	}).Value(&externalGateway)

	if externalGateway.Id == nil || *externalGateway.Id == "" {
		return diag.Errorf("No External Gateway with name %s found", name)
	}

	d.SetId(*externalGateway.Id)
	d.Set("name", *externalGateway.Name)
	d.Set("description", *externalGateway.Description)
	d.Set("gateway_type", *externalGateway.GatewayType)
	if externalGateway.IsTelcoGateway != nil {
		d.Set("is_telco_gateway", *externalGateway.IsTelcoGateway)
	}
	if externalGateway.VrfName != nil {
		d.Set("vrf_name", *externalGateway.VrfName)
	}
	return nil
}
