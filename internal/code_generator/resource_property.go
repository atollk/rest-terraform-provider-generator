package code_generator

import (
	_ "embed"
	"fmt"
	"slices"

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

func (p *augmentedPropertySchema) RenderSchemaCreation() string {
	return fmt.Sprintf(
		"schema.%s { Validators: []validator.%s { &OpenApiSchemaValidator{ operationPath: \"%s\", operationMethod: \"%s\", propertyName: \"%s\" } } },\n",
		p.GetSchemaType(),
		p.GetValidatorType(),
		p.parent.GetCreatePath(),
		p.parent.GetCreateMethod(),
		p.Name,
	)
}
