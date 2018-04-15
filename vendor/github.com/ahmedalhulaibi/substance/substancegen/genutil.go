package substancegen

import (
	"sort"
)

/*StringInSlice returns true if a string is an element within a slice*/
func StringInSlice(searchVal string, list []string) bool {
	for _, val := range list {
		if val == searchVal {
			return true
		}
	}
	return false
}

/*AddJSONTagsToProperties adds json go tags to each property for each object*/
func AddJSONTagsToProperties(gqlObjectTypes map[string]GenObjectType) {

	for _, value := range gqlObjectTypes {
		for _, propVal := range value.Properties {
			propVal.Tags["json"] = append(propVal.Tags["json"], propVal.ScalarName)
		}
	}
}

/*SearchForKeyColumnByKeyType returns a string containing the name of the column of a certain key type*/
func SearchForKeyColumnByKeyType(gqlObjectType GenObjectType, searchKeyType string) string {
	keyColumn := ""
	//Loop through all properties in alphabetic order (key sorted)
	//This prevents different keys being identified across multiple runs using the same input data
	propKeys := make([]string, 0)
	for propKey := range gqlObjectType.Properties {
		propKeys = append(propKeys, propKey)
	}
	sort.Strings(propKeys)
	for _, propKey := range propKeys {
		propVal := gqlObjectType.Properties[propKey]
		if StringInSlice(searchKeyType, propVal.KeyType) {
			keyColumn = propVal.ScalarNameUpper
			break
		}
	}
	return keyColumn
}
