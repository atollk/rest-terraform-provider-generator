use crate::oas_parser::parser::OASVersionedSpec;
use crate::provider_spec::rest_api_provider_schema;
use std::path::Path;
use crate::code_generator::templates::ResourceInfo;

mod templates;

pub fn render_spec(
    output_path: &Path,
    provider_spec: &rest_api_provider_schema::RestApiProviderConfiguration,
    api_spec: &OASVersionedSpec,
) -> anyhow::Result<()> {
    let provider_info = templates::ProviderInfo {
        name_kebab: "petstore".to_string(),
        author: "foo".to_string(),
        name_caps: "PETSTORE".to_string(),
    };
    let resources = vec![
        ResourceInfo {
            name_snake: "res_one".to_string(),
            name_pascal: "ResOne".to_string(),
        }
    ];
    let data_sources = vec![];

    let makefile_template = templates::MakefileTemplate {
        provider_info: &provider_info,
    };
    let main_go_template = templates::MainGoTemplate {
        provider_info: &provider_info,
    };
    let go_mod_template = templates::GoModTemplate {
        provider_info: &provider_info,
    };
    let provider_go_template = templates::ProviderGoTemplate {
        provider_info: &provider_info,
        resources: &resources,
        data_sources: &data_sources,
    };
    let template_files = maplit::hashmap! {
        "Makefile" => &makefile_template as &dyn askama::DynTemplate,
        "go.mod" => &go_mod_template,
        "main.go" => &main_go_template,
        "internal/provider/provider.go" => &provider_go_template,
    };
    for (path, template) in template_files {
        let complete_path = output_path.join(Path::new(path));
        if let Some(parent_path) = complete_path.parent() {
            std::fs::create_dir_all(parent_path)?;
        }
        let content = template.dyn_render()?;
        std::fs::write(complete_path, content)?;
    }
    Ok(())
}
