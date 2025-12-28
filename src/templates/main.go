package templates

import (
	"atollk/terraform-api-provider-generator/src/oas_parser"
	"atollk/terraform-api-provider-generator/src/provider_spec"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"
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
	/*
		resources := []ResourceInfo{
			{nameSnake: "res_one", namePascal: "ResOne"},
		}
		var dataSources []DataSourceInfo
	*/

	// Set up template objects
	makefileTemplate, err := getMakefileTemplate()
	if err != nil {
		return fmt.Errorf("error while loading Makefile template: %w", err)
	}

	// Map output file names to templates
	type TemplateArgPair = struct {
		*template.Template
		any
	}
	templates := make(map[string]TemplateArgPair)
	templates["Makefile"] = TemplateArgPair{makefileTemplate, &providerInfo}

	// Write out files
	for filename, tmpl := range templates {
		completePath := path.Join(output_path, filename)
		err := os.MkdirAll(path.Dir(completePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create directories: %w", err)
		}
		file, err := os.Create(completePath)
		if err != nil {
			return fmt.Errorf("cannot create file: %w", err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Panicf("cannot close file: %v", err)
			}
		}()
		err = tmpl.Template.Execute(file, tmpl.any)
		if err != nil {
			return fmt.Errorf("cannot execute template: %w", err)
		}
	}
	return nil
}
