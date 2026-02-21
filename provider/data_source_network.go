package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &networkDataSource{}

func NewNetworkDataSource() datasource.DataSource {
	return &networkDataSource{}
}

type networkDataSource struct {
	BaseDataSource
}

type networkDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	VlanID  types.Int64  `tfsdk:"vlan_id"`
	Subnet  types.String `tfsdk:"subnet"`
	Purpose types.String `tfsdk:"purpose"`
}

func (d *networkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (d *networkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the network. If provided, will be used for lookup.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the network. If provided, will be used for lookup.",
			},
			"vlan_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The VLAN ID of the network.",
			},
			"subnet": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The subnet of the network.",
			},
			"purpose": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The purpose of the network.",
			},
		},
	}
}

func (d *networkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data networkDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networks, err := d.Client.ListNetworks(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing networks", err.Error())
		return
	}

	found := false
	for _, n := range networks {
		if (!data.ID.IsNull() && n.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && n.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(n.ID)
			data.Name = types.StringValue(n.Name)
			if n.VLAN != nil {
				data.VlanID = types.Int64Value(int64(*n.VLAN))
			}
			data.Subnet = types.StringValue(n.IPSubnet)
			data.Purpose = types.StringValue(n.Purpose)
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("Network not found", "Could not find a network with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
