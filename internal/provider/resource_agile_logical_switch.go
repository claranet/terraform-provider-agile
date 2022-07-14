package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	uuid "github.com/nu7hatch/gouuid"
)

func resourceAgileLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages Logical Switches.",
		CreateContext: resourceAgileLogicalSwitchCreate,
		ReadContext:   resourceAgileLogicalSwitchRead,
		UpdateContext: resourceAgileLogicalSwitchUpdate,
		DeleteContext: resourceAgileLogicalSwitchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileLogicalSwitchImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Logical switch ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Logical switch name.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 255),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Logical switch description.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"logic_network_id": {
				Type:         schema.TypeString,
				Description:  "Logical network where a logical switch is located.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"vni": {
				Type:         schema.TypeInt,
				Description:  "Logical switch VNI.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16000000),
			},
			"bd": {
				Type:         schema.TypeInt,
				Description:  "BD ID of a logical switch.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16000000),
			},
			"mac_address": {
				Type:         schema.TypeString,
				Description:  "MAC address of a logical switch.",
				Optional:     true,
				ValidateFunc: validation.IsMACAddress,
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Description:  "Tenant ID. In the northbound direction, the value can be either specified or not. The controller can automatically obtain the tenant ID from a logical network.",
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},
			"storm_suppress": {
				Type:        schema.TypeList,
				Description: "Storm Suppress Settings.",
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broadcast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the broadcast function",
							Optional:    true,
							Default:     false,
						},
						"multicast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the multicast function.",
							Optional:    true,
							Default:     false,
						},
						"unicast_enable": {
							Type:        schema.TypeBool,
							Description: "Whether to enable the unicast function. ",
							Optional:    true,
							Default:     false,
						},
						"broadcast_cbs": {
							Type:         schema.TypeInt,
							Description:  "CBS of broadcast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"broadcast_cbs_unit": {
							Type:         schema.TypeString,
							Description:  "CBS unit of broadcast packets. The value can be bytes, Kbytes, or Mbytes.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbytes", "mbytes"}, false),
						},
						"broadcast_cir": {
							Type:         schema.TypeInt,
							Description:  "Broadcast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"broadcast_cir_unit": {
							Type:         schema.TypeString,
							Description:  "CIR unit of broadcast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbps", "mbytes"}, false),
						},
						"unicast_cbs": {
							Type:         schema.TypeInt,
							Description:  "CBS of unicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"unicast_cbs_unit": {
							Type:         schema.TypeString,
							Description:  "CBS unit of unicast packets. The value can be byte, Kbytes, or Mbytes.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbyte", "mbytes"}, false),
						},
						"unicast_cir": {
							Type:         schema.TypeInt,
							Description:  "Unicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"unicast_cir_unit": {
							Type:         schema.TypeString,
							Description:  "CIR unit of unicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbps", "mbytes"}, false),
						},
						"multicast_cbs": {
							Type:         schema.TypeInt,
							Description:  "CBS of multicast packets. The value range is from 10000 to 4294967295 in bytes, 9 to 4194303 in Kbytes, or 1 to 4095 in Mbytes.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"multicast_cbs_unit": {
							Type:         schema.TypeString,
							Description:  "CBS unit of multicast packets. The value can be bytes, Kbytes, or Mbytes.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbytes", "mbytes"}, false),
						},
						"multicast_cir": {
							Type:         schema.TypeInt,
							Description:  "Multicast CIR. The value range is from 0 to 4294967295 in kbit/s, 0 to 4294967 in Mbit/s, or 0 to 4294 in Gbit/s.",
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"multicast_cir_unit": {
							Type:         schema.TypeString,
							Description:  "CIR unit of multicast packets. The value can be Gbit/s, Mbit/s, or kbit/s.",
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"byte", "kbps", "mbytes"}, false),
						},
					},
				},
			},
			"additional": {
				Type:        schema.TypeList,
				Description: "Additional Settings.",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"producer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used.",
							Default:     "default",
							ForceNew:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringLenBetween(0, 36),
							),
						},
					},
				},
			},
		},
	}
}

func resourceAgileLogicalSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Switcg: Beginning Creation")

	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	logicalSwitch, err := NewLogicalSwitchAttributes(d)

	if err != nil {
		return err
	}

	if err := agileClient.CreateLogicalSwitch(agile.String(id.String()), agile.String(name), logicalSwitch); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAgileLogicalSwitchRead(ctx, d, meta)

}

func resourceAgileLogicalSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalSwitch, err := agileClient.GetLogicalSwitch(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	if _, err := setLogicalSwitchAttributes(logicalSwitch, d); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileLogicalSwitchUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Logical Switch: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	logicalSwitchAttr, err := NewLogicalSwitchAttributes(d)

	if err != nil {
		return err
	}

	if _, err := agileClient.UpdateLogicalSwitch(agile.String(d.Id()), agile.String(name), logicalSwitchAttr); err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileLogicalSwitchRead(ctx, d, meta)
}

func resourceAgileLogicalSwitchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	if err := agileClient.DeleteLogicalSwitch(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func resourceAgileLogicalSwitchImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalSwitch, err := agileClient.GetLogicalSwitch(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setLogicalSwitchAttributes(logicalSwitch, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func NewLogicalSwitchAttributes(d *schema.ResourceData) (*models.LogicalSwitchAttributes, diag.Diagnostics) {
	logicalSwitchAttr := models.LogicalSwitchAttributes{}

	if _, ok := d.GetOk("description"); ok {
		logicalSwitchAttr.Description = agile.String(d.Get("description").(string))
	}

	if _, ok := d.GetOk("logic_network_id"); ok {
		logicalSwitchAttr.LogicNetworkId = agile.String(d.Get("logic_network_id").(string))
	}

	if _, ok := d.GetOk("vni"); ok {
		logicalSwitchAttr.Vni = agile.Int32(int32(d.Get("vni").(int)))
	}

	if _, ok := d.GetOk("bd"); ok {
		logicalSwitchAttr.Bd = agile.Int32(int32(d.Get("bd").(int)))
	}

	if _, ok := d.GetOk("mac_address"); ok {
		logicalSwitchAttr.MacAddress = agile.String(d.Get("mac_address").(string))
	}

	if _, ok := d.GetOk("tenant_id"); ok {
		logicalSwitchAttr.TenantId = agile.String(d.Get("tenant_id").(string))
	}

	if val, ok := d.GetOk("storm_suppress"); ok {
		stormSuppress := val.([]interface{})[0].(map[string]interface{})
		logicalSwitchAttr.StormSuppress = &models.LogicalSwitchStormSuppress{}

		if val, ok := stormSuppress["broadcast_enable"]; ok {
			logicalSwitchAttr.StormSuppress.BroadcastEnable = agile.Bool(val.(bool))
			if *logicalSwitchAttr.StormSuppress.BroadcastEnable {
				if stormSuppress["broadcast_cbs"].(int) == 0 || len(stormSuppress["broadcast_cbs_unit"].(string)) == 0 || stormSuppress["broadcast_cir"].(int) == 0 || len(stormSuppress["broadcast_cir_unit"].(string)) == 0 {
					return nil, diag.Errorf("broadcast_cbs, broadcast_cbs_unit, broadcast_cir, broadcast_cir_unit parameters must be set when the broadcast is enabled.")
				}
				logicalSwitchAttr.StormSuppress.BroadcastCbs = agile.String(fmt.Sprint(stormSuppress["broadcast_cbs"].(int)))
				logicalSwitchAttr.StormSuppress.BroadcastCbsUnit = agile.String(stormSuppress["broadcast_cbs_unit"].(string))
				logicalSwitchAttr.StormSuppress.BroadcastCir = agile.Int64(int64(stormSuppress["broadcast_cir"].(int)))
				logicalSwitchAttr.StormSuppress.BroadcastCirUnit = agile.String(stormSuppress["broadcast_cir_unit"].(string))
			} else if stormSuppress["broadcast_cbs"].(int) != 0 || len(stormSuppress["broadcast_cbs_unit"].(string)) != 0 || stormSuppress["broadcast_cir"].(int) != 0 || len(stormSuppress["broadcast_cir_unit"].(string)) != 0 {
				return nil, diag.Errorf("broadcast_cbs, broadcast_cbs_unit, broadcast_cir, broadcast_cir_unit parameters cannot be set when the broadcast is not enabled.")
			}
		}

		if val, ok := stormSuppress["multicast_enable"]; ok {
			logicalSwitchAttr.StormSuppress.MulticastEnable = agile.Bool(val.(bool))
			if *logicalSwitchAttr.StormSuppress.MulticastEnable {
				if stormSuppress["multicast_cbs"].(int) == 0 || len(stormSuppress["multicast_cbs_unit"].(string)) == 0 || stormSuppress["multicast_cir"].(int) == 0 || len(stormSuppress["multicast_cir_unit"].(string)) == 0 {
					return nil, diag.Errorf("multicast_cbs, multicast_cbs_unit, multicast_cir, multicast_cir_unit parameters must be set when the multicast is enabled.")
				}
				logicalSwitchAttr.StormSuppress.MulticastCbs = agile.String(fmt.Sprint(stormSuppress["multicast_cbs"].(int)))
				logicalSwitchAttr.StormSuppress.MulticastCbsUnit = agile.String(stormSuppress["multicast_cbs_unit"].(string))
				logicalSwitchAttr.StormSuppress.MulticastCir = agile.Int64(int64(stormSuppress["multicast_cir"].(int)))
				logicalSwitchAttr.StormSuppress.MulticastCirUnit = agile.String(stormSuppress["multicast_cir_unit"].(string))
			} else if stormSuppress["multicast_cbs"].(int) != 0 || len(stormSuppress["multicast_cbs_unit"].(string)) != 0 || stormSuppress["multicast_cir"].(int) != 0 || len(stormSuppress["multicast_cir_unit"].(string)) != 0 {
				return nil, diag.Errorf("multicast_cbs, multicast_cbs_unit, multicast_cir, multicast_cir_unit parameters cannot be set when the multicast is not enabled.")
			}
		}

		if val, ok := stormSuppress["unicast_enable"]; ok {
			logicalSwitchAttr.StormSuppress.UnicastEnable = agile.Bool(val.(bool))
			if *logicalSwitchAttr.StormSuppress.UnicastEnable {
				if stormSuppress["unicast_cbs"].(int) == 0 || len(stormSuppress["unicast_cbs_unit"].(string)) == 0 || stormSuppress["unicast_cir"].(int) == 0 || len(stormSuppress["unicast_cir_unit"].(string)) == 0 {
					return nil, diag.Errorf("unicast_cbs, unicast_cbs_unit, unicast_cir, unicast_cir_unit parameters must be set when the unicast is enabled.")
				}
				logicalSwitchAttr.StormSuppress.UnicastCbs = agile.String(fmt.Sprint(stormSuppress["unicast_cbs"].(int)))
				logicalSwitchAttr.StormSuppress.UnicastCbsUnit = agile.String(stormSuppress["unicast_cbs_unit"].(string))
				logicalSwitchAttr.StormSuppress.UnicastCir = agile.Int64(int64(stormSuppress["unicast_cir"].(int)))
				logicalSwitchAttr.StormSuppress.UnicastCirUnit = agile.String(stormSuppress["unicast_cir_unit"].(string))
			} else if stormSuppress["unicast_cbs"].(int) != 0 || len(stormSuppress["unicast_cbs_unit"].(string)) != 0 || stormSuppress["unicast_cir"].(int) != 0 || len(stormSuppress["unicast_cir_unit"].(string)) != 0 {
				return nil, diag.Errorf("unicast_cbs, unicast_cbs_unit, unicast_cir, unicast_cir_unit parameters cannot be set when the unicast is not enabled.")
			}
		}
	}

	if val, ok := d.GetOk("additional"); ok {
		additional := val.([]interface{})[0].(map[string]interface{})
		logicalSwitchAttr.Additional = &models.LogicalSwitchAdditional{}
		if val, ok := additional["producer"]; ok {
			logicalSwitchAttr.Additional.Producer = agile.String(val.(string))
		}
	}

	return &logicalSwitchAttr, nil

}

func setLogicalSwitchAttributes(logicalSwitch *models.LogicalSwitch, d *schema.ResourceData) (*schema.ResourceData, error) {
	if err := d.Set("name", logicalSwitch.Name); err != nil {
		return nil, err
	}
	if err := d.Set("description", logicalSwitch.Description); err != nil {
		return nil, err
	}
	if err := d.Set("logic_network_id", logicalSwitch.LogicNetworkId); err != nil {
		return nil, err
	}
	if err := d.Set("vni", logicalSwitch.Vni); err != nil {
		return nil, err
	}
	if err := d.Set("bd", logicalSwitch.Bd); err != nil {
		return nil, err
	}
	if err := d.Set("tenant_id", logicalSwitch.TenantId); err != nil {
		return nil, err
	}
	if err := d.Set("mac_address", logicalSwitch.MacAddress); err != nil {
		return nil, err
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
				return nil, err
			}
		}

		if logicalSwitch.StormSuppress.UnicastCbs != nil {
			if val, err := strconv.ParseInt(*logicalSwitch.StormSuppress.UnicastCbs, 10, 64); err == nil {
				stormSuppress[0].(map[string]interface{})["unicast_cbs"] = val
			} else {
				return nil, err
			}
		}

		if logicalSwitch.StormSuppress.MulticastCbs != nil {
			if val, err := strconv.ParseInt(*logicalSwitch.StormSuppress.MulticastCbs, 10, 64); err == nil {
				stormSuppress[0].(map[string]interface{})["multicast_cbs"] = val
			} else {
				return nil, err
			}
		}

		if err := d.Set("storm_suppress", stormSuppress); err != nil {
			return nil, err
		}
	}

	if logicalSwitch.Additional != nil {
		err := d.Set("additional", []interface{}{
			map[string]interface{}{
				"producer": logicalSwitch.Additional.Producer,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}
