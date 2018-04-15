package gqlschema

var graphqlSchemaTypesTemplate = `{{range $key, $value := . }}type {{.Name}} { {{range .Properties}}
	{{.ScalarName}}: {{if .IsList}}[{{index .AltScalarType "gqlschema"}}]{{else}}{{index .AltScalarType "gqlschema"}}{{end}}{{if .Nullable}}{{else}}!{{end}}{{end}}
}
{{end}}`

var graphqlSchemaQueriesTemplate = `{{range .}}{{end}}`

var graphqlSchemaGetQueriesTemplate = `{{range $key, $value := . }}
	# {{.Name}} returns first {{.Name}} in database table
	{{.Name}}: {{.Name}}
	# Get{{.Name}} takes the properties of {{.Name}} as search parameters. It will return all {{.Name}} rows found that matches the search criteria. Null input paramters are valid.
	Get{{.Name}}({{range .Properties}}{{.ScalarName}}: {{if .IsList}}[{{index .AltScalarType "gqlschema"}}]{{else}}{{index .AltScalarType "gqlschema"}}{{end}}, {{end}}): [{{.Name}}]{{end}}
`
