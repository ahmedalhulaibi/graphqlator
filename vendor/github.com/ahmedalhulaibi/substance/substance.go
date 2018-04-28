package substance

import "database/sql"

/*SubstanceInterface defines the functions that must be implemented*/
type SubstanceInterface interface {
	GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error)
	DescribeDatabaseFunc(dbType string, connectionString string) ([]ColumnDescription, error)
	DescribeTableFunc(dbType string, connectionString string, tableName string) ([]ColumnDescription, error)
	DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]ColumnRelationship, error)
	DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]ColumnConstraint, error)
	GetGoDataType(sqlType string) (string, error)
}

/*substance plugin map*/
var substancePlugins = make(map[string]SubstanceInterface)

/*Register registers a sbustance plugin which implements the Substance interface*/
func Register(pluginName string, pluginInterface SubstanceInterface) {
	//fmt.Println(substancePlugins)
	substancePlugins[pluginName] = pluginInterface
}

/*ColumnDescription Structure to store properties of each column in a table */
type ColumnDescription struct {
	DatabaseName string
	TableName    string
	PropertyName string
	PropertyType string
	KeyType      string
	DefaultValue string
	Nullable     bool
}

/*ColumnRelationship Structure to store relationships between tables*/
type ColumnRelationship struct {
	TableName           string
	ColumnName          string
	ReferenceTableName  string
	ReferenceColumnName string
}

/*ColumnConstraint Struct to store column constraint types*/
type ColumnConstraint struct {
	TableName      string
	ColumnName     string
	ConstraintType string
}

/*QueryResult Struct to store results from ExecuteQuery*/
type QueryResult struct {
	Rows     *sql.Rows
	Columns  []string
	Values   []interface{}
	ScanArgs []interface{}
	Err      error
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func GetCurrentDatabaseName(dbType string, connectionString string) (string, error) {
	return substancePlugins[dbType].GetCurrentDatabaseNameFunc(dbType, connectionString)
}

/*DescribeDatabase returns tables in database*/
func DescribeDatabase(dbType string, connectionString string) ([]ColumnDescription, error) {
	return substancePlugins[dbType].DescribeDatabaseFunc(dbType, connectionString)
}

/*DescribeTable returns columns of a table*/
func DescribeTable(dbType string, connectionString string, tableName string) ([]ColumnDescription, error) {
	return substancePlugins[dbType].DescribeTableFunc(dbType, connectionString, tableName)
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func DescribeTableRelationship(dbType string, connectionString string, tableName string) ([]ColumnRelationship, error) {
	return substancePlugins[dbType].DescribeTableRelationshipFunc(dbType, connectionString, tableName)
}

/*DescribeTableConstraints returns all column constraints in a database table*/
func DescribeTableConstraints(dbType string, connectionString string, tableName string) ([]ColumnConstraint, error) {
	return substancePlugins[dbType].DescribeTableConstraintsFunc(dbType, connectionString, tableName)
}

/*ExecuteQuery executes a sql query with one or no tableName, specific to mysqlsubstnace and pgsqlsubstance*/
func ExecuteQuery(dbType string, connectionString string, tableName string, query string) QueryResult {
	db, err := sql.Open(dbType, connectionString)
	defer db.Close()
	if err != nil {
		return QueryResult{Err: err}
	}
	var rows *sql.Rows
	if tableName == "" {
		rows, err = db.Query(query)
	} else {
		rows, err = db.Query(query, tableName)
	}
	if err != nil {
		return QueryResult{Err: err, Rows: rows}
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return QueryResult{Err: err, Rows: rows, Columns: columns}
	}
	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	return QueryResult{Err: err, Rows: rows, Columns: columns, ScanArgs: scanArgs, Values: values}
}
