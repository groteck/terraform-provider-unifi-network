package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &networkResource{}
var _ resource.ResourceWithImportState = &networkResource{}

func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

type networkResource struct {
	BaseResource
}

type networkResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Purpose types.String `tfsdk:"purpose"`
	VlanID  types.Int64  `tfsdk:"vlan_id"`
	Subnet  types.String `tfsdk:"subnet"`
}

func (r *networkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi network (VLAN).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the network.",
			},
			"purpose": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The purpose of the network (e.g., corporate, guest). Defaults to 'corporate'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The VLAN ID for the network.",
			},
			"subnet": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The subnet for the network (CIDR format).",
			},
		},
	}
}

func (r *networkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data networkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vlanEnabled := true
	network := &unifi.Network{
		Name:    data.Name.ValueString(),
		Purpose: utils.StringOrEmpty(data.Purpose),
		NetworkVLAN: unifi.NetworkVLAN{
			VLAN:        utils.Int64Ptr(data.VlanID),
			VLANEnabled: &vlanEnabled,
			IPSubnet:    data.Subnet.ValueString(),
		},
	}

	if network.Purpose == "" {
		network.Purpose = "corporate"
	}

	created, err := r.Client.CreateNetwork(ctx, network)
	if err != nil {
		resp.Diagnostics.AddError("Error creating network", err.Error())
		return
	}

	data.ID = types.StringValue(created.ID)
	data.Purpose = types.StringValue(created.Purpose)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data networkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	network, err := r.Client.GetNetwork(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading network", err.Error())
		return
	}

	data.Name = types.StringValue(network.Name)
	data.Purpose = types.StringValue(network.Purpose)
	data.VlanID = utils.Int64Value(network.VLAN)
	data.Subnet = types.StringValue(network.IPSubnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data networkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vlanEnabled := true
	network := &unifi.Network{
		ID:      data.ID.ValueString(),
		Name:    data.Name.ValueString(),
		Purpose: data.Purpose.ValueString(),
		NetworkVLAN: unifi.NetworkVLAN{
			VLAN:        utils.Int64Ptr(data.VlanID),
			VLANEnabled: &vlanEnabled,
			IPSubnet:    data.Subnet.ValueString(),
		},
	}

	updated, err := r.Client.UpdateNetwork(ctx, data.ID.ValueString(), network)
	if err != nil {
		resp.Diagnostics.AddError("Error updating network", err.Error())
		return
	}

	data.Name = types.StringValue(updated.Name)
	data.Purpose = types.StringValue(updated.Purpose)
	data.VlanID = utils.Int64Value(updated.VLAN)
	data.Subnet = types.StringValue(updated.IPSubnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *networkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data networkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteNetwork(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting network", err.Error())
		return
	}
}

func (r *networkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
