package stdgolibs

import (
	pkg "go/printer"

	"reflect"
)

func init() {
	registerValues("go/printer", map[string]reflect.Value{
		// Functions
		"Fprint": reflect.ValueOf(pkg.Fprint),

		// Consts

		"RawFormat": reflect.ValueOf(pkg.RawFormat),
		"TabIndent": reflect.ValueOf(pkg.TabIndent),
		"UseSpaces": reflect.ValueOf(pkg.UseSpaces),
		"SourcePos": reflect.ValueOf(pkg.SourcePos),

		// Variables

	})
	registerTypes("go/printer", map[string]reflect.Type{
		// Non interfaces

		"Mode":          reflect.TypeOf((*pkg.Mode)(nil)).Elem(),
		"Config":        reflect.TypeOf((*pkg.Config)(nil)).Elem(),
		"CommentedNode": reflect.TypeOf((*pkg.CommentedNode)(nil)).Elem(),
	})
}
