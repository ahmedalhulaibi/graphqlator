package substancegen

/*GraphqlatorInterface placeholder comment */
type SubstanceGenInterface interface {
	GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GenObjectType
	ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, genObjects map[string]GenObjectType) map[string]GenObjectType
	OutputCodeFunc(map[string]GenObjectType)
}

var substanceGenPlugins = make(map[string]SubstanceGenInterface)

/*Register placeholder comment */
func Register(pluginName string, pluginInterface SubstanceGenInterface) {
	substanceGenPlugins[pluginName] = pluginInterface
}

/*GenObjectProperty placeholder comment */
type GenObjectProperty struct {
	ScalarName string
	ScalarType string
	IsList     bool
	Nullable   bool
	KeyType    string
}

/*GenObjectProperties placeholder comment */
type GenObjectProperties map[string]GenObjectProperty

/*GenObjectType placeholder comment */
type GenObjectType struct {
	Name       string
	Properties GenObjectProperties
}

/*Generate placeholder comment */
func Generate(generatorName string, dbType string, connectionString string, tableNames []string) {
	substanceGenPlugins[generatorName].OutputCodeFunc(substanceGenPlugins[generatorName].GetObjectTypesFunc(dbType, connectionString, tableNames))
}
