package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Graphqlator",
	Long:  `All software has versions. This is Graphqlators's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Graphqlator GraphQL Generator v0.1")
	},
}
