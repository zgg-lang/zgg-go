package stdgolibs

import (
	pkg "go/scanner"

	"reflect"
)

func init() {
	registerValues("go/scanner", map[string]reflect.Value{
		// Functions
		"PrintError": reflect.ValueOf(pkg.PrintError),

		// Consts

		"ScanComments": reflect.ValueOf(pkg.ScanComments),

		// Variables

	})
	registerTypes("go/scanner", map[string]reflect.Type{
		// Non interfaces

		"Error":        reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"ErrorList":    reflect.TypeOf((*pkg.ErrorList)(nil)).Elem(),
		"ErrorHandler": reflect.TypeOf((*pkg.ErrorHandler)(nil)).Elem(),
		"Scanner":      reflect.TypeOf((*pkg.Scanner)(nil)).Elem(),
		"Mode":         reflect.TypeOf((*pkg.Mode)(nil)).Elem(),
	})
}
