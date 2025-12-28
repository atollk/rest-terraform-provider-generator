package main

import (
	"atollk/terraform-api-provider-generator/src/oas_parser"
	"atollk/terraform-api-provider-generator/src/provider_spec"
	"fmt"
	"log"
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
	fmt.Printf("oadoc: %v\n", oadoc)
	fmt.Printf("providerSpec: %v\n", providerSpec)
}
