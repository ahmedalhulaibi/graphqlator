package graphqlgo

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
	"unicode"

	"github.com/jinzhu/inflection"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

func (g Gql) GenPackageImports(dbType string, buff *bytes.Buffer) {
	buff.WriteString("package main\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\t\"github.com/graphql-go/graphql\"\n\t\"github.com/graphql-go/handler\"")

	if importVal, exists := g.GraphqlDbTypeImports[dbType]; exists {
		buff.WriteString(importVal)
	}
	buff.WriteString("\n)")
}

func (g Gql) GenGraphqlGoTypeFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	a := []rune(gqlObjectType.Name)
	a[0] = unicode.ToLower(a[0])
	gqlObjectTypeNameLowCamel := string(a)
	buff.WriteString(fmt.Sprintf("\nvar %sType = graphql.NewObject(\n\tgraphql.ObjectConfig{\n\t\tName: \"%s\",\n\t\tFields: graphql.Fields{\n\t\t\t", gqlObjectTypeNameLowCamel, gqlObjectType.Name))

	for _, property := range gqlObjectType.Properties {
		g.GenGraphqlGoTypePropertyFunc(*property, buff)
	}

	buff.WriteString(fmt.Sprintf("\n\t\t},\n\t},\n)\n"))
}

func (g Gql) GenGraphqlGoTypePropertyFunc(gqlObjectProperty substancegen.GenObjectProperty, buff *bytes.Buffer) {
	gqlPropertyTypeName := g.ResolveGraphqlGoFieldType(gqlObjectProperty)
	buff.WriteString(fmt.Sprintf("\n\t\t\t\"%s\": &graphql.Field{\n\t\t\t\tType: %s,\n\t\t\t},", gqlObjectProperty.ScalarName, gqlPropertyTypeName))
}

func (g Gql) ResolveGraphqlGoFieldType(gqlObjectProperty substancegen.GenObjectProperty) string {
	var gqlPropertyTypeName string

	if gqlObjectProperty.IsObjectType {
		a := []rune(inflection.Singular(gqlObjectProperty.ScalarName))
		a[0] = unicode.ToLower(a[0])
		gqlPropertyTypeName = fmt.Sprintf("%sType", string(a))
	} else {
		gqlPropertyTypeName = g.GraphqlDataTypes[gqlObjectProperty.ScalarType]
	}

	if gqlObjectProperty.IsList {
		gqlPropertyTypeName = fmt.Sprintf("graphql.NewList(%s)", gqlPropertyTypeName)
	}

	if !gqlObjectProperty.Nullable {
		gqlPropertyTypeName = fmt.Sprintf("graphql.NewNonNull(%s)", gqlPropertyTypeName)
	}

	return gqlPropertyTypeName
}

func GenGraphqlGoMainFunc(dbType string, connectionString string, gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	buff.WriteString(fmt.Sprintf("\nvar DB *gorm.DB\n\n"))
	buff.WriteString(fmt.Sprintf("\nfunc main() {\n\n\tDB, _ = gorm.Open(\"%s\",\"%s\")\n\tdefer DB.Close()\n\n\t", dbType, connectionString))
	sampleQuery := GenGraphqlGoSampleQuery(gqlObjectTypes)
	buff.WriteString(fmt.Sprintf("\n\tfmt.Println(\"Test with Get\t: curl -g 'http://localhost:8080/graphql?query={%s}'\")", sampleQuery.String()))

	buff.WriteString(GraphqlGoMainConfig)

	buff.WriteString("\n}\n")
}

func GenGraphqlGoFieldsFunc(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer

	buff.WriteString("\n\tvar Fields = graphql.Fields{")
	graphqlQGoFieldsTemplate := "{{$name := .Name}}\n\t\t\"{{.Name}}\": &graphql.Field{\n\t\t\tType: {{.LowerName}}Type,\n\t\t\tResolve: func(p graphql.ResolveParams) (interface{}, error) {\n\t\t\t\t{{.Name}}Obj := {{.Name}}{}\n\t\t\t\tDB.First(&{{.Name}}Obj){{range .Properties}}{{if .IsObjectType}}\n\t\t\t\t{{.ScalarName}}Obj := {{if .IsList}}[]{{end}}{{.ScalarType}}{}\n\t\t\t\tDB.Model(&{{$name}}Obj).Association(\"{{.ScalarName}}\").Find(&{{.ScalarName}}Obj)\n\t\t\t\t{{$name}}Obj.{{.ScalarName}} = {{if .IsList}}append({{$name}}Obj.{{.ScalarName}}, {{.ScalarName}}Obj...){{else}}{{.ScalarName}}Obj{{end}}{{end}}{{end}}\n\t\t\t\treturn {{$name}}Obj, nil\n\t\t\t},\n\t\t},"
	tmpl := template.New("graphqlFields")
	tmpl, err := tmpl.Parse(graphqlQGoFieldsTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	for _, value := range gqlObjectTypes {
		err1 := tmpl.Execute(&buff, value)
		if err1 != nil {
			log.Fatal("Execute: ", err1)
			return buff
		}
	}
	buff.WriteString("\n}\n")
	return buff
}

func GenGraphqlGoSampleQuery(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer
	graphqlQueryTemplate := `{{.Name}} { {{range .Properties}}{{if not .IsObjectType}}{{.ScalarName}},{{end}} {{end}}},`
	tmpl := template.New("graphqlQuery")
	tmpl, err := tmpl.Parse(graphqlQueryTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	for _, value := range gqlObjectTypes {
		err1 := tmpl.Execute(&buff, value)
		if err1 != nil {
			log.Fatal("Execute: ", err1)
			return buff
		}
	}
	bufferString := buff.String()
	bufferString = strings.Replace(bufferString, " ", "", -1)
	buff.Reset()
	buff.WriteString(bufferString)
	return buff
}

/*OutputGraphqlSchema Returns a buffer containing a GraphQL schema in the standard GraphQL schema syntax*/
func OutputGraphqlSchema(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer

	graphqlSchemaTemplate := "type {{.Name}} {\n {{range .Properties}}\t{{.ScalarName}}: {{if .IsList}}[{{.ScalarType}}]{{else}}{{.ScalarType}}{{end}}{{if .Nullable}}{{else}}!{{end}}\n{{end}}}\n"
	tmpl := template.New("graphqlSchema")
	tmpl, err := tmpl.Parse(graphqlSchemaTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	for _, value := range gqlObjectTypes {
		err1 := tmpl.Execute(&buff, value)
		if err1 != nil {
			log.Fatal("Execute: ", err1)
			return buff
		}
	}
	return buff
}
