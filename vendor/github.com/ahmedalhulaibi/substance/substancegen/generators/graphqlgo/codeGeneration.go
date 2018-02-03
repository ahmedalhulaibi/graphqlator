package graphqlgo

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

func (g gql) OutputCodeFunc(dbType string, connectionString string, gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer

	g.GenPackageImports(dbType, &buff)
	//print schema
	for _, value := range gqlObjectTypes {
		for _, propVal := range value.Properties {
			propVal.Tags["json"] = append(propVal.Tags["json"], propVal.ScalarName)
		}
		g.GenObjectTypeToStringFunc(value, &buff)
		g.GenGormObjectTableNameOverrideFunc(value, &buff)
		g.GenGraphqlGoTypeFunc(value, &buff)
	}
	buff.WriteString(GraphqlGoExecuteQueryFunc)
	g.GenGraphqlGoMainFunc(dbType, connectionString, gqlObjectTypes, &buff)
	fmt.Print(buff.String())
	return buff
}

func (g gql) GenPackageImports(dbType string, buff *bytes.Buffer) {
	buff.WriteString("package main\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\t\"github.com/graphql-go/graphql\"")

	if importVal, exists := g.GraphqlDbTypeImports[dbType]; exists {
		buff.WriteString(importVal)
	}
	buff.WriteString("\n)")
}

func (g gql) GenObjectTypeToStringFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := strings.TrimSuffix(gqlObjectType.Name, "s")
	buff.WriteString(fmt.Sprintf("\ntype %s struct {\n", gqlObjectTypeNameSingular))
	for _, property := range gqlObjectType.Properties {
		g.GenObjectPropertyToStringFunc(property, buff)
	}
	buff.WriteString("}\n")
}

func (g gql) GenObjectPropertyToStringFunc(gqlObjectProperty substancegen.GenObjectProperty, buff *bytes.Buffer) {

	a := []rune(gqlObjectProperty.ScalarName)
	a[0] = unicode.ToUpper(a[0])
	gqlObjectPropertyNameUpper := string(a)
	if gqlObjectProperty.IsList {
		buff.WriteString(fmt.Sprintf("\t%s\t[]%s\t", gqlObjectPropertyNameUpper, gqlObjectProperty.ScalarType))
	} else {
		buff.WriteString(fmt.Sprintf("\t%s\t%s\t", gqlObjectPropertyNameUpper, gqlObjectProperty.ScalarType))
	}
	g.GenObjectTagToStringFunc(gqlObjectProperty.Tags, buff)
	buff.WriteString("\n")
}

func (g gql) GenObjectTagToStringFunc(genObjectTags substancegen.GenObjectTag, buff *bytes.Buffer) {
	buff.WriteString("`")
	for key, tags := range genObjectTags {
		buff.WriteString(fmt.Sprintf("%s:\"", key))
		for _, tag := range tags {
			buff.WriteString(fmt.Sprintf("%s", tag))
		}
		buff.WriteString("\" ")
	}
	buff.WriteString("`")
}

func (g gql) GenGormObjectTableNameOverrideFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := strings.TrimSuffix(gqlObjectType.Name, "s")
	buff.WriteString(fmt.Sprintf("\nfunc (%s) TableName() string {\n\treturn \"%s\"\n}\n", gqlObjectTypeNameSingular, gqlObjectType.Name))
}

func (g gql) GenGraphqlGoTypeFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	a := []rune(strings.TrimSuffix(gqlObjectType.Name, "s"))
	a[0] = unicode.ToLower(a[0])
	gqlObjectTypeNameLowCamel := string(a)
	gqlObjectTypeNameSingular := strings.TrimSuffix(gqlObjectType.Name, "s")
	buff.WriteString(fmt.Sprintf("\nvar %sType = graphql.NewObject(\n\tgraphql.ObjectConfig{\n\t\tName: \"%s\",\n\t\tFields: graphql.Fields{\n\t\t\t", gqlObjectTypeNameLowCamel, gqlObjectTypeNameSingular))

	for _, property := range gqlObjectType.Properties {
		g.GenGraphqlGoTypePropertyFunc(property, buff)
	}

	buff.WriteString(fmt.Sprintf("\n\t\t},\n\t},\n)\n"))
}

func (g gql) GenGraphqlGoTypePropertyFunc(gqlObjectProperty substancegen.GenObjectProperty, buff *bytes.Buffer) {
	var gqlPropertyTypeName string
	if gqlObjectProperty.IsObjectType {
		a := []rune(strings.TrimSuffix(gqlObjectProperty.ScalarName, "s"))
		a[0] = unicode.ToLower(a[0])
		gqlPropertyTypeName = fmt.Sprintf("%sType", string(a))
	} else {
		gqlPropertyTypeName = g.GraphqlDataTypes[gqlObjectProperty.ScalarType]
	}

	if gqlObjectProperty.IsList {
		buff.WriteString(fmt.Sprintf("\n\t\t\t\"%s\": &graphql.Field{\n\t\t\t\tType: graphql.NewList(%s),\n\t\t\t},", gqlObjectProperty.ScalarName, gqlPropertyTypeName))
	} else {
		buff.WriteString(fmt.Sprintf("\n\t\t\t\"%s\": &graphql.Field{\n\t\t\t\tType: %s,\n\t\t\t},", gqlObjectProperty.ScalarName, gqlPropertyTypeName))

	}
}

func (g gql) GenGraphqlGoMainFunc(dbType string, connectionString string, gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	buff.WriteString(fmt.Sprintf("\nfunc main() {\n\n\tdb, err := gorm.Open(\"%s\",\"%s\")\n\tdefer db.Close()\n\n\t", dbType, connectionString))
	sampleQuery := g.GenGraphqlGoSampleQuery(gqlObjectTypes)
	buff.WriteString(fmt.Sprintf("\n\tfmt.Println(\"Test with Get\t: curl -g 'http://localhost:8080/graphql?query={%s}'\")", sampleQuery.String()))

	buff.WriteString("\n\tfields := graphql.Fields{")
	for _, value := range gqlObjectTypes {
		g.GenGraphqlGoQueryFieldsFunc(value, buff)
	}
	buff.WriteString("\n\t\t}")
	buff.WriteString(GraphqlGoMainConfig)

	buff.WriteString("\n}\n")
}

func (g gql) GenGraphqlGoQueryFieldsFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := strings.TrimSuffix(gqlObjectType.Name, "s")
	a := []rune(strings.TrimSuffix(gqlObjectType.Name, "s"))
	a[0] = unicode.ToLower(a[0])
	gqlObjectTypeNameLowCamel := string(a)
	buff.WriteString(fmt.Sprintf("\n\t\t\"%s\": &graphql.Field{\n\t\t\tType: %sType,", gqlObjectTypeNameSingular, gqlObjectTypeNameLowCamel))
	buff.WriteString(fmt.Sprintf("\n\t\t\tResolve: func(p graphql.ResolveParams) (interface{}, error) {"))
	buff.WriteString(fmt.Sprintf("\n\t\t\t\t%s := %s{}", gqlObjectTypeNameLowCamel, gqlObjectTypeNameSingular))
	buff.WriteString(fmt.Sprintf("\n\t\t\t\tdb.First(&%s)", gqlObjectTypeNameLowCamel))

	for _, propVal := range gqlObjectType.Properties {
		if propVal.IsObjectType {
			a := []rune(propVal.ScalarName)
			a[0] = unicode.ToLower(a[0])
			propValNameLowCamel := string(a)
			b := []rune(propVal.ScalarName)
			b[0] = unicode.ToUpper(b[0])
			propValNameUpperCamel := string(b)
			if propVal.IsList {
				buff.WriteString(fmt.Sprintf("\n\t\t\t\t%s := []%s{}", propValNameLowCamel, propVal.ScalarType))

				buff.WriteString(fmt.Sprintf("\n\t\t\t\tdb.Model(&%s).Association(\"%s\").Find(&%s)", gqlObjectTypeNameLowCamel, propVal.ScalarName, propValNameLowCamel))

				buff.WriteString(fmt.Sprintf("\n\t\t\t\t%s.%s = append(%s.%s, %s...)", gqlObjectTypeNameLowCamel, propValNameUpperCamel, gqlObjectTypeNameLowCamel, propValNameUpperCamel, propValNameLowCamel))
			} else {
				buff.WriteString(fmt.Sprintf("\n\t\t\t\t%s := %s{}", propValNameLowCamel, propVal.ScalarType))

				buff.WriteString(fmt.Sprintf("\n\t\t\t\tdb.Model(&%s).Association(\"%s\").Find(&%s)", gqlObjectTypeNameLowCamel, propVal.ScalarName, propValNameLowCamel))

				buff.WriteString(fmt.Sprintf("\n\t\t\t\t%s.%s = %s", gqlObjectTypeNameLowCamel, propValNameUpperCamel, propValNameLowCamel))
			}
		}
	}
	buff.WriteString(fmt.Sprintf("\n\t\t\t\treturn %s, nil", gqlObjectTypeNameLowCamel))
	buff.WriteString("\n\t\t\t},")
	buff.WriteString("\n\t\t},")
}

func (g gql) GenGraphqlGoSampleQuery(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer
	for _, gqlObjectType := range gqlObjectTypes {
		g.GenGraphlGoSampleObjectQuery(gqlObjectTypes, gqlObjectType, &buff)
	}
	return buff
}

func (g gql) GenGraphlGoSampleObjectQuery(gqlObjectTypes map[string]substancegen.GenObjectType, gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := strings.TrimSuffix(gqlObjectType.Name, "s")
	buff.WriteString(fmt.Sprintf("%s{", gqlObjectTypeNameSingular))
	for _, propVal := range gqlObjectType.Properties {
		if !propVal.IsObjectType {
			buff.WriteString(fmt.Sprintf("%s,", propVal.ScalarName))
		}
	}
	buff.WriteString("},")
}
