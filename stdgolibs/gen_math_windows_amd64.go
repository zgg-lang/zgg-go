package stdgolibs

import (
	pkg "math"

	"reflect"
)

func init() {
	registerValues("math", map[string]reflect.Value{
		// Functions
		"Ldexp":           reflect.ValueOf(pkg.Ldexp),
		"Logb":            reflect.ValueOf(pkg.Logb),
		"Ilogb":           reflect.ValueOf(pkg.Ilogb),
		"Modf":            reflect.ValueOf(pkg.Modf),
		"Asinh":           reflect.ValueOf(pkg.Asinh),
		"Floor":           reflect.ValueOf(pkg.Floor),
		"Ceil":            reflect.ValueOf(pkg.Ceil),
		"Trunc":           reflect.ValueOf(pkg.Trunc),
		"Round":           reflect.ValueOf(pkg.Round),
		"RoundToEven":     reflect.ValueOf(pkg.RoundToEven),
		"FMA":             reflect.ValueOf(pkg.FMA),
		"Atanh":           reflect.ValueOf(pkg.Atanh),
		"Copysign":        reflect.ValueOf(pkg.Copysign),
		"Cos":             reflect.ValueOf(pkg.Cos),
		"Sin":             reflect.ValueOf(pkg.Sin),
		"Gamma":           reflect.ValueOf(pkg.Gamma),
		"Log":             reflect.ValueOf(pkg.Log),
		"Mod":             reflect.ValueOf(pkg.Mod),
		"Sqrt":            reflect.ValueOf(pkg.Sqrt),
		"Atan2":           reflect.ValueOf(pkg.Atan2),
		"Cbrt":            reflect.ValueOf(pkg.Cbrt),
		"Dim":             reflect.ValueOf(pkg.Dim),
		"Max":             reflect.ValueOf(pkg.Max),
		"Min":             reflect.ValueOf(pkg.Min),
		"Expm1":           reflect.ValueOf(pkg.Expm1),
		"Pow10":           reflect.ValueOf(pkg.Pow10),
		"Atan":            reflect.ValueOf(pkg.Atan),
		"Inf":             reflect.ValueOf(pkg.Inf),
		"NaN":             reflect.ValueOf(pkg.NaN),
		"IsNaN":           reflect.ValueOf(pkg.IsNaN),
		"IsInf":           reflect.ValueOf(pkg.IsInf),
		"Nextafter32":     reflect.ValueOf(pkg.Nextafter32),
		"Nextafter":       reflect.ValueOf(pkg.Nextafter),
		"J0":              reflect.ValueOf(pkg.J0),
		"Y0":              reflect.ValueOf(pkg.Y0),
		"Lgamma":          reflect.ValueOf(pkg.Lgamma),
		"Log10":           reflect.ValueOf(pkg.Log10),
		"Log2":            reflect.ValueOf(pkg.Log2),
		"Hypot":           reflect.ValueOf(pkg.Hypot),
		"Log1p":           reflect.ValueOf(pkg.Log1p),
		"Sinh":            reflect.ValueOf(pkg.Sinh),
		"Cosh":            reflect.ValueOf(pkg.Cosh),
		"Tanh":            reflect.ValueOf(pkg.Tanh),
		"Asin":            reflect.ValueOf(pkg.Asin),
		"Acos":            reflect.ValueOf(pkg.Acos),
		"Erf":             reflect.ValueOf(pkg.Erf),
		"Erfc":            reflect.ValueOf(pkg.Erfc),
		"Frexp":           reflect.ValueOf(pkg.Frexp),
		"Exp":             reflect.ValueOf(pkg.Exp),
		"Exp2":            reflect.ValueOf(pkg.Exp2),
		"J1":              reflect.ValueOf(pkg.J1),
		"Y1":              reflect.ValueOf(pkg.Y1),
		"Jn":              reflect.ValueOf(pkg.Jn),
		"Yn":              reflect.ValueOf(pkg.Yn),
		"Pow":             reflect.ValueOf(pkg.Pow),
		"Signbit":         reflect.ValueOf(pkg.Signbit),
		"Abs":             reflect.ValueOf(pkg.Abs),
		"Acosh":           reflect.ValueOf(pkg.Acosh),
		"Sincos":          reflect.ValueOf(pkg.Sincos),
		"Float32bits":     reflect.ValueOf(pkg.Float32bits),
		"Float32frombits": reflect.ValueOf(pkg.Float32frombits),
		"Float64bits":     reflect.ValueOf(pkg.Float64bits),
		"Float64frombits": reflect.ValueOf(pkg.Float64frombits),
		"Erfinv":          reflect.ValueOf(pkg.Erfinv),
		"Erfcinv":         reflect.ValueOf(pkg.Erfcinv),
		"Remainder":       reflect.ValueOf(pkg.Remainder),
		"Tan":             reflect.ValueOf(pkg.Tan),

		// Consts

		"E":                      reflect.ValueOf(pkg.E),
		"Pi":                     reflect.ValueOf(pkg.Pi),
		"Phi":                    reflect.ValueOf(pkg.Phi),
		"Sqrt2":                  reflect.ValueOf(pkg.Sqrt2),
		"SqrtE":                  reflect.ValueOf(pkg.SqrtE),
		"SqrtPi":                 reflect.ValueOf(pkg.SqrtPi),
		"SqrtPhi":                reflect.ValueOf(pkg.SqrtPhi),
		"Ln2":                    reflect.ValueOf(pkg.Ln2),
		"Log2E":                  reflect.ValueOf(pkg.Log2E),
		"Ln10":                   reflect.ValueOf(pkg.Ln10),
		"Log10E":                 reflect.ValueOf(pkg.Log10E),
		"MaxFloat32":             reflect.ValueOf(pkg.MaxFloat32),
		"SmallestNonzeroFloat32": reflect.ValueOf(pkg.SmallestNonzeroFloat32),
		"MaxFloat64":             reflect.ValueOf(pkg.MaxFloat64),
		"SmallestNonzeroFloat64": reflect.ValueOf(pkg.SmallestNonzeroFloat64),
		"MaxInt8":                reflect.ValueOf(pkg.MaxInt8),
		"MinInt8":                reflect.ValueOf(pkg.MinInt8),
		"MaxInt16":               reflect.ValueOf(pkg.MaxInt16),
		"MinInt16":               reflect.ValueOf(pkg.MinInt16),
		"MaxInt32":               reflect.ValueOf(pkg.MaxInt32),
		"MinInt32":               reflect.ValueOf(pkg.MinInt32),
		"MaxInt64":               reflect.ValueOf(pkg.MaxInt64),
		"MinInt64":               reflect.ValueOf(pkg.MinInt64),
		"MaxUint8":               reflect.ValueOf(pkg.MaxUint8),
		"MaxUint16":              reflect.ValueOf(pkg.MaxUint16),
		"MaxUint32":              reflect.ValueOf(pkg.MaxUint32),
		"MaxUint64":              reflect.ValueOf(uint64(pkg.MaxUint64)),

		// Variables

	})
	registerTypes("math", map[string]reflect.Type{
		// Non interfaces

	})
}
