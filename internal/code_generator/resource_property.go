package code_generator

import (
	_ "embed"
	"fmt"
	"slices"

	"github.com/danielgtaylor/casing"
	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

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

func (p *augmentedPropertySchema) IsNullable() bool {
	hasNullableType := slices.Contains(p.Schema.Type, "null")
	return (p.Schema.Nullable != nil && *p.Schema.Nullable) || hasNullableType
}

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

func (p *augmentedPropertySchema) RenderModelDataFields() string {
	return fmt.Sprintf("%s types.%s `tfsdk:\"%s\"`", casing.Camel(p.Name), p.GetTypeType(), casing.Snake(p.Name))
}

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
