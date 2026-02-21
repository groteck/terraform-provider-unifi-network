package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ datasource.DataSource = &userGroupDataSource{}

func NewUserGroupDataSource() datasource.DataSource {
	return &userGroupDataSource{}
}

type userGroupDataSource struct {
	BaseDataSource
}

type userGroupDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DownloadLimit types.Int64  `tfsdk:"download_limit"`
	UploadLimit   types.Int64  `tfsdk:"upload_limit"`
}

func (d *userGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *userGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi user group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the user group.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the user group.",
			},
			"download_limit": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The download limit in Kbps.",
			},
			"upload_limit": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The upload limit in Kbps.",
			},
		},
	}
}

func (d *userGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := d.Client.ListUserGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing user groups", err.Error())
		return
	}

	found := false
	for _, g := range groups {
		if (!data.ID.IsNull() && g.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && g.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(g.ID)
			data.Name = types.StringValue(g.Name)
			data.DownloadLimit = utils.Int64Value(g.QosRateMaxDown)
			data.UploadLimit = utils.Int64Value(g.QosRateMaxUp)
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("User Group not found", "Could not find a user group with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
