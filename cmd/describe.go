package cmd

import (
	"fmt"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/sqlsubstance"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(describe)
}

var describe = &cobra.Command{
	Use:   "describe [database type] [connection string] [table name]",
	Short: "Describe database or table",
	Long:  `Describe database listing tables. If table name is supplied, the fields of the table will be described.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(args)
		if len(args) > 2 {
			describleTable(args[0], args[1], args[2:len(args)])
		} else {
			describeDatabase(args[0], args[1])
		}
	},
}

func describeDatabase(dbType string, connectionString string) {
	fmt.Println("------TABLES FOUND------")
	results, err := sqlsubstance.DescribeDatabase(dbType, connectionString)
	if err != nil {
		panic(err)
	}
	for _, result := range results {

		fmt.Println(result)
	}
	fmt.Println("------------------------")
}

func describleTable(dbType string, connectionString string, tableNames []string) {
	tableDesc := []sqlsubstance.ColumnDescription{}
	for _, tableName := range tableNames {
		_results, _ := sqlsubstance.DescribeTable(dbType, connectionString, tableName)
		tableDesc = append(tableDesc, _results...)
	}
	for _, colDesc := range tableDesc {
		fmt.Println(colDesc)
		fmt.Println("Table Name:\t", colDesc.TableName)
		fmt.Println("Property Name:\t", colDesc.PropertyName)
		fmt.Println("Property Type:\t", colDesc.PropertyType)
		fmt.Println("Key Type:\t", colDesc.KeyType)
		fmt.Println("Nullable:\t", colDesc.Nullable)
	}
	fmt.Println("------------------------")
	relationshipDesc := []sqlsubstance.ColumnRelationship{}

	for _, tableName := range tableNames {
		_results, _ := sqlsubstance.DescribeTableRelationship(dbType, connectionString, tableName)
		relationshipDesc = append(relationshipDesc, _results...)
	}
	for _, colRel := range relationshipDesc {
		fmt.Println(colRel)
		fmt.Println("Table Name:\t", colRel.TableName)
		fmt.Println("Column Name:\t", colRel.ColumnName)
		fmt.Println("Ref Table Name:\t", colRel.ReferenceTableName)
		fmt.Println("Ref Col Name:\t", colRel.ReferenceColumnName)
	}
	fmt.Println("------------------------")
}
