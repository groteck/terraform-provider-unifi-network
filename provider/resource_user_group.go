package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &userGroupResource{}
var _ resource.ResourceWithImportState = &userGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &userGroupResource{}
}

type userGroupResource struct {
	BaseResource
}

type userGroupResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DownloadLimit types.Int64  `tfsdk:"download_limit"`
	UploadLimit   types.Int64  `tfsdk:"upload_limit"`
}

func (r *userGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *userGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi user group (bandwidth profile).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the user group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the user group.",
			},
			"download_limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The download limit in Kbps.",
			},
			"upload_limit": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The upload limit in Kbps.",
			},
		},
	}
}

func (r *userGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &client.UserGroup{
		Name:           data.Name.ValueString(),
		QosRateMaxDown: utils.Int64Ptr(data.DownloadLimit),
		QosRateMaxUp:   utils.Int64Ptr(data.UploadLimit),
	}

	created, err := r.Client.CreateUserGroup(ctx, group)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user group", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.Client.GetUserGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user group", err.Error())
		return
	}

	r.syncState(&data, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &client.UserGroup{
		ID:             data.ID.ValueString(),
		Name:           data.Name.ValueString(),
		QosRateMaxDown: utils.Int64Ptr(data.DownloadLimit),
		QosRateMaxUp:   utils.Int64Ptr(data.UploadLimit),
	}

	updated, err := r.Client.UpdateUserGroup(ctx, data.ID.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user group", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteUserGroup(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting user group", err.Error())
		return
	}
}

func (r *userGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userGroupResource) syncState(data *userGroupResourceModel, group *client.UserGroup) {
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.DownloadLimit = utils.Int64Value(group.QosRateMaxDown)
	data.UploadLimit = utils.Int64Value(group.QosRateMaxUp)
}
