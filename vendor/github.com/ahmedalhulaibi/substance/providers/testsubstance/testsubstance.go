package testsubstance

import (
	"github.com/ahmedalhulaibi/substance"
)

func init() {
	testPlugin := testsql{}
	substance.Register("test", &testPlugin)
}

type testsql struct {
	name string
}

/*GetCurrentDatabaseName returns currrent database schema name as string*/
func (t testsql) GetCurrentDatabaseNameFunc(dbType string, connectionString string) (string, error) {
	returnValue := "testDatabase"
	var err error
	return returnValue, err
}

/*DescribeDatabase returns tables in database*/
func (t testsql) DescribeDatabaseFunc(dbType string, connectionString string) ([]substance.ColumnDescription, error) {
	columnDesc := []substance.ColumnDescription{}
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		PropertyType: "Table",
		PropertyName: "TableNumberOne",
		TableName:    "TableNumberOne",
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		PropertyType: "Table",
		PropertyName: "TableNumberTwo",
		TableName:    "TableNumberTwo",
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		PropertyType: "Table",
		PropertyName: "TableNumberThree",
		TableName:    "TableNumberThree",
	})
	return columnDesc, nil
}

/*DescribeTable returns columns in database*/
func (t testsql) DescribeTableFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnDescription, error) {
	columnDesc := []substance.ColumnDescription{}
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberOne",
		PropertyType: "int32",
		PropertyName: "UniqueIdOne",
		KeyType:      "p",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberOne",
		PropertyType: "string",
		PropertyName: "Name",
		KeyType:      "",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberOne",
		PropertyType: "float64",
		PropertyName: "Salary",
		KeyType:      "",
		Nullable:     true,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberTwo",
		PropertyType: "UniqueIdTwo",
		PropertyName: "int32",
		KeyType:      "",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberTwo",
		PropertyType: "ForeignIdOne",
		PropertyName: "int32",
		KeyType:      "f",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberThree",
		PropertyType: "UniqueIdThree",
		PropertyName: "int32",
		KeyType:      "",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberThree",
		PropertyType: "ForeignIdOne",
		PropertyName: "int32",
		KeyType:      "f",
		Nullable:     false,
	})
	columnDesc = append(columnDesc, substance.ColumnDescription{
		DatabaseName: "testDatabase",
		TableName:    "TableNumberThree",
		PropertyType: "ForeignIdTwo",
		PropertyName: "int32",
		KeyType:      "f",
		Nullable:     true,
	})
	return columnDesc, nil
}

/*DescribeTableRelationship returns all foreign column references in database table*/
func (t testsql) DescribeTableRelationshipFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnRelationship, error) {
	columnRel := []substance.ColumnRelationship{}
	columnRel = append(columnRel, substance.ColumnRelationship{
		TableName:           "TableNumberTwo",
		ColumnName:          "ForeignIdOne",
		ReferenceTableName:  "TableNumberOne",
		ReferenceColumnName: "UniqueIdOne",
	})
	columnRel = append(columnRel, substance.ColumnRelationship{
		TableName:           "TableNumberThree",
		ColumnName:          "ForeignIdOne",
		ReferenceTableName:  "TableNumberOne",
		ReferenceColumnName: "UniqueIdOne",
	})
	columnRel = append(columnRel, substance.ColumnRelationship{
		TableName:           "TableNumberThree",
		ColumnName:          "ForeignIdTwo",
		ReferenceTableName:  "TableNumberTwo",
		ReferenceColumnName: "UniqueIdTwo",
	})
	return columnRel, nil
}

func (t testsql) DescribeTableConstraintsFunc(dbType string, connectionString string, tableName string) ([]substance.ColumnConstraint, error) {
	columnDesc := []substance.ColumnConstraint{}
	return columnDesc, nil
}

func (t testsql) GetGoDataType(sqlType string) (string, error) {
	return "", nil
}
