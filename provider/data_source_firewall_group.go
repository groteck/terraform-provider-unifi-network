package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &firewallGroupDataSource{}

func NewFirewallGroupDataSource() datasource.DataSource {
	return &firewallGroupDataSource{}
}

type firewallGroupDataSource struct {
	BaseDataSource
}

type firewallGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	GroupType    types.String `tfsdk:"group_type"`
	GroupMembers types.List   `tfsdk:"group_members"`
}

func (d *firewallGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_group"
}

func (d *firewallGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a UniFi firewall group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the firewall group.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the firewall group.",
			},
			"group_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the firewall group.",
			},
			"group_members": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The members of the firewall group.",
			},
		},
	}
}

func (d *firewallGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data firewallGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := d.Client.ListFirewallGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing firewall groups", err.Error())
		return
	}

	found := false
	for _, g := range groups {
		if (!data.ID.IsNull() && g.ID == data.ID.ValueString()) ||
			(!data.Name.IsNull() && g.Name == data.Name.ValueString()) {
			data.ID = types.StringValue(g.ID)
			data.Name = types.StringValue(g.Name)
			data.GroupType = types.StringValue(g.GroupType)

			members, _ := types.ListValueFrom(ctx, types.StringType, g.GroupMembers)
			data.GroupMembers = members
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("Firewall Group not found", "Could not find a firewall group with the provided ID or Name")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
