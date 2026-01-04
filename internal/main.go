// Package main is the entry point for the Terraform API provider generator.
// It parses OpenAPI and provider specifications and generates Terraform provider code.
package main

import (
	"atollk/terraform-api-provider-generator/internal/code_generator"
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"fmt"
	"os"
)

func main() {
	oadoc, err := oas_parser.Parse("example/openapi.json")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	providerSpec, err := provider_spec.ParseSpecFromFile("example/genspec.yaml")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	err = code_generator.RenderSpec("example/out", providerSpec, oadoc)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
