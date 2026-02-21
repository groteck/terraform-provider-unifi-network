package provider

import (
	"context"
	client "github.com/jlopez/terraform-provider-unifi-network/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &trafficRuleResource{}
var _ resource.ResourceWithImportState = &trafficRuleResource{}

func NewTrafficRuleResource() resource.Resource {
	return &trafficRuleResource{}
}

type trafficRuleResource struct {
	BaseResource
}

type trafficRuleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	Action         types.String `tfsdk:"action"`
	MatchingTarget types.String `tfsdk:"matching_target"`
	Description    types.String `tfsdk:"description"`
}

func (r *trafficRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_traffic_rule"
}

func (r *trafficRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi traffic management rule (v2 API).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the traffic rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the traffic rule. Note: Newer UniFi versions may not store this field; it will be kept in state for convenience.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the traffic rule is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The action for the traffic rule (e.g., BLOCK, ALLOW).",
			},
			"matching_target": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The matching target for the traffic rule (e.g., INTERNET, IP, DOMAIN, APP).",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A description for the traffic rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *trafficRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data trafficRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.TrafficRule{
		Name:           data.Name.ValueString(),
		Enabled:        utils.BoolPtr(data.Enabled),
		Action:         data.Action.ValueString(),
		MatchingTarget: data.MatchingTarget.ValueString(),
		Description:    data.Description.ValueString(),
	}

	created, err := r.Client.CreateTrafficRule(ctx, rule)
	if err != nil {
		resp.Diagnostics.AddError("Error creating traffic rule", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *trafficRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data trafficRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.Client.GetTrafficRule(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading traffic rule", err.Error())
		return
	}

	r.syncState(&data, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *trafficRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data trafficRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &client.TrafficRule{
		ID:             data.ID.ValueString(),
		Name:           data.Name.ValueString(),
		Enabled:        utils.BoolPtr(data.Enabled),
		Action:         data.Action.ValueString(),
		MatchingTarget: data.MatchingTarget.ValueString(),
		Description:    data.Description.ValueString(),
	}

	updated, err := r.Client.UpdateTrafficRule(ctx, data.ID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error updating traffic rule", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *trafficRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data trafficRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteTrafficRule(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting traffic rule", err.Error())
		return
	}
}

func (r *trafficRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *trafficRuleResource) syncState(data *trafficRuleResourceModel, rule *client.TrafficRule) {
	data.ID = types.StringValue(rule.ID)

	// Handle Name: if API returns empty, keep existing if it exists
	if rule.Name != "" {
		data.Name = types.StringValue(rule.Name)
	} else if data.Name.IsNull() || data.Name.IsUnknown() {
		data.Name = types.StringNull()
	}
	// Otherwise leave data.Name as is (from plan or state)

	data.Enabled = utils.BoolValue(rule.Enabled)
	data.Action = types.StringValue(rule.Action)
	data.MatchingTarget = types.StringValue(rule.MatchingTarget)

	// Handle Description: same logic
	if rule.Description != "" {
		data.Description = types.StringValue(rule.Description)
	} else if data.Description.IsNull() || data.Description.IsUnknown() {
		data.Description = types.StringNull()
	}
}
