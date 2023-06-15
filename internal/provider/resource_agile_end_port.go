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
)

func resourceAgileEndPort() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages End Ports.",
		CreateContext: resourceAgileEndPortCreate,
		ReadContext:   resourceAgileEndPortRead,
		UpdateContext: resourceAgileEndPortUpdate,
		DeleteContext: resourceAgileEndPortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAgileEndPortImport,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "End port ID.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "End port name.",
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
				Description: "End port description.",
				Optional:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"logic_port_id": {
				Type:         schema.TypeString,
				Description:  "ID of the associated logical port",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"logic_network_id": {
				Type:         schema.TypeString,
				Description:  "ID of the logical network to which the end port belongs.",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"location": {
				Type:        schema.TypeString,
				Description: "Location.",
				Optional:    true,
				ForceNew:    false,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 242),
				),
			},
			"vm_name": {
				Type:        schema.TypeString,
				Description: "Terminal name.",
				Optional:    true,
				ForceNew:    false,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringLenBetween(0, 255),
				),
			},
			"ipv4": {
				Type:        schema.TypeString,
				Description: "End Port IPv4.",
				Optional:    true,
				ForceNew:    false,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.IsIPv4Address,
				),
			},
			"ipv6": {
				Type:        schema.TypeString,
				Description: "End Port IPv6.",
				Optional:    true,
				ForceNew:    false,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.IsIPv6Address,
				),
			},
		},
	}
}

func resourceAgileEndPortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] End Port: Beginning Creation")

	agileClient := meta.(*agile.Client)

	id, _ := uuid.NewV4()

	name := d.Get("name").(string)

	endPort, err := NewEndPortAttributes(d)

	if err != nil {
		return err
	}

	if err := agileClient.CreateEndPort(agile.String(id.String()), agile.String(name), endPort); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.String())
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAgileEndPortRead(ctx, d, meta)
}

func resourceAgileEndPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read End Port", d.Id())
	agileClient := meta.(*agile.Client)
	id := d.Id()
	endPort, err := agileClient.GetEndPort(id)

	if err != nil {
		d.SetId("")
		return nil
	}

	if _, err := setEndPortAttributes(endPort, d); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAgileEndPortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: End Port: Beginning Update", d.Id())
	agileClient := meta.(*agile.Client)

	name := d.Get("name").(string)

	endPortAttr, errAttr := NewEndPortAttributes(d)

	if errAttr != nil {
		return errAttr
	}

	if _, err := agileClient.UpdateEndPort(agile.String(d.Id()), agile.String(name), endPortAttr); err != nil {
		return diag.FromErr(err)
	}

	return resourceAgileEndPortRead(ctx, d, meta)
}

func resourceAgileEndPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	agileClient := meta.(*agile.Client)

	err := agileClient.DeleteEndPort(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")

	return diag.FromErr(nil)
}

func NewEndPortAttributes(d *schema.ResourceData) (*models.EndPortAttributes, diag.Diagnostics) {
	endPortAttr := models.EndPortAttributes{}

	if _, ok := d.GetOk("description"); ok {
		endPortAttr.Description = agile.String(d.Get("description").(string))
	}

	if _, ok := d.GetOk("logic_port_id"); ok {
		endPortAttr.LogicPortId = agile.String(d.Get("logic_port_id").(string))
	}

	if _, ok := d.GetOk("logic_network_id"); ok {
		endPortAttr.LogicNetworkId = agile.String(d.Get("logic_network_id").(string))
	}

	if _, ok := d.GetOk("location"); ok {
		endPortAttr.Location = agile.String(d.Get("location").(string))
	}

	if _, ok := d.GetOk("vm_name"); ok {
		endPortAttr.VmName = agile.String(d.Get("vm_name").(string))
	}

	if _, ok := d.GetOk("ipv4"); ok {
		endPortAttr.Ipv4 = []*string{agile.String(d.Get("ipv4").(string))}
	}

	if _, ok := d.GetOk("ipv6"); ok {
		endPortAttr.Ipv6 = []*string{agile.String(d.Get("ipv6").(string))}
	}

	return &endPortAttr, nil

}

func resourceAgileEndPortImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	agileClient := meta.(*agile.Client)

	id := d.Id()
	endPort, err := agileClient.GetEndPort(id)

	if err != nil {
		return nil, err
	}

	schemaFilled, err := setEndPortAttributes(endPort, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func setEndPortAttributes(endPort *models.EndPort, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("name", *endPort.Name)
	d.Set("description", *endPort.Description)
	d.Set("logic_port_id", *endPort.LogicPortId)
	d.Set("logic_network_id", endPort.LogicNetworkId)
	d.Set("location", *endPort.Location)
	d.Set("vm_name", *endPort.VmName)
	d.Set("ipv4", *endPort.Ipv4[0])
	d.Set("ipv6", *endPort.Ipv6[0])
	return d, nil
}
