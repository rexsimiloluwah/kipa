package utils

// Utility function to check if a string slice contains an element
func StrSliceContains(s []string, el string) bool {
	var result bool = false
	for i := 0; i < len(s); i++ {
		if el == s[i] {
			result = true
			break
		}
	}

	return result
}
