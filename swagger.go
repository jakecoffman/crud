package crud

type Swagger struct {
	Swagger  string `json:"swagger"`
	Info     Info   `json:"info"`
	BasePath string `json:"basePath,omitempty"`

	Paths       map[string]*Path      `json:"paths"`
	Definitions map[string]JsonSchema `json:"definitions"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

type JsonSchema struct {
	Type        string                `json:"type,omitempty"`
	Format      string                `json:"format,omitempty"`
	Properties  map[string]JsonSchema `json:"properties,omitempty"`
	Items       *JsonSchema           `json:"items,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Example     interface{}           `json:"example,omitempty"`
	Description string                `json:"description,omitempty"`
	Minimum     float64               `json:"minimum,omitempty"`
	Maximum     float64               `json:"maximum,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty"`
	Default     interface{}           `json:"default,omitempty"`
	Pattern     string                `json:"pattern,omitempty"`
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

	// one of path, query, header, body, or form
	Type   string `json:"type,omitempty"`
	Schema *Ref   `json:"schema,omitempty"`

	Required         *bool         `json:"required,omitempty"`
	Description      string        `json:"description,omitempty"`
	Minimum          *float64      `json:"minimum,omitempty"`
	Maximum          *float64      `json:"maximum,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Items            *JsonSchema   `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

type Response struct {
	Schema      JsonSchema `json:"schema"`
	Description string     `json:"description"`

	Example interface{}       `json:"interface,omitempty"`
	Ref     *Ref              `json:"$ref,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

var defaultResponse = map[string]Response{
	"default": {
		Schema:      JsonSchema{Type: "string"},
		Description: "Successful",
	},
}
