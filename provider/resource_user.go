package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &userResource{}
var _ resource.ResourceWithImportState = &userResource{}

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	BaseResource
}

type userResourceModel struct {
	ID          types.String `tfsdk:"id"`
	MAC         types.String `tfsdk:"mac"`
	Name        types.String `tfsdk:"name"`
	Note        types.String `tfsdk:"note"`
	UseFixedIP  types.Bool   `tfsdk:"use_fixedip"`
	FixedIP     types.String `tfsdk:"fixed_ip"`
	NetworkID   types.String `tfsdk:"network_id"`
	UserGroupID types.String `tfsdk:"user_group_id"`
	Blocked     types.Bool   `tfsdk:"blocked"`
	IsWired     types.Bool   `tfsdk:"is_wired"`
	IsGuest     types.Bool   `tfsdk:"is_guest"`
	OUI         types.String `tfsdk:"oui"`
	Noted       types.Bool   `tfsdk:"noted"`
	SiteID      types.String `tfsdk:"site_id"`
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi user (client device record, DHCP reservation).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the user record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The MAC address of the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"note": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A note for the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_fixedip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to use a fixed IP address for the device.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fixed_ip": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The fixed IP address for the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the network for the fixed IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_group_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the user group for the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the device is blocked.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_wired": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the device is wired.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_guest": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the device is a guest.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"oui": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Organizationally Unique Identifier of the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"noted": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the device has a note.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &client.User{
		MAC:         data.MAC.ValueString(),
		Name:        data.Name.ValueString(),
		Note:        data.Note.ValueString(),
		UseFixedIP:  utils.BoolPtr(data.UseFixedIP),
		FixedIP:     data.FixedIP.ValueString(),
		NetworkID:   data.NetworkID.ValueString(),
		UsergroupID: data.UserGroupID.ValueString(),
		Blocked:     utils.BoolPtr(data.Blocked),
	}

	created, err := r.Client.CreateUser(ctx, user)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.Client.GetUser(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	r.syncState(&data, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &client.User{
		ID:          data.ID.ValueString(),
		MAC:         data.MAC.ValueString(),
		Name:        data.Name.ValueString(),
		Note:        data.Note.ValueString(),
		UseFixedIP:  utils.BoolPtr(data.UseFixedIP),
		FixedIP:     data.FixedIP.ValueString(),
		NetworkID:   data.NetworkID.ValueString(),
		UsergroupID: data.UserGroupID.ValueString(),
		Blocked:     utils.BoolPtr(data.Blocked),
	}

	updated, err := r.Client.UpdateUser(ctx, data.ID.ValueString(), user)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteUser(ctx, data.MAC.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userResource) syncState(data *userResourceModel, user *client.User) {
	data.ID = types.StringValue(user.ID)
	data.MAC = types.StringValue(user.MAC)
	data.Name = utils.StringToValue(user.Name)
	data.Note = utils.StringToValue(user.Note)
	data.UseFixedIP = utils.BoolValue(user.UseFixedIP)
	data.FixedIP = utils.StringToValue(user.FixedIP)
	data.NetworkID = utils.StringToValue(user.NetworkID)
	data.UserGroupID = utils.StringToValue(user.UsergroupID)
	data.Blocked = utils.BoolValue(user.Blocked)
	data.IsWired = utils.BoolValue(user.IsWired)
	data.IsGuest = utils.BoolValue(user.IsGuest)
	data.OUI = utils.StringToValue(user.OUI)
	data.Noted = utils.BoolValue(user.Noted)
	data.SiteID = types.StringValue(user.SiteID)
}
