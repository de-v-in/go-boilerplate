package jsonschema

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func GenJSONSchema(data interface{}) string {
    t := reflect.TypeOf(data)
	kind := t.Kind()
    if kind != reflect.Struct {
        // return "{}"
		switch kind {
		case reflect.String:
			return "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return "integer"
		case reflect.Float32, reflect.Float64:
			return "number"
		case reflect.Bool:
			return "boolean"
		default:
			return "{}"
		}
    }

	schema := make(map[string]interface{})
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldName := field.Tag.Get("json")

        if fieldName == "" {
            continue
        }
        bindingTag := field.Tag.Get("binding")
        fieldSchema := map[string]interface{}{
            "type": getFieldType(field.Type),
        }
        parseBindingTag(bindingTag, fieldSchema)
        fmt.Printf("Field: %s, Schema: %v\n", fieldName, fieldSchema)
        schema[fieldName] = fieldSchema
    }

    schemaBytes, _ := json.Marshal(schema)
    return string(schemaBytes)
}

func parseBindingTag(bindingTag string, fieldSchema map[string]interface{}) {
	tags := strings.Split(bindingTag, ",")
	for _, tag := range tags {
		if tag == "required" {
			fieldSchema["required"] = true
		} else if strings.HasPrefix(tag, "max=") {
			maxValue := strings.TrimPrefix(tag, "max=")
			fieldSchema["max"] = maxValue
		} else if strings.HasPrefix(tag, "min=") {
            minValue := strings.TrimPrefix(tag, "min=")
            fieldSchema["min"] = minValue
        }
	}
}


func getFieldType(fieldType reflect.Type) interface{} {
    kind := fieldType.Kind()
    switch kind {
    case reflect.String:
        return "string"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return "integer"
    case reflect.Float32, reflect.Float64:
        return "number"
    case reflect.Bool:
        return "boolean"
    case reflect.Struct:
        // Nested struct, generate schema recursively
        return GenJSONSchema(reflect.New(fieldType).Elem().Interface())
    default:
        return "unknown"
    }
}