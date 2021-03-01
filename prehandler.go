package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/url"
	"strconv"
)

// this is where the validation happens!
func validationMiddleware(spec Spec) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := spec.Validate
		var query url.Values
		var body interface{}
		var path map[string]string

		if val.Path.Initialized() {
			path = map[string]string{}
			for _, param := range c.Params {
				path[param.Key] = param.Value
			}
		}

		if val.Body.Initialized() && val.Body.kind != KindFile {
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatusJSON(400, err.Error())
				return
			}
			// TODO strip unknown/unexpected fields option
			data, _ := json.Marshal(body)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
		}

		if val.Query.Initialized() {
			query = c.Request.URL.Query()
		}

		if err := validate(val, query, body, path); err != nil {
			c.AbortWithStatusJSON(400, err.Error())
		}
	}
}

func validate(val Validate, query url.Values, body interface{}, path map[string]string) error {
	if val.Query.kind == KindObject { // not sure how any other type makes sense
		for field, schema := range val.Query.obj {
			// query values are always strings, so we must try to convert
			queryValue := query[field]

			if len(queryValue) == 0 {
				if schema.required != nil && *schema.required {
					return fmt.Errorf("query validation failed for field %v: %v", field, errRequired)
				}
			} else if len(queryValue) > 1 {
				if schema.arr == nil {
					return fmt.Errorf("query validation failed for field %v: %v", field, errWrongType)
				}
				// TODO validate each item in the array
			} else {
				convertedValue, err := convert(queryValue[0], schema)
				if err != nil {
					return fmt.Errorf("query validation failed for field %v: %v", field, err.Error())
				}
				if err = schema.Validate(convertedValue); err != nil {
					return fmt.Errorf("query validation failed for field %v: %v", field, err.Error())
				}
			}
		}
	}

	if val.Body.Initialized() && val.Body.kind != KindFile {
		err := validateBody(&val.Body, body)
		if err != nil {
			return err
		}
	}

	if val.Path.kind == KindObject {
		for field, schema := range val.Path.obj {
			param := path[field]

			convertedValue, err := convert(param, schema)
			if err != nil {
				return fmt.Errorf("query validation failed for field %v: %v", field, err.Error())
			}
			if err = schema.Validate(convertedValue); err != nil {
				return fmt.Errorf("query validation failed for field %v: %v", field, err.Error())
			}
		}
	}

	return nil
}

func validateBody(field *Field, body interface{}) error {
	switch v := body.(type) {
	case nil:
		if field.required != nil && *field.required {
			return fmt.Errorf("body validation failed for field %v: %v", field, errRequired)
		}
	case string:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed: %v", err.Error())
		}
	case bool:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed: %v", err.Error())
		}
	case float64:
		if field.kind == KindInteger {
			// JSON doesn't have integers, so Go treats these fields as float64.
			// Need to convert to integer before validating it.
			if v != float64(int64(v)) {
				return fmt.Errorf("body validation failed for field %v: %v", field, errWrongType)
			}
			if err := field.Validate(int(v)); err != nil {
				return fmt.Errorf("body validation failed: %v", err.Error())
			}
		} else {
			if err := field.Validate(v); err != nil {
				return fmt.Errorf("body validation failed: %v", err.Error())
			}
		}
	case []interface{}:
		if err := field.Validate(v); err != nil {
			return fmt.Errorf("body validation failed: %v", err.Error())
		}
	case map[string]interface{}:
		for name, field := range field.obj {
			if err := validateBody(&field, v[name]); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("body validation failed: %v", errWrongType)
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
	case KindArray:
		// TODO convert each item in the array
	default:
		return nil, fmt.Errorf("unknown kind: %v", schema.kind)
	}
	return convertedValue, nil
}
