use std::fs;

mod parser;
mod rest_api_provider_schema;

fn load_provider_spec(
    path: &std::path::Path,
) -> anyhow::Result<rest_api_provider_schema::RestApiProviderConfiguration> {
    Ok(serde_saphyr::from_str(&fs::read_to_string(path)?)?)
}

fn main() {
    let api_spec = parser::OASVersionedSpec::parse_from_file(
        std::path::Path::new("./example/openapi.json"),
        parser::OASVersion::V3_0,
    )
    .unwrap();
    let provider_spec = load_provider_spec(std::path::Path::new("example/genspec.yaml")).unwrap();
    //println!("{:?}", api_spec);
    println!("{:?}", provider_spec);
}
