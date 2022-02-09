package stdgolibs

import (
	pkg "testing"

	"reflect"
)

func init() {
	registerValues("testing", map[string]reflect.Value{
		// Functions
		"AllocsPerRun":  reflect.ValueOf(pkg.AllocsPerRun),
		"RunBenchmarks": reflect.ValueOf(pkg.RunBenchmarks),
		"Benchmark":     reflect.ValueOf(pkg.Benchmark),
		"Coverage":      reflect.ValueOf(pkg.Coverage),
		"RegisterCover": reflect.ValueOf(pkg.RegisterCover),
		"RunExamples":   reflect.ValueOf(pkg.RunExamples),
		"Init":          reflect.ValueOf(pkg.Init),
		"Short":         reflect.ValueOf(pkg.Short),
		"CoverMode":     reflect.ValueOf(pkg.CoverMode),
		"Verbose":       reflect.ValueOf(pkg.Verbose),
		"Main":          reflect.ValueOf(pkg.Main),
		"MainStart":     reflect.ValueOf(pkg.MainStart),
		"RunTests":      reflect.ValueOf(pkg.RunTests),

		// Consts

		// Variables

	})
	registerTypes("testing", map[string]reflect.Type{
		// Non interfaces

		"InternalBenchmark": reflect.TypeOf((*pkg.InternalBenchmark)(nil)).Elem(),
		"B":                 reflect.TypeOf((*pkg.B)(nil)).Elem(),
		"BenchmarkResult":   reflect.TypeOf((*pkg.BenchmarkResult)(nil)).Elem(),
		"PB":                reflect.TypeOf((*pkg.PB)(nil)).Elem(),
		"CoverBlock":        reflect.TypeOf((*pkg.CoverBlock)(nil)).Elem(),
		"Cover":             reflect.TypeOf((*pkg.Cover)(nil)).Elem(),
		"InternalExample":   reflect.TypeOf((*pkg.InternalExample)(nil)).Elem(),
		"T":                 reflect.TypeOf((*pkg.T)(nil)).Elem(),
		"InternalTest":      reflect.TypeOf((*pkg.InternalTest)(nil)).Elem(),
		"M":                 reflect.TypeOf((*pkg.M)(nil)).Elem(),
	})
}
