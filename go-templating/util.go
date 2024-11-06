package main

import (
	"reflect"
	"strings"
)

// Converts a struct object to a map. It also recursively converts nested structs to maps.
// Inputs
// - object (any): The struct to convert to a map.
// Outputs
// - map[string]any: The map representation of the object.
func ConvertStructToMap(object any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(object)
	typ := reflect.TypeOf(object)
	// Loop through all fields in the struct and convert them to map
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		// Get corresponding tag value for variable consistency.
		tag := structField.Tag.Get("yaml")
		if tag == "" {
			tag = structField.Tag.Get("name")
			if tag == "" {
				// If no predermined tag, convert first letter of variable name to lowercase
				// to match yaml format.
				tag = strings.ToLower(structField.Name[:1]) + structField.Name[1:]
			}
		}
		if field.Kind() == reflect.Struct {
			// Recursively convert nested struct to map
			result[tag] = ConvertStructToMap(field.Interface())
		} else {
			result[tag] = field.Interface()
		}
	}
	return result
}
