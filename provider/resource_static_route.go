package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &staticRouteResource{}
var _ resource.ResourceWithImportState = &staticRouteResource{}

func NewStaticRouteResource() resource.Resource {
	return &staticRouteResource{}
}

type staticRouteResource struct {
	BaseResource
}

type staticRouteResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Type     types.String `tfsdk:"type"`
	Network  types.String `tfsdk:"network"`
	Nexthop  types.String `tfsdk:"nexthop"`
	Distance types.Int64  `tfsdk:"distance"`
}

func (r *staticRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_route"
}

func (r *staticRouteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi static route.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the static route.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the static route.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the static route is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The type of the static route (e.g., static-route, interface-route).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The destination network in CIDR format.",
			},
			"nexthop": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The next hop IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"distance": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The administrative distance of the route. Must be between 1 and 255.",
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
		},
	}
}

func (r *staticRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data staticRouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route := &unifi.Routing{
		Name:                data.Name.ValueString(),
		Enabled:             utils.BoolPtr(data.Enabled),
		Type:                data.Type.ValueString(),
		StaticRouteNetwork:  data.Network.ValueString(),
		StaticRouteNexthop:  data.Nexthop.ValueString(),
		StaticRouteDistance: utils.Int64Ptr(data.Distance),
	}

	if route.Type == "" {
		route.Type = "static-route"
	}

	created, err := r.Client.CreateStaticRoute(ctx, route)
	if err != nil {
		resp.Diagnostics.AddError("Error creating static route", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data staticRouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route, err := r.Client.GetStaticRoute(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading static route", err.Error())
		return
	}

	r.syncState(&data, route)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data staticRouteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	route := &unifi.Routing{
		ID:                  data.ID.ValueString(),
		Name:                data.Name.ValueString(),
		Enabled:             utils.BoolPtr(data.Enabled),
		Type:                data.Type.ValueString(),
		StaticRouteNetwork:  data.Network.ValueString(),
		StaticRouteNexthop:  data.Nexthop.ValueString(),
		StaticRouteDistance: utils.Int64Ptr(data.Distance),
	}

	updated, err := r.Client.UpdateStaticRoute(ctx, data.ID.ValueString(), route)
	if err != nil {
		resp.Diagnostics.AddError("Error updating static route", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data staticRouteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteStaticRoute(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting static route", err.Error())
		return
	}
}

func (r *staticRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *staticRouteResource) syncState(data *staticRouteResourceModel, route *unifi.Routing) {
	data.ID = types.StringValue(route.ID)
	data.Name = types.StringValue(route.Name)
	data.Enabled = utils.BoolValue(route.Enabled)
	data.Type = types.StringValue(route.Type)
	data.Network = types.StringValue(route.StaticRouteNetwork)
	data.Nexthop = types.StringValue(route.StaticRouteNexthop)
	data.Distance = utils.Int64Value(route.StaticRouteDistance)
}
