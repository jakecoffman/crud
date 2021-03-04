package crud

import (
	"fmt"
	"strings"
)

// Field allows specification of swagger or json schema types using the builder pattern.
type Field struct {
	kind        string
	obj         map[string]Field
	max         *float64
	min         *float64
	required    *bool
	example     interface{}
	description string
	enum        enum
	_default    interface{}
	arr         *Field
	allow       enum
}

func (f Field) String() string {
	var str strings.Builder
	str.WriteString(fmt.Sprintf("{Field: '%v'", f.kind))
	if f.required != nil {
		str.WriteString(" required")
	}
	str.WriteString("}")
	return str.String()
}

// Initialized returns true if the field has been initialized with Number, String, etc.
// When the Swagger is being built, often an uninitialized field will be ignored.
func (f Field) Initialized() bool {
	return f.kind != ""
}

type enum []interface{}

func (e enum) has(needle interface{}) bool {
	for _, value := range e {
		if value == needle {
			return true
		}
	}
	return false
}

var (
	errRequired     = fmt.Errorf("value is required")
	errWrongType    = fmt.Errorf("wrong type passed")
	errMaximum      = fmt.Errorf("maximum exceeded")
	errMinimum      = fmt.Errorf("minimum exceeded")
	errEnumNotFound = fmt.Errorf("value not in enum")
)

// Validate is used in the validation middleware to tell if the value passed
// into the controller meets the restrictions set on the field.
func (f *Field) Validate(value interface{}) error {
	if value == nil && f.required != nil && *f.required {
		return errRequired
	}
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int:
		if f.kind != "integer" {
			return errWrongType
		}
		if f.max != nil && float64(v) > *f.max {
			return errMaximum
		}
		if f.min != nil && float64(v) < *f.min {
			return errMinimum
		}
	case float64:
		if f.kind != "number" {
			return errWrongType
		}
		if f.max != nil && v > *f.max {
			return errMaximum
		}
		if f.min != nil && v < *f.min {
			return errMinimum
		}
	case string:
		if f.kind != "string" {
			return errWrongType
		}
		if f.required != nil && *f.required && v == "" && !f.allow.has("") {
			return errRequired
		}
	case bool:
		if f.kind != "boolean" {
			return errWrongType
		}
	case []interface{}:
		if f.kind != "array" {
			return errWrongType
		}
		if f.min != nil && float64(len(v)) < *f.min {
			return errMinimum
		}
		if f.max != nil && float64(len(v)) > *f.max {
			return errMaximum
		}
	default:
		return fmt.Errorf("unhandled type %v", v)
	}

	if f.enum != nil && !f.enum.has(value) {
		return errEnumNotFound
	}

	return nil
}

// These kinds correlate to swagger and json types.
const (
	KindNumber  = "number"
	KindString  = "string"
	KindBoolean = "boolean"
	KindObject  = "object"
	KindArray   = "array"
	KindFile    = "file"
	KindInteger = "integer"
)

// Number creates a field with floating point type
func Number() Field {
	return Field{kind: KindNumber}
}

// String creates a field with string type
func String() Field {
	return Field{kind: KindString}
}

// Boolean creates a field with boolean type
func Boolean() Field {
	return Field{kind: KindBoolean}
}

// Object creates a field with object type
func Object(obj map[string]Field) Field {
	return Field{kind: KindObject, obj: obj}
}

// Array creates a field with array type
func Array() Field {
	return Field{kind: KindArray}
}

// File creates a field with file type
func File() Field {
	return Field{kind: KindFile}
}

// Integer creates a field with integer type
func Integer() Field {
	return Field{kind: KindInteger}
}

// Min specifies a minimum value for this field
func (f Field) Min(min float64) Field {
	f.min = &min
	return f
}

// Max specifies a maximum value for this field
func (f Field) Max(max float64) Field {
	f.max = &max
	return f
}

// Required specifies the field must be provided
func (f Field) Required() Field {
	required := true
	f.required = &required
	return f
}

// Example specifies an example value for the swagger to display
func (f Field) Example(ex interface{}) Field {
	f.example = ex
	return f
}

// Description specifies a human-readable explanation of the field
func (f Field) Description(description string) Field {
	f.description = description
	return f
}

// Enum restricts the field's values to the set of values specified
func (f Field) Enum(values ...interface{}) Field {
	f.enum = values
	return f
}

// Items specifies the type of elements in an array
func (f Field) Items(item Field) Field {
	if f.kind != KindArray {
		panic("Items can only be used with array types")
	}
	f.arr = &item
	return f
}

// Allow lets you break rules
// For example, String().Required() excludes "", unless you Allow("")
func (f Field) Allow(values ...interface{}) Field {
	f.allow = append(f.allow, values...)
	return f
}

// ToSwaggerParameters transforms a field into a slice of Parameter.
func (f *Field) ToSwaggerParameters(in string) (parameters []Parameter) {
	switch f.kind {
	case KindArray:
		items := f.arr.ToJsonSchema()
		parameters = append(parameters, Parameter{
			In:               in,
			Type:             f.kind,
			Items:            &items,
			CollectionFormat: "multi",
			Required:         f.required,
			Description:      f.description,
			Default:          f._default,
		})
	case KindObject:
		for name, field := range f.obj {
			param := Parameter{
				In:          in,
				Name:        name,
				Type:        field.kind,
				Required:    field.required,
				Description: field.description,
				Default:     field._default,
				Enum:        field.enum,
				Minimum:     field.min,
				Maximum:     field.max,
			}
			if field.kind == KindArray {
				temp := field.arr.ToJsonSchema()
				param.Items = &temp
				param.CollectionFormat = "multi"
			}
			if field.kind == KindObject {
				// TODO
			}
			parameters = append(parameters, param)
		}
	}
	return
}

// ToJsonSchema transforms a field into a JsonSchema.
func (f *Field) ToJsonSchema() JsonSchema {
	schema := JsonSchema{
		Type: f.kind,
	}

	switch f.kind {
	case KindArray:
		items := f.arr.ToJsonSchema()
		schema.Items = &items
	case KindObject:
		schema.Properties = map[string]JsonSchema{}
		for name, field := range f.obj {
			prop := JsonSchema{
				Type:        field.kind,
				Example:     field.example,
				Description: field.description,
				Default:     field._default,
			}
			if field.required != nil && *field.required {
				schema.Required = append(schema.Required, name)
			}
			if field.min != nil {
				prop.Minimum = *field.min
			}
			if field.max != nil {
				prop.Maximum = *field.max
			}
			if prop.Type == KindArray {
				items := field.arr.ToJsonSchema()
				prop.Items = &items
			}
			schema.Properties[name] = prop
		}
	}
	return schema
}
