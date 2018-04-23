package graphqlgo

import (
	"bytes"

	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/gorm"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/gostruct"
)

func init() {
	gqlPlugin := Gql{}
	gqlPlugin.GraphqlDataTypes = make(map[string]string)
	InitGraphqlDataTypes(&gqlPlugin)
	gqlPlugin.GraphqlDbTypeImports = make(map[string]string)
	gqlPlugin.GraphqlDbTypeImports["mysql"] = "\n\t\"github.com/jinzhu/gorm\"\n\t_ \"github.com/jinzhu/gorm/dialects/mysql\""
	gqlPlugin.GraphqlDbTypeImports["postgres"] = "\n\t\"github.com/jinzhu/gorm\"\n\t_ \"github.com/jinzhu/gorm/dialects/postgres\""
	substancegen.Register("graphql-go", gqlPlugin)
}

/*InitGraphqlDataTypes initializes gqlPlugin data for go to graphql-go type mapping*/
func InitGraphqlDataTypes(gqlPlugin *Gql) {
	gqlPlugin.GraphqlDataTypes = make(map[string]string)
	gqlPlugin.GraphqlDataTypes["int"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int8"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int16"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int32"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int64"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint8"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint16"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint32"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint64"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["byte"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["rune"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["bool"] = "graphql.Boolean"
	gqlPlugin.GraphqlDataTypes["string"] = "graphql.String"
	gqlPlugin.GraphqlDataTypes["float32"] = "graphql.Float"
	gqlPlugin.GraphqlDataTypes["float64"] = "graphql.Float"
}

/*Gql plugin struct implementation of substancegen Generator interface*/
type Gql struct {
	Name                  string
	GraphqlDataTypes      map[string]string
	GraphqlDbTypeGormFlag map[string]bool
	GraphqlDbTypeImports  map[string]string
}

/*OutputCodeFunc returns a buffer with a graphql-go implementation in a single file*/
func (g Gql) OutputCodeFunc(dbType string, connectionString string, gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer

	g.GenPackageImports(dbType, &buff)
	//print schema
	substancegen.AddJSONTagsToProperties(gqlObjectTypes)
	gostruct.GenObjectTypeToStructFunc(gqlObjectTypes, &buff)
	for _, value := range gqlObjectTypes {
		gorm.GenGormObjectTableNameOverrideFunc(value, &buff)
	}
	g.GenerateGraphqlGoTypesFunc(gqlObjectTypes, &buff)
	buff.WriteString(GraphqlGoExecuteQueryFunc)
	var graphqlFieldsBuff bytes.Buffer
	GenGraphqlGoFieldsFunc(gqlObjectTypes, &graphqlFieldsBuff)
	buff.Write(graphqlFieldsBuff.Bytes())
	GenGraphqlGoMainFunc(dbType, connectionString, gqlObjectTypes, &buff)
	return buff
}
