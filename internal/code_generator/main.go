package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"fmt"
	"os"
	"path"
)

func RenderSpec(
	outputPath string,
	providerSpec provider_spec.RESTAPIProviderConfiguration,
	apiSpec oas_parser.OADoc,
) error {
	// Prepare input data
	providerInfo := ProviderInfo{
		Author:       "foo",
		Name:         "pet_store",
		SpecDefaults: providerSpec.GlobalDefaults,
	}
	resources := []ResourceInfo{
		{
			Name:         "res_one",
			ResourceSpec: providerSpec.Resources.OtherProps["pet"],
			OADoc:        apiSpec,
		},
	}
	var dataSources []DataSourceInfo

	// Map output file names to templates
	templates := []templateRenderer{
		getMakefileTemplate(&providerInfo),
		getMainGoTemplate(&providerInfo),
		getGoModTemplate(&providerInfo),
		getSharedGoTemplate(),
		getOasJsonTemplate(apiSpec),
		getProviderGoTemplate(&providerInfo, resources, dataSources),
	}
	for _, resource := range resources {
		templates = append(templates, getResourceGoTemplate(&providerInfo, &resource))
	}

	// Write out files
	for _, renderer := range templates {
		completePath := path.Join(outputPath, renderer.Name())
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
