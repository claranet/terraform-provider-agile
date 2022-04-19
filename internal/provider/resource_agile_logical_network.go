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
	"terraform-provider-agile/tools"
)

func resourceAgileLogicalNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages Logical Networks.",
		CreateContext: resourceAgileLogicalNetworkCreate,
		ReadContext:   resourceAgileLogicalNetworkRead,
		UpdateContext: resourceAgileLogicalNetworkUpdate,
		DeleteContext: resourceAgileLogicalNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileLogicalNetworkImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Logical network ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Logical network name.",
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
				Description: "Logical network description.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"tenant_id": {
				Type:         schema.TypeString,
				Description:  "Tenant to which a logical network (VPC) belongs. If this parameter is left empty, it is a public VPC.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"fabrics_id": {
				Type:        schema.TypeList,
				Description: "ID of the fabrics associated with the logical network. If the VPC is not a public VPC, the fabric must have been added by the tenant.",
				Optional:    true,
				MinItems:    0,
				MaxItems:    4000,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"multicast_capability": {
				Type:        schema.TypeBool,
				Description: "Whether the multicast capability is supported.",
				Optional:    true,
				Default:     false,
			},
			"type": {
				Description: "Logical network type, which can be Instance or Transit. If this parameter is left empty, the default value is Instance.",
				Type:        schema.TypeString,
				Default:     "Instance",
				Optional:    true,
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice([]string{"Transit", "Instance"}, false),
				),
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
			"is_vpc_deployed": {
				Type:        schema.TypeBool,
				Description: "Indicates if VPC is deployed",
				Computed:    true,
			},
		},
	}
}

func resourceAgileLogicalNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Logical Network: Beginning Creation")

	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	logicalNetwork, errLogicalNetwork := NewLogicalNetworkAttributes(d)

	if errLogicalNetwork != nil {
		return errLogicalNetwork
	}

	err := agileClient.CreateLogicalNetwork(agile.String(id.String()), agile.String(name), logicalNetwork)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAgileLogicalNetworkRead(ctx, d, meta)

}

func resourceAgileLogicalNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalNetwork, err := agileClient.GetLogicalNetwork(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	_, errAttr := setLogicalNetworkAttributes(logicalNetwork, d)
	if err != nil {
		return diag.FromErr(errAttr)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileLogicalNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Logical Network: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	logicalNetworkAttr, errAttr := NewLogicalNetworkAttributes(d)

	if errAttr != nil {
		return errAttr
	}

	_, err := agileClient.UpdateLogicalNetwork(agile.String(d.Id()), agile.String(name), logicalNetworkAttr)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileLogicalNetworkRead(ctx, d, meta)
}

func resourceAgileLogicalNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	err := agileClient.DeleteLogicalNetwork(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func NewLogicalNetworkAttributes(d *schema.ResourceData) (*models.LogicalNetworkAttributes, diag.Diagnostics) {
	logicalNetworkAttr := models.LogicalNetworkAttributes{
		MulticastCapability: agile.Bool(d.Get("multicast_capability").(bool)),
		Type:                agile.String(d.Get("type").(string)),
	}

	if _, ok := d.GetOk("description"); ok {
		logicalNetworkAttr.Description = agile.String(d.Get("description").(string))
	}

	if _, ok := d.GetOk("tenant_id"); ok {
		logicalNetworkAttr.TenantId = agile.String(d.Get("tenant_id").(string))
	}

	if val, ok := d.GetOk("fabrics_id"); ok {
		logicalNetworkAttr.FabricId = tools.ExtractSliceOfStrings(val.([]interface{}))
	}

	if _, ok := d.GetOk("additional"); ok {
		additional := d.Get("additional").(*schema.Set).List()[0].(map[string]interface{})
		additionalAttr := models.LogicalNetworkAdditional{}
		if val, ok := additional["producer"]; ok {
			additionalAttr.Producer = agile.String(val.(string))
		}
		logicalNetworkAttr.Additional = &additionalAttr
	}

	return &logicalNetworkAttr, nil

}

func resourceAgileLogicalNetworkImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()
	logicalNetwork, err := agileClient.GetLogicalNetwork(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setLogicalNetworkAttributes(logicalNetwork, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func setLogicalNetworkAttributes(logicalNetwork *models.LogicalNetwork, d *schema.ResourceData) (*schema.ResourceData, error) {
	//d.SetId(tenant.Id)
	d.Set("name", *logicalNetwork.Name)
	d.Set("description", *logicalNetwork.Description)
	d.Set("tenant_id", *logicalNetwork.TenantId)
	d.Set("fabrics_id", tools.CreateSliceOfStrings(logicalNetwork.FabricId))
	d.Set("multicast_capability", *logicalNetwork.MulticastCapability)
	d.Set("type", *logicalNetwork.Type)
	d.Set("is_vpc_deployed", *logicalNetwork.IsVpcDeployed)
	if _, ok := d.GetOk("additional"); ok {
		d.Set("additional", []interface{}{
			map[string]string{
				"producer": *logicalNetwork.Additional.Producer,
			},
		})
	}
	return d, nil
}
