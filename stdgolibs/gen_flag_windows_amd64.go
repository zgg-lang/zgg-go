package stdgolibs

import (
	pkg "flag"

	"reflect"
)

func init() {
	registerValues("flag", map[string]reflect.Value{
		// Functions
		"VisitAll":      reflect.ValueOf(pkg.VisitAll),
		"Visit":         reflect.ValueOf(pkg.Visit),
		"Lookup":        reflect.ValueOf(pkg.Lookup),
		"Set":           reflect.ValueOf(pkg.Set),
		"UnquoteUsage":  reflect.ValueOf(pkg.UnquoteUsage),
		"PrintDefaults": reflect.ValueOf(pkg.PrintDefaults),
		"NFlag":         reflect.ValueOf(pkg.NFlag),
		"Arg":           reflect.ValueOf(pkg.Arg),
		"NArg":          reflect.ValueOf(pkg.NArg),
		"Args":          reflect.ValueOf(pkg.Args),
		"BoolVar":       reflect.ValueOf(pkg.BoolVar),
		"Bool":          reflect.ValueOf(pkg.Bool),
		"IntVar":        reflect.ValueOf(pkg.IntVar),
		"Int":           reflect.ValueOf(pkg.Int),
		"Int64Var":      reflect.ValueOf(pkg.Int64Var),
		"Int64":         reflect.ValueOf(pkg.Int64),
		"UintVar":       reflect.ValueOf(pkg.UintVar),
		"Uint":          reflect.ValueOf(pkg.Uint),
		"Uint64Var":     reflect.ValueOf(pkg.Uint64Var),
		"Uint64":        reflect.ValueOf(pkg.Uint64),
		"StringVar":     reflect.ValueOf(pkg.StringVar),
		"String":        reflect.ValueOf(pkg.String),
		"Float64Var":    reflect.ValueOf(pkg.Float64Var),
		"Float64":       reflect.ValueOf(pkg.Float64),
		"DurationVar":   reflect.ValueOf(pkg.DurationVar),
		"Duration":      reflect.ValueOf(pkg.Duration),
		"Func":          reflect.ValueOf(pkg.Func),
		"Var":           reflect.ValueOf(pkg.Var),
		"Parse":         reflect.ValueOf(pkg.Parse),
		"Parsed":        reflect.ValueOf(pkg.Parsed),
		"NewFlagSet":    reflect.ValueOf(pkg.NewFlagSet),

		// Consts

		"ContinueOnError": reflect.ValueOf(pkg.ContinueOnError),
		"ExitOnError":     reflect.ValueOf(pkg.ExitOnError),
		"PanicOnError":    reflect.ValueOf(pkg.PanicOnError),

		// Variables

		"ErrHelp":     reflect.ValueOf(&pkg.ErrHelp),
		"Usage":       reflect.ValueOf(&pkg.Usage),
		"CommandLine": reflect.ValueOf(&pkg.CommandLine),
	})
	registerTypes("flag", map[string]reflect.Type{
		// Non interfaces

		"ErrorHandling": reflect.TypeOf((*pkg.ErrorHandling)(nil)).Elem(),
		"FlagSet":       reflect.TypeOf((*pkg.FlagSet)(nil)).Elem(),
		"Flag":          reflect.TypeOf((*pkg.Flag)(nil)).Elem(),
	})
}
