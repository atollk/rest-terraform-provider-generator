package code_generator

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/danielgtaylor/casing"
	"github.com/pb33f/libopenapi/datamodel/high/base"
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

func (r *resourceTemplateRenderer) getOperationBodies(path string, operation string) (*base.Schema, *base.Schema, error) {
	pathObject, present := r.ResourceInfo.OADoc.Model.Paths.PathItems.Get(path)
	if !present {
		return nil, nil, fmt.Errorf("could not find expected path %v", path)
	}
	op, present := pathObject.GetOperations().Get(strings.ToLower(operation))
	if !present {
		return nil, nil, fmt.Errorf("could not find expected operation %v", operation)
	}
	requestContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, fmt.Errorf("could not find expected request content type %v", "application/json")
	}
	responseContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, fmt.Errorf("could not find expected response content type %v", "application/json")
	}
	requestSchema := requestContent.Schema.Schema()
	if requestSchema == nil {
		return nil, nil, fmt.Errorf("could not build schema: %w", requestContent.Schema.GetBuildError())
	}
	responseSchema := responseContent.Schema.Schema()
	if responseSchema == nil {
		return nil, nil, fmt.Errorf("could not build schema: %w", responseContent.Schema.GetBuildError())
	}
	return requestSchema, responseSchema, nil
}

func (r *resourceTemplateRenderer) getCreateBodies() (*base.Schema, *base.Schema, error) {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Create.Path
	}
	op := r.ProviderInfo.SpecDefaults.CreateMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Create.Method
	}
	return r.getOperationBodies(path, op)
}

func (r *resourceTemplateRenderer) getUpdateBodies() (*base.Schema, *base.Schema, error) {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Update.Path
	}
	op := r.ProviderInfo.SpecDefaults.UpdateMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Update.Method
	}
	return r.getOperationBodies(path, op)
}

func (r *resourceTemplateRenderer) RenderModelDataFields() (string, error) {
	// Retrieve the relevant schemas
	result := strings.Builder{}
	createRequestBody, createResponseBody, err := r.getCreateBodies()
	if err != nil {
		return "", fmt.Errorf("could not get request/response bodies for create: %w", err)
	}
	updateRequestBody, updateResponseBody, err := r.getUpdateBodies()
	if err != nil {
		return "", fmt.Errorf("could not get request/response bodies for update: %w", err)
	}
	if !slices.Contains(createRequestBody.Type, "object") || !slices.Contains(createResponseBody.Type, "object") || !slices.Contains(updateRequestBody.Type, "object") || !slices.Contains(updateResponseBody.Type, "object") {
		return "", fmt.Errorf("only object types are supported for request/response bodies")
	}

	// TODO: merge all four bodies

	// Render properties
	properties := createRequestBody.Properties
	for propName, propSchema := range properties.FromOldest() {
		schema := propSchema.Schema()
		if schema == nil {
			return "", fmt.Errorf("could not build schema of property %s: %w", propName, propSchema.GetBuildError())
		}
		attributeType := newPropertyType(schema)
		builderWriteStrings(
			&result,
			casing.Camel(propName),
			" types.",
			attributeType.GetTypeType(),
			" `tfsdk:\"",
			casing.Snake(propName),
			"\"`\n",
		)
	}
	return result.String(), nil
}

func (r *resourceTemplateRenderer) RenderAttributeDefinitions() (string, error) {
	// Retrieve the relevant schemas
	createRequestBody, createResponseBody, err := r.getCreateBodies()
	if err != nil {
		return "", fmt.Errorf("could not get request/response bodies for create: %w", err)
	}
	updateRequestBody, updateResponseBody, err := r.getUpdateBodies()
	if err != nil {
		return "", fmt.Errorf("could not get request/response bodies for update: %w", err)
	}
	if !slices.Contains(createRequestBody.Type, "object") || !slices.Contains(createResponseBody.Type, "object") || !slices.Contains(updateRequestBody.Type, "object") || !slices.Contains(updateResponseBody.Type, "object") {
		return "", fmt.Errorf("only object types are supported for request/response bodies")
	}

	// TODO: merge all four bodies

	// Render properties
	properties := createRequestBody.Properties
	result := strings.Builder{}
	for propName, propSchema := range properties.FromOldest() {
		schema := propSchema.Schema()
		if schema == nil {
			return "", fmt.Errorf("could not build schema of property %s: %w", propName, propSchema.GetBuildError())
		}
		attributeType := newPropertyType(schema)
		builderWriteStrings(
			&result,
			"\"",
			propName,
			"\": ",
			attributeType.RenderSchemaCreation(),
		)
	}
	return result.String(), nil
}

func builderWriteStrings(builder *strings.Builder, strs ...string) {
	for _, s := range strs {
		builder.WriteString(s)
	}
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
