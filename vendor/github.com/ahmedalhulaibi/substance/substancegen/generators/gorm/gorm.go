package gorm

import (
	"bytes"
	"fmt"

	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/ahmedalhulaibi/substance/substancegen/generators/genutil"
	"github.com/jinzhu/inflection"
)

/*GenGormObjectTableNameOverrideFunc generates a function to override the GORM default table name
See examples over override at http://doc.gorm.io/models.html#conventions*/
func GenGormObjectTableNameOverrideFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(gqlObjectType.Name)
	buff.WriteString(fmt.Sprintf("\nfunc (%s) TableName() string {\n\treturn \"%s\"\n}\n", gqlObjectTypeNameSingular, gqlObjectType.Name))
}

/*GenObjectGormCreateFunc generates functions for basic CRUD Create using gorm and writes it to a buffer*/
func GenObjectGormCreateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(gqlObjectType.Name)

	buff.WriteString(fmt.Sprintf("\n\nfunc Create%s (db *gorm.DB, new%s %s) {\n\tdb.Create(&new%s)\n}",
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular))
}

/*GenObjectGormReadFunc generates functions for basic CRUD Read/Get using gorm and writes it to a buffer*/
func GenObjectGormReadFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(gqlObjectType.Name)

	buff.WriteString(fmt.Sprintf("\n\nfunc Get%s (db *gorm.DB, query%s %s, result%s *%s) {\n\tdb.Where(&query%s).First(result%s)\n}",
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular))
}

/*GenObjectGormUpdateFunc generates functions for basic CRUD Update using gorm and writes it to a buffer*/
func GenObjectGormUpdateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(gqlObjectType.Name)
	var primaryKeyColumn string
	for index, propVal := range gqlObjectType.Properties {
		if genutil.StringInSlice("p", propVal.KeyType) || genutil.StringInSlice("PRIMARY KEY", propVal.KeyType) {
			primaryKeyColumn = index
			break
		}
	}

	buff.WriteString(fmt.Sprintf("\n\nfunc Update%s (db *gorm.DB, old%s %s, new%s %s, result%s *%s) {\n\tvar oldResult%s %s\n\tdb.Where(&old%s).First(&oldResult%s)\n\tif oldResult%s.%s == new%s.%s {\n\t\toldResult%s = new%s\n\t\tdb.Save(oldResult%s)\n\t}\n\tGet%s(db, new%s, result%s)\n}",
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		primaryKeyColumn,
		gqlObjectTypeNameSingular,
		primaryKeyColumn,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular))
}

/*GenObjectGormDeleteFunc generates functions for basic CRUD Delete using gorm and writes it to a buffer*/
func GenObjectGormDeleteFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(gqlObjectType.Name)

	buff.WriteString(fmt.Sprintf("\n\nfunc Delete%s (db *gorm.DB, old%s %s) {\n\tdb.Delete(&old%s)\n}",
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular,
		gqlObjectTypeNameSingular))
}

/*GenObjectGormCrud generates functions for basic CRUD operations using gorm and writes it to a buffer*/
func GenObjectGormCrud(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	GenObjectGormCreateFunc(gqlObjectType, buff)

	GenObjectGormReadFunc(gqlObjectType, buff)

	GenObjectGormUpdateFunc(gqlObjectType, buff)

	GenObjectGormDeleteFunc(gqlObjectType, buff)
}
