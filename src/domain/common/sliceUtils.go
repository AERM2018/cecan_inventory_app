package common

import (
	"reflect"
)

func FindElementInSlice(element interface{}, slice []string) bool {
	for _, indexElement := range slice {
		if indexElement == element {
			return true
		}
	}
	return false
}

func FilterSlice(arr interface{}, cond func(interface{}) bool) interface{} {
	contentType := reflect.TypeOf(arr)
	contentValue := reflect.ValueOf(arr)

	newContent := reflect.MakeSlice(contentType, 0, 0)
	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); cond(content.Interface()) {
			newContent = reflect.Append(newContent, content)
		}
	}
	return newContent.Interface()
}
