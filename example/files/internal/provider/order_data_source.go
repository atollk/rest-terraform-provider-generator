package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OrderDataSource{}

func NewOrderDataSource() datasource.DataSource {
	return &OrderDataSource{}
}

// OrderDataSource defines the data source implementation.
type OrderDataSource struct {
	client *client.Client
}

// OrderDataSourceModel describes the data source data model.
type OrderDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	PetID    types.Int64  `tfsdk:"pet_id"`
	Quantity types.Int64  `tfsdk:"quantity"`
	ShipDate types.String `tfsdk:"ship_date"`
	Status   types.String `tfsdk:"status"`
	Complete types.Bool   `tfsdk:"complete"`
}

func (d *OrderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

func (d *OrderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Order data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Order identifier",
				Required:            true,
			},
			"pet_id": schema.Int64Attribute{
				MarkdownDescription: "Pet identifier",
				Computed:            true,
			},
			"quantity": schema.Int64Attribute{
				MarkdownDescription: "Quantity of pets ordered",
				Computed:            true,
			},
			"ship_date": schema.StringAttribute{
				MarkdownDescription: "Ship date (RFC3339 format)",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Order status",
				Computed:            true,
			},
			"complete": schema.BoolAttribute{
				MarkdownDescription: "Whether the order is complete",
				Computed:            true,
			},
		},
	}
}

func (d *OrderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *OrderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrderDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API
	order, err := d.client.GetOrder(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read order, got error: %s", err))
		return
	}

	// Update model with API data
	data.PetID = types.Int64Value(order.PetID)
	data.Quantity = types.Int64Value(int64(order.Quantity))
	data.Status = types.StringValue(order.Status)
	data.Complete = types.BoolValue(order.Complete)

	if !order.ShipDate.IsZero() {
		data.ShipDate = types.StringValue(order.ShipDate.Format(time.RFC3339))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
