package substancegen

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/ahmedalhulaibi/substance"
	"github.com/apex/log"
	"github.com/jinzhu/inflection"
)

/*GeneratorInterface describes the implementation required to generate code from substance objects*/
type GeneratorInterface interface {
	OutputCodeFunc(dbType string, connectionString string, gqlObjectTypes map[string]GenObjectType) bytes.Buffer
}

/*SubstanceGenPlugins is a map storing a reference to the current plugins
Key: pluginName
Value: reference to an implementation of SubstanceGenInterface*/
var SubstanceGenPlugins = make(map[string]GeneratorInterface)

/*Register registers a GeneratorInterface plugin */
func Register(pluginName string, pluginInterface GeneratorInterface) {
	SubstanceGenPlugins[pluginName] = pluginInterface
}

/*GenObjectTag stores a key value pair of go struct a tag and their value(s)
Example:
Key: gorm
Tabs: {'primary_key','column_name'}*/
type GenObjectTag map[string][]string

/*TODO: Create new type to store KeyType to map [string]string
This will require changes in generators/graphqlgo pkg
This will require changes in generators/gorm.go pkg + gorm_test.go
This will require changes in generators/gostruct_test.go*/

/*GenObjectProperty represents a property of an object (aka a field of a struct) */
type GenObjectProperty struct {
	ScalarName      string `json:"scalarName"`
	ScalarNameUpper string
	ScalarType      string `json:"scalarType"`
	AltScalarType   map[string]string
	IsList          bool         `json:"isList"`
	Nullable        bool         `json:"nullable"`
	KeyType         []string     `json:"keyType"`
	Tags            GenObjectTag `json:"tags"`
	IsObjectType    bool         `json:"isObjectType"`
}

/*GenObjectProperties a type defining a map of GenObjectProperty
Key: PropertyName
Value: GenObjectProperty */
type GenObjectProperties map[string]*GenObjectProperty

/*GenObjectType represents an object (aka a struct) */
type GenObjectType struct {
	Name            string `json:"objectName"`
	SourceTableName string `json:"sourceTableName"`
	LowerName       string
	Properties      GenObjectProperties `json:"properties"`
}

/*Generate is a one stop function to quickly generate code */
func Generate(generatorName string, dbType string, connectionString string, tableNames []string) bytes.Buffer {
	return SubstanceGenPlugins[generatorName].OutputCodeFunc(dbType, connectionString, GetObjectTypesFunc(dbType, connectionString, tableNames))
}

/*GetObjectTypesFunc returns all object definitions as a map given tableNames and connectionString*/
func GetObjectTypesFunc(dbType string, connectionString string, tableNames []string) map[string]GenObjectType {
	//init array of column descriptions for all tables
	tableDesc := []substance.ColumnDescription{}

	//init array of graphql types
	gqlObjectTypes := make(map[string]GenObjectType)

	//for each table name add a new graphql type and init its properties
	for _, tableName := range tableNames {
		a := []rune(inflection.Singular(tableName))
		a[0] = unicode.ToUpper(a[0])
		genObjectTypeNameUpper := string(a)
		a[0] = unicode.ToLower(a[0])
		genObjectTypeNameLower := string(a)
		newGqlObj := GenObjectType{Name: genObjectTypeNameUpper, LowerName: genObjectTypeNameLower, SourceTableName: tableName}
		newGqlObj.Properties = make(GenObjectProperties)
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
		newGqlObjProperty := GenObjectProperty{
			ScalarName:      colDesc.PropertyName,
			ScalarNameUpper: colDescPropNameUpper,
			ScalarType:      colDesc.PropertyType,
			AltScalarType:   make(map[string]string),
			Nullable:        colDesc.Nullable,
			KeyType:         []string{""},
		}
		newGqlObjProperty.Tags = make(GenObjectTag)
		newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], "column:"+newGqlObjProperty.ScalarName+";")

		gqlObjectTypes[colDesc.TableName].Properties[colDesc.PropertyName] = &newGqlObjProperty
	}
	//resolve relationships
	gqlObjectTypes = ResolveRelationshipsFunc(dbType,
		connectionString,
		tableNames,
		gqlObjectTypes)

	return gqlObjectTypes
}

/*ResolveRelationshipsFunc calls out to multiple functions to orchestrate the resolution of constraints and relationships*/
func ResolveRelationshipsFunc(dbType string, connectionString string, tableNames []string, gqlObjectTypes map[string]GenObjectType) map[string]GenObjectType {
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

		ResolveConstraintsFunc(dbType, constraintDesc, gqlObjectTypes)
	}

	ResolveForeignRefsFunc(dbType, relationshipDesc, gqlObjectTypes)
	return gqlObjectTypes
}

/*ResolveConstraintsFunc maps assigns key constraints received from ColumnConstraint and appends them to an array for the corresponding property
This function also appends gorm tags for primary_key constraints*/
func ResolveConstraintsFunc(dbType string, constraintDesc []substance.ColumnConstraint, gqlObjectTypes map[string]GenObjectType) {
	for _, constraint := range constraintDesc {
		gqlKeyTypes := &gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType
		gqlTags := &gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].Tags
		//fmt.Println("GQL Key Type ", constraint.TableName, constraint.ColumnName, gqlKeyTypes)
		for _, gqlKeyType := range *gqlKeyTypes {
			switch {
			case gqlKeyType == "" || gqlKeyType == " ":
				gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType = []string{constraint.ConstraintType}
				isPrimary := (StringInSlice("p", gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType) ||
					StringInSlice("PRIMARY KEY", gqlObjectTypes[constraint.TableName].Properties[constraint.ColumnName].KeyType))
				if isPrimary {
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

/*ResolveForeignRefsFunc resolves the foreign key relationships between objects and inserts properties that have associations
This is able to handle many-to-one and one-to-many
Untested is the many-to-many assocation using a cross-reference table. For example:
	Customer table
		ID
		NAME
	Account table
		AccountNumber
		Product
	CustomerToAccount table
		CustomerID <- Foreign Key
		AccountNumber <- Foreign Key
		RelationshipType
	The current expected result is that this would resolve to:
	struct Customer type{
		...
		CustomerToAccount []CustomerToAccountType
	}
	struct Account type{
		...
		CustomerToAccount []CustomerToAccountType
	}
*/
func ResolveForeignRefsFunc(dbType string, relationshipDesc []substance.ColumnRelationship, gqlObjectTypes map[string]GenObjectType) {
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
			isUnique := (StringInSlice("u", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || StringInSlice("UNIQUE", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
			isPrimary := (StringInSlice("p", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || StringInSlice("PRIMARY KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))
			isForeign := (StringInSlice("f", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType) || StringInSlice("FOREIGN KEY", gqlObjectTypes[colRel.TableName].Properties[colRel.ColumnName].KeyType))

			if isForeign && !isPrimary && !isUnique {
				gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
				gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
				newGqlObjProperty := GenObjectProperty{
					ScalarName:      inflection.Plural(gqlObjectTypes[colRel.TableName].Name),
					ScalarNameUpper: inflection.Plural(gqlObjectTypes[colRel.TableName].Name),
					ScalarType:      gqlObjectTypes[colRel.TableName].Name,
					Nullable:        true,
					IsList:          true,
					IsObjectType:    true,
					AltScalarType:   make(map[string]string),
				}
				newGqlObjProperty.Tags = make(GenObjectTag)
				newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)

				gqlObjectTypes[colRel.ReferenceTableName].Properties[colRel.TableName] = &newGqlObjProperty
			} else if (isUnique || isPrimary) && isForeign {
				gormTagForeign := "ForeignKey:" + colRel.ColumnName + ";"
				gormTagAssociationForeign := "AssociationForeignKey:" + colRel.ReferenceColumnName + ";"
				newGqlObjProperty := GenObjectProperty{
					ScalarName:      gqlObjectTypes[colRel.TableName].Name,
					ScalarNameUpper: gqlObjectTypes[colRel.TableName].Name,
					ScalarType:      gqlObjectTypes[colRel.TableName].Name,
					Nullable:        true,
					IsList:          false,
					IsObjectType:    true,
					AltScalarType:   make(map[string]string),
				}
				newGqlObjProperty.Tags = make(GenObjectTag)
				newGqlObjProperty.Tags["gorm"] = append(newGqlObjProperty.Tags["gorm"], gormTagForeign, gormTagAssociationForeign)
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
