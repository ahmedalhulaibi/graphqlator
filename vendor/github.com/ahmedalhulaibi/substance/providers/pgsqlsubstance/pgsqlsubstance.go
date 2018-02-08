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

	rows, err := db.Query(GetCurrentDatabaseNameQuery)
	if err != nil {
		return "", err
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return "", err
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
			return "", err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case []byte:
				switch columns[i] {
				case "current_database":
					returnValue = string(value.([]byte))
				}
			}
		}
		//fmt.Println("-----------------------------------")
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

	//setup query
	rows, err := db.Query(DescribeDatabaseQuery)
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

	//setup array of column descriptions
	columnDesc := []substance.ColumnDescription{}

	//get database name
	databaseName, err := substance.GetCurrentDatabaseName(dbType, connectionString)
	if err != nil {
		return nil, err
	}

	//newColDesc to be added to columnDesc array
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, PropertyType: "Table"}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case []byte:
				switch columns[i] {
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

	rows, err := db.Query(DescribeTableQuery, tableName)
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

	databaseName, err := p.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, TableName: tableName}

	//get all column constraints to determine key type
	//columnConstraints, err := subsInterface.DescribeTableConstraintsFunc(dbType, connectionString, tableName)
	// if err != nil {
	// 	return nil, err
	// }

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case bool:
				switch columns[i] {
				case "isNotNull":
					newColDesc.Nullable = !value.(bool)
				}
			case []byte:
				switch columns[i] {
				case "Field":
					newColDesc.PropertyName = string(value.([]byte))
				case "Type":
					newColDesc.PropertyType, _ = p.GetGoDataType(string(value.([]byte)))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
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

	if err != nil {
		return nil, err
	}
	rows, err := db.Query(DescribeTableRelationshipQuery, tableName)
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

	columnTableDesc, err := substance.DescribeTable(dbType, connectionString, tableName)
	if err != nil {
		return nil, err
	}
	columnDesc := []substance.ColumnRelationship{}
	newColDesc := substance.ColumnRelationship{}
	//newColDesc.TableName = tableName
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Print data
		for i, value := range values {
			//fmt.Printf("DescribeTableRelationshipFunc Value %T ", value)
			switch value.(type) {
			case string:
				//fmt.Println("\t", columns[i], ": ", value)
				switch columns[i] {
				case "table_name":
					newColDesc.TableName = string(value.(string))
				case "column":
					newColDesc.ColumnName = string(value.(string))
				}
			case []byte:
				//fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

				switch columns[i] {
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

	if err != nil {
		return nil, err
	}
	rows, err := db.Query(DescribeTableConstraintsQuery, tableName)
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
			switch value.(type) {
			case string:
				//fmt.Println("\t", columns[i], ": ", string(value.(string)))

				switch columns[i] {
				case "table_name":
					newColDesc.TableName = string(value.(string))
				case "column":
					newColDesc.ColumnName = string(value.(string))
				case "contype":
					newColDesc.ConstraintType = string(value.(string))
				}
			default:
				//fmt.Println("\t", columns[i], ": ", value)
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
