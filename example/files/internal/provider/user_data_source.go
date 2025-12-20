package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	client *client.Client
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Username   types.String `tfsdk:"username"`
	FirstName  types.String `tfsdk:"first_name"`
	LastName   types.String `tfsdk:"last_name"`
	Email      types.String `tfsdk:"email"`
	Phone      types.String `tfsdk:"phone"`
	UserStatus types.Int64  `tfsdk:"user_status"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "User data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "User identifier",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username",
				Required:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address",
				Computed:            true,
			},
			"phone": schema.StringAttribute{
				MarkdownDescription: "Phone number",
				Computed:            true,
			},
			"user_status": schema.Int64Attribute{
				MarkdownDescription: "User status",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get user from API
	user, err := d.client.GetUser(data.Username.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
		return
	}

	// Update model with API data
	data.ID = types.Int64Value(user.ID)
	data.Username = types.StringValue(user.Username)
	data.FirstName = types.StringValue(user.FirstName)
	data.LastName = types.StringValue(user.LastName)
	data.Email = types.StringValue(user.Email)
	data.Phone = types.StringValue(user.Phone)
	data.UserStatus = types.Int64Value(int64(user.UserStatus))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
