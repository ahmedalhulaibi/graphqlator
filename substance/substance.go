package substance

/*SubstanceInterface defines the functions that must be implemented*/
type SubstanceInterface interface {
	GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error)
	DescribeDatabaseFunc(dbType string, connectionString string) ([]ColumnDescription, error)
	DescribeTableFunc(dbType string, connectionString string, tableName string) ([]ColumnDescription, error)
	DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]ColumnRelationship, error)
	DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]ColumnConstraint, error)
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

func GetCurrentDatabaseName(dbType string, connectionString string) (string, error) {
	return substancePlugins[dbType].GetCurrentDatabaseNameFunc(dbType, connectionString)
}

func DescribeDatabase(dbType string, connectionString string) ([]ColumnDescription, error) {
	return substancePlugins[dbType].DescribeDatabaseFunc(dbType, connectionString)
}
func DescribeTable(dbType string, connectionString string, tableName string) ([]ColumnDescription, error) {
	return substancePlugins[dbType].DescribeTableFunc(dbType, connectionString, tableName)
}
func DescribeTableRelationship(dbType string, connectionString string, tableName string) ([]ColumnRelationship, error) {
	return substancePlugins[dbType].DescribeTableRelationshipFunc(dbType, connectionString, tableName)
}

func DescribeTableConstraints(dbType string, connectionString string, tableName string) ([]ColumnConstraint, error) {
	return substancePlugins[dbType].DescribeTableConstraintsFunc(dbType, connectionString, tableName)
}
