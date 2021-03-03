package crud

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func query(query string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://example.com"+query, nil)
	return w, c
}

func TestQueryValidation(t *testing.T) {
	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected int
	}{
		{
			Schema: map[string]Field{
				"testquery": String(),
			},
			Input:    "",
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "",
			Expected: 400,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "?testquery=",
			Expected: 400,
		}, {
			Schema: map[string]Field{
				"testquery": String().Required(),
			},
			Input:    "?testquery=ok",
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"testquery": Number(),
			},
			Input:    "",
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "?testquery=1",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Number().Required(),
			},
			Input:    "?testquery=1.1",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Number(),
			},
			Input:    "?testquery=a",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "?testquery=true",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "?testquery=false",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Boolean(),
			},
			Input:    "?testquery=1",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "?testquery=1",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Max(1),
			},
			Input:    "?testquery=2",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Min(5),
			},
			Input:    "?testquery=4",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "?testquery=1.1",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer(),
			},
			Input:    "?testquery=a",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Enum(1, 2),
			},
			Input:    "?testquery=2",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": Integer().Enum(1, 2),
			},
			Input:    "?testquery=3",
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"testquery": String().Enum("a"),
			},
			Input:    "?testquery=a",
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"testquery": String().Enum("a"),
			},
			Input:    "?testquery=b",
			Expected: 400,
		},
	}

	for _, test := range tests {
		handler := validationMiddleware(Spec{
			Validate: Validate{Query: Object(test.Schema)},
		})

		w, c := query(test.Input)
		handler(c)

		if w.Result().StatusCode != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, w.Code, test.Input, test.Schema)
		}
	}
}

func body(payload string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://example.com", strings.NewReader(payload))
	return w, c
}

func TestSimpleBodyValidation(t *testing.T) {
	tests := []struct {
		Schema   Field
		Input    string
		Expected int
	}{
		{
			Schema:   Number(),
			Input:    "1",
			Expected: 200,
		},
		{
			Schema:   Number(),
			Input:    "a",
			Expected: 400,
		},
		{
			Schema:   String(),
			Input:    `"2"`,
			Expected: 200,
		},
		{
			Schema:   Boolean(),
			Input:    `true`,
			Expected: 200,
		},
		{
			Schema:   Boolean(),
			Input:    `false`,
			Expected: 200,
		},
		{
			Schema:   Boolean(),
			Input:    `1`,
			Expected: 400,
		},
	}

	for _, test := range tests {
		handler := validationMiddleware(Spec{
			Validate: Validate{Body: test.Schema},
		})

		w, c := body(test.Input)
		handler(c)

		if w.Result().StatusCode != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, w.Code, test.Input, test.Schema)
		}
	}
}

func TestBodyValidation(t *testing.T) {
	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected int
	}{
		{
			Schema: map[string]Field{
				"str": String().Required(),
			},
			Input:    `{"str":""}`,
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"int": Integer(),
			},
			Input:    `{}`,
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{}`,
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{"int":1}`,
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    `{"int":1.9}`,
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"float": Number().Required(),
			},
			Input:    `{"float":-1}`,
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"float": Number().Required(),
			},
			Input:    `{"float":1.1}`,
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"obj1": Object(map[string]Field{
					"inner": Number().Required(),
				}),
			},
			Input:    `{"obj1":{"inner":1}}`,
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"obj2": Object(map[string]Field{
					"inner": Number().Required(),
				}),
			},
			Input:    `{"obj2":{"inner":"not a number"}}`,
			Expected: 400,
		}, {
			Schema: map[string]Field{
				"arr1": Array().Items(Number()),
			},
			Input:    `{"arr1":[1]}`,
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"arr2": Array().Items(Number()),
			},
			Input:    `{"arr2":["a"]}`,
			Expected: 400,
		}, {
			Schema: map[string]Field{
				"complex1": Object(map[string]Field{
					"array": Array().Required().Items(Object(map[string]Field{
						"id": Number().Required(),
					})),
				}).Required(),
			},
			Input:    `{"complex1":{"array":[{"id":1}]}}`,
			Expected: 200,
		}, {
			Schema: map[string]Field{
				"complex2": Object(map[string]Field{
					"array": Array().Required().Items(Object(map[string]Field{
						"id": Number().Required(),
					})),
				}).Required(),
			},
			Input:    `{"complex2":{"array":[{"id":"a"}]}}`,
			Expected: 400,
		},
	}

	for _, test := range tests {
		handler := validationMiddleware(Spec{
			Validate: Validate{Body: Object(test.Schema)},
		})

		w, c := body(test.Input)
		handler(c)

		if w.Result().StatusCode != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, w.Code, test.Input, test.Schema)
		}
	}
}

func path(pathValue string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "http://example.com", nil)
	c.Params = []gin.Param{{Key: "id", Value: pathValue}}
	return w, c
}

func TestPathValidation(t *testing.T) {
	tests := []struct {
		Schema   map[string]Field
		Input    string
		Expected int
	}{
		{
			Schema: map[string]Field{
				"id": Integer(),
			},
			Input:    ``,
			Expected: 200,
		},
		{
			Schema: map[string]Field{
				"int": Integer().Required(),
			},
			Input:    ``,
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"id": Integer().Required(),
			},
			Input:    `a`,
			Expected: 400,
		},
		{
			Schema: map[string]Field{
				"id": Integer().Required(),
			},
			Input:    `1`,
			Expected: 200,
		},
	}

	for _, test := range tests {
		handler := validationMiddleware(Spec{
			Validate: Validate{Path: Object(test.Schema)},
		})

		w, c := path(test.Input)
		handler(c)

		if w.Result().StatusCode != test.Expected {
			t.Errorf("expected '%v' got '%v'. input: '%v'. schema: '%v'", test.Expected, w.Code, test.Input, test.Schema)
		}
	}
}
