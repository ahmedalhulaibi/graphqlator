package substancegen

import (
	"bytes"
)

/*GraphqlatorInterface placeholder comment */
type SubstanceGenInterface interface {
	GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GenObjectType
	ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, genObjects map[string]GenObjectType) map[string]GenObjectType
	OutputCodeFunc(dbType string, connectionString string, gqlObjectTypes map[string]GenObjectType) bytes.Buffer
	GenObjectTypeToStringFunc(GenObjectType, *bytes.Buffer)
	GenObjectPropertyToStringFunc(GenObjectProperty, *bytes.Buffer)
	GenObjectTagToStringFunc(GenObjectTag, *bytes.Buffer)
}

var substanceGenPlugins = make(map[string]SubstanceGenInterface)

/*Register placeholder comment */
func Register(pluginName string, pluginInterface SubstanceGenInterface) {
	substanceGenPlugins[pluginName] = pluginInterface
}

type GenObjectTag map[string][]string

/*GenObjectProperty placeholder comment */
type GenObjectProperty struct {
	ScalarName   string
	ScalarType   string
	IsList       bool
	Nullable     bool
	KeyType      []string
	Tags         GenObjectTag
	IsObjectType bool
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
	substanceGenPlugins[generatorName].OutputCodeFunc(dbType, connectionString, substanceGenPlugins[generatorName].GetObjectTypesFunc(dbType, connectionString, tableNames))
}
