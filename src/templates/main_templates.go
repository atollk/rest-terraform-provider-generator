package templates

import (
	_ "embed"
	"text/template"
)

type ProviderInfo struct {
	Author    string
	NameKebab string
	NameCaps  string
}

type ResourceInfo struct {
	NameSnake  string
	NamePascal string
}

type DataSourceInfo struct {
	NameSnake  string
	NamePascal string
}

// -------------------------------------------------------------------------------------------------

//go:embed templates/main/_Makefile
var makefileTemplate string

func getMakefileTemplate() (*template.Template, error) {
	return template.New("Makefile").Parse(makefileTemplate)
}

/*
	struct MakefileTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct GoModTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct MainGoTemplate<'a> {
		provider_info &'a ProviderInfo,
	}

	struct ProviderGoTemplate<'a> {
		provider_info &'a ProviderInfo,
		resources &'a [ResourceInfo],
		data_sources &'a [DataSourceInfo],
	}

	struct ResourceGoTemplate<'a> {
		resource_info &'a ResourceInfo,
	}

	impl
crate::code_generator::helper_templates::ResourceGoTemplate < '_> {
	fn
	get_attributes_definition(&self)- > Vec < (string, GoCode)> {
	Vec::new()
	}
}

*/
