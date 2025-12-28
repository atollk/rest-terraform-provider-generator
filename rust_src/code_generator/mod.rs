use crate::code_generator::main_templates::ResourceInfo;
use crate::oas_parser::parser::OASVersionedSpec;
use crate::provider_spec::rest_api_provider_schema;
use std::path::Path;

mod helper_templates;
mod main_templates;

pub(crate) struct GoCode(String);

pub fn render_spec(
    output_path: &Path,
    provider_spec: &rest_api_provider_schema::RestApiProviderConfiguration,
    api_spec: &OASVersionedSpec,
) -> anyhow::Result<()> {
    // Prepare input data
    let provider_info = main_templates::ProviderInfo {
        name_kebab: "petstore".to_string(),
        author: "foo".to_string(),
        name_caps: "PETSTORE".to_string(),
    };
    let x = if let OASVersionedSpec::V3_0(x) = api_spec {
        let a = x.paths.iter().next().unwrap().1;
        a.
    } else { panic!() };
    let resources = vec![ResourceInfo {
        name_snake: "res_one".to_string(),
        name_pascal: "ResOne".to_string(),
    }];
    let data_sources = vec![];

    // Set up template objects
    let makefile_template = main_templates::MakefileTemplate {
        provider_info: &provider_info,
    };
    let main_go_template = main_templates::MainGoTemplate {
        provider_info: &provider_info,
    };
    let go_mod_template = main_templates::GoModTemplate {
        provider_info: &provider_info,
    };
    let provider_go_template = main_templates::ProviderGoTemplate {
        provider_info: &provider_info,
        resources: &resources,
        data_sources: &data_sources,
    };
    let resource_templates: Vec<_> = resources
        .iter()
        .map(|resource_info| main_templates::ResourceGoTemplate { resource_info })
        .collect();

    // Map output file names to templates
    let template_files = {
        let mut template_files = maplit::hashmap! {
            "Makefile".to_string() => &makefile_template as &dyn askama::DynTemplate,
            "go.mod".to_string() => &go_mod_template,
            "main.go".to_string() => &main_go_template,
            "internal/provider/provider.go".to_string() => &provider_go_template,
        };
        for resource_template in resource_templates.iter() {
            template_files.insert(
                format!(
                    "internal/provider/{}.go",
                    resource_template.resource_info.name_snake
                ),
                resource_template,
            );
        }
        template_files
    };

    // Write out files
    for (path, template) in template_files {
        let complete_path = output_path.join(Path::new(&path));
        if let Some(parent_path) = complete_path.parent() {
            std::fs::create_dir_all(parent_path)?;
        }
        let content = template.dyn_render()?;
        std::fs::write(complete_path, content)?;
    }
    Ok(())
}
