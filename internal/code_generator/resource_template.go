package code_generator

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

//go:embed templates/main/internal/provider/resource.go.tmpl
var resourceGoTemplate string

// resourceTemplateRenderer implements templateRenderer for generating Terraform resource code.
type resourceTemplateRenderer struct {
	name         string
	ProviderInfo *ProviderInfo
	ResourceInfo *ResourceInfo
}

// Name returns the output file name for this resource template.
func (r *resourceTemplateRenderer) Name() string {
	return r.name
}

// getOperationBodies extracts the request and response schemas for a given OpenAPI operation.
// It returns the request schema, response schema, and any error encountered.
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

// GetCreatePath returns the API path for creating resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetCreatePath() string {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Create.Path
	}
	return path
}

// GetCreateMethod returns the HTTP method for creating resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetCreateMethod() string {
	op := r.ProviderInfo.SpecDefaults.CreateMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Create.Method
	}
	return op
}

// GetUpdatePath returns the API path for updating resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetUpdatePath() string {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Update.Path
	}
	return path
}

// GetUpdateMethod returns the HTTP method for updating resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetUpdateMethod() string {
	op := r.ProviderInfo.SpecDefaults.UpdateMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Update.Method
	}
	return op
}

// GetDestroyPath returns the API path for deleting resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetDestroyPath() string {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Destroy.Path
	}
	return path
}

// GetDestroyMethod returns the HTTP method for deleting resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetDestroyMethod() string {
	op := r.ProviderInfo.SpecDefaults.DestroyMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Destroy.Method
	}
	return op
}

// GetReadPath returns the API path for reading resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetReadPath() string {
	path := r.ResourceInfo.ResourceSpec.Path
	if path == "" {
		path = r.ResourceInfo.ResourceSpec.Read.Path
	}
	return path
}

// GetReadMethod returns the HTTP method for reading resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetReadMethod() string {
	op := r.ProviderInfo.SpecDefaults.ReadMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec.Read.Method
	}
	return op
}

// getPropertiesFromBodies extracts and merges properties from create and update request/response bodies.
// It returns a list of augmented property schemas with metadata about which bodies contain each property.
func (r *resourceTemplateRenderer) getPropertiesFromBodies() ([]augmentedPropertySchema, error) {
	createRequestBody, createResponseBody, err := r.getOperationBodies(r.GetCreatePath(), r.GetCreateMethod())
	if err != nil {
		return nil, errors.Errorf("could not get request/response bodies for create: %w", err)
	}
	updateRequestBody, updateResponseBody, err := r.getOperationBodies(r.GetUpdatePath(), r.GetUpdateMethod())
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

// renderForEachProp applies a rendering function to each property and concatenates the results.
func (r *resourceTemplateRenderer) renderForEachProp(f func(*augmentedPropertySchema) string) (string, error) {
	properties, err := r.getPropertiesFromBodies()
	if err != nil {
		return "", errors.Errorf("could not get body properties: %w", err)
	}

	result := strings.Builder{}
	for _, prop := range properties {
		result.WriteString(f(&prop))
		result.WriteRune('\n')
	}
	return result.String(), nil
}

// RenderModelDataFields generates Go struct field declarations for the Terraform resource model.
func (r *resourceTemplateRenderer) RenderModelDataFields() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderModelDataFields()
		},
	)
}

// RenderAttributeDefinitions generates Terraform schema attribute definitions for all resource properties.
func (r *resourceTemplateRenderer) RenderAttributeDefinitions() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderAttributeDefinitions()
		},
	)
}

// RenderFillCreateBody generates code to populate the API request body from Terraform state during resource creation.
func (r *resourceTemplateRenderer) RenderFillCreateBody() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderFillCreateBody()
		},
	)
}

// RenderUpdateDataWithCreateResponse generates code to update Terraform state from API response data after resource creation.
func (r *resourceTemplateRenderer) RenderUpdateDataWithCreateResponse() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderUpdateDataWithCreateResponse()
		},
	)
}

// Render executes the resource template and returns the generated Go code.
func (r *resourceTemplateRenderer) Render() ([]byte, error) {
	return renderTemplateAs(r.name, resourceGoTemplate, r).Render()
}

// getResourceGoTemplate creates a template renderer for generating a Terraform resource implementation.
func getResourceGoTemplate(providerInfo *ProviderInfo, resourceInfo *ResourceInfo) templateRenderer {
	return &resourceTemplateRenderer{
		name:         "internal/provider/resource.go",
		ProviderInfo: providerInfo,
		ResourceInfo: resourceInfo,
	}
}
