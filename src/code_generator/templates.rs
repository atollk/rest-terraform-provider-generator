use askama::Template;

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
pub(crate) struct DatSourceInfo {
    pub(crate) name_snake: String,
    pub(crate) name_pascal: String,
}

#[derive(Template)]
#[template(path = "Makefile")]
pub(crate) struct MakefileTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "go.mod")]
pub(crate) struct GoModTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "main.go")]
pub(crate) struct MainGoTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
}

#[derive(Template)]
#[template(path = "internal/provider/provider.go")]
pub(crate) struct ProviderGoTemplate<'a> {
    pub(crate) provider_info: &'a ProviderInfo,
    pub(crate) resources: &'a [ResourceInfo],
    pub(crate) data_sources: &'a [DatSourceInfo],
}
