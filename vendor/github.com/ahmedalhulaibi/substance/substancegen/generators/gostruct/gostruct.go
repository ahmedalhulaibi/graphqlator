package gostruct

import (
	"bytes"
	"log"
	"text/template"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

/*GenObjectTypeToStructFunc takes a GenObjectType and writes it to a buffer as a go struct*/
func GenObjectTypeToStructFunc(genObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gostructTemplate := "\ntype {{.Name}} struct { {{range .Properties}}\n\t{{.ScalarNameUpper}}\t{{if .IsList}}[]{{end}}{{.ScalarType}}\t`{{range $index, $element := .Tags}}{{$index}}:\"{{range $element}}{{.}}{{end}}\" {{end}}`{{end}}\n}\n"

	tmpl := template.New("gostruct")
	tmpl, err := tmpl.Parse(gostructTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, genObjectType)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}
