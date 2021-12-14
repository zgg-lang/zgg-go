package stdgolibs

import (
	pkg "database/sql/driver"

	"reflect"
)

func init() {
	registerValues("database/sql/driver", map[string]reflect.Value{
		// Functions
		"IsValue":     reflect.ValueOf(pkg.IsValue),
		"IsScanValue": reflect.ValueOf(pkg.IsScanValue),

		// Consts

		// Variables

		"ErrSkip":                   reflect.ValueOf(&pkg.ErrSkip),
		"ErrBadConn":                reflect.ValueOf(&pkg.ErrBadConn),
		"ErrRemoveArgument":         reflect.ValueOf(&pkg.ErrRemoveArgument),
		"ResultNoRows":              reflect.ValueOf(&pkg.ResultNoRows),
		"Bool":                      reflect.ValueOf(&pkg.Bool),
		"Int32":                     reflect.ValueOf(&pkg.Int32),
		"String":                    reflect.ValueOf(&pkg.String),
		"DefaultParameterConverter": reflect.ValueOf(&pkg.DefaultParameterConverter),
	})
	registerTypes("database/sql/driver", map[string]reflect.Type{
		// Non interfaces

		"NamedValue":     reflect.TypeOf((*pkg.NamedValue)(nil)).Elem(),
		"IsolationLevel": reflect.TypeOf((*pkg.IsolationLevel)(nil)).Elem(),
		"TxOptions":      reflect.TypeOf((*pkg.TxOptions)(nil)).Elem(),
		"RowsAffected":   reflect.TypeOf((*pkg.RowsAffected)(nil)).Elem(),
		"Null":           reflect.TypeOf((*pkg.Null)(nil)).Elem(),
		"NotNull":        reflect.TypeOf((*pkg.NotNull)(nil)).Elem(),
	})
}
