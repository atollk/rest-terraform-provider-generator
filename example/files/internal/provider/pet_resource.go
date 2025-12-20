package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PetResource{}
var _ resource.ResourceWithImportState = &PetResource{}

func NewPetResource() resource.Resource {
	return &PetResource{}
}

// PetResource defines the resource implementation.
type PetResource struct {
	client *client.Client
}

// PetResourceModel describes the resource data model.
type PetResourceModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PhotoURLs types.List   `tfsdk:"photo_urls"`
	Status    types.String `tfsdk:"status"`
	Category  types.Object `tfsdk:"category"`
	Tags      types.List   `tfsdk:"tags"`
}

type CategoryModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type TagModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *PetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pet"
}

func (r *PetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Pet resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Pet identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Pet name",
				Required:            true,
			},
			"photo_urls": schema.ListAttribute{
				MarkdownDescription: "List of photo URLs",
				Required:            true,
				ElementType:         types.StringType,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Pet status in the store (available, pending, sold)",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.SingleNestedAttribute{
				MarkdownDescription: "Pet category",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "Category ID",
						Optional:            true,
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Category name",
						Optional:            true,
					},
				},
			},
			"tags": schema.ListNestedAttribute{
				MarkdownDescription: "List of tags",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Tag ID",
							Optional:            true,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Tag name",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *PetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	pet := &client.Pet{
		Name: data.Name.ValueString(),
	}

	// Handle photo URLs
	var photoURLs []string
	resp.Diagnostics.Append(data.PhotoURLs.ElementsAs(ctx, &photoURLs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pet.PhotoURLs = photoURLs

	// Handle status
	if !data.Status.IsNull() {
		pet.Status = data.Status.ValueString()
	}

	// Handle category
	if !data.Category.IsNull() {
		var category CategoryModel
		resp.Diagnostics.Append(data.Category.As(ctx, &category, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		pet.Category = &client.Category{
			Name: category.Name.ValueString(),
		}
		if !category.ID.IsNull() {
			pet.Category.ID = category.ID.ValueInt64()
		}
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tags []TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		pet.Tags = make([]client.Tag, len(tags))
		for i, tag := range tags {
			pet.Tags[i] = client.Tag{
				Name: tag.Name.ValueString(),
			}
			if !tag.ID.IsNull() {
				pet.Tags[i].ID = tag.ID.ValueInt64()
			}
		}
	}

	// Create the pet
	createdPet, err := r.client.CreatePet(pet)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pet, got error: %s", err))
		return
	}

	// Update model with response data
	data.ID = types.Int64Value(createdPet.ID)
	data.Name = types.StringValue(createdPet.Name)
	
	if createdPet.Status != "" {
		data.Status = types.StringValue(createdPet.Status)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get pet from API
	pet, err := r.client.GetPet(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pet, got error: %s", err))
		return
	}

	// Update model with fresh data
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	pet := &client.Pet{
		ID:   data.ID.ValueInt64(),
		Name: data.Name.ValueString(),
	}

	// Handle photo URLs
	var photoURLs []string
	resp.Diagnostics.Append(data.PhotoURLs.ElementsAs(ctx, &photoURLs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pet.PhotoURLs = photoURLs

	// Handle status
	if !data.Status.IsNull() {
		pet.Status = data.Status.ValueString()
	}

	// Handle category
	if !data.Category.IsNull() {
		var category CategoryModel
		resp.Diagnostics.Append(data.Category.As(ctx, &category, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		pet.Category = &client.Category{
			Name: category.Name.ValueString(),
		}
		if !category.ID.IsNull() {
			pet.Category.ID = category.ID.ValueInt64()
		}
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tags []TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		pet.Tags = make([]client.Tag, len(tags))
		for i, tag := range tags {
			pet.Tags[i] = client.Tag{
				Name: tag.Name.ValueString(),
			}
			if !tag.ID.IsNull() {
				pet.Tags[i].ID = tag.ID.ValueInt64()
			}
		}
	}

	// Update the pet
	updatedPet, err := r.client.UpdatePet(pet)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pet, got error: %s", err))
		return
	}

	// Update model with response data
	data.Name = types.StringValue(updatedPet.Name)
	data.Status = types.StringValue(updatedPet.Status)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the pet
	err := r.client.DeletePet(data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pet, got error: %s", err))
		return
	}
}

func (r *PetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
