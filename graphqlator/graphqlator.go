package graphqlator

/*GraphqlatorInterface placeholder comment */
type GraphqlatorInterface interface {
	GetGqlObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GqlObjectType
	ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjects map[string]GqlObjectType) map[string]GqlObjectType
	OutputCodeFunc(map[string]GqlObjectType)
}

var graphqlatorPlugins = make(map[string]GraphqlatorInterface)

/*Register placeholder comment */
func Register(pluginName string, pluginInterface GraphqlatorInterface) {
	graphqlatorPlugins[pluginName] = pluginInterface
}

/*GqlObjectProperty placeholder comment */
type GqlObjectProperty struct {
	ScalarName string
	ScalarType string
	IsList     bool
	Nullable   bool
	KeyType    string
}

/*GqlObjectProperties placeholder comment */
type GqlObjectProperties map[string]GqlObjectProperty

/*GqlObjectType placeholder comment */
type GqlObjectType struct {
	Name       string
	Properties GqlObjectProperties
}

/*Graphqlate placeholder comment */
func Graphqlate(gqlType string, dbType string, connectionString string, tableNames []string) {
	graphqlatorPlugins[gqlType].OutputCodeFunc(graphqlatorPlugins[gqlType].GetGqlObjectTypesFunc(dbType, connectionString, tableNames))
}
