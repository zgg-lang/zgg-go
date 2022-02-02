package stdgolibs

import (
	pkg "go/constant"

	"reflect"
)

func init() {
	registerValues("go/constant", map[string]reflect.Value{
		// Functions
		"MakeUnknown":     reflect.ValueOf(pkg.MakeUnknown),
		"MakeBool":        reflect.ValueOf(pkg.MakeBool),
		"MakeString":      reflect.ValueOf(pkg.MakeString),
		"MakeInt64":       reflect.ValueOf(pkg.MakeInt64),
		"MakeUint64":      reflect.ValueOf(pkg.MakeUint64),
		"MakeFloat64":     reflect.ValueOf(pkg.MakeFloat64),
		"MakeFromLiteral": reflect.ValueOf(pkg.MakeFromLiteral),
		"BoolVal":         reflect.ValueOf(pkg.BoolVal),
		"StringVal":       reflect.ValueOf(pkg.StringVal),
		"Int64Val":        reflect.ValueOf(pkg.Int64Val),
		"Uint64Val":       reflect.ValueOf(pkg.Uint64Val),
		"Float32Val":      reflect.ValueOf(pkg.Float32Val),
		"Float64Val":      reflect.ValueOf(pkg.Float64Val),
		"Val":             reflect.ValueOf(pkg.Val),
		"Make":            reflect.ValueOf(pkg.Make),
		"BitLen":          reflect.ValueOf(pkg.BitLen),
		"Sign":            reflect.ValueOf(pkg.Sign),
		"Bytes":           reflect.ValueOf(pkg.Bytes),
		"MakeFromBytes":   reflect.ValueOf(pkg.MakeFromBytes),
		"Num":             reflect.ValueOf(pkg.Num),
		"Denom":           reflect.ValueOf(pkg.Denom),
		"MakeImag":        reflect.ValueOf(pkg.MakeImag),
		"Real":            reflect.ValueOf(pkg.Real),
		"Imag":            reflect.ValueOf(pkg.Imag),
		"ToInt":           reflect.ValueOf(pkg.ToInt),
		"ToFloat":         reflect.ValueOf(pkg.ToFloat),
		"ToComplex":       reflect.ValueOf(pkg.ToComplex),
		"UnaryOp":         reflect.ValueOf(pkg.UnaryOp),
		"BinaryOp":        reflect.ValueOf(pkg.BinaryOp),
		"Shift":           reflect.ValueOf(pkg.Shift),
		"Compare":         reflect.ValueOf(pkg.Compare),

		// Consts

		"Unknown": reflect.ValueOf(pkg.Unknown),
		"Bool":    reflect.ValueOf(pkg.Bool),
		"String":  reflect.ValueOf(pkg.String),
		"Int":     reflect.ValueOf(pkg.Int),
		"Float":   reflect.ValueOf(pkg.Float),
		"Complex": reflect.ValueOf(pkg.Complex),

		// Variables

	})
	registerTypes("go/constant", map[string]reflect.Type{
		// Non interfaces

		"Kind": reflect.TypeOf((*pkg.Kind)(nil)).Elem(),
	})
}
