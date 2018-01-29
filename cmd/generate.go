package cmd

import (
	"github.com/ahmedalhulaibi/substance/substancegen"
	_ "github.com/ahmedalhulaibi/substance/substancegen/generators/graphqlgo"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generate)
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
	substancegen.Generate("graphql-go", dbType, connectionString, tableNames)
}
