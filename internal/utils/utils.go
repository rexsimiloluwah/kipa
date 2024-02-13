package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Utility function to check if a slice contains a value
func Contains[T comparable](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// checks if an interface is a JSON data
func IsJSON(data interface{}) bool {
	if _, ok := data.(string); !ok {
		return false
	}
	if !strings.ContainsAny(data.(string), "{") || !strings.ContainsAny(data.(string), "}") {
		return false
	}
	var js json.RawMessage
	return json.Unmarshal([]byte(data.(string)), &js) == nil
}

// returns the underlying type of an interface
func TypeOf(data interface{}) string {
	if IsJSON(data) {
		return "json"
	}
	switch t := data.(type) {
	case string:
		return "string"
	case int:
		return "number"
	case float64:
		return "number"
	case bool:
		return "bool"
	default:
		r := reflect.TypeOf(t)
		if r == reflect.TypeOf(make([]interface{}, 0)) {
			return "array"
		} else if r == reflect.TypeOf(make(map[string]interface{}, 0)) {
			return "object"
		}
		return r.String()
	}
}
