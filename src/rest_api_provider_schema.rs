#![allow(clippy::redundant_closure_call)]
#![allow(clippy::needless_lifetimes)]
#![allow(clippy::match_single_binding)]
#![allow(clippy::clone_on_copy)]

#[doc = r" Error types."]
pub mod error {
    #[doc = r" Error from a `TryFrom` or `FromStr` implementation."]
    pub struct ConversionError(::std::borrow::Cow<'static, str>);
    impl ::std::error::Error for ConversionError {}
    impl ::std::fmt::Display for ConversionError {
        fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> Result<(), ::std::fmt::Error> {
            ::std::fmt::Display::fmt(&self.0, f)
        }
    }
    impl ::std::fmt::Debug for ConversionError {
        fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> Result<(), ::std::fmt::Error> {
            ::std::fmt::Debug::fmt(&self.0, f)
        }
    }
    impl From<&'static str> for ConversionError {
        fn from(value: &'static str) -> Self {
            Self(value.into())
        }
    }
    impl From<String> for ConversionError {
        fn from(value: String) -> Self {
            Self(value.into())
        }
    }
}
#[doc = "Configuration schema for REST API provider"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"title\": \"REST API Provider Configuration\","]
#[doc = "  \"description\": \"Configuration schema for REST API provider\","]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"global\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"cert_file\": {"]
#[doc = "          \"description\": \"When set with the key_file parameter, the provider will load a client certificate as a file for mTLS authentication.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"cert_string\": {"]
#[doc = "          \"description\": \"When set with the key_string parameter, the provider will load a client certificate as a string for mTLS authentication.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"copy_keys\": {"]
#[doc = "          \"description\": \"When set, any PUT to the API for an object will copy these keys from the data the provider has gathered about the object. This is useful if internal API information must also be provided with updates, such as the revision of the object.\","]
#[doc = "          \"type\": \"array\","]
#[doc = "          \"items\": {"]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"create_method\": {"]
#[doc = "          \"description\": \"Defaults to POST. The HTTP method used to CREATE objects of this type on the API server.\","]
#[doc = "          \"default\": \"POST\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"create_returns_object\": {"]
#[doc = "          \"description\": \"Set this when the API returns the object created only on creation operations (POST). This is used by the provider to refresh internal data structures.\","]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"debug\": {"]
#[doc = "          \"description\": \"Enabling this will cause lots of debug information to be printed to STDOUT by the API client.\","]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"destroy_method\": {"]
#[doc = "          \"description\": \"Defaults to DELETE. The HTTP method used to DELETE objects of this type on the API server.\","]
#[doc = "          \"default\": \"DELETE\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"headers\": {"]
#[doc = "          \"description\": \"A map of header names and values to set on all outbound requests. This is useful if you want to use a script via the 'external' provider or provide a pre-approved token or change Content-Type from application/json. If username and password are set and Authorization is one of the headers defined here, the BASIC auth credentials take precedence.\","]
#[doc = "          \"type\": \"object\","]
#[doc = "          \"additionalProperties\": {"]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"id_attribute\": {"]
#[doc = "          \"description\": \"When set, this key will be used to operate on REST objects. For example, if the ID is set to 'name', changes to the API object will be to http://foo.com/bar/VALUE_OF_NAME. This value may also be a '/'-delimited path to the id attribute if it is multiple levels deep in the data (such as attributes/id in the case of an object { \\\"attributes\\\": { \\\"id\\\": 1234 }, \\\"config\\\": { \\\"name\\\": \\\"foo\\\", \\\"something\\\": \\\"bar\\\" } }\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"insecure\": {"]
#[doc = "          \"description\": \"When using https, this disables TLS verification of the host.\","]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"key_file\": {"]
#[doc = "          \"description\": \"When set with the cert_file parameter, the provider will load a client certificate as a file for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"key_string\": {"]
#[doc = "          \"description\": \"When set with the cert_string parameter, the provider will load a client certificate as a string for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"oauth_client_credentials\": {"]
#[doc = "          \"description\": \"Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation\","]
#[doc = "          \"type\": \"array\","]
#[doc = "          \"items\": {"]
#[doc = "            \"description\": \"Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation\","]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"properties\": {"]
#[doc = "              \"endpoint_params\": {"]
#[doc = "                \"description\": \"Additional key/values to pass to the underlying Oauth client library (as EndpointParams)\","]
#[doc = "                \"type\": \"object\","]
#[doc = "                \"additionalProperties\": true"]
#[doc = "              },"]
#[doc = "              \"oauth_client_id\": {"]
#[doc = "                \"description\": \"The OAuth client ID\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"oauth_client_secret\": {"]
#[doc = "                \"description\": \"The OAuth client secret\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"oauth_scopes\": {"]
#[doc = "                \"description\": \"OAuth scopes to request\","]
#[doc = "                \"type\": \"array\","]
#[doc = "                \"items\": {"]
#[doc = "                  \"type\": \"string\""]
#[doc = "                }"]
#[doc = "              },"]
#[doc = "              \"oauth_token_endpoint\": {"]
#[doc = "                \"description\": \"The OAuth token endpoint URL\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              }"]
#[doc = "            }"]
#[doc = "          },"]
#[doc = "          \"maxItems\": 1"]
#[doc = "        },"]
#[doc = "        \"password\": {"]
#[doc = "          \"description\": \"When set, will use this password for BASIC auth to the API.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"rate_limit\": {"]
#[doc = "          \"description\": \"Set this to limit the number of requests per second made to the API.\","]
#[doc = "          \"type\": \"number\""]
#[doc = "        },"]
#[doc = "        \"read_method\": {"]
#[doc = "          \"description\": \"Defaults to GET. The HTTP method used to READ objects of this type on the API server.\","]
#[doc = "          \"default\": \"GET\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"root_ca_file\": {"]
#[doc = "          \"description\": \"When set, the provider will load a root CA certificate as a file for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"root_ca_string\": {"]
#[doc = "          \"description\": \"When set, the provider will load a root CA certificate as a string for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"test_path\": {"]
#[doc = "          \"description\": \"If set, the provider will issue a read_method request to this path after instantiation requiring a 200 OK response before proceeding. This is useful if your API provides a no-op endpoint that can signal if this provider is configured correctly. Response data will be ignored.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"timeout\": {"]
#[doc = "          \"description\": \"When set, will cause requests taking longer than this time (in seconds) to be aborted.\","]
#[doc = "          \"type\": \"number\""]
#[doc = "        },"]
#[doc = "        \"update_method\": {"]
#[doc = "          \"description\": \"Defaults to PUT. The HTTP method used to UPDATE objects of this type on the API server.\","]
#[doc = "          \"default\": \"PUT\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"uri\": {"]
#[doc = "          \"description\": \"URI of the REST API endpoint. This serves as the base of all requests.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"use_cookies\": {"]
#[doc = "          \"description\": \"Enable cookie jar to persist session.\","]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"username\": {"]
#[doc = "          \"description\": \"When set, will use this username for BASIC auth to the API.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"write_returns_object\": {"]
#[doc = "          \"description\": \"Set this when the API returns the object created on all write operations (POST, PUT). This is used by the provider to refresh internal data structures.\","]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"xssi_prefix\": {"]
#[doc = "          \"description\": \"Trim the xssi prefix from response string, if present, before parsing.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    },"]
#[doc = "    \"resources\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"generate_data_source\": {"]
#[doc = "          \"description\": \"Defaults to true. Whether to generate a Terraform data source type for this API object.\","]
#[doc = "          \"default\": true,"]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        },"]
#[doc = "        \"generate_resource\": {"]
#[doc = "          \"description\": \"Defaults to true. Whether to generate a Terraform resource type for this API object.\","]
#[doc = "          \"default\": true,"]
#[doc = "          \"type\": \"boolean\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": {"]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"required\": ["]
#[doc = "          \"path\""]
#[doc = "        ],"]
#[doc = "        \"properties\": {"]
#[doc = "          \"create\": {"]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"properties\": {"]
#[doc = "              \"method\": {"]
#[doc = "                \"description\": \"Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"path\": {"]
#[doc = "                \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              }"]
#[doc = "            },"]
#[doc = "            \"additionalProperties\": false"]
#[doc = "          },"]
#[doc = "          \"debug\": {"]
#[doc = "            \"description\": \"Whether to emit verbose debug output while working with the API object on the server.\","]
#[doc = "            \"type\": \"boolean\""]
#[doc = "          },"]
#[doc = "          \"destroy\": {"]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"properties\": {"]
#[doc = "              \"method\": {"]
#[doc = "                \"description\": \"Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"path\": {"]
#[doc = "                \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              }"]
#[doc = "            },"]
#[doc = "            \"additionalProperties\": false"]
#[doc = "          },"]
#[doc = "          \"force_new\": {"]
#[doc = "            \"description\": \"Any changes to these values will result in recreating the resource instead of updating.\","]
#[doc = "            \"type\": \"array\","]
#[doc = "            \"items\": {"]
#[doc = "              \"type\": \"string\""]
#[doc = "            }"]
#[doc = "          },"]
#[doc = "          \"id_attribute\": {"]
#[doc = "            \"description\": \"Defaults to id_attribute set on the provider. Allows per-resource override of id_attribute (see id_attribute provider config documentation)\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"ignore_all_server_changes\": {"]
#[doc = "            \"description\": \"By default Terraform will attempt to revert changes to remote resources. Set this to 'true' to ignore any remote changes. Default: false\","]
#[doc = "            \"default\": false,"]
#[doc = "            \"type\": \"boolean\""]
#[doc = "          },"]
#[doc = "          \"ignore_changes_to\": {"]
#[doc = "            \"description\": \"A list of fields to which remote changes will be ignored. For example, an API might add or remove metadata, such as a 'last_modified' field, which Terraform should not attempt to correct. To ignore changes to nested fields, use the dot syntax: 'metadata.timestamp'\","]
#[doc = "            \"type\": \"array\","]
#[doc = "            \"items\": {"]
#[doc = "              \"type\": \"string\""]
#[doc = "            }"]
#[doc = "          },"]
#[doc = "          \"object_id\": {"]
#[doc = "            \"description\": \"Defaults to the id learned by the provider during normal operations and id_attribute. Allows you to set the id manually. This is used in conjunction with the *_path attributes.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"path\": {"]
#[doc = "            \"description\": \"The API path on top of the base URL set in the provider that represents objects of this type on the API server.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"query_string\": {"]
#[doc = "            \"description\": \"Query string to be included in the path\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"read\": {"]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"properties\": {"]
#[doc = "              \"method\": {"]
#[doc = "                \"description\": \"Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"path\": {"]
#[doc = "                \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"search\": {"]
#[doc = "                \"description\": \"Custom search for read_path.\","]
#[doc = "                \"type\": \"object\","]
#[doc = "                \"required\": ["]
#[doc = "                  \"search_key\","]
#[doc = "                  \"search_value\""]
#[doc = "                ],"]
#[doc = "                \"properties\": {"]
#[doc = "                  \"query_string\": {"]
#[doc = "                    \"description\": \"An optional query string to send when performing the search.\","]
#[doc = "                    \"type\": \"string\""]
#[doc = "                  },"]
#[doc = "                  \"results_key\": {"]
#[doc = "                    \"description\": \"When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.\","]
#[doc = "                    \"type\": \"string\""]
#[doc = "                  },"]
#[doc = "                  \"search_key\": {"]
#[doc = "                    \"description\": \"When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.\","]
#[doc = "                    \"type\": \"string\""]
#[doc = "                  },"]
#[doc = "                  \"search_path\": {"]
#[doc = "                    \"description\": \"The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.\","]
#[doc = "                    \"type\": \"string\""]
#[doc = "                  },"]
#[doc = "                  \"search_value\": {"]
#[doc = "                    \"description\": \"The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.\","]
#[doc = "                    \"type\": \"string\""]
#[doc = "                  }"]
#[doc = "                },"]
#[doc = "                \"additionalProperties\": false"]
#[doc = "              }"]
#[doc = "            },"]
#[doc = "            \"additionalProperties\": false"]
#[doc = "          },"]
#[doc = "          \"update\": {"]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"properties\": {"]
#[doc = "              \"method\": {"]
#[doc = "                \"description\": \"Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"path\": {"]
#[doc = "                \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              }"]
#[doc = "            },"]
#[doc = "            \"additionalProperties\": false"]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"additionalProperties\": false"]
#[doc = "      }"]
#[doc = "    }"]
#[doc = "  }"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
pub struct RestApiProviderConfiguration {
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub global: ::std::option::Option<RestApiProviderConfigurationGlobal>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub resources: ::std::option::Option<RestApiProviderConfigurationResources>,
}
impl ::std::convert::From<&RestApiProviderConfiguration> for RestApiProviderConfiguration {
    fn from(value: &RestApiProviderConfiguration) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfiguration {
    fn default() -> Self {
        Self {
            global: Default::default(),
            resources: Default::default(),
        }
    }
}
impl RestApiProviderConfiguration {
    pub fn builder() -> builder::RestApiProviderConfiguration {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationGlobal`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"cert_file\": {"]
#[doc = "      \"description\": \"When set with the key_file parameter, the provider will load a client certificate as a file for mTLS authentication.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"cert_string\": {"]
#[doc = "      \"description\": \"When set with the key_string parameter, the provider will load a client certificate as a string for mTLS authentication.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"copy_keys\": {"]
#[doc = "      \"description\": \"When set, any PUT to the API for an object will copy these keys from the data the provider has gathered about the object. This is useful if internal API information must also be provided with updates, such as the revision of the object.\","]
#[doc = "      \"type\": \"array\","]
#[doc = "      \"items\": {"]
#[doc = "        \"type\": \"string\""]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"create_method\": {"]
#[doc = "      \"description\": \"Defaults to POST. The HTTP method used to CREATE objects of this type on the API server.\","]
#[doc = "      \"default\": \"POST\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"create_returns_object\": {"]
#[doc = "      \"description\": \"Set this when the API returns the object created only on creation operations (POST). This is used by the provider to refresh internal data structures.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"debug\": {"]
#[doc = "      \"description\": \"Enabling this will cause lots of debug information to be printed to STDOUT by the API client.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"destroy_method\": {"]
#[doc = "      \"description\": \"Defaults to DELETE. The HTTP method used to DELETE objects of this type on the API server.\","]
#[doc = "      \"default\": \"DELETE\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"headers\": {"]
#[doc = "      \"description\": \"A map of header names and values to set on all outbound requests. This is useful if you want to use a script via the 'external' provider or provide a pre-approved token or change Content-Type from application/json. If username and password are set and Authorization is one of the headers defined here, the BASIC auth credentials take precedence.\","]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"additionalProperties\": {"]
#[doc = "        \"type\": \"string\""]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"id_attribute\": {"]
#[doc = "      \"description\": \"When set, this key will be used to operate on REST objects. For example, if the ID is set to 'name', changes to the API object will be to http://foo.com/bar/VALUE_OF_NAME. This value may also be a '/'-delimited path to the id attribute if it is multiple levels deep in the data (such as attributes/id in the case of an object { \\\"attributes\\\": { \\\"id\\\": 1234 }, \\\"config\\\": { \\\"name\\\": \\\"foo\\\", \\\"something\\\": \\\"bar\\\" } }\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"insecure\": {"]
#[doc = "      \"description\": \"When using https, this disables TLS verification of the host.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"key_file\": {"]
#[doc = "      \"description\": \"When set with the cert_file parameter, the provider will load a client certificate as a file for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"key_string\": {"]
#[doc = "      \"description\": \"When set with the cert_string parameter, the provider will load a client certificate as a string for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"oauth_client_credentials\": {"]
#[doc = "      \"description\": \"Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation\","]
#[doc = "      \"type\": \"array\","]
#[doc = "      \"items\": {"]
#[doc = "        \"description\": \"Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation\","]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"properties\": {"]
#[doc = "          \"endpoint_params\": {"]
#[doc = "            \"description\": \"Additional key/values to pass to the underlying Oauth client library (as EndpointParams)\","]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"additionalProperties\": true"]
#[doc = "          },"]
#[doc = "          \"oauth_client_id\": {"]
#[doc = "            \"description\": \"The OAuth client ID\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"oauth_client_secret\": {"]
#[doc = "            \"description\": \"The OAuth client secret\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"oauth_scopes\": {"]
#[doc = "            \"description\": \"OAuth scopes to request\","]
#[doc = "            \"type\": \"array\","]
#[doc = "            \"items\": {"]
#[doc = "              \"type\": \"string\""]
#[doc = "            }"]
#[doc = "          },"]
#[doc = "          \"oauth_token_endpoint\": {"]
#[doc = "            \"description\": \"The OAuth token endpoint URL\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"maxItems\": 1"]
#[doc = "    },"]
#[doc = "    \"password\": {"]
#[doc = "      \"description\": \"When set, will use this password for BASIC auth to the API.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"rate_limit\": {"]
#[doc = "      \"description\": \"Set this to limit the number of requests per second made to the API.\","]
#[doc = "      \"type\": \"number\""]
#[doc = "    },"]
#[doc = "    \"read_method\": {"]
#[doc = "      \"description\": \"Defaults to GET. The HTTP method used to READ objects of this type on the API server.\","]
#[doc = "      \"default\": \"GET\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"root_ca_file\": {"]
#[doc = "      \"description\": \"When set, the provider will load a root CA certificate as a file for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"root_ca_string\": {"]
#[doc = "      \"description\": \"When set, the provider will load a root CA certificate as a string for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"test_path\": {"]
#[doc = "      \"description\": \"If set, the provider will issue a read_method request to this path after instantiation requiring a 200 OK response before proceeding. This is useful if your API provides a no-op endpoint that can signal if this provider is configured correctly. Response data will be ignored.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"timeout\": {"]
#[doc = "      \"description\": \"When set, will cause requests taking longer than this time (in seconds) to be aborted.\","]
#[doc = "      \"type\": \"number\""]
#[doc = "    },"]
#[doc = "    \"update_method\": {"]
#[doc = "      \"description\": \"Defaults to PUT. The HTTP method used to UPDATE objects of this type on the API server.\","]
#[doc = "      \"default\": \"PUT\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"uri\": {"]
#[doc = "      \"description\": \"URI of the REST API endpoint. This serves as the base of all requests.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"use_cookies\": {"]
#[doc = "      \"description\": \"Enable cookie jar to persist session.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"username\": {"]
#[doc = "      \"description\": \"When set, will use this username for BASIC auth to the API.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"write_returns_object\": {"]
#[doc = "      \"description\": \"Set this when the API returns the object created on all write operations (POST, PUT). This is used by the provider to refresh internal data structures.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"xssi_prefix\": {"]
#[doc = "      \"description\": \"Trim the xssi prefix from response string, if present, before parsing.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationGlobal {
    #[doc = "When set with the key_file parameter, the provider will load a client certificate as a file for mTLS authentication."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub cert_file: ::std::option::Option<::std::string::String>,
    #[doc = "When set with the key_string parameter, the provider will load a client certificate as a string for mTLS authentication."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub cert_string: ::std::option::Option<::std::string::String>,
    #[doc = "When set, any PUT to the API for an object will copy these keys from the data the provider has gathered about the object. This is useful if internal API information must also be provided with updates, such as the revision of the object."]
    #[serde(default, skip_serializing_if = "::std::vec::Vec::is_empty")]
    pub copy_keys: ::std::vec::Vec<::std::string::String>,
    #[doc = "Defaults to POST. The HTTP method used to CREATE objects of this type on the API server."]
    #[serde(default = "defaults::rest_api_provider_configuration_global_create_method")]
    pub create_method: ::std::string::String,
    #[doc = "Set this when the API returns the object created only on creation operations (POST). This is used by the provider to refresh internal data structures."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub create_returns_object: ::std::option::Option<bool>,
    #[doc = "Enabling this will cause lots of debug information to be printed to STDOUT by the API client."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub debug: ::std::option::Option<bool>,
    #[doc = "Defaults to DELETE. The HTTP method used to DELETE objects of this type on the API server."]
    #[serde(default = "defaults::rest_api_provider_configuration_global_destroy_method")]
    pub destroy_method: ::std::string::String,
    #[doc = "A map of header names and values to set on all outbound requests. This is useful if you want to use a script via the 'external' provider or provide a pre-approved token or change Content-Type from application/json. If username and password are set and Authorization is one of the headers defined here, the BASIC auth credentials take precedence."]
    #[serde(
        default,
        skip_serializing_if = ":: std :: collections :: HashMap::is_empty"
    )]
    pub headers: ::std::collections::HashMap<::std::string::String, ::std::string::String>,
    #[doc = "When set, this key will be used to operate on REST objects. For example, if the ID is set to 'name', changes to the API object will be to http://foo.com/bar/VALUE_OF_NAME. This value may also be a '/'-delimited path to the id attribute if it is multiple levels deep in the data (such as attributes/id in the case of an object { \"attributes\": { \"id\": 1234 }, \"config\": { \"name\": \"foo\", \"something\": \"bar\" } }"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub id_attribute: ::std::option::Option<::std::string::String>,
    #[doc = "When using https, this disables TLS verification of the host."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub insecure: ::std::option::Option<bool>,
    #[doc = "When set with the cert_file parameter, the provider will load a client certificate as a file for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub key_file: ::std::option::Option<::std::string::String>,
    #[doc = "When set with the cert_string parameter, the provider will load a client certificate as a string for mTLS authentication. Note that this mechanism simply delegates to golang's tls.LoadX509KeyPair which does not support passphrase protected private keys. The most robust security protections available to the key_file are simple file system permissions."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub key_string: ::std::option::Option<::std::string::String>,
    #[doc = "Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation"]
    #[serde(default, skip_serializing_if = "::std::vec::Vec::is_empty")]
    pub oauth_client_credentials:
        ::std::vec::Vec<RestApiProviderConfigurationGlobalOauthClientCredentialsItem>,
    #[doc = "When set, will use this password for BASIC auth to the API."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub password: ::std::option::Option<::std::string::String>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub rate_limit: ::std::option::Option<f64>,
    #[doc = "Defaults to GET. The HTTP method used to READ objects of this type on the API server."]
    #[serde(default = "defaults::rest_api_provider_configuration_global_read_method")]
    pub read_method: ::std::string::String,
    #[doc = "When set, the provider will load a root CA certificate as a file for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub root_ca_file: ::std::option::Option<::std::string::String>,
    #[doc = "When set, the provider will load a root CA certificate as a string for mTLS authentication. This is useful when the API server is using a self-signed certificate and the client needs to trust it."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub root_ca_string: ::std::option::Option<::std::string::String>,
    #[doc = "If set, the provider will issue a read_method request to this path after instantiation requiring a 200 OK response before proceeding. This is useful if your API provides a no-op endpoint that can signal if this provider is configured correctly. Response data will be ignored."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub test_path: ::std::option::Option<::std::string::String>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub timeout: ::std::option::Option<f64>,
    #[doc = "Defaults to PUT. The HTTP method used to UPDATE objects of this type on the API server."]
    #[serde(default = "defaults::rest_api_provider_configuration_global_update_method")]
    pub update_method: ::std::string::String,
    #[doc = "URI of the REST API endpoint. This serves as the base of all requests."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub uri: ::std::option::Option<::std::string::String>,
    #[doc = "Enable cookie jar to persist session."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub use_cookies: ::std::option::Option<bool>,
    #[doc = "When set, will use this username for BASIC auth to the API."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub username: ::std::option::Option<::std::string::String>,
    #[doc = "Set this when the API returns the object created on all write operations (POST, PUT). This is used by the provider to refresh internal data structures."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub write_returns_object: ::std::option::Option<bool>,
    #[doc = "Trim the xssi prefix from response string, if present, before parsing."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub xssi_prefix: ::std::option::Option<::std::string::String>,
}
impl ::std::convert::From<&RestApiProviderConfigurationGlobal>
    for RestApiProviderConfigurationGlobal
{
    fn from(value: &RestApiProviderConfigurationGlobal) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationGlobal {
    fn default() -> Self {
        Self {
            cert_file: Default::default(),
            cert_string: Default::default(),
            copy_keys: Default::default(),
            create_method: defaults::rest_api_provider_configuration_global_create_method(),
            create_returns_object: Default::default(),
            debug: Default::default(),
            destroy_method: defaults::rest_api_provider_configuration_global_destroy_method(),
            headers: Default::default(),
            id_attribute: Default::default(),
            insecure: Default::default(),
            key_file: Default::default(),
            key_string: Default::default(),
            oauth_client_credentials: Default::default(),
            password: Default::default(),
            rate_limit: Default::default(),
            read_method: defaults::rest_api_provider_configuration_global_read_method(),
            root_ca_file: Default::default(),
            root_ca_string: Default::default(),
            test_path: Default::default(),
            timeout: Default::default(),
            update_method: defaults::rest_api_provider_configuration_global_update_method(),
            uri: Default::default(),
            use_cookies: Default::default(),
            username: Default::default(),
            write_returns_object: Default::default(),
            xssi_prefix: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationGlobal {
    pub fn builder() -> builder::RestApiProviderConfigurationGlobal {
        Default::default()
    }
}
#[doc = "Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"description\": \"Configuration for oauth client credential flow using the https://pkg.go.dev/golang.org/x/oauth2 implementation\","]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"endpoint_params\": {"]
#[doc = "      \"description\": \"Additional key/values to pass to the underlying Oauth client library (as EndpointParams)\","]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"additionalProperties\": true"]
#[doc = "    },"]
#[doc = "    \"oauth_client_id\": {"]
#[doc = "      \"description\": \"The OAuth client ID\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"oauth_client_secret\": {"]
#[doc = "      \"description\": \"The OAuth client secret\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"oauth_scopes\": {"]
#[doc = "      \"description\": \"OAuth scopes to request\","]
#[doc = "      \"type\": \"array\","]
#[doc = "      \"items\": {"]
#[doc = "        \"type\": \"string\""]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"oauth_token_endpoint\": {"]
#[doc = "      \"description\": \"The OAuth token endpoint URL\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  }"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
pub struct RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
    #[doc = "Additional key/values to pass to the underlying Oauth client library (as EndpointParams)"]
    #[serde(default, skip_serializing_if = "::serde_json::Map::is_empty")]
    pub endpoint_params: ::serde_json::Map<::std::string::String, ::serde_json::Value>,
    #[doc = "The OAuth client ID"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub oauth_client_id: ::std::option::Option<::std::string::String>,
    #[doc = "The OAuth client secret"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub oauth_client_secret: ::std::option::Option<::std::string::String>,
    #[doc = "OAuth scopes to request"]
    #[serde(default, skip_serializing_if = "::std::vec::Vec::is_empty")]
    pub oauth_scopes: ::std::vec::Vec<::std::string::String>,
    #[doc = "The OAuth token endpoint URL"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub oauth_token_endpoint: ::std::option::Option<::std::string::String>,
}
impl ::std::convert::From<&RestApiProviderConfigurationGlobalOauthClientCredentialsItem>
    for RestApiProviderConfigurationGlobalOauthClientCredentialsItem
{
    fn from(value: &RestApiProviderConfigurationGlobalOauthClientCredentialsItem) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
    fn default() -> Self {
        Self {
            endpoint_params: Default::default(),
            oauth_client_id: Default::default(),
            oauth_client_secret: Default::default(),
            oauth_scopes: Default::default(),
            oauth_token_endpoint: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
    pub fn builder() -> builder::RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResources`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"generate_data_source\": {"]
#[doc = "      \"description\": \"Defaults to true. Whether to generate a Terraform data source type for this API object.\","]
#[doc = "      \"default\": true,"]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"generate_resource\": {"]
#[doc = "      \"description\": \"Defaults to true. Whether to generate a Terraform resource type for this API object.\","]
#[doc = "      \"default\": true,"]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": {"]
#[doc = "    \"type\": \"object\","]
#[doc = "    \"required\": ["]
#[doc = "      \"path\""]
#[doc = "    ],"]
#[doc = "    \"properties\": {"]
#[doc = "      \"create\": {"]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"properties\": {"]
#[doc = "          \"method\": {"]
#[doc = "            \"description\": \"Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"path\": {"]
#[doc = "            \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"additionalProperties\": false"]
#[doc = "      },"]
#[doc = "      \"debug\": {"]
#[doc = "        \"description\": \"Whether to emit verbose debug output while working with the API object on the server.\","]
#[doc = "        \"type\": \"boolean\""]
#[doc = "      },"]
#[doc = "      \"destroy\": {"]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"properties\": {"]
#[doc = "          \"method\": {"]
#[doc = "            \"description\": \"Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"path\": {"]
#[doc = "            \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"additionalProperties\": false"]
#[doc = "      },"]
#[doc = "      \"force_new\": {"]
#[doc = "        \"description\": \"Any changes to these values will result in recreating the resource instead of updating.\","]
#[doc = "        \"type\": \"array\","]
#[doc = "        \"items\": {"]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"id_attribute\": {"]
#[doc = "        \"description\": \"Defaults to id_attribute set on the provider. Allows per-resource override of id_attribute (see id_attribute provider config documentation)\","]
#[doc = "        \"type\": \"string\""]
#[doc = "      },"]
#[doc = "      \"ignore_all_server_changes\": {"]
#[doc = "        \"description\": \"By default Terraform will attempt to revert changes to remote resources. Set this to 'true' to ignore any remote changes. Default: false\","]
#[doc = "        \"default\": false,"]
#[doc = "        \"type\": \"boolean\""]
#[doc = "      },"]
#[doc = "      \"ignore_changes_to\": {"]
#[doc = "        \"description\": \"A list of fields to which remote changes will be ignored. For example, an API might add or remove metadata, such as a 'last_modified' field, which Terraform should not attempt to correct. To ignore changes to nested fields, use the dot syntax: 'metadata.timestamp'\","]
#[doc = "        \"type\": \"array\","]
#[doc = "        \"items\": {"]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"object_id\": {"]
#[doc = "        \"description\": \"Defaults to the id learned by the provider during normal operations and id_attribute. Allows you to set the id manually. This is used in conjunction with the *_path attributes.\","]
#[doc = "        \"type\": \"string\""]
#[doc = "      },"]
#[doc = "      \"path\": {"]
#[doc = "        \"description\": \"The API path on top of the base URL set in the provider that represents objects of this type on the API server.\","]
#[doc = "        \"type\": \"string\""]
#[doc = "      },"]
#[doc = "      \"query_string\": {"]
#[doc = "        \"description\": \"Query string to be included in the path\","]
#[doc = "        \"type\": \"string\""]
#[doc = "      },"]
#[doc = "      \"read\": {"]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"properties\": {"]
#[doc = "          \"method\": {"]
#[doc = "            \"description\": \"Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"path\": {"]
#[doc = "            \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"search\": {"]
#[doc = "            \"description\": \"Custom search for read_path.\","]
#[doc = "            \"type\": \"object\","]
#[doc = "            \"required\": ["]
#[doc = "              \"search_key\","]
#[doc = "              \"search_value\""]
#[doc = "            ],"]
#[doc = "            \"properties\": {"]
#[doc = "              \"query_string\": {"]
#[doc = "                \"description\": \"An optional query string to send when performing the search.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"results_key\": {"]
#[doc = "                \"description\": \"When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"search_key\": {"]
#[doc = "                \"description\": \"When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"search_path\": {"]
#[doc = "                \"description\": \"The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              },"]
#[doc = "              \"search_value\": {"]
#[doc = "                \"description\": \"The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.\","]
#[doc = "                \"type\": \"string\""]
#[doc = "              }"]
#[doc = "            },"]
#[doc = "            \"additionalProperties\": false"]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"additionalProperties\": false"]
#[doc = "      },"]
#[doc = "      \"update\": {"]
#[doc = "        \"type\": \"object\","]
#[doc = "        \"properties\": {"]
#[doc = "          \"method\": {"]
#[doc = "            \"description\": \"Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          },"]
#[doc = "          \"path\": {"]
#[doc = "            \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "            \"type\": \"string\""]
#[doc = "          }"]
#[doc = "        },"]
#[doc = "        \"additionalProperties\": false"]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"additionalProperties\": false"]
#[doc = "  }"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
pub struct RestApiProviderConfigurationResources {
    #[doc = "Defaults to true. Whether to generate a Terraform data source type for this API object."]
    #[serde(default = "defaults::default_bool::<true>")]
    pub generate_data_source: bool,
    #[doc = "Defaults to true. Whether to generate a Terraform resource type for this API object."]
    #[serde(default = "defaults::default_bool::<true>")]
    pub generate_resource: bool,
    #[serde(flatten)]
    pub extra: ::std::collections::HashMap<
        ::std::string::String,
        RestApiProviderConfigurationResourcesExtraValue,
    >,
}
impl ::std::convert::From<&RestApiProviderConfigurationResources>
    for RestApiProviderConfigurationResources
{
    fn from(value: &RestApiProviderConfigurationResources) -> Self {
        value.clone()
    }
}
impl RestApiProviderConfigurationResources {
    pub fn builder() -> builder::RestApiProviderConfigurationResources {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResourcesExtraValue`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"required\": ["]
#[doc = "    \"path\""]
#[doc = "  ],"]
#[doc = "  \"properties\": {"]
#[doc = "    \"create\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"method\": {"]
#[doc = "          \"description\": \"Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"path\": {"]
#[doc = "          \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    },"]
#[doc = "    \"debug\": {"]
#[doc = "      \"description\": \"Whether to emit verbose debug output while working with the API object on the server.\","]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"destroy\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"method\": {"]
#[doc = "          \"description\": \"Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"path\": {"]
#[doc = "          \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    },"]
#[doc = "    \"force_new\": {"]
#[doc = "      \"description\": \"Any changes to these values will result in recreating the resource instead of updating.\","]
#[doc = "      \"type\": \"array\","]
#[doc = "      \"items\": {"]
#[doc = "        \"type\": \"string\""]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"id_attribute\": {"]
#[doc = "      \"description\": \"Defaults to id_attribute set on the provider. Allows per-resource override of id_attribute (see id_attribute provider config documentation)\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"ignore_all_server_changes\": {"]
#[doc = "      \"description\": \"By default Terraform will attempt to revert changes to remote resources. Set this to 'true' to ignore any remote changes. Default: false\","]
#[doc = "      \"default\": false,"]
#[doc = "      \"type\": \"boolean\""]
#[doc = "    },"]
#[doc = "    \"ignore_changes_to\": {"]
#[doc = "      \"description\": \"A list of fields to which remote changes will be ignored. For example, an API might add or remove metadata, such as a 'last_modified' field, which Terraform should not attempt to correct. To ignore changes to nested fields, use the dot syntax: 'metadata.timestamp'\","]
#[doc = "      \"type\": \"array\","]
#[doc = "      \"items\": {"]
#[doc = "        \"type\": \"string\""]
#[doc = "      }"]
#[doc = "    },"]
#[doc = "    \"object_id\": {"]
#[doc = "      \"description\": \"Defaults to the id learned by the provider during normal operations and id_attribute. Allows you to set the id manually. This is used in conjunction with the *_path attributes.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"path\": {"]
#[doc = "      \"description\": \"The API path on top of the base URL set in the provider that represents objects of this type on the API server.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"query_string\": {"]
#[doc = "      \"description\": \"Query string to be included in the path\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"read\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"method\": {"]
#[doc = "          \"description\": \"Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"path\": {"]
#[doc = "          \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"search\": {"]
#[doc = "          \"description\": \"Custom search for read_path.\","]
#[doc = "          \"type\": \"object\","]
#[doc = "          \"required\": ["]
#[doc = "            \"search_key\","]
#[doc = "            \"search_value\""]
#[doc = "          ],"]
#[doc = "          \"properties\": {"]
#[doc = "            \"query_string\": {"]
#[doc = "              \"description\": \"An optional query string to send when performing the search.\","]
#[doc = "              \"type\": \"string\""]
#[doc = "            },"]
#[doc = "            \"results_key\": {"]
#[doc = "              \"description\": \"When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.\","]
#[doc = "              \"type\": \"string\""]
#[doc = "            },"]
#[doc = "            \"search_key\": {"]
#[doc = "              \"description\": \"When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.\","]
#[doc = "              \"type\": \"string\""]
#[doc = "            },"]
#[doc = "            \"search_path\": {"]
#[doc = "              \"description\": \"The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.\","]
#[doc = "              \"type\": \"string\""]
#[doc = "            },"]
#[doc = "            \"search_value\": {"]
#[doc = "              \"description\": \"The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.\","]
#[doc = "              \"type\": \"string\""]
#[doc = "            }"]
#[doc = "          },"]
#[doc = "          \"additionalProperties\": false"]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    },"]
#[doc = "    \"update\": {"]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"properties\": {"]
#[doc = "        \"method\": {"]
#[doc = "          \"description\": \"Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"path\": {"]
#[doc = "          \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValue {
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub create: ::std::option::Option<RestApiProviderConfigurationResourcesExtraValueCreate>,
    #[doc = "Whether to emit verbose debug output while working with the API object on the server."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub debug: ::std::option::Option<bool>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub destroy: ::std::option::Option<RestApiProviderConfigurationResourcesExtraValueDestroy>,
    #[doc = "Any changes to these values will result in recreating the resource instead of updating."]
    #[serde(default, skip_serializing_if = "::std::vec::Vec::is_empty")]
    pub force_new: ::std::vec::Vec<::std::string::String>,
    #[doc = "Defaults to id_attribute set on the provider. Allows per-resource override of id_attribute (see id_attribute provider config documentation)"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub id_attribute: ::std::option::Option<::std::string::String>,
    #[doc = "By default Terraform will attempt to revert changes to remote resources. Set this to 'true' to ignore any remote changes. Default: false"]
    #[serde(default)]
    pub ignore_all_server_changes: bool,
    #[doc = "A list of fields to which remote changes will be ignored. For example, an API might add or remove metadata, such as a 'last_modified' field, which Terraform should not attempt to correct. To ignore changes to nested fields, use the dot syntax: 'metadata.timestamp'"]
    #[serde(default, skip_serializing_if = "::std::vec::Vec::is_empty")]
    pub ignore_changes_to: ::std::vec::Vec<::std::string::String>,
    #[doc = "Defaults to the id learned by the provider during normal operations and id_attribute. Allows you to set the id manually. This is used in conjunction with the *_path attributes."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub object_id: ::std::option::Option<::std::string::String>,
    #[doc = "The API path on top of the base URL set in the provider that represents objects of this type on the API server."]
    pub path: ::std::string::String,
    #[doc = "Query string to be included in the path"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub query_string: ::std::option::Option<::std::string::String>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub read: ::std::option::Option<RestApiProviderConfigurationResourcesExtraValueRead>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub update: ::std::option::Option<RestApiProviderConfigurationResourcesExtraValueUpdate>,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValue>
    for RestApiProviderConfigurationResourcesExtraValue
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValue) -> Self {
        value.clone()
    }
}
impl RestApiProviderConfigurationResourcesExtraValue {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValue {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResourcesExtraValueCreate`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"method\": {"]
#[doc = "      \"description\": \"Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"path\": {"]
#[doc = "      \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValueCreate {
    #[doc = "Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub method: ::std::option::Option<::std::string::String>,
    #[doc = "Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub path: ::std::option::Option<::std::string::String>,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValueCreate>
    for RestApiProviderConfigurationResourcesExtraValueCreate
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValueCreate) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueCreate {
    fn default() -> Self {
        Self {
            method: Default::default(),
            path: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationResourcesExtraValueCreate {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValueCreate {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResourcesExtraValueDestroy`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"method\": {"]
#[doc = "      \"description\": \"Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"path\": {"]
#[doc = "      \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValueDestroy {
    #[doc = "Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub method: ::std::option::Option<::std::string::String>,
    #[doc = "Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub path: ::std::option::Option<::std::string::String>,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValueDestroy>
    for RestApiProviderConfigurationResourcesExtraValueDestroy
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValueDestroy) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueDestroy {
    fn default() -> Self {
        Self {
            method: Default::default(),
            path: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationResourcesExtraValueDestroy {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValueDestroy {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResourcesExtraValueRead`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"method\": {"]
#[doc = "      \"description\": \"Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"path\": {"]
#[doc = "      \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"search\": {"]
#[doc = "      \"description\": \"Custom search for read_path.\","]
#[doc = "      \"type\": \"object\","]
#[doc = "      \"required\": ["]
#[doc = "        \"search_key\","]
#[doc = "        \"search_value\""]
#[doc = "      ],"]
#[doc = "      \"properties\": {"]
#[doc = "        \"query_string\": {"]
#[doc = "          \"description\": \"An optional query string to send when performing the search.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"results_key\": {"]
#[doc = "          \"description\": \"When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"search_key\": {"]
#[doc = "          \"description\": \"When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"search_path\": {"]
#[doc = "          \"description\": \"The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        },"]
#[doc = "        \"search_value\": {"]
#[doc = "          \"description\": \"The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.\","]
#[doc = "          \"type\": \"string\""]
#[doc = "        }"]
#[doc = "      },"]
#[doc = "      \"additionalProperties\": false"]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValueRead {
    #[doc = "Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub method: ::std::option::Option<::std::string::String>,
    #[doc = "Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub path: ::std::option::Option<::std::string::String>,
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub search: ::std::option::Option<RestApiProviderConfigurationResourcesExtraValueReadSearch>,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValueRead>
    for RestApiProviderConfigurationResourcesExtraValueRead
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValueRead) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueRead {
    fn default() -> Self {
        Self {
            method: Default::default(),
            path: Default::default(),
            search: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationResourcesExtraValueRead {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValueRead {
        Default::default()
    }
}
#[doc = "Custom search for read_path."]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"description\": \"Custom search for read_path.\","]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"required\": ["]
#[doc = "    \"search_key\","]
#[doc = "    \"search_value\""]
#[doc = "  ],"]
#[doc = "  \"properties\": {"]
#[doc = "    \"query_string\": {"]
#[doc = "      \"description\": \"An optional query string to send when performing the search.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"results_key\": {"]
#[doc = "      \"description\": \"When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"search_key\": {"]
#[doc = "      \"description\": \"When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"search_path\": {"]
#[doc = "      \"description\": \"The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"search_value\": {"]
#[doc = "      \"description\": \"The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValueReadSearch {
    #[doc = "An optional query string to send when performing the search."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub query_string: ::std::option::Option<::std::string::String>,
    #[doc = "When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub results_key: ::std::option::Option<::std::string::String>,
    #[doc = "When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object."]
    pub search_key: ::std::string::String,
    #[doc = "The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub search_path: ::std::option::Option<::std::string::String>,
    #[doc = "The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used."]
    pub search_value: ::std::string::String,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValueReadSearch>
    for RestApiProviderConfigurationResourcesExtraValueReadSearch
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValueReadSearch) -> Self {
        value.clone()
    }
}
impl RestApiProviderConfigurationResourcesExtraValueReadSearch {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValueReadSearch {
        Default::default()
    }
}
#[doc = "`RestApiProviderConfigurationResourcesExtraValueUpdate`"]
#[doc = r""]
#[doc = r" <details><summary>JSON schema</summary>"]
#[doc = r""]
#[doc = r" ```json"]
#[doc = "{"]
#[doc = "  \"type\": \"object\","]
#[doc = "  \"properties\": {"]
#[doc = "    \"method\": {"]
#[doc = "      \"description\": \"Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    },"]
#[doc = "    \"path\": {"]
#[doc = "      \"description\": \"Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.\","]
#[doc = "      \"type\": \"string\""]
#[doc = "    }"]
#[doc = "  },"]
#[doc = "  \"additionalProperties\": false"]
#[doc = "}"]
#[doc = r" ```"]
#[doc = r" </details>"]
#[derive(:: serde :: Deserialize, :: serde :: Serialize, Clone, Debug)]
#[serde(deny_unknown_fields)]
pub struct RestApiProviderConfigurationResourcesExtraValueUpdate {
    #[doc = "Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)"]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub method: ::std::option::Option<::std::string::String>,
    #[doc = "Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute."]
    #[serde(default, skip_serializing_if = "::std::option::Option::is_none")]
    pub path: ::std::option::Option<::std::string::String>,
}
impl ::std::convert::From<&RestApiProviderConfigurationResourcesExtraValueUpdate>
    for RestApiProviderConfigurationResourcesExtraValueUpdate
{
    fn from(value: &RestApiProviderConfigurationResourcesExtraValueUpdate) -> Self {
        value.clone()
    }
}
impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueUpdate {
    fn default() -> Self {
        Self {
            method: Default::default(),
            path: Default::default(),
        }
    }
}
impl RestApiProviderConfigurationResourcesExtraValueUpdate {
    pub fn builder() -> builder::RestApiProviderConfigurationResourcesExtraValueUpdate {
        Default::default()
    }
}
#[doc = r" Types for composing complex structures."]
pub mod builder {
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfiguration {
        global: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationGlobal>,
            ::std::string::String,
        >,
        resources: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResources>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfiguration {
        fn default() -> Self {
            Self {
                global: Ok(Default::default()),
                resources: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfiguration {
        pub fn global<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<super::RestApiProviderConfigurationGlobal>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.global = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for global: {}", e));
            self
        }
        pub fn resources<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<super::RestApiProviderConfigurationResources>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.resources = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for resources: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfiguration> for super::RestApiProviderConfiguration {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfiguration,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                global: value.global?,
                resources: value.resources?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfiguration> for RestApiProviderConfiguration {
        fn from(value: super::RestApiProviderConfiguration) -> Self {
            Self {
                global: Ok(value.global),
                resources: Ok(value.resources),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationGlobal {
        cert_file: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        cert_string: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        copy_keys:
            ::std::result::Result<::std::vec::Vec<::std::string::String>, ::std::string::String>,
        create_method: ::std::result::Result<::std::string::String, ::std::string::String>,
        create_returns_object:
            ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        debug: ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        destroy_method: ::std::result::Result<::std::string::String, ::std::string::String>,
        headers: ::std::result::Result<
            ::std::collections::HashMap<::std::string::String, ::std::string::String>,
            ::std::string::String,
        >,
        id_attribute: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        insecure: ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        key_file: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        key_string: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        oauth_client_credentials: ::std::result::Result<
            ::std::vec::Vec<super::RestApiProviderConfigurationGlobalOauthClientCredentialsItem>,
            ::std::string::String,
        >,
        password: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        rate_limit: ::std::result::Result<::std::option::Option<f64>, ::std::string::String>,
        read_method: ::std::result::Result<::std::string::String, ::std::string::String>,
        root_ca_file: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        root_ca_string: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        test_path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        timeout: ::std::result::Result<::std::option::Option<f64>, ::std::string::String>,
        update_method: ::std::result::Result<::std::string::String, ::std::string::String>,
        uri: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        use_cookies: ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        username: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        write_returns_object:
            ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        xssi_prefix: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationGlobal {
        fn default() -> Self {
            Self {
                cert_file: Ok(Default::default()),
                cert_string: Ok(Default::default()),
                copy_keys: Ok(Default::default()),
                create_method: Ok(
                    super::defaults::rest_api_provider_configuration_global_create_method(),
                ),
                create_returns_object: Ok(Default::default()),
                debug: Ok(Default::default()),
                destroy_method: Ok(
                    super::defaults::rest_api_provider_configuration_global_destroy_method(),
                ),
                headers: Ok(Default::default()),
                id_attribute: Ok(Default::default()),
                insecure: Ok(Default::default()),
                key_file: Ok(Default::default()),
                key_string: Ok(Default::default()),
                oauth_client_credentials: Ok(Default::default()),
                password: Ok(Default::default()),
                rate_limit: Ok(Default::default()),
                read_method: Ok(
                    super::defaults::rest_api_provider_configuration_global_read_method(),
                ),
                root_ca_file: Ok(Default::default()),
                root_ca_string: Ok(Default::default()),
                test_path: Ok(Default::default()),
                timeout: Ok(Default::default()),
                update_method: Ok(
                    super::defaults::rest_api_provider_configuration_global_update_method(),
                ),
                uri: Ok(Default::default()),
                use_cookies: Ok(Default::default()),
                username: Ok(Default::default()),
                write_returns_object: Ok(Default::default()),
                xssi_prefix: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationGlobal {
        pub fn cert_file<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.cert_file = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for cert_file: {}", e));
            self
        }
        pub fn cert_string<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.cert_string = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for cert_string: {}", e));
            self
        }
        pub fn copy_keys<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::vec::Vec<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.copy_keys = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for copy_keys: {}", e));
            self
        }
        pub fn create_method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.create_method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for create_method: {}", e));
            self
        }
        pub fn create_returns_object<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.create_returns_object = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for create_returns_object: {}",
                    e
                )
            });
            self
        }
        pub fn debug<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.debug = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for debug: {}", e));
            self
        }
        pub fn destroy_method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.destroy_method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for destroy_method: {}", e));
            self
        }
        pub fn headers<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::collections::HashMap<::std::string::String, ::std::string::String>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.headers = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for headers: {}", e));
            self
        }
        pub fn id_attribute<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.id_attribute = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for id_attribute: {}", e));
            self
        }
        pub fn insecure<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.insecure = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for insecure: {}", e));
            self
        }
        pub fn key_file<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.key_file = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for key_file: {}", e));
            self
        }
        pub fn key_string<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.key_string = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for key_string: {}", e));
            self
        }
        pub fn oauth_client_credentials<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::vec::Vec<
                    super::RestApiProviderConfigurationGlobalOauthClientCredentialsItem,
                >,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.oauth_client_credentials = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for oauth_client_credentials: {}",
                    e
                )
            });
            self
        }
        pub fn password<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.password = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for password: {}", e));
            self
        }
        pub fn rate_limit<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<f64>>,
            T::Error: ::std::fmt::Display,
        {
            self.rate_limit = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for rate_limit: {}", e));
            self
        }
        pub fn read_method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.read_method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for read_method: {}", e));
            self
        }
        pub fn root_ca_file<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.root_ca_file = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for root_ca_file: {}", e));
            self
        }
        pub fn root_ca_string<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.root_ca_string = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for root_ca_string: {}", e));
            self
        }
        pub fn test_path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.test_path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for test_path: {}", e));
            self
        }
        pub fn timeout<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<f64>>,
            T::Error: ::std::fmt::Display,
        {
            self.timeout = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for timeout: {}", e));
            self
        }
        pub fn update_method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.update_method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for update_method: {}", e));
            self
        }
        pub fn uri<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.uri = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for uri: {}", e));
            self
        }
        pub fn use_cookies<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.use_cookies = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for use_cookies: {}", e));
            self
        }
        pub fn username<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.username = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for username: {}", e));
            self
        }
        pub fn write_returns_object<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.write_returns_object = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for write_returns_object: {}",
                    e
                )
            });
            self
        }
        pub fn xssi_prefix<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.xssi_prefix = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for xssi_prefix: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationGlobal>
        for super::RestApiProviderConfigurationGlobal
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationGlobal,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                cert_file: value.cert_file?,
                cert_string: value.cert_string?,
                copy_keys: value.copy_keys?,
                create_method: value.create_method?,
                create_returns_object: value.create_returns_object?,
                debug: value.debug?,
                destroy_method: value.destroy_method?,
                headers: value.headers?,
                id_attribute: value.id_attribute?,
                insecure: value.insecure?,
                key_file: value.key_file?,
                key_string: value.key_string?,
                oauth_client_credentials: value.oauth_client_credentials?,
                password: value.password?,
                rate_limit: value.rate_limit?,
                read_method: value.read_method?,
                root_ca_file: value.root_ca_file?,
                root_ca_string: value.root_ca_string?,
                test_path: value.test_path?,
                timeout: value.timeout?,
                update_method: value.update_method?,
                uri: value.uri?,
                use_cookies: value.use_cookies?,
                username: value.username?,
                write_returns_object: value.write_returns_object?,
                xssi_prefix: value.xssi_prefix?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationGlobal>
        for RestApiProviderConfigurationGlobal
    {
        fn from(value: super::RestApiProviderConfigurationGlobal) -> Self {
            Self {
                cert_file: Ok(value.cert_file),
                cert_string: Ok(value.cert_string),
                copy_keys: Ok(value.copy_keys),
                create_method: Ok(value.create_method),
                create_returns_object: Ok(value.create_returns_object),
                debug: Ok(value.debug),
                destroy_method: Ok(value.destroy_method),
                headers: Ok(value.headers),
                id_attribute: Ok(value.id_attribute),
                insecure: Ok(value.insecure),
                key_file: Ok(value.key_file),
                key_string: Ok(value.key_string),
                oauth_client_credentials: Ok(value.oauth_client_credentials),
                password: Ok(value.password),
                rate_limit: Ok(value.rate_limit),
                read_method: Ok(value.read_method),
                root_ca_file: Ok(value.root_ca_file),
                root_ca_string: Ok(value.root_ca_string),
                test_path: Ok(value.test_path),
                timeout: Ok(value.timeout),
                update_method: Ok(value.update_method),
                uri: Ok(value.uri),
                use_cookies: Ok(value.use_cookies),
                username: Ok(value.username),
                write_returns_object: Ok(value.write_returns_object),
                xssi_prefix: Ok(value.xssi_prefix),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
        endpoint_params: ::std::result::Result<
            ::serde_json::Map<::std::string::String, ::serde_json::Value>,
            ::std::string::String,
        >,
        oauth_client_id: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        oauth_client_secret: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        oauth_scopes:
            ::std::result::Result<::std::vec::Vec<::std::string::String>, ::std::string::String>,
        oauth_token_endpoint: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
        fn default() -> Self {
            Self {
                endpoint_params: Ok(Default::default()),
                oauth_client_id: Ok(Default::default()),
                oauth_client_secret: Ok(Default::default()),
                oauth_scopes: Ok(Default::default()),
                oauth_token_endpoint: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationGlobalOauthClientCredentialsItem {
        pub fn endpoint_params<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::serde_json::Map<::std::string::String, ::serde_json::Value>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.endpoint_params = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for endpoint_params: {}", e));
            self
        }
        pub fn oauth_client_id<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.oauth_client_id = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for oauth_client_id: {}", e));
            self
        }
        pub fn oauth_client_secret<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.oauth_client_secret = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for oauth_client_secret: {}",
                    e
                )
            });
            self
        }
        pub fn oauth_scopes<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::vec::Vec<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.oauth_scopes = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for oauth_scopes: {}", e));
            self
        }
        pub fn oauth_token_endpoint<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.oauth_token_endpoint = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for oauth_token_endpoint: {}",
                    e
                )
            });
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationGlobalOauthClientCredentialsItem>
        for super::RestApiProviderConfigurationGlobalOauthClientCredentialsItem
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationGlobalOauthClientCredentialsItem,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                endpoint_params: value.endpoint_params?,
                oauth_client_id: value.oauth_client_id?,
                oauth_client_secret: value.oauth_client_secret?,
                oauth_scopes: value.oauth_scopes?,
                oauth_token_endpoint: value.oauth_token_endpoint?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationGlobalOauthClientCredentialsItem>
        for RestApiProviderConfigurationGlobalOauthClientCredentialsItem
    {
        fn from(
            value: super::RestApiProviderConfigurationGlobalOauthClientCredentialsItem,
        ) -> Self {
            Self {
                endpoint_params: Ok(value.endpoint_params),
                oauth_client_id: Ok(value.oauth_client_id),
                oauth_client_secret: Ok(value.oauth_client_secret),
                oauth_scopes: Ok(value.oauth_scopes),
                oauth_token_endpoint: Ok(value.oauth_token_endpoint),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResources {
        generate_data_source: ::std::result::Result<bool, ::std::string::String>,
        generate_resource: ::std::result::Result<bool, ::std::string::String>,
        extra: ::std::result::Result<
            ::std::collections::HashMap<
                ::std::string::String,
                super::RestApiProviderConfigurationResourcesExtraValue,
            >,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResources {
        fn default() -> Self {
            Self {
                generate_data_source: Ok(super::defaults::default_bool::<true>()),
                generate_resource: Ok(super::defaults::default_bool::<true>()),
                extra: Err("no value supplied for extra".to_string()),
            }
        }
    }
    impl RestApiProviderConfigurationResources {
        pub fn generate_data_source<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<bool>,
            T::Error: ::std::fmt::Display,
        {
            self.generate_data_source = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for generate_data_source: {}",
                    e
                )
            });
            self
        }
        pub fn generate_resource<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<bool>,
            T::Error: ::std::fmt::Display,
        {
            self.generate_resource = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for generate_resource: {}",
                    e
                )
            });
            self
        }
        pub fn extra<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::collections::HashMap<
                    ::std::string::String,
                    super::RestApiProviderConfigurationResourcesExtraValue,
                >,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.extra = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for extra: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResources>
        for super::RestApiProviderConfigurationResources
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResources,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                generate_data_source: value.generate_data_source?,
                generate_resource: value.generate_resource?,
                extra: value.extra?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResources>
        for RestApiProviderConfigurationResources
    {
        fn from(value: super::RestApiProviderConfigurationResources) -> Self {
            Self {
                generate_data_source: Ok(value.generate_data_source),
                generate_resource: Ok(value.generate_resource),
                extra: Ok(value.extra),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValue {
        create: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueCreate>,
            ::std::string::String,
        >,
        debug: ::std::result::Result<::std::option::Option<bool>, ::std::string::String>,
        destroy: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueDestroy>,
            ::std::string::String,
        >,
        force_new:
            ::std::result::Result<::std::vec::Vec<::std::string::String>, ::std::string::String>,
        id_attribute: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        ignore_all_server_changes: ::std::result::Result<bool, ::std::string::String>,
        ignore_changes_to:
            ::std::result::Result<::std::vec::Vec<::std::string::String>, ::std::string::String>,
        object_id: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        path: ::std::result::Result<::std::string::String, ::std::string::String>,
        query_string: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        read: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueRead>,
            ::std::string::String,
        >,
        update: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueUpdate>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValue {
        fn default() -> Self {
            Self {
                create: Ok(Default::default()),
                debug: Ok(Default::default()),
                destroy: Ok(Default::default()),
                force_new: Ok(Default::default()),
                id_attribute: Ok(Default::default()),
                ignore_all_server_changes: Ok(Default::default()),
                ignore_changes_to: Ok(Default::default()),
                object_id: Ok(Default::default()),
                path: Err("no value supplied for path".to_string()),
                query_string: Ok(Default::default()),
                read: Ok(Default::default()),
                update: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValue {
        pub fn create<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueCreate>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.create = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for create: {}", e));
            self
        }
        pub fn debug<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<bool>>,
            T::Error: ::std::fmt::Display,
        {
            self.debug = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for debug: {}", e));
            self
        }
        pub fn destroy<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<
                    super::RestApiProviderConfigurationResourcesExtraValueDestroy,
                >,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.destroy = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for destroy: {}", e));
            self
        }
        pub fn force_new<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::vec::Vec<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.force_new = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for force_new: {}", e));
            self
        }
        pub fn id_attribute<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.id_attribute = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for id_attribute: {}", e));
            self
        }
        pub fn ignore_all_server_changes<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<bool>,
            T::Error: ::std::fmt::Display,
        {
            self.ignore_all_server_changes = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for ignore_all_server_changes: {}",
                    e
                )
            });
            self
        }
        pub fn ignore_changes_to<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::vec::Vec<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.ignore_changes_to = value.try_into().map_err(|e| {
                format!(
                    "error converting supplied value for ignore_changes_to: {}",
                    e
                )
            });
            self
        }
        pub fn object_id<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.object_id = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for object_id: {}", e));
            self
        }
        pub fn path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for path: {}", e));
            self
        }
        pub fn query_string<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.query_string = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for query_string: {}", e));
            self
        }
        pub fn read<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueRead>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.read = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for read: {}", e));
            self
        }
        pub fn update<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueUpdate>,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.update = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for update: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValue>
        for super::RestApiProviderConfigurationResourcesExtraValue
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValue,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                create: value.create?,
                debug: value.debug?,
                destroy: value.destroy?,
                force_new: value.force_new?,
                id_attribute: value.id_attribute?,
                ignore_all_server_changes: value.ignore_all_server_changes?,
                ignore_changes_to: value.ignore_changes_to?,
                object_id: value.object_id?,
                path: value.path?,
                query_string: value.query_string?,
                read: value.read?,
                update: value.update?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValue>
        for RestApiProviderConfigurationResourcesExtraValue
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValue) -> Self {
            Self {
                create: Ok(value.create),
                debug: Ok(value.debug),
                destroy: Ok(value.destroy),
                force_new: Ok(value.force_new),
                id_attribute: Ok(value.id_attribute),
                ignore_all_server_changes: Ok(value.ignore_all_server_changes),
                ignore_changes_to: Ok(value.ignore_changes_to),
                object_id: Ok(value.object_id),
                path: Ok(value.path),
                query_string: Ok(value.query_string),
                read: Ok(value.read),
                update: Ok(value.update),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValueCreate {
        method: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueCreate {
        fn default() -> Self {
            Self {
                method: Ok(Default::default()),
                path: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValueCreate {
        pub fn method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for method: {}", e));
            self
        }
        pub fn path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for path: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValueCreate>
        for super::RestApiProviderConfigurationResourcesExtraValueCreate
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValueCreate,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                method: value.method?,
                path: value.path?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValueCreate>
        for RestApiProviderConfigurationResourcesExtraValueCreate
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValueCreate) -> Self {
            Self {
                method: Ok(value.method),
                path: Ok(value.path),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValueDestroy {
        method: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueDestroy {
        fn default() -> Self {
            Self {
                method: Ok(Default::default()),
                path: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValueDestroy {
        pub fn method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for method: {}", e));
            self
        }
        pub fn path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for path: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValueDestroy>
        for super::RestApiProviderConfigurationResourcesExtraValueDestroy
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValueDestroy,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                method: value.method?,
                path: value.path?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValueDestroy>
        for RestApiProviderConfigurationResourcesExtraValueDestroy
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValueDestroy) -> Self {
            Self {
                method: Ok(value.method),
                path: Ok(value.path),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValueRead {
        method: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        search: ::std::result::Result<
            ::std::option::Option<super::RestApiProviderConfigurationResourcesExtraValueReadSearch>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueRead {
        fn default() -> Self {
            Self {
                method: Ok(Default::default()),
                path: Ok(Default::default()),
                search: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValueRead {
        pub fn method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for method: {}", e));
            self
        }
        pub fn path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for path: {}", e));
            self
        }
        pub fn search<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<
                ::std::option::Option<
                    super::RestApiProviderConfigurationResourcesExtraValueReadSearch,
                >,
            >,
            T::Error: ::std::fmt::Display,
        {
            self.search = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for search: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValueRead>
        for super::RestApiProviderConfigurationResourcesExtraValueRead
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValueRead,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                method: value.method?,
                path: value.path?,
                search: value.search?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValueRead>
        for RestApiProviderConfigurationResourcesExtraValueRead
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValueRead) -> Self {
            Self {
                method: Ok(value.method),
                path: Ok(value.path),
                search: Ok(value.search),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValueReadSearch {
        query_string: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        results_key: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        search_key: ::std::result::Result<::std::string::String, ::std::string::String>,
        search_path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        search_value: ::std::result::Result<::std::string::String, ::std::string::String>,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueReadSearch {
        fn default() -> Self {
            Self {
                query_string: Ok(Default::default()),
                results_key: Ok(Default::default()),
                search_key: Err("no value supplied for search_key".to_string()),
                search_path: Ok(Default::default()),
                search_value: Err("no value supplied for search_value".to_string()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValueReadSearch {
        pub fn query_string<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.query_string = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for query_string: {}", e));
            self
        }
        pub fn results_key<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.results_key = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for results_key: {}", e));
            self
        }
        pub fn search_key<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.search_key = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for search_key: {}", e));
            self
        }
        pub fn search_path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.search_path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for search_path: {}", e));
            self
        }
        pub fn search_value<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::string::String>,
            T::Error: ::std::fmt::Display,
        {
            self.search_value = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for search_value: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValueReadSearch>
        for super::RestApiProviderConfigurationResourcesExtraValueReadSearch
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValueReadSearch,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                query_string: value.query_string?,
                results_key: value.results_key?,
                search_key: value.search_key?,
                search_path: value.search_path?,
                search_value: value.search_value?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValueReadSearch>
        for RestApiProviderConfigurationResourcesExtraValueReadSearch
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValueReadSearch) -> Self {
            Self {
                query_string: Ok(value.query_string),
                results_key: Ok(value.results_key),
                search_key: Ok(value.search_key),
                search_path: Ok(value.search_path),
                search_value: Ok(value.search_value),
            }
        }
    }
    #[derive(Clone, Debug)]
    pub struct RestApiProviderConfigurationResourcesExtraValueUpdate {
        method: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
        path: ::std::result::Result<
            ::std::option::Option<::std::string::String>,
            ::std::string::String,
        >,
    }
    impl ::std::default::Default for RestApiProviderConfigurationResourcesExtraValueUpdate {
        fn default() -> Self {
            Self {
                method: Ok(Default::default()),
                path: Ok(Default::default()),
            }
        }
    }
    impl RestApiProviderConfigurationResourcesExtraValueUpdate {
        pub fn method<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.method = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for method: {}", e));
            self
        }
        pub fn path<T>(mut self, value: T) -> Self
        where
            T: ::std::convert::TryInto<::std::option::Option<::std::string::String>>,
            T::Error: ::std::fmt::Display,
        {
            self.path = value
                .try_into()
                .map_err(|e| format!("error converting supplied value for path: {}", e));
            self
        }
    }
    impl ::std::convert::TryFrom<RestApiProviderConfigurationResourcesExtraValueUpdate>
        for super::RestApiProviderConfigurationResourcesExtraValueUpdate
    {
        type Error = super::error::ConversionError;
        fn try_from(
            value: RestApiProviderConfigurationResourcesExtraValueUpdate,
        ) -> ::std::result::Result<Self, super::error::ConversionError> {
            Ok(Self {
                method: value.method?,
                path: value.path?,
            })
        }
    }
    impl ::std::convert::From<super::RestApiProviderConfigurationResourcesExtraValueUpdate>
        for RestApiProviderConfigurationResourcesExtraValueUpdate
    {
        fn from(value: super::RestApiProviderConfigurationResourcesExtraValueUpdate) -> Self {
            Self {
                method: Ok(value.method),
                path: Ok(value.path),
            }
        }
    }
}
#[doc = r" Generation of default values for serde."]
pub mod defaults {
    pub(super) fn default_bool<const V: bool>() -> bool {
        V
    }
    pub(super) fn rest_api_provider_configuration_global_create_method() -> ::std::string::String {
        "POST".to_string()
    }
    pub(super) fn rest_api_provider_configuration_global_destroy_method() -> ::std::string::String {
        "DELETE".to_string()
    }
    pub(super) fn rest_api_provider_configuration_global_read_method() -> ::std::string::String {
        "GET".to_string()
    }
    pub(super) fn rest_api_provider_configuration_global_update_method() -> ::std::string::String {
        "PUT".to_string()
    }
}
