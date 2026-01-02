package code_generator

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/danielgtaylor/casing"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
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

func (r *resourceTemplateRenderer) getOperationBodies(path string, operation string) (*base.Schema, *base.Schema, error) {
	pathObject, present := r.ResourceInfo.OADoc.Model.Paths.PathItems.Get(path)
	if !present {
		return nil, nil, errors.Errorf("could not find expected path %v", path)
	}
	op, present := pathObject.GetOperations().Get(strings.ToLower(operation))
	if !present {
		return nil, nil, errors.Errorf("could not find expected operation %v", operation)
	}
	requestContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, errors.Errorf("could not find expected request content type %v", "application/json")
	}
	responseContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, errors.Errorf("could not find expected response content type %v", "application/json")
	}
	requestSchema := requestContent.Schema.Schema()
	if requestSchema == nil {
		return nil, nil, errors.Errorf("could not build schema: %w", requestContent.Schema.GetBuildError())
	}
	responseSchema := responseContent.Schema.Schema()
	if responseSchema == nil {
		return nil, nil, errors.Errorf("could not build schema: %w", responseContent.Schema.GetBuildError())
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

func (r *resourceTemplateRenderer) getPropertiesFromBodies() ([]augmentedPropertySchema, error) {
	createRequestBody, createResponseBody, err := r.getCreateBodies()
	if err != nil {
		return nil, errors.Errorf("could not get request/response bodies for create: %w", err)
	}
	updateRequestBody, updateResponseBody, err := r.getUpdateBodies()
	if err != nil {
		return nil, errors.Errorf("could not get request/response bodies for update: %w", err)
	}
	if !slices.Contains(createRequestBody.Type, "object") || !slices.Contains(createResponseBody.Type, "object") || !slices.Contains(updateRequestBody.Type, "object") || !slices.Contains(updateResponseBody.Type, "object") {
		return nil, errors.Errorf("only object types are supported for request/response bodies")
	}
	for bodyName, schemaType := range map[string][]string{
		"create request":  createRequestBody.Type,
		"create response": createResponseBody.Type,
		"update request":  updateRequestBody.Type,
		"update response": updateResponseBody.Type,
	} {
		if !slices.Equal(schemaType, []string{"object"}) || slices.Equal(schemaType, []string{"object", "null"}) || slices.Equal(schemaType, []string{"null", "object"}) {
			logger.Warn(fmt.Sprintf("%s body has unexpected schema type %v; will default to dynamic type", bodyName, schemaType))
			return nil, nil
		}
	}

	propertyMap := orderedmap.New[string, *augmentedPropertySchema]()
	parsePropertiesForBody := func(bodyName string, bodySchema *base.Schema, flag int) error {
		for propertyName, propertySchemaProxy := range createRequestBody.Properties.FromOldest() {
			propertySchema := propertySchemaProxy.Schema()
			if propertySchema == nil {
				return errors.Errorf("could not get schema for property %s in %s body", propertyName, bodyName)
			}
			entry, exists := propertyMap.Get(propertyName)
			if !exists {
				entry = &augmentedPropertySchema{Name: propertyName, Schema: propertySchema, parent: r}
				propertyMap.Set(propertyName, entry)
			}
			entry.containedInBodyFlag = entry.containedInBodyFlag | flag
		}
		return nil
	}
	err = parsePropertiesForBody("create request", createRequestBody, augmentedPropertySchemaCreateRequest)
	if err != nil {
		return nil, err
	}
	err = parsePropertiesForBody("create response", createResponseBody, augmentedPropertySchemaCreateResponse)
	if err != nil {
		return nil, err
	}
	err = parsePropertiesForBody("update request", updateRequestBody, augmentedPropertySchemaUpdateRequest)
	if err != nil {
		return nil, err
	}
	err = parsePropertiesForBody("update response", updateResponseBody, augmentedPropertySchemaUpdateResponse)
	if err != nil {
		return nil, err
	}

	var result []augmentedPropertySchema
	for prop := range propertyMap.ValuesFromOldest() {
		result = append(result, *prop)
	}
	return result, nil
}

func (r *resourceTemplateRenderer) RenderModelDataFields() (string, error) {
	properties, err := r.getPropertiesFromBodies()
	if err != nil {
		return "", errors.Errorf("could not get body properties: %w", err)
	}

	result := strings.Builder{}
	for _, prop := range properties {
		builderWriteStrings(
			&result,
			casing.Camel(prop.Name),
			" types.",
			prop.GetTypeType(),
			" `tfsdk:\"",
			casing.Snake(prop.Name),
			"\"`\n",
		)
	}
	return result.String(), nil
}

func (r *resourceTemplateRenderer) RenderAttributeValidators() (string, error) {
	properties, err := r.getPropertiesFromBodies()
	if err != nil {
		return "", errors.Errorf("could not get body properties: %w", err)
	}

	result := strings.Builder{}
	for _, prop := range properties {
		result.WriteString(prop.RenderValidatorType())
	}
	return result.String(), nil

}

func (r *resourceTemplateRenderer) RenderAttributeDefinitions() (string, error) {
	properties, err := r.getPropertiesFromBodies()
	if err != nil {
		return "", errors.Errorf("could not get body properties: %w", err)
	}

	result := strings.Builder{}
	for _, prop := range properties {
		builderWriteStrings(
			&result,
			"\"",
			prop.Name,
			"\": ",
			prop.RenderSchemaCreation(),
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
