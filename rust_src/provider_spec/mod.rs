use std::fs;

pub(crate) mod rest_api_provider_schema;

pub fn load_provider_spec(
    path: &std::path::Path,
) -> anyhow::Result<rest_api_provider_schema::RestApiProviderConfiguration> {
    Ok(serde_saphyr::from_str(&fs::read_to_string(path)?)?)
}
