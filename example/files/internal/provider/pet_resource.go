package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PetResource{}
var _ resource.ResourceWithImportState = &PetResource{}

func NewPetResource() resource.Resource {
	return &PetResource{}
}

// PetResource defines the resource implementation.
type PetResource struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
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

// HTTPConfig holds HTTP client configuration
type HTTPConfig struct {
	BaseURL string
	APIKey  string
}

func (r *PetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*HTTPConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *HTTPConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.baseURL = config.BaseURL
	r.apiKey = config.APIKey
	r.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}
}

// doRequest performs an HTTP request and returns the response body
func (r *PetResource) doRequest(method, url string, body map[string]interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(data)
	}

	httpReq, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if r.apiKey != "" {
		httpReq.Header.Set("api_key", r.apiKey)
	}

	res, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.Errorf("API request failed with status %d: %s", res.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *PetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build request body as a map
	petData := map[string]interface{}{
		"name": data.Name.ValueString(),
	}

	// Handle photo URLs
	var photoURLs []string
	resp.Diagnostics.Append(data.PhotoURLs.ElementsAs(ctx, &photoURLs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	petData["photoUrls"] = photoURLs

	// Handle status
	if !data.Status.IsNull() {
		petData["status"] = data.Status.ValueString()
	}

	// Handle category
	if !data.Category.IsNull() {
		var category CategoryModel
		resp.Diagnostics.Append(data.Category.As(ctx, &category, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		categoryData := map[string]interface{}{
			"name": category.Name.ValueString(),
		}
		if !category.ID.IsNull() {
			categoryData["id"] = category.ID.ValueInt64()
		}
		petData["category"] = categoryData
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tags []TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tagsData := make([]map[string]interface{}, len(tags))
		for i, tag := range tags {
			tagData := map[string]interface{}{
				"name": tag.Name.ValueString(),
			}
			if !tag.ID.IsNull() {
				tagData["id"] = tag.ID.ValueInt64()
			}
			tagsData[i] = tagData
		}
		petData["tags"] = tagsData
	}

	// Create the pet
	createdPet, err := r.doRequest("POST", fmt.Sprintf("%s/pet", r.baseURL), petData)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pet, got error: %s", err))
		return
	}

	// Update model with response data
	if id, ok := createdPet["id"].(float64); ok {
		data.ID = types.Int64Value(int64(id))
	}
	if name, ok := createdPet["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if status, ok := createdPet["status"].(string); ok && status != "" {
		data.Status = types.StringValue(status)
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
	pet, err := r.doRequest("GET", fmt.Sprintf("%s/pet/%d", r.baseURL, data.ID.ValueInt64()), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pet, got error: %s", err))
		return
	}

	// Update model with fresh data
	if name, ok := pet["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if status, ok := pet["status"].(string); ok {
		data.Status = types.StringValue(status)
	}

	// Convert photo URLs
	if photoURLsIface, ok := pet["photoUrls"].([]interface{}); ok {
		photoURLs := make([]types.String, len(photoURLsIface))
		for i, url := range photoURLsIface {
			if urlStr, ok := url.(string); ok {
				photoURLs[i] = types.StringValue(urlStr)
			}
		}
		photoURLsList, diags := types.ListValueFrom(ctx, types.StringType, photoURLs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.PhotoURLs = photoURLsList
	}

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

	// Build request body as a map
	petData := map[string]interface{}{
		"id":   data.ID.ValueInt64(),
		"name": data.Name.ValueString(),
	}

	// Handle photo URLs
	var photoURLs []string
	resp.Diagnostics.Append(data.PhotoURLs.ElementsAs(ctx, &photoURLs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	petData["photoUrls"] = photoURLs

	// Handle status
	if !data.Status.IsNull() {
		petData["status"] = data.Status.ValueString()
	}

	// Handle category
	if !data.Category.IsNull() {
		var category CategoryModel
		resp.Diagnostics.Append(data.Category.As(ctx, &category, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		categoryData := map[string]interface{}{
			"name": category.Name.ValueString(),
		}
		if !category.ID.IsNull() {
			categoryData["id"] = category.ID.ValueInt64()
		}
		petData["category"] = categoryData
	}

	// Handle tags
	if !data.Tags.IsNull() {
		var tags []TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tagsData := make([]map[string]interface{}, len(tags))
		for i, tag := range tags {
			tagData := map[string]interface{}{
				"name": tag.Name.ValueString(),
			}
			if !tag.ID.IsNull() {
				tagData["id"] = tag.ID.ValueInt64()
			}
			tagsData[i] = tagData
		}
		petData["tags"] = tagsData
	}

	// Update the pet
	updatedPet, err := r.doRequest("PUT", fmt.Sprintf("%s/pet", r.baseURL), petData)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pet, got error: %s", err))
		return
	}

	// Update model with response data
	if name, ok := updatedPet["name"].(string); ok {
		data.Name = types.StringValue(name)
	}
	if status, ok := updatedPet["status"].(string); ok {
		data.Status = types.StringValue(status)
	}

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
	_, err := r.doRequest("DELETE", fmt.Sprintf("%s/pet/%d", r.baseURL, data.ID.ValueInt64()), nil)
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
