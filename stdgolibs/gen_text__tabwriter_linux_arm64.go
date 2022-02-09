package stdgolibs

import (
	pkg "text/tabwriter"

	"reflect"
)

func init() {
	registerValues("text/tabwriter", map[string]reflect.Value{
		// Functions
		"NewWriter": reflect.ValueOf(pkg.NewWriter),

		// Consts

		"FilterHTML":          reflect.ValueOf(pkg.FilterHTML),
		"StripEscape":         reflect.ValueOf(pkg.StripEscape),
		"AlignRight":          reflect.ValueOf(pkg.AlignRight),
		"DiscardEmptyColumns": reflect.ValueOf(pkg.DiscardEmptyColumns),
		"TabIndent":           reflect.ValueOf(pkg.TabIndent),
		"Debug":               reflect.ValueOf(pkg.Debug),
		"Escape":              reflect.ValueOf(pkg.Escape),

		// Variables

	})
	registerTypes("text/tabwriter", map[string]reflect.Type{
		// Non interfaces

		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
	})
}
