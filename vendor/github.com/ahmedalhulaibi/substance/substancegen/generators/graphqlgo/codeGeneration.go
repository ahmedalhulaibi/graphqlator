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

func (g Gql) GenerateGraphqlGoTypesFunc(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	for _, value := range gqlObjectTypes {
		for _, propVal := range value.Properties {
			if propVal.IsObjectType {
				a := []rune(inflection.Singular(propVal.ScalarName))
				a[0] = unicode.ToLower(a[0])
				propVal.AltScalarType = fmt.Sprintf("%sType", string(a))
			} else {
				propVal.AltScalarType = g.GraphqlDataTypes[propVal.ScalarType]
			}

			if propVal.IsList {
				propVal.AltScalarType = fmt.Sprintf("graphql.NewList(%s)", propVal.AltScalarType)
			}

			if !propVal.Nullable {
				propVal.AltScalarType = fmt.Sprintf("graphql.NewNonNull(%s)", propVal.AltScalarType)
			}
		}
	}
	graphqlTypesTemplate := "{{range $key, $value := . }}\nvar {{.LowerName}}Type = graphql.NewObject(\n\tgraphql.ObjectConfig{\n\t\tName: \"{{.Name}}\",\n\t\tFields: graphql.Fields{\n\t\t\t{{range .Properties}}\"{{.ScalarNameUpper}}\":&graphql.Field{\n\t\t\t\tType: {{.AltScalarType}},\n\t\t\t},\n\t\t\t{{end}}\n\t\t},\n\t},\n)\n{{end}}"

	tmpl := template.New("graphqlTypes")
	tmpl, err := tmpl.Parse(graphqlTypesTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
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

	graphqlQGoFieldsTemplate := "\n\tvar Fields = graphql.Fields{ {{range $key, $value := . }}{{$name := .Name}}\n\t\t\"{{.Name}}\": &graphql.Field{\n\t\t\tType: {{.LowerName}}Type,\n\t\t\tResolve: func(p graphql.ResolveParams) (interface{}, error) {\n\t\t\t\t{{.Name}}Obj := {{.Name}}{}\n\t\t\t\tDB.First(&{{.Name}}Obj){{range .Properties}}{{if .IsObjectType}}\n\t\t\t\t{{.ScalarName}}Obj := {{if .IsList}}[]{{end}}{{.ScalarType}}{}\n\t\t\t\tDB.Model(&{{$name}}Obj).Association(\"{{.ScalarName}}\").Find(&{{.ScalarName}}Obj)\n\t\t\t\t{{$name}}Obj.{{.ScalarName}} = {{if .IsList}}append({{$name}}Obj.{{.ScalarName}}, {{.ScalarName}}Obj...){{else}}{{.ScalarName}}Obj{{end}}{{end}}{{end}}\n\t\t\t\treturn {{$name}}Obj, nil\n\t\t\t},\n\t\t},{{end}}\n}\n"
	tmpl := template.New("graphqlFields")
	tmpl, err := tmpl.Parse(graphqlQGoFieldsTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	err1 := tmpl.Execute(&buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return buff
	}
	return buff
}

func GenGraphqlGoSampleQuery(gqlObjectTypes map[string]substancegen.GenObjectType) bytes.Buffer {
	var buff bytes.Buffer
	graphqlQueryTemplate := `{{range $key, $value := . }}{{.Name}} { {{range .Properties}}{{if not .IsObjectType}}{{.ScalarName}},{{end}} {{end}}},{{end}}`
	tmpl := template.New("graphqlQuery")
	tmpl, err := tmpl.Parse(graphqlQueryTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	err1 := tmpl.Execute(&buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return buff
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

	graphqlSchemaTemplate := "{{range $key, $value := . }}type {{.Name}} {\n {{range .Properties}}\t{{.ScalarName}}: {{if .IsList}}[{{.ScalarType}}]{{else}}{{.ScalarType}}{{end}}{{if .Nullable}}{{else}}!{{end}}\n{{end}}}\n{{end}}"
	tmpl := template.New("graphqlSchema")
	tmpl, err := tmpl.Parse(graphqlSchemaTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return buff
	}
	//print schema
	err1 := tmpl.Execute(&buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return buff
	}

	return buff
}
