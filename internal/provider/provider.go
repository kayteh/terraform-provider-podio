package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kayteh/podio-go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.Provider = &provider{}

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	client *podio.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	TrustLevel   types.Int64  `tfsdk:"trust_level"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ClientID.Null || data.ClientSecret.Null || data.Username.Null || data.Password.Null {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"While creating the provider, some configuration values were missing. Please check the documentation for required configuration values.",
		)
		return
	}

	p.client = podio.NewClient(podio.ClientOptions{
		ApiKey:    data.ClientID.Value,
		ApiSecret: data.ClientSecret.Value,
		UserAgent: fmt.Sprintf("terraform-provider-podio/%s", p.version),
	})

	err := p.client.AuthenticateWithCredentials(data.Username.Value, data.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to authenticate with Podio",
			fmt.Sprintf("Failed to authenticate with Podio: %s", err),
		)
		return
	}

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"podio_space": spaceResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"podio_organization": organizationDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"client_id": {
				MarkdownDescription: "Client ID for Podio",
				Type:                types.StringType,
				Sensitive:           true,
				Required:            true,
			},
			"client_secret": {
				MarkdownDescription: "Client Secret for Podio",
				Type:                types.StringType,
				Sensitive:           true,
				Required:            true,
			},
			"username": {
				MarkdownDescription: "Username for Podio",
				Type:                types.StringType,
				Sensitive:           true,
				Optional:            true,
				Computed:            true,
			},
			"password": {
				MarkdownDescription: "Password for Podio",
				Type:                types.StringType,
				Sensitive:           true,
				Required:            true,
			},
			"trust_level": {
				MarkdownDescription: "rust level for Podio API key. Turns on guard-rails when your key can't be used for certain operations for lesser trust levels. `2` is the default, allowing all public API methods.",
				Type:                types.Int64Type,
				Optional:            true,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
