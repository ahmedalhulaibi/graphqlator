package substancegen

import (
	"bytes"
)

/*GeneratorInterface describes the implementation required to generate code from substance objects*/
type GeneratorInterface interface {
	GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GenObjectType
	ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, genObjects map[string]GenObjectType) map[string]GenObjectType
	OutputCodeFunc(dbType string, connectionString string, gqlObjectTypes map[string]GenObjectType) bytes.Buffer
}

/*SubstanceGenPlugins is a map storing a reference to the current plugins
Key: pluginName
Value: reference to an implementation of SubstanceGenInterface*/
var SubstanceGenPlugins = make(map[string]GeneratorInterface)

/*Register registers a GeneratorInterface plugin */
func Register(pluginName string, pluginInterface GeneratorInterface) {
	SubstanceGenPlugins[pluginName] = pluginInterface
}

/*GenObjectTag stores a key value pair of go struct a tag and their value(s)
Example:
Key: gorm
Tabs: {'primary_key','column_name'}*/
type GenObjectTag map[string][]string

/*TODO: Create new type to store KeyType to map [string]string
This will require changes in generators/graphqlgo pkg
This will require changes in generators/gorm.go pkg + gorm_test.go
This will require changes in generators/gostruct_test.go*/

/*GenObjectProperty represents a property of an object (aka a field of a struct) */
type GenObjectProperty struct {
	ScalarName      string `json:"scalarName"`
	ScalarNameUpper string
	ScalarType      string       `json:"scalarType"`
	IsList          bool         `json:"isList"`
	Nullable        bool         `json:"nullable"`
	KeyType         []string     `json:"keyType"`
	Tags            GenObjectTag `json:"tags"`
	IsObjectType    bool         `json:"isObjectType"`
}

/*GenObjectProperties a type defining a map of GenObjectProperty
Key: PropertyName
Value: GenObjectProperty */
type GenObjectProperties map[string]*GenObjectProperty

/*GenObjectType represents an object (aka a struct) */
type GenObjectType struct {
	Name            string `json:"objectName"`
	SourceTableName string `json:"sourceTableName"`
	LowerName       string
	Properties      GenObjectProperties `json:"properties"`
}

/*Generate is a one stop function to quickly generate code */
func Generate(generatorName string, dbType string, connectionString string, tableNames []string) bytes.Buffer {
	return SubstanceGenPlugins[generatorName].OutputCodeFunc(dbType, connectionString, SubstanceGenPlugins[generatorName].GetObjectTypesFunc(dbType, connectionString, tableNames))
}
