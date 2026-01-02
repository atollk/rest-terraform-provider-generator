package provider_spec

import (
	_ "embed"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-yaml"
	"github.com/kaptinlin/jsonschema"

	"os"
)

//go:embed rest_api_provider_schema.json
var restApiProviderSchemaJson []byte

func ParseSpecFromFile(filename string) (RESTAPIProviderConfiguration, error) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return RESTAPIProviderConfiguration{}, errors.Errorf("error while reading file: %w", err)
	}
	return ParseSpec(fileContents)
}

func ParseSpec(fileContents []byte) (RESTAPIProviderConfiguration, error) {
	spec := RESTAPIProviderConfiguration{}
	rawJson, err := yaml.YAMLToJSON(fileContents)
	if err != nil {
		return spec, errors.Errorf("error while converting YAML to JSON: %w", err)
	}
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(restApiProviderSchemaJson)
	if err != nil {
		return spec, errors.Errorf("error while parsing JSON schema: %w", err)
	}
	validateResult := schema.ValidateJSON(rawJson)
	if !validateResult.IsValid() {
		return spec, errors.Errorf("invalid json for expected schema: %s", validateResult.Error())
	}
	err = schema.Unmarshal(&spec, rawJson)
	if err != nil {
		return spec, errors.Errorf("error while unmarshalling: %w", err)
	}
	return spec, nil
}
