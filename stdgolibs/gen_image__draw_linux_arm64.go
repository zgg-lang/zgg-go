package stdgolibs

import (
	pkg "image/draw"

	"reflect"
)

func init() {
	registerValues("image/draw", map[string]reflect.Value{
		// Functions
		"Draw":     reflect.ValueOf(pkg.Draw),
		"DrawMask": reflect.ValueOf(pkg.DrawMask),

		// Consts

		"Over": reflect.ValueOf(pkg.Over),
		"Src":  reflect.ValueOf(pkg.Src),

		// Variables

		"FloydSteinberg": reflect.ValueOf(&pkg.FloydSteinberg),
	})
	registerTypes("image/draw", map[string]reflect.Type{
		// Non interfaces

		"Op": reflect.TypeOf((*pkg.Op)(nil)).Elem(),
	})
}
