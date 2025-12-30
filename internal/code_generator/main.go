package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"fmt"
	"os"
	"path"
)

func RenderSpec(
	output_path string,
	provider_spec provider_spec.RESTAPIProviderConfiguration,
	api_spec oas_parser.OADoc,
) error {
	// Prepare input data
	providerInfo := ProviderInfo{
		Author:    "foo",
		NameKebab: "petstore",
		NameCaps:  "PETSTORE",
	}
	resources := []ResourceInfo{
		{NameSnake: "res_one", NamePascal: "ResOne", ResourceSpec: provider_spec.Resources.OtherProps["pet"], OADoc: api_spec},
	}
	var dataSources []DataSourceInfo

	// Map output file names to templates
	templates := []templateRenderer{
		getMakefileTemplate(&providerInfo),
		getMainGoTemplate(&providerInfo),
		getGoModTemplate(&providerInfo),
		getSharedGoTemplate(),
		getOasJsonTemplate(api_spec),
		getProviderGoTemplate(&providerInfo, resources, dataSources),
	}
	for _, resource := range resources {
		templates = append(templates, getResourceGoTemplate(&providerInfo, &resource))
	}

	// Write out files
	for _, renderer := range templates {
		completePath := path.Join(output_path, renderer.Name())
		err := os.MkdirAll(path.Dir(completePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create directories: %w", err)
		}
		output, err := renderer.Render()
		if err != nil {
			return fmt.Errorf("cannot execute template: %w", err)
		}
		err = os.WriteFile(completePath, output, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot write file: %w", err)
		}
	}

	return nil
}
