package provider

import (
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &portForwardResource{}
var _ resource.ResourceWithImportState = &portForwardResource{}

func NewPortForwardResource() resource.Resource {
	return &portForwardResource{}
}

type portForwardResource struct {
	BaseResource
}

type portForwardResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Protocol      types.String `tfsdk:"protocol"`
	Src           types.String `tfsdk:"src"`
	DstPort       types.String `tfsdk:"dst_port"`
	Fwd           types.String `tfsdk:"fwd"`
	FwdPort       types.String `tfsdk:"fwd_port"`
	PfwdInterface types.String `tfsdk:"pfwd_interface"`
}

func (r *portForwardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forward"
}

func (r *portForwardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi port forwarding rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the port forwarding rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the port forwarding rule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the port forwarding rule is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The protocol for the port forwarding rule (e.g., tcp, udp, tcp_udp).",
			},
			"src": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The source IP or network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dst_port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The destination port or port range.",
			},
			"fwd": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The forward-to IP address.",
			},
			"fwd_port": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The forward-to port or port range.",
			},
			"pfwd_interface": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The interface for the port forwarding rule (e.g., wan, wan2, both).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *portForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data portForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	forward := &client.PortForward{
		Name:          data.Name.ValueString(),
		Enabled:       utils.BoolPtr(data.Enabled),
		Proto:         data.Protocol.ValueString(),
		Src:           data.Src.ValueString(),
		DstPort:       data.DstPort.ValueString(),
		Fwd:           data.Fwd.ValueString(),
		FwdPort:       data.FwdPort.ValueString(),
		PfwdInterface: data.PfwdInterface.ValueString(),
	}

	if forward.PfwdInterface == "" {
		forward.PfwdInterface = "wan"
	}
	if forward.Src == "" {
		forward.Src = "any"
	}

	created, err := r.Client.CreatePortForward(ctx, forward)
	if err != nil {
		resp.Diagnostics.AddError("Error creating port forward", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data portForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	forward, err := r.Client.GetPortForward(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading port forward", err.Error())
		return
	}

	r.syncState(&data, forward)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data portForwardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	forward := &client.PortForward{
		ID:            data.ID.ValueString(),
		Name:          data.Name.ValueString(),
		Enabled:       utils.BoolPtr(data.Enabled),
		Proto:         data.Protocol.ValueString(),
		Src:           data.Src.ValueString(),
		DstPort:       data.DstPort.ValueString(),
		Fwd:           data.Fwd.ValueString(),
		FwdPort:       data.FwdPort.ValueString(),
		PfwdInterface: data.PfwdInterface.ValueString(),
	}

	updated, err := r.Client.UpdatePortForward(ctx, data.ID.ValueString(), forward)
	if err != nil {
		resp.Diagnostics.AddError("Error updating port forward", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *portForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data portForwardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeletePortForward(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting port forward", err.Error())
		return
	}
}

func (r *portForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *portForwardResource) syncState(data *portForwardResourceModel, forward *client.PortForward) {
	data.ID = types.StringValue(forward.ID)
	data.Name = types.StringValue(forward.Name)
	data.Enabled = utils.BoolValue(forward.Enabled)
	data.Protocol = types.StringValue(forward.Proto)
	data.Src = types.StringValue(forward.Src)
	data.DstPort = types.StringValue(forward.DstPort)
	data.Fwd = types.StringValue(forward.Fwd)
	data.FwdPort = types.StringValue(forward.FwdPort)
	data.PfwdInterface = types.StringValue(forward.PfwdInterface)
}
