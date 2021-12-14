package stdgolibs

import (
	pkg "path"

	"reflect"
)

func init() {
	registerValues("path", map[string]reflect.Value{
		// Functions
		"Match": reflect.ValueOf(pkg.Match),
		"Clean": reflect.ValueOf(pkg.Clean),
		"Split": reflect.ValueOf(pkg.Split),
		"Join":  reflect.ValueOf(pkg.Join),
		"Ext":   reflect.ValueOf(pkg.Ext),
		"Base":  reflect.ValueOf(pkg.Base),
		"IsAbs": reflect.ValueOf(pkg.IsAbs),
		"Dir":   reflect.ValueOf(pkg.Dir),

		// Consts

		// Variables

		"ErrBadPattern": reflect.ValueOf(&pkg.ErrBadPattern),
	})
	registerTypes("path", map[string]reflect.Type{
		// Non interfaces

	})
}
