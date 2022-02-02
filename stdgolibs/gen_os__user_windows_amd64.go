package stdgolibs

import (
	pkg "os/user"

	"reflect"
)

func init() {
	registerValues("os/user", map[string]reflect.Value{
		// Functions
		"Current":       reflect.ValueOf(pkg.Current),
		"Lookup":        reflect.ValueOf(pkg.Lookup),
		"LookupId":      reflect.ValueOf(pkg.LookupId),
		"LookupGroup":   reflect.ValueOf(pkg.LookupGroup),
		"LookupGroupId": reflect.ValueOf(pkg.LookupGroupId),

		// Consts

		// Variables

	})
	registerTypes("os/user", map[string]reflect.Type{
		// Non interfaces

		"User":                reflect.TypeOf((*pkg.User)(nil)).Elem(),
		"Group":               reflect.TypeOf((*pkg.Group)(nil)).Elem(),
		"UnknownUserIdError":  reflect.TypeOf((*pkg.UnknownUserIdError)(nil)).Elem(),
		"UnknownUserError":    reflect.TypeOf((*pkg.UnknownUserError)(nil)).Elem(),
		"UnknownGroupIdError": reflect.TypeOf((*pkg.UnknownGroupIdError)(nil)).Elem(),
		"UnknownGroupError":   reflect.TypeOf((*pkg.UnknownGroupError)(nil)).Elem(),
	})
}
