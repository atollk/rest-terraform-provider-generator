use std::path::Path;

mod code_generator;
mod oas_parser;
mod provider_spec;

fn main() {
    let api_spec = oas_parser::parser::OASVersionedSpec::parse_from_file(
        std::path::Path::new("./example/openapi.json"),
        oas_parser::parser::OASVersion::V3_0,
    )
    .unwrap();
    let provider_spec =
        provider_spec::load_provider_spec(std::path::Path::new("example/genspec.yaml")).unwrap();

    code_generator::render_spec(Path::new("example/out"), &provider_spec, &api_spec).unwrap();
}
