package cmd

import (
	_ "github.com/ahmedalhulaibi/go-graphqlator-cli/gqlschema"
	"github.com/ahmedalhulaibi/go-graphqlator-cli/substancegen"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generate)
}

type gqlObjectProperty struct {
	scalarName string
	scalarType string
	isList     bool
	nullable   bool
	keyType    string
}

type gqlObjectProperties map[string]gqlObjectProperty

type gqlObjectType struct {
	name       string
	properties gqlObjectProperties
}

var generate = &cobra.Command{
	Use:   "generate [database type] [connection string] [table names...]",
	Short: "Generate GraphQL type schema from database table(s).",
	Long:  `Generate GraphQL type schema from database table(s).`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "mariadb":
			args[0] = "mysql"
			break
		}
		generateGqlSchema(args[0], args[1], args[2:len(args)])
	},
}

func generateGqlSchema(dbType string, connectionString string, tableNames []string) {
	substancegen.Graphqlate("graphql-go", dbType, connectionString, tableNames)
}
