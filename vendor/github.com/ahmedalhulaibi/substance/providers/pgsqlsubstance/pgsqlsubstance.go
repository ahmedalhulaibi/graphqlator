package pgsqlsubstance

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ahmedalhulaibi/substance"
	/*blank import to load postgres driver*/
	_ "github.com/lib/pq"
)

func init() {
	pgsqlPlugin := pgsql{}
	substance.Register("postgres", &pgsqlPlugin)
}

type pgsql struct {
	name string
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func (p pgsql) GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error) {
	returnValue := "placeholder"
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return "", err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, "", GetCurrentDatabaseNameQuery)
	if queryResult.Err != nil {
		return "", queryResult.Err
	}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return "", err
		}

		// Print data
		for i, value := range queryResult.Values {
			switch value.(type) {
			case []byte:
				switch queryResult.Columns[i] {
				case "current_database":
					returnValue = string(value.([]byte))
				}
			}
		}

	}

	return returnValue, err
}

/*DescribeDatabase returns tables in database*/
func (p pgsql) DescribeDatabaseFunc(dbType string, connectionString string) ([]substance.ColumnDescription, error) {
	//opening connection
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, "", DescribeDatabaseQuery)

	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	//setup array of column descriptions
	columnDesc := []substance.ColumnDescription{}

	//get database name
	databaseName, err := p.GetCurrentDatabaseNameFunc(dbType, connectionString)
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
			case []byte:
				switch queryResult.Columns[i] {
				case "tablename":
					newColDesc.TableName = string(value.([]byte))
					newColDesc.PropertyName = string(value.([]byte))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
	}
	return columnDesc, nil
}

/*DescribeTable returns columns in database*/
func (p pgsql) DescribeTableFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnDescription, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, tableName, DescribeTableQuery)

	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnDesc := []substance.ColumnDescription{}

	databaseName, err := p.GetCurrentDatabaseNameFunc(dbType, connectionString)
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
			case bool:
				switch queryResult.Columns[i] {
				case "isNotNull":
					newColDesc.Nullable = !value.(bool)
				}
			case []byte:
				switch queryResult.Columns[i] {
				case "Field":
					newColDesc.PropertyName = string(value.([]byte))
				case "Type":
					newColDesc.PropertyType, _ = p.GetGoDataType(string(value.([]byte)))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)

	}
	return columnDesc, nil
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func (p pgsql) DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnRelationship, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, tableName, DescribeTableRelationshipQuery)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}

	columnTableDesc, err := substance.DescribeTable(dbType, connectionString, tableName)
	if err != nil {
		return nil, err
	}
	columnDesc := []substance.ColumnRelationship{}
	newColDesc := substance.ColumnRelationship{}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {

			switch value.(type) {
			case string:

				switch queryResult.Columns[i] {
				case "table_name":
					newColDesc.TableName = string(value.(string))
				case "column":
					newColDesc.ColumnName = string(value.(string))
				}
			case []byte:

				switch queryResult.Columns[i] {
				case "ref_table":
					newColDesc.ReferenceTableName = string(value.([]byte))
					columnTableDesc, err = substance.DescribeTable(dbType, connectionString, newColDesc.ReferenceTableName)
					if err != nil {
						return nil, err
					}
				case "ref_columnNum":
					//this gets returned as {1} a reference to the column number in the table
					//this has to be replaced with the column name

					refColumnNumStr := strings.Replace(strings.Replace(string(value.([]byte)), "{", "", -1), "}", "", -1)

					refColumnNum, err := strconv.Atoi(refColumnNumStr)
					if err != nil {
						return nil, err
					}

					newColDesc.ReferenceColumnName = columnTableDesc[refColumnNum-1].PropertyName

				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
	}
	return columnDesc, nil
}

func (p pgsql) DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnConstraint, error) {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	queryResult := substance.ExecuteQuery(dbType, connectionString, tableName, DescribeTableConstraintsQuery)
	if queryResult.Err != nil {
		return nil, queryResult.Err
	}
	columnDesc := []substance.ColumnConstraint{}
	newColDesc := substance.ColumnConstraint{}

	for queryResult.Rows.Next() {
		err = queryResult.Rows.Scan(queryResult.ScanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range queryResult.Values {
			switch value.(type) {
			case string:

				switch queryResult.Columns[i] {
				case "table_name":
					newColDesc.TableName = string(value.(string))
				case "column":
					newColDesc.ColumnName = string(value.(string))
				case "contype":
					newColDesc.ConstraintType = string(value.(string))
				}
			default:

			}
		}
		columnDesc = append(columnDesc, newColDesc)
	}
	return columnDesc, nil
}

func (p pgsql) GetGoDataType(sqlType string) (string, error) {
	if regexDataTypePatterns == nil {
		regexDataTypePatterns["bit.*"] = "int64"
		regexDataTypePatterns["bool.*|tinyint\\(1\\)"] = "bool"
		regexDataTypePatterns["tinyint.*"] = "int8"
		regexDataTypePatterns["unsigned\\stinyint.*"] = "uint8"
		regexDataTypePatterns["smallint.*"] = "int16"
		regexDataTypePatterns["unsigned\\ssmallint.*"] = "uint16"
		regexDataTypePatterns["(mediumint.*|int.*)"] = "int32"
		regexDataTypePatterns["unsigned\\s(mediumint.*|int.*)"] = "uint32"
		regexDataTypePatterns["bigint.*"] = "int64"
		regexDataTypePatterns["unsigned\\sbigint.*"] = "uint64"
		regexDataTypePatterns["(unsigned\\s){0,1}(double.*|float.*|dec.*)"] = "float64"
		regexDataTypePatterns["varchar.*|date.*|time.*|year.*|char.*|.*text.*|enum.*|set.*|.*blob.*|.*binary.*"] = "string"
	}

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
