package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "graphqlator",
	Short: "Graphqlator helps you generate a GraphQL type schema. Type 'graphqlator help' to see usage.",
	Long:  `A command line tool that generates a GraphQL type schema from a database table schema. Complete documentation is available at https://github.com/ahmedalhulaibi/graphqlator-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println(cmd.Short)
	},
}
