package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kayteh/podio-go"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"podio_organization": dataSourceOrganization(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"podio_space": resourceSpace(),
			},
			Schema: map[string]*schema.Schema{
				"client_id": {
					Description: "Client ID for Podio",
					Required:    true,
					Type:        schema.TypeString,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("PODIO_CLIENT_ID", ""),
				},
				"client_secret": {
					Description: "Client Secret for Podio",
					Required:    true,
					Type:        schema.TypeString,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("PODIO_CLIENT_SECRET", ""),
				},
				"username": {
					Description: "Username for Podio (Required for now)",
					Required:    true,
					Type:        schema.TypeString,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("PODIO_USERNAME", ""),
				},
				"password": {
					Description: "Password for Podio",
					Type:        schema.TypeString,
					Sensitive:   true,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PODIO_PASSWORD", ""),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client := podio.NewClient(podio.ClientOptions{
			ApiKey:    d.Get("client_id").(string),
			ApiSecret: d.Get("client_secret").(string),
			UserAgent: p.UserAgent("terraform-provider-podio", version),
		})

		if username, ok := d.GetOk("username"); ok {
			err := client.AuthenticateWithCredentials(username.(string), d.Get("password").(string))
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		return client, nil
	}
}
