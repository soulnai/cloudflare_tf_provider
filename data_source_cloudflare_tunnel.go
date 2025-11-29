package main

import (
	"context"
	"fmt"
	"terraform-provider-cloudflare-tunnel/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &cloudflareTunnelDataSource{}
)

type cloudflareTunnelDataSource struct {
	client *client.Client
}

type cloudflareTunnelDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	AccountId   types.String `tfsdk:"account_id"`
	TunnelToken types.String `tfsdk:"tunnel_token"`
}

func CloudflareTunnelDataSource() datasource.DataSource {
	return &cloudflareTunnelDataSource{}
}

func (d *cloudflareTunnelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *cloudflareTunnelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true, // Input
			},
			"name": schema.StringAttribute{
				Computed: true, // Output
			},
			"account_id": schema.StringAttribute{
				Computed: true, // Output
			},
			"tunnel_token": schema.StringAttribute{
				Computed:  true, // Output
				Sensitive: true,
			},
		},
	}
}

func (d *cloudflareTunnelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (d *cloudflareTunnelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudflareTunnelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	tunnel, err := d.client.GetTunnel(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tunnel",
			"Could not read tunnel ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update the model
	data.Name = types.StringValue(tunnel.Name)
	// Note: Secret might not be returned by GET, handle accordingly (leave null or set if returned)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
