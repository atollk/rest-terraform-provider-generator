package code_generator

import (
	_ "embed"
	"fmt"
	"slices"

	"github.com/danielgtaylor/casing"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// augmentedPropertySchema represents a resource property with metadata about where it appears in request/response bodies.
type augmentedPropertySchema struct {
	Name                string
	Schema              *base.Schema
	containedInBodyFlag int
	parent              *resourceTemplateRenderer
}

const (
	augmentedPropertySchemaCreateRequest  = 1 << iota
	augmentedPropertySchemaCreateResponse = 1 << iota
	augmentedPropertySchemaUpdateRequest  = 1 << iota
	augmentedPropertySchemaUpdateResponse = 1 << iota
)

const (
	propertyTypeBool   = "bool"
	propertyTypeInt    = "integer"
	propertyTypeFloat  = "float"
	propertyTypeString = "string"
	propertyTypeAny    = "any"
)

// mapJsonSchemaToInternal converts a JSON schema type to an internal property type constant.
func mapJsonSchemaToInternal(jsonSchemaType string) string {
	switch jsonSchemaType {
	case "null":
		return propertyTypeAny
	case "boolean":
		return propertyTypeBool
	case "integer":
		return propertyTypeInt
	case "number":
		return propertyTypeFloat
	case "string":
		return propertyTypeString
	case "object":
		// TODO
		return propertyTypeAny
	case "array":
		// TODO
		return propertyTypeAny
	default:
		return propertyTypeAny
	}
}

// GetTopSchemaType returns the primary type of the property, handling nullable types appropriately.
func (p *augmentedPropertySchema) GetTopSchemaType() string {
	if slices.Contains(p.Schema.Type, "null") {
		if len(p.Schema.Type) != 2 {
			return propertyTypeAny
		}
		if p.Schema.Type[0] == "null" {
			return mapJsonSchemaToInternal(p.Schema.Type[1])
		} else {
			return mapJsonSchemaToInternal(p.Schema.Type[0])
		}
	} else {
		if len(p.Schema.Type) != 1 {
			return propertyTypeAny
		} else {
			return mapJsonSchemaToInternal(p.Schema.Type[0])
		}
	}
}

// IsNullable returns true if the property can be null.
func (p *augmentedPropertySchema) IsNullable() bool {
	hasNullableType := slices.Contains(p.Schema.Type, "null")
	return (p.Schema.Nullable != nil && *p.Schema.Nullable) || hasNullableType
}

// GetTypeType returns the Terraform types package type name for this property (e.g., "String", "Int64").
func (p *augmentedPropertySchema) GetTypeType() string {
	switch p.GetTopSchemaType() {
	case propertyTypeBool:
		return "Bool"
	case propertyTypeInt:
		return "Int64"
	case propertyTypeFloat:
		return "Float64"
	case propertyTypeString:
		return "String"
	case propertyTypeAny:
		return "Dynamic"
	default:
		logger.Warn(fmt.Sprintf("Invalid property type enum found: %s  . Defaulting to 'any' type.", p.GetTopSchemaType()))
		return "Dynamic"
	}
}

// GetSchemaType returns the Terraform schema attribute type for this property (e.g., "StringAttribute").
func (p *augmentedPropertySchema) GetSchemaType() string {
	switch p.GetTopSchemaType() {
	case propertyTypeBool:
		return "BoolAttribute"
	case propertyTypeInt:
		return "Int64Attribute"
	case propertyTypeFloat:
		return "Float64Attribute"
	case propertyTypeString:
		return "StringAttribute"
	case propertyTypeAny:
		return "DynamicAttribute"
	default:
		logger.Warn(fmt.Sprintf("Invalid property type enum found: %s  . Defaulting to 'any' type.", p.GetTopSchemaType()))
		return "DynamicAttribute"
	}
}

// GetValidatorType returns the Terraform validator type for this property (e.g., "String", "Int64").
func (p *augmentedPropertySchema) GetValidatorType() string {
	switch p.GetTopSchemaType() {
	case propertyTypeBool:
		return "Bool"
	case propertyTypeInt:
		return "Int64"
	case propertyTypeFloat:
		return "Float64"
	case propertyTypeString:
		return "String"
	case propertyTypeAny:
		return "Dynamic"
	default:
		logger.Warn(fmt.Sprintf("Invalid property type enum found: %s  . Defaulting to 'any' type.", p.GetTopSchemaType()))
		return "Dynamic"
	}
}

// GetGoType returns the native Go type for this property (e.g., "string", "int64").
func (p *augmentedPropertySchema) GetGoType() string {
	switch p.GetTopSchemaType() {
	case propertyTypeBool:
		return "bool"
	case propertyTypeInt:
		return "int64"
	case propertyTypeFloat:
		return "float64"
	case propertyTypeString:
		return "string"
	case propertyTypeAny:
		return "any"
	default:
		logger.Warn(fmt.Sprintf("Invalid property type enum found: %s  . Defaulting to 'any' type.", p.GetTopSchemaType()))
		return "any"
	}
}

// RenderModelDataFields generates a Go struct field declaration for this property in the Terraform resource model.
func (p *augmentedPropertySchema) RenderModelDataFields() string {
	return fmt.Sprintf("%s types.%s `tfsdk:\"%s\"`", casing.Camel(p.Name), p.GetTypeType(), casing.Snake(p.Name))
}

// RenderAttributeDefinitions generates Terraform schema attribute definition code for this property.
func (p *augmentedPropertySchema) RenderAttributeDefinitions() string {
	schemaDef := fmt.Sprintf(
		"schema.%s { Validators: []validator.%s { &OpenApiSchemaValidator{ operationPath: \"%s\", operationMethod: \"%s\", propertyName: \"%s\" } } },",
		p.GetSchemaType(),
		p.GetValidatorType(),
		p.parent.GetCreatePath(),
		p.parent.GetCreateMethod(),
		p.Name,
	)
	return fmt.Sprintf(`"%s": %s`, p.Name, schemaDef)
}

// RenderFillCreateBody generates code to populate this property in the API request body during resource creation.
func (p *augmentedPropertySchema) RenderFillCreateBody() string {
	switch p.GetTopSchemaType() {
	case propertyTypeAny:
		fmtStr := `%[2]sUnpacked, err := UnpackDynamicType(&data.%[3]s, ctx)
if err != nil {
  resp.Diagnostics.AddError("cannot unpack data", fmt.Sprintf("cannot unpack data of property '%[1]s' due to error: %%v", err))
}
requestBody["%[1]s"] = %[2]sUnpacked`
		return fmt.Sprintf(fmtStr, p.Name, casing.LowerCamel(p.Name), casing.Camel(p.Name))
	default:
		fmtStr := `requestBody["%s"] = data.%s.Value%s()`
		return fmt.Sprintf(fmtStr, p.Name, casing.Camel(p.Name), p.GetTypeType())
	}
}

// RenderUpdateDataWithCreateResponse generates code to update Terraform state with this property's value from the API response.
func (p *augmentedPropertySchema) RenderUpdateDataWithCreateResponse() string {
	var fmtStr string
	switch p.GetTopSchemaType() {
	case propertyTypeAny:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, err := anyToDynamic(%[2]sRaw); err == nil {
		data.%[3]s = %[2]sValue
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", err.Error())
	}
}`
	default:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, ok := %[2]sRaw.(%[5]s); ok {
		data.%[3]s = types.%[4]sValue(%[2]sValue)
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", "'%[1]s' field is not of the expected type")
	}
}`
	}
	return fmt.Sprintf(fmtStr, p.Name, casing.LowerCamel(p.Name), casing.Camel(p.Name), p.GetTypeType(), p.GetGoType())
}

// RenderFillUpdateBody generates code to populate this property in the API request body during resource creation.
func (p *augmentedPropertySchema) RenderFillUpdateBody() string {
	switch p.GetTopSchemaType() {
	case propertyTypeAny:
		fmtStr := `%[2]sUnpacked, err := UnpackDynamicType(&data.%[3]s, ctx)
if err != nil {
  resp.Diagnostics.AddError("cannot unpack data", fmt.Sprintf("cannot unpack data of property '%[1]s' due to error: %%v", err))
}
requestBody["%[1]s"] = %[2]sUnpacked`
		return fmt.Sprintf(fmtStr, p.Name, casing.LowerCamel(p.Name), casing.Camel(p.Name))
	default:
		fmtStr := `requestBody["%s"] = data.%s.Value%s()`
		return fmt.Sprintf(fmtStr, p.Name, casing.Camel(p.Name), p.GetTypeType())
	}
}

// RenderUpdateDataWithUpdateResponse generates code to update Terraform state with this property's value from the API response.
func (p *augmentedPropertySchema) RenderUpdateDataWithUpdateResponse() string {
	var fmtStr string
	switch p.GetTopSchemaType() {
	case propertyTypeAny:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, err := anyToDynamic(%[2]sRaw); err == nil {
		data.%[3]s = %[2]sValue
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", err.Error())
	}
}`
	default:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, ok := %[2]sRaw.(%[5]s); ok {
		data.%[3]s = types.%[4]sValue(%[2]sValue)
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", "'%[1]s' field is not of the expected type")
	}
}`
	}
	return fmt.Sprintf(fmtStr, p.Name, casing.LowerCamel(p.Name), casing.Camel(p.Name), p.GetTypeType(), p.GetGoType())
}

// RenderUpdateDataWithReadResponse generates code to update Terraform state with this property's value from the API response.
func (p *augmentedPropertySchema) RenderUpdateDataWithReadResponse() string {
	var fmtStr string
	switch p.GetTopSchemaType() {
	case propertyTypeAny:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, err := anyToDynamic(%[2]sRaw); err == nil {
		data.%[3]s = %[2]sValue
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", err.Error())
	}
}`
	default:
		fmtStr = `if %[2]sRaw, ok := responseBody["%[1]s"]; ok {
	if %[2]sValue, ok := %[2]sRaw.(%[5]s); ok {
		data.%[3]s = types.%[4]sValue(%[2]sValue)
	} else {
		resp.Diagnostics.AddError("failure to set data '%[1]s'", "'%[1]s' field is not of the expected type")
	}
}`
	}
	return fmt.Sprintf(fmtStr, p.Name, casing.LowerCamel(p.Name), casing.Camel(p.Name), p.GetTypeType(), p.GetGoType())
}
