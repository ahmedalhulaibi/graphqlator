package gostruct

import (
	"bytes"
	"log"
	"text/template"

	"github.com/ahmedalhulaibi/substance/substancegen"
)

/*GenObjectTypeToStructFunc takes a GenObjectType and writes it to a buffer as a go struct*/
func GenObjectTypeToStructFunc(genObjectTypes map[string]substancegen.GenObjectType, buff *bytes.Buffer) {
	gostructTemplate := "{{range $key, $value := . }}\ntype {{.Name}} struct { {{range .Properties}}\n\t{{.ScalarNameUpper}}\t{{if .IsList}}[]{{end}}{{.ScalarType}}\t`{{range $index, $element := .Tags}}{{$index}}:\"{{range $element}}{{.}}{{end}}\" {{end}}`{{end}}\n}\n{{end}}"

	tmpl := template.New("gostruct")
	tmpl, err := tmpl.Parse(gostructTemplate)
	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}
	err1 := tmpl.Execute(buff, genObjectTypes)
	if err1 != nil {
		log.Fatal("Execute: ", err1)
		return
	}
}
