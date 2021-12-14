package stdgolibs

import (
	pkg "reflect"

	"reflect"
)

func init() {
	registerValues("reflect", map[string]reflect.Value{
		// Functions
		"DeepEqual":       reflect.ValueOf(pkg.DeepEqual),
		"MakeFunc":        reflect.ValueOf(pkg.MakeFunc),
		"Swapper":         reflect.ValueOf(pkg.Swapper),
		"TypeOf":          reflect.ValueOf(pkg.TypeOf),
		"PtrTo":           reflect.ValueOf(pkg.PtrTo),
		"ChanOf":          reflect.ValueOf(pkg.ChanOf),
		"MapOf":           reflect.ValueOf(pkg.MapOf),
		"FuncOf":          reflect.ValueOf(pkg.FuncOf),
		"SliceOf":         reflect.ValueOf(pkg.SliceOf),
		"StructOf":        reflect.ValueOf(pkg.StructOf),
		"ArrayOf":         reflect.ValueOf(pkg.ArrayOf),
		"Append":          reflect.ValueOf(pkg.Append),
		"AppendSlice":     reflect.ValueOf(pkg.AppendSlice),
		"Copy":            reflect.ValueOf(pkg.Copy),
		"Select":          reflect.ValueOf(pkg.Select),
		"MakeSlice":       reflect.ValueOf(pkg.MakeSlice),
		"MakeChan":        reflect.ValueOf(pkg.MakeChan),
		"MakeMap":         reflect.ValueOf(pkg.MakeMap),
		"MakeMapWithSize": reflect.ValueOf(pkg.MakeMapWithSize),
		"Indirect":        reflect.ValueOf(pkg.Indirect),
		"ValueOf":         reflect.ValueOf(pkg.ValueOf),
		"Zero":            reflect.ValueOf(pkg.Zero),
		"New":             reflect.ValueOf(pkg.New),
		"NewAt":           reflect.ValueOf(pkg.NewAt),

		// Consts

		"Invalid":       reflect.ValueOf(pkg.Invalid),
		"Bool":          reflect.ValueOf(pkg.Bool),
		"Int":           reflect.ValueOf(pkg.Int),
		"Int8":          reflect.ValueOf(pkg.Int8),
		"Int16":         reflect.ValueOf(pkg.Int16),
		"Int32":         reflect.ValueOf(pkg.Int32),
		"Int64":         reflect.ValueOf(pkg.Int64),
		"Uint":          reflect.ValueOf(pkg.Uint),
		"Uint8":         reflect.ValueOf(pkg.Uint8),
		"Uint16":        reflect.ValueOf(pkg.Uint16),
		"Uint32":        reflect.ValueOf(pkg.Uint32),
		"Uint64":        reflect.ValueOf(pkg.Uint64),
		"Uintptr":       reflect.ValueOf(pkg.Uintptr),
		"Float32":       reflect.ValueOf(pkg.Float32),
		"Float64":       reflect.ValueOf(pkg.Float64),
		"Complex64":     reflect.ValueOf(pkg.Complex64),
		"Complex128":    reflect.ValueOf(pkg.Complex128),
		"Array":         reflect.ValueOf(pkg.Array),
		"Chan":          reflect.ValueOf(pkg.Chan),
		"Func":          reflect.ValueOf(pkg.Func),
		"Interface":     reflect.ValueOf(pkg.Interface),
		"Map":           reflect.ValueOf(pkg.Map),
		"Ptr":           reflect.ValueOf(pkg.Ptr),
		"Slice":         reflect.ValueOf(pkg.Slice),
		"String":        reflect.ValueOf(pkg.String),
		"Struct":        reflect.ValueOf(pkg.Struct),
		"UnsafePointer": reflect.ValueOf(pkg.UnsafePointer),
		"RecvDir":       reflect.ValueOf(pkg.RecvDir),
		"SendDir":       reflect.ValueOf(pkg.SendDir),
		"BothDir":       reflect.ValueOf(pkg.BothDir),
		"SelectSend":    reflect.ValueOf(pkg.SelectSend),
		"SelectRecv":    reflect.ValueOf(pkg.SelectRecv),
		"SelectDefault": reflect.ValueOf(pkg.SelectDefault),

		// Variables

	})
	registerTypes("reflect", map[string]reflect.Type{
		// Non interfaces

		"Kind":         reflect.TypeOf((*pkg.Kind)(nil)).Elem(),
		"ChanDir":      reflect.TypeOf((*pkg.ChanDir)(nil)).Elem(),
		"Method":       reflect.TypeOf((*pkg.Method)(nil)).Elem(),
		"StructField":  reflect.TypeOf((*pkg.StructField)(nil)).Elem(),
		"StructTag":    reflect.TypeOf((*pkg.StructTag)(nil)).Elem(),
		"Value":        reflect.TypeOf((*pkg.Value)(nil)).Elem(),
		"ValueError":   reflect.TypeOf((*pkg.ValueError)(nil)).Elem(),
		"MapIter":      reflect.TypeOf((*pkg.MapIter)(nil)).Elem(),
		"StringHeader": reflect.TypeOf((*pkg.StringHeader)(nil)).Elem(),
		"SliceHeader":  reflect.TypeOf((*pkg.SliceHeader)(nil)).Elem(),
		"SelectDir":    reflect.TypeOf((*pkg.SelectDir)(nil)).Elem(),
		"SelectCase":   reflect.TypeOf((*pkg.SelectCase)(nil)).Elem(),
	})
}
