package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kayteh/terraform-provider-podio/validators"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = appResourceType{}
var _ tfsdk.Resource = appResource{}

type appResourceType struct{}

func (t appResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "An app within a space in Podio",

		Attributes: map[string]tfsdk.Attribute{
			"space_id": {
				MarkdownDescription: "ID of the space",
				Type:                types.Int64Type,
				Required:            true,
			},
			"app_id": {
				MarkdownDescription: "ID of the app",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"name": {
				MarkdownDescription: "Name of the app",
				Type:                types.StringType,
				Required:            true,
			},
			"type": {
				MarkdownDescription: "Type of the app. One of: `standard`, `meeting`, `contact`",
				Type:                types.StringType,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSliceValidator{"standard", "meeting", "contact"},
				},
			},
			"item_name": {
				MarkdownDescription: "Name of the item type to use for the app",
				Type:                types.StringType,
				Optional:            true,
			},
			"description": {
				MarkdownDescription: "Description of the app",
				Type:                types.StringType,
				Optional:            true,
			},
			"usage": {
				MarkdownDescription: "How the app should be used.",
				Type:                types.StringType,
				Optional:            true,
			},
			"icon": {
				MarkdownDescription: "Icon of the app",
				Type:                types.StringType,
				Optional:            true,
			},
			"allow_edit": {
				MarkdownDescription: "Whether the app should be editable",
				Type:                types.BoolType,
				Optional:            true,
			},
			"allow_attachments": {
				MarkdownDescription: "True if attachment of files to an item is allowed",
				Type:                types.BoolType,
				Optional:            true,
			},
			"allow_comments": {
				MarkdownDescription: "True if comments are allowed",
				Type:                types.BoolType,
				Optional:            true,
			},
			"silent_creates": {
				MarkdownDescription: "True if item creates should not be posted to the stream",
				Type:                types.BoolType,
				Optional:            true,
			},
			"silent_edits": {
				MarkdownDescription: "True if item edits should not be posted to the stream",
				Type:                types.BoolType,
				Optional:            true,
			},
		},
	}, nil
}

func (t appResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return appResource{
		provider: provider,
	}, diags
}

type appResourceData struct {
	SpaceID types.Int64 `tfsdk:"space_id"`
}

type appResource struct {
	provider provider
}

func (r appResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data appResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// do stuff

	tflog.Trace(ctx, "created a space in Podio")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r appResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data appResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// do stuff

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r appResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data appResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// do stuff

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r appResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data appResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// do stuff

	resp.State.RemoveResource(ctx)
}

func (r appResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("app_id"), req, resp)
}
