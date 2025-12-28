package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

type ProviderInfo struct {
	Author    string
	NameKebab string
	NameCaps  string
}

type ResourceInfo struct {
	NameSnake    string
	NamePascal   string
	ResourceSpec provider_spec.ResourceSchema
	OADoc        oas_parser.OADoc
}

type DataSourceInfo struct {
	NameSnake  string
	NamePascal string
}

// -------------------------------------------------------------------------------------------------

type templateRenderer interface {
	Name() string
	Render() ([]byte, error)
}

type basicTemplateRenderer struct {
	name   string
	render func() ([]byte, error)
}

func (r *basicTemplateRenderer) Name() string {
	return r.name
}

func (r *basicTemplateRenderer) Render() ([]byte, error) {
	return r.render()
}

//go:embed templates/main/_Makefile
var makefileTemplate string

func getMakefileTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("Makefile", makefileTemplate, providerInfo)
}

//go:embed templates/main/_main.go
var mainGoTemplate string

func getMainGoTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("main.go", mainGoTemplate, providerInfo)
}

//go:embed templates/main/_go.mod
var goModTemplate string

func getGoModTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("go.mod", goModTemplate, providerInfo)
}

//go:embed templates/main/internal/provider/_provider.go
var providerGoTemplate string

func getProviderGoTemplate(providerInfo *ProviderInfo, resources []ResourceInfo, dataSources []DataSourceInfo) templateRenderer {
	return renderTemplateAs("internal/provider/provider.go", providerGoTemplate, struct {
		ProviderInfo *ProviderInfo
		Resources    []ResourceInfo
		DataSources  []DataSourceInfo
	}{
		providerInfo, resources, dataSources,
	})
}

//go:embed templates/main/internal/provider/_resource.go
var resourceGoTemplate string

type resourceTemplateRenderer struct {
	name         string
	ProviderInfo *ProviderInfo
	ResourceInfo *ResourceInfo
}

func (r *resourceTemplateRenderer) Name() string {
	return r.name
}

func (r *resourceTemplateRenderer) findRequestBodyProperties() (*orderedmap.Map[string, *base.SchemaProxy], error) {
	path, present := r.ResourceInfo.OADoc.Model.Paths.PathItems.Get(r.ResourceInfo.ResourceSpec.Path)
	if !present {
		return nil, fmt.Errorf("could not find expected path %v", r.ResourceInfo.ResourceSpec.Path)
	}
	op, present := path.GetOperations().Get("post")
	if !present {
		return nil, fmt.Errorf("could not find expected operation %v", "POST")
	}
	content, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, fmt.Errorf("could not find expected content type %v", "application/json")
	}
	schema := content.Schema.Schema()
	if schema == nil {
		return nil, fmt.Errorf("could not build schema: %w", content.Schema.GetBuildError())
	}
	if !slices.Contains(schema.Type, "object") {
		return nil, fmt.Errorf("only object types are supported for request bodies")
	}
	return schema.Properties, nil
}

func (r *resourceTemplateRenderer) renderAttributeDefinition(propertyName string, schemaProxy *base.SchemaProxy) (string, error) {
	schema := schemaProxy.Schema()
	if schema == nil {
		return "", fmt.Errorf("could not build schema: %w", schemaProxy.GetBuildError())
	}
	if len(schema.Type) != 1 {
		return "", fmt.Errorf("property schemas have to have exactly one type, but was: %v", schema.Type)
	}
	var attributeType string
	switch schema.Type[0] {
	case "boolean":
		attributeType = "Bool"
	case "integer":
		attributeType = "Int64"
	case "number":
		attributeType = "Float64"
	default:
		attributeType = "Dynamic"
	}
	return fmt.Sprintf("schema.%vAttribute {}", attributeType), nil
}

func (r *resourceTemplateRenderer) RenderAttributeDefinitions() (string, error) {
	result := strings.Builder{}
	properties, err := r.findRequestBodyProperties()
	if err != nil {
		return "", err
	}
	for propName, propSchema := range properties.FromOldest() {
		attributeDefinition, err := r.renderAttributeDefinition(propName, propSchema)
		if err != nil {
			return "", fmt.Errorf("could not render attribute definition of property: %w", err)
		}
		result.WriteString(fmt.Sprintf("\"%s\"", propName))
		result.WriteString(attributeDefinition)
		result.WriteString(",")
	}
	return result.String(), nil
}

func (r *resourceTemplateRenderer) Render() ([]byte, error) {
	return renderTemplateAs(r.name, resourceGoTemplate, r).Render()
}

func getResourceGoTemplate(providerInfo *ProviderInfo, resourceInfo *ResourceInfo) templateRenderer {
	return &resourceTemplateRenderer{
		name:         "internal/provider/resource.go",
		ProviderInfo: providerInfo,
		ResourceInfo: resourceInfo,
	}
}

/*
	struct MakefileTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct GoModTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct MainGoTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct ProviderGoTemplate<'a> {
		provider_info &'a ProviderInfo,
		resources &'a [ResourceInfo],
		data_sources &'a [DataSourceInfo],
	}

	struct ResourceGoTemplate<'a> {
		resource_info &'a ResourceInfo,
	}

	impl
crate::code_generator::helper_templates::ResourceGoTemplate < '_> {
	fn
	get_attributes_definition(&self)- > Vec < (string, GoCode)> {
	Vec::new()
	}
}

*/

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
