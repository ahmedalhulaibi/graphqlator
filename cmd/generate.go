package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ahmedalhulaibi/substance/substancegen/generators/gorm"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/gostruct"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/graphqlgo"

	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/spf13/cobra"
)

var updateSchema bool
var updateGormQueries bool
var updateGqlFields bool
var updateGqlTypes bool
var updateModel bool
var updateMain bool
var updateAll bool

func init() {
	generate.Flags().BoolVarP(&updateSchema, "update-schema", "u", false, "update and overwrite schema.graphql")
	generate.Flags().BoolVarP(&updateGormQueries, "update-gormQueries", "q", false, "update and overwrite gormQueries.go")
	generate.Flags().BoolVarP(&updateGqlFields, "update-gqlFields", "g", false, "update and overwrite graphqlFields.go")
	generate.Flags().BoolVarP(&updateGqlTypes, "update-gqlTypes", "t", false, "update and overwrite graphqlTypes.go")
	generate.Flags().BoolVarP(&updateMain, "update-main", "m", false, "update and overwrite main.go")
	generate.Flags().BoolVarP(&updateModel, "update-gormStruct", "s", false, "update and overwrite model.go")
	generate.Flags().BoolVarP(&updateAll, "update-all", "a", false, "update and overwrite all files")
	RootCmd.AddCommand(generate)
}

var generate = &cobra.Command{
	Use:   "generate",
	Short: "Generate GraphQL-Go API implementation using grapqhlator-pkg.json.",
	Long: `Generate GraphQL-Go API implementation from database information schema and tables defined in grapqhlator-pkg.json
Run 'graphqlator init' before running 'graphqlator generate'`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		gqlPkg := getGraphqlatorPkgFile()
		gqlGen := substancegen.SubstanceGenPlugins["graphql-go"].(graphqlgo.Gql)
		gqlObjectTypes := gqlGen.GetObjectTypesFunc(gqlPkg.DatabaseType, gqlPkg.ConnectionString, gqlPkg.TableNames)
		gqlGen.AddJSONTagsToProperties(gqlObjectTypes)

		if gqlPkg.GenMode == "graphql-go" {
			if !updateGormQueries && !updateGqlFields && !updateGqlTypes && !updateModel && !updateSchema && !updateMain && !updateAll {
				cmd.Help()
			}
			if updateMain || updateAll {
				mainFile := createFile("main.go", true)

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

			if updateGqlTypes || updateAll {
				graphqlTypesFile := createFile("graphqlTypes.go", true)

				var graphqlTypesFileBuff bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &graphqlTypesFileBuff)
				for _, value := range gqlObjectTypes {
					gqlGen.GenGraphqlGoTypeFunc(value, &graphqlTypesFileBuff)
				}
				_, err := graphqlTypesFile.Write(graphqlTypesFileBuff.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				graphqlTypesFile.Close()
			}

			if updateModel || updateAll {
				dataModelFile := createFile("model.go", true)

				var dataModelFileBuff bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &dataModelFileBuff)
				for _, value := range gqlObjectTypes {
					gostruct.GenObjectTypeToStructFunc(value, &dataModelFileBuff)
					gorm.GenGormObjectTableNameOverrideFunc(value, &dataModelFileBuff)
				}
				_, err := dataModelFile.Write(dataModelFileBuff.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				dataModelFile.Close()
			}

			if updateGqlFields || updateAll {
				gqlFieldsFile := createFile("graphqlFields.go", true)

				var gqlFieldsFileBuff bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &gqlFieldsFileBuff)
				gqlGen.GenGraphqlGoRootQueryFunc(gqlObjectTypes, &gqlFieldsFileBuff)
				_, err := gqlFieldsFile.Write(gqlFieldsFileBuff.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				gqlFieldsFile.Close()
			}

			if updateGormQueries || updateAll {
				gormQueriesFile := createFile("gormQueries.go", true)

				var gormQueriesFileBuff bytes.Buffer
				gqlGen.GenPackageImports(gqlPkg.DatabaseType, &gormQueriesFileBuff)
				for _, value := range gqlObjectTypes {
					gorm.GenObjectGormCrud(value, &gormQueriesFileBuff)
				}
				_, err := gormQueriesFile.Write(gormQueriesFileBuff.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				gormQueriesFile.Close()
			}

			if updateSchema || updateAll {
				graphqlSchemaFile := createFile("schema.graphql", true)
				graphqlSchemaFileBuffer := gqlGen.OutputGraphqlSchema(gqlObjectTypes)
				_, err := graphqlSchemaFile.Write(graphqlSchemaFileBuffer.Bytes())
				if err != nil {
					fmt.Println(err.Error())
				}
				graphqlSchemaFile.Close()
			}

			{
				formatFile := createFile("format.sh", true)
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

func createFile(filepath string, overwrite bool) *os.File {
	file, err := os.Open(filepath)
	if err == nil {
		if overwrite {
			file.Close()
			os.Remove(file.Name())
		} else {
			return file
		}
	}
	file, err = os.Create(filepath)
	if err != nil {
		check(err, fmt.Sprintf("Problem creating file %s", filepath))
	}
	return file
}
