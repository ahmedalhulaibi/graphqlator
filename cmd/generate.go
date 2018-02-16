package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ahmedalhulaibi/substance/substancegen/generators/graphqlgo"

	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generate)
}

var generate = &cobra.Command{
	Use:   "generate",
	Short: "Generate GraphQL type schema using grapqhlator-pkg.json.",
	Long: `Generate GraphQL type schema from database information schema and tables defined in grapqhlator-pkg.json
Run 'graphqlator init' before running 'graphqlator generate'`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		gqlPkg := getGraphqlatorPkgFile()
		gqlGen := substancegen.SubstanceGenPlugins["graphql-go"].(graphqlgo.Gql)
		gqlObjectTypes := gqlGen.GetObjectTypesFunc(gqlPkg.DatabaseType, gqlPkg.ConnectionString, gqlPkg.TableNames)
		gqlGen.AddJSONTagsToProperties(gqlObjectTypes)

		if gqlPkg.GenMode == "graphql-go" {
			{
				dataModelFile := createFile("model.go")

				var dataModelFileBuff bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &dataModelFileBuff)
				for _, value := range gqlObjectTypes {
					gqlGen.GenObjectTypeToStringFunc(value, &dataModelFileBuff)
					gqlGen.GenGormObjectTableNameOverrideFunc(value, &dataModelFileBuff)
					gqlGen.GenGraphqlGoTypeFunc(value, &dataModelFileBuff)
				}
				_, err := dataModelFile.Write(dataModelFileBuff.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				dataModelFile.Close()
			}

			{
				mainFile := createFile("main.go")

				var mainFileBuffer bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &mainFileBuffer)
				mainFileBuffer.WriteString(graphqlgo.GraphqlGoExecuteQueryFunc)
				gqlGen.GenGraphqlGoMainFunc(gqlPkg.DatabaseType, gqlPkg.ConnectionString, gqlObjectTypes, &mainFileBuffer)
				_, err := mainFile.Write(mainFileBuffer.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				mainFile.Close()
			}

			{
				graphqlSchemaFile := createFile("schema.graphql")
				graphqlSchemaFileBuffer := gqlGen.OutputGraphqlSchema(gqlObjectTypes)
				_, err := graphqlSchemaFile.Write(graphqlSchemaFileBuffer.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				graphqlSchemaFile.Close()
			}
		}

		{
			formatFile := createFile("format.sh")
			var formatFileBuffer bytes.Buffer
			formatFileBuffer.WriteString("#!usr/bin/env bash\n")
			formatFileBuffer.WriteString("gofmt -w ./*.go\n")
			formatFileBuffer.WriteString("goreturns -w ./*.go\n")
			_, err := formatFile.Write(formatFileBuffer.Bytes())
			if err != nil {
				fmt.Println(err.Error())
			}
			formatFile.Close()
		}
		check(exec.Command("bash", "format.sh").Run(), "format failed")

	},
}

func getGraphqlatorPkgFile() gqlpackage {
	f, err := ioutil.ReadFile("./graphqlator-pkg.json")
	check(err, "Problem opening graphqlator-pkg.json make sure it exists.")
	var gqlPkg gqlpackage
	json.Unmarshal(f, &gqlPkg)
	return gqlPkg
}

func createFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err == nil {
		file.Close()
		os.Remove(file.Name())
	}
	file, err = os.Create(filepath)
	if err != nil {
		check(err, fmt.Sprintf("Problem creating file %s", filepath))
	}
	return file
}
