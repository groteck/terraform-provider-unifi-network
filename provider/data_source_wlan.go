package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ datasource.DataSource = &wlanDataSource{}

func NewWLANDataSource() datasource.DataSource {
	return &wlanDataSource{}
}

type wlanDataSource struct {
	BaseDataSource
}

type wlanDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Security    types.String `tfsdk:"security"`
	NetworkID   types.String `tfsdk:"network_id"`
	APGroupIDs  types.List   `tfsdk:"ap_group_ids"`
	UserGroupID types.String `tfsdk:"user_group_id"`
}

func (d *wlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wlan"
}

func (d *wlanDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi wireless network (SSID).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the WLAN.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The SSID of the wireless network.",
			},
			"enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the WLAN is enabled.",
			},
			"security": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The security protocol for the wireless network.",
			},
			"network_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the network configuration.",
			},
			"ap_group_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The IDs of the AP groups that broadcast this SSID.",
			},
			"user_group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the user group for the WLAN.",
			},
		},
	}
}

func (d *wlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data wlanDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wlans, err := d.Client.ListWLANs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing WLANs", err.Error())
		return
	}

	found := false
	for _, w := range wlans {
		if (!data.ID.IsNull() && w.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && w.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(w.ID)
			data.Name = types.StringValue(w.Name)
			data.Enabled = utils.BoolValue(w.Enabled)
			data.Security = types.StringValue(w.Security)
			data.NetworkID = types.StringValue(w.NetworkConfID)
			data.UserGroupID = types.StringValue(w.Usergroup)

			ids, _ := types.ListValueFrom(ctx, types.StringType, w.APGroupIDs)
			data.APGroupIDs = ids
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("WLAN not found", "Could not find a WLAN with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
