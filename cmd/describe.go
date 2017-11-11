package cmd

import (
	"database/sql"
	"fmt"

	_ "gopkg.in/mgo.v2"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(describe)
}

var describe = &cobra.Command{
	Use:   "describe [database type] [connection string] [table name or collection name]",
	Short: "Describe database, collection, or table",
	Long:  `Describe database listing tables/collections. If table name/collection name is supplied, the fields of the table/document will be described.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(args)
		switch args[0] {
		case "mysql":
			if len(args) > 2 {
				describleTable(args[0], args[1], args[2:len(args)])
			} else {
				describeDatabase(args[0], args[1])
			}
			break
		}
	},
}

func describeDatabase(dbType string, connectionString string) {
	fmt.Println("------TABLES FOUND------")
	runQuery(dbType, connectionString, "SHOW TABLES")
	fmt.Println("------------------------")
}

func describleTable(dbType string, connectionString string, tableNames []string) {
	for _, tableName := range tableNames {
		fmt.Printf("------TABLE %s DESCRIPTION------\n", tableName)
		query := fmt.Sprintf("DESCRIBE %s", tableName)
		runQuery(dbType, connectionString, query)
		fmt.Println("------------------------")
	}
}

func runQuery(dbType string, connectionString string, queryString string) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		panic(err.Error())
	}
	// Connect and check the server version

	rows, err := db.Query(queryString)
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

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				fmt.Println("\t", columns[i], ": NULL")

			case []byte:
				fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

			default:
				fmt.Println("\t", columns[i], ": ", value)
			}
		}
		fmt.Println("-----------------------------------")
	}
}
