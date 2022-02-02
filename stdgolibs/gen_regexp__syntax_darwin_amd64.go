package stdgolibs

import (
	pkg "regexp/syntax"

	"reflect"
)

func init() {
	registerValues("regexp/syntax", map[string]reflect.Value{
		// Functions
		"EmptyOpContext": reflect.ValueOf(pkg.EmptyOpContext),
		"IsWordChar":     reflect.ValueOf(pkg.IsWordChar),
		"Compile":        reflect.ValueOf(pkg.Compile),
		"Parse":          reflect.ValueOf(pkg.Parse),

		// Consts

		"InstAlt":                  reflect.ValueOf(pkg.InstAlt),
		"InstAltMatch":             reflect.ValueOf(pkg.InstAltMatch),
		"InstCapture":              reflect.ValueOf(pkg.InstCapture),
		"InstEmptyWidth":           reflect.ValueOf(pkg.InstEmptyWidth),
		"InstMatch":                reflect.ValueOf(pkg.InstMatch),
		"InstFail":                 reflect.ValueOf(pkg.InstFail),
		"InstNop":                  reflect.ValueOf(pkg.InstNop),
		"InstRune":                 reflect.ValueOf(pkg.InstRune),
		"InstRune1":                reflect.ValueOf(pkg.InstRune1),
		"InstRuneAny":              reflect.ValueOf(pkg.InstRuneAny),
		"InstRuneAnyNotNL":         reflect.ValueOf(pkg.InstRuneAnyNotNL),
		"EmptyBeginLine":           reflect.ValueOf(pkg.EmptyBeginLine),
		"EmptyEndLine":             reflect.ValueOf(pkg.EmptyEndLine),
		"EmptyBeginText":           reflect.ValueOf(pkg.EmptyBeginText),
		"EmptyEndText":             reflect.ValueOf(pkg.EmptyEndText),
		"EmptyWordBoundary":        reflect.ValueOf(pkg.EmptyWordBoundary),
		"EmptyNoWordBoundary":      reflect.ValueOf(pkg.EmptyNoWordBoundary),
		"OpNoMatch":                reflect.ValueOf(pkg.OpNoMatch),
		"OpEmptyMatch":             reflect.ValueOf(pkg.OpEmptyMatch),
		"OpLiteral":                reflect.ValueOf(pkg.OpLiteral),
		"OpCharClass":              reflect.ValueOf(pkg.OpCharClass),
		"OpAnyCharNotNL":           reflect.ValueOf(pkg.OpAnyCharNotNL),
		"OpAnyChar":                reflect.ValueOf(pkg.OpAnyChar),
		"OpBeginLine":              reflect.ValueOf(pkg.OpBeginLine),
		"OpEndLine":                reflect.ValueOf(pkg.OpEndLine),
		"OpBeginText":              reflect.ValueOf(pkg.OpBeginText),
		"OpEndText":                reflect.ValueOf(pkg.OpEndText),
		"OpWordBoundary":           reflect.ValueOf(pkg.OpWordBoundary),
		"OpNoWordBoundary":         reflect.ValueOf(pkg.OpNoWordBoundary),
		"OpCapture":                reflect.ValueOf(pkg.OpCapture),
		"OpStar":                   reflect.ValueOf(pkg.OpStar),
		"OpPlus":                   reflect.ValueOf(pkg.OpPlus),
		"OpQuest":                  reflect.ValueOf(pkg.OpQuest),
		"OpRepeat":                 reflect.ValueOf(pkg.OpRepeat),
		"OpConcat":                 reflect.ValueOf(pkg.OpConcat),
		"OpAlternate":              reflect.ValueOf(pkg.OpAlternate),
		"ErrInternalError":         reflect.ValueOf(pkg.ErrInternalError),
		"ErrInvalidCharClass":      reflect.ValueOf(pkg.ErrInvalidCharClass),
		"ErrInvalidCharRange":      reflect.ValueOf(pkg.ErrInvalidCharRange),
		"ErrInvalidEscape":         reflect.ValueOf(pkg.ErrInvalidEscape),
		"ErrInvalidNamedCapture":   reflect.ValueOf(pkg.ErrInvalidNamedCapture),
		"ErrInvalidPerlOp":         reflect.ValueOf(pkg.ErrInvalidPerlOp),
		"ErrInvalidRepeatOp":       reflect.ValueOf(pkg.ErrInvalidRepeatOp),
		"ErrInvalidRepeatSize":     reflect.ValueOf(pkg.ErrInvalidRepeatSize),
		"ErrInvalidUTF8":           reflect.ValueOf(pkg.ErrInvalidUTF8),
		"ErrMissingBracket":        reflect.ValueOf(pkg.ErrMissingBracket),
		"ErrMissingParen":          reflect.ValueOf(pkg.ErrMissingParen),
		"ErrMissingRepeatArgument": reflect.ValueOf(pkg.ErrMissingRepeatArgument),
		"ErrTrailingBackslash":     reflect.ValueOf(pkg.ErrTrailingBackslash),
		"ErrUnexpectedParen":       reflect.ValueOf(pkg.ErrUnexpectedParen),
		"FoldCase":                 reflect.ValueOf(pkg.FoldCase),
		"Literal":                  reflect.ValueOf(pkg.Literal),
		"ClassNL":                  reflect.ValueOf(pkg.ClassNL),
		"DotNL":                    reflect.ValueOf(pkg.DotNL),
		"OneLine":                  reflect.ValueOf(pkg.OneLine),
		"NonGreedy":                reflect.ValueOf(pkg.NonGreedy),
		"PerlX":                    reflect.ValueOf(pkg.PerlX),
		"UnicodeGroups":            reflect.ValueOf(pkg.UnicodeGroups),
		"WasDollar":                reflect.ValueOf(pkg.WasDollar),
		"Simple":                   reflect.ValueOf(pkg.Simple),
		"MatchNL":                  reflect.ValueOf(pkg.MatchNL),
		"Perl":                     reflect.ValueOf(pkg.Perl),
		"POSIX":                    reflect.ValueOf(pkg.POSIX),

		// Variables

	})
	registerTypes("regexp/syntax", map[string]reflect.Type{
		// Non interfaces

		"Prog":      reflect.TypeOf((*pkg.Prog)(nil)).Elem(),
		"InstOp":    reflect.TypeOf((*pkg.InstOp)(nil)).Elem(),
		"EmptyOp":   reflect.TypeOf((*pkg.EmptyOp)(nil)).Elem(),
		"Inst":      reflect.TypeOf((*pkg.Inst)(nil)).Elem(),
		"Regexp":    reflect.TypeOf((*pkg.Regexp)(nil)).Elem(),
		"Op":        reflect.TypeOf((*pkg.Op)(nil)).Elem(),
		"Error":     reflect.TypeOf((*pkg.Error)(nil)).Elem(),
		"ErrorCode": reflect.TypeOf((*pkg.ErrorCode)(nil)).Elem(),
		"Flags":     reflect.TypeOf((*pkg.Flags)(nil)).Elem(),
	})
}
