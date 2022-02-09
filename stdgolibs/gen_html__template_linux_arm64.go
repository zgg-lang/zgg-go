package stdgolibs

import (
	pkg "html/template"

	"reflect"
)

func init() {
	registerValues("html/template", map[string]reflect.Value{
		// Functions
		"New":              reflect.ValueOf(pkg.New),
		"Must":             reflect.ValueOf(pkg.Must),
		"ParseFiles":       reflect.ValueOf(pkg.ParseFiles),
		"ParseGlob":        reflect.ValueOf(pkg.ParseGlob),
		"IsTrue":           reflect.ValueOf(pkg.IsTrue),
		"ParseFS":          reflect.ValueOf(pkg.ParseFS),
		"HTMLEscape":       reflect.ValueOf(pkg.HTMLEscape),
		"HTMLEscapeString": reflect.ValueOf(pkg.HTMLEscapeString),
		"HTMLEscaper":      reflect.ValueOf(pkg.HTMLEscaper),
		"JSEscape":         reflect.ValueOf(pkg.JSEscape),
		"JSEscapeString":   reflect.ValueOf(pkg.JSEscapeString),
		"JSEscaper":        reflect.ValueOf(pkg.JSEscaper),
		"URLQueryEscaper":  reflect.ValueOf(pkg.URLQueryEscaper),

		// Consts

		"OK":                   reflect.ValueOf(pkg.OK),
		"ErrAmbigContext":      reflect.ValueOf(pkg.ErrAmbigContext),
		"ErrBadHTML":           reflect.ValueOf(pkg.ErrBadHTML),
		"ErrBranchEnd":         reflect.ValueOf(pkg.ErrBranchEnd),
		"ErrEndContext":        reflect.ValueOf(pkg.ErrEndContext),
		"ErrNoSuchTemplate":    reflect.ValueOf(pkg.ErrNoSuchTemplate),
		"ErrOutputContext":     reflect.ValueOf(pkg.ErrOutputContext),
		"ErrPartialCharset":    reflect.ValueOf(pkg.ErrPartialCharset),
		"ErrPartialEscape":     reflect.ValueOf(pkg.ErrPartialEscape),
		"ErrRangeLoopReentry":  reflect.ValueOf(pkg.ErrRangeLoopReentry),
		"ErrSlashAmbig":        reflect.ValueOf(pkg.ErrSlashAmbig),
		"ErrPredefinedEscaper": reflect.ValueOf(pkg.ErrPredefinedEscaper),

		// Variables

	})
	registerTypes("html/template", map[string]reflect.Type{
		// Non interfaces

		"Error":     reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"ErrorCode": reflect.TypeOf((*pkg.ErrorCode)(nil)).Elem(),
		"Template":  reflect.TypeOf((*pkg.Template)(nil)).Elem(),
		"FuncMap":   reflect.TypeOf((*pkg.FuncMap)(nil)).Elem(),
		"CSS":       reflect.TypeOf((*pkg.CSS)(nil)).Elem(),
		"HTML":      reflect.TypeOf((*pkg.HTML)(nil)).Elem(),
		"HTMLAttr":  reflect.TypeOf((*pkg.HTMLAttr)(nil)).Elem(),
		"JS":        reflect.TypeOf((*pkg.JS)(nil)).Elem(),
		"JSStr":     reflect.TypeOf((*pkg.JSStr)(nil)).Elem(),
		"URL":       reflect.TypeOf((*pkg.URL)(nil)).Elem(),
		"Srcset":    reflect.TypeOf((*pkg.Srcset)(nil)).Elem(),
	})
}
