package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &portProfileDataSource{}

func NewPortProfileDataSource() datasource.DataSource {
	return &portProfileDataSource{}
}

type portProfileDataSource struct {
	BaseDataSource
}

type portProfileDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	NativeNetworkID  types.String `tfsdk:"native_network_id"`
	TaggedNetworkIDs types.List   `tfsdk:"tagged_network_ids"`
	Forward          types.String `tfsdk:"forward"`
}

func (d *portProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_profile"
}

func (d *portProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi port profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the port profile.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the port profile.",
			},
			"native_network_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the native network for the port profile.",
			},
			"tagged_network_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The IDs of the tagged networks for the port profile.",
			},
			"forward": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The forwarding mode for the port profile.",
			},
		},
	}
}

func (d *portProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data portProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profiles, err := d.Client.ListPortProfiles(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing port profiles", err.Error())
		return
	}

	found := false
	for _, p := range profiles {
		if (!data.ID.IsNull() && p.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && p.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(p.ID)
			data.Name = types.StringValue(p.Name)
			data.NativeNetworkID = types.StringValue(p.NativeNetworkconfID)
			data.Forward = types.StringValue(p.Forward)

			taggedIDs, _ := types.ListValueFrom(ctx, types.StringType, p.TaggedNetworkconfIDs)
			data.TaggedNetworkIDs = taggedIDs
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("Port Profile not found", "Could not find a port profile with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
