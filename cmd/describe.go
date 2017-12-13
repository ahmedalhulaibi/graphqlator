package cmd

import (
	"fmt"

	. "github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(describe)
}

var describe = &cobra.Command{
	Use:   "describe [database type] [connection string] [table names...]",
	Short: "List database tables or describe columns of table(s)",
	Long:  `List database tables or describe columns of table(s). If no table names given, it will list tables in database. If table names supplied, the columns of the tables will be described.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			describleTable(args[0], args[1], args[2:len(args)])
		} else {
			describeDatabase(args[0], args[1])
		}
	},
}

func describeDatabase(dbType string, connectionString string) {
	results, err := DescribeDatabase(dbType, connectionString)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		fmt.Println("Database: ", results[0].DatabaseName)
	}
	for _, result := range results {
		fmt.Printf("Table: %s\n", result.TableName)
	}
	fmt.Println("=====================")
}

func describleTable(dbType string, connectionString string, tableNames []string) {
	tableDesc := []ColumnDescription{}
	for _, tableName := range tableNames {
		_results, _ := DescribeTable(dbType, connectionString, tableName)
		tableDesc = append(tableDesc, _results...)
	}
	for _, colDesc := range tableDesc {
		//fmt.Println(colDesc)
		fmt.Println("Table Name:\t", colDesc.TableName)
		fmt.Println("Property Name:\t", colDesc.PropertyName)
		fmt.Println("Property Type:\t", colDesc.PropertyType)
		fmt.Println("Key Type:\t", colDesc.KeyType)
		fmt.Println("Nullable:\t", colDesc.Nullable)
		fmt.Println("------------------------")
	}
	fmt.Println("=====================")
	relationshipDesc := []ColumnRelationship{}

	for _, tableName := range tableNames {
		_results, _ := DescribeTableRelationship(dbType, connectionString, tableName)
		relationshipDesc = append(relationshipDesc, _results...)
	}
	for _, colRel := range relationshipDesc {
		fmt.Println(colRel)
		fmt.Println("Table Name:\t", colRel.TableName)
		fmt.Println("Column Name:\t", colRel.ColumnName)
		fmt.Println("Ref Table Name:\t", colRel.ReferenceTableName)
		fmt.Println("Ref Col Name:\t", colRel.ReferenceColumnName)
	}
	fmt.Println("=====================")
	contraintDesc := []ColumnConstraint{}

	for _, tableName := range tableNames {
		_results, _ := DescribeTableConstraints(dbType, connectionString, tableName)
		contraintDesc = append(contraintDesc, _results...)
	}
	for _, colCon := range contraintDesc {
		fmt.Println(colCon)
		fmt.Println("Table Name:\t", colCon.TableName)
		fmt.Println("Column Name:\t", colCon.ColumnName)
		fmt.Println("Constraint:\t", colCon.ConstraintType)
	}
	fmt.Println("=====================")
}
