package stdgolibs

import (
	pkg "net/mail"

	"reflect"
)

func init() {
	registerValues("net/mail", map[string]reflect.Value{
		// Functions
		"ReadMessage":      reflect.ValueOf(pkg.ReadMessage),
		"ParseDate":        reflect.ValueOf(pkg.ParseDate),
		"ParseAddress":     reflect.ValueOf(pkg.ParseAddress),
		"ParseAddressList": reflect.ValueOf(pkg.ParseAddressList),

		// Consts

		// Variables

		"ErrHeaderNotPresent": reflect.ValueOf(&pkg.ErrHeaderNotPresent),
	})
	registerTypes("net/mail", map[string]reflect.Type{
		// Non interfaces

		"Message":       reflect.TypeOf((*pkg.Message)(nil)).Elem(),
		"Header":        reflect.TypeOf((*pkg.Header)(nil)).Elem(),
		"Address":       reflect.TypeOf((*pkg.Address)(nil)).Elem(),
		"AddressParser": reflect.TypeOf((*pkg.AddressParser)(nil)).Elem(),
	})
}
