package stdgolibs

import (
	pkg "math"

	"reflect"
)

func init() {
	registerValues("math", map[string]reflect.Value{
		// Functions
		"Modf":            reflect.ValueOf(pkg.Modf),
		"Sqrt":            reflect.ValueOf(pkg.Sqrt),
		"Tan":             reflect.ValueOf(pkg.Tan),
		"Tanh":            reflect.ValueOf(pkg.Tanh),
		"Logb":            reflect.ValueOf(pkg.Logb),
		"Ilogb":           reflect.ValueOf(pkg.Ilogb),
		"Erf":             reflect.ValueOf(pkg.Erf),
		"Erfc":            reflect.ValueOf(pkg.Erfc),
		"Exp":             reflect.ValueOf(pkg.Exp),
		"Exp2":            reflect.ValueOf(pkg.Exp2),
		"Log1p":           reflect.ValueOf(pkg.Log1p),
		"Signbit":         reflect.ValueOf(pkg.Signbit),
		"Dim":             reflect.ValueOf(pkg.Dim),
		"Max":             reflect.ValueOf(pkg.Max),
		"Min":             reflect.ValueOf(pkg.Min),
		"Asin":            reflect.ValueOf(pkg.Asin),
		"Acos":            reflect.ValueOf(pkg.Acos),
		"Atanh":           reflect.ValueOf(pkg.Atanh),
		"Inf":             reflect.ValueOf(pkg.Inf),
		"NaN":             reflect.ValueOf(pkg.NaN),
		"IsNaN":           reflect.ValueOf(pkg.IsNaN),
		"IsInf":           reflect.ValueOf(pkg.IsInf),
		"Hypot":           reflect.ValueOf(pkg.Hypot),
		"J0":              reflect.ValueOf(pkg.J0),
		"Y0":              reflect.ValueOf(pkg.Y0),
		"Acosh":           reflect.ValueOf(pkg.Acosh),
		"Expm1":           reflect.ValueOf(pkg.Expm1),
		"Jn":              reflect.ValueOf(pkg.Jn),
		"Yn":              reflect.ValueOf(pkg.Yn),
		"Abs":             reflect.ValueOf(pkg.Abs),
		"FMA":             reflect.ValueOf(pkg.FMA),
		"Frexp":           reflect.ValueOf(pkg.Frexp),
		"Gamma":           reflect.ValueOf(pkg.Gamma),
		"Lgamma":          reflect.ValueOf(pkg.Lgamma),
		"Log":             reflect.ValueOf(pkg.Log),
		"Pow":             reflect.ValueOf(pkg.Pow),
		"Sincos":          reflect.ValueOf(pkg.Sincos),
		"Atan2":           reflect.ValueOf(pkg.Atan2),
		"Floor":           reflect.ValueOf(pkg.Floor),
		"Ceil":            reflect.ValueOf(pkg.Ceil),
		"Trunc":           reflect.ValueOf(pkg.Trunc),
		"Round":           reflect.ValueOf(pkg.Round),
		"RoundToEven":     reflect.ValueOf(pkg.RoundToEven),
		"Mod":             reflect.ValueOf(pkg.Mod),
		"Nextafter32":     reflect.ValueOf(pkg.Nextafter32),
		"Nextafter":       reflect.ValueOf(pkg.Nextafter),
		"Remainder":       reflect.ValueOf(pkg.Remainder),
		"Copysign":        reflect.ValueOf(pkg.Copysign),
		"Ldexp":           reflect.ValueOf(pkg.Ldexp),
		"Cbrt":            reflect.ValueOf(pkg.Cbrt),
		"Atan":            reflect.ValueOf(pkg.Atan),
		"Erfinv":          reflect.ValueOf(pkg.Erfinv),
		"Erfcinv":         reflect.ValueOf(pkg.Erfcinv),
		"J1":              reflect.ValueOf(pkg.J1),
		"Y1":              reflect.ValueOf(pkg.Y1),
		"Log10":           reflect.ValueOf(pkg.Log10),
		"Log2":            reflect.ValueOf(pkg.Log2),
		"Pow10":           reflect.ValueOf(pkg.Pow10),
		"Cos":             reflect.ValueOf(pkg.Cos),
		"Sin":             reflect.ValueOf(pkg.Sin),
		"Sinh":            reflect.ValueOf(pkg.Sinh),
		"Cosh":            reflect.ValueOf(pkg.Cosh),
		"Asinh":           reflect.ValueOf(pkg.Asinh),
		"Float32bits":     reflect.ValueOf(pkg.Float32bits),
		"Float32frombits": reflect.ValueOf(pkg.Float32frombits),
		"Float64bits":     reflect.ValueOf(pkg.Float64bits),
		"Float64frombits": reflect.ValueOf(pkg.Float64frombits),

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
