package crud

type Swagger struct {
	OpenAPI string `json:"openapi"`
	Info    Info   `json:"info"`

	Paths      map[string]*Path `json:"paths"`
	Components Components       `json:"components"`
}

// Info provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in
// editing or documentation generation tools for convenience.
type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`

	Description    string `json:"description,omitempty"`
	TermsOfService string `json:"termsOfService,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the OAS. All objects defined within the
// components object will have no effect on the API unless they are explicitly referenced from properties outside the
// components object.
type Components struct {
	Schemas         map[string]*JsonSchema `json:"schemas,omitempty"`
	Responses       map[string]*JsonSchema `json:"responses,omitempty"`
	Parameters      map[string]*JsonSchema `json:"parameters,omitempty"`
	Examples        map[string]*JsonSchema `json:"examples,omitempty"`
	RequestBodies   map[string]*JsonSchema `json:"requestBodies,omitempty"`
	Headers         map[string]*JsonSchema `json:"headers,omitempty"`
	SecuritySchemas map[string]*JsonSchema `json:"securitySchemas,omitempty"`
	Links           map[string]*JsonSchema `json:"links,omitempty"`
	Callbacks       map[string]*JsonSchema `json:"callbacks,omitempty"`
}

type JsonSchema struct {
	Type        string                `json:"type,omitempty"`
	Properties  map[string]JsonSchema `json:"properties,omitempty"`
	Items       *JsonSchema           `json:"items,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Example     interface{}           `json:"example,omitempty"`
	Description string                `json:"description,omitempty"`
	Minimum     *float64              `json:"minimum,omitempty"`
	Maximum     *float64              `json:"maximum,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty"`
	Default     interface{}           `json:"default,omitempty"`
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
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Description string              `json:"description,omitempty"`
	Summary     string              `json:"summary,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Content     map[string]MediaType `json:"content"`
	Required    *bool                `json:"required,omitempty"`
}

type Parameter struct {
	In   string `json:"in"`
	Name string `json:"name"`

	Ref    string      `json:"$ref,omitempty"`
	Schema *JsonSchema `json:"schema,omitempty"`

	Required    *bool  `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

type Response struct {
	Description string `json:"description"`

	// Ref to a response in the Components
	Ref     string            `json:"$ref,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`

	Content map[string]MediaType `json:"content,omitempty"` // map content-type to content
}

type MediaType struct {
	// Schema must be a *Ref (reference to components/schemas) or JsonSchema
	// TODO make this type safe
	Schema interface{} `json:"schema"`
	//Encoding map[string]Encoding `json:"encoding"` TODO
	Example  interface{}        `json:"example,omitempty"`
	Examples map[string]Example `json:"examples,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref"`
}

type Example struct {
	Summary       string      `json:"summary,omitempty"`
	Description   string      `json:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty"`
}

var defaultResponse = map[string]Response{
	"default": {
		Description: "Successful",
	},
}
