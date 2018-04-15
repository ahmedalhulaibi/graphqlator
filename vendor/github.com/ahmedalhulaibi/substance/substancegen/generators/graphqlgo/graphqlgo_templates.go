package graphqlgo

/*GraphqlGoExecuteQueryFunc boilerplate string for graphql-go function to execute a graphql query*/
var GraphqlGoExecuteQueryFunc = `
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
`

var graphqlGoMainFunc = `
var DB *gorm.DB


func main() {

	DB, _ = gorm.Open("{{.DbType}}","{{.ConnectionString}}")
	defer DB.Close()


	fmt.Println("Test with Get	:	curl -g 'http://localhost:8080/graphql?query={ {{.SampleQuery}} }'")

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: Fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	gHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", gHandler)

	fmt.Println("Now server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}
`

var graphqlTypesTemplate = `{{range $key, $value := . }}
var {{.LowerName}}Type = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "{{.Name}}",
		Fields: graphql.Fields{ {{range .Properties}}
			"{{.ScalarNameUpper}}":&graphql.Field{
				Type: {{index .AltScalarType "graphql-go"}},
			},{{end}}
		},
	},
)
{{end}}`

var graphqlGoFieldsTemplate = `
	var QueryFields = graphql.Fields{ {{range $key, $value := . }}{{$name := .Name}}
		"{{.Name}}": &graphql.Field{
			Type: {{.LowerName}}Type,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				{{.Name}}Obj := {{.Name}}{}
				DB.First(&{{.Name}}Obj){{range .Properties}}{{if .IsObjectType}}
				{{.ScalarName}}Obj := {{if .IsList}}[]{{end}}{{.ScalarType}}{}
				DB.Model(&{{$name}}Obj).Association("{{.ScalarName}}").Find(&{{.ScalarName}}Obj)
				{{$name}}Obj.{{.ScalarName}} = {{if .IsList}}append({{$name}}Obj.{{.ScalarName}}, {{.ScalarName}}Obj...){{else}}{{.ScalarName}}Obj{{end}}{{end}}{{end}}
				return {{$name}}Obj, nil
			},
		},{{end}}
}
`

var graphqlQueryTemplate = `{{range $key, $value := . }}{{.Name}} { {{range .Properties}}{{if not .IsObjectType}}{{.ScalarName}},{{end}} {{end}}},{{end}}`
