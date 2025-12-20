package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/yourusername/terraform-provider-petstore/internal/client"
)

// Ensure PetstoreProvider satisfies various provider interfaces.
var _ provider.Provider = &PetstoreProvider{}

// PetstoreProvider defines the provider implementation.
type PetstoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PetstoreProviderModel describes the provider data model.
type PetstoreProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
	APIKey  types.String `tfsdk:"api_key"`
}

func (p *PetstoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "petstore"
	resp.Version = p.version
}

func (p *PetstoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for the Petstore API. May also be provided via PETSTORE_BASE_URL environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authentication. May also be provided via PETSTORE_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PetstoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PetstoreProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// If configuration values are not yet known, the provider should not
	// attempt to configure the client.
	if data.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Petstore API Base URL",
			"The provider cannot create the Petstore API client as there is an unknown configuration value for the Petstore API base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PETSTORE_BASE_URL environment variable.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Petstore API Key",
			"The provider cannot create the Petstore API client as there is an unknown configuration value for the Petstore API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the PETSTORE_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	baseURL := os.Getenv("PETSTORE_BASE_URL")
	apiKey := os.Getenv("PETSTORE_API_KEY")

	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	}

	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing Petstore API Base URL",
			"The provider cannot create the Petstore API client as there is a missing or empty value for the Petstore API base URL. "+
				"Set the base_url value in the configuration or use the PETSTORE_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Petstore client using the configuration values
	client := client.NewClient(baseURL, apiKey)

	// Make the Petstore client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PetstoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPetResource,
		NewUserResource,
		NewOrderResource,
	}
}

func (p *PetstoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPetDataSource,
		NewUserDataSource,
		NewOrderDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PetstoreProvider{
			version: version,
		}
	}
}
