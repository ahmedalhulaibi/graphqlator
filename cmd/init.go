package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

func check(e error, message string) {
	if e != nil {
		fmt.Println(message)
		fmt.Println(e.Error())
		panic(e)
	}
}

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Generate graphqlator config file.",
	Long:  `Walks you through the generation of a graphqlator config file.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open("graphqlator-pkg.json")
		if err == nil {
			check(fmt.Errorf(""), "graphqlator-pkg.json File already exists")
		}

		f, err = os.Create("graphqlator-pkg.json")
		defer f.Close()
		check(err, "Failed to create file.")

		reader := bufio.NewReader(os.Stdin)

		newGqlPackage := gqlpackage{}

		if len(args) > 0 {
			newGqlPackage.ProjectName = args[0]
		} else {
			fmt.Printf("Input project name (enter to continue): ")
			projectName, err := reader.ReadString('\n')
			check(err, "Problem reading input project name")
			projectName = strings.Replace(projectName, "\n", "", -1)
			projectName = strings.Replace(projectName, " ", "-", -1)
			newGqlPackage.ProjectName = projectName
		}

		//input database_type
		{
			fmt.Printf("Input database type e.g. mysql (enter to continue): ")
			dbType, err := reader.ReadString('\n')
			check(err, "Problem reading database type")
			dbType = strings.Replace(dbType, "\n", "", -1)
			newGqlPackage.DatabaseType = strings.ToLower(dbType)
		}

		//input connection_string
		{
			fmt.Printf(`Input Database Connection String
				MySql Example: username:password@tcp(localhost:3306)/schemaname
				Postgresql Example: postgres://username:password@localhost:5432/dbname
Input database connection string (enter to continue): `)
			connString, err := reader.ReadString('\n')
			check(err, "Problem reading connection string")
			connString = strings.Replace(connString, "\n", "", -1)
			newGqlPackage.ConnectionString = strings.ToLower(connString)
		}

		//input git_repo link
		{
			fmt.Printf("Input git repo url (enter to continue): ")
			gitRepo, err := reader.ReadString('\n')
			check(err, "Problem reading git repo")
			gitRepo = strings.Replace(gitRepo, "\n", "", -1)
			newGqlPackage.GitRepo = gitRepo
		}

		//input table_names
		{
			fmt.Println("Input table names - Must be EXACT spelling and case (enter without input to skip)")
			tableNames := []string{}
			i := 1
			for {
				fmt.Printf("Table #%d : ", i)
				tableName, err := reader.ReadString('\n')
				check(err, "Problem reading table name")
				tableName = strings.Replace(tableName, "\n", "", -1)
				tableName = strings.Replace(tableName, " ", "", -1)
				if tableName != "" {
					tableNames = append(tableNames, tableName)
				} else {
					break
				}
				i++
			}
			newGqlPackage.TableNames = tableNames
		}

		{
			fmt.Println("What graphql implementation will you be using? gqlgen or graphql-go?")
			fmt.Println("gqlgen - https://github.com/vektah/gqlgen")
			fmt.Println("graphql-go - https://github.com/graphql-go/graphql")
			genMode, err := reader.ReadString('\n')
			check(err, "Problem reading gen mode")
			genMode = strings.Replace(genMode, "\n", "", -1)
			newGqlPackage.GenMode = genMode
		}

		newGqlPackageContent, err := json.Marshal(newGqlPackage)
		check(err, "Failed to encode JSON to byte string")
		f.Write(newGqlPackageContent)

	},
}
