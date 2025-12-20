use anyhow::{anyhow, bail};
use std::fs;

fn main() {
    println!("Hello, world!");
}

enum OASVersion {
    V2_0,
    V3_0,
    V3_1,
}

enum OASVersionedSpec {
    V2_0(roas::v2::spec::Spec),
    V3_0(roas::v3_0::spec::Spec),
    V3_1(roas::v3_1::spec::Spec),
}

fn parse_openapi_from_str<T: serde_core::de::DeserializeOwned>(
    path: &std::path::Path,
) -> anyhow::Result<T> {
    if let Some(os_str) = path.extension() {
        let extension = os_str
            .to_str()
            .ok_or(anyhow!(
                "Path does not have a valid file extension: {:?}",
                path
            ))?
            .to_lowercase();
        let file_contents = fs::read_to_string(path)?;
        let result: T = match extension.as_str() {
            "yaml" => serde_saphyr::from_str(&file_contents)?,
            "json" => serde_json::from_str(&file_contents)?,
            _ => bail!("File extension is not supported: {}", extension),
        };
        Ok(result)
    } else {
        bail!("Path does not have a valid file extension: {:?}", path)
    }
}

fn parse_openapi(
    path: &std::path::Path,
    oas_version: OASVersion,
) -> anyhow::Result<OASVersionedSpec> {
    match oas_version {
        OASVersion::V2_0 => {
            let spec: roas::v2::spec::Spec = parse_openapi_from_str(path)?;
            Ok(OASVersionedSpec::V2_0(spec))
        }
        OASVersion::V3_0 => {
            let spec: roas::v3_0::spec::Spec = parse_openapi_from_str(path)?;
            Ok(OASVersionedSpec::V3_0(spec))
        }
        OASVersion::V3_1 => {
            let spec: roas::v3_1::spec::Spec = parse_openapi_from_str(path)?;
            Ok(OASVersionedSpec::V3_1(spec))
        }
    }
}
