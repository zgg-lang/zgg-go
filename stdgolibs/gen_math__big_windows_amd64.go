package stdgolibs

import (
	pkg "math/big"

	"reflect"
)

func init() {
	registerValues("math/big", map[string]reflect.Value{
		// Functions
		"NewFloat":   reflect.ValueOf(pkg.NewFloat),
		"NewInt":     reflect.ValueOf(pkg.NewInt),
		"Jacobi":     reflect.ValueOf(pkg.Jacobi),
		"NewRat":     reflect.ValueOf(pkg.NewRat),
		"ParseFloat": reflect.ValueOf(pkg.ParseFloat),

		// Consts

		"MaxExp":        reflect.ValueOf(pkg.MaxExp),
		"MinExp":        reflect.ValueOf(pkg.MinExp),
		"MaxPrec":       reflect.ValueOf(pkg.MaxPrec),
		"ToNearestEven": reflect.ValueOf(pkg.ToNearestEven),
		"ToNearestAway": reflect.ValueOf(pkg.ToNearestAway),
		"ToZero":        reflect.ValueOf(pkg.ToZero),
		"AwayFromZero":  reflect.ValueOf(pkg.AwayFromZero),
		"ToNegativeInf": reflect.ValueOf(pkg.ToNegativeInf),
		"ToPositiveInf": reflect.ValueOf(pkg.ToPositiveInf),
		"Below":         reflect.ValueOf(pkg.Below),
		"Exact":         reflect.ValueOf(pkg.Exact),
		"Above":         reflect.ValueOf(pkg.Above),
		"MaxBase":       reflect.ValueOf(pkg.MaxBase),

		// Variables

	})
	registerTypes("math/big", map[string]reflect.Type{
		// Non interfaces

		"Float":        reflect.TypeOf((*pkg.Float)(nil)).Elem(),
		"ErrNaN":       reflect.TypeOf((*pkg.ErrNaN)(nil)).Elem(),
		"RoundingMode": reflect.TypeOf((*pkg.RoundingMode)(nil)).Elem(),
		"Accuracy":     reflect.TypeOf((*pkg.Accuracy)(nil)).Elem(),
		"Int":          reflect.TypeOf((*pkg.Int)(nil)).Elem(),
		"Rat":          reflect.TypeOf((*pkg.Rat)(nil)).Elem(),
		"Word":         reflect.TypeOf((*pkg.Word)(nil)).Elem(),
	})
}
