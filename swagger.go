package crud

type Swagger struct {
	// Swagger is the version of OpenAPI (should always be 2.0)
	Swagger string `json:"swagger"`
	// Info is metadata for this API
	Info Info `json:"info"`
	// Host can specify the host for cases where the docs are on a different server.
	Host string `json:"host,omitempty"`
	// BasePath defines the base path to be used by the docs. Useful for hosting multiple sites on the same host.
	BasePath string `json:"basePath,omitempty"`

	// Schemes must be: "http", "https", "ws", "wss".
	// If the schemes is not included, the default scheme to be used is the one used to access the
	// Swagger definition itself.
	Schemes []string `json:"schemes,omitempty"`
	// Consumes is a list of MIME types this service consumes. Paths may override this.
	Consumes []string `json:"consumes,omitempty"`
	// Produces is a list of MIME types this service produces. Paths may override this.
	Produces []string `json:"produces,omitempty"`

	// Paths is the paths supported by this API. The key must start with a slash.
	Paths map[string]*Path `json:"paths"`
	// Definitions contains all schemas defined by crud.
	Definitions map[string]Schema `json:"definitions"`

	// SecurityDefinitions defines schemes used by this API.
	SecurityDefinitions map[string]SecurityScheme `json:"securityDefinitions,omitempty"`
	// Security defines the security that needs to be used across the entire API.
	// The map key is the name that corresponds to a name in SecurityDefinitions.
	// If oauth2 is used, an array of required scopes is listed, otherwise
	// use an empty list.
	Security []map[string][]string `json:"security,omitempty"`

	// Tags can be used to provide metadata to the tags.
	Tags []Tag `json:"tags,omitempty"`
	// ExternalDocs defines additional docs.
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`

	// These features are not supported as it seems to be only useful when writing swagger by hand.
	//Parameters map[string]Parameter `json:"parameters,omitempty"`
	//Responses  map[string]Response  `json:"responses,omitempty"`
}

type Info struct {
	// Required fields:

	// Title is the name of this API
	Title string `json:"title"`
	// Version is the version of this API
	Version string `json:"version"`

	// Optional fields:

	// Description of this API
	Description string `json:"description,omitempty"`
	// TermsOfService for this API
	TermsOfService string `json:"termsOfService,omitempty"`
	// Contact for the API
	Contact *Contact `json:"contact,omitempty"`
	// License for using the API
	License *License `json:"license,omitempty"`
}

type Contact struct {
	// Name is the contact's name
	Name string `json:"name,omitempty"`
	// URL to contact information
	URL string `json:"url,omitempty"`
	// Email to the contact
	Email string `json:"email,omitempty"`
}

type License struct {
	// Name is the name of the license
	Name string `json:"name"`
	// URL is a url to the license of the API
	URL string `json:"url,omitempty"`
}

// Schema is a superset of JSON Schema, that OpenAPI uses.
type Schema struct {
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Description string            `json:"description,omitempty"`
	Minimum     float64           `json:"minimum,omitempty"`
	Maximum     float64           `json:"maximum,omitempty"`
	Enum        []interface{}     `json:"enum,omitempty"`
	Default     interface{}       `json:"default,omitempty"`
	Pattern     string            `json:"pattern,omitempty"`
}

// JsonSchema is deprecated...
// Deprecated: use Schema instead
type JsonSchema = Schema

type Path struct {
	Get     *Operation `json:"get,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
	Options *Operation `json:"options,omitempty"`

	// These are not supported as they seem to be conveniences for writing swagger by hand.
	//Ref *Ref `json:"$ref,omitempty"`
	//Parameters interface{} `json:"parameters,omitempty"`
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
	Items            *Schema       `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

type Response struct {
	Schema      Schema `json:"schema"`
	Description string `json:"description"`

	Example interface{}       `json:"interface,omitempty"`
	Ref     *Ref              `json:"$ref,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

var defaultResponse = map[string]Response{
	"default": {
		Schema:      Schema{Type: "string"},
		Description: "Successful",
	},
}

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs *ExternalDoc `json:"externalDocs,omitempty"`
}

type ExternalDoc struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

type SecurityScheme struct {
	// Type is required and must be one of: basic", "apiKey" or "oauth2"
	Type string `json:"type"`
	// Description describes the security scheme
	Description string `json:"description,omitempty"`

	// Required fields for apiKey
	// Name of the header or query param containing the api key.
	Name string
	// In should be one of "query" or "header" to define where to look for the key.
	In string

	// Required fields for oauth2
	// Flow should be: "implicit", "password", "application" or "accessCode"
	Flow string
	// AuthorizationURL is required for "implicit" or "accessCode" flows.
	AuthorizationURL string
	// TokenURL is required for "password", "application", or "accessCode" flows.
	TokenURL string
	// Scopes lists the available scopes for an OAuth2 security scheme. Key is the scope, Value is a description.
	Scopes map[string]string
}
