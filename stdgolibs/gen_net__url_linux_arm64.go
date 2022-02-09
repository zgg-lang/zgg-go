package stdgolibs

import (
	pkg "net/url"

	"reflect"
)

func init() {
	registerValues("net/url", map[string]reflect.Value{
		// Functions
		"QueryUnescape":   reflect.ValueOf(pkg.QueryUnescape),
		"PathUnescape":    reflect.ValueOf(pkg.PathUnescape),
		"QueryEscape":     reflect.ValueOf(pkg.QueryEscape),
		"PathEscape":      reflect.ValueOf(pkg.PathEscape),
		"User":            reflect.ValueOf(pkg.User),
		"UserPassword":    reflect.ValueOf(pkg.UserPassword),
		"Parse":           reflect.ValueOf(pkg.Parse),
		"ParseRequestURI": reflect.ValueOf(pkg.ParseRequestURI),
		"ParseQuery":      reflect.ValueOf(pkg.ParseQuery),

		// Consts

		// Variables

	})
	registerTypes("net/url", map[string]reflect.Type{
		// Non interfaces

		"Error":            reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"EscapeError":      reflect.TypeOf((*pkg.EscapeError)(nil)).Elem(),
		"InvalidHostError": reflect.TypeOf((*pkg.InvalidHostError)(nil)).Elem(),
		"URL":              reflect.TypeOf((*pkg.URL)(nil)).Elem(),
		"Userinfo":         reflect.TypeOf((*pkg.Userinfo)(nil)).Elem(),
		"Values":           reflect.TypeOf((*pkg.Values)(nil)).Elem(),
	})
}
