package utils

// StringInSlice checks for a string in a slice of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
