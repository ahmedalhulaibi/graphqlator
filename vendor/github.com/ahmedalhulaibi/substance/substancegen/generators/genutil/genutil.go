package genutil

/*StringInSlice returns true if a string is an element within a slice*/
func StringInSlice(searchVal string, list []string) bool {
	for _, val := range list {
		if val == searchVal {
			return true
		}
	}
	return false
}
