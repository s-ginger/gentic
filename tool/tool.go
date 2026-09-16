package tool

import "context"

// Tool represents an executable tool that can be exposed to an LLM or
// invoked directly by an application.
//
// A Tool provides a name, a description, an input schema, and an execution
// method.
//
// Example:
//
//	type Calculator struct{}
//
//	func (Calculator) Name() string {
//		return "calculator"
//	}
//
//	func (Calculator) Description() string {
//		return "Performs basic arithmetic operations."
//	}
//
//	func (Calculator) InputSchema() tool.Schema {
//		return tool.Schema{
//			Type: "object",
//			Properties: map[string]tool.Property{
//				"a": {
//					Type:        "number",
//					Description: "First number.",
//				},
//				"b": {
//					Type:        "number",
//					Description: "Second number.",
//				},
//			},
//			Required: []string{"a", "b"},
//		}
//	}
//
//	func (Calculator) Call(
//		ctx context.Context,
//		input map[string]any,
//	) (any, error) {
//		// Execute the tool.
//		return nil, nil
//	}
type Tool interface {
	// Name returns the unique name of the tool.
	Name() string

	// Description returns a human-readable description of what the tool does.
	// This description can be provided to an LLM when selecting tools.
	Description() string

	// InputSchema returns the schema describing the tool's input.
	InputSchema() Schema

	// Call executes the tool with the provided input.
	Call(
		ctx context.Context,
		input map[string]any,
	) (any, error)
}

// Schema describes the input expected by a Tool.
//
// Type usually follows JSON Schema type names such as "object", "string",
// "number", "integer", "boolean", or "array".
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