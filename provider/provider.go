package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jlopez/terraform-provider-unifi-network/internal/client"
)

var _ provider.Provider = &unifiProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &unifiProvider{
			version: version,
		}
	}
}

type unifiProvider struct {
	version string
}

type unifiProviderModel struct {
	Host         types.String `tfsdk:"host"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	APIKey       types.String `tfsdk:"api_key"`
	Site         types.String `tfsdk:"site"`
	Insecure     types.Bool   `tfsdk:"allow_insecure"`
	IsStandalone types.Bool   `tfsdk:"is_standalone"`
}

func (p *unifiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "unifi"
	resp.Version = p.version
}

func (p *unifiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for UniFi Network API.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The UniFi controller host URL. Defaults to https://localhost:8443.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "UniFi controller username. Can also be set via UNIFI_USERNAME environment variable.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "UniFi controller password. Can also be set via UNIFI_PASSWORD environment variable.",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "UniFi Integration API Key. Can also be set via UNIFI_API_KEY environment variable.",
			},
			"site": schema.StringAttribute{
				Optional:    true,
				Description: "UniFi site ID. Defaults to 'default'.",
			},
			"allow_insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow insecure SSL connections. Defaults to false.",
			},
			"is_standalone": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true if using a standalone UniFi Network Application (no /proxy/network prefix). Defaults to false.",
			},
		},
	}
}

func (p *unifiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data unifiProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration defaults and env vars
	host := os.Getenv("UNIFI_HOST")
	if host == "" {
		host = "https://localhost:8443"
	}
	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	username := os.Getenv("UNIFI_USERNAME")
	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	password := os.Getenv("UNIFI_PASSWORD")
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	apiKey := os.Getenv("UNIFI_API_KEY")
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	site := "default"
	if !data.Site.IsNull() {
		site = data.Site.ValueString()
	}

	insecure := false
	if !data.Insecure.IsNull() {
		insecure = data.Insecure.ValueBool()
	}

	isStandalone := false
	if !data.IsStandalone.IsNull() {
		isStandalone = data.IsStandalone.ValueBool()
	}

	c, err := client.NewClient(host, username, password, apiKey, site, insecure, isStandalone)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Failed to create unifi client: "+err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *unifiProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworkResource,
		NewFirewallRuleResource,
		NewPortProfileResource,
		NewUserGroupResource,
		NewAPGroupResource,
		NewWLANResource,
		NewFirewallGroupResource,
		NewUserResource,
		NewRADIUSProfileResource,
		NewPortForwardResource,
		NewStaticRouteResource,
		NewStaticDNSResource,
		NewTrafficRuleResource,
	}
}

func (p *unifiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNetworkDataSource,
		NewAPGroupDataSource,
		NewUserGroupDataSource,
		NewWLANDataSource,
		NewFirewallGroupDataSource,
		NewRADIUSProfileDataSource,
		NewPortProfileDataSource,
	}
}
