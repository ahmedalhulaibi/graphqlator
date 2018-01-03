package gqlgraphqlator

import (
	"strings"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/graphqlator"
	"github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
)

func init() {
	gqlPlugin := gql{}
	graphqlator.Register("gql", gqlPlugin)
}

type gql struct {
	name string
}

func (g gql) GetGqlObjectTypesFunc(dbType string, connectionString string, tableNames []string) []graphqlator.GqlObjectType {
	//init array of column descriptions for all tables
	tableDesc := []substance.ColumnDescription{}

	//init array of graphql types
	gqlObjectTypes := make(map[string]graphqlator.GqlObjectType)

	//for each table name add a new graphql type and init its properties
	for _, tableName := range tableNames {
		newGqlObj := graphqlator.GqlObjectType{Name: tableName}
		newGqlObj.Properties = make(graphqlator.GqlObjectProperties)
		gqlObjectTypes[tableName] = newGqlObj
		//describe each table
		_results, err := substance.DescribeTable(dbType, connectionString, tableName)
		if err != nil {
			panic(err)
		}
		//append results to tableDesc
		tableDesc = append(tableDesc, _results...)
	}

	//map types
	for _, colDesc := range tableDesc {
		propertyType := ""
		switch {
		case strings.Contains(colDesc.PropertyType, "tinyint(1)") || strings.Contains(colDesc.PropertyType, "bit"):
			propertyType = "Boolean"
			break
		case strings.Contains(colDesc.PropertyType, "varchar"):
			propertyType = "String"
			break
		case strings.Contains(colDesc.PropertyType, "int"):
			propertyType = "Int"
			break
		case strings.Contains(colDesc.PropertyType, "double") || strings.Contains(colDesc.PropertyType, "float") || strings.Contains(colDesc.PropertyType, "decimal") || strings.Contains(colDesc.PropertyType, "numeric"):
			propertyType = "Float"
			break
		}
		newGqlObjProperty := graphqlator.GqlObjectProperty{
			ScalarName: colDesc.PropertyName,
			ScalarType: propertyType,
			Nullable:   colDesc.Nullable,
			KeyType:    colDesc.KeyType}
		gqlObjectTypes[colDesc.TableName].Properties[colDesc.PropertyName] = newGqlObjProperty
	}
	//resolve relationships
	gqlObjectTypes = g.ResolveRelationshipsFunc(dbType,
		connectionString,
		tableNames,
		gqlObjectTypes)

	return nil
}

func (g gql) ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjectTypes map[string]graphqlator.GqlObjectType) map[string]graphqlator.GqlObjectType {
	relationshipDesc := []substance.ColumnRelationship{}
	constraintDesc := []substance.ColumnConstraint{}

	for _, tableName := range tableNames {
		relResults, err := substance.DescribeTableRelationship(dbType, connectionString, tableName)

		if err != nil {
			panic(err)
		}
		relationshipDesc = append(relationshipDesc, relResults...)

		constraintResults, err := substance.DescribeTableConstraints(dbType, connectionString, tableName)

		if err != nil {
			panic(err)
		}
		constraintDesc = append(constraintDesc, constraintResults...)
	}

	for _, colRel := range relationshipDesc {
		//search constraintDesc for columns that are both unique and foreign, or only foreign
		//replace the type info with the appropriate object
		//Example:
		//CREATE TABLE Persons (
		// 	PersonID int PRIMARY KEY,
		// 	LastName varchar(255),
		// 	FirstName varchar(255),
		// 	Address varchar(255),
		// 	City varchar(255)
		// );

		// CREATE TABLE Orders (
		// 	OrderID int UNIQUE NOT NULL,
		// 	OrderNumber int NOT NULL,
		// 	PersonID int DEFAULT NULL,
		// 	PRIMARY KEY (OrderID),
		// 	FOREIGN KEY (PersonID) REFERENCES Persons(PersonID)
		// );
		//
		//The above table would result in an Order object which has a Person object
		//The Person object would have an array of Order objects to reflect the one-to-many relationship

		//Replace column foreign key reference with the Object type (Order has a Person)
		newGqlObjProperty := graphqlator.GqlObjectProperty{
			ScalarName: colRel.ReferenceTableName,
			ScalarType: colRel.ReferenceTableName,
			Nullable:   gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].Nullable,
			KeyType:    gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType}

		gqlObjectTypes[colRel.TableName].Properties[colRel.ReferenceTableName] = newGqlObjProperty

		//Add a new property to table
		//Persons have many orders
		if gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType == "MUL" {
			newGqlObjProperty := graphqlator.GqlObjectProperty{
				ScalarName: colRel.TableName,
				ScalarType: colRel.TableName,
				Nullable:   true,
				IsList:     true}
			gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = newGqlObjProperty
		}
		//remove old property
		delete(gqlObjectTypes[colRel.TableName].Properties, colRel.ColumnName)
		//fmt.Println(gqlObjectTypes)
	}

	return gqlObjectTypes
}

func (g gql) OutputCodeFunc(gqlObjectTypes []graphqlator.GqlObjectType) {

}
