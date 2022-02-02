package stdgolibs

import (
	pkg "archive/tar"

	"reflect"
)

func init() {
	registerValues("archive/tar", map[string]reflect.Value{
		// Functions
		"NewReader":      reflect.ValueOf(pkg.NewReader),
		"NewWriter":      reflect.ValueOf(pkg.NewWriter),
		"FileInfoHeader": reflect.ValueOf(pkg.FileInfoHeader),

		// Consts

		"TypeReg":           reflect.ValueOf(pkg.TypeReg),
		"TypeRegA":          reflect.ValueOf(pkg.TypeRegA),
		"TypeLink":          reflect.ValueOf(pkg.TypeLink),
		"TypeSymlink":       reflect.ValueOf(pkg.TypeSymlink),
		"TypeChar":          reflect.ValueOf(pkg.TypeChar),
		"TypeBlock":         reflect.ValueOf(pkg.TypeBlock),
		"TypeDir":           reflect.ValueOf(pkg.TypeDir),
		"TypeFifo":          reflect.ValueOf(pkg.TypeFifo),
		"TypeCont":          reflect.ValueOf(pkg.TypeCont),
		"TypeXHeader":       reflect.ValueOf(pkg.TypeXHeader),
		"TypeXGlobalHeader": reflect.ValueOf(pkg.TypeXGlobalHeader),
		"TypeGNUSparse":     reflect.ValueOf(pkg.TypeGNUSparse),
		"TypeGNULongName":   reflect.ValueOf(pkg.TypeGNULongName),
		"TypeGNULongLink":   reflect.ValueOf(pkg.TypeGNULongLink),
		"FormatUnknown":     reflect.ValueOf(pkg.FormatUnknown),
		"FormatUSTAR":       reflect.ValueOf(pkg.FormatUSTAR),
		"FormatPAX":         reflect.ValueOf(pkg.FormatPAX),
		"FormatGNU":         reflect.ValueOf(pkg.FormatGNU),

		// Variables

		"ErrHeader":          reflect.ValueOf(&pkg.ErrHeader),
		"ErrWriteTooLong":    reflect.ValueOf(&pkg.ErrWriteTooLong),
		"ErrFieldTooLong":    reflect.ValueOf(&pkg.ErrFieldTooLong),
		"ErrWriteAfterClose": reflect.ValueOf(&pkg.ErrWriteAfterClose),
	})
	registerTypes("archive/tar", map[string]reflect.Type{
		// Non interfaces

		"Reader": reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"Writer": reflect.TypeOf((*pkg.Writer)(nil)).Elem(),
		"Header": reflect.TypeOf((*pkg.Header)(nil)).Elem(),
		"Format": reflect.TypeOf((*pkg.Format)(nil)).Elem(),
	})
}
