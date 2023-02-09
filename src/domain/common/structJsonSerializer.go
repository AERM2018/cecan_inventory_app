package common

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

func StructJsonSerializer(filters interface{}, tagName string, logicalOperator string) string {
	fmt.Println(reflect.TypeOf(filters))
	filtersJson := make(map[string]interface{})
	filterAsMap := structs.Map(filters)

	// Convert struct property names to json tags and remove the ones which are empty
	for _, field := range structs.Fields(filters) {
		if field.Kind().String() == "string" || field.Kind().String() == "int" {
			if !field.IsZero() {
				jsonTag := field.Tag(tagName)
				if strings.Contains(jsonTag, ",") {
					jsonTag = strings.Split(jsonTag, ",")[0]
				}
				filtersJson[jsonTag] = fmt.Sprintf("%v", filterAsMap[field.Name()])
			}
		}
	}
	return jsonToConditionString(filtersJson, logicalOperator)
}

func jsonToConditionString(json map[string]interface{}, logicalOperator string) string {
	conditionString := ""
	filterCounter := 0
	keys := reflect.ValueOf(json).MapKeys()
	// rawKeys := maps.Keys(json)
	if len(keys) > 0 {
		fmt.Println(keys)
		// fmt.Println(filtersJson[keys[0].String()])

		includeLogicalAndOperator := len(keys) > 1
		for _, k := range keys {
			conditionString += fmt.Sprintf("%v::varchar ILIKE %v%v%v", k, "'%", json[k.String()], "%'")
			if includeLogicalAndOperator && filterCounter+1 < len(json) {
				conditionString += fmt.Sprintf(" %v ", logicalOperator)
			}
			filterCounter += 1
		}
	}
	return conditionString
}
