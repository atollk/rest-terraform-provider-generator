package code_generator

import (
	"fmt"
	"slices"

	"github.com/kaptinlin/messageformat-go/pkg/logger"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// property represents any kind of data object that is defined by a schema in OpenAPI.
type property struct {
	schema *base.Schema
}

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

func (p *property) GetTopSchemaType() string {
	if slices.Contains(p.schema.Type, "null") {
		if len(p.schema.Type) != 2 {
			return propertyTypeAny
		}
		if p.schema.Type[0] == "null" {
			return mapJsonSchemaToInternal(p.schema.Type[1])
		} else {
			return mapJsonSchemaToInternal(p.schema.Type[0])
		}
	} else {
		if len(p.schema.Type) != 1 {
			return propertyTypeAny
		} else {
			return mapJsonSchemaToInternal(p.schema.Type[0])
		}
	}
}

func (p *property) IsNullable() bool {
	hasNullableType := slices.Contains(p.schema.Type, "null")
	return (p.schema.Nullable != nil && *p.schema.Nullable) || hasNullableType
}

func (p *property) GetTypeType() string {
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

func (p *property) GetSchemaType() string {
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

func (p *property) GetValidatorType() string {
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

func (p *property) RenderSchemaCreation() string {
	return fmt.Sprintf(
		"schema.%s { Validators: []validator.%s { /* TODO */ } },\n",
		p.GetSchemaType(),
		p.GetValidatorType(),
	)
}

func newPropertyType(schema *base.Schema) property {
	return property{schema}
}
