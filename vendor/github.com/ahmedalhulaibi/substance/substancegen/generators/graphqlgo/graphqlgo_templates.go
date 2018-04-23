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

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: QueryFields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: MutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
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

var graphqlGoFieldsQueryTemplate = `
var QueryFields graphql.Fields

func init() {
	QueryFields = make(graphql.Fields,1)
	{{template "graphqlFieldsGet" .}}
}
`

var graphqlQueryTemplate = `{{range $key, $value := . }}Get{{.Name}} { {{range .Properties}}{{if not .IsObjectType}}{{.ScalarName}},{{end}} {{end}}},{{end}}`

var graphqlGoQueryFieldsGetTemplate = `{{define "graphqlFieldsGet"}}{{range $key, $value := . }}
	QueryFields["Get{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}{{$name := .Name}}
			var Result{{$name}}Obj {{.Name}}
			Get{{.Name}}(DB,Query{{.Name}}Obj,&Result{{$name}}Obj){{range .Properties}}{{if .IsObjectType}}
			{{.ScalarName}}Obj := {{if .IsList}}[]{{end}}{{.ScalarType}}{}
			DB.Model(&Result{{$name}}Obj).Association("{{.ScalarName}}").Find(&{{.ScalarName}}Obj)
			Result{{$name}}Obj.{{.ScalarName}} = {{if .IsList}}append(Result{{$name}}Obj.{{.ScalarName}}, {{.ScalarName}}Obj...){{else}}{{.ScalarName}}Obj{{end}}{{end}}{{end}}
			return Result{{$name}}Obj, nil
		},
	}
{{end}}{{end}}
`

var graphqlGoFieldsMutationTemplate = `
var MutationFields graphql.Fields

func init() {
	MutationFields = make(graphql.Fields,1)
	{{template "graphqlFieldsCreate" .}}
}
`

var graphqlGoMutationCreateTemplate = `{{define "graphqlFieldsCreate"}}{{range $key, $value := . }}
	MutationFields["Create{{.Name}}"] = &graphql.Field{
		Type: {{.LowerName}}Type,
		Args: graphql.FieldConfigArgument{
			{{range .Properties}}{{if not .IsObjectType}}"{{.ScalarName}}": &graphql.ArgumentConfig{
					Type: {{index .AltScalarType "graphql-go"}},
			},
			{{end}}{{end}}
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			Query{{.Name}}Obj := {{.Name}}{}
		{{range .Properties}}	{{if not .IsObjectType}}if val, ok := p.Args["{{.ScalarName}}"]; ok {
				Query{{$value.Name}}Obj.{{.ScalarName}} = {{$type := goType .ScalarType}}{{if eq .ScalarType  $type}}val.({{.ScalarType}}){{else}} {{.ScalarType}}(val.({{$type}})){{end}}
			}
		{{end}}{{end}}
			Create{{.Name}}(DB,Query{{.Name}}Obj)
			var Result{{.Name}}Obj {{.Name}}
			Get{{.Name}}(DB,Query{{.Name}}Obj,&Result{{.Name}}Obj)

			return Result{{.Name}}Obj, nil
		},
	}
{{end}}{{end}}
`
