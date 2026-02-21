package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &firewallRuleResource{}
var _ resource.ResourceWithImportState = &firewallRuleResource{}

func NewFirewallRuleResource() resource.Resource {
	return &firewallRuleResource{}
}

type firewallRuleResource struct {
	BaseResource
}

type firewallRuleResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Ruleset          types.String `tfsdk:"ruleset"`
	Action           types.String `tfsdk:"action"`
	Protocol         types.String `tfsdk:"protocol"`
	SrcNetworkID     types.String `tfsdk:"src_network_id"`
	SrcNetworkType   types.String `tfsdk:"src_network_type"`
	SrcAddress       types.String `tfsdk:"src_address"`
	DstNetworkID     types.String `tfsdk:"dst_network_id"`
	DstNetworkType   types.String `tfsdk:"dst_network_type"`
	DstAddress       types.String `tfsdk:"dst_address"`
	StateEstablished types.Bool   `tfsdk:"state_established"`
	StateInvalid     types.Bool   `tfsdk:"state_invalid"`
	StateNew         types.Bool   `tfsdk:"state_new"`
	StateRelated     types.Bool   `tfsdk:"state_related"`
	IPSec            types.String `tfsdk:"ipsec"`
	RuleIndex        types.Int64  `tfsdk:"rule_index"`
	Logging          types.Bool   `tfsdk:"logging"`
}

func (r *firewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (r *firewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi firewall rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the firewall rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the firewall rule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the firewall rule is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ruleset": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ruleset for the firewall rule (e.g., LAN_IN, WAN_OUT).",
			},
			"action": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("drop"),
				MarkdownDescription: "The action for the firewall rule (e.g., accept, drop, reject). Defaults to 'drop'.",
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The protocol for the firewall rule (e.g., all, tcp, udp).",
			},
			"src_network_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Source network configuration ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_network_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Source network configuration type (e.g., ADDRv4, NETv4).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Source IP address or CIDR.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dst_network_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Destination network configuration ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dst_network_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Destination network configuration type (e.g., ADDRv4, NETv4).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dst_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Destination IP address or CIDR.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state_established": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Match established connections.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state_invalid": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Match invalid connections.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state_new": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Match new connections.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state_related": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Match related connections.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ipsec": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Match IPSec traffic (e.g., match-ipsec, match-none).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_index": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The index of the firewall rule.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"logging": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Enable logging for this rule.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *firewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data firewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &unifi.FirewallRule{
		Name:               data.Name.ValueString(),
		Enabled:            utils.BoolPtr(data.Enabled),
		Ruleset:            data.Ruleset.ValueString(),
		Action:             data.Action.ValueString(),
		Protocol:           data.Protocol.ValueString(),
		SrcNetworkConfID:   data.SrcNetworkID.ValueString(),
		SrcNetworkConfType: data.SrcNetworkType.ValueString(),
		SrcAddress:         data.SrcAddress.ValueString(),
		DstNetworkConfID:   data.DstNetworkID.ValueString(),
		DstNetworkConfType: data.DstNetworkType.ValueString(),
		DstAddress:         data.DstAddress.ValueString(),
		StateEstablished:   utils.BoolPtr(data.StateEstablished),
		StateInvalid:       utils.BoolPtr(data.StateInvalid),
		StateNew:           utils.BoolPtr(data.StateNew),
		StateRelated:       utils.BoolPtr(data.StateRelated),
		IPSec:              data.IPSec.ValueString(),
		Logging:            utils.BoolPtr(data.Logging),
		RuleIndex:          utils.Int64Ptr(data.RuleIndex),
	}

	created, err := r.Client.CreateFirewallRule(ctx, rule)
	if err != nil {
		resp.Diagnostics.AddError("Error creating firewall rule", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data firewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.Client.GetFirewallRule(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rule", err.Error())
		return
	}

	r.syncState(&data, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data firewallRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &unifi.FirewallRule{
		ID:                 data.ID.ValueString(),
		Name:               data.Name.ValueString(),
		Enabled:            utils.BoolPtr(data.Enabled),
		Ruleset:            data.Ruleset.ValueString(),
		Action:             data.Action.ValueString(),
		Protocol:           data.Protocol.ValueString(),
		SrcNetworkConfID:   data.SrcNetworkID.ValueString(),
		SrcNetworkConfType: data.SrcNetworkType.ValueString(),
		SrcAddress:         data.SrcAddress.ValueString(),
		DstNetworkConfID:   data.DstNetworkID.ValueString(),
		DstNetworkConfType: data.DstNetworkType.ValueString(),
		DstAddress:         data.DstAddress.ValueString(),
		StateEstablished:   utils.BoolPtr(data.StateEstablished),
		StateInvalid:       utils.BoolPtr(data.StateInvalid),
		StateNew:           utils.BoolPtr(data.StateNew),
		StateRelated:       utils.BoolPtr(data.StateRelated),
		IPSec:              data.IPSec.ValueString(),
		Logging:            utils.BoolPtr(data.Logging),
		RuleIndex:          utils.Int64Ptr(data.RuleIndex),
	}

	updated, err := r.Client.UpdateFirewallRule(ctx, data.ID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error updating firewall rule", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *firewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data firewallRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteFirewallRule(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting firewall rule", err.Error())
		return
	}
}

func (r *firewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *firewallRuleResource) syncState(data *firewallRuleResourceModel, rule *unifi.FirewallRule) {
	data.ID = types.StringValue(rule.ID)
	data.Name = types.StringValue(rule.Name)
	data.Enabled = utils.BoolValue(rule.Enabled)
	data.Ruleset = types.StringValue(rule.Ruleset)
	data.Action = types.StringValue(rule.Action)
	data.Protocol = types.StringValue(rule.Protocol)
	data.SrcNetworkID = types.StringValue(rule.SrcNetworkConfID)
	data.SrcNetworkType = types.StringValue(rule.SrcNetworkConfType)
	data.SrcAddress = types.StringValue(rule.SrcAddress)
	data.DstNetworkID = types.StringValue(rule.DstNetworkConfID)
	data.DstNetworkType = types.StringValue(rule.DstNetworkConfType)
	data.DstAddress = types.StringValue(rule.DstAddress)
	data.StateEstablished = utils.BoolValue(rule.StateEstablished)
	data.StateInvalid = utils.BoolValue(rule.StateInvalid)
	data.StateNew = utils.BoolValue(rule.StateNew)
	data.StateRelated = utils.BoolValue(rule.StateRelated)
	data.IPSec = types.StringValue(rule.IPSec)
	data.RuleIndex = utils.Int64Value(rule.RuleIndex)
	data.Logging = utils.BoolValue(rule.Logging)
}
