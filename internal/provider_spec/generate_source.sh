OUTFILE=src/provider_spec/provider_schema.go
go run github.com/RyoJerryYu/go-jsonschema/cmd/jsonschemagen src/provider_spec/rest_api_provider_schema.json --with-additional-properties -n provider_spec > $OUTFILE
sed -i 's/`json:"-"`/`json:",inline"`/' $OUTFILE