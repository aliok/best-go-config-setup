package util

import (
	"github.com/invopop/jsonschema"
	"log"
	"strconv"
	"strings"
)

// VisitSchema visits all the schemas in the schema tree and calls the visitor function for the schemas that have the given propType.
func VisitSchema(schema *jsonschema.Schema, propType string, visitor func(*jsonschema.Schema)) {
	if schema.Type == propType {
		visitor(schema)
	}

	for pair := schema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		VisitSchema(pair.Value, propType, visitor)
	}
	for _, def := range schema.Definitions {
		VisitSchema(def, propType, visitor)
	}
}

// FixArrayDefaultValues fixes the default values of array fields in a JSON schema.
// go-defaultz expects the default values of array fields to be in the form of a space-separated string as in "a b c" or "1.2 2.5 -21.3".
// This function converts the default values of array fields to the appropriate type, such as []string{"a", "b", "c"} or []int{1, 2, 3}.
func FixArrayDefaultValues(schema *jsonschema.Schema) {
	if schema.Default == nil {
		return
	}

	var ok bool

	// convert schema.Default to an array
	// since the type is array and we don't actually provide an array but a space separated string (such as "a b c"),
	// the first item will be the array as string like []string{"a b c"}
	var asArray []interface{}
	if asArray, ok = schema.Default.([]interface{}); !ok {
		return
	}

	if len(asArray) == 0 {
		return
	}

	var defaultStr string
	if defaultStr, ok = asArray[0].(string); !ok {
		return
	}

	// now we have the default value as a string
	// defaultStr="a b c"
	// OR defaultStr="1.2 2.5 -21.3"
	// we now need to create a real array, based on the item type
	// for example, if the item type is string, we need to convert it to []string{"a", "b", "c"}
	// if the item type is number, we need to convert it to []int{1.2, 2.5, -21.3}
	//
	// Spec says
	//	   String values MUST be one of the six primitive types ("null", "boolean", "object", "array", "number", or "string"),
	//	   or "integer" which matches any number with a zero fractional part.
	// https://json-schema.org/draft/2020-12/json-schema-validation#name-type
	switch schema.Items.Type {
	case "string":
		schema.Default = strings.Fields(defaultStr)
	case "integer":
		parts := strings.Fields(defaultStr)
		arr := make([]int, 0)
		for _, part := range parts {
			i, err := strconv.Atoi(part)
			if err != nil {
				log.Fatalf("Failed to convert default value to integer: %v", err)
			}
			arr = append(arr, i)
		}
		schema.Default = arr
	case "number":
		parts := strings.Fields(defaultStr)
		arr := make([]float64, 0)
		for _, part := range parts {
			f, err := strconv.ParseFloat(part, 64)
			if err != nil {
				log.Fatalf("Failed to convert default value to float64: %v", err)
			}
			arr = append(arr, f)
		}
		schema.Default = arr
	case "boolean":
		parts := strings.Fields(defaultStr)
		arr := make([]bool, 0)
		for _, part := range parts {
			b, err := strconv.ParseBool(part)
			if err != nil {
				log.Fatalf("Failed to convert default value to bool: %v", err)
			}
			arr = append(arr, b)
		}
		schema.Default = arr
	default:
		log.Fatalf("Unsupported array item type: %v", schema.Items.Type)
	}
}
