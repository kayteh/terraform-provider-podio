package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kayteh/podio-go"
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
				Computed:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringInSliceValidator{"standard", "meeting", "contact"},
				},
			},
			"item_name": {
				MarkdownDescription: "Name of the item type to use for the app",
				Type:                types.StringType,
				Required:            true,
			},
			"description": {
				MarkdownDescription: "Description of the app",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"usage": {
				MarkdownDescription: "How the app should be used.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"icon": {
				MarkdownDescription: "Icon of the app. Must be in the format `12.png`. You might want to use `podio_icon_search` data source to pick one as the numbers are essentially useless.",
				Type:                types.StringType,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					validators.StringMatchesRegexpValidator{
						Regexp: regexp.MustCompile(`^\d+\.png$`),
					},
				},
			},
			"allow_edit": {
				MarkdownDescription: "Whether the app should be editable",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"allow_attachments": {
				MarkdownDescription: "True if attachment of files to an item is allowed",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"allow_comments": {
				MarkdownDescription: "True if comments are allowed",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"silent_creates": {
				MarkdownDescription: "True if item creates should not be posted to the stream",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"silent_edits": {
				MarkdownDescription: "True if item edits should not be posted to the stream",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
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
	SpaceID          types.Int64  `tfsdk:"space_id"`
	AppID            types.Int64  `tfsdk:"app_id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	ItemName         types.String `tfsdk:"item_name"`
	Description      types.String `tfsdk:"description"`
	Usage            types.String `tfsdk:"usage"`
	Icon             types.String `tfsdk:"icon"`
	AllowEdit        types.Bool   `tfsdk:"allow_edit"`
	AllowAttachments types.Bool   `tfsdk:"allow_attachments"`
	AllowComments    types.Bool   `tfsdk:"allow_comments"`
	SilentCreates    types.Bool   `tfsdk:"silent_creates"`
	SilentEdits      types.Bool   `tfsdk:"silent_edits"`
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

	app, err := r.provider.client.CreateApplication(
		strconv.Itoa(int(data.SpaceID.Value)),
		podio.CreateApplicationParams{
			Config: podio.AppConfig{
				Name:             data.Name.Value,
				Type:             data.Type.Value,
				ItemName:         data.ItemName.Value,
				Description:      data.Description.Value,
				Usage:            data.Usage.Value,
				Icon:             data.Icon.Value,
				AllowEdit:        data.AllowEdit.Value,
				AllowAttachments: data.AllowAttachments.Value,
				AllowComments:    data.AllowComments.Value,
				SilentCreates:    data.SilentCreates.Value,
				SilentEdits:      data.SilentEdits.Value,
			},
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create app: %s", err))
		return
	}

	data.AppID = types.Int64{Value: int64(app.AppID)}
	data.SpaceID = types.Int64{Value: int64(app.SpaceID)}
	data.Name = types.String{Value: app.Config.Name}
	data.Type = types.String{Value: app.Config.Type}
	data.ItemName = types.String{Value: app.Config.ItemName}
	data.Description = types.String{Value: app.Config.Description}
	data.Usage = types.String{Value: app.Config.Usage}
	data.Icon = types.String{Value: app.Config.Icon}
	data.AllowEdit = types.Bool{Value: app.Config.AllowEdit}
	data.AllowAttachments = types.Bool{Value: app.Config.AllowAttachments}
	data.AllowComments = types.Bool{Value: app.Config.AllowComments}
	data.SilentCreates = types.Bool{Value: app.Config.SilentCreates}
	data.SilentEdits = types.Bool{Value: app.Config.SilentEdits}

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

	app, err := r.provider.client.GetApplication(
		strconv.Itoa(int(data.AppID.Value)),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get app: %s", err))
		return
	}

	data.AppID = types.Int64{Value: int64(app.AppID)}
	data.SpaceID = types.Int64{Value: int64(app.SpaceID)}
	data.Name = types.String{Value: app.Config.Name}
	data.Type = types.String{Value: app.Config.Type}
	data.ItemName = types.String{Value: app.Config.ItemName}
	data.Description = types.String{Value: app.Config.Description}
	data.Usage = types.String{Value: app.Config.Usage}
	data.Icon = types.String{Value: app.Config.Icon}
	data.AllowEdit = types.Bool{Value: app.Config.AllowEdit}
	data.AllowAttachments = types.Bool{Value: app.Config.AllowAttachments}
	data.AllowComments = types.Bool{Value: app.Config.AllowComments}
	data.SilentCreates = types.Bool{Value: app.Config.SilentCreates}
	data.SilentEdits = types.Bool{Value: app.Config.SilentEdits}

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

	app, err := r.provider.client.UpdateApplication(
		strconv.Itoa(int(data.AppID.Value)),
		podio.CreateApplicationParams{
			Config: podio.AppConfig{
				Name:             data.Name.Value,
				Type:             data.Type.Value,
				ItemName:         data.ItemName.Value,
				Description:      data.Description.Value,
				Usage:            data.Usage.Value,
				Icon:             data.Icon.Value,
				AllowEdit:        data.AllowEdit.Value,
				AllowAttachments: data.AllowAttachments.Value,
				AllowComments:    data.AllowComments.Value,
				SilentCreates:    data.SilentCreates.Value,
				SilentEdits:      data.SilentEdits.Value,
			},
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update app: %s", err))
		return
	}

	data.AppID = types.Int64{Value: int64(app.AppID)}
	data.SpaceID = types.Int64{Value: int64(app.SpaceID)}
	data.Name = types.String{Value: app.Config.Name}
	data.Type = types.String{Value: app.Config.Type}
	data.ItemName = types.String{Value: app.Config.ItemName}
	data.Description = types.String{Value: app.Config.Description}
	data.Usage = types.String{Value: app.Config.Usage}
	data.Icon = types.String{Value: app.Config.Icon}
	data.AllowEdit = types.Bool{Value: app.Config.AllowEdit}
	data.AllowAttachments = types.Bool{Value: app.Config.AllowAttachments}
	data.AllowComments = types.Bool{Value: app.Config.AllowComments}
	data.SilentCreates = types.Bool{Value: app.Config.SilentCreates}
	data.SilentEdits = types.Bool{Value: app.Config.SilentEdits}

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

	err := r.provider.client.DeleteApplication(
		strconv.Itoa(int(data.AppID.Value)),
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete app: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r appResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("app_id"), req, resp)
}
