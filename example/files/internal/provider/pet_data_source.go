package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PetDataSource{}

func NewPetDataSource() datasource.DataSource {
	return &PetDataSource{}
}

// PetDataSource defines the data source implementation.
type PetDataSource struct {
	client *client.Client
}

// PetDataSourceModel describes the data source data model.
type PetDataSourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PhotoURLs types.List   `tfsdk:"photo_urls"`
	Status    types.String `tfsdk:"status"`
	Category  types.Object `tfsdk:"category"`
	Tags      types.List   `tfsdk:"tags"`
}

func (d *PetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pet"
}

func (d *PetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Pet data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Pet identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Pet name",
				Computed:            true,
			},
			"photo_urls": schema.ListAttribute{
				MarkdownDescription: "List of photo URLs",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Pet status in the store",
				Computed:            true,
			},
			"category": schema.SingleNestedAttribute{
				MarkdownDescription: "Pet category",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "Category ID",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Category name",
						Computed:            true,
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "List of tags",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Tag ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Tag name",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *PetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get pet from API
	pet, err := d.client.GetPet(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pet, got error: %s", err))
		return
	}

	// Update model with API data
	data.Name = types.StringValue(pet.Name)
	data.Status = types.StringValue(pet.Status)

	// Convert photo URLs
	photoURLs := make([]types.String, len(pet.PhotoURLs))
	for i, url := range pet.PhotoURLs {
		photoURLs[i] = types.StringValue(url)
	}
	photoURLsList, diags := types.ListValueFrom(ctx, types.StringType, photoURLs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.PhotoURLs = photoURLsList

	// Convert category if present
	if pet.Category != nil {
		categoryObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"id":   types.Int64Type,
			"name": types.StringType,
		}, map[string]interface{}{
			"id":   pet.Category.ID,
			"name": pet.Category.Name,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Category = categoryObj
	}

	// Convert tags if present
	if len(pet.Tags) > 0 {
		tags := make([]TagModel, len(pet.Tags))
		for i, tag := range pet.Tags {
			tags[i] = TagModel{
				ID:   types.Int64Value(tag.ID),
				Name: types.StringValue(tag.Name),
			}
		}
		tagsList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.Int64Type,
				"name": types.StringType,
			},
		}, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsList
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
