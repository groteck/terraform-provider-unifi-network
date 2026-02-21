package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ datasource.DataSource = &apGroupDataSource{}

func NewAPGroupDataSource() datasource.DataSource {
	return &apGroupDataSource{}
}

type apGroupDataSource struct {
	BaseDataSource
}

type apGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DeviceMACs  types.List   `tfsdk:"device_macs"`
	ForWLANConf types.Bool   `tfsdk:"for_wlanconf"`
}

func (d *apGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ap_group"
}

func (d *apGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi AP group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the AP group.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the AP group.",
			},
			"device_macs": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The MAC addresses of the devices in the AP group.",
			},
			"for_wlanconf": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the AP group is used for WLAN configuration.",
			},
		},
	}
}

func (d *apGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data apGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := d.Client.ListAPGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing AP groups", err.Error())
		return
	}

	found := false
	for _, g := range groups {
		if (!data.ID.IsNull() && g.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && g.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(g.ID)
			data.Name = types.StringValue(g.Name)
			data.ForWLANConf = utils.BoolValue(g.ForWLANConf)

			macs, _ := types.ListValueFrom(ctx, types.StringType, g.DeviceMACs)
			data.DeviceMACs = macs
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("AP Group not found", "Could not find an AP group with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
