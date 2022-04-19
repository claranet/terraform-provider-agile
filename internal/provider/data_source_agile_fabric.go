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

func dataSourceAgileFabric() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve fabric by name.",
		ReadContext: dataSourceAgileFabricRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Fabric name.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 128),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"id": {
				Description: "Fabric ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Fabric description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"network_type": {
				Description: "Fabric VXLAN type. The value is Distributed or Centralized.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"physical_network_mode": {
				Description: "Fabric networking type. Only VXLAN is supported.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"micro_segment": {
				Description: "Whether to enable microsegmentation.",
				Type:        schema.TypeBool,
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

func dataSourceAgileFabricRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	fabrics, err := agileClient.ListFabrics(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var fabric models.Fabric
	underscore.Chain(fabrics).Find(func(f models.Fabric, _ int) bool {
		return *f.Name == name
	}).Value(&fabric)

	if fabric.Id == nil || *fabric.Id == "" {
		return diag.Errorf("No Fabric with name %s found", name)
	}

	d.SetId(*fabric.Id)
	d.Set("name", *fabric.Name)
	d.Set("description", *fabric.Description)
	d.Set("network_type", *fabric.NetworkType)
	d.Set("physical_network_mode", *fabric.PhysicalNetworkMode)
	d.Set("multicast_capability", *fabric.MulticastCapability)
	d.Set("micro_segment", *fabric.MicroSegmentCapability)
	return nil
}
