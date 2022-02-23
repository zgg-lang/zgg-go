package stdgolibs

import (
	pkg "go/ast"

	"reflect"
)

func init() {
	registerValues("go/ast", map[string]reflect.Value{
		// Functions
		"NewPackage":        reflect.ValueOf(pkg.NewPackage),
		"NewScope":          reflect.ValueOf(pkg.NewScope),
		"NewObj":            reflect.ValueOf(pkg.NewObj),
		"Walk":              reflect.ValueOf(pkg.Walk),
		"Inspect":           reflect.ValueOf(pkg.Inspect),
		"NewIdent":          reflect.ValueOf(pkg.NewIdent),
		"IsExported":        reflect.ValueOf(pkg.IsExported),
		"NewCommentMap":     reflect.ValueOf(pkg.NewCommentMap),
		"FileExports":       reflect.ValueOf(pkg.FileExports),
		"PackageExports":    reflect.ValueOf(pkg.PackageExports),
		"FilterDecl":        reflect.ValueOf(pkg.FilterDecl),
		"FilterFile":        reflect.ValueOf(pkg.FilterFile),
		"FilterPackage":     reflect.ValueOf(pkg.FilterPackage),
		"MergePackageFiles": reflect.ValueOf(pkg.MergePackageFiles),
		"SortImports":       reflect.ValueOf(pkg.SortImports),
		"NotNilFilter":      reflect.ValueOf(pkg.NotNilFilter),
		"Fprint":            reflect.ValueOf(pkg.Fprint),
		"Print":             reflect.ValueOf(pkg.Print),

		// Consts

		"Bad":                        reflect.ValueOf(pkg.Bad),
		"Pkg":                        reflect.ValueOf(pkg.Pkg),
		"Con":                        reflect.ValueOf(pkg.Con),
		"Typ":                        reflect.ValueOf(pkg.Typ),
		"Var":                        reflect.ValueOf(pkg.Var),
		"Fun":                        reflect.ValueOf(pkg.Fun),
		"Lbl":                        reflect.ValueOf(pkg.Lbl),
		"SEND":                       reflect.ValueOf(pkg.SEND),
		"RECV":                       reflect.ValueOf(pkg.RECV),
		"FilterFuncDuplicates":       reflect.ValueOf(pkg.FilterFuncDuplicates),
		"FilterUnassociatedComments": reflect.ValueOf(pkg.FilterUnassociatedComments),
		"FilterImportDuplicates":     reflect.ValueOf(pkg.FilterImportDuplicates),

		// Variables

	})
	registerTypes("go/ast", map[string]reflect.Type{
		// Non interfaces

		"Importer":       reflect.TypeOf((*pkg.Importer)(nil)).Elem(),
		"Scope":          reflect.TypeOf((*pkg.Scope)(nil)).Elem(),
		"Object":         reflect.TypeOf((*pkg.Object)(nil)).Elem(),
		"ObjKind":        reflect.TypeOf((*pkg.ObjKind)(nil)).Elem(),
		"Comment":        reflect.TypeOf((*pkg.Comment)(nil)).Elem(),
		"CommentGroup":   reflect.TypeOf((*pkg.CommentGroup)(nil)).Elem(),
		"Field":          reflect.TypeOf((*pkg.Field)(nil)).Elem(),
		"FieldList":      reflect.TypeOf((*pkg.FieldList)(nil)).Elem(),
		"BadExpr":        reflect.TypeOf((*pkg.BadExpr)(nil)).Elem(),
		"Ident":          reflect.TypeOf((*pkg.Ident)(nil)).Elem(),
		"Ellipsis":       reflect.TypeOf((*pkg.Ellipsis)(nil)).Elem(),
		"BasicLit":       reflect.TypeOf((*pkg.BasicLit)(nil)).Elem(),
		"FuncLit":        reflect.TypeOf((*pkg.FuncLit)(nil)).Elem(),
		"CompositeLit":   reflect.TypeOf((*pkg.CompositeLit)(nil)).Elem(),
		"ParenExpr":      reflect.TypeOf((*pkg.ParenExpr)(nil)).Elem(),
		"SelectorExpr":   reflect.TypeOf((*pkg.SelectorExpr)(nil)).Elem(),
		"IndexExpr":      reflect.TypeOf((*pkg.IndexExpr)(nil)).Elem(),
		"SliceExpr":      reflect.TypeOf((*pkg.SliceExpr)(nil)).Elem(),
		"TypeAssertExpr": reflect.TypeOf((*pkg.TypeAssertExpr)(nil)).Elem(),
		"CallExpr":       reflect.TypeOf((*pkg.CallExpr)(nil)).Elem(),
		"StarExpr":       reflect.TypeOf((*pkg.StarExpr)(nil)).Elem(),
		"UnaryExpr":      reflect.TypeOf((*pkg.UnaryExpr)(nil)).Elem(),
		"BinaryExpr":     reflect.TypeOf((*pkg.BinaryExpr)(nil)).Elem(),
		"KeyValueExpr":   reflect.TypeOf((*pkg.KeyValueExpr)(nil)).Elem(),
		"ChanDir":        reflect.TypeOf((*pkg.ChanDir)(nil)).Elem(),
		"ArrayType":      reflect.TypeOf((*pkg.ArrayType)(nil)).Elem(),
		"StructType":     reflect.TypeOf((*pkg.StructType)(nil)).Elem(),
		"FuncType":       reflect.TypeOf((*pkg.FuncType)(nil)).Elem(),
		"InterfaceType":  reflect.TypeOf((*pkg.InterfaceType)(nil)).Elem(),
		"MapType":        reflect.TypeOf((*pkg.MapType)(nil)).Elem(),
		"ChanType":       reflect.TypeOf((*pkg.ChanType)(nil)).Elem(),
		"BadStmt":        reflect.TypeOf((*pkg.BadStmt)(nil)).Elem(),
		"DeclStmt":       reflect.TypeOf((*pkg.DeclStmt)(nil)).Elem(),
		"EmptyStmt":      reflect.TypeOf((*pkg.EmptyStmt)(nil)).Elem(),
		"LabeledStmt":    reflect.TypeOf((*pkg.LabeledStmt)(nil)).Elem(),
		"ExprStmt":       reflect.TypeOf((*pkg.ExprStmt)(nil)).Elem(),
		"SendStmt":       reflect.TypeOf((*pkg.SendStmt)(nil)).Elem(),
		"IncDecStmt":     reflect.TypeOf((*pkg.IncDecStmt)(nil)).Elem(),
		"AssignStmt":     reflect.TypeOf((*pkg.AssignStmt)(nil)).Elem(),
		"GoStmt":         reflect.TypeOf((*pkg.GoStmt)(nil)).Elem(),
		"DeferStmt":      reflect.TypeOf((*pkg.DeferStmt)(nil)).Elem(),
		"ReturnStmt":     reflect.TypeOf((*pkg.ReturnStmt)(nil)).Elem(),
		"BranchStmt":     reflect.TypeOf((*pkg.BranchStmt)(nil)).Elem(),
		"BlockStmt":      reflect.TypeOf((*pkg.BlockStmt)(nil)).Elem(),
		"IfStmt":         reflect.TypeOf((*pkg.IfStmt)(nil)).Elem(),
		"CaseClause":     reflect.TypeOf((*pkg.CaseClause)(nil)).Elem(),
		"SwitchStmt":     reflect.TypeOf((*pkg.SwitchStmt)(nil)).Elem(),
		"TypeSwitchStmt": reflect.TypeOf((*pkg.TypeSwitchStmt)(nil)).Elem(),
		"CommClause":     reflect.TypeOf((*pkg.CommClause)(nil)).Elem(),
		"SelectStmt":     reflect.TypeOf((*pkg.SelectStmt)(nil)).Elem(),
		"ForStmt":        reflect.TypeOf((*pkg.ForStmt)(nil)).Elem(),
		"RangeStmt":      reflect.TypeOf((*pkg.RangeStmt)(nil)).Elem(),
		"ImportSpec":     reflect.TypeOf((*pkg.ImportSpec)(nil)).Elem(),
		"ValueSpec":      reflect.TypeOf((*pkg.ValueSpec)(nil)).Elem(),
		"TypeSpec":       reflect.TypeOf((*pkg.TypeSpec)(nil)).Elem(),
		"BadDecl":        reflect.TypeOf((*pkg.BadDecl)(nil)).Elem(),
		"GenDecl":        reflect.TypeOf((*pkg.GenDecl)(nil)).Elem(),
		"FuncDecl":       reflect.TypeOf((*pkg.FuncDecl)(nil)).Elem(),
		"File":           reflect.TypeOf((*pkg.File)(nil)).Elem(),
		"Package":        reflect.TypeOf((*pkg.Package)(nil)).Elem(),
		"CommentMap":     reflect.TypeOf((*pkg.CommentMap)(nil)).Elem(),
		"Filter":         reflect.TypeOf((*pkg.Filter)(nil)).Elem(),
		"MergeMode":      reflect.TypeOf((*pkg.MergeMode)(nil)).Elem(),
		"FieldFilter":    reflect.TypeOf((*pkg.FieldFilter)(nil)).Elem(),
	})
}
