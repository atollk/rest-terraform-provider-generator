package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"bytes"
	_ "embed"
	"fmt"
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

type ResourceDataSourceInfo interface {
	ParentProviderInfo() *ProviderInfo
	Name() string
	NameSnake() string
	NamePascal() string
	MainTypeName() string
	OADoc() oas_parser.OADoc
	ResourceSpec() *provider_spec.ResourceSchema
}

// ResourceInfo contains metadata and configuration for a Terraform resource.
type ResourceInfo struct {
	name         string
	resourceSpec provider_spec.ResourceSchema
	oadoc        oas_parser.OADoc
	providerInfo *ProviderInfo
}

func (r *ResourceInfo) Name() string {
	return r.name
}

// NameSnake returns the resource name in snake_case format.
func (r *ResourceInfo) NameSnake() string {
	return casing.Snake(r.name)
}

// NamePascal returns the resource name in PascalCase format.
func (r *ResourceInfo) NamePascal() string {
	return casing.Camel(r.name)
}

func (r *ResourceInfo) MainTypeName() string { return fmt.Sprintf("%sResource", r.NamePascal()) }

func (r *ResourceInfo) OADoc() oas_parser.OADoc {
	return r.oadoc
}

func (r *ResourceInfo) ParentProviderInfo() *ProviderInfo {
	return r.providerInfo
}

func (r *ResourceInfo) ResourceSpec() *provider_spec.ResourceSchema {
	return &r.resourceSpec
}

var _ ResourceDataSourceInfo = &ResourceInfo{}

// DataSourceInfo contains metadata for a Terraform data source.
type DataSourceInfo struct {
	name         string
	resourceSpec provider_spec.ResourceSchema
	oadoc        oas_parser.OADoc
	providerInfo *ProviderInfo
}

func (d *DataSourceInfo) Name() string {
	return d.name
}

// NameSnake returns the resource name in snake_case format.
func (d *DataSourceInfo) NameSnake() string {
	return casing.Snake(d.name)
}

// NamePascal returns the resource name in PascalCase format.
func (d *DataSourceInfo) NamePascal() string {
	return casing.Camel(d.name)
}

func (d *DataSourceInfo) MainTypeName() string { return fmt.Sprintf("%sDataSource", d.NamePascal()) }

func (d *DataSourceInfo) OADoc() oas_parser.OADoc {
	return d.oadoc
}

func (d *DataSourceInfo) ParentProviderInfo() *ProviderInfo {
	return d.providerInfo
}

func (d *DataSourceInfo) ResourceSpec() *provider_spec.ResourceSchema {
	return &d.resourceSpec
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
