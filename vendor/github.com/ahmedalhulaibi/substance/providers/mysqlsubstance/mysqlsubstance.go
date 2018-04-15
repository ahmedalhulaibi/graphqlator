package mysqlsubstance

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/ahmedalhulaibi/substance"
	/*blank import to load mysql driver*/
	_ "github.com/go-sql-driver/mysql"
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

	queryResult := substance.ExecuteQuery(dbType, connectionString, "", GetCurrentDatabaseNameQuery)
	if queryResult.Err != nil {
		return "", queryResult.Err
	}

	var returnValue string
	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return "nil", err
		}

		// Print data
		for _, value := range queryResult.Values {
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

	queryResult := substance.ExecuteQuery(dbType, connectionString, "", DescribeDatabaseQuery)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnDesc := []substance.ColumnDescription{}

	databaseName, err := m.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, PropertyType: "Table"}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")

				err := fmt.Errorf("Null column value found at column: '%s' index: '%d'", queryResult.Columns[i], i)
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

/*DescribeTable returns columns of a table*/
func (m mysql) DescribeTableFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnDescription, error) {

	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(DescribeTableQuery, tableName)
	queryResult := substance.ExecuteQuery(dbType, connectionString, "", query)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnDesc := []substance.ColumnDescription{}

	databaseName, err := m.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, TableName: tableName}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {
			switch value.(type) {
			case []byte:

				switch queryResult.Columns[i] {
				case "Field":
					newColDesc.PropertyName = string(value.([]byte))
				case "Type":
					newColDesc.PropertyType, err = m.GetGoDataType(string(value.([]byte)))
					if err != nil {
						fmt.Printf("Warning: %s", err.Error())
					}
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

	databaseName, err := m.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(DescribeTableRelationshipQuery, databaseName)

	queryResult := substance.ExecuteQuery(dbType, connectionString, tableName, query)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnRel := []substance.ColumnRelationship{}
	newColRel := substance.ColumnRelationship{}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {
			switch value.(type) {
			case nil:
				//fmt.Println("\t", columns[i], ": NULL")
				err := fmt.Errorf("Null column value found at column: '%s' index: '%d'", queryResult.Columns[i], i)
				return nil, error(err)
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

				switch queryResult.Columns[i] {
				case "TABLE_NAME":
					newColRel.TableName = string(value.([]byte))
				case "COLUMN_NAME":
					newColRel.ColumnName = string(value.([]byte))
				case "REFERENCED_TABLE_NAME":
					newColRel.ReferenceTableName = string(value.([]byte))
				case "REFERENCED_COLUMN_NAME":
					newColRel.ReferenceColumnName = string(value.([]byte))
				}
			default:
				//fmt.Println("\t", columns[i], ": ", value)
			}
		}
		columnRel = append(columnRel, newColRel)
		//fmt.Println("-----------------------------------")
	}
	return columnRel, nil
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func (m mysql) DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnConstraint, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, tableName, DescribeTableConstraintsQuery)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnCon := []substance.ColumnConstraint{}
	newColCon := substance.ColumnConstraint{}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {
			newColCon.TableName = tableName
			switch value.(type) {
			case []byte:
				switch queryResult.Columns[i] {
				case "Column":
					newColCon.ColumnName = string(value.([]byte))
				case "Constraint":
					newColCon.ConstraintType = string(value.([]byte))
				}
			}
		}
		columnCon = append(columnCon, newColCon)
		//fmt.Println("-----------------------------------")
	}
	return columnCon, nil
}

/*GetGoDataType returns the go data type for the equivalent mysql data type*/
func (m mysql) GetGoDataType(sqlType string) (string, error) {
	for pattern, value := range regexDataTypePatterns {
		match, err := regexp.MatchString(pattern, sqlType)
		if match && err == nil {
			result := value
			return result, nil
		}
	}
	err := fmt.Errorf("No match found for data type %s", sqlType)
	fmt.Println(err)
	return sqlType, err
}
