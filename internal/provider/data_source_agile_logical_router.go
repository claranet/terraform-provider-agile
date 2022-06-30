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
	"regexp"
)

func dataSourceAgileLogicalRouter() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve Logical Router by name.",
		ReadContext: dataSourceAgileLogicalRouterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Logical router name.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 255),
						validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_]*$`), ""),
					),
				),
			},
			"id": {
				Description: "Logical router ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Logical router description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"logic_network_id": {
				Type:        schema.TypeString,
				Description: "Logical network where a logical router is located.",
				Computed:    true,
			},
			"type": {
				Description: "Logical router type, which can be Normal, Nfvi, MultiActive, Transit, or Connect. This field cannot be updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vni": {
				Type:        schema.TypeInt,
				Description: "Online/offline status of a device.",
				Computed:    true,
			},
			"vrf_name": {
				Type:        schema.TypeString,
				Description: "VRF name.",
				Computed:    true,
			},
			"router_locations": {
				Type:        schema.TypeSet,
				Description: "Router Locations Settings.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fabric_id": {
							Type:        schema.TypeString,
							Description: "Fabric ID",
							Computed:    true,
						},
						"fabric_role": {
							Type:        schema.TypeString,
							Description: "Fabric role, which can be master or backup.",
							Computed:    true,
						},
						"fabric_name": {
							Type:        schema.TypeString,
							Description: "Fabric name.",
							Computed:    true,
						},
						"device_group": {
							Type:        schema.TypeSet,
							Description: "Device group. Devices in the list must belong to the same device group.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:        schema.TypeString,
										Description: "Specified physical device.",
										Computed:    true,
									},
									"device_ip": {
										Type:        schema.TypeString,
										Description: "Device management IP address.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAgileLogicalRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	logicalRouters, err := agileClient.ListLogicalRouters(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var logicalRouter models.LogicalRouter
	underscore.Chain(logicalRouters).Find(func(l models.LogicalRouter, _ int) bool {
		return *l.Name == name
	}).Value(&logicalRouter)

	if logicalRouter.Id == nil || *logicalRouter.Id == "" {
		return diag.Errorf("No Logical router with name %s found", name)
	}

	d.SetId(*logicalRouter.Id)
	d.Set("name", logicalRouter.Name)
	d.Set("description", *logicalRouter.Description)
	d.Set("logic_network_id", *logicalRouter.LogicNetworkId)
	d.Set("vni", *logicalRouter.Vni)
	d.Set("vrf_name", *logicalRouter.VrfName)
	d.Set("type", *logicalRouter.Type)

	var routerLocations []interface{}
	var deviceGroups []interface{}
	for _, deviceGroup := range logicalRouter.RouterLocations[0].DeviceGroup {
		deviceGroups = append(deviceGroups, map[string]interface{}{
			"device_id": *deviceGroup.DeviceId,
			"device_ip": *deviceGroup.DeviceIp,
		})
	}

	routerLocations = append(routerLocations, map[string]interface{}{
		"fabric_role":  *logicalRouter.RouterLocations[0].FabricRole,
		"fabric_id":    *logicalRouter.RouterLocations[0].FabricId,
		"fabric_name":  *logicalRouter.RouterLocations[0].FabricName,
		"device_group": deviceGroups,
	})

	d.Set("router_locations", routerLocations)
	return nil
}
