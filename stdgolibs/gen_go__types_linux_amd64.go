package stdgolibs

import (
	pkg "go/types"

	"reflect"
)

func init() {
	registerValues("go/types", map[string]reflect.Value{
		// Functions
		"Id":                      reflect.ValueOf(pkg.Id),
		"NewPkgName":              reflect.ValueOf(pkg.NewPkgName),
		"NewConst":                reflect.ValueOf(pkg.NewConst),
		"NewTypeName":             reflect.ValueOf(pkg.NewTypeName),
		"NewVar":                  reflect.ValueOf(pkg.NewVar),
		"NewParam":                reflect.ValueOf(pkg.NewParam),
		"NewField":                reflect.ValueOf(pkg.NewField),
		"NewFunc":                 reflect.ValueOf(pkg.NewFunc),
		"NewLabel":                reflect.ValueOf(pkg.NewLabel),
		"ObjectString":            reflect.ValueOf(pkg.ObjectString),
		"RelativeTo":              reflect.ValueOf(pkg.RelativeTo),
		"TypeString":              reflect.ValueOf(pkg.TypeString),
		"WriteType":               reflect.ValueOf(pkg.WriteType),
		"WriteSignature":          reflect.ValueOf(pkg.WriteSignature),
		"IsInterface":             reflect.ValueOf(pkg.IsInterface),
		"Comparable":              reflect.ValueOf(pkg.Comparable),
		"Default":                 reflect.ValueOf(pkg.Default),
		"NewArray":                reflect.ValueOf(pkg.NewArray),
		"NewSlice":                reflect.ValueOf(pkg.NewSlice),
		"NewStruct":               reflect.ValueOf(pkg.NewStruct),
		"NewPointer":              reflect.ValueOf(pkg.NewPointer),
		"NewTuple":                reflect.ValueOf(pkg.NewTuple),
		"NewSignature":            reflect.ValueOf(pkg.NewSignature),
		"NewInterface":            reflect.ValueOf(pkg.NewInterface),
		"NewInterfaceType":        reflect.ValueOf(pkg.NewInterfaceType),
		"NewMap":                  reflect.ValueOf(pkg.NewMap),
		"NewChan":                 reflect.ValueOf(pkg.NewChan),
		"NewNamed":                reflect.ValueOf(pkg.NewNamed),
		"Eval":                    reflect.ValueOf(pkg.Eval),
		"CheckExpr":               reflect.ValueOf(pkg.CheckExpr),
		"LookupFieldOrMethod":     reflect.ValueOf(pkg.LookupFieldOrMethod),
		"MissingMethod":           reflect.ValueOf(pkg.MissingMethod),
		"ExprString":              reflect.ValueOf(pkg.ExprString),
		"WriteExpr":               reflect.ValueOf(pkg.WriteExpr),
		"DefPredeclaredTestFuncs": reflect.ValueOf(pkg.DefPredeclaredTestFuncs),
		"NewScope":                reflect.ValueOf(pkg.NewScope),
		"SelectionString":         reflect.ValueOf(pkg.SelectionString),
		"SizesFor":                reflect.ValueOf(pkg.SizesFor),
		"NewChecker":              reflect.ValueOf(pkg.NewChecker),
		"NewMethodSet":            reflect.ValueOf(pkg.NewMethodSet),
		"NewPackage":              reflect.ValueOf(pkg.NewPackage),
		"AssertableTo":            reflect.ValueOf(pkg.AssertableTo),
		"AssignableTo":            reflect.ValueOf(pkg.AssignableTo),
		"ConvertibleTo":           reflect.ValueOf(pkg.ConvertibleTo),
		"Implements":              reflect.ValueOf(pkg.Implements),
		"Identical":               reflect.ValueOf(pkg.Identical),
		"IdenticalIgnoreTags":     reflect.ValueOf(pkg.IdenticalIgnoreTags),

		// Consts

		"Invalid":        reflect.ValueOf(pkg.Invalid),
		"Bool":           reflect.ValueOf(pkg.Bool),
		"Int":            reflect.ValueOf(pkg.Int),
		"Int8":           reflect.ValueOf(pkg.Int8),
		"Int16":          reflect.ValueOf(pkg.Int16),
		"Int32":          reflect.ValueOf(pkg.Int32),
		"Int64":          reflect.ValueOf(pkg.Int64),
		"Uint":           reflect.ValueOf(pkg.Uint),
		"Uint8":          reflect.ValueOf(pkg.Uint8),
		"Uint16":         reflect.ValueOf(pkg.Uint16),
		"Uint32":         reflect.ValueOf(pkg.Uint32),
		"Uint64":         reflect.ValueOf(pkg.Uint64),
		"Uintptr":        reflect.ValueOf(pkg.Uintptr),
		"Float32":        reflect.ValueOf(pkg.Float32),
		"Float64":        reflect.ValueOf(pkg.Float64),
		"Complex64":      reflect.ValueOf(pkg.Complex64),
		"Complex128":     reflect.ValueOf(pkg.Complex128),
		"String":         reflect.ValueOf(pkg.String),
		"UnsafePointer":  reflect.ValueOf(pkg.UnsafePointer),
		"UntypedBool":    reflect.ValueOf(pkg.UntypedBool),
		"UntypedInt":     reflect.ValueOf(pkg.UntypedInt),
		"UntypedRune":    reflect.ValueOf(pkg.UntypedRune),
		"UntypedFloat":   reflect.ValueOf(pkg.UntypedFloat),
		"UntypedComplex": reflect.ValueOf(pkg.UntypedComplex),
		"UntypedString":  reflect.ValueOf(pkg.UntypedString),
		"UntypedNil":     reflect.ValueOf(pkg.UntypedNil),
		"Byte":           reflect.ValueOf(pkg.Byte),
		"Rune":           reflect.ValueOf(pkg.Rune),
		"IsBoolean":      reflect.ValueOf(pkg.IsBoolean),
		"IsInteger":      reflect.ValueOf(pkg.IsInteger),
		"IsUnsigned":     reflect.ValueOf(pkg.IsUnsigned),
		"IsFloat":        reflect.ValueOf(pkg.IsFloat),
		"IsComplex":      reflect.ValueOf(pkg.IsComplex),
		"IsString":       reflect.ValueOf(pkg.IsString),
		"IsUntyped":      reflect.ValueOf(pkg.IsUntyped),
		"IsOrdered":      reflect.ValueOf(pkg.IsOrdered),
		"IsNumeric":      reflect.ValueOf(pkg.IsNumeric),
		"IsConstType":    reflect.ValueOf(pkg.IsConstType),
		"SendRecv":       reflect.ValueOf(pkg.SendRecv),
		"SendOnly":       reflect.ValueOf(pkg.SendOnly),
		"RecvOnly":       reflect.ValueOf(pkg.RecvOnly),
		"FieldVal":       reflect.ValueOf(pkg.FieldVal),
		"MethodVal":      reflect.ValueOf(pkg.MethodVal),
		"MethodExpr":     reflect.ValueOf(pkg.MethodExpr),

		// Variables

		"Universe": reflect.ValueOf(&pkg.Universe),
		"Unsafe":   reflect.ValueOf(&pkg.Unsafe),
		"Typ":      reflect.ValueOf(&pkg.Typ),
	})
	registerTypes("go/types", map[string]reflect.Type{
		// Non interfaces

		"PkgName":       reflect.TypeOf((*pkg.PkgName)(nil)).Elem(),
		"Const":         reflect.TypeOf((*pkg.Const)(nil)).Elem(),
		"TypeName":      reflect.TypeOf((*pkg.TypeName)(nil)).Elem(),
		"Var":           reflect.TypeOf((*pkg.Var)(nil)).Elem(),
		"Func":          reflect.TypeOf((*pkg.Func)(nil)).Elem(),
		"Label":         reflect.TypeOf((*pkg.Label)(nil)).Elem(),
		"Builtin":       reflect.TypeOf((*pkg.Builtin)(nil)).Elem(),
		"Nil":           reflect.TypeOf((*pkg.Nil)(nil)).Elem(),
		"Qualifier":     reflect.TypeOf((*pkg.Qualifier)(nil)).Elem(),
		"BasicKind":     reflect.TypeOf((*pkg.BasicKind)(nil)).Elem(),
		"BasicInfo":     reflect.TypeOf((*pkg.BasicInfo)(nil)).Elem(),
		"Basic":         reflect.TypeOf((*pkg.Basic)(nil)).Elem(),
		"Array":         reflect.TypeOf((*pkg.Array)(nil)).Elem(),
		"Slice":         reflect.TypeOf((*pkg.Slice)(nil)).Elem(),
		"Struct":        reflect.TypeOf((*pkg.Struct)(nil)).Elem(),
		"Pointer":       reflect.TypeOf((*pkg.Pointer)(nil)).Elem(),
		"Tuple":         reflect.TypeOf((*pkg.Tuple)(nil)).Elem(),
		"Signature":     reflect.TypeOf((*pkg.Signature)(nil)).Elem(),
		"Interface":     reflect.TypeOf((*pkg.Interface)(nil)).Elem(),
		"Map":           reflect.TypeOf((*pkg.Map)(nil)).Elem(),
		"Chan":          reflect.TypeOf((*pkg.Chan)(nil)).Elem(),
		"ChanDir":       reflect.TypeOf((*pkg.ChanDir)(nil)).Elem(),
		"Named":         reflect.TypeOf((*pkg.Named)(nil)).Elem(),
		"Scope":         reflect.TypeOf((*pkg.Scope)(nil)).Elem(),
		"SelectionKind": reflect.TypeOf((*pkg.SelectionKind)(nil)).Elem(),
		"Selection":     reflect.TypeOf((*pkg.Selection)(nil)).Elem(),
		"StdSizes":      reflect.TypeOf((*pkg.StdSizes)(nil)).Elem(),
		"Checker":       reflect.TypeOf((*pkg.Checker)(nil)).Elem(),
		"MethodSet":     reflect.TypeOf((*pkg.MethodSet)(nil)).Elem(),
		"Package":       reflect.TypeOf((*pkg.Package)(nil)).Elem(),
		"Error":         reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"ImportMode":    reflect.TypeOf((*pkg.ImportMode)(nil)).Elem(),
		"Config":        reflect.TypeOf((*pkg.Config)(nil)).Elem(),
		"Info":          reflect.TypeOf((*pkg.Info)(nil)).Elem(),
		"TypeAndValue":  reflect.TypeOf((*pkg.TypeAndValue)(nil)).Elem(),
		"Initializer":   reflect.TypeOf((*pkg.Initializer)(nil)).Elem(),
	})
}
