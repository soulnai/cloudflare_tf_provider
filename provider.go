package main

import (
	"context"
	"terraform-provider-cloudflare-tunnel/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &cloudflareTunnelProvider{}
)

// New returns a new provider factory.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cloudflareTunnelProvider{
			version: version,
		}
	}
}

type cloudflareTunnelProvider struct {
	version string
}

type cloudflareTunnelProviderModel struct {
	ApiToken  types.String `tfsdk:"api_token"`
	AccountId types.String `tfsdk:"account_id"`
	BaseURL   types.String `tfsdk:"base_url"`
}

func (p *cloudflareTunnelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudflare-tunnel"
	resp.Version = p.version
}

func (p *cloudflareTunnelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"account_id": schema.StringAttribute{
				Required: true,
			},
			"base_url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *cloudflareTunnelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data cloudflareTunnelProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if configuration is unknown (e.g. during terraform plan)
	if data.ApiToken.IsUnknown() || data.AccountId.IsUnknown() {
		return
	}

	// Create the client
	c := client.NewClient(data.ApiToken.ValueString(), data.AccountId.ValueString())

	// If base_url is set, update it (you might need to update NewClient or set it directly if you export the field)
	if !data.BaseURL.IsNull() {
		c.BaseURL = data.BaseURL.ValueString()
	}

	// Make the client available to resources/data sources
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *cloudflareTunnelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		CloudflareTunnelResource,
	}
}

func (p *cloudflareTunnelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		CloudflareTunnelDataSource,
	}
}
