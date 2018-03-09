package graphqlgo

import (
	"fmt"
	"unicode"

	"github.com/ahmedalhulaibi/substance"
	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/genutil"
	"github.com/apex/log"
	"github.com/jinzhu/inflection"
)

func init() {
	gqlPlugin := Gql{}
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

type Gql struct {
	Name                  string
	GraphqlDataTypes      map[string]string
	GraphqlDbTypeGormFlag map[string]bool
	GraphqlDbTypeImports  map[string]string
}

/*GetObjectTypesFunc*/
func (g Gql) GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]substancegen.GenObjectType {
	//init array of column descriptions for all tables
	tableDesc := []substance.ColumnDescription{}

	//init array of graphql types
	gqlObjectTypes := make(map[string]substancegen.GenObjectType)

	//for each table name add a new graphql type and init its properties
	for _, tableName := range tableNames {
		a := []rune(inflection.Singular(tableName))
		a[0] = unicode.ToUpper(a[0])
		genObjectTypeNameUpper := string(a)
		a[0] = unicode.ToLower(a[0])
		genObjectTypeNameLower := string(a)
		newGqlObj := substancegen.GenObjectType{Name: genObjectTypeNameUpper, LowerName: genObjectTypeNameLower, SourceTableName: tableName}
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
		a := []rune(inflection.Singular(colDesc.PropertyName))
		a[0] = unicode.ToUpper(a[0])
		colDescPropNameUpper := string(a)
		newGqlObjProperty := substancegen.GenObjectProperty{
			ScalarName:      colDesc.PropertyName,
			ScalarNameUpper: colDescPropNameUpper,
			ScalarType:      colDesc.PropertyType,
			Nullable:        colDesc.Nullable,
			KeyType:         []string{""},
		}
		newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
		if _, ok := g.GraphqlDbTypeGormFlag[dbType]; ok {
			newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], "column:"+newGqlObjProperty.ScalarName+";")
		}
		gqlObjectTypes[colDesc.TableName].Properties[colDesc.PropertyName] = &newGqlObjProperty
	}
	//resolve relationships
	gqlObjectTypes = g.ResolveRelationshipsFunc(dbType,
		connectionString,
		tableNames,
		gqlObjectTypes)

	return gqlObjectTypes
}

func (g Gql) ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjectTypes map[string]substancegen.GenObjectType) map[string]substancegen.GenObjectType {
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

		g.ResolveConstraintsFunc(dbType, constraintDesc, gqlObjectTypes)
	}

	g.ResolveForeignRefsFunc(dbType, relationshipDesc, gqlObjectTypes)
	return gqlObjectTypes
}

func (g Gql) ResolveConstraintsFunc(dbType string, constraintDesc []substance.ColumnConstraint, gqlObjectTypes map[string]substancegen.GenObjectType) {
	for _, constraint := range constraintDesc {
		gqlKeyTypes := &gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType
		gqlTags := &gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags
		//fmt.Println("GQL Key Type ", constraint.TableName, constraint.ColumnName, gqlKeyTypes)
		for _, gqlKeyType := range *gqlKeyTypes {
			switch {
			case gqlKeyType == "" || gqlKeyType == " ":
				gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = []string{constraint.ConstraintType}
				isPrimary := (genutil.StringInSlice("p", gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType) ||
					genutil.StringInSlice("PRIMARY KEY", gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType))
				if isPrimary && g.GraphqlDbTypeGormFlag[dbType] {
					(*gqlTags)["gorm"] = append((*gqlTags)["gorm"], "primary_key"+";")
				}
			case gqlKeyType == "p" || gqlKeyType == "PRIMARY KEY":
				if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = append(gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType, "f")
				}
			case gqlKeyType == "u" || gqlKeyType == "UNIQUE":
				if constraint.ConstraintType == "f" || constraint.ConstraintType == "FOREIGN KEY" {
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = append(gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType, "f")
				}
			case gqlKeyType == "f" || gqlKeyType == "FOREIGN KEY":
				if constraint.ConstraintType == "p" || constraint.ConstraintType == "PRIMARY KEY" {
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = append(gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType, "p")

				} else if constraint.ConstraintType == "u" || constraint.ConstraintType == "UNIQUE" {
					gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = append(gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType, "u")

				}
			}
		}

	}
}

func (g Gql) ResolveForeignRefsFunc(dbType string, relationshipDesc []substance.ColumnRelationship, gqlObjectTypes map[string]substancegen.GenObjectType) {
	for _, colRel := range relationshipDesc {
		_, colRelTableOk := gqlObjectTypes[colRel.TableName]
		_, colRelRefTableOk := gqlObjectTypes[colRel.ReferenceTableName]
		if colRelTableOk && colRelRefTableOk {

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
			isUnique := (genutil.StringInSlice("u", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || genutil.StringInSlice("UNIQUE", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
			isPrimary := (genutil.StringInSlice("p", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || genutil.StringInSlice("PRIMARY KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
			isForeign := (genutil.StringInSlice("f", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || genutil.StringInSlice("FOREIGN KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))

			if isForeign && !isPrimary && !isUnique {
				gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
				gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
				newGqlObjProperty := substancegen.GenObjectProperty{
					ScalarName:      inflection.Plural(gqlObjectTypes[colRel.TableName].Name),
					ScalarNameUpper: inflection.Plural(gqlObjectTypes[colRel.TableName].Name),
					ScalarType:      gqlObjectTypes[colRel.TableName].Name,
					Nullable:        true,
					IsList:          true,
					IsObjectType:    true,
				}
				newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
				if g.GraphqlDbTypeGormFlag[dbType] {
					newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)
				}
				gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = &newGqlObjProperty
			} else if (isUnique || isPrimary) && isForeign {
				gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
				gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
				newGqlObjProperty := substancegen.GenObjectProperty{
					ScalarName:      gqlObjectTypes[colRel.TableName].Name,
					ScalarNameUpper: gqlObjectTypes[colRel.TableName].Name,
					ScalarType:      gqlObjectTypes[colRel.TableName].Name,
					Nullable:        true,
					IsList:          false,
					IsObjectType:    true,
				}
				newGqlObjProperty.Tags = make(substancegen.GenObjectTag)
				if g.GraphqlDbTypeGormFlag[dbType] {
					newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)
				}
				gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = &newGqlObjProperty

			}
		}

		if !colRelTableOk {
			log.Errorf(fmt.Sprintf("%s Table definition not found in ResolveForeignRefsFunc gqlObjectTypes", colRel.TableName))
		}

		if !colRelRefTableOk {
			log.Errorf(fmt.Sprintf("%s Table definition not found in ResolveForeignRefsFunc gqlObjectTypes", colRel.ReferenceTableName))
		}
	}
}
