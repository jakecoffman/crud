package crud

type Swagger struct {
	Swagger string `json:"swagger"`
	Info    Info   `json:"info"`

	Paths       map[string]*Path      `json:"paths"`
	Definitions map[string]JsonSchema `json:"definitions"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type JsonSchema struct {
	Type        string                `json:"type,omitempty"`
	Properties  map[string]JsonSchema `json:"properties,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Example     interface{}           `json:"example,omitempty"`
	Description string                `json:"description,omitempty"`
	Minimum     float64               `json:"minimum,omitempty"`
	Maximum     float64               `json:"maximum,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty"`
	Default     interface{}           `json:"default"`
}

type Path struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`
}

type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Description string              `json:"description"`
	Summary     string              `json:"summary"`
}

type Parameter struct {
	In   string `json:"in"`
	Name string `json:"name"`

	Type   string `json:"type,omitempty"`
	Schema *Ref   `json:"schema,omitempty"`

	Required    *bool         `json:"required,omitempty"`
	Description string        `json:"description,omitempty"`
	Minimum     *float64      `json:"minimum,omitempty"`
	Maximum     *float64      `json:"maximum,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

type Response struct {
	Schema      JsonSchema `json:"schema"`
	Description string     `json:"description"`
}

var DefaultResponse = map[string]Response{
	"default": {
		Schema:      JsonSchema{Type: "string"},
		Description: "Successful",
	},
}
