package stdgolibs

import (
	pkg "bufio"

	"reflect"
)

func init() {
	registerValues("bufio", map[string]reflect.Value{
		// Functions
		"NewReaderSize": reflect.ValueOf(pkg.NewReaderSize),
		"NewReader":     reflect.ValueOf(pkg.NewReader),
		"NewWriterSize": reflect.ValueOf(pkg.NewWriterSize),
		"NewWriter":     reflect.ValueOf(pkg.NewWriter),
		"NewReadWriter": reflect.ValueOf(pkg.NewReadWriter),
		"NewScanner":    reflect.ValueOf(pkg.NewScanner),
		"ScanBytes":     reflect.ValueOf(pkg.ScanBytes),
		"ScanRunes":     reflect.ValueOf(pkg.ScanRunes),
		"ScanLines":     reflect.ValueOf(pkg.ScanLines),
		"ScanWords":     reflect.ValueOf(pkg.ScanWords),

		// Consts

		"MaxScanTokenSize": reflect.ValueOf(pkg.MaxScanTokenSize),

		// Variables

		"ErrInvalidUnreadByte": reflect.ValueOf(&pkg.ErrInvalidUnreadByte),
		"ErrInvalidUnreadRune": reflect.ValueOf(&pkg.ErrInvalidUnreadRune),
		"ErrBufferFull":        reflect.ValueOf(&pkg.ErrBufferFull),
		"ErrNegativeCount":     reflect.ValueOf(&pkg.ErrNegativeCount),
		"ErrTooLong":           reflect.ValueOf(&pkg.ErrTooLong),
		"ErrNegativeAdvance":   reflect.ValueOf(&pkg.ErrNegativeAdvance),
		"ErrAdvanceTooFar":     reflect.ValueOf(&pkg.ErrAdvanceTooFar),
		"ErrBadReadCount":      reflect.ValueOf(&pkg.ErrBadReadCount),
		"ErrFinalToken":        reflect.ValueOf(&pkg.ErrFinalToken),
	})
	registerTypes("bufio", map[string]reflect.Type{
		// Non interfaces

		"Reader":     reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer":     reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"ReadWriter": reflect.TypeOf((*pkg.ReadWriter)(nil)).Elem(),
		"Scanner":    reflect.TypeOf((*pkg.Scanner)(nil)).Elem(),
		"SplitFunc":  reflect.TypeOf((*pkg.SplitFunc)(nil)).Elem(),
	})
}
