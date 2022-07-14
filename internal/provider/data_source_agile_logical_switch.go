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
	"strconv"
)

func dataSourceAgileLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Description: "Data source can be used to retrieve Logical Switch by name.",
		ReadContext: dataSourceAgileLogicalSwitchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Logical switch name.",
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
				Description: "Logical Switch ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Logical switch description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"logic_network_id": {
				Type:        schema.TypeString,
				Description: "Logical network where a logical switch is located.",
				Computed:    true,
			},
			"bd": {
				Description: "BD ID of a logical switch.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"vni": {
				Type:        schema.TypeInt,
				Description: "Logical switch VNI.",
				Computed:    true,
			},
			"mac_address": {
				Type:        schema.TypeString,
				Description: "MAC address of a logical switch.",
				Computed:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Tenant ID. In the northbound direction, the value can be either specified or not. The controller can automatically obtain the tenant ID from a logical network.",
				Computed:    true,
			},
			"storm_suppress": {
				Type:        schema.TypeList,
				Description: "Storm Suppress Settings.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broadcast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the broadcast function",
							Computed:    true,
						},
						"multicast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the multicast function.",
							Computed:    true,
						},
						"unicast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the unicast function. ",
							Computed:    true,
						},
						"broadcast_cbs": {
							Type:        schema.TypeInt,
							Description: "CBS of broadcast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Computed:    true,
						},
						"broadcast_cbs_unit": {
							Type:        schema.TypeString,
							Description: "CBS unit of broadcast packets. The value can be bytes, Kbytes, or Mbytes.",
							Computed:    true,
						},
						"broadcast_cir": {
							Type:        schema.TypeInt,
							Description: "Broadcast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Computed:    true,
						},
						"broadcast_cir_unit": {
							Type:        schema.TypeString,
							Description: "CIR unit of broadcast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Computed:    true,
						},
						"unicast_cbs": {
							Type:        schema.TypeInt,
							Description: "CBS of unicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Computed:    true,
						},
						"unicast_cbs_unit": {
							Type:        schema.TypeString,
							Description: "CBS unit of unicast packets. The value can be byte, Kbytes, or Mbytes.",
							Computed:    true,
						},
						"unicast_cir": {
							Type:        schema.TypeInt,
							Description: "Unicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Computed:    true,
						},
						"unicast_cir_unit": {
							Type:        schema.TypeString,
							Description: "CIR unit of unicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Computed:    true,
						},
						"multicast_cbs": {
							Type:        schema.TypeInt,
							Description: "CBS of multicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Computed:    true,
						},
						"multicast_cbs_unit": {
							Type:        schema.TypeString,
							Description: "CBS unit of multicast packets. The value can be bytes, Kbytes, or Mbytes.",
							Computed:    true,
						},
						"multicast_cir": {
							Type:        schema.TypeInt,
							Description: "Multicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Computed:    true,
						},
						"multicast_cir_unit": {
							Type:        schema.TypeString,
							Description: "CIR unit of multicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Computed:    true,
						},
					},
				},
			},
			"additional": {
				Type:        schema.TypeList,
				Description: "Additional Settings.",
				Computed:    true,
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
		},
	}
}

func dataSourceAgileLogicalSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] %s: Beginning Read", name)

	agileClient := meta.(*agile.Client)

	logicalSwitches, err := agileClient.ListLogicalSwitches(nil)

	if err != nil {
		return diag.FromErr(err)
	}

	var logicalSwitch models.LogicalSwitch
	underscore.Chain(logicalSwitches).Find(func(l models.LogicalSwitch, _ int) bool {
		return *l.Name == name
	}).Value(&logicalSwitch)

	if logicalSwitch.Id == nil || *logicalSwitch.Id == "" {
		return diag.Errorf("No Logical switch with name %s found", name)
	}

	d.SetId(*logicalSwitch.Id)
	if err := d.Set("name", logicalSwitch.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", *logicalSwitch.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("logic_network_id", *logicalSwitch.LogicNetworkId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vni", *logicalSwitch.Vni); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("bd", *logicalSwitch.Bd); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mac_address", *logicalSwitch.MacAddress); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tenant_id", *logicalSwitch.TenantId); err != nil {
		return diag.FromErr(err)
	}

	if logicalSwitch.Additional != nil {
		if err := d.Set("additional", []interface{}{
			map[string]*string{
				"producer": logicalSwitch.Additional.Producer,
			},
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	if logicalSwitch.StormSuppress != nil {
		stormSuppress := []interface{}{
			map[string]interface{}{
				"broadcast_enable":   logicalSwitch.StormSuppress.BroadcastEnable,
				"multicast_enable":   logicalSwitch.StormSuppress.MulticastEnable,
				"unicast_enable":     logicalSwitch.StormSuppress.UnicastEnable,
				"broadcast_cbs_unit": logicalSwitch.StormSuppress.BroadcastCbsUnit,
				"broadcast_cir":      logicalSwitch.StormSuppress.BroadcastCir,
				"broadcast_cir_unit": logicalSwitch.StormSuppress.BroadcastCirUnit,
				"unicast_cbs_unit":   logicalSwitch.StormSuppress.UnicastCbsUnit,
				"unicast_cir":        logicalSwitch.StormSuppress.UnicastCir,
				"unicast_cir_unit":   logicalSwitch.StormSuppress.UnicastCirUnit,
				"multicast_cbs_unit": logicalSwitch.StormSuppress.MulticastCbsUnit,
				"multicast_cir":      logicalSwitch.StormSuppress.MulticastCir,
				"multicast_cir_unit": logicalSwitch.StormSuppress.MulticastCirUnit,
			},
		}

		if logicalSwitch.StormSuppress.BroadcastCbs != nil {
			if val, err := strconv.ParseInt(*logicalSwitch.StormSuppress.BroadcastCbs, 10, 64); err == nil {
				stormSuppress[0].(map[string]interface{})["broadcast_cbs"] = val
			} else {
				return diag.FromErr(err)
			}
		}

		if logicalSwitch.StormSuppress.UnicastCbs != nil {
			if val, err := strconv.ParseInt(*logicalSwitch.StormSuppress.UnicastCbs, 10, 64); err == nil {
				stormSuppress[0].(map[string]interface{})["unicast_cbs"] = val
			} else {
				return diag.FromErr(err)
			}
		}

		if logicalSwitch.StormSuppress.MulticastCbs != nil {
			if val, err := strconv.ParseInt(*logicalSwitch.StormSuppress.MulticastCbs, 10, 64); err == nil {
				stormSuppress[0].(map[string]interface{})["multicast_cbs"] = val
			} else {
				return diag.FromErr(err)
			}
		}

		if err := d.Set("storm_suppress", stormSuppress); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
