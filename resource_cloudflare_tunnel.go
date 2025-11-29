package main

import (
	"context"
	"fmt"

	"terraform-provider-cloudflare-tunnel/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &cloudflareTunnelResource{}
)

func CloudflareTunnelResource() resource.Resource {
	return &cloudflareTunnelResource{}
}

type cloudflareTunnelResource struct {
	client *client.Client
}

type cloudflareTunnelResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	TunnelToken types.String `tfsdk:"tunnel_token"`
}

func (r *cloudflareTunnelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *cloudflareTunnelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (r *cloudflareTunnelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"tunnel_token": schema.StringAttribute{
				Required:  false,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *cloudflareTunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data cloudflareTunnelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	tunnel, err := r.client.CreateTunnel(data.Name.ValueString(), data.TunnelToken.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tunnel",
			"Could not create tunnel, unexpected error: "+err.Error(),
		)
		return
	}

	// Update the model with values from the API
	data.ID = types.StringValue(tunnel.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudflareTunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data cloudflareTunnelResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	tunnel, err := r.client.GetTunnel(data.ID.ValueString())
	if err != nil {
		// If 404, remove from state
		// Note: You might need to check the error message or status code specifically
		// For now, assuming any error means we can't read it, but ideally check for "not found"
		resp.Diagnostics.AddError(
			"Error reading tunnel",
			"Could not read tunnel ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update the model with values from the API
	data.Name = types.StringValue(tunnel.Name)
	// data.Secret is sensitive and usually not returned by GET, so we keep the existing value
	// unless we want to mark it as unknown if it changed. For now, keep as is.

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudflareTunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data cloudflareTunnelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	_, err := r.client.UpdateTunnel(data.ID.ValueString(), data.Name.ValueString(), data.TunnelToken.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tunnel",
			"Could not update tunnel ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudflareTunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data cloudflareTunnelResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	_, err := r.client.DeleteTunnel(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tunnel",
			"Could not delete tunnel ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
