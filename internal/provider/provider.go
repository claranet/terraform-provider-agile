package provider

import (
	"context"
	"fmt"
	"github.com/claranet/agilec-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AGILE_USERNAME", nil),
					Description: "User name for the Huawei Agile controller API. Can be specified with the `AGILE_USERNAME` " +
						"environment variable.",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AGILE_PASSWORD", nil),
					Description: "Password for the user accessing the API. Can be specified with the `AGILE_PASSWORD` " +
						"environment variable.",
				},
				"api_url": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AGILE_API", nil),
					Description: "URL of the Huawei Agile controller API. Can be specified with the `AGILE_API` environment variable. ",
				},
				"allow_insecure": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Skip verification of TLS certificates of API requests. You may need to set this to `true` " +
						"if you are using your local API without setting up a signed certificate. Can be specified with the " +
						"`AGILE_INSECURE` environment variable.",
					DefaultFunc: schema.EnvDefaultFunc("AGILE_INSECURE", false),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"agile_fabric":           dataSourceAgileFabric(),
				"agile_external_gateway": dataSourceAgileExternalGateway(),
				"agile_dhcp_group":       dataSourceAgileDhcpGroup(),
				"agile_tenant":           dataSourceAgileTenant(),
				"agile_logical_network":  dataSourceAgileLogicalNetwork(),
				"agile_logical_router":   dataSourceAgileLogicalRouter(),
				"agile_logical_switch":   dataSourceAgileLogicalSwitch(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"agile_tenant":          resourceAgileTenant(),
				"agile_logical_network": resourceAgileLogicalNetwork(),
				"agile_logical_port":    resourceAgileLogicalPort(),
				"agile_logical_router":  resourceAgileLogicalRouter(),
				"agile_logical_switch":  resourceAgileLogicalSwitch(),
				"agile_end_port":        resourceAgileEndPort(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type Config struct {
	Username   string
	Password   string
	URL        string
	IsInsecure bool
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		config := Config{
			Username:   d.Get("username").(string),
			Password:   d.Get("password").(string),
			URL:        d.Get("api_url").(string),
			IsInsecure: d.Get("allow_insecure").(bool),
		}

		if err := config.Valid(); err != nil {
			return nil, diag.FromErr(err)
		}

		return config.getClient(), nil
	}
}

func (c Config) Valid() error {

	if c.Username == "" {
		return fmt.Errorf("username must be provided for the AGILE provider")
	}

	if c.Password == "" {
		return fmt.Errorf("password must be provided for the AGILE provider")
	}

	if c.URL == "" {
		return fmt.Errorf("URL must be provided for the AGILE provider")
	}

	return nil
}

func (c Config) getClient() interface{} {
	return client.GetClient(c.URL, c.Username, c.Password, client.Insecure(c.IsInsecure))
}
