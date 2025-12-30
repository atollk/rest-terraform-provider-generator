package code_generator

import (
	"atollk/terraform-api-provider-generator/internal/oas_parser"
	"atollk/terraform-api-provider-generator/internal/provider_spec"
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
	NameSnake    string
	NamePascal   string
	ResourceSpec provider_spec.ResourceSchema
	OADoc        oas_parser.OADoc
}

type DataSourceInfo struct {
	NameSnake  string
	NamePascal string
}

// -------------------------------------------------------------------------------------------------

type templateRenderer interface {
	Name() string
	Render() ([]byte, error)
}

type basicTemplateRenderer struct {
	name   string
	render func() ([]byte, error)
}

func (r *basicTemplateRenderer) Name() string {
	return r.name
}

func (r *basicTemplateRenderer) Render() ([]byte, error) {
	return r.render()
}

//go:embed templates/main/Makefile.tmpl
var makefileTemplate string

func getMakefileTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("Makefile", makefileTemplate, providerInfo)
}

//go:embed templates/main/main.go.tmpl
var mainGoTemplate string

func getMainGoTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("main.go", mainGoTemplate, providerInfo)
}

//go:embed templates/main/go.mod.tmpl
var goModTemplate string

func getGoModTemplate(providerInfo *ProviderInfo) templateRenderer {
	return renderTemplateAs("go.mod", goModTemplate, providerInfo)
}

//go:embed templates/main/internal/provider/shared.go.tmpl
var sharedGoTemplate string

func getSharedGoTemplate() templateRenderer {
	return renderTemplateAs("internal/provider/shared.go", sharedGoTemplate, nil)
}

func getOasJsonTemplate(oadoc oas_parser.OADoc) templateRenderer {
	renderFunc := func() ([]byte, error) {
		return oadoc.Model.RenderJSON(" ")
	}
	return &basicTemplateRenderer{name: "internal/provider/oas.json", render: renderFunc}
}

//go:embed templates/main/internal/provider/provider.go.tmpl
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
	return &basicTemplateRenderer{name: templateName, render: renderFunc}
}
