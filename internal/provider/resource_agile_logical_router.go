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
	"regexp"
)

func resourceAgileLogicalRouter() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages Logical Routers.",
		CreateContext: resourceAgileLogicalRouterCreate,
		ReadContext:   resourceAgileLogicalRouterRead,
		UpdateContext: resourceAgileLogicalRouterUpdate,
		DeleteContext: resourceAgileLogicalRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileLogicalRouterImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Logical router ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Logical router name.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 255),
						validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_]*$`), ""),
					),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Logical router description.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"logic_network_id": {
				Type:         schema.TypeString,
				Description:  "Logical network where a logical router is located.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"type": {
				Description: "Logical router type, which can be Normal, Nfvi, MultiActive, Transit, or Connect. This field cannot be updated.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice([]string{"Normal", "Nfvi", "MultiActive", "Transit", "Connect"}, false),
				),
			},
			"vni": {
				Type:         schema.TypeInt,
				Description:  "Online/offline status of a device.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 16000000),
			},
			"vrf_name": {
				Type:        schema.TypeString,
				Description: "VRF name.",
				Optional:    true,
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(1, 31),
						validation.StringDoesNotContainAny(" "),
					),
				),
			},
			"router_locations": {
				Type:        schema.TypeSet,
				Description: "Router Locations Settings.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fabric_id": {
							Type:         schema.TypeString,
							Description:  "Fabric ID",
							ForceNew:     true,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"fabric_role": {
							Type:        schema.TypeString,
							Description: "Fabric role, which can be master or backup.",
							ForceNew:    true,
							Optional:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{"master", "backup"}, true),
							),
						},
						"fabric_name": {
							Type:        schema.TypeString,
							Description: "Fabric name.",
							Computed:    true,
						},
						"device_group": {
							Type:        schema.TypeSet,
							Description: "Device group. Devices in the list must belong to the same device group.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:         schema.TypeString,
										Description:  "Specified physical device.",
										Required:     true,
										ValidateFunc: validation.IsUUID,
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
			//"additional": {
			//	Type:        schema.TypeSet,
			//	Description: "Additional Settings.",
			//	Optional:    true,
			//	MaxItems:    1,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"producer": {
			//				Type:        schema.TypeString,
			//				Optional:    true,
			//				Description: "This parameter is optional. If it is specified by the user, the specified value is used. The character string starting with component is reserved. If no value is specified, the default value default is used.",
			//				Default:     "default",
			//				ValidateDiagFunc: validation.ToDiagFunc(
			//					validation.StringLenBetween(0, 36),
			//				),
			//			},
			//		},
			//	},
			//},
		},
	}
}

func resourceAgileLogicalRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Router: Beginning Creation")

	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	logicalRouter, err := NewLogicalRouterAttributes(d)

	if err != nil {
		return err
	}

	if err := agileClient.CreateLogicalRouter(agile.String(id.String()), agile.String(name), logicalRouter); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAgileLogicalRouterRead(ctx, d, meta)

}

func resourceAgileLogicalRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalRouter, err := agileClient.GetLogicalRouter(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	if _, err := setLogicalRouterAttributes(logicalRouter, d); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileLogicalRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Logical Router: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	logicalRouterAttr, err := NewLogicalRouterAttributes(d)

	if err != nil {
		return err
	}

	if _, err := agileClient.UpdateLogicalRouter(agile.String(d.Id()), agile.String(name), logicalRouterAttr); err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileLogicalRouterRead(ctx, d, meta)
}

func resourceAgileLogicalRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	if err := agileClient.DeleteLogicalRouter(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func NewLogicalRouterAttributes(d *schema.ResourceData) (*models.LogicalRouterAttributes, diag.Diagnostics) {
	logicalRouterAttr := models.LogicalRouterAttributes{}

	if _, ok := d.GetOk("description"); ok {
		logicalRouterAttr.Description = agile.String(d.Get("description").(string))
	}

	if _, ok := d.GetOk("logic_network_id"); ok {
		logicalRouterAttr.LogicNetworkId = agile.String(d.Get("logic_network_id").(string))
	}

	if _, ok := d.GetOk("type"); ok {
		logicalRouterAttr.Type = agile.String(d.Get("type").(string))
	}

	if _, ok := d.GetOk("vni"); ok {
		logicalRouterAttr.Vni = agile.Int32(int32(d.Get("vni").(int)))
	}

	if _, ok := d.GetOk("vrf_name"); ok {
		logicalRouterAttr.VrfName = agile.String(d.Get("vrf_name").(string))
	}

	if val, ok := d.GetOk("router_locations"); ok {
		logicalRouterAttr.RouterLocations = make([]*models.LogicalRouterLocations, 0)
		location := val.(*schema.Set).List()[0].(map[string]interface{})
		var locationItem models.LogicalRouterLocations

		if locationVal, ok := location["fabric_id"]; ok {
			locationItem.FabricId = agile.String(locationVal.(string))
		}

		if locationVal, ok := location["fabric_role"]; ok {
			locationItem.FabricRole = agile.String(locationVal.(string))
		}

		if deviceGroupVal, ok := location["device_group"]; ok {
			routerDeviceGroup := make([]*models.LogicalRouterLocationsDeviceGroup, 0)
			for _, deviceGroupItem := range deviceGroupVal.(*schema.Set).List() {
				deviceGroup := deviceGroupItem.(map[string]interface{})
				if val, ok := deviceGroup["device_id"]; ok {
					routerDeviceGroup = append(routerDeviceGroup, &models.LogicalRouterLocationsDeviceGroup{
						DeviceIp: agile.String(val.(string)),
					})
				}
			}
			locationItem.DeviceGroup = routerDeviceGroup
		}

		logicalRouterAttr.RouterLocations = append(logicalRouterAttr.RouterLocations, &locationItem)
	}

	return &logicalRouterAttr, nil

}

func resourceAgileLogicalRouterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalRouter, err := agileClient.GetLogicalRouter(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setLogicalRouterAttributes(logicalRouter, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func setLogicalRouterAttributes(logicalRouter *models.LogicalRouter, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("name", *logicalRouter.Name)
	d.Set("description", *logicalRouter.Description)
	d.Set("logic_network_id", *logicalRouter.LogicNetworkId)
	d.Set("vni", *logicalRouter.Vni)
	d.Set("vrf_name", *logicalRouter.VrfName)
	d.Set("type", *logicalRouter.Type)

	var routerLocations []interface{}
	for _, location := range logicalRouter.RouterLocations {
		var deviceGroups []interface{}
		for _, deviceGroup := range location.DeviceGroup {
			deviceGroups = append(deviceGroups, map[string]interface{}{
				"device_id": *deviceGroup.DeviceId,
				"device_ip": *deviceGroup.DeviceIp,
			})
		}

		routerLocations = append(routerLocations, map[string]interface{}{
			"fabric_role":  *location.FabricRole,
			"fabric_id":    *location.FabricId,
			"fabric_name":  *location.FabricName,
			"device_group": deviceGroups,
		})
	}

	d.Set("router_locations", routerLocations)
	//if _, ok := d.GetOk("additional"); ok {
	//	d.Set("additional", []interface{}{
	//		map[string]string{
	//			"producer": *logicalRouter.Additional.Producer,
	//		},
	//	})
	//}
	return d, nil
}
