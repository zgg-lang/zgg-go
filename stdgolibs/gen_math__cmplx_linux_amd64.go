package stdgolibs

import (
	pkg "math/cmplx"

	"reflect"
)

func init() {
	registerValues("math/cmplx", map[string]reflect.Value{
		// Functions
		"Exp":   reflect.ValueOf(pkg.Exp),
		"IsInf": reflect.ValueOf(pkg.IsInf),
		"Inf":   reflect.ValueOf(pkg.Inf),
		"Sin":   reflect.ValueOf(pkg.Sin),
		"Sinh":  reflect.ValueOf(pkg.Sinh),
		"Cos":   reflect.ValueOf(pkg.Cos),
		"Cosh":  reflect.ValueOf(pkg.Cosh),
		"Asin":  reflect.ValueOf(pkg.Asin),
		"Asinh": reflect.ValueOf(pkg.Asinh),
		"Acos":  reflect.ValueOf(pkg.Acos),
		"Acosh": reflect.ValueOf(pkg.Acosh),
		"Atan":  reflect.ValueOf(pkg.Atan),
		"Atanh": reflect.ValueOf(pkg.Atanh),
		"Phase": reflect.ValueOf(pkg.Phase),
		"Polar": reflect.ValueOf(pkg.Polar),
		"Abs":   reflect.ValueOf(pkg.Abs),
		"IsNaN": reflect.ValueOf(pkg.IsNaN),
		"NaN":   reflect.ValueOf(pkg.NaN),
		"Pow":   reflect.ValueOf(pkg.Pow),
		"Conj":  reflect.ValueOf(pkg.Conj),
		"Log":   reflect.ValueOf(pkg.Log),
		"Log10": reflect.ValueOf(pkg.Log10),
		"Rect":  reflect.ValueOf(pkg.Rect),
		"Sqrt":  reflect.ValueOf(pkg.Sqrt),
		"Tan":   reflect.ValueOf(pkg.Tan),
		"Tanh":  reflect.ValueOf(pkg.Tanh),
		"Cot":   reflect.ValueOf(pkg.Cot),

		// Consts

		// Variables

	})
	registerTypes("math/cmplx", map[string]reflect.Type{
		// Non interfaces

	})
}
