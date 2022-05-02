package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kayteh/podio-go"
	"github.com/kayteh/terraform-provider-podio/validators"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = spaceResourceType{}
var _ tfsdk.Resource = spaceResource{}

type spaceResourceType struct{}

func (t spaceResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A space/workspace within a Podio organization.",

		Attributes: map[string]tfsdk.Attribute{
			"space_id": {
				MarkdownDescription: "ID of the space",
				Type:                types.Int64Type,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"org_id": {
				MarkdownDescription: "ID of the organization",
				Type:                types.Int64Type,
				Required:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"name": {
				MarkdownDescription: "Name of the space. Changing this does not affect the ID or URL of the space.",
				Type:                types.StringType,
				Required:            true,
			},
			"url": {
				MarkdownDescription: "URL of the space",
				Type:                types.StringType,
				Computed:            true,
			},
			"privacy": {
				MarkdownDescription: "Privacy of the space, one of: `open` or `closed`. Defaults to `closed`.",
				Type:                types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []tfsdk.Validator{
					validators.StringInSliceValidator{"open", "closed"},
				},
			},
			"auto_join": {
				MarkdownDescription: "If true, new employees automatically join this space. Defaults to `false`",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"post_on_new_app": {
				MarkdownDescription: "If true, new apps are posted as a status update to this space. Defaults to `false`",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"post_on_new_member": {
				MarkdownDescription: "If true, new members are posted as a status update to this space. Defaults to `false`",
				Type:                types.BoolType,
				Optional:            true,
				Computed:            true,
			},
		},
	}, nil
}

func (t spaceResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return spaceResource{
		provider: provider,
	}, diags
}

type spaceResourceData struct {
	SpaceID         types.Int64  `tfsdk:"space_id"`
	OrgID           types.Int64  `tfsdk:"org_id"`
	Name            types.String `tfsdk:"name"`
	URL             types.String `tfsdk:"url"`
	Privacy         types.String `tfsdk:"privacy"`
	AutoJoin        types.Bool   `tfsdk:"auto_join"`
	PostOnNewApp    types.Bool   `tfsdk:"post_on_new_app"`
	PostOnNewMember types.Bool   `tfsdk:"post_on_new_member"`
}

type spaceResource struct {
	provider provider
}

func (r spaceResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data spaceResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.provider.client.CreateSpace(podio.CreateSpaceParams{
		Name:            data.Name.Value,
		OrgID:           int(data.OrgID.Value),
		Privacy:         data.Privacy.Value,
		AutoJoin:        data.AutoJoin.Value,
		PostOnNewApp:    data.PostOnNewApp.Value,
		PostOnNewMember: data.PostOnNewMember.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space, got error: %s", err))
		return
	}

	data.Name = types.String{Value: space.Name}
	data.URL = types.String{Value: space.URL}
	data.SpaceID = types.Int64{Value: int64(space.ID)}
	data.OrgID = types.Int64{Value: int64(space.OrgID)}
	data.Privacy = types.String{Value: space.Privacy}
	data.AutoJoin = types.Bool{Value: space.AutoJoin}
	data.PostOnNewApp = types.Bool{Value: space.PostOnNewApp}
	data.PostOnNewMember = types.Bool{Value: space.PostOnNewMember}

	tflog.Trace(ctx, "created a space in Podio")

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r spaceResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data spaceResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.provider.client.GetSpace(fmt.Sprintf("%d", data.SpaceID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get space, got error: %s", err))
		return
	}

	data.Name = types.String{Value: space.Name}
	data.URL = types.String{Value: space.URL}
	data.SpaceID = types.Int64{Value: int64(space.ID)}
	data.OrgID = types.Int64{Value: int64(space.OrgID)}
	data.Privacy = types.String{Value: space.Privacy}
	data.AutoJoin = types.Bool{Value: space.AutoJoin}
	data.PostOnNewApp = types.Bool{Value: space.PostOnNewApp}
	data.PostOnNewMember = types.Bool{Value: space.PostOnNewMember}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r spaceResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data spaceResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	space, err := r.provider.client.UpdateSpace(fmt.Sprintf("%d", data.SpaceID.Value), podio.CreateSpaceParams{
		Name:            data.Name.Value,
		Privacy:         data.Privacy.Value,
		AutoJoin:        data.AutoJoin.Value,
		PostOnNewApp:    data.PostOnNewApp.Value,
		PostOnNewMember: data.PostOnNewMember.Value,
	})

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get space, got error: %s", err))
		return
	}

	data.Name = types.String{Value: space.Name}
	data.URL = types.String{Value: space.URL}
	data.SpaceID = types.Int64{Value: int64(space.ID)}
	data.OrgID = types.Int64{Value: int64(space.OrgID)}
	data.Privacy = types.String{Value: space.Privacy}
	data.AutoJoin = types.Bool{Value: space.AutoJoin}
	data.PostOnNewApp = types.Bool{Value: space.PostOnNewApp}
	data.PostOnNewMember = types.Bool{Value: space.PostOnNewMember}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r spaceResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data spaceResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.client.DeleteSpace(fmt.Sprintf("%d", data.SpaceID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r spaceResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("space_id"), req, resp)
}
