use crate::oas_parser::parser::OASVersionedSpec;
use crate::provider_spec::rest_api_provider_schema;
use std::path::Path;

mod templates;

pub fn render_spec(
    output_path: &Path,
    provider_spec: &rest_api_provider_schema::RestApiProviderConfiguration,
    api_spec: &OASVersionedSpec,
) -> anyhow::Result<()> {
    let makefile_template = templates::MakefileTemplate {
        provider_name: "petstore",
        provider_author: "foo",
    };
    let main_go_template = templates::MainGoTemplate {
        provider_name: "petstore",
        provider_author: "foo",
    };
    let go_mod_template = templates::GoModTemplate {
        provider_name: "petstore",
        provider_author: "foo",
    };
    let provider_go_template = templates::ProviderGoTemplate {
        provider_name: "petstore",
        provider_author: "foo",
        provider_name_caps: "PETSTORE",
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
