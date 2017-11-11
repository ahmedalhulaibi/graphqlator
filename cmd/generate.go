package cmd

import (
	"database/sql"
	"fmt"

	_ "gopkg.in/mgo.v2"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generate)
}

/*
  Below query returns forein key restraints on a table
  SELECT
  TABLE_NAME,
  COLUMN_NAME,
  CONSTRAINT_NAME,
  REFERENCED_TABLE_NAME,
  REFERENCED_COLUMN_NAME
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
  REFERENCED_TABLE_NAME = 'My_Table';

  When using DESCRIBE table_name the Key column tells us if it is unique or not
         Field :  PersonID
         Type :  int(11)
         Null :  YES
         Key :  MUL
         Default : NULL
         Extra :
*/

type columnDescriptions struct {
	TableName           string
	PropertyName        string
	PropertyType        string
	KeyType             string
	ReferenceTableName  string
	ReferenceColumnName string
}

var schemaTable []columnDescriptions

var generate = &cobra.Command{
	Use:   "generate [database type] [connection string] [table name or collection name...]",
	Short: "Generate GraphQL type schema from database collection or table",
	Long:  `Describe database listing tables/collections. If table name/collection name is supplied, the fields of the table/document will be described.`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(args)
		switch args[0] {
		case "mysql":
			populateSchemaTable(args[0], args[1], args[2:len(args)])
			fmt.Println(schemaTable)
			break
		}
	},
}

func populateSchemaTable(dbType string, connectionString string, tableNames []string) {
	for i := range tableNames {
		processTable(dbType, connectionString, tableNames[i])
	}
}

func processTable(dbType string, connectionString string, tableName string) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	// Connect and check the server version

	query := fmt.Sprintf("DESCRIBE %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		newCol := columnDescriptions{TableName: tableName}
		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				fmt.Println("\t", columns[i], ": NULL")
			case []byte:
				fmt.Println("\t", columns[i], ": ", string(value.([]byte)))
				switch columns[i] {
				case "Field":
					newCol.PropertyName = string(value.([]byte))
				case "Type":
					newCol.PropertyType = string(value.([]byte))
				case "Key":
					newCol.KeyType = string(value.([]byte))
				}
			default:
				fmt.Println("\t", columns[i], ": ", value)
			}
		}
		schemaTable = append(schemaTable, newCol)
		fmt.Println("-----------------------------------")
	}
}
