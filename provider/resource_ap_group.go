package provider

import (
	"context"
	client "github.com/jlopez/terraform-provider-unifi-network/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
)

var _ resource.Resource = &apGroupResource{}
var _ resource.ResourceWithImportState = &apGroupResource{}

func NewAPGroupResource() resource.Resource {
	return &apGroupResource{}
}

type apGroupResource struct {
	BaseResource
}

type apGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DeviceMACs  types.List   `tfsdk:"device_macs"`
	ForWLANConf types.Bool   `tfsdk:"for_wlanconf"`
}

func (r *apGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ap_group"
}

func (r *apGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi access point group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the AP group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the AP group.",
			},
			"device_macs": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The MAC addresses of the devices in the AP group.",
			},
			"for_wlanconf": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the AP group is used for WLAN configuration.",
			},
		},
	}
}

func (r *apGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data apGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceMACs := []string{}
	if !data.DeviceMACs.IsNull() && !data.DeviceMACs.IsUnknown() {
		resp.Diagnostics.Append(data.DeviceMACs.ElementsAs(ctx, &deviceMACs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	group := &client.APGroup{
		Name:        data.Name.ValueString(),
		DeviceMACs:  deviceMACs,
		ForWLANConf: utils.BoolPtr(data.ForWLANConf),
	}

	created, err := r.Client.CreateAPGroup(ctx, group)
	if err != nil {
		resp.Diagnostics.AddError("Error creating AP group", err.Error())
		return
	}

	r.syncState(ctx, &data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data apGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.Client.GetAPGroup(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading AP group", err.Error())
		return
	}

	r.syncState(ctx, &data, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data apGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceMACs := []string{}
	if !data.DeviceMACs.IsNull() && !data.DeviceMACs.IsUnknown() {
		resp.Diagnostics.Append(data.DeviceMACs.ElementsAs(ctx, &deviceMACs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	group := &client.APGroup{
		ID:          data.ID.ValueString(),
		Name:        data.Name.ValueString(),
		DeviceMACs:  deviceMACs,
		ForWLANConf: utils.BoolPtr(data.ForWLANConf),
	}

	updated, err := r.Client.UpdateAPGroup(ctx, data.ID.ValueString(), group)
	if err != nil {
		resp.Diagnostics.AddError("Error updating AP group", err.Error())
		return
	}

	r.syncState(ctx, &data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data apGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteAPGroup(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting AP group", err.Error())
		return
	}
}

func (r *apGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *apGroupResource) syncState(ctx context.Context, data *apGroupResourceModel, group *client.APGroup) {
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.ForWLANConf = utils.BoolValue(group.ForWLANConf)

	macs, _ := types.ListValueFrom(ctx, types.StringType, group.DeviceMACs)
	data.DeviceMACs = macs
}
