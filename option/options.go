package option

// Option configures a router option. Use the convenience constructors below.
type Option struct {
	StripUnknown *bool
	AllowUnknown *bool
}

// StripUnknown will remove unknown fields if true, leave them if false. Defaults to true.
func StripUnknown(v bool) Option {
	return Option{StripUnknown: &v}
}

// AllowUnknown false will cause the validation to fail if it encounters an unknown field. Defaults to true.
func AllowUnknown(v bool) Option {
	return Option{AllowUnknown: &v}
}
