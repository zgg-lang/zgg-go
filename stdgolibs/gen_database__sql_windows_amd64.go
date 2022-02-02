package stdgolibs

import (
	pkg "database/sql"

	"reflect"
)

func init() {
	registerValues("database/sql", map[string]reflect.Value{
		// Functions
		"Register": reflect.ValueOf(pkg.Register),
		"Drivers":  reflect.ValueOf(pkg.Drivers),
		"Named":    reflect.ValueOf(pkg.Named),
		"OpenDB":   reflect.ValueOf(pkg.OpenDB),
		"Open":     reflect.ValueOf(pkg.Open),

		// Consts

		"LevelDefault":         reflect.ValueOf(pkg.LevelDefault),
		"LevelReadUncommitted": reflect.ValueOf(pkg.LevelReadUncommitted),
		"LevelReadCommitted":   reflect.ValueOf(pkg.LevelReadCommitted),
		"LevelWriteCommitted":  reflect.ValueOf(pkg.LevelWriteCommitted),
		"LevelRepeatableRead":  reflect.ValueOf(pkg.LevelRepeatableRead),
		"LevelSnapshot":        reflect.ValueOf(pkg.LevelSnapshot),
		"LevelSerializable":    reflect.ValueOf(pkg.LevelSerializable),
		"LevelLinearizable":    reflect.ValueOf(pkg.LevelLinearizable),

		// Variables

		"ErrNoRows":   reflect.ValueOf(&pkg.ErrNoRows),
		"ErrConnDone": reflect.ValueOf(&pkg.ErrConnDone),
		"ErrTxDone":   reflect.ValueOf(&pkg.ErrTxDone),
	})
	registerTypes("database/sql", map[string]reflect.Type{
		// Non interfaces

		"NamedArg":       reflect.TypeOf((*pkg.NamedArg)(nil)).Elem(),
		"IsolationLevel": reflect.TypeOf((*pkg.IsolationLevel)(nil)).Elem(),
		"TxOptions":      reflect.TypeOf((*pkg.TxOptions)(nil)).Elem(),
		"RawBytes":       reflect.TypeOf((*pkg.RawBytes)(nil)).Elem(),
		"NullString":     reflect.TypeOf((*pkg.NullString)(nil)).Elem(),
		"NullInt64":      reflect.TypeOf((*pkg.NullInt64)(nil)).Elem(),
		"NullInt32":      reflect.TypeOf((*pkg.NullInt32)(nil)).Elem(),
		"NullFloat64":    reflect.TypeOf((*pkg.NullFloat64)(nil)).Elem(),
		"NullBool":       reflect.TypeOf((*pkg.NullBool)(nil)).Elem(),
		"NullTime":       reflect.TypeOf((*pkg.NullTime)(nil)).Elem(),
		"Out":            reflect.TypeOf((*pkg.Out)(nil)).Elem(),
		"DB":             reflect.TypeOf((*pkg.DB)(nil)).Elem(),
		"DBStats":        reflect.TypeOf((*pkg.DBStats)(nil)).Elem(),
		"Conn":           reflect.TypeOf((*pkg.Conn)(nil)).Elem(),
		"Tx":             reflect.TypeOf((*pkg.Tx)(nil)).Elem(),
		"Stmt":           reflect.TypeOf((*pkg.Stmt)(nil)).Elem(),
		"Rows":           reflect.TypeOf((*pkg.Rows)(nil)).Elem(),
		"ColumnType":     reflect.TypeOf((*pkg.ColumnType)(nil)).Elem(),
		"Row":            reflect.TypeOf((*pkg.Row)(nil)).Elem(),
	})
}
