package provider_spec

type RESTAPIProviderConfiguration struct {
	GlobalDefaults *GlobalDefaults  `json:"global_defaults,omitempty"`
	Resources      *ResourcesSchema `json:"resources,omitempty"`
}

type GlobalDefaults struct {
	CreateMethod  string `json:"create_method,omitempty"`  // Defaults to POST. The HTTP method used to CREATE objects of this type on the API server.
	Debug         bool   `json:"debug,omitempty"`          // Enabling this will cause lots of debug information to be printed to STDOUT by the API client.
	DestroyMethod string `json:"destroy_method,omitempty"` // Defaults to DELETE. The HTTP method used to DELETE objects of this type on the API server.
	Headers       *struct {
		// Additional properties, not valided now
		OtherProps map[string]string `json:",inline"`
	} `json:"headers,omitempty"` // A map of header names and values to set on all outbound requests. This is useful if you want to use a script via the 'external' provider or provide a pre-approved token or change Content-Type from application/json. If username and password are set and Authorization is one of the headers defined here, the BASIC auth credentials take precedence.
	IdAttribute  string `json:"id_attribute,omitempty"`  // When set, this key will be used to operate on REST objects. For example, if the ID is set to 'name', changes to the API object will be to http://foo.com/bar/VALUE_OF_NAME. This value may also be a '/'-delimited path to the id attribute if it is multiple levels deep in the data (such as attributes/id in the case of an object { "attributes": { "id": 1234 }, "config": { "name": "foo", "something": "bar" } }
	ReadMethod   string `json:"read_method,omitempty"`   // Defaults to GET. The HTTP method used to READ objects of this type on the API server.
	UpdateMethod string `json:"update_method,omitempty"` // Defaults to PUT. The HTTP method used to UPDATE objects of this type on the API server.
	Uri          string `json:"uri,omitempty"`           // URI of the REST API endpoint. This serves as the base of all requests.
}

type ResourcesSchema struct {
	GenerateDataSource bool `json:"generate_data_source,omitempty"` // Defaults to true. Whether to generate a Terraform data source type for this API object.
	GenerateResource   bool `json:"generate_resource,omitempty"`    // Defaults to true. Whether to generate a Terraform resource type for this API object.

	// Additional properties, not valided now
	OtherProps map[string]ResourceSchema `json:",inline"`
}

type ResourceSchema struct {
	Create *struct {
		Method string `json:"method,omitempty"` // Defaults to global {create_method}. Allows per-resource override of create_method (see create_method config documentation)
		Path   string `json:"path,omitempty"`   // Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.
	} `json:"create,omitempty"`
	Debug   bool `json:"debug,omitempty"` // Whether to emit verbose debug output while working with the API object on the server.
	Destroy *struct {
		Method string `json:"method,omitempty"` // Defaults to global {destroy_method}. Allows per-resource override of create_method (see create_method config documentation)
		Path   string `json:"path,omitempty"`   // Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.
	} `json:"destroy,omitempty"`
	ForceNew               []string `json:"force_new,omitempty"`                 // Any changes to these values will result in recreating the resource instead of updating.
	ForceRecreate          bool     `json:"force_recreate,omitempty"`            // If set to true, any changes to the resource will recreate it instead of updating.
	IdAttribute            string   `json:"id_attribute,omitempty"`              // Defaults to id_attribute set on the provider. Allows per-resource override of id_attribute (see id_attribute provider config documentation)
	IgnoreAllServerChanges bool     `json:"ignore_all_server_changes,omitempty"` // By default Terraform will attempt to revert changes to remote resources. Set this to 'true' to ignore any remote changes. Default: false
	IgnoreChangesTo        []string `json:"ignore_changes_to,omitempty"`         // A list of fields to which remote changes will be ignored. For example, an API might add or remove metadata, such as a 'last_modified' field, which Terraform should not attempt to correct. To ignore changes to nested fields, use the dot syntax: 'metadata.timestamp'
	ObjectId               string   `json:"object_id,omitempty"`                 // Defaults to the id learned by the provider during normal operations and id_attribute. Allows you to set the id manually. This is used in conjunction with the *_path attributes.
	Path                   string   `json:"path"`                                // The API path on top of the base URL set in the provider that represents objects of this type on the API server.
	QueryString            string   `json:"query_string,omitempty"`              // Query string to be included in the path
	Read                   *struct {
		Method string `json:"method,omitempty"` // Defaults to global {read_method}. Allows per-resource override of create_method (see create_method config documentation)
		Path   string `json:"path,omitempty"`   // Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.
		Search *struct {
			QueryString string `json:"query_string,omitempty"` // An optional query string to send when performing the search.
			ResultsKey  string `json:"results_key,omitempty"`  // When issuing a GET to the path, this JSON key is used to locate the results array. The format is 'field/field/field'. Example: 'results/values'. If omitted, it is assumed the results coming back are already an array and are to be used exactly as-is.
			SearchKey   string `json:"search_key"`             // When reading search results from the API, this key is used to identify the specific record to read. This should be a unique record such as 'name'. Similar to results_key, the value may be in the format of 'field/field/field' to search for data deeper in the returned object.
			SearchPath  string `json:"search_path,omitempty"`  // The API path on top of the base URL set in the provider that represents the location to search for objects of this type on the API server. If not set, defaults to the value of path.
			SearchValue string `json:"search_value"`           // The value of 'search_key' will be compared to this value to determine if the correct object was found. Example: if 'search_key' is 'name' and 'search_value' is 'foo', the record in the array returned by the API with name=foo will be used.
		} `json:"search,omitempty"` // Custom search for read_path.
	} `json:"read,omitempty"`
	Update *struct {
		Method string `json:"method,omitempty"` // Defaults to global {update_method}. Allows per-resource override of create_method (see create_method config documentation)
		Path   string `json:"path,omitempty"`   // Defaults to {path}. The API path that represents where to CREATE (POST) objects of this type on the API server. The string {id} will be replaced with the terraform ID of the object if the data contains the id_attribute.
	} `json:"update,omitempty"`
}
