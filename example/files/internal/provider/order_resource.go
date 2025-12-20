package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OrderResource{}
var _ resource.ResourceWithImportState = &OrderResource{}

func NewOrderResource() resource.Resource {
	return &OrderResource{}
}

// OrderResource defines the resource implementation.
type OrderResource struct {
	client *client.Client
}

// OrderResourceModel describes the resource data model.
type OrderResourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	PetID    types.Int64  `tfsdk:"pet_id"`
	Quantity types.Int64  `tfsdk:"quantity"`
	ShipDate types.String `tfsdk:"ship_date"`
	Status   types.String `tfsdk:"status"`
	Complete types.Bool   `tfsdk:"complete"`
}

func (r *OrderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

func (r *OrderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Order resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Order identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pet_id": schema.Int64Attribute{
				MarkdownDescription: "Pet identifier",
				Required:            true,
			},
			"quantity": schema.Int64Attribute{
				MarkdownDescription: "Quantity of pets ordered",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ship_date": schema.StringAttribute{
				MarkdownDescription: "Ship date (RFC3339 format)",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Order status (placed, approved, delivered)",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"complete": schema.BoolAttribute{
				MarkdownDescription: "Whether the order is complete",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OrderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *OrderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	order := &client.Order{
		PetID: data.PetID.ValueInt64(),
	}

	if !data.Quantity.IsNull() {
		order.Quantity = int32(data.Quantity.ValueInt64())
	}

	if !data.ShipDate.IsNull() {
		shipDate, err := time.Parse(time.RFC3339, data.ShipDate.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Ship Date", fmt.Sprintf("Unable to parse ship_date, got error: %s", err))
			return
		}
		order.ShipDate = shipDate
	}

	if !data.Status.IsNull() {
		order.Status = data.Status.ValueString()
	}

	if !data.Complete.IsNull() {
		order.Complete = data.Complete.ValueBool()
	}

	// Create the order
	createdOrder, err := r.client.CreateOrder(order)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create order, got error: %s", err))
		return
	}

	// Update model with response data
	data.ID = types.Int64Value(createdOrder.ID)
	data.PetID = types.Int64Value(createdOrder.PetID)
	data.Quantity = types.Int64Value(int64(createdOrder.Quantity))
	data.Status = types.StringValue(createdOrder.Status)
	data.Complete = types.BoolValue(createdOrder.Complete)

	if !createdOrder.ShipDate.IsZero() {
		data.ShipDate = types.StringValue(createdOrder.ShipDate.Format(time.RFC3339))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get order from API
	order, err := r.client.GetOrder(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read order, got error: %s", err))
		return
	}

	// Update model with fresh data
	data.PetID = types.Int64Value(order.PetID)
	data.Quantity = types.Int64Value(int64(order.Quantity))
	data.Status = types.StringValue(order.Status)
	data.Complete = types.BoolValue(order.Complete)

	if !order.ShipDate.IsZero() {
		data.ShipDate = types.StringValue(order.ShipDate.Format(time.RFC3339))
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The Petstore API doesn't support updating orders
	// You would need to delete and recreate
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"The Petstore API does not support updating orders. To modify an order, delete and recreate it.",
	)
}

func (r *OrderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the order
	err := r.client.DeleteOrder(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete order, got error: %s", err))
		return
	}
}

func (r *OrderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID as an integer
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("The import ID must be a valid integer, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
