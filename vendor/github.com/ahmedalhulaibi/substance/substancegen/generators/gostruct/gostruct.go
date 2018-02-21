package gostruct

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/ahmedalhulaibi/substance/substancegen"
	"github.com/jinzhu/inflection"
)

/*GenObjectTypeToStructFunc takes a GenObjectType and writes it to a buffer as a go struct*/
func GenObjectTypeToStructFunc(genObjectType substancegen.GenObjectType, buff *bytes.Buffer) {
	gqlObjectTypeNameSingular := inflection.Singular(genObjectType.Name)
	buff.WriteString(fmt.Sprintf("\ntype %s struct {\n", gqlObjectTypeNameSingular))
	for _, property := range genObjectType.Properties {
		GenObjectPropertyToStringFunc(*property, buff)
	}
	buff.WriteString("}\n")
}

/*GenObjectPropertyToStringFunc takes a GenObjectProperty and writes it to a buffer as a go struct instance variable*/
func GenObjectPropertyToStringFunc(genObjectType substancegen.GenObjectProperty, buff *bytes.Buffer) {

	a := []rune(genObjectType.ScalarName)
	a[0] = unicode.ToUpper(a[0])
	gqlObjectPropertyNameUpper := string(a)
	if genObjectType.IsList {
		buff.WriteString(fmt.Sprintf("\t%s\t[]%s\t", gqlObjectPropertyNameUpper, genObjectType.ScalarType))
	} else {
		buff.WriteString(fmt.Sprintf("\t%s\t%s\t", gqlObjectPropertyNameUpper, genObjectType.ScalarType))
	}
	GenObjectTagToStringFunc(genObjectType.Tags, buff)
	buff.WriteString("\n")
}

/*GenObjectTagToStringFunc takes a map GenObjectTap and writes it to a buffer as go struct instance variable tags*/
func GenObjectTagToStringFunc(genObjectTags substancegen.GenObjectTag, buff *bytes.Buffer) {
	buff.WriteString("`")
	for key, tags := range genObjectTags {
		buff.WriteString(fmt.Sprintf("%s:\"", key))
		for _, tag := range tags {
			buff.WriteString(fmt.Sprintf("%s", tag))
		}
		buff.WriteString("\" ")
	}
	buff.WriteString("`")
}
