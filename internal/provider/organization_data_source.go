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
var _ tfsdk.DataSourceType = organizationDataSourceType{}
var _ tfsdk.DataSource = organizationDataSource{}

type organizationDataSourceType struct{}

func (t organizationDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A Podio organization",

		Attributes: map[string]tfsdk.Attribute{
			"url_label": {
				Type:        types.StringType,
				Description: "The URL label/slug for the organization, e.g. the `citrix` part of `https://podio.com/citrix`. Mutually exclusive with `org_id`.",
				Optional:    true,
			},
			"org_id": {
				Type:        types.Int64Type,
				Description: "The numeric ID of the organization. Mutually exclusive with `url_label`.",
				Optional:    true,
			},
			"url": {
				Description: "URL of the Podio organization",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (t organizationDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return organizationDataSource{
		provider: provider,
	}, diags
}

type organizationDataSourceData struct {
	URLLabel types.String `tfsdk:"url_label"`
	OrgID    types.Int64  `tfsdk:"org_id"`
	URL      types.String `tfsdk:"url"`
}

type organizationDataSource struct {
	provider provider
}

func (d organizationDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data organizationDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	org := &podio.Organization{}
	var err error

	// Error if both `url_label` and `org_id` are set
	if !data.URLLabel.Null && !data.OrgID.Null {
		resp.Diagnostics.AddError("Ambiguous search pattern", "Only set one of `url_label` or `org_id`, not both.")
		return
	}

	if !data.URLLabel.Null {
		org, err = d.provider.client.GetOrganizationBySlug(data.URLLabel.Value)
	} else if !data.URLLabel.Null {
		org, err = d.provider.client.GetOrganization(string(data.OrgID.Value))
	} else {
		resp.Diagnostics.AddError("No URL or Org ID specified", "Either `url_label` or `org_id` must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error fetching organization", fmt.Sprintf("Unable to fetch organization, got error: %s", err))
		return
	}

	data.URL = types.String{Value: org.URL}
	data.URLLabel = types.String{Value: org.URLLabel}
	data.OrgID = types.Int64{Value: int64(org.ID)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
