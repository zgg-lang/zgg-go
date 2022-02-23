package stdgolibs

import (
	pkg "image/gif"

	"reflect"
)

func init() {
	registerValues("image/gif", map[string]reflect.Value{
		// Functions
		"Decode":       reflect.ValueOf(pkg.Decode),
		"DecodeAll":    reflect.ValueOf(pkg.DecodeAll),
		"DecodeConfig": reflect.ValueOf(pkg.DecodeConfig),
		"EncodeAll":    reflect.ValueOf(pkg.EncodeAll),
		"Encode":       reflect.ValueOf(pkg.Encode),

		// Consts

		"DisposalNone":       reflect.ValueOf(pkg.DisposalNone),
		"DisposalBackground": reflect.ValueOf(pkg.DisposalBackground),
		"DisposalPrevious":   reflect.ValueOf(pkg.DisposalPrevious),

		// Variables

	})
	registerTypes("image/gif", map[string]reflect.Type{
		// Non interfaces

		"GIF":     reflect.TypeOf((*pkg.GIF)(nil)).Elem(),
		"Options": reflect.TypeOf((*pkg.Options)(nil)).Elem(),
	})
}
