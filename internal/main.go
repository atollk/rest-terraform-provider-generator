// Package main is the entry point for the Terraform API provider generator.
// It parses OpenAPI and provider specifications and generates Terraform provider code.
package main

import (
	"atollk/terraform-api-provider-generator/internal/code_generator"
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"errors"
	"flag"
	"fmt"
	"os"
)

type cliConfig struct {
	openApiPath         string
	providerSpecPath    string
	outputDirectoryPath string
}

func parseCliConfig() (cliConfig, error) {
	var openApiPath string
	var providerSpecPath string
	var outputDirectoryPath string
	flag.StringVar(&openApiPath, "open_api", "", "path to the OpenAPI JSON or YAML file")
	flag.StringVar(&providerSpecPath, "provider_spec", "", "path to the YAML file configuring the provider generation (see rest_api_provider_schema.json for schema)")
	flag.StringVar(&outputDirectoryPath, "output_directory", "", "path to a directory where all output files are written to")
	flag.Parse()
	var errs []error
	if openApiPath == "" {
		errs = append(errs, errors.New("open_api needs to be given"))
	}
	if providerSpecPath == "" {
		errs = append(errs, errors.New("provider_spec needs to be given"))
	}
	if outputDirectoryPath == "" {
		errs = append(errs, errors.New("output_directory needs to be given"))
	}
	if errs != nil {
		return cliConfig{}, errors.Join(errs...)
	}
	return cliConfig{
		openApiPath:         openApiPath,
		providerSpecPath:    providerSpecPath,
		outputDirectoryPath: outputDirectoryPath,
	}, nil
}

func main() {
	config, err := parseCliConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	oadoc, err := oas_parser.Parse(config.openApiPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	providerSpec, err := provider_spec.ParseSpecFromFile(config.providerSpecPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	err = code_generator.RenderSpec(config.outputDirectoryPath, providerSpec, oadoc)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
