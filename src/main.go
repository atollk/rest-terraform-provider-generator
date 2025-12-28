package main

import (
	"atollk/terraform-api-provider-generator/src/oas_parser"
	"fmt"
	"log"
)

func main() {
	oadoc, err := oas_parser.Parse("example/openapi.json")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v", oadoc)
}
