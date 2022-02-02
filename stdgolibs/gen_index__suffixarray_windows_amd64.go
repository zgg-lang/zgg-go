package stdgolibs

import (
	pkg "index/suffixarray"

	"reflect"
)

func init() {
	registerValues("index/suffixarray", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),

		// Consts

		// Variables

	})
	registerTypes("index/suffixarray", map[string]reflect.Type{
		// Non interfaces

		"Index": reflect.TypeOf((*pkg.Index)(nil)).Elem(),
	})
}
