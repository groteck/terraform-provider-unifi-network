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
)

var _ resource.Resource = &firewallGroupResource{}
var _ resource.ResourceWithImportState = &firewallGroupResource{}

func NewFirewallGroupResource() resource.Resource {
	return &firewallGroupResource{}
}

type firewallGroupResource struct {
	BaseResource
}

type firewallGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	GroupType    types.String `tfsdk:"group_type"`
	GroupMembers types.List   `tfsdk:"group_members"`
}

func (r *firewallGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_group"
}

func (r *firewallGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi firewall group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the firewall group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the firewall group.",
			},
			"group_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the firewall group (e.g., address-group, port-group, ipv6-address-group).",
			},
			"group_members": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The members of the firewall group.",
			},
		},
	}
}

func (r *firewallGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data firewallGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var members []string
	resp.Diagnostics.Append(data.GroupMembers.ElementsAs(ctx, &members, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &client.FirewallGroup{
		Name:         data.Name.ValueString(),
		GroupType:    data.GroupType.ValueString(),
		GroupMembers: members,
	}

	created, err := r.Client.CreateFirewallGroup(ctx, group)
	if err != nil {
		resp.Diagnostics.AddError("Error creating firewall group", err.Error())
		return
	}

	r.syncState(ctx, &data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data firewallGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.Client.GetFirewallGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall group", err.Error())
		return
	}

	r.syncState(ctx, &data, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data firewallGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var members []string
	resp.Diagnostics.Append(data.GroupMembers.ElementsAs(ctx, &members, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group := &client.FirewallGroup{
		ID:           data.ID.ValueString(),
		Name:         data.Name.ValueString(),
		GroupType:    data.GroupType.ValueString(),
		GroupMembers: members,
	}

	updated, err := r.Client.UpdateFirewallGroup(ctx, data.ID.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError("Error updating firewall group", err.Error())
		return
	}

	r.syncState(ctx, &data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data firewallGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteFirewallGroup(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall group", err.Error())
		return
	}
}

func (r *firewallGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *firewallGroupResource) syncState(ctx context.Context, data *firewallGroupResourceModel, group *client.FirewallGroup) {
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.GroupType = types.StringValue(group.GroupType)

	members, _ := types.ListValueFrom(ctx, types.StringType, group.GroupMembers)
	data.GroupMembers = members
}
