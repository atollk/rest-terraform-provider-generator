package oas_parser

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/pb33f/libopenapi"
	_ "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
)

type OADoc = *libopenapi.DocumentModel[v3.Document]

func Parse(path string) (OADoc, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Errorf("error while reading file: %w", err)
	}

	document, err := libopenapi.NewDocument(bytes)
	if err != nil {
		return nil, errors.Errorf("error while parsing OAS: %w", err)
	}

	model, err := document.BuildV3Model()
	if err != nil {
		panic(errors.Errorf("cannot create v3 model from document: %w", err))
	}

	return model, nil
}
