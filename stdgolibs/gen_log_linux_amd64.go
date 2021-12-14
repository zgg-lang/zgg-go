package stdgolibs

import (
	pkg "log"

	"reflect"
)

func init() {
	registerValues("log", map[string]reflect.Value{
		// Functions
		"New":       reflect.ValueOf(pkg.New),
		"Default":   reflect.ValueOf(pkg.Default),
		"SetOutput": reflect.ValueOf(pkg.SetOutput),
		"Flags":     reflect.ValueOf(pkg.Flags),
		"SetFlags":  reflect.ValueOf(pkg.SetFlags),
		"Prefix":    reflect.ValueOf(pkg.Prefix),
		"SetPrefix": reflect.ValueOf(pkg.SetPrefix),
		"Writer":    reflect.ValueOf(pkg.Writer),
		"Print":     reflect.ValueOf(pkg.Print),
		"Printf":    reflect.ValueOf(pkg.Printf),
		"Println":   reflect.ValueOf(pkg.Println),
		"Fatal":     reflect.ValueOf(pkg.Fatal),
		"Fatalf":    reflect.ValueOf(pkg.Fatalf),
		"Fatalln":   reflect.ValueOf(pkg.Fatalln),
		"Panic":     reflect.ValueOf(pkg.Panic),
		"Panicf":    reflect.ValueOf(pkg.Panicf),
		"Panicln":   reflect.ValueOf(pkg.Panicln),
		"Output":    reflect.ValueOf(pkg.Output),

		// Consts

		"Ldate":         reflect.ValueOf(pkg.Ldate),
		"Ltime":         reflect.ValueOf(pkg.Ltime),
		"Lmicroseconds": reflect.ValueOf(pkg.Lmicroseconds),
		"Llongfile":     reflect.ValueOf(pkg.Llongfile),
		"Lshortfile":    reflect.ValueOf(pkg.Lshortfile),
		"LUTC":          reflect.ValueOf(pkg.LUTC),
		"Lmsgprefix":    reflect.ValueOf(pkg.Lmsgprefix),
		"LstdFlags":     reflect.ValueOf(pkg.LstdFlags),

		// Variables

	})
	registerTypes("log", map[string]reflect.Type{
		// Non interfaces

		"Logger": reflect.TypeOf((*pkg.Logger)(nil)).Elem(),
	})
}
