package stdgolibs

import (
	pkg "image/color/palette"

	"reflect"
)

func init() {
	registerValues("image/color/palette", map[string]reflect.Value{
		// Functions

		// Consts

		// Variables

		"Plan9":   reflect.ValueOf(&pkg.Plan9),
		"WebSafe": reflect.ValueOf(&pkg.WebSafe),
	})
	registerTypes("image/color/palette", map[string]reflect.Type{
		// Non interfaces

	})
}
