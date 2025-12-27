use askama::Template;
use clap::builder::Str;
use crate::code_generator::GoCode;

#[derive(Clone)]
pub(crate) struct ProviderInfo {
    pub(crate) author: String,
    pub(crate) name_kebab: String,
    pub(crate) name_caps: String,
}

#[derive(Clone)]
pub(crate) struct ResourceInfo {
    pub(crate) name_snake: String,
    pub(crate) name_pascal: String,
}

#[derive(Clone)]
pub(crate) struct DataSourceInfo {
    pub(crate) name_snake: String,
    pub(crate) name_pascal: String,
}

// -------------------------------------------------------------------------------------------------


#[derive(Template)]
#[template(path = "main/Makefile")]
pub(crate) struct MakefileTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "main/go.mod")]
pub(crate) struct GoModTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "main/main.go")]
pub(crate) struct MainGoTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "main/internal/provider/provider.go")]
pub(crate) struct ProviderGoTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
    pub(crate) resources: &'a [ResourceInfo],
    pub(crate) data_sources: &'a [DataSourceInfo],
}

#[derive(Template)]
#[template(path = "main/internal/provider/resource.go")]
pub(crate) struct ResourceGoTemplate<'a> {
    pub(crate) resource_info: &'a ResourceInfo,
}

impl crate::code_generator::helper_templates::ResourceGoTemplate<'_> {
    fn get_attributes_definition(&self) -> Vec<(String, GoCode)> {
        Vec::new()
    }
}