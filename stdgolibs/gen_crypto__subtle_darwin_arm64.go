package stdgolibs

import (
	pkg "crypto/subtle"

	"reflect"
)

func init() {
	registerValues("crypto/subtle", map[string]reflect.Value{
		// Functions
		"ConstantTimeCompare":  reflect.ValueOf(pkg.ConstantTimeCompare),
		"ConstantTimeSelect":   reflect.ValueOf(pkg.ConstantTimeSelect),
		"ConstantTimeByteEq":   reflect.ValueOf(pkg.ConstantTimeByteEq),
		"ConstantTimeEq":       reflect.ValueOf(pkg.ConstantTimeEq),
		"ConstantTimeCopy":     reflect.ValueOf(pkg.ConstantTimeCopy),
		"ConstantTimeLessOrEq": reflect.ValueOf(pkg.ConstantTimeLessOrEq),

		// Consts

		// Variables

	})
	registerTypes("crypto/subtle", map[string]reflect.Type{
		// Non interfaces

	})
}
