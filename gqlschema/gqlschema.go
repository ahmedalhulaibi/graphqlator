package gqlgraphqlator

import (
	"fmt"
	"strings"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
	"github.com/ahmedalhulaibi/go-graphqlator-cli/substancegen"
)

func init() {
	gqlPlugin := gql{}
	substancegen.Register("graphql-go", gqlPlugin)
}

type gql struct {
	name string
}

func (g gql) GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]substancegen.GqlObjectType {
	//init array of column descriptions for all tables
	tableDesc := []substance.ColumnDescription{}

	//init array of graphql types
	gqlObjectTypes := make(map[string]substancegen.GqlObjectType)

	//for each table name add a new graphql type and init its properties
	for _, tableName := range tableNames {
		newGqlObj := substancegen.GqlObjectType{Name: tableName}
		newGqlObj.Properties = make(substancegen.GqlObjectProperties)
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
	//TODO: remove this mapping once substance supports go data types in the PropertyType field
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
		newGqlObjProperty := substancegen.GqlObjectProperty{
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

	return gqlObjectTypes
}

func (g gql) ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjectTypes map[string]substancegen.GqlObjectType) map[string]substancegen.GqlObjectType {
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

		for _, constraint := range constraintDesc {
			gqlKeyType := gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType
			fmt.Println("GQL Key Type ", constraint.TableName, constraint.ColumnName, gqlKeyType)
			switch {
			case gqlKeyType == "":
				newGqlObjProperty := substancegen.GqlObjectProperty{
					ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
					ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
					Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
					KeyType:    constraint.ConstraintType}
				gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty
			case gqlKeyType == "p" || gqlKeyType == "PRIMARY KEY":
				if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {
					newGqlObjProperty := substancegen.GqlObjectProperty{
						ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
						ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
						Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
						KeyType:    "UFO"}
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty
				}
			case gqlKeyType == "u" || gqlKeyType == "UNIQUE":
				if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {
					newGqlObjProperty := substancegen.GqlObjectProperty{
						ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
						ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
						Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
						KeyType:    "UFO"}
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty
				}
			case gqlKeyType == "f" || gqlKeyType == "FOREIGN KEY":
				if constraint.ConstraintType == "p" || constraint.ConstraintType == "PRIMARY KEY" || constraint.ConstraintType == "u" || constraint.ConstraintType == "UNIQUE" {
					newGqlObjProperty := substancegen.GqlObjectProperty{
						ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
						ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
						Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
						KeyType:    "UFO"}
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty
				}
			}
		}

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
		newGqlObjProperty := substancegen.GqlObjectProperty{
			ScalarName: colRel.ReferenceTableName,
			ScalarType: colRel.ReferenceTableName,
			Nullable:   gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].Nullable,
			KeyType:    gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType}

		gqlObjectTypes[colRel.TableName].Properties[colRel.ReferenceTableName] = newGqlObjProperty

		//Add a new property to table
		//Persons have many orders
		if gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType == "FOREIGN KEY" ||
			gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType == "f" {
			newGqlObjProperty := substancegen.GqlObjectProperty{
				ScalarName: colRel.TableName,
				ScalarType: colRel.TableName,
				Nullable:   true,
				IsList:     true}
			gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = newGqlObjProperty
		} else if gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType == "UFO" {
			newGqlObjProperty := substancegen.GqlObjectProperty{
				ScalarName: colRel.TableName,
				ScalarType: colRel.TableName,
				Nullable:   true,
				IsList:     false}
			gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = newGqlObjProperty
		}
		//remove old property
		delete(gqlObjectTypes[colRel.TableName].Properties, colRel.ColumnName)
		//fmt.Println(gqlObjectTypes)
	}

	return gqlObjectTypes
}

func (g gql) OutputCodeFunc(gqlObjectTypes map[string]substancegen.GqlObjectType) {

	//print schema
	for _, value := range gqlObjectTypes {
		fmt.Printf("type %s {\n", value.Name)
		for _, propVal := range value.Properties {
			nullSymbol := "!"
			if propVal.Nullable {
				nullSymbol = ""
			}
			if propVal.IsList {
				fmt.Printf("\t %s: [%s]%s\n", propVal.ScalarName, propVal.ScalarType, nullSymbol)
			} else {
				fmt.Printf("\t %s: %s%s\n", propVal.ScalarName, propVal.ScalarType, nullSymbol)
			}
		}
		fmt.Println("}")
	}
}
