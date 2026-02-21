package provider

import (
	client "github.com/jlopez/terraform-provider-unifi-network/internal/client"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &radiusProfileResource{}
var _ resource.ResourceWithImportState = &radiusProfileResource{}

func NewRADIUSProfileResource() resource.Resource {
	return &radiusProfileResource{}
}

type radiusProfileResource struct {
	BaseResource
}

type radiusProfileResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AuthServers types.List   `tfsdk:"auth_servers"`
}

type radiusServerModel struct {
	IP     types.String `tfsdk:"ip"`
	Port   types.Int64  `tfsdk:"port"`
	Secret types.String `tfsdk:"secret"`
}

func (r *radiusProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_profile"
}

func (r *radiusProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi RADIUS profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the RADIUS profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the RADIUS profile.",
			},
			"auth_servers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required: true,
						},
						"port": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"secret": schema.StringAttribute{
							Required:  true,
							Sensitive: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func (r *radiusProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data radiusProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var authServers []radiusServerModel
	resp.Diagnostics.Append(data.AuthServers.ElementsAs(ctx, &authServers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	servers := make([]client.RADIUSServer, len(authServers))
	for i, s := range authServers {
		port := int(s.Port.ValueInt64())
		if port == 0 {
			port = 1812
		}
		servers[i] = client.RADIUSServer{
			IP:      s.IP.ValueString(),
			Port:    &port,
			XSecret: s.Secret.ValueString(),
		}
	}

	profile := &client.RADIUSProfile{
		Name:        data.Name.ValueString(),
		AuthServers: servers,
	}

	created, err := r.Client.CreateRADIUSProfile(ctx, profile)
	if err != nil {
		resp.Diagnostics.AddError("Error creating RADIUS profile", err.Error())
		return
	}

	r.syncState(ctx, &data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data radiusProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.Client.GetRADIUSProfile(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading RADIUS profile", err.Error())
		return
	}

	r.syncState(ctx, &data, profile)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data radiusProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var authServers []radiusServerModel
	resp.Diagnostics.Append(data.AuthServers.ElementsAs(ctx, &authServers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	servers := make([]client.RADIUSServer, len(authServers))
	for i, s := range authServers {
		port := int(s.Port.ValueInt64())
		servers[i] = client.RADIUSServer{
			IP:      s.IP.ValueString(),
			Port:    &port,
			XSecret: s.Secret.ValueString(),
		}
	}

	profile := &client.RADIUSProfile{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		AuthServers: servers,
	}

	updated, err := r.Client.UpdateRADIUSProfile(ctx, data.ID.ValueString(), profile)
	if err != nil {
		resp.Diagnostics.AddError("Error updating RADIUS profile", err.Error())
		return
	}

	r.syncState(ctx, &data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *radiusProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data radiusProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteRADIUSProfile(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting RADIUS profile", err.Error())
		return
	}
}

func (r *radiusProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *radiusProfileResource) syncState(ctx context.Context, data *radiusProfileResourceModel, profile *client.RADIUSProfile) {
	data.ID = types.StringValue(profile.ID)
	data.Name = types.StringValue(profile.Name)

	servers := make([]radiusServerModel, len(profile.AuthServers))
	for i, s := range profile.AuthServers {
		servers[i] = radiusServerModel{
			IP:     types.StringValue(s.IP),
			Port:   utils.Int64Value(s.Port),
			Secret: types.StringValue(s.XSecret),
		}
	}

	newServers, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"ip":     types.StringType,
		"port":   types.Int64Type,
		"secret": types.StringType,
	}}, servers)
	data.AuthServers = newServers
}
