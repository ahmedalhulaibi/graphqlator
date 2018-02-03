package graphqlgo

import (
	"strings"

	"github.com/ahmedalhulaibi/substance"
	"github.com/ahmedalhulaibi/substance/substancegen"
)

func init() {
	gqlPlugin := gql{}
	gqlPlugin.GraphqlDataTypes = make(map[string]string)
	gqlPlugin.GraphqlDataTypes["int"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int8"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int16"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int32"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["int64"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint8"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint16"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint32"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["uint64"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["byte"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["rune"] = "graphql.Int"
	gqlPlugin.GraphqlDataTypes["bool"] = "graphql.Boolean"
	gqlPlugin.GraphqlDataTypes["string"] = "graphql.String"
	gqlPlugin.GraphqlDataTypes["float32"] = "graphql.Float"
	gqlPlugin.GraphqlDataTypes["float64"] = "graphql.Float"
	gqlPlugin.GraphqlDbTypeGormFlag = make(map[string]bool)
	gqlPlugin.GraphqlDbTypeGormFlag["mysql"] = true
	gqlPlugin.GraphqlDbTypeGormFlag["postgres"] = true
	gqlPlugin.GraphqlDbTypeImports = make(map[string]string)
	gqlPlugin.GraphqlDbTypeImports["mysql"] = "\n\t\"github.com/jinzhu/gorm\"\n\t_ \"github.com/jinzhu/gorm/dialects/mysql\""
	gqlPlugin.GraphqlDbTypeImports["postgres"] = "\n\t\"github.com/jinzhu/gorm\"\n\t_ \"github.com/jinzhu/gorm/dialects/postgres\""
	substancegen.Register("graphql-go", gqlPlugin)
}

type gql struct {
	Name                  string
	GraphqlDataTypes      map[string]string
	GraphqlDbTypeGormFlag map[string]bool
	GraphqlDbTypeImports  map[string]string
}

func (g gql) GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]substancegen.GenObjectType {
	//init array of column descriptions for all tables
	tableDesc := []substance.ColumnDescription{}

	//init array of graphql types
	gqlObjectTypes := make(map[string]substancegen.GenObjectType)

	//for each table name add a new graphql type and init its properties
	for _, tableName := range tableNames {
		newGqlObj := substancegen.GenObjectType{Name: tableName}
		newGqlObj.Properties = make(substancegen.GenObjectProperties)
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
		newGqlObjProperty := substancegen.GenObjectProperty{
			ScalarName: colDesc.PropertyName,
			ScalarType: colDesc.PropertyType,
			Nullable:   colDesc.Nullable,
			KeyType:    []string{colDesc.KeyType},
		}
		newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
		if g.GraphqlDbTypeGormFlag[dbType] {
			newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], "column:"+newGqlObjProperty.ScalarName+";")
		}
		gqlObjectTypes[colDesc.TableName].Properties[colDesc.PropertyName] = newGqlObjProperty
	}
	//resolve relationships
	gqlObjectTypes = g.ResolveRelationshipsFunc(dbType,
		connectionString,
		tableNames,
		gqlObjectTypes)

	return gqlObjectTypes
}

func (g gql) ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjectTypes map[string]substancegen.GenObjectType) map[string]substancegen.GenObjectType {
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
			gqlKeyTypes := gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType
			//fmt.Println("GQL Key Type ", constraint.TableName, constraint.ColumnName, gqlKeyTypes)
			for _, gqlKeyType := range gqlKeyTypes {
				switch {
				case gqlKeyType == "" || gqlKeyType == " ":

					newGqlObjProperty := substancegen.GenObjectProperty{
						ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
						ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
						Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
						KeyType:    []string{constraint.ConstraintType},
						Tags:       gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags,
					}
					isPrimary := (stringInSlice("p", newGqlObjProperty.KeyType) || stringInSlice("PRIMARY KEY", newGqlObjProperty.KeyType))
					if isPrimary && g.GraphqlDbTypeGormFlag[dbType] {
						newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], "primary_key"+";")
					}

					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty

				case gqlKeyType == "p" || gqlKeyType == "PRIMARY KEY":
					if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {
						newGqlObjProperty := substancegen.GenObjectProperty{
							ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
							ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
							Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
							KeyType:    append(gqlKeyTypes, "f"),
							Tags:       gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags,
						}

						gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty

					}
				case gqlKeyType == "u" || gqlKeyType == "UNIQUE":
					if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {

						newGqlObjProperty := substancegen.GenObjectProperty{
							ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
							ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
							Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
							KeyType:    append(gqlKeyTypes, "f"),
							Tags:       gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags,
						}

						gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty

					}
				case gqlKeyType == "f" || gqlKeyType == "FOREIGN KEY":
					if constraint.ConstraintType == "p" || constraint.ConstraintType == "PRIMARY KEY" {

						newGqlObjProperty := substancegen.GenObjectProperty{
							ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
							ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
							Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
							KeyType:    append(gqlKeyTypes, "p"),
							Tags:       gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags,
						}

						gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty

					} else if constraint.ConstraintType == "u" || constraint.ConstraintType == "UNIQUE" {

						newGqlObjProperty := substancegen.GenObjectProperty{
							ScalarName: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarName,
							ScalarType: gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].ScalarType,
							Nullable:   gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Nullable,
							KeyType:    append(gqlKeyTypes, "u"),
							Tags:       gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags,
						}

						gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName] = newGqlObjProperty

					}
				}
			}

		}

	}

	for _, colRel := range relationshipDesc {
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
		//The Person object would have an array of Order objects to reflect the one-to-many relationship
		//Add a new property to table
		//Persons have many orders
		isUnique := (stringInSlice("u", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || stringInSlice("UNIQUE", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
		isPrimary := (stringInSlice("p", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || stringInSlice("PRIMARY KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
		isForeign := (stringInSlice("f", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || stringInSlice("FOREIGN KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))

		if isForeign && !isPrimary && !isUnique {
			gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
			gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
			newGqlObjProperty := substancegen.GenObjectProperty{
				ScalarName:   colRel.TableName,
				ScalarType:   strings.TrimSuffix(colRel.TableName, "s"),
				Nullable:     true,
				IsList:       true,
				IsObjectType: true,
			}
			newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
			if g.GraphqlDbTypeGormFlag[dbType] {
				newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)
			}
			gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = newGqlObjProperty
		} else if (isUnique || isPrimary) && isForeign {
			gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
			gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
			newGqlObjProperty := substancegen.GenObjectProperty{
				ScalarName:   strings.TrimSuffix(colRel.TableName, "s"),
				ScalarType:   strings.TrimSuffix(colRel.TableName, "s"),
				Nullable:     true,
				IsList:       false,
				IsObjectType: true,
			}
			newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
			if g.GraphqlDbTypeGormFlag[dbType] {
				newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)
			}
			gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = newGqlObjProperty

		}
	}
	return gqlObjectTypes
}

func stringInSlice(searchVal string, list []string) bool {
	for _, val := range list {
		if val == searchVal {
			return true
		}
	}
	return false
}
