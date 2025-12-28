//go:build tools
// +build tools

package tools

import (
	_ "github.com/RyoJerryYu/go-jsonschema/cmd/jsonschemagen"
	_ "github.com/atombender/go-jsonschema"
	_ "github.com/kaptinlin/jsonschema/cmd/schemagen"
)
