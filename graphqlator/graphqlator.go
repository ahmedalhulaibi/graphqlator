package graphqlator

type GraphqlatorInterface interface {
}

type gqlObjectProperty struct {
	scalarName string
	scalarType string
	isList     bool
	nullable   bool
	keyType    string
}

type gqlObjectProperties map[string]gqlObjectProperty

type gqlObjectType struct {
	name       string
	properties gqlObjectProperties
}
