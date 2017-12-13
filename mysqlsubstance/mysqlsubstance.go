package mysqlsubstance

import (
	"database/sql"
	"fmt"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
)

func init() {
	mysqlPlugin := mysql{}
	substance.Register("mysql", &mysqlPlugin)
}

type mysql struct {
	name string
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func (m mysql) GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return "nil", err
	}
	rows, err := db.Query("SELECT DATABASE()")
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
				err := fmt.Errorf("No database found make sure connection string includes database. e.g. user:pass@localhost:port/database")
				return "nil", error(err)
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
func (m mysql) DescribeDatabaseFunc(dbType string, connectionString string) ([]substance.ColumnDescription, error) {
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

	columnDesc := []substance.ColumnDescription{}
	var subsInterface = mysql{}
	databaseName, err := subsInterface.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, PropertyType: "Table"}

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

				err := fmt.Errorf("Null column value found at column: '%s' index: '%d'", columns[i], i)
				return nil, error(err)
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
func (m mysql) DescribeTableFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnDescription, error) {

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

	columnDesc := []substance.ColumnDescription{}
	var subsInterface = mysql{}
	databaseName, err := subsInterface.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, TableName: tableName}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				//IGNORE NIL VALUE
				//fmt.Println("\t", columns[i], ": NULL")
				//err := fmt.Errorf("Null column value found at column: '%s' index: '%d'", columns[i], i)
				//return nil, error(err)
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
func (m mysql) DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnRelationship, error) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	subsInterface := mysql{}
	databaseName, err := subsInterface.GetCurrentDatabaseNameFunc(dbType, connectionString)
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

	columnDesc := []substance.ColumnRelationship{}
	newColDesc := substance.ColumnRelationship{}

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
				err := fmt.Errorf("Null column value found at column: '%s' index: '%d'", columns[i], i)
				return nil, error(err)
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

/*DescribeTableRelationship returns all foreign column references in database table*/
func (m mysql) DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnConstraint, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`SELECT DISTINCT kcu.column_name as 'Column', tc.constraint_type as 'Constraint'
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE as kcu
		JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS as tc on tc.constraint_name = kcu.constraint_name
		WHERE kcu.table_name = '%s'
		order by kcu.column_name, tc.constraint_type;`, tableName)
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

	columnDesc := []substance.ColumnConstraint{}
	newColDesc := substance.ColumnConstraint{}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			newColDesc.TableName = tableName
			switch value.(type) {
			case []byte:
				switch columns[i] {
				case "Column":
					newColDesc.ColumnName = string(value.([]byte))
				case "Constraint":
					newColDesc.ConstraintType = string(value.([]byte))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
	}
	return columnDesc, nil
}
