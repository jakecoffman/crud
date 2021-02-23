package crud

import (
	"fmt"
)

type Field struct {
	Type       string      `json:"type"`
	Maximum    *float64    `json:"maximum,omitempty"`
	Minimum    *float64    `json:"minimum,omitempty"`
	IsRequired *bool       `json:"required,omitempty"`
	Ex         interface{} `json:"example,omitempty"`
}

var (
	ErrRequired  = fmt.Errorf("value is required")
	ErrWrongType = fmt.Errorf("wrong type passed")
	ErrMaximum   = fmt.Errorf("maximum exceeded")
	ErrMinimum   = fmt.Errorf("minumum exceeded")
)

func (f *Field) Validate(value interface{}) error {
	if value == nil && *f.IsRequired {
		return ErrRequired
	}

	switch v := value.(type) {
	case int:
		if f.Type != "number" {
			return ErrWrongType
		}
		if f.Maximum != nil && float64(v) > *f.Maximum {
			return ErrMaximum
		}
		if f.Minimum != nil && float64(v) < *f.Minimum {
			return ErrMinimum
		}
	case float64:
		if f.Type != "number" {
			return ErrWrongType
		}
		if f.Maximum != nil && v > *f.Maximum {
			return ErrMaximum
		}
		if f.Minimum != nil && v < *f.Minimum {
			return ErrMinimum
		}
	case string:
		if f.Type != "string" {
			return ErrWrongType
		}
	default:
		return ErrWrongType
	}

	return nil
}

func Number() Field {
	return Field{Type: "number"}
}

func String() Field {
	return Field{Type: "string"}
}

func (f Field) Min(min float64) Field {
	f.Minimum = &min
	return f
}

func (f Field) Max(max float64) Field {
	f.Maximum = &max
	return f
}

func (f Field) Required() Field {
	required := true
	f.IsRequired = &required
	return f
}

func (f Field) Example(ex interface{}) Field {
	f.Ex = ex
	return f
}

func ToJsonSchema(fields map[string]Field) JsonSchema {
	schema := JsonSchema{
		Type:       "object",
		Properties: map[string]JsonSchema{},
	}

	for name, field := range fields {
		schema.Properties[name] = JsonSchema{
			Type:    field.Type,
			Example: field.Ex,
		}
		if field.IsRequired != nil && *field.IsRequired {
			schema.Required = append(schema.Required, name)
		}
	}

	return schema
}
