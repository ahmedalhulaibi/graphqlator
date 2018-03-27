package substancegen

/*StringInSlice returns true if a string is an element within a slice*/
func StringInSlice(searchVal string, list []string) bool {
	for _, val := range list {
		if val == searchVal {
			return true
		}
	}
	return false
}

func AddJSONTagsToProperties(gqlObjectTypes map[string]GenObjectType) {

	for _, value := range gqlObjectTypes {
		for _, propVal := range value.Properties {
			propVal.Tags["json"] = append(propVal.Tags["json"], propVal.ScalarName)
		}
	}
}
