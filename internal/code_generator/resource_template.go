package code_generator

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

//go:embed templates/main/internal/provider/resource.go.tmpl
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

func getPropertyTypeStrings(propertySchema *base.Schema) (string, error) {
	var t string
	if len(propertySchema.Type) != 1 {
		return "", fmt.Errorf("property schemas have to have exactly one type, but was: %v", propertySchema.Type)
	}
	switch propertySchema.Type[0] {
	case "boolean":
		t = "Bool"
	case "integer":
		t = "Int64"
	case "number":
		t = "Float64"
	case "string":
		t = "String"
	default:
		t = "Dynamic"
	}
	return t, nil
}

func (r *resourceTemplateRenderer) renderModelDataField(propertyName string, schemaProxy *base.SchemaProxy) (string, error) {
	schema := schemaProxy.Schema()
	if schema == nil {
		return "", fmt.Errorf("could not build schema: %w", schemaProxy.GetBuildError())
	}
	attributeType, err := getPropertyTypeStrings(schema)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s types.%s `tfsdk:\"%s\"`", casing.Camel(propertyName), attributeType, casing.Snake(propertyName)), nil
}

func (r *resourceTemplateRenderer) renderModelDataFields() (string, error) {
	result := strings.Builder{}
	properties, err := r.findRequestBodyProperties()
	if err != nil {
		return "", err
	}
	for propName, propSchema := range properties.FromOldest() {
		attributeDefinition, err := r.renderModelDataField(propName, propSchema)
		if err != nil {
			return "", fmt.Errorf("could not render attribute definition of property: %w", err)
		}
		result.WriteString(fmt.Sprintf("\"%s\": ", propName))
		result.WriteString(attributeDefinition)
		result.WriteString(",\n")
	}
	return result.String(), nil
}

func (r *resourceTemplateRenderer) renderAttributeDefinition(propertyName string, schemaProxy *base.SchemaProxy) (string, error) {
	schema := schemaProxy.Schema()
	if schema == nil {
		return "", fmt.Errorf("could not build schema: %w", schemaProxy.GetBuildError())
	}
	attributeType, err := getPropertyTypeStrings(schema)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("schema.%vAttribute { Validators: []validator.%v { /* TODO */ } }", attributeType, attributeType), nil
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
		result.WriteString(fmt.Sprintf("\"%s\": ", propName))
		result.WriteString(attributeDefinition)
		result.WriteString(",\n")
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
