package stdgolibs

import (
	pkg "strconv"

	"reflect"
)

func init() {
	registerValues("strconv", map[string]reflect.Value{
		// Functions
		"FormatUint":               reflect.ValueOf(pkg.FormatUint),
		"FormatInt":                reflect.ValueOf(pkg.FormatInt),
		"Itoa":                     reflect.ValueOf(pkg.Itoa),
		"AppendInt":                reflect.ValueOf(pkg.AppendInt),
		"AppendUint":               reflect.ValueOf(pkg.AppendUint),
		"Quote":                    reflect.ValueOf(pkg.Quote),
		"AppendQuote":              reflect.ValueOf(pkg.AppendQuote),
		"QuoteToASCII":             reflect.ValueOf(pkg.QuoteToASCII),
		"AppendQuoteToASCII":       reflect.ValueOf(pkg.AppendQuoteToASCII),
		"QuoteToGraphic":           reflect.ValueOf(pkg.QuoteToGraphic),
		"AppendQuoteToGraphic":     reflect.ValueOf(pkg.AppendQuoteToGraphic),
		"QuoteRune":                reflect.ValueOf(pkg.QuoteRune),
		"AppendQuoteRune":          reflect.ValueOf(pkg.AppendQuoteRune),
		"QuoteRuneToASCII":         reflect.ValueOf(pkg.QuoteRuneToASCII),
		"AppendQuoteRuneToASCII":   reflect.ValueOf(pkg.AppendQuoteRuneToASCII),
		"QuoteRuneToGraphic":       reflect.ValueOf(pkg.QuoteRuneToGraphic),
		"AppendQuoteRuneToGraphic": reflect.ValueOf(pkg.AppendQuoteRuneToGraphic),
		"CanBackquote":             reflect.ValueOf(pkg.CanBackquote),
		"UnquoteChar":              reflect.ValueOf(pkg.UnquoteChar),
		"Unquote":                  reflect.ValueOf(pkg.Unquote),
		"IsPrint":                  reflect.ValueOf(pkg.IsPrint),
		"IsGraphic":                reflect.ValueOf(pkg.IsGraphic),
		"ParseBool":                reflect.ValueOf(pkg.ParseBool),
		"FormatBool":               reflect.ValueOf(pkg.FormatBool),
		"AppendBool":               reflect.ValueOf(pkg.AppendBool),
		"ParseComplex":             reflect.ValueOf(pkg.ParseComplex),
		"FormatComplex":            reflect.ValueOf(pkg.FormatComplex),
		"FormatFloat":              reflect.ValueOf(pkg.FormatFloat),
		"AppendFloat":              reflect.ValueOf(pkg.AppendFloat),
		"ParseFloat":               reflect.ValueOf(pkg.ParseFloat),
		"ParseUint":                reflect.ValueOf(pkg.ParseUint),
		"ParseInt":                 reflect.ValueOf(pkg.ParseInt),
		"Atoi":                     reflect.ValueOf(pkg.Atoi),

		// Consts

		"IntSize": reflect.ValueOf(pkg.IntSize),

		// Variables

		"ErrRange":  reflect.ValueOf(&pkg.ErrRange),
		"ErrSyntax": reflect.ValueOf(&pkg.ErrSyntax),
	})
	registerTypes("strconv", map[string]reflect.Type{
		// Non interfaces

		"NumError": reflect.TypeOf((*pkg.NumError)(nil)).Elem(),
	})
}
