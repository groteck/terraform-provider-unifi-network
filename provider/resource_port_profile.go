package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &portProfileResource{}
var _ resource.ResourceWithImportState = &portProfileResource{}

func NewPortProfileResource() resource.Resource {
	return &portProfileResource{}
}

type portProfileResource struct {
	BaseResource
}

type portProfileResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	NativeNetworkID  types.String `tfsdk:"native_network_id"`
	TaggedNetworkIDs types.List   `tfsdk:"tagged_network_ids"`
	Forward          types.String `tfsdk:"forward"`
}

func (r *portProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_profile"
}

func (r *portProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi port profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the port profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the port profile.",
			},
			"native_network_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the native network for the port profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tagged_network_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The IDs of the tagged networks for the port profile.",
			},
			"forward": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The forwarding mode for the port profile (e.g., all, native, customize).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *portProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data portProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taggedNetworkIDs []string
	if !data.TaggedNetworkIDs.IsNull() && !data.TaggedNetworkIDs.IsUnknown() {
		resp.Diagnostics.Append(data.TaggedNetworkIDs.ElementsAs(ctx, &taggedNetworkIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	forward := "native"
	if !data.Forward.IsNull() {
		forward = data.Forward.ValueString()
	}

	profile := &unifi.PortConf{
		Name:                 data.Name.ValueString(),
		NativeNetworkconfID:  data.NativeNetworkID.ValueString(),
		TaggedNetworkconfIDs: taggedNetworkIDs,
		Forward:              forward,
	}

	created, err := r.Client.CreatePortProfile(ctx, profile)
	if err != nil {
		resp.Diagnostics.AddError("Error creating port profile", err.Error())
		return
	}

	r.syncState(ctx, &data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data portProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.Client.GetPortProfile(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading port profile", err.Error())
		return
	}

	r.syncState(ctx, &data, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data portProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var taggedNetworkIDs []string
	if !data.TaggedNetworkIDs.IsNull() && !data.TaggedNetworkIDs.IsUnknown() {
		resp.Diagnostics.Append(data.TaggedNetworkIDs.ElementsAs(ctx, &taggedNetworkIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	profile := &unifi.PortConf{
		ID:                   data.ID.ValueString(),
		Name:                 data.Name.ValueString(),
		NativeNetworkconfID:  data.NativeNetworkID.ValueString(),
		TaggedNetworkconfIDs: taggedNetworkIDs,
		Forward:              data.Forward.ValueString(),
	}

	updated, err := r.Client.UpdatePortProfile(ctx, data.ID.ValueString(), profile)
	if err != nil {
		resp.Diagnostics.AddError("Error updating port profile", err.Error())
		return
	}

	r.syncState(ctx, &data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data portProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeletePortProfile(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting port profile", err.Error())
		return
	}
}

func (r *portProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *portProfileResource) syncState(ctx context.Context, data *portProfileResourceModel, profile *unifi.PortConf) {
	data.ID = types.StringValue(profile.ID)
	data.Name = types.StringValue(profile.Name)
	data.NativeNetworkID = types.StringValue(profile.NativeNetworkconfID)
	data.Forward = types.StringValue(profile.Forward)

	taggedIDs, _ := types.ListValueFrom(ctx, types.StringType, profile.TaggedNetworkconfIDs)
	data.TaggedNetworkIDs = taggedIDs
}
