package oas_parser

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/pb33f/libopenapi"
	_ "github.com/pb33f/libopenapi-validator"
	"github.com/pb33f/libopenapi/datamodel/high/v3"
)

// OADoc is a type alias for an OpenAPI v3 document model.
type OADoc = *libopenapi.DocumentModel[v3.Document]

// Parse reads and parses an OpenAPI specification file from the given path.
// It returns a v3 document model or an error if parsing fails.
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
