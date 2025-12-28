package code_generator

import (
	"bytes"
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

type templateRenderer struct {
	name   string
	render func() ([]byte, error)
}

//go:embed templates/main/_Makefile
var makefileTemplate string

func getMakefileTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("Makefile", makefileTemplate, providerInfo)
}

//go:embed templates/main/_main.go
var mainGoTemplate string

func getMainGoTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("main.go", mainGoTemplate, providerInfo)
}

//go:embed templates/main/_go.mod
var goModTemplate string

func getGoModTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("go.mod", goModTemplate, providerInfo)
}

//go:embed templates/main/internal/provider/_provider.go
var providerGoTemplate string

func getProviderGoTemplate(providerInfo *ProviderInfo, resources []ResourceInfo, dataSources []DataSourceInfo) templateRenderer {
	return renderTemplateAs("internal/provider/provider.go", providerGoTemplate, struct {
		ProviderInfo *ProviderInfo
		Resources    []ResourceInfo
		DataSources  []DataSourceInfo
	}{
		providerInfo, resources, dataSources,
	})
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

func renderTemplateAs(templateName string, templateData string, args any) templateRenderer {
	renderFunc := func() ([]byte, error) {
		templ, err := template.New(templateName).Parse(templateData)
		if err != nil {
			return nil, err
		}
		buffer := bytes.Buffer{}
		err = templ.Execute(&buffer, args)
		if err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}
	return templateRenderer{name: templateName, render: renderFunc}
}
