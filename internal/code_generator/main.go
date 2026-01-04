package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"os"
	"path"

	"github.com/cockroachdb/errors"
)

// RenderSpec generates Terraform provider code from the given provider and API specifications.
// It creates all necessary files in the outputPath directory including templates, providers, and resources.
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
	var resources []ResourceInfo
	var dataSources []DataSourceInfo
	for name, spec := range providerSpec.Resources.OtherProps {
		resource := ResourceInfo{
			name:         name,
			resourceSpec: spec,
			oadoc:        apiSpec,
			providerInfo: &providerInfo,
		}
		dataSource := DataSourceInfo{
			name:         name,
			resourceSpec: spec,
			oadoc:        apiSpec,
			providerInfo: &providerInfo,
		}
		if spec.GenerateResource == nil || *spec.GenerateResource {
			resources = append(resources, resource)
		}
		if spec.GenerateDataSource == nil || *spec.GenerateDataSource {
			dataSources = append(dataSources, dataSource)
		}
	}

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
		templates = append(templates, getResourceGoTemplate(&providerInfo, &resource, false))
	}
	for _, dataSource := range dataSources {
		templates = append(templates, getResourceGoTemplate(&providerInfo, &dataSource, true))
	}

	// Write out files
	for _, renderer := range templates {
		completePath := path.Join(outputPath, renderer.Name())
		err := os.MkdirAll(path.Dir(completePath), os.ModePerm)
		if err != nil {
			return errors.Errorf("cannot create directories: %w", err)
		}
		output, err := renderer.Render()
		if err != nil {
			return errors.Errorf("cannot execute template: %w", err)
		}
		err = os.WriteFile(completePath, output, os.ModePerm)
		if err != nil {
			return errors.Errorf("cannot write file: %w", err)
		}
	}

	return nil
}
