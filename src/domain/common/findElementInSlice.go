package common

func FindElementInSlice(element interface{}, slice []string) bool {
	for _, indexElement := range slice {
		if indexElement == element {
			return true
		}
	}
	return false
}
