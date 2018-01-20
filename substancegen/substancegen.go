package substancegen

/*GraphqlatorInterface placeholder comment */
type SubstanceGenInterface interface {
	GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GqlObjectType
	ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjects map[string]GqlObjectType) map[string]GqlObjectType
	OutputCodeFunc(map[string]GqlObjectType)
}

var graphqlatorPlugins = make(map[string]SubstanceGenInterface)

/*Register placeholder comment */
func Register(pluginName string, pluginInterface SubstanceGenInterface) {
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
	graphqlatorPlugins[gqlType].OutputCodeFunc(graphqlatorPlugins[gqlType].GetObjectTypesFunc(dbType, connectionString, tableNames))
}
