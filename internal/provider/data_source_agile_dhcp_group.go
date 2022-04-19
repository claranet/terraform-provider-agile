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

func dataSourceAgileDhcpGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve DHCP Group by name.",
		ReadContext: dataSourceAgileDhcpGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "DHCP server name.",
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
				Description: "DHCP server ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "DHCP server description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"producer": {
				Description: "Producer.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"logic_router_id": {
				Description: "ID of the logical router that the DHCP server belongs to.",
				Type:        schema.TypeString,
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

func dataSourceAgileDhcpGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	dhcpGroups, err := agileClient.ListDHCPGroups(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var dhcpGroup models.DHCPGroup
	underscore.Chain(dhcpGroups).Find(func(d models.DHCPGroup, _ int) bool {
		return *d.Name == name
	}).Value(&dhcpGroup)

	if dhcpGroup.Id == nil || *dhcpGroup.Id == "" {
		return diag.Errorf("No Dhcp Group with name %s found", name)
	}

	d.SetId(*dhcpGroup.Id)
	d.Set("name", *dhcpGroup.Name)
	d.Set("description", *dhcpGroup.Description)
	d.Set("producer", *dhcpGroup.Producer)
	d.Set("logic_router_id", *dhcpGroup.LogicRouterId)
	d.Set("vrf_name", *dhcpGroup.VrfName)
	return nil
}
