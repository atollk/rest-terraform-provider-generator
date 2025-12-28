package oas_parser

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	_ "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
)

type OADoc = *libopenapi.DocumentModel[v3.Document]

func Parse(path string) (OADoc, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	document, err := libopenapi.NewDocument(bytes)
	if err != nil {
		return nil, fmt.Errorf("error while parsing OAS: %w", err)
	}

	model, err := document.BuildV3Model()
	if err != nil {
		panic(fmt.Errorf("cannot create v3 model from document: %w", err))
	}

	return model, nil
}
