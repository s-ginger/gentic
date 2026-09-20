package tools

// Schema describes the input expected by a Tool.
//
// Type follows JSON Schema type names such as:
//
//   - object
//   - string
//   - number
//   - integer
//   - boolean
//   - array
type Schema struct {
	// Type specifies the type of the input.
	Type string

	// Properties describes the fields accepted by the tool.
	Properties map[string]Property

	// Required contains the names of fields that must be provided.
	Required []string
}

// Property describes a single field in a Tool input schema.
type Property struct {
	// Type specifies the type of the property.
	Type string

	// Description explains what the property is used for.
	Description string
}