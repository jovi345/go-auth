package helper

import "reflect"

func GetJSONFieldName(fieldName string, inputStruct interface{}) string {
	t := reflect.TypeOf(inputStruct)
	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	return jsonTag
}
