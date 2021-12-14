package stdgolibs

import (
	pkg "math/cmplx"

	"reflect"
)

func init() {
	registerValues("math/cmplx", map[string]reflect.Value{
		// Functions
		"Asin":  reflect.ValueOf(pkg.Asin),
		"Asinh": reflect.ValueOf(pkg.Asinh),
		"Acos":  reflect.ValueOf(pkg.Acos),
		"Acosh": reflect.ValueOf(pkg.Acosh),
		"Atan":  reflect.ValueOf(pkg.Atan),
		"Atanh": reflect.ValueOf(pkg.Atanh),
		"Conj":  reflect.ValueOf(pkg.Conj),
		"Phase": reflect.ValueOf(pkg.Phase),
		"Sin":   reflect.ValueOf(pkg.Sin),
		"Sinh":  reflect.ValueOf(pkg.Sinh),
		"Cos":   reflect.ValueOf(pkg.Cos),
		"Cosh":  reflect.ValueOf(pkg.Cosh),
		"Tan":   reflect.ValueOf(pkg.Tan),
		"Tanh":  reflect.ValueOf(pkg.Tanh),
		"Cot":   reflect.ValueOf(pkg.Cot),
		"IsInf": reflect.ValueOf(pkg.IsInf),
		"Inf":   reflect.ValueOf(pkg.Inf),
		"IsNaN": reflect.ValueOf(pkg.IsNaN),
		"NaN":   reflect.ValueOf(pkg.NaN),
		"Polar": reflect.ValueOf(pkg.Polar),
		"Log":   reflect.ValueOf(pkg.Log),
		"Log10": reflect.ValueOf(pkg.Log10),
		"Rect":  reflect.ValueOf(pkg.Rect),
		"Abs":   reflect.ValueOf(pkg.Abs),
		"Exp":   reflect.ValueOf(pkg.Exp),
		"Pow":   reflect.ValueOf(pkg.Pow),
		"Sqrt":  reflect.ValueOf(pkg.Sqrt),

		// Consts

		// Variables

	})
	registerTypes("math/cmplx", map[string]reflect.Type{
		// Non interfaces

	})
}
