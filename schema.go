package crud

import (
	"fmt"
)

type Field struct {
	kind        string
	max         *float64
	min         *float64
	required    *bool
	example     interface{}
	description string
	enum        []interface{}
}

var (
	ErrRequired       = fmt.Errorf("value is required")
	ErrWrongType      = fmt.Errorf("wrong type passed")
	ErrMaximum        = fmt.Errorf("maximum exceeded")
	ErrMinimum        = fmt.Errorf("minumum exceeded")
	ErrNotImplemented = fmt.Errorf("not implemented")
)

func (f *Field) Validate(value interface{}) error {
	if value == nil && f.required != nil && *f.required {
		return ErrRequired
	}
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int:
		if f.kind != "integer" {
			return ErrWrongType
		}
		if f.max != nil && float64(v) > *f.max {
			return ErrMaximum
		}
		if f.min != nil && float64(v) < *f.min {
			return ErrMinimum
		}
	case float64:
		if f.kind != "number" {
			return ErrWrongType
		}
		if f.max != nil && v > *f.max {
			return ErrMaximum
		}
		if f.min != nil && v < *f.min {
			return ErrMinimum
		}
	case string:
		if f.kind != "string" {
			return ErrWrongType
		}
	case bool:
		if f.kind != "boolean" {
			return ErrWrongType
		}
	case []interface{}:
		if f.kind != "array" {
			return ErrWrongType
		}
	default:
		return fmt.Errorf("unhandled type %v", v)
	}

	return nil
}

const (
	KindNumber  = "number"
	KindString  = "string"
	KindBoolean = "boolean"
	KindArray   = "array"
	KindFile    = "file"
	KindInteger = "integer"
)

func Number() Field {
	return Field{kind: KindNumber}
}

func String() Field {
	return Field{kind: KindString}
}

func Boolean() Field {
	return Field{kind: KindBoolean}
}

func Array() Field {
	return Field{kind: KindArray}
}

func File() Field {
	return Field{kind: KindFile}
}

func Integer() Field {
	return Field{kind: KindInteger}
}

func (f Field) Min(min float64) Field {
	f.min = &min
	return f
}

func (f Field) Max(max float64) Field {
	f.max = &max
	return f
}

func (f Field) Required() Field {
	required := true
	f.required = &required
	return f
}

func (f Field) Example(ex interface{}) Field {
	f.example = ex
	return f
}

func (f Field) Description(description string) Field {
	f.description = description
	return f
}

func (f Field) Enum(values ...interface{}) Field {
	f.enum = values
	return f
}

func ToJsonSchema(fields map[string]Field) JsonSchema {
	schema := JsonSchema{
		Type:       "object",
		Properties: map[string]JsonSchema{},
	}

	for name, field := range fields {
		prop := JsonSchema{
			Type:        field.kind,
			Example:     field.example,
			Description: field.description,
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
		schema.Properties[name] = prop
	}

	return schema
}
