package crud

import (
	"fmt"
	"net/url"
	"strconv"
)

// Validate checks the spec against the inputs and returns an error if it finds one.
func (r *Router) Validate(val Validate, query url.Values, body interface{}, path map[string]string) error {
	if val.Query.kind == KindObject { // not sure how any other type makes sense
		for field, schema := range val.Query.obj {
			// query values are always strings, so we must try to convert
			queryValue := query[field]

			if len(queryValue) == 0 {
				if schema.required != nil && *schema.required {
					return fmt.Errorf("query validation failed for field %v: %w", field, errRequired)
				}
				if schema._default != nil {
					query[field] = []string{fmt.Sprint(schema._default)}
				}
				continue
			}
			if len(queryValue) > 1 {
				if schema.kind != KindArray {
					return fmt.Errorf("query validation failed for field %v: %w", field, errWrongType)
				}
			}
			if schema.kind == KindArray {
				// sadly we have to convert to a []interface{} to simplify the validation code
				var intray []interface{}
				for _, v := range queryValue {
					intray = append(intray, v)
				}
				if err := schema.Validate(intray); err != nil {
					return fmt.Errorf("query validation failed for field %v: %w", field, err)
				}
				if schema.arr != nil {
					for _, v := range queryValue {
						convertedValue, err := convert(v, *schema.arr)
						if err != nil {
							return fmt.Errorf("query validation failed for field %v: %w", field, err)
						}
						if err = schema.arr.Validate(convertedValue); err != nil {
							return fmt.Errorf("query validation failed for field %v: %w", field, err)
						}
					}
				}
			} else {
				convertedValue, err := convert(queryValue[0], schema)
				if err != nil {
					return fmt.Errorf("query validation failed for field %v: %w", field, err)
				}
				if err = schema.Validate(convertedValue); err != nil {
					return fmt.Errorf("query validation failed for field %v: %w", field, err)
				}
			}
		}
	}

	if val.Body.Initialized() && val.Body.kind != KindFile {
		// use router defaults if the object doesn't have anything set
		f := val.Body
		if f.strip == nil {
			f = f.Strip(r.stripUnknown)
		}
		if f.unknown == nil {
			f = f.Unknown(r.allowUnknown)
		}
		if err := f.Validate(body); err != nil {
			return err
		}
	}

	if val.Path.kind == KindObject {
		for field, schema := range val.Path.obj {
			param := path[field]

			convertedValue, err := convert(param, schema)
			if err != nil {
				return fmt.Errorf("path validation failed for field %v: %w", field, err)
			}
			if err = schema.Validate(convertedValue); err != nil {
				return fmt.Errorf("path validation failed for field %v: %w", field, err)
			}
		}
	}

	return nil
}

// For certain types of data passed like Query and Header, the value is always
// a string. So this function attempts to convert the string into the desired field kind.
func convert(inputValue string, schema Field) (interface{}, error) {
	// don't try to convert if the field is empty
	if inputValue == "" {
		if schema.required != nil && *schema.required {
			return nil, errRequired
		}
		return nil, nil
	}
	var convertedValue interface{}
	switch schema.kind {
	case KindBoolean:
		if inputValue == "true" {
			convertedValue = true
		} else if inputValue == "false" {
			convertedValue = false
		} else {
			return nil, errWrongType
		}
	case KindString:
		convertedValue = inputValue
	case KindNumber:
		var err error
		convertedValue, err = strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return nil, errWrongType
		}
	case KindInteger:
		var err error
		convertedValue, err = strconv.Atoi(inputValue)
		if err != nil {
			return nil, errWrongType
		}
	default:
		return nil, fmt.Errorf("unknown kind: %v", schema.kind)
	}
	return convertedValue, nil
}
