package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

// this is where the validation happens!
func preHandler(spec Spec) gin.HandlerFunc {
	return func(c *gin.Context) {
		val := spec.Validate
		if val.Query != nil {
			values := c.Request.URL.Query()
			for field, schema := range val.Query {
				// query values are always strings, so we must try to convert
				queryValue := values.Get(field)

				// don't try to convert if the field is empty
				if queryValue == "" {
					if schema.required != nil && *schema.required {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrRequired))
					}
					return
				}
				var convertedValue interface{}
				switch schema.kind {
				case KindBoolean:
					if queryValue == "true" {
						convertedValue = true
					} else if queryValue == "false" {
						convertedValue = false
					} else {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindString:
					convertedValue = queryValue
				case KindNumber:
					var err error
					convertedValue, err = strconv.ParseFloat(queryValue, 64)
					if err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindInteger:
					var err error
					convertedValue, err = strconv.Atoi(queryValue)
					if err != nil {
						c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, ErrWrongType))
						return
					}
				case KindArray:
					// TODO I'm not sure how this works yet
					c.AbortWithStatusJSON(http.StatusNotImplemented, "TODO")
					return
				default:
					c.AbortWithStatusJSON(400, fmt.Sprintf("Validation not possible due to kind: %v", schema.kind))
				}
				if err := schema.Validate(convertedValue); err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Query validation failed for field %v: %v", field, err.Error()))
					return
				}
			}
		}

		if val.Body != nil {
			var body map[string]interface{}
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatusJSON(400, err.Error())
				return
			}
			for field, schema := range val.Body {
				value := body[field]
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
			// TODO strip unknown/unexpected fields
			data, _ := json.Marshal(body)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
		}

		if val.Path != nil {
			for field, schema := range val.Path {
				path := c.Param(field)
				if schema.required != nil && *schema.required && path == "" {
					c.AbortWithStatusJSON(400, fmt.Sprintf("Missing path param"))
					return
				}
			}
		}
	}
}
