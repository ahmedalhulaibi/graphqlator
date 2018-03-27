package gorm

import (
	"bytes"
	"log"
	"text/template"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

/*GenGormObjectTableNameOverrideFunc generates a function to override the GORM default table name
See examples over override at http://doc.gorm.io/models.html#conventions*/
func GenGormObjectTableNameOverrideFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormObjTblNameOverrideTemplate := "\nfunc ({{.Name}}) TableName() string {\n\treturn \"{{.SourceTableName}}\"\n}\n"
	tmpl := template.New("gormObjTblNameOverride")
	tmpl, err := tmpl.Parse(gormObjTblNameOverrideTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormCreateFunc generates functions for basic CRUD Create using gorm and writes it to a buffer*/
func GenObjectGormCreateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormCreateFuncTemplate := "\n\nfunc Create{{.Name}} (db *gorm.DB, new{{.Name}} {{.Name}}) {\n\tdb.Create(&new{{.Name}})\n}"
	tmpl := template.New("gormCreateFunc")
	tmpl, err := tmpl.Parse(gormCreateFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormReadFunc generates functions for basic CRUD Read/Get using gorm and writes it to a buffer*/
func GenObjectGormReadFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormReadFuncTemplate := "\n\nfunc Get{{.Name}} (db *gorm.DB, query{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) {\n\tdb.Where(&query{{.Name}}).First(result{{.Name}})\n}"
	tmpl := template.New("gormReadFunc")
	tmpl, err := tmpl.Parse(gormReadFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormUpdateFunc generates functions for basic CRUD Update using gorm and writes it to a buffer*/
func GenObjectGormUpdateFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	primaryKeyColumn := ""
	for _, propVal := range gqlObjectType.Properties {
		if substancegen.StringInSlice("p", propVal.KeyType) || substancegen.StringInSlice("PRIMARY KEY", propVal.KeyType) {
			primaryKeyColumn = propVal.ScalarNameUpper
			break
		}
		if substancegen.StringInSlice("u", propVal.KeyType) || substancegen.StringInSlice("UNIQUE", propVal.KeyType) {
			primaryKeyColumn = propVal.ScalarNameUpper
			break
		}
	}
	var templateData = struct {
		Name string
		Key  string
	}{
		gqlObjectType.Name,
		primaryKeyColumn,
	}

	gormUpdateFuncTemplate := "\n\nfunc Update{{.Name}} (db *gorm.DB, old{{.Name}} {{.Name}}, new{{.Name}} {{.Name}}, result{{.Name}} *{{.Name}}) {\n\tvar oldResult{{.Name}} {{.Name}}\n\tdb.Where(&old{{.Name}}).First(&oldResult{{.Name}})\n\tif oldResult{{.Name}}.{{.Key}} == new{{.Name}}.{{.Key}} {\n\t\toldResult{{.Name}} = new{{.Name}}\n\t\tdb.Save(oldResult{{.Name}})\n\t}\n\tGet{{.Name}}(db, new{{.Name}}, result{{.Name}})\n}"
	tmpl := template.New("gormUpdateFunc")
	tmpl, err := tmpl.Parse(gormUpdateFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, templateData)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormDeleteFunc generates functions for basic CRUD Delete using gorm and writes it to a buffer*/
func GenObjectGormDeleteFunc(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gormDeleteFuncTemplate := "\n\nfunc Delete{{.Name}} (db *gorm.DB, old{{.Name}} {{.Name}}) {\n\tdb.Delete(&old{{.Name}})\n}"
	tmpl := template.New("gormReadFunc")
	tmpl, err := tmpl.Parse(gormDeleteFuncTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, gqlObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}

/*GenObjectGormCrud generates functions for basic CRUD operations using gorm and writes it to a buffer*/
func GenObjectGormCrud(gqlObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	GenObjectGormCreateFunc(gqlObjectType, buff)

	GenObjectGormReadFunc(gqlObjectType, buff)

	GenObjectGormUpdateFunc(gqlObjectType, buff)

	GenObjectGormDeleteFunc(gqlObjectType, buff)
}
