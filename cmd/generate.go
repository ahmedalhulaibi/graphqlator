package cmd

import (
	"fmt"
	"strings"

	"github.com/ahmedalhulaibi/go-graphqlator-cli/substance"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generate)
}

type gqlObjectProperty struct {
	scalarName string
	scalarType string
	isList     bool
	nullable   bool
	keyType    string
}

type gqlObjectProperties map[string]gqlObjectProperty

type gqlObjectType struct {
	name       string
	properties gqlObjectProperties
}

var generate = &cobra.Command{
	Use:   "generate [database type] [connection string] [table names...]",
	Short: "Generate GraphQL type schema from database table(s).",
	Long:  `Generate GraphQL type schema from database table(s).`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "mariadb":
			args[0] = "mysql"
			break
		}
		generateGqlSchema(args[0], args[1], args[2:len(args)])
	},
}

func generateGqlSchema(dbType string, connectionString string, tableNames []string) {
	tableDesc := []substance.ColumnDescription{}
	gqlObjectTypes := make(map[string]gqlObjectType)
	for _, tableName := range tableNames {
		newGqlObj := gqlObjectType{name: tableName}
		newGqlObj.properties = make(gqlObjectProperties)
		gqlObjectTypes[tableName] = newGqlObj
		_results, err := substance.DescribeTable(dbType, connectionString, tableName)
		if err != nil {
			panic(err)
		}
		tableDesc = append(tableDesc, _results...)
	}
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
		newGqlObjProperty := gqlObjectProperty{scalarName: colDesc.PropertyName, scalarType: propertyType, nullable: colDesc.Nullable, keyType: colDesc.KeyType}
		gqlObjectTypes[colDesc.TableName].properties[colDesc.PropertyName] = newGqlObjProperty

		//fmt.Println(gqlObjectTypes[colDesc.TableName])
	}
	relationshipDesc := []substance.ColumnRelationship{}

	for _, tableName := range tableNames {
		_results, err := substance.DescribeTableRelationship(dbType, connectionString, tableName)
		if err != nil {
			panic(err)
		}
		relationshipDesc = append(relationshipDesc, _results...)
	}
	for _, colRel := range relationshipDesc {
		//creating a type for col
		/*Example:
		Given the below sql layout
		Orders Table
		---------------
		OrderID		int
		OrderNumber int
		PersonID	int <------ this is a foreign key reference
		===============
		Persons Table
		---------------
		ID 			int
		Name		string
		===============

		The relationship above in graphql schema would be a has-a relationship: An order has a person associated with it
		Additionally, depending on the key type for PersonID in orders, Persons could have a one-to-one or one-to-many relationship to Orders

		This code creates a Persons relationship for Orders, and removes the PersonID reference
		*/
		newGqlObjProperty := gqlObjectProperty{
			scalarName: colRel.ReferenceTableName,
			scalarType: colRel.ReferenceTableName,
			nullable:   gqlObjectTypes[colRel.TableName].properties[colRel.ColumnName].nullable,
			keyType:    gqlObjectTypes[colRel.TableName].properties[colRel.ColumnName].keyType}

		gqlObjectTypes[colRel.TableName].properties[colRel.ReferenceTableName] = newGqlObjProperty

		if gqlObjectTypes[colRel.TableName].properties[colRel.ColumnName].keyType == "MUL" {
			newGqlObjProperty := gqlObjectProperty{
				scalarName: colRel.TableName,
				scalarType: colRel.TableName,
				nullable:   true,
				isList:     true}
			gqlObjectTypes[colRel.ReferenceTableName].properties[colRel.TableName] = newGqlObjProperty
		}
		//remove old property
		delete(gqlObjectTypes[colRel.TableName].properties, colRel.ColumnName)
		//fmt.Println(gqlObjectTypes)
	}

	//print schema
	for _, value := range gqlObjectTypes {
		fmt.Printf("type %s {\n", value.name)
		for _, propVal := range value.properties {
			nullSymbol := "!"
			if propVal.nullable {
				nullSymbol = ""
			}
			if propVal.isList {
				fmt.Printf("\t %s: [%s]%s\n", propVal.scalarName, propVal.scalarType, nullSymbol)
			} else {
				fmt.Printf("\t %s: %s%s\n", propVal.scalarName, propVal.scalarType, nullSymbol)
			}
		}
		fmt.Println("}")
	}
}
