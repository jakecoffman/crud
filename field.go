package crud

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Field allows specification of swagger or json schema types using the builder pattern.
type Field struct {
	kind        string
	format      string
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
	strip       *bool
	unknown     *bool
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

// Kind returns the kind of the field.
func (f Field) Kind() string {
	return f.kind
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
	errUnknown      = fmt.Errorf("unknown value")
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
		if f.kind != KindInteger {
			return errWrongType
		}
		if f.max != nil && float64(v) > *f.max {
			return errMaximum
		}
		if f.min != nil && float64(v) < *f.min {
			return errMinimum
		}
	case float64:
		if f.kind == KindInteger {
			// since JSON is unmarshalled as float64 always
			if float64(int(v)) != v {
				return errWrongType
			}
		} else if f.kind != KindNumber {
			return errWrongType
		}
		if f.max != nil && v > *f.max {
			return errMaximum
		}
		if f.min != nil && v < *f.min {
			return errMinimum
		}
	case string:
		if f.kind != KindString {
			return errWrongType
		}
		if f.required != nil && *f.required && v == "" && !f.allow.has("") {
			return errRequired
		}
		if f.max != nil && len(v) > int(*f.max) {
			return errMaximum
		}
		if f.min != nil && len(v) < int(*f.min) {
			return errMinimum
		}
		switch f.format {
		case FormatDateTime:
			_, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}
		case FormatDate:
			_, err := time.Parse(fullDate, v)
			if err != nil {
				return err
			}
		}
	case bool:
		if f.kind != KindBoolean {
			return errWrongType
		}
	case []interface{}:
		if f.kind != KindArray {
			return errWrongType
		}
		if f.min != nil && float64(len(v)) < *f.min {
			return errMinimum
		}
		if f.max != nil && float64(len(v)) > *f.max {
			return errMaximum
		}
		if f.arr != nil {
			// child fields inherit parent's settings, unless specified on child
			if f.arr.strip == nil {
				f.arr.strip = f.strip
			}
			if f.arr.unknown == nil {
				f.arr.unknown = f.unknown
			}
			for _, item := range v {
				if err := f.arr.Validate(item); err != nil {
					return err
				}
			}
		}
	case map[string]interface{}:
		if f.kind != KindObject {
			return errWrongType
		}
		return validateObject("", f, v)
	default:
		return fmt.Errorf("unhandled type %v", v)
	}

	if f.enum != nil && !f.enum.has(value) {
		return errEnumNotFound
	}

	return nil
}

// validateObject is a recursive function that validates the field values in the object. It also
// performs stripping of values, or erroring when unexpected fields are present, depending on the
// options on the fields.
func validateObject(name string, field *Field, input interface{}) error {
	switch v := input.(type) {
	case nil:
		if field.required != nil && *field.required {
			return fmt.Errorf("object validation failed for field %v: %w", name, errRequired)
		}
	case string, bool:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("object validation failed for field %v: %w", name, err)
		}
	case float64:
		if field.kind == KindInteger {
			// JSON doesn't have integers, so Go treats these fields as float64.
			// Need to convert to integer before validating it.
			if v != float64(int64(v)) {
				return fmt.Errorf("object validation failed for field %v: %w", name, errWrongType)
			}
			if err := field.Validate(int(v)); err != nil {
				return fmt.Errorf("object validation failed for field %v: %w", name, err)
			}
		} else {
			if err := field.Validate(v); err != nil {
				return fmt.Errorf("object validation failed for field %v: %w", name, err)
			}
		}
	case []interface{}:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("object validation failed for field %v: %w", name, err)
		}
		if field.arr != nil {
			for i, item := range v {
				if err := validateObject(fmt.Sprintf("%v[%v]", name, i), field.arr, item); err != nil {
					return err
				}
			}
		}
	case map[string]interface{}:
		if !field.isAllowUnknown() {
			for key := range v {
				if _, ok := field.obj[key]; !ok {
					return fmt.Errorf("unknown field in object: %v %w", key, errUnknown)
				}
			}
		}

		if field.isStripUnknown() {
			for key := range v {
				if _, ok := field.obj[key]; !ok {
					delete(v, key)
				}
			}
		}

		for childName, childField := range field.obj {
			// child fields inherit parent's settings, unless specified on child
			if childField.strip == nil {
				childField.strip = field.strip
			}
			if childField.unknown == nil {
				childField.unknown = field.unknown
			}

			newV := v[childName]
			if newV == nil && childField.required != nil && *childField.required {
				return fmt.Errorf("object validation failed for field %v.%v: %w", name, childName, errRequired)
			} else if newV == nil && childField._default != nil {
				v[childName] = childField._default
			} else if err := validateObject(name+"."+childName, &childField, v[childName]); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("object validation failed for type %v: %w", reflect.TypeOf(v), errWrongType)
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

// DateTime creates a field with dateTime type
func DateTime() Field {
	return Field{kind: KindString, format: FormatDateTime}
}

// https://xml2rfc.tools.ietf.org/public/rfc/html/rfc3339.html#anchor14
const fullDate = "2006-01-02"

// Date creates a field with date type
func Date() Field {
	return Field{kind: KindString, format: FormatDate}
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
	if f.max != nil && *f.max < min {
		panic("min cannot be larger than max")
	}
	return f
}

// Max specifies a maximum value for this field
func (f Field) Max(max float64) Field {
	f.max = &max
	if f.min != nil && *f.min > max {
		panic("min cannot be larger than max")
	}
	return f
}

// Required specifies the field must be provided. Can't be used with Default.
func (f Field) Required() Field {
	if f._default != nil {
		panic("required and default cannot be used together")
	}

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

// Default specifies a default value to use if the field is nil. Can't be used with Required.
func (f Field) Default(value interface{}) Field {
	if f.required != nil && *f.required {
		panic("default and required cannot be used together")
	}

	switch value.(type) {
	case int:
		if f.kind != KindInteger {
			panic("wrong type passed default")
		}
	case float64:
		if f.kind != KindNumber {
			panic("wrong type passed default")
		}
	case string:
		if f.kind != KindString {
			panic("wrong type passed default")
		}
	case bool:
		if f.kind != KindBoolean {
			panic("wrong type passed default")
		}
	default:
		panic("default must be an int, float64, bool or string")
	}
	f._default = value
	return f
}

// Enum restricts the field's values to the set of values specified
func (f Field) Enum(values ...interface{}) Field {
	if f.kind == KindArray || f.kind == KindObject || f.kind == KindFile {
		panic("Enum cannot be used on arrays, objects, files")
	}
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

const (
	FormatDate     = "date"
	FormatDateTime = "dateTime"
)

// Format is used to set custom format types. Note that formats with special
// validation in this library also have their own constructor. See DateTime for example.
func (f Field) Format(format string) Field {
	f.format = format
	return f
}

// Allow lets you break rules
// For example, String().Required() excludes "", unless you Allow("")
func (f Field) Allow(values ...interface{}) Field {
	f.allow = append(f.allow, values...)
	return f
}

// Strip overrides the global "strip unknown" setting just for this field, and all children of this field
func (f Field) Strip(strip bool) Field {
	f.strip = &strip
	return f
}

// Unknown overrides the global "allow unknown" setting just for this field, and all children of this field
func (f Field) Unknown(allow bool) Field {
	f.unknown = &allow
	return f
}

// ToSwaggerParameters transforms a field into a slice of Parameter.
func (f *Field) ToSwaggerParameters(in string) (parameters []Parameter) {
	switch f.kind {
	case KindArray:
		p := Parameter{
			In:               in,
			Type:             f.kind,
			CollectionFormat: "multi",
			Required:         f.required,
			Description:      f.description,
			Default:          f._default,
		}
		if f.arr != nil {
			items := f.arr.ToJsonSchema()
			p.Items = &items
		}
		parameters = append(parameters, p)
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
				if field.arr != nil {
					temp := field.arr.ToJsonSchema()
					param.Items = &temp
				}
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

// ToJsonSchema transforms a field into a Swagger Schema.
// TODO this is an extension of JsonSchema, rename in v2 ToSchema() Schema
func (f *Field) ToJsonSchema() JsonSchema {
	schema := JsonSchema{
		Type: f.kind,
	}

	switch f.kind {
	case KindArray:
		if f.arr != nil {
			items := f.arr.ToJsonSchema()
			schema.Items = &items
		}
	case KindObject:
		schema.Properties = map[string]JsonSchema{}
		for name, field := range f.obj {
			prop := JsonSchema{
				Type:        field.kind,
				Format:      field.format,
				Example:     field.example,
				Description: field.description,
				Default:     field._default,
			}
			if field.example == nil {
				if field.kind == KindString {
					switch field.format {
					case FormatDateTime:
						prop.Example = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)
					case FormatDate:
						prop.Example = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local).Format(fullDate)
					}
				}
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
				if field.arr != nil {
					items := field.arr.ToJsonSchema()
					prop.Items = &items
				}
			}
			schema.Properties[name] = prop
		}
	}
	return schema
}

func (f Field) isAllowUnknown() bool {
	if f.unknown == nil {
		return true // by default allow unknown
	}
	return *f.unknown
}

func (f Field) isStripUnknown() bool {
	if f.strip == nil {
		return true // by default strip
	}
	return *f.strip
}
