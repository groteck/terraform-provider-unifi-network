package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/provider/utils"
	"github.com/resnickio/unifi-go-sdk/pkg/unifi"
)

var _ resource.Resource = &staticDNSResource{}
var _ resource.ResourceWithImportState = &staticDNSResource{}

func NewStaticDNSResource() resource.Resource {
	return &staticDNSResource{}
}

type staticDNSResource struct {
	BaseResource
}

type staticDNSResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Key        types.String `tfsdk:"key"`
	Value      types.String `tfsdk:"value"`
	RecordType types.String `tfsdk:"record_type"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	TTL        types.Int64  `tfsdk:"ttl"`
}

func (r *staticDNSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static_dns"
}

func (r *staticDNSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a UniFi static DNS record.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the DNS record.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The hostname for the DNS record.",
			},
			"value": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The value (IP or hostname) for the DNS record.",
			},
			"record_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The type of the DNS record (e.g., A, CNAME). Defaults to 'A'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the DNS record is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ttl": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The TTL for the DNS record.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *staticDNSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data staticDNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record := &unifi.StaticDNS{
		Key:        data.Key.ValueString(),
		Value:      data.Value.ValueString(),
		RecordType: utils.StringOrEmpty(data.RecordType),
		Enabled:    utils.BoolPtr(data.Enabled),
		TTL:        utils.Int64Ptr(data.TTL),
	}

	if record.RecordType == "" {
		record.RecordType = "A"
	}

	created, err := r.Client.CreateStaticDNS(ctx, record)
	if err != nil {
		resp.Diagnostics.AddError("Error creating static DNS", err.Error())
		return
	}

	r.syncState(&data, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticDNSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data staticDNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := r.Client.GetStaticDNS(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading static DNS", err.Error())
		return
	}

	r.syncState(&data, record)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticDNSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data staticDNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record := &unifi.StaticDNS{
		ID:         data.ID.ValueString(),
		Key:        data.Key.ValueString(),
		Value:      data.Value.ValueString(),
		RecordType: data.RecordType.ValueString(),
		Enabled:    utils.BoolPtr(data.Enabled),
		TTL:        utils.Int64Ptr(data.TTL),
	}

	updated, err := r.Client.UpdateStaticDNS(ctx, data.ID.ValueString(), record)
	if err != nil {
		resp.Diagnostics.AddError("Error updating static DNS", err.Error())
		return
	}

	r.syncState(&data, updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *staticDNSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data staticDNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.Client.DeleteStaticDNS(ctx, data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting static DNS", err.Error())
		return
	}
}

func (r *staticDNSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *staticDNSResource) syncState(data *staticDNSResourceModel, record *unifi.StaticDNS) {
	data.ID = types.StringValue(record.ID)
	data.Key = types.StringValue(record.Key)
	data.Value = types.StringValue(record.Value)
	data.RecordType = types.StringValue(record.RecordType)
	data.Enabled = utils.BoolValue(record.Enabled)
	data.TTL = utils.Int64Value(record.TTL)
}
