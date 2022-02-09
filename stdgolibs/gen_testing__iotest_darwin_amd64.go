package stdgolibs

import (
	pkg "testing/iotest"

	"reflect"
)

func init() {
	registerValues("testing/iotest", map[string]reflect.Value{
		// Functions
		"OneByteReader":  reflect.ValueOf(pkg.OneByteReader),
		"HalfReader":     reflect.ValueOf(pkg.HalfReader),
		"DataErrReader":  reflect.ValueOf(pkg.DataErrReader),
		"TimeoutReader":  reflect.ValueOf(pkg.TimeoutReader),
		"ErrReader":      reflect.ValueOf(pkg.ErrReader),
		"TestReader":     reflect.ValueOf(pkg.TestReader),
		"TruncateWriter": reflect.ValueOf(pkg.TruncateWriter),
		"NewWriteLogger": reflect.ValueOf(pkg.NewWriteLogger),
		"NewReadLogger":  reflect.ValueOf(pkg.NewReadLogger),

		// Consts

		// Variables

		"ErrTimeout": reflect.ValueOf(&pkg.ErrTimeout),
	})
	registerTypes("testing/iotest", map[string]reflect.Type{
		// Non interfaces

	})
}
