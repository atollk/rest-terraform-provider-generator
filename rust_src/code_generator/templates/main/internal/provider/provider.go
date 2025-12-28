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
	"github.com/{{provider_info.author}}/terraform-provider-{{provider_info.name_kebab}}/internal/client"
)

// Ensure Provider satisfies various provider interfaces.
var _ provider.Provider = &Provider{}

// Provider defines the provider implementation.
type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	BaseURL types.String `tfsdk:"base_url"`
	APIKey  types.String `tfsdk:"api_key"`
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "{{provider_info.name_kebab}}"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL for the Petstore API. May also be provided via {{provider_info.name_caps}}_BASE_URL environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for authentication. May also be provided via {{provider_info.name_caps}}_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

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
			"Unknown {{provider_info.name_kebab}} API Base URL",
			"The provider cannot create the {{provider_info.name_kebab}} API client as there is an unknown configuration value for the {{provider_info.name_kebab}} API base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the {{provider_info.name_caps}}_BASE_URL environment variable.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown {{provider_info.name_kebab}} API Key",
			"The provider cannot create the {{provider_info.name_kebab}} API client as there is an unknown configuration value for the {{provider_info.name_kebab}} API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the {{provider_info.name_caps}}_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	baseURL := os.Getenv("{{provider_info.name_caps}}_BASE_URL")
	apiKey := os.Getenv("{{provider_info.name_caps}}_API_KEY")

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
			"Missing {{provider_info.name_kebab}} API Base URL",
			"The provider cannot create the {{provider_info.name_kebab}} API client as there is a missing or empty value for the {{provider_info.name_kebab}} API base URL. "+
				"Set the base_url value in the configuration or use the {{provider_info.name_caps}}_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new {{provider_info.name_kebab}} client using the configuration values
	client := client.NewClient(baseURL, apiKey)

	// Make the {{provider_info.name_kebab}} client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
	{% for resource in resources %}
	    New{{resource.name_pascal}},
    {% endfor %}
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	{% for data_source in data_sources %}
	    New{{data_source.name_pascal}},
    {% endfor %}
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}
