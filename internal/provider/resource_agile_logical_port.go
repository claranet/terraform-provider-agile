package provider

import (
	"context"
	agile "github.com/claranet/agilec-go-client/client"
	"github.com/claranet/agilec-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"strings"
)

func resourceAgileLogicalPort() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages Logical Ports.",
		CreateContext: resourceAgileLogicalPortCreate,
		ReadContext:   resourceAgileLogicalPortRead,
		UpdateContext: resourceAgileLogicalPortUpdate,
		DeleteContext: resourceAgileLogicalPortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileLogicalPortImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Logical port ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Logical port name.",
				Required:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 255),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Logical port description.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 242),
				),
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Description:  "Tenant to which a logical port belongs. This parameter is automatically obtained by the controller.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"fabric_id": {
				Type:         schema.TypeString,
				Description:  "Fabric to which the logical port belongs.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"logic_switch_id": {
				Type:         schema.TypeString,
				Description:  "Logical switch to which a logical port belongs.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"access_info": {
				Type:        schema.TypeSet,
				Description: "Access info Settings.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Description: "Port mode, which can be UNI or NNI.",
							ForceNew:    true,
							Required:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{"UNI", "NNI"}, true),
							),
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Logical port type, which can be DOT1Q, DEFAULT, UNTAG, or QINQ.",
							ForceNew:    true,
							Required:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{"DOT1Q", "DEFAULT", "UNTAG", "QINQ"}, true),
							),
						},
						"vlan": {
							Type:        schema.TypeInt,
							Description: "Access VLAN ID. This parameter is mandatory when the access type is dot1q.",
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(2, 4094),
							),
						},
						"qinq": {
							Type:        schema.TypeSet,
							Description: "Qinq Settings.",
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inner_vid_begin": {
										Type:        schema.TypeInt,
										Description: "Start VLAN ID of the inner VLAN tag for QinQ.",
										Required:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.IntBetween(2, 4094),
										),
									},
									"inner_vid_end": {
										Type:        schema.TypeInt,
										Description: "End VLAN ID of the inner VLAN tag for QinQ.",
										Optional:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.IntBetween(2, 4094),
										),
									},
									"outer_vid_begin": {
										Type:        schema.TypeInt,
										Description: "Start VLAN ID of the outer VLAN tag for QinQ.",
										Required:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.IntBetween(2, 4094),
										),
									},
									"outer_vid_end": {
										Type:        schema.TypeInt,
										Description: "End VLAN ID of the outer VLAN tag for QinQ.",
										Optional:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.IntBetween(2, 4094),
										),
									},
									"rewrite_action": {
										Type:        schema.TypeString,
										Description: "Rewrite action of QinQ, which can be POPDOUBLE or PASSTHROUGH.",
										Required:    true,
										//ForceNew:    true,
										ValidateDiagFunc: validation.ToDiagFunc(
											validation.StringInSlice([]string{"POPDOUBLE", "PASSTHROUGH"}, true),
										),
									},
								},
							},
						},
						"location": {
							Type:        schema.TypeSet,
							Description: "Location Settings.",
							MaxItems:    4000,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_group_id": {
										Type:         schema.TypeString,
										Description:  "Device group ID of a physical device.",
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsUUID,
									},
									"device_id": {
										Type:         schema.TypeString,
										Description:  "Specified physical device.",
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsUUID,
									},
									"port_id": {
										Type:         schema.TypeString,
										Description:  "Specified physical port.",
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IsUUID,
									},
									"port_name": {
										Type:        schema.TypeString,
										Description: "Port name.",
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
						"subinterface_number": {
							Type:        schema.TypeInt,
							Description: "Number of an access sub-interface.",
							ForceNew:    true,
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntBetween(1, 4094),
							),
						},
					},
				},
			},
			"additional": {
				Type:        schema.TypeSet,
				Description: "Additional Settings.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"producer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used.",
							Default:     "default",
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

func resourceAgileLogicalPortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Port: Beginning Creation")

	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	logicalPort, err := NewLogicalPortAttributes(d)

	if err != nil {
		return err
	}

	if err := agileClient.CreateLogicalPort(agile.String(id.String()), agile.String(name), logicalPort); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAgileLogicalPortRead(ctx, d, meta)
}

func resourceAgileLogicalPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read Logical Port", d.Id())
	agileClient := meta.(*agile.Client)
	id := d.Id()
	logicalPort, err := agileClient.GetLogicalPort(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	if _, err := setLogicalPortAttributes(logicalPort, d); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileLogicalPortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Logical Port: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	logicalPortAttr, errAttr := NewLogicalPortAttributes(d)

	if errAttr != nil {
		return errAttr
	}

	_, err := agileClient.UpdateLogicalPort(agile.String(d.Id()), agile.String(name), logicalPortAttr)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileLogicalPortRead(ctx, d, meta)
}

func resourceAgileLogicalPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	err := agileClient.DeleteLogicalPort(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func NewLogicalPortAttributes(d *schema.ResourceData) (*models.LogicalPortAttributes, diag.Diagnostics) {
	logicalPortAttr := models.LogicalPortAttributes{}

	if _, ok := d.GetOk("description"); ok {
		logicalPortAttr.Description = agile.String(d.Get("description").(string))
	}

	if _, ok := d.GetOk("tenant_id"); ok {
		logicalPortAttr.TenantId = agile.String(d.Get("tenant_id").(string))
	}

	if _, ok := d.GetOk("fabric_id"); ok {
		logicalPortAttr.FabricId = agile.String(d.Get("fabric_id").(string))
	}

	if _, ok := d.GetOk("logic_switch_id"); ok {
		logicalPortAttr.LogicSwitchId = agile.String(d.Get("logic_switch_id").(string))
	}

	if _, ok := d.GetOk("access_info"); ok {
		accessInfo := d.Get("access_info").(*schema.Set).List()[0].(map[string]interface{})
		logicalPortAttr.AccessInfo = &models.LogicalPortAccessInfo{}
		logicalPortAttr.AccessInfo.Mode = agile.String(accessInfo["mode"].(string))
		logicalPortAttr.AccessInfo.Type = agile.String(accessInfo["type"].(string))

		if val, ok := accessInfo["vlan"]; ok {
			logicalPortAttr.AccessInfo.Vlan = agile.Int32(int32(val.(int)))
		} else if strings.ToUpper(*logicalPortAttr.AccessInfo.Type) == strings.ToUpper("Dot1q") {
			return nil, diag.Errorf("Vlan parameter is mandatory when the access type is dot1q.")
		}

		if val, ok := accessInfo["subinterface_number"]; ok {
			logicalPortAttr.AccessInfo.SubinterfaceNumber = agile.Int32(int32(val.(int)))
		}

		if val, ok := accessInfo["qinq"]; ok {
			qinq := val.(*schema.Set).List()[0].(map[string]interface{})
			logicalPortAttr.AccessInfo.Qinq = &models.LogicalPortAccessInfoQinq{}
			if qinqVal, ok := qinq["inner_vid_begin"]; ok {
				logicalPortAttr.AccessInfo.Qinq.InnerVidBegin = agile.Int32(int32(qinqVal.(int)))
			}

			if qinqVal, ok := qinq["inner_vid_end"]; ok {
				logicalPortAttr.AccessInfo.Qinq.InnerVidEnd = agile.Int32(int32(qinqVal.(int)))
			}

			if qinqVal, ok := qinq["outer_vid_begin"]; ok {
				logicalPortAttr.AccessInfo.Qinq.OuterVidBegin = agile.Int32(int32(qinqVal.(int)))
			}

			if qinqVal, ok := qinq["outer_vid_end"]; ok {
				logicalPortAttr.AccessInfo.Qinq.OuterVidEnd = agile.Int32(int32(qinqVal.(int)))
			}

			if qinqVal, ok := qinq["rewrite_action"]; ok {
				logicalPortAttr.AccessInfo.Qinq.RewriteAction = agile.String(qinqVal.(string))
			}
		}

		if val, ok := accessInfo["location"]; ok {
			logicalPortAttr.AccessInfo.Location = make([]*models.LogicalPortAccessInfoLocation, 0)
			for _, item := range val.(*schema.Set).List() {
				location := item.(map[string]interface{})
				var locationItem models.LogicalPortAccessInfoLocation

				if locationVal, ok := location["device_group_id"]; ok {
					locationItem.DeviceGroupId = agile.String(locationVal.(string))
				}

				if location["device_id"] != "" {
					locationItem.DeviceId = agile.String(location["device_id"].(string))
				}

				if location["port_id"] != "" {
					locationItem.PortId = agile.String(location["port_id"].(string))
				}

				logicalPortAttr.AccessInfo.Location = append(logicalPortAttr.AccessInfo.Location, &locationItem)
			}
		}
	}

	if _, ok := d.GetOk("additional"); ok {
		additional := d.Get("additional").(*schema.Set).List()[0].(map[string]interface{})
		logicalPortAttr.Additional = &models.LogicalPortAdditional{}
		if val, ok := additional["producer"]; ok {
			logicalPortAttr.Additional.Producer = agile.String(val.(string))
		}
	}

	return &logicalPortAttr, nil

}

func resourceAgileLogicalPortImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalPort, err := agileClient.GetLogicalPort(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setLogicalPortAttributes(logicalPort, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func setLogicalPortAttributes(logicalPort *models.LogicalPort, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("name", *logicalPort.Name)
	d.Set("description", *logicalPort.Description)
	d.Set("tenant_id", *logicalPort.TenantId)
	d.Set("fabric_id", logicalPort.FabricId)
	d.Set("logic_switch_id", *logicalPort.LogicSwitchId)

	accessInfo := []interface{}{
		map[string]interface{}{
			"mode": *logicalPort.AccessInfo.Mode,
			"type": *logicalPort.AccessInfo.Type,
			"vlan": *logicalPort.AccessInfo.Vlan,
			"qinq": []interface{}{
				map[string]interface{}{
					"inner_vid_begin": *logicalPort.AccessInfo.Qinq.InnerVidBegin,
					"inner_vid_end":   *logicalPort.AccessInfo.Qinq.InnerVidEnd,
					"outer_vid_begin": *logicalPort.AccessInfo.Qinq.OuterVidBegin,
					"outer_vid_end":   *logicalPort.AccessInfo.Qinq.OuterVidEnd,
					"rewrite_action":  *logicalPort.AccessInfo.Qinq.RewriteAction,
				}},
			"location":            []interface{}{},
			"subinterface_number": *logicalPort.AccessInfo.SubinterfaceNumber,
		},
	}

	for _, location := range logicalPort.AccessInfo.Location {
		accessInfo[0].(map[string]interface{})["location"] = append(accessInfo[0].(map[string]interface{})["location"].([]interface{}), map[string]interface{}{
			"device_group_id": *location.DeviceGroupId,
			"device_id":       *location.DeviceId,
			"port_id":         *location.PortId,
			"port_name":       *location.PortName,
			"device_ip":       *location.DeviceIp,
		})
	}
	d.Set("access_info", accessInfo)

	if _, ok := d.GetOk("additional"); ok || logicalPort.AccessInfo != nil {
		d.Set("additional", []interface{}{
			map[string]string{
				"producer": *logicalPort.Additional.Producer,
			},
		})
	}
	return d, nil
}
