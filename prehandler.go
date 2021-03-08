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
					query[field] = []string{schema._default.(string)}
				}
				continue
			}
			if len(queryValue) > 1 {
				if schema.kind != KindArray {
					return fmt.Errorf("query validation failed for field %v: %w", field, errWrongType)
				}
			}
			if schema.kind == KindArray {
				// sadly we have to convert to a []interface{} to simplify the validate code
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
		err := r.validateBody("body", &val.Body, body)
		if err != nil {
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

func (r *Router) validateBody(name string, field *Field, body interface{}) error {
	switch v := body.(type) {
	case nil:
		if field.required != nil && *field.required {
			return fmt.Errorf("body validation failed for field %v: %w", name, errRequired)
		}
	case string:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed for field %v: %w", name, err)
		}
	case bool:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed for field %v: %w", name, err)
		}
	case float64:
		if field.kind == KindInteger {
			// JSON doesn't have integers, so Go treats these fields as float64.
			// Need to convert to integer before validating it.
			if v != float64(int64(v)) {
				return fmt.Errorf("body validation failed for field %v: %w", name, errWrongType)
			}
			if err := field.Validate(int(v)); err != nil {
				return fmt.Errorf("body validation failed for field %v: %w", name, err)
			}
		} else {
			if err := field.Validate(v); err != nil {
				return fmt.Errorf("body validation failed for field %v: %w", name, err)
			}
		}
	case []interface{}:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed for field %v: %w", name, err)
		}
		if field.arr != nil {
			for i, item := range v {
				if err := r.validateBody(fmt.Sprintf("%v[%v]", name, i), field.arr, item); err != nil {
					return err
				}
			}
		}
	case map[string]interface{}:
		if !r.allowUnknown {
			for key := range v {
				if _, ok := field.obj[key]; !ok {
					return fmt.Errorf("unknown field in body: %v %w", key, errUnknown)
				}
			}
		}

		if r.stripUnknown {
			for key := range v {
				if _, ok := field.obj[key]; !ok {
					delete(v, key)
				}
			}
		}

		for name, field := range field.obj {
			newV := v[name]
			if newV == nil && field.required != nil && *field.required {
				return fmt.Errorf("body validation failed for field %v: %w", name, errRequired)
			} else if newV == nil && field._default != nil {
				v[name] = field._default
			} else if err := r.validateBody(name, &field, v[name]); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("body validation failed: %w", errWrongType)
	}
	return nil
}

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
