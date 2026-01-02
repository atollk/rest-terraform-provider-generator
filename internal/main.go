package main

import (
	"atollk/terraform-api-provider-generator/internal/code_generator"
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
	"fmt"
	"log"
	"os"
)

func main() {
	oadoc, err := oas_parser.Parse("example/openapi.json")
	if err != nil {
		log.Panic(err)
	}
	providerSpec, err := provider_spec.ParseSpecFromFile("example/genspec.yaml")
	if err != nil {
		log.Panic(err)
	}
	err = code_generator.RenderSpec("example/out", providerSpec, oadoc)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
