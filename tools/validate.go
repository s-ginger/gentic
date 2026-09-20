package tools

import (
	"fmt"
	"math"
)

// Validate validates input against the schema.
//
// The current implementation supports object schemas with
// string, number, integer, boolean, array, and object properties.
func (s Schema) Validate(input map[string]any) error {
	if s.Type != "object" {
		return fmt.Errorf("root schema must be object")
	}

	if input == nil {
		return fmt.Errorf("input must not be nil")
	}

	// Check required fields.
	for _, name := range s.Required {
		if _, ok := input[name]; !ok {
			return fmt.Errorf("missing required field %q", name)
		}
	}

	// Check provided fields.
	for name, value := range input {
		property, ok := s.Properties[name]
		if !ok {
			return fmt.Errorf("unknown field %q", name)
		}

		if err := property.Validate(name, value); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates a single property value.
func (p Property) Validate(name string, value any) error {
	if value == nil {
		return fmt.Errorf("%q must not be null", name)
	}

	switch p.Type {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%q must be a string", name)
		}

	case "number":
		if !isNumber(value) {
			return fmt.Errorf("%q must be a number", name)
		}

	case "integer":
		if !isInteger(value) {
			return fmt.Errorf("%q must be an integer", name)
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%q must be a boolean", name)
		}

	case "array":
		if !isArray(value) {
			return fmt.Errorf("%q must be an array", name)
		}

	case "object":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("%q must be an object", name)
		}

	default:
		return fmt.Errorf(
			"unsupported property type %q for field %q",
			p.Type,
			name,
		)
	}

	return nil
}

func isNumber(value any) bool {
	switch value.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}

func isInteger(value any) bool {
	switch v := value.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true

	// JSON numbers are decoded into float64 by encoding/json.
	case float32:
		return !math.IsNaN(float64(v)) &&
			!math.IsInf(float64(v), 0) &&
			v == float32(int64(v))

	case float64:
		return !math.IsNaN(v) &&
			!math.IsInf(v, 0) &&
			v == math.Trunc(v)

	default:
		return false
	}
}

func isArray(value any) bool {
	switch value.(type) {
	case []any:
		return true

	case []string:
		return true

	case []int:
		return true

	case []int64:
		return true

	case []float64:
		return true

	case []bool:
		return true

	default:
		return false
	}
}