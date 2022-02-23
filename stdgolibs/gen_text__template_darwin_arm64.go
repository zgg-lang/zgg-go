package stdgolibs

import (
	pkg "text/template"

	"reflect"
)

func init() {
	registerValues("text/template", map[string]reflect.Value{
		// Functions
		"IsTrue":           reflect.ValueOf(pkg.IsTrue),
		"HTMLEscape":       reflect.ValueOf(pkg.HTMLEscape),
		"HTMLEscapeString": reflect.ValueOf(pkg.HTMLEscapeString),
		"HTMLEscaper":      reflect.ValueOf(pkg.HTMLEscaper),
		"JSEscape":         reflect.ValueOf(pkg.JSEscape),
		"JSEscapeString":   reflect.ValueOf(pkg.JSEscapeString),
		"JSEscaper":        reflect.ValueOf(pkg.JSEscaper),
		"URLQueryEscaper":  reflect.ValueOf(pkg.URLQueryEscaper),
		"Must":             reflect.ValueOf(pkg.Must),
		"ParseFiles":       reflect.ValueOf(pkg.ParseFiles),
		"ParseGlob":        reflect.ValueOf(pkg.ParseGlob),
		"ParseFS":          reflect.ValueOf(pkg.ParseFS),
		"New":              reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("text/template", map[string]reflect.Type{
		// Non interfaces

		"ExecError": reflect.TypeOf((*pkg.ExecError)(nil)).Elem(),
		"FuncMap":   reflect.TypeOf((*pkg.FuncMap)(nil)).Elem(),
		"Template":  reflect.TypeOf((*pkg.Template)(nil)).Elem(),
	})
}
