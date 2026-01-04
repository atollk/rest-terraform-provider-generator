package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"bytes"
	_ "embed"
	"strings"
	"text/template"

	"github.com/danielgtaylor/casing"
)

// ProviderInfo contains metadata and configuration for a Terraform provider.
type ProviderInfo struct {
	Name         string
	Author       string
	SpecDefaults *provider_spec.GlobalDefaults
}

// NameKebab returns the provider name in kebab-case format.
func (p *ProviderInfo) NameKebab() string {
	return casing.Kebab(p.Name)
}

// NameCaps returns the provider name in uppercase.
func (p *ProviderInfo) NameCaps() string {
	return strings.ToUpper(p.NameKebab())
}

// ResourceInfo contains metadata and configuration for a Terraform resource.
type ResourceInfo struct {
	Name         string
	ResourceSpec provider_spec.ResourceSchema
	OADoc        oas_parser.OADoc
}

// NameSnake returns the resource name in snake_case format.
func (r *ResourceInfo) NameSnake() string {
	return casing.Snake(r.Name)
}

// NamePascal returns the resource name in PascalCase format.
func (r *ResourceInfo) NamePascal() string {
	return casing.Camel(r.Name)
}

// DataSourceInfo contains metadata for a Terraform data source.
type DataSourceInfo struct {
	Name string
}

// NameSnake returns the data source name in snake_case format.
func (d *DataSourceInfo) NameSnake() string {
	return casing.Snake(d.Name)
}

// NamePascal returns the data source name in PascalCase format.
func (d *DataSourceInfo) NamePascal() string {
	return casing.Camel(d.Name)
}

// -------------------------------------------------------------------------------------------------

// templateRenderer defines an interface for rendering templates to files.
type templateRenderer interface {
	Name() string
	Render() ([]byte, error)
}

// basicTemplateRenderer is a simple implementation of templateRenderer.
type basicTemplateRenderer struct {
	name   string
	render func() ([]byte, error)
}

// Name returns the output file name for this template.
func (r *basicTemplateRenderer) Name() string {
	return r.name
}

// Render executes the template and returns the rendered content.
func (r *basicTemplateRenderer) Render() ([]byte, error) {
	return r.render()
}

//go:embed templates/main/Makefile.tmpl
var makefileTemplate string

// getMakefileTemplate creates a template renderer for the Makefile.
func getMakefileTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("Makefile", makefileTemplate, providerInfo)
}

//go:embed templates/main/main.go.tmpl
var mainGoTemplate string

// getMainGoTemplate creates a template renderer for the main.go file.
func getMainGoTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("main.go", mainGoTemplate, providerInfo)
}

//go:embed templates/main/go.mod.tmpl
var goModTemplate string

// getGoModTemplate creates a template renderer for the go.mod file.
func getGoModTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("go.mod", goModTemplate, providerInfo)
}

//go:embed templates/main/internal/provider/shared.go.tmpl
var sharedGoTemplate string

// getSharedGoTemplate creates a template renderer for the shared.go file.
func getSharedGoTemplate() templateRenderer {
	return renderTemplateAs("internal/provider/shared.go", sharedGoTemplate, nil)
}

// getOasJsonTemplate creates a template renderer for the OpenAPI specification JSON file.
func getOasJsonTemplate(oadoc oas_parser.OADoc) templateRenderer {
	renderFunc := func() ([]byte, error) {
		return oadoc.Model.RenderJSON(" ")
	}
	return &basicTemplateRenderer{name: "internal/provider/oas.json", render: renderFunc}
}

//go:embed templates/main/internal/provider/provider.go.tmpl
var providerGoTemplate string

// getProviderGoTemplate creates a template renderer for the provider.go file with the given provider info, resources, and data sources.
func getProviderGoTemplate(providerInfo *ProviderInfo, resources []ResourceInfo, dataSources []DataSourceInfo) templateRenderer {
	return renderTemplateAs("internal/provider/provider.go", providerGoTemplate, struct {
		ProviderInfo *ProviderInfo
		Resources    []ResourceInfo
		DataSources  []DataSourceInfo
	}{
		providerInfo, resources, dataSources,
	})
}

// renderTemplateAs creates a template renderer that parses and executes the given template with the provided arguments.
func renderTemplateAs(templateName string, templateData string, args any) templateRenderer {
	renderFunc := func() ([]byte, error) {
		templ, err := template.New(templateName).Parse(templateData)
		if err != nil {
			return nil, err
		}
		buffer := bytes.Buffer{}
		err = templ.Execute(&buffer, args)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return &basicTemplateRenderer{name: templateName, render: renderFunc}
}
