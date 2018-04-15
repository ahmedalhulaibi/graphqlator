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

/*GenPackageImports writes a predefined package and import statement to a buffer*/
func (g Gql) GenPackageImports(dbType string, buff *bytes.Buffer) {
	buff.WriteString("package main\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n\t\"log\"\n\t\"net/http\"\n\t\"github.com/graphql-go/graphql\"\n\t\"github.com/graphql-go/handler\"")

	if importVal, exists := g.GraphqlDbTypeImports[dbType]; exists {
		buff.WriteString(importVal)
	}
	buff.WriteString("\n)")
}

/*GenerateGraphqlGoTypesFunc takes a map of gen objects and outputs graphql-go types to a buffer*/
func (g Gql) GenerateGraphqlGoTypesFunc(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	for _, value := range gqlObjectTypes {
		for _, propVal := range value.Properties {
			if propVal.IsObjectType {
				a := []rune(inflection.Singular(propVal.ScalarName))
				a[0] = unicode.ToLower(a[0])
				propVal.AltScalarType["graphql-go"] = fmt.Sprintf("%sType", string(a))
			} else {
				propVal.AltScalarType["graphql-go"] = g.GraphqlDataTypes[propVal.ScalarType]
			}

			if propVal.IsList {
				propVal.AltScalarType["graphql-go"] = fmt.Sprintf("graphql.NewList(%s)", propVal.AltScalarType["graphql-go"])
			}

			if !propVal.Nullable {
				propVal.AltScalarType["graphql-go"] = fmt.Sprintf("graphql.NewNonNull(%s)", propVal.AltScalarType["graphql-go"])
			}
		}
	}
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

/*GenGraphqlGoMainFunc generates the main function (entrypoint) for the graphql-go server*/
func GenGraphqlGoMainFunc(dbType string, connectionString string, gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	var sampleQuery bytes.Buffer
	GenGraphqlGoSampleQuery(gqlObjectTypes, &sampleQuery)

	// buff.WriteString(fmt.Sprintf("\nvar DB *gorm.DB\n\n"))
	// buff.WriteString(fmt.Sprintf("\nfunc main() {\n\n\tDB, _ = gorm.Open(\"%s\",\"%s\")\n\tdefer DB.Close()\n\n\t", dbType, connectionString))

	// buff.WriteString(fmt.Sprintf("\n\tfmt.Println(\"Test with Get\t: curl -g 'http://localhost:8080/graphql?query={%s}'\")", sampleQuery.String()))

	// buff.WriteString(graphqlGoMainFunc)

	// buff.WriteString("\n}\n")

	mainData := struct {
		DbType           string
		ConnectionString string
		SampleQuery      string
	}{
		dbType,
		connectionString,
		sampleQuery.String(),
	}
	tmpl := template.New("graphqlGoMainFunc")
	tmpl, err := tmpl.Parse(graphqlGoMainFunc)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, mainData)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
	}
}

/*GenGraphqlGoFieldsFunc generates a basic graphql-go queries
to retrieve the first element of each object type (and its associations) from a database*/
func GenGraphqlGoFieldsFunc(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	tmpl := template.New("graphqlFields")
	tmpl, err := tmpl.Parse(graphqlGoFieldsTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
	}
}

/*GenGraphqlGoSampleQuery generates a sample graphql query based on the given objects*/
func GenGraphqlGoSampleQuery(gqlObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	tmpl := template.New("graphqlQuery")
	tmpl, err := tmpl.Parse(graphqlQueryTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	//print schema
	err1 := tmpl.Execute(buff, gqlObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}

	bufferString := buff.String()
	bufferString = strings.Replace(bufferString, " ", "", -1)
	buff.Reset()
	buff.WriteString(bufferString)
}
