package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ datasource.DataSource = &radiusProfileDataSource{}

func NewRADIUSProfileDataSource() datasource.DataSource {
	return &radiusProfileDataSource{}
}

type radiusProfileDataSource struct {
	BaseDataSource
}

type radiusProfileDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AuthServers types.List   `tfsdk:"auth_servers"`
}

func (d *radiusProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_radius_profile"
}

func (d *radiusProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi RADIUS profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the RADIUS profile.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the RADIUS profile.",
			},
			"auth_servers": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Computed: true,
						},
						"port": schema.Int64Attribute{
							Computed: true,
						},
						"secret": schema.StringAttribute{
							Computed:  true,
							Sensitive: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *radiusProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data radiusProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profiles, err := d.Client.ListRADIUSProfiles(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing RADIUS profiles", err.Error())
		return
	}

	found := false
	for _, p := range profiles {
		if (!data.ID.IsNull() && p.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && p.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(p.ID)
			data.Name = types.StringValue(p.Name)

			servers := make([]radiusServerModel, len(p.AuthServers))
			for i, s := range p.AuthServers {
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
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("RADIUS Profile not found", "Could not find a RADIUS profile with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
