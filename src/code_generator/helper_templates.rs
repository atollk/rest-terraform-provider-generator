use askama::Template;
use crate::code_generator::GoCode;

#[derive(Template)]
#[template(path = "helper/attribute_definition.go")]
pub(crate) struct AttributeDefinitionGoTemplate<'a> {
    name: &'a str,
    markdown_description: &'a str,
    data_type: &'a str,
    is_required: bool,
    is_computed: bool,
    is_sensitive: bool,
    default_value: GoCode,
}

