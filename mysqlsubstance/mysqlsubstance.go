package mysqlsubstance

import (
	"database/sql"
	"fmt"
)

/*ColumnDescription Structure to store properties of each column in a table */
type ColumnDescription struct {
	DatabaseName string
	TableName    string
	PropertyName string
	PropertyType string
	KeyType      string
	Nullable     bool
}

/*ColumnRelationship Structure to store relationships between tables*/
type ColumnRelationship struct {
	TableName           string
	ColumnName          string
	ReferenceTableName  string
	ReferenceColumnName string
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func GetCurrentDatabaseName(dbType string, connectionString string) (string, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return "nil", err
	}
	var query string
	switch dbType {
	case "mysql":
		query = "SELECT DATABASE()"
		break
	case "mariadb":
		query = "SELECT DATABASE()"
		break
	}
	rows, err := db.Query(query)
	if err != nil {
		return "nil", err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return "nil", err
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var returnValue string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return "nil", err
		}

		// Print data
		for _, value := range values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")
				return "nil", err
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))
				returnValue = string(value.([]byte))
			default:
				//fmt.Println("\t", columns[i], ": ", value)
				returnValue = string(value.([]byte))
			}
		}
		//fmt.Println("-----------------------------------")
	}
	return returnValue, err
}

/*DescribeDatabase returns tables in database*/
func DescribeDatabase(dbType string, connectionString string) ([]ColumnDescription, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	columnDesc := []ColumnDescription{}
	databaseName, err := GetCurrentDatabaseName(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := ColumnDescription{DatabaseName: databaseName, PropertyType: "Table"}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for _, value := range values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))
				newColDesc.TableName = string(value.([]byte))
				newColDesc.PropertyName = string(value.([]byte))

			default:
				//fmt.Println("\t", columns[i], ": ", value)
				newColDesc.TableName = string(value.([]byte))
				newColDesc.PropertyName = string(value.([]byte))
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
	}
	return columnDesc, nil
}

/*DescribeTable returns columns in database*/
func DescribeTable(dbType string, connectionString string, tableName string) ([]ColumnDescription, error) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("DESCRIBE %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	columnDesc := []ColumnDescription{}
	databaseName, err := GetCurrentDatabaseName(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := ColumnDescription{DatabaseName: databaseName, TableName: tableName}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

				switch columns[i] {
				case "Field":
					newColDesc.PropertyName = string(value.([]byte))
				case "Type":
					newColDesc.PropertyType = string(value.([]byte))
				case "Key":
					newColDesc.KeyType = string(value.([]byte))
				case "Null":
					if string(value.([]byte)) == "YES" {
						newColDesc.Nullable = true
					} else {
						newColDesc.Nullable = false
					}
				}
			default:
				//fmt.Println("\t", columns[i], ": ", value)
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
	}
	return columnDesc, nil
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func DescribeTableRelationship(dbType string, connectionString string, tableName string) ([]ColumnRelationship, error) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	databaseName, err := GetCurrentDatabaseName(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT 
		TABLE_NAME,COLUMN_NAME,CONSTRAINT_NAME, REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME
	  FROM
		INFORMATION_SCHEMA.KEY_COLUMN_USAGE
	  WHERE
		REFERENCED_TABLE_SCHEMA = '%s' AND
		REFERENCED_TABLE_NAME = '%s';`, databaseName, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	columnDesc := []ColumnRelationship{}
	newColDesc := ColumnRelationship{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

				switch columns[i] {
				case "TABLE_NAME":
					newColDesc.TableName = string(value.([]byte))
				case "COLUMN_NAME":
					newColDesc.ColumnName = string(value.([]byte))
				case "REFERENCED_TABLE_NAME":
					newColDesc.ReferenceTableName = string(value.([]byte))
				case "REFERENCED_COLUMN_NAME":
					newColDesc.ReferenceColumnName = string(value.([]byte))
				}
			default:
				//fmt.Println("\t", columns[i], ": ", value)
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
	}
	return columnDesc, nil
}
