package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
)

// this is where the validation happens!
func validationMiddleware(spec Spec) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := spec.Validate
		if val.Query.kind == KindObject { // not sure how any other type makes sense
			values := c.Request.URL.Query()
			for field, schema := range val.Query.obj {
				// query values are always strings, so we must try to convert
				queryValue := values.Get(field)

				convertedValue, err := convert(queryValue, schema)
				if err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
					return
				}
				if err = schema.Validate(convertedValue); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
					return
				}
			}
		}

		if val.Body.Initialized() && val.Body.kind != KindFile {
			var body interface{}
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatusJSON(400, err.Error())
				return
			}
			switch v := body.(type) {
			case string:
				if err := val.Body.Validate(v); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for body: %v", err.Error()))
					return
				}
			case bool:
				if err := val.Body.Validate(v); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for body: %v", err.Error()))
					return
				}
			case float64:
				if err := val.Body.Validate(v); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for body: %v", err.Error()))
					return
				}
			case []interface{}:
				if err := val.Body.Validate(v); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for body: %v", err.Error()))
					return
				}
			case map[string]interface{}:
				for field, schema := range val.Body.obj {
					value := v[field]
					if value == nil {
						if schema.required != nil && *schema.required {
							c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for field %v: %v", field, ErrRequired))
							return
						}
						continue
					}

					if schema.kind == KindInteger {
						// JSON doesn't have integers, so Go treats these fields as float64.
						// Need to convert to integer before validating it.
						switch value.(type) {
						case float64:
							v := value.(float64)
							// check to see if the number can be represented as an integer
							if v != float64(int64(v)) {
								c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for field %v: %v", field, ErrWrongType))
								return
							}
							value = int(value.(float64))
						default:
							c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for field %v: %v", field, ErrWrongType))
							return
						}
					}
					if err := schema.Validate(value); err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed for field %v: %v", field, err.Error()))
						return
					}
				}
			default:
				c.AbortWithStatusJSON(400, fmt.Sprintf("Body validation failed: %v", ErrWrongType))
				return
			}
			// TODO strip unknown/unexpected fields option
			data, _ := json.Marshal(body)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
		}

		if val.Path.kind == KindObject {
			for field, schema := range val.Path.obj {
				path := c.Param(field)

				convertedValue, err := convert(path, schema)
				if err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
					return
				}
				if err = schema.Validate(convertedValue); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
					return
				}
			}
		}
	}
}

func convert(inputValue string, schema Field) (interface{}, error) {
	// don't try to convert if the field is empty
	if inputValue == "" {
		if schema.required != nil && *schema.required {
			return nil, ErrRequired
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
			return nil, ErrWrongType
		}
	case KindString:
		convertedValue = inputValue
	case KindNumber:
		var err error
		convertedValue, err = strconv.ParseFloat(inputValue, 64)
		if err != nil {
			return nil, ErrWrongType
		}
	case KindInteger:
		var err error
		convertedValue, err = strconv.Atoi(inputValue)
		if err != nil {
			return nil, ErrWrongType
		}
	case KindArray:
		// TODO I'm not sure how this works yet
		return nil, ErrNotImplemented
	default:
		return nil, fmt.Errorf("unknown kind: %v", schema.kind)
	}
	return convertedValue, nil
}
