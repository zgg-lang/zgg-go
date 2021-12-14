package stdgolibs

import (
	pkg "hash/adler32"

	"reflect"
)

func init() {
	registerValues("hash/adler32", map[string]reflect.Value{
		// Functions
		"New":      reflect.ValueOf(pkg.New),
		"Checksum": reflect.ValueOf(pkg.Checksum),

		// Consts

		"Size": reflect.ValueOf(pkg.Size),

		// Variables

	})
	registerTypes("hash/adler32", map[string]reflect.Type{
		// Non interfaces

	})
}
