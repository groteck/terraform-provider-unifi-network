package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &wlanResource{}
var _ resource.ResourceWithImportState = &wlanResource{}

func NewWLANResource() resource.Resource {
	return &wlanResource{}
}

type wlanResource struct {
	BaseResource
}

type wlanResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Passphrase  types.String `tfsdk:"passphrase"`
	Security    types.String `tfsdk:"security"`
	NetworkID   types.String `tfsdk:"network_id"`
	APGroupIDs  types.List   `tfsdk:"ap_group_ids"`
	UserGroupID types.String `tfsdk:"user_group_id"`
}

func (r *wlanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wlan"
}

func (r *wlanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi wireless network (SSID).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the WLAN.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The SSID of the wireless network.",
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the WLAN is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The passphrase for the wireless network.",
			},
			"security": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The security protocol for the wireless network (e.g., wpapsk, wpaeap).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the network configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ap_group_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The IDs of the AP groups that should broadcast this SSID.",
			},
			"user_group_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The ID of the user group for the WLAN.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *wlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data wlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apGroupIDs []string
	if !data.APGroupIDs.IsNull() && !data.APGroupIDs.IsUnknown() {
		resp.Diagnostics.Append(data.APGroupIDs.ElementsAs(ctx, &apGroupIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	wlan := &unifi.WLANConf{
		Name:          data.Name.ValueString(),
		Enabled:       utils.BoolPtr(data.Enabled),
		XPassphrase:   data.Passphrase.ValueString(),
		Security:      utils.StringOrEmpty(data.Security),
		NetworkConfID: data.NetworkID.ValueString(),
		APGroupIDs:    apGroupIDs,
		Usergroup:     data.UserGroupID.ValueString(),
	}

	if wlan.Security == "" {
		wlan.Security = "wpapsk"
	}

	created, err := r.Client.CreateWLAN(ctx, wlan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating WLAN", err.Error())
		return
	}

	r.syncState(ctx, &data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data wlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	wlan, err := r.Client.GetWLAN(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading WLAN", err.Error())
		return
	}

	r.syncState(ctx, &data, wlan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data wlanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var apGroupIDs []string
	if !data.APGroupIDs.IsNull() && !data.APGroupIDs.IsUnknown() {
		resp.Diagnostics.Append(data.APGroupIDs.ElementsAs(ctx, &apGroupIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	wlan := &unifi.WLANConf{
		ID:            data.ID.ValueString(),
		Name:          data.Name.ValueString(),
		Enabled:       utils.BoolPtr(data.Enabled),
		XPassphrase:   data.Passphrase.ValueString(),
		Security:      data.Security.ValueString(),
		NetworkConfID: data.NetworkID.ValueString(),
		APGroupIDs:    apGroupIDs,
		Usergroup:     data.UserGroupID.ValueString(),
	}

	updated, err := r.Client.UpdateWLAN(ctx, data.ID.ValueString(), wlan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating WLAN", err.Error())
		return
	}

	r.syncState(ctx, &data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *wlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data wlanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteWLAN(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting WLAN", err.Error())
		return
	}
}

func (r *wlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *wlanResource) syncState(ctx context.Context, data *wlanResourceModel, wlan *unifi.WLANConf) {
	data.ID = types.StringValue(wlan.ID)
	data.Name = types.StringValue(wlan.Name)
	data.Enabled = utils.BoolValue(wlan.Enabled)
	data.Security = types.StringValue(wlan.Security)
	data.NetworkID = types.StringValue(wlan.NetworkConfID)
	data.UserGroupID = types.StringValue(wlan.Usergroup)

	ids, _ := types.ListValueFrom(ctx, types.StringType, wlan.APGroupIDs)
	data.APGroupIDs = ids
}
