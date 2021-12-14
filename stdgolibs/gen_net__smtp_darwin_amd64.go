package stdgolibs

import (
	pkg "net/smtp"

	"reflect"
)

func init() {
	registerValues("net/smtp", map[string]reflect.Value{
		// Functions
		"PlainAuth":   reflect.ValueOf(pkg.PlainAuth),
		"CRAMMD5Auth": reflect.ValueOf(pkg.CRAMMD5Auth),
		"Dial":        reflect.ValueOf(pkg.Dial),
		"NewClient":   reflect.ValueOf(pkg.NewClient),
		"SendMail":    reflect.ValueOf(pkg.SendMail),

		// Consts

		// Variables

	})
	registerTypes("net/smtp", map[string]reflect.Type{
		// Non interfaces

		"ServerInfo": reflect.TypeOf((*pkg.ServerInfo)(nil)).Elem(),
		"Client":     reflect.TypeOf((*pkg.Client)(nil)).Elem(),
	})
}
