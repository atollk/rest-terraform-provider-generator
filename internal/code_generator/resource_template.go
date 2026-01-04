package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/danielgtaylor/casing"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/samber/lo"
)

//go:embed templates/main/internal/provider/resource.go.tmpl
var resourceGoTemplate string

// resourceTemplateRenderer implements templateRenderer for generating Terraform resource code.
type resourceTemplateRenderer struct {
	name         string
	ProviderInfo *ProviderInfo
	ResourceInfo ResourceDataSourceInfo
	IsDataSource bool
}

// Name returns the output file name for this resource template.
func (r *resourceTemplateRenderer) Name() string {
	return r.name
}

// getOperationBodies extracts the request and response schemas for a given OpenAPI operation.
// It returns the request schema, response schema, and any error encountered.
func (r *resourceTemplateRenderer) getOperationBodies(path string, operation string) (*base.Schema, *base.Schema, error) {
	pathObject, present := r.ResourceInfo.OADoc().Model.Paths.PathItems.Get(path)
	if !present {
		return nil, nil, errors.Errorf("could not find expected path %s", path)
	}
	opName := strings.ToLower(operation)
	op, present := pathObject.GetOperations().Get(opName)
	if !present {
		return nil, nil, errors.Errorf("could not find expected operation %s at path %s", opName, path)
	}
	requestContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, errors.Errorf("could not find expected request content type %s at operation %s at path %s", "application/json", opName, path)
	}
	responseContent, present := op.RequestBody.Content.Get("application/json")
	if !present {
		return nil, nil, errors.Errorf("could not find expected response content type %s at operation %s at path %s", "application/json", opName, path)
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
	path := ""
	opSpec := r.ResourceInfo.ResourceSpec().Create
	if opSpec != nil {
		path = opSpec.Path
	}
	if path == "" {
		path = r.ResourceInfo.ResourceSpec().Path
	}
	return path
}

// GetCreateMethod returns the HTTP method for creating resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetCreateMethod() string {
	op := r.ProviderInfo.SpecDefaults.CreateMethod
	if op == "" {
		op = r.ResourceInfo.ResourceSpec().Create.Method
	}
	return op
}

// GetUpdatePath returns the API path for updating resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetUpdatePath() string {
	return r.ResourceInfo.ResourceSpec().GetOperationPath(provider_spec.Update, r.ProviderInfo.SpecDefaults)
}

// GetUpdateMethod returns the HTTP method for updating resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetUpdateMethod() string {
	return r.ResourceInfo.ResourceSpec().GetOperationMethod(provider_spec.Update, r.ProviderInfo.SpecDefaults)
}

// GetDestroyPath returns the API path for deleting resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetDestroyPath() string {
	return r.ResourceInfo.ResourceSpec().GetOperationPath(provider_spec.Delete, r.ProviderInfo.SpecDefaults)
}

// GetDestroyMethod returns the HTTP method for deleting resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetDestroyMethod() string {
	return r.ResourceInfo.ResourceSpec().GetOperationMethod(provider_spec.Delete, r.ProviderInfo.SpecDefaults)
}

// GetReadPath returns the API path for reading resources, using the resource-specific path if defined.
func (r *resourceTemplateRenderer) GetReadPath() string {
	return r.ResourceInfo.ResourceSpec().GetOperationPath(provider_spec.Read, r.ProviderInfo.SpecDefaults)
}

// GetReadMethod returns the HTTP method for reading resources, preferring resource-specific configuration over provider defaults.
func (r *resourceTemplateRenderer) GetReadMethod() string {
	return r.ResourceInfo.ResourceSpec().GetOperationMethod(provider_spec.Read, r.ProviderInfo.SpecDefaults)
}

// getPropertiesFromBodies extracts and merges properties from create and update request/response bodies.
// It returns a list of augmented property schemas with metadata about which bodies contain each property.
func (r *resourceTemplateRenderer) getPropertiesFromBodies() ([]augmentedPropertySchema, error) {
	createRequestBody, createResponseBody, err := r.getOperationBodies(r.GetCreatePath(), r.GetCreateMethod())
	if err != nil {
		return nil, errors.Errorf("could not get request/response bodies for create: %w", err)
	}
	if !slices.Contains(createRequestBody.Type, "object") || !slices.Contains(createResponseBody.Type, "object") {
		return nil, errors.Errorf("only object types are supported for request/response bodies")
	}
	var updateRequestBody *base.Schema
	var updateResponseBody *base.Schema
	if !r.ResourceInfo.ResourceSpec().ForceRecreate {
		updateRequestBody, updateResponseBody, err = r.getOperationBodies(r.GetUpdatePath(), r.GetUpdateMethod())
		if err != nil {
			return nil, errors.Errorf("could not get request/response bodies for update: %w", err)
		}
		if !slices.Contains(updateRequestBody.Type, "object") || !slices.Contains(updateResponseBody.Type, "object") {
			return nil, errors.Errorf("only object types are supported for request/response bodies")
		}
	}
	dynamicTypeChecks := map[string][]string{
		"create request":  createRequestBody.Type,
		"create response": createResponseBody.Type,
	}
	if updateRequestBody != nil {
		dynamicTypeChecks["update request"] = updateRequestBody.Type
		dynamicTypeChecks["update response"] = updateResponseBody.Type
	}
	for bodyName, schemaType := range dynamicTypeChecks {
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
	if updateRequestBody != nil {
		err = parsePropertiesForBody("update request", updateRequestBody, augmentedPropertySchemaUpdateRequest)
		if err != nil {
			return nil, err
		}
		err = parsePropertiesForBody("update response", updateResponseBody, augmentedPropertySchemaUpdateResponse)
		if err != nil {
			return nil, err
		}
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

func (r *resourceTemplateRenderer) RenderRequestHeaders() (string, error) {
	result := &strings.Builder{}
	headers := r.ProviderInfo.SpecDefaults.Headers
	if headers != nil {
		for name, value := range headers.OtherProps {
			result.WriteString(fmt.Sprintf(`httpReq.Header.Set("%s", "%s")`, name, value))
			result.WriteRune('\n')
		}
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

func (r *resourceTemplateRenderer) RenderCreateRequestUrlExpression() (string, error) {
	properties, err := r.getPropertiesFromBodies()
	if err != nil {
		return "", errors.Errorf("could not get body properties: %w", err)
	}

	idAttribute := r.ResourceInfo.ResourceSpec().IdAttribute
	if idAttribute == "" {
		idAttribute = r.ProviderInfo.SpecDefaults.IdAttribute
	}
	if idAttribute == "" {
		idAttribute = "id"
	}
	idProp, found := lo.Find(properties, func(p augmentedPropertySchema) bool { return p.Name == idAttribute })
	if !found {
		return "", fmt.Errorf("could not find property with id_attribute name %s", idAttribute)
	}

	fmtStr := `strings.Replace(fmt.Sprintf("%%s%s", r.baseURL), "{%s}", fmt.Sprintf("%%v", data.%s.Value%s()), -1)`
	result := fmt.Sprintf(fmtStr, r.ResourceInfo.ResourceSpec().GetOperationPath(provider_spec.Create, r.ProviderInfo.SpecDefaults), r.ResourceInfo.ResourceSpec().IdAttributePath, casing.Camel(idProp.Name), idProp.GetTypeType())
	return result, nil
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

// RenderFillCreateBody generates code to populate the API request body from Terraform state during resource creation.
func (r *resourceTemplateRenderer) RenderFillUpdateBody() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderFillUpdateBody()
		},
	)
}

// RenderUpdateDataWithCreateResponse generates code to update Terraform state from API response data after resource creation.
func (r *resourceTemplateRenderer) RenderUpdateDataWithUpdateResponse() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderUpdateDataWithUpdateResponse()
		},
	)
}

// RenderUpdateDataWithReadResponse generates code to update Terraform state from API response data after resource creation.
func (r *resourceTemplateRenderer) RenderUpdateDataWithReadResponse() (string, error) {
	return r.renderForEachProp(
		func(prop *augmentedPropertySchema) string {
			return prop.RenderUpdateDataWithReadResponse()
		},
	)
}

// Render executes the resource template and returns the generated Go code.
func (r *resourceTemplateRenderer) Render() ([]byte, error) {
	return renderTemplateAs(r.name, resourceGoTemplate, r).Render()
}

// getResourceGoTemplate creates a template renderer for generating a Terraform resource implementation.
func getResourceGoTemplate(providerInfo *ProviderInfo, resourceInfo ResourceDataSourceInfo, isDataSource bool) templateRenderer {
	var namePrefix string
	if isDataSource {
		namePrefix = "data_source"
	} else {
		namePrefix = "resource"
	}
	return &resourceTemplateRenderer{
		name:         fmt.Sprintf("internal/provider/%s_%s.go", namePrefix, resourceInfo.NameSnake()),
		ProviderInfo: providerInfo,
		ResourceInfo: resourceInfo,
		IsDataSource: isDataSource,
	}
}
