use anyhow::{anyhow, bail};
use std::fs;

#[derive(Copy, Clone, Debug)]
pub enum OASVersion {
    V2_0,
    V3_0,
    V3_1,
}

#[derive(Copy, Clone, Debug)]
pub enum OASFormat {
    YAML,
    JSON,
}

#[derive(Debug)]
pub enum OASVersionedSpec {
    V2_0(roas::v2::spec::Spec),
    V3_0(roas::v3_0::spec::Spec),
    V3_1(roas::v3_1::spec::Spec),
}

impl OASVersionedSpec {
    pub fn parse(
        spec_content: &str,
        spec_format: OASFormat,
        oas_version: OASVersion,
    ) -> anyhow::Result<OASVersionedSpec> {
        match oas_version {
            OASVersion::V2_0 => Ok(OASVersionedSpec::V2_0(match spec_format {
                OASFormat::YAML => serde_saphyr::from_str(&spec_content)?,
                OASFormat::JSON => serde_json::from_str(&spec_content)?,
            })),
            OASVersion::V3_0 => Ok(OASVersionedSpec::V3_0(match spec_format {
                OASFormat::YAML => serde_saphyr::from_str(&spec_content)?,
                OASFormat::JSON => serde_json::from_str(&spec_content)?,
            })),
            OASVersion::V3_1 => Ok(OASVersionedSpec::V3_1(match spec_format {
                OASFormat::YAML => serde_saphyr::from_str(&spec_content)?,
                OASFormat::JSON => serde_json::from_str(&spec_content)?,
            })),
        }
    }

    pub fn parse_from_file(
        path: &std::path::Path,
        oas_version: OASVersion,
    ) -> anyhow::Result<OASVersionedSpec> {
        let format = if let Some(os_str) = path.extension() {
            let extension = os_str
                .to_str()
                .ok_or(anyhow!(
                    "Path does not have a valid file extension: {:?}",
                    path
                ))?
                .to_lowercase();
            match extension.as_str() {
                "yaml" => OASFormat::YAML,
                "json" => OASFormat::JSON,
                _ => bail!("File extension is not supported: {}", extension),
            }
        } else {
            bail!("Path does not have a valid file extension: {:?}", path)
        };
        OASVersionedSpec::parse(&fs::read_to_string(path)?, format, oas_version)
    }
}
