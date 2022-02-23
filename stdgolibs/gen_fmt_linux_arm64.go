package stdgolibs

import (
	pkg "fmt"

	"reflect"
)

func init() {
	registerValues("fmt", map[string]reflect.Value{
		// Functions
		"Errorf":   reflect.ValueOf(pkg.Errorf),
		"Fprintf":  reflect.ValueOf(pkg.Fprintf),
		"Printf":   reflect.ValueOf(pkg.Printf),
		"Sprintf":  reflect.ValueOf(pkg.Sprintf),
		"Fprint":   reflect.ValueOf(pkg.Fprint),
		"Print":    reflect.ValueOf(pkg.Print),
		"Sprint":   reflect.ValueOf(pkg.Sprint),
		"Fprintln": reflect.ValueOf(pkg.Fprintln),
		"Println":  reflect.ValueOf(pkg.Println),
		"Sprintln": reflect.ValueOf(pkg.Sprintln),
		"Scan":     reflect.ValueOf(pkg.Scan),
		"Scanln":   reflect.ValueOf(pkg.Scanln),
		"Scanf":    reflect.ValueOf(pkg.Scanf),
		"Sscan":    reflect.ValueOf(pkg.Sscan),
		"Sscanln":  reflect.ValueOf(pkg.Sscanln),
		"Sscanf":   reflect.ValueOf(pkg.Sscanf),
		"Fscan":    reflect.ValueOf(pkg.Fscan),
		"Fscanln":  reflect.ValueOf(pkg.Fscanln),
		"Fscanf":   reflect.ValueOf(pkg.Fscanf),

		// Consts

		// Variables

	})
	registerTypes("fmt", map[string]reflect.Type{
		// Non interfaces

	})
}
