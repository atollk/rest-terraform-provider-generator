use askama::Template;

#[derive(Template)]
#[template(path = "Makefile")]
pub(crate) struct MakefileTemplate<'a> {
    pub(crate) provider_author: &'a str,
    pub(crate) provider_name: &'a str,
}

#[derive(Template)]
#[template(path = "go.mod")]
pub(crate) struct GoModTemplate<'a> {
    pub(crate) provider_author: &'a str,
    pub(crate) provider_name: &'a str,
}

#[derive(Template)]
#[template(path = "main.go")]
pub(crate) struct MainGoTemplate<'a> {
    pub(crate) provider_author: &'a str,
    pub(crate) provider_name: &'a str,
}

#[derive(Template)]
#[template(path = "internal/provider/provider.go")]
pub(crate) struct ProviderGoTemplate<'a> {
    pub(crate) provider_author: &'a str,
    pub(crate) provider_name: &'a str,
    pub(crate) provider_name_caps: &'a str,
}
