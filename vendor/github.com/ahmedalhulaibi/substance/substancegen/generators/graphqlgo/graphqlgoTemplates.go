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

/*GraphqlGoMainConfig boilerplate string to setup graphl-go and http handler*/
var GraphqlGoMainConfig = `
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
`
