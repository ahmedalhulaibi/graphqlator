package pgsqlsubstance

import (
	"database/sql"
	"fmt"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
)

func init() {
	pgsqlPlugin := pgsql{}
	substance.Register("postgres", &pgsqlPlugin)
}

type pgsql struct {
	name string
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func (m pgsql) GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error) {
	returnValue := "postgres"
	var err error
	return returnValue, err
}

/*DescribeDatabase returns tables in database*/
func (m pgsql) DescribeDatabaseFunc(dbType string, connectionString string) ([]substance.ColumnDescription, error) {
	postgresString := "postgres://"
	connString := postgresString + connectionString
	db, err := sql.Open(dbType, connString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM pg_catalog.pg_tables where schemaname not in ('pg_catalog','information_schema');")
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
	var subsInterface = pgsql{}
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
			case []byte:
				switch columns[i] {
				case "tablename":
					newColDesc.TableName = string(value.([]byte))
				case "schemaname":
					newColDesc.PropertyName = string(value.([]byte))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
	}
	return columnDesc, nil
}

/*DescribeTable returns columns in database*/
func (m pgsql) DescribeTableFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnDescription, error) {
	postgresString := "postgres://"
	connString := postgresString + connectionString
	db, err := sql.Open(dbType, connString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`select
		att.attrelid as "classId",
		class.relname as "Table",
		att.attname as "Field",
		dsc.description as "description",
		typ.typname as "Type",
		att.attnum as "num",
		att.attnotnull as "isNotNull",
		att.atthasdef as "hasDefault"
	  from
		pg_catalog.pg_attribute as att
		left join pg_catalog.pg_description as dsc on dsc.objoid = att.attrelid and dsc.objsubid = att.attnum
		left join pg_type as typ on typ.oid = att.atttypid
		left join pg_catalog.pg_class as class on class.oid = att.attrelid
	  where
		att.attrelid in (
			select rel.oid as "id"
			from pg_catalog.pg_class as rel
			left join pg_catalog.pg_description as dsc on dsc.objoid = rel.oid and dsc.objsubid = 0
			where 
			class.relname = $1 and
			rel.relpersistence in ('p') and
			rel.relkind in ('r', 'v', 'm', 'c', 'f')
		) and
		att.attnum > 0 and
		not att.attisdropped
	  order by
		att.attrelid, att.attnum;`, tableName)
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
	var subsInterface = pgsql{}
	databaseName, err := subsInterface.GetCurrentDatabaseNameFunc(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	newColDesc := substance.ColumnDescription{DatabaseName: databaseName, TableName: tableName}

	//get all column constraints to determine key type
	//columnConstraints, err := subsInterface.DescribeTableConstraintsFunc(dbType, connectionString, tableName)
	if err != nil {
		return nil, err
	}

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
				case "Field":
					newColDesc.PropertyName = string(value.([]byte))
				case "Type":
					newColDesc.PropertyType = string(value.([]byte))
				case "isNotNull":
					if string(value.([]byte)) == "f" {
						newColDesc.Nullable = true
					} else {
						newColDesc.Nullable = false
					}
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
		//fmt.Println("-----------------------------------")
	}
	return columnDesc, nil
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func (m pgsql) DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnRelationship, error) {
	postgresString := "postgres://"
	connString := postgresString + connectionString
	db, err := sql.Open(dbType, connString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`select
							tc.table_name as "table_name",
							kcu.column_name as "column",
							class.relname as "ref_table",
							con.confkey as "ref_columnNum"
							  from
								pg_catalog.pg_constraint as con
								left join information_schema.table_constraints as tc on tc.constraint_name = con.conname
								left join information_schema.key_column_usage as kcu on kcu.constraint_name = con.conname
								left join pg_catalog.pg_class as class on class.oid = con.confrelid
							  where
								tc.table_name = '%s' and
								con.contype = 'f'
							  	;`, tableName)
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
				fmt.Println("\t", columns[i], ": ", value)
				switch columns[i] {
				case "table_name":
					newColDesc.TableName = string(value.(string))
				case "column":
					newColDesc.ColumnName = string(value.(string))
				}
			case []byte:
				fmt.Println("\t", columns[i], ": ", string(value.([]byte)))

				switch columns[i] {
				case "ref_table":
					newColDesc.ReferenceTableName = string(value.([]byte))
				case "ref_columnNum":
					newColDesc.ReferenceColumnName = string(value.([]byte))
				}
			}
		}
		columnDesc = append(columnDesc, newColDesc)
	}
	return columnDesc, nil
}

func (m pgsql) DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnConstraint, error) {
	postgresString := "postgres://"
	connString := postgresString + connectionString
	db, err := sql.Open(dbType, connString)
	defer db.Close()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`select distinct on (con.conrelid, con.conkey, con.confrelid, con.confkey)
	tc.table_name,
	kcu.column_name as "column",
	contype
  from
	pg_catalog.pg_constraint as con
	left join information_schema.table_constraints as tc on tc.constraint_name = con.conname
	left join information_schema.key_column_usage as kcu on kcu.constraint_name = con.conname
	left join pg_catalog.pg_class as class on class.oid = con.confrelid
  where
		tc.table_name = $1
  ;`, tableName)
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
