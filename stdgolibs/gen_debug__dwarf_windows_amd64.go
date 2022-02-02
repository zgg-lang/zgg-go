package stdgolibs

import (
	pkg "debug/dwarf"

	"reflect"
)

func init() {
	registerValues("debug/dwarf", map[string]reflect.Value{
		// Functions
		"New": reflect.ValueOf(pkg.New),

		// Consts

		"AttrSibling":               reflect.ValueOf(pkg.AttrSibling),
		"AttrLocation":              reflect.ValueOf(pkg.AttrLocation),
		"AttrName":                  reflect.ValueOf(pkg.AttrName),
		"AttrOrdering":              reflect.ValueOf(pkg.AttrOrdering),
		"AttrByteSize":              reflect.ValueOf(pkg.AttrByteSize),
		"AttrBitOffset":             reflect.ValueOf(pkg.AttrBitOffset),
		"AttrBitSize":               reflect.ValueOf(pkg.AttrBitSize),
		"AttrStmtList":              reflect.ValueOf(pkg.AttrStmtList),
		"AttrLowpc":                 reflect.ValueOf(pkg.AttrLowpc),
		"AttrHighpc":                reflect.ValueOf(pkg.AttrHighpc),
		"AttrLanguage":              reflect.ValueOf(pkg.AttrLanguage),
		"AttrDiscr":                 reflect.ValueOf(pkg.AttrDiscr),
		"AttrDiscrValue":            reflect.ValueOf(pkg.AttrDiscrValue),
		"AttrVisibility":            reflect.ValueOf(pkg.AttrVisibility),
		"AttrImport":                reflect.ValueOf(pkg.AttrImport),
		"AttrStringLength":          reflect.ValueOf(pkg.AttrStringLength),
		"AttrCommonRef":             reflect.ValueOf(pkg.AttrCommonRef),
		"AttrCompDir":               reflect.ValueOf(pkg.AttrCompDir),
		"AttrConstValue":            reflect.ValueOf(pkg.AttrConstValue),
		"AttrContainingType":        reflect.ValueOf(pkg.AttrContainingType),
		"AttrDefaultValue":          reflect.ValueOf(pkg.AttrDefaultValue),
		"AttrInline":                reflect.ValueOf(pkg.AttrInline),
		"AttrIsOptional":            reflect.ValueOf(pkg.AttrIsOptional),
		"AttrLowerBound":            reflect.ValueOf(pkg.AttrLowerBound),
		"AttrProducer":              reflect.ValueOf(pkg.AttrProducer),
		"AttrPrototyped":            reflect.ValueOf(pkg.AttrPrototyped),
		"AttrReturnAddr":            reflect.ValueOf(pkg.AttrReturnAddr),
		"AttrStartScope":            reflect.ValueOf(pkg.AttrStartScope),
		"AttrStrideSize":            reflect.ValueOf(pkg.AttrStrideSize),
		"AttrUpperBound":            reflect.ValueOf(pkg.AttrUpperBound),
		"AttrAbstractOrigin":        reflect.ValueOf(pkg.AttrAbstractOrigin),
		"AttrAccessibility":         reflect.ValueOf(pkg.AttrAccessibility),
		"AttrAddrClass":             reflect.ValueOf(pkg.AttrAddrClass),
		"AttrArtificial":            reflect.ValueOf(pkg.AttrArtificial),
		"AttrBaseTypes":             reflect.ValueOf(pkg.AttrBaseTypes),
		"AttrCalling":               reflect.ValueOf(pkg.AttrCalling),
		"AttrCount":                 reflect.ValueOf(pkg.AttrCount),
		"AttrDataMemberLoc":         reflect.ValueOf(pkg.AttrDataMemberLoc),
		"AttrDeclColumn":            reflect.ValueOf(pkg.AttrDeclColumn),
		"AttrDeclFile":              reflect.ValueOf(pkg.AttrDeclFile),
		"AttrDeclLine":              reflect.ValueOf(pkg.AttrDeclLine),
		"AttrDeclaration":           reflect.ValueOf(pkg.AttrDeclaration),
		"AttrDiscrList":             reflect.ValueOf(pkg.AttrDiscrList),
		"AttrEncoding":              reflect.ValueOf(pkg.AttrEncoding),
		"AttrExternal":              reflect.ValueOf(pkg.AttrExternal),
		"AttrFrameBase":             reflect.ValueOf(pkg.AttrFrameBase),
		"AttrFriend":                reflect.ValueOf(pkg.AttrFriend),
		"AttrIdentifierCase":        reflect.ValueOf(pkg.AttrIdentifierCase),
		"AttrMacroInfo":             reflect.ValueOf(pkg.AttrMacroInfo),
		"AttrNamelistItem":          reflect.ValueOf(pkg.AttrNamelistItem),
		"AttrPriority":              reflect.ValueOf(pkg.AttrPriority),
		"AttrSegment":               reflect.ValueOf(pkg.AttrSegment),
		"AttrSpecification":         reflect.ValueOf(pkg.AttrSpecification),
		"AttrStaticLink":            reflect.ValueOf(pkg.AttrStaticLink),
		"AttrType":                  reflect.ValueOf(pkg.AttrType),
		"AttrUseLocation":           reflect.ValueOf(pkg.AttrUseLocation),
		"AttrVarParam":              reflect.ValueOf(pkg.AttrVarParam),
		"AttrVirtuality":            reflect.ValueOf(pkg.AttrVirtuality),
		"AttrVtableElemLoc":         reflect.ValueOf(pkg.AttrVtableElemLoc),
		"AttrAllocated":             reflect.ValueOf(pkg.AttrAllocated),
		"AttrAssociated":            reflect.ValueOf(pkg.AttrAssociated),
		"AttrDataLocation":          reflect.ValueOf(pkg.AttrDataLocation),
		"AttrStride":                reflect.ValueOf(pkg.AttrStride),
		"AttrEntrypc":               reflect.ValueOf(pkg.AttrEntrypc),
		"AttrUseUTF8":               reflect.ValueOf(pkg.AttrUseUTF8),
		"AttrExtension":             reflect.ValueOf(pkg.AttrExtension),
		"AttrRanges":                reflect.ValueOf(pkg.AttrRanges),
		"AttrTrampoline":            reflect.ValueOf(pkg.AttrTrampoline),
		"AttrCallColumn":            reflect.ValueOf(pkg.AttrCallColumn),
		"AttrCallFile":              reflect.ValueOf(pkg.AttrCallFile),
		"AttrCallLine":              reflect.ValueOf(pkg.AttrCallLine),
		"AttrDescription":           reflect.ValueOf(pkg.AttrDescription),
		"AttrBinaryScale":           reflect.ValueOf(pkg.AttrBinaryScale),
		"AttrDecimalScale":          reflect.ValueOf(pkg.AttrDecimalScale),
		"AttrSmall":                 reflect.ValueOf(pkg.AttrSmall),
		"AttrDecimalSign":           reflect.ValueOf(pkg.AttrDecimalSign),
		"AttrDigitCount":            reflect.ValueOf(pkg.AttrDigitCount),
		"AttrPictureString":         reflect.ValueOf(pkg.AttrPictureString),
		"AttrMutable":               reflect.ValueOf(pkg.AttrMutable),
		"AttrThreadsScaled":         reflect.ValueOf(pkg.AttrThreadsScaled),
		"AttrExplicit":              reflect.ValueOf(pkg.AttrExplicit),
		"AttrObjectPointer":         reflect.ValueOf(pkg.AttrObjectPointer),
		"AttrEndianity":             reflect.ValueOf(pkg.AttrEndianity),
		"AttrElemental":             reflect.ValueOf(pkg.AttrElemental),
		"AttrPure":                  reflect.ValueOf(pkg.AttrPure),
		"AttrRecursive":             reflect.ValueOf(pkg.AttrRecursive),
		"AttrSignature":             reflect.ValueOf(pkg.AttrSignature),
		"AttrMainSubprogram":        reflect.ValueOf(pkg.AttrMainSubprogram),
		"AttrDataBitOffset":         reflect.ValueOf(pkg.AttrDataBitOffset),
		"AttrConstExpr":             reflect.ValueOf(pkg.AttrConstExpr),
		"AttrEnumClass":             reflect.ValueOf(pkg.AttrEnumClass),
		"AttrLinkageName":           reflect.ValueOf(pkg.AttrLinkageName),
		"AttrStringLengthBitSize":   reflect.ValueOf(pkg.AttrStringLengthBitSize),
		"AttrStringLengthByteSize":  reflect.ValueOf(pkg.AttrStringLengthByteSize),
		"AttrRank":                  reflect.ValueOf(pkg.AttrRank),
		"AttrStrOffsetsBase":        reflect.ValueOf(pkg.AttrStrOffsetsBase),
		"AttrAddrBase":              reflect.ValueOf(pkg.AttrAddrBase),
		"AttrRnglistsBase":          reflect.ValueOf(pkg.AttrRnglistsBase),
		"AttrDwoName":               reflect.ValueOf(pkg.AttrDwoName),
		"AttrReference":             reflect.ValueOf(pkg.AttrReference),
		"AttrRvalueReference":       reflect.ValueOf(pkg.AttrRvalueReference),
		"AttrMacros":                reflect.ValueOf(pkg.AttrMacros),
		"AttrCallAllCalls":          reflect.ValueOf(pkg.AttrCallAllCalls),
		"AttrCallAllSourceCalls":    reflect.ValueOf(pkg.AttrCallAllSourceCalls),
		"AttrCallAllTailCalls":      reflect.ValueOf(pkg.AttrCallAllTailCalls),
		"AttrCallReturnPC":          reflect.ValueOf(pkg.AttrCallReturnPC),
		"AttrCallValue":             reflect.ValueOf(pkg.AttrCallValue),
		"AttrCallOrigin":            reflect.ValueOf(pkg.AttrCallOrigin),
		"AttrCallParameter":         reflect.ValueOf(pkg.AttrCallParameter),
		"AttrCallPC":                reflect.ValueOf(pkg.AttrCallPC),
		"AttrCallTailCall":          reflect.ValueOf(pkg.AttrCallTailCall),
		"AttrCallTarget":            reflect.ValueOf(pkg.AttrCallTarget),
		"AttrCallTargetClobbered":   reflect.ValueOf(pkg.AttrCallTargetClobbered),
		"AttrCallDataLocation":      reflect.ValueOf(pkg.AttrCallDataLocation),
		"AttrCallDataValue":         reflect.ValueOf(pkg.AttrCallDataValue),
		"AttrNoreturn":              reflect.ValueOf(pkg.AttrNoreturn),
		"AttrAlignment":             reflect.ValueOf(pkg.AttrAlignment),
		"AttrExportSymbols":         reflect.ValueOf(pkg.AttrExportSymbols),
		"AttrDeleted":               reflect.ValueOf(pkg.AttrDeleted),
		"AttrDefaulted":             reflect.ValueOf(pkg.AttrDefaulted),
		"AttrLoclistsBase":          reflect.ValueOf(pkg.AttrLoclistsBase),
		"TagArrayType":              reflect.ValueOf(pkg.TagArrayType),
		"TagClassType":              reflect.ValueOf(pkg.TagClassType),
		"TagEntryPoint":             reflect.ValueOf(pkg.TagEntryPoint),
		"TagEnumerationType":        reflect.ValueOf(pkg.TagEnumerationType),
		"TagFormalParameter":        reflect.ValueOf(pkg.TagFormalParameter),
		"TagImportedDeclaration":    reflect.ValueOf(pkg.TagImportedDeclaration),
		"TagLabel":                  reflect.ValueOf(pkg.TagLabel),
		"TagLexDwarfBlock":          reflect.ValueOf(pkg.TagLexDwarfBlock),
		"TagMember":                 reflect.ValueOf(pkg.TagMember),
		"TagPointerType":            reflect.ValueOf(pkg.TagPointerType),
		"TagReferenceType":          reflect.ValueOf(pkg.TagReferenceType),
		"TagCompileUnit":            reflect.ValueOf(pkg.TagCompileUnit),
		"TagStringType":             reflect.ValueOf(pkg.TagStringType),
		"TagStructType":             reflect.ValueOf(pkg.TagStructType),
		"TagSubroutineType":         reflect.ValueOf(pkg.TagSubroutineType),
		"TagTypedef":                reflect.ValueOf(pkg.TagTypedef),
		"TagUnionType":              reflect.ValueOf(pkg.TagUnionType),
		"TagUnspecifiedParameters":  reflect.ValueOf(pkg.TagUnspecifiedParameters),
		"TagVariant":                reflect.ValueOf(pkg.TagVariant),
		"TagCommonDwarfBlock":       reflect.ValueOf(pkg.TagCommonDwarfBlock),
		"TagCommonInclusion":        reflect.ValueOf(pkg.TagCommonInclusion),
		"TagInheritance":            reflect.ValueOf(pkg.TagInheritance),
		"TagInlinedSubroutine":      reflect.ValueOf(pkg.TagInlinedSubroutine),
		"TagModule":                 reflect.ValueOf(pkg.TagModule),
		"TagPtrToMemberType":        reflect.ValueOf(pkg.TagPtrToMemberType),
		"TagSetType":                reflect.ValueOf(pkg.TagSetType),
		"TagSubrangeType":           reflect.ValueOf(pkg.TagSubrangeType),
		"TagWithStmt":               reflect.ValueOf(pkg.TagWithStmt),
		"TagAccessDeclaration":      reflect.ValueOf(pkg.TagAccessDeclaration),
		"TagBaseType":               reflect.ValueOf(pkg.TagBaseType),
		"TagCatchDwarfBlock":        reflect.ValueOf(pkg.TagCatchDwarfBlock),
		"TagConstType":              reflect.ValueOf(pkg.TagConstType),
		"TagConstant":               reflect.ValueOf(pkg.TagConstant),
		"TagEnumerator":             reflect.ValueOf(pkg.TagEnumerator),
		"TagFileType":               reflect.ValueOf(pkg.TagFileType),
		"TagFriend":                 reflect.ValueOf(pkg.TagFriend),
		"TagNamelist":               reflect.ValueOf(pkg.TagNamelist),
		"TagNamelistItem":           reflect.ValueOf(pkg.TagNamelistItem),
		"TagPackedType":             reflect.ValueOf(pkg.TagPackedType),
		"TagSubprogram":             reflect.ValueOf(pkg.TagSubprogram),
		"TagTemplateTypeParameter":  reflect.ValueOf(pkg.TagTemplateTypeParameter),
		"TagTemplateValueParameter": reflect.ValueOf(pkg.TagTemplateValueParameter),
		"TagThrownType":             reflect.ValueOf(pkg.TagThrownType),
		"TagTryDwarfBlock":          reflect.ValueOf(pkg.TagTryDwarfBlock),
		"TagVariantPart":            reflect.ValueOf(pkg.TagVariantPart),
		"TagVariable":               reflect.ValueOf(pkg.TagVariable),
		"TagVolatileType":           reflect.ValueOf(pkg.TagVolatileType),
		"TagDwarfProcedure":         reflect.ValueOf(pkg.TagDwarfProcedure),
		"TagRestrictType":           reflect.ValueOf(pkg.TagRestrictType),
		"TagInterfaceType":          reflect.ValueOf(pkg.TagInterfaceType),
		"TagNamespace":              reflect.ValueOf(pkg.TagNamespace),
		"TagImportedModule":         reflect.ValueOf(pkg.TagImportedModule),
		"TagUnspecifiedType":        reflect.ValueOf(pkg.TagUnspecifiedType),
		"TagPartialUnit":            reflect.ValueOf(pkg.TagPartialUnit),
		"TagImportedUnit":           reflect.ValueOf(pkg.TagImportedUnit),
		"TagMutableType":            reflect.ValueOf(pkg.TagMutableType),
		"TagCondition":              reflect.ValueOf(pkg.TagCondition),
		"TagSharedType":             reflect.ValueOf(pkg.TagSharedType),
		"TagTypeUnit":               reflect.ValueOf(pkg.TagTypeUnit),
		"TagRvalueReferenceType":    reflect.ValueOf(pkg.TagRvalueReferenceType),
		"TagTemplateAlias":          reflect.ValueOf(pkg.TagTemplateAlias),
		"TagCoarrayType":            reflect.ValueOf(pkg.TagCoarrayType),
		"TagGenericSubrange":        reflect.ValueOf(pkg.TagGenericSubrange),
		"TagDynamicType":            reflect.ValueOf(pkg.TagDynamicType),
		"TagAtomicType":             reflect.ValueOf(pkg.TagAtomicType),
		"TagCallSite":               reflect.ValueOf(pkg.TagCallSite),
		"TagCallSiteParameter":      reflect.ValueOf(pkg.TagCallSiteParameter),
		"TagSkeletonUnit":           reflect.ValueOf(pkg.TagSkeletonUnit),
		"TagImmutableType":          reflect.ValueOf(pkg.TagImmutableType),
		"ClassUnknown":              reflect.ValueOf(pkg.ClassUnknown),
		"ClassAddress":              reflect.ValueOf(pkg.ClassAddress),
		"ClassBlock":                reflect.ValueOf(pkg.ClassBlock),
		"ClassConstant":             reflect.ValueOf(pkg.ClassConstant),
		"ClassExprLoc":              reflect.ValueOf(pkg.ClassExprLoc),
		"ClassFlag":                 reflect.ValueOf(pkg.ClassFlag),
		"ClassLinePtr":              reflect.ValueOf(pkg.ClassLinePtr),
		"ClassLocListPtr":           reflect.ValueOf(pkg.ClassLocListPtr),
		"ClassMacPtr":               reflect.ValueOf(pkg.ClassMacPtr),
		"ClassRangeListPtr":         reflect.ValueOf(pkg.ClassRangeListPtr),
		"ClassReference":            reflect.ValueOf(pkg.ClassReference),
		"ClassReferenceSig":         reflect.ValueOf(pkg.ClassReferenceSig),
		"ClassString":               reflect.ValueOf(pkg.ClassString),
		"ClassReferenceAlt":         reflect.ValueOf(pkg.ClassReferenceAlt),
		"ClassStringAlt":            reflect.ValueOf(pkg.ClassStringAlt),
		"ClassAddrPtr":              reflect.ValueOf(pkg.ClassAddrPtr),
		"ClassLocList":              reflect.ValueOf(pkg.ClassLocList),
		"ClassRngList":              reflect.ValueOf(pkg.ClassRngList),
		"ClassRngListsPtr":          reflect.ValueOf(pkg.ClassRngListsPtr),
		"ClassStrOffsetsPtr":        reflect.ValueOf(pkg.ClassStrOffsetsPtr),

		// Variables

		"ErrUnknownPC": reflect.ValueOf(&pkg.ErrUnknownPC),
	})
	registerTypes("debug/dwarf", map[string]reflect.Type{
		// Non interfaces

		"Attr":            reflect.TypeOf((*pkg.Attr)(nil)).Elem(),
		"Tag":             reflect.TypeOf((*pkg.Tag)(nil)).Elem(),
		"Entry":           reflect.TypeOf((*pkg.Entry)(nil)).Elem(),
		"Field":           reflect.TypeOf((*pkg.Field)(nil)).Elem(),
		"Class":           reflect.TypeOf((*pkg.Class)(nil)).Elem(),
		"Offset":          reflect.TypeOf((*pkg.Offset)(nil)).Elem(),
		"Reader":          reflect.TypeOf((*pkg.Reader)(nil)).Elem(),
		"LineReader":      reflect.TypeOf((*pkg.LineReader)(nil)).Elem(),
		"LineEntry":       reflect.TypeOf((*pkg.LineEntry)(nil)).Elem(),
		"LineFile":        reflect.TypeOf((*pkg.LineFile)(nil)).Elem(),
		"LineReaderPos":   reflect.TypeOf((*pkg.LineReaderPos)(nil)).Elem(),
		"DecodeError":     reflect.TypeOf((*pkg.DecodeError)(nil)).Elem(),
		"Data":            reflect.TypeOf((*pkg.Data)(nil)).Elem(),
		"CommonType":      reflect.TypeOf((*pkg.CommonType)(nil)).Elem(),
		"BasicType":       reflect.TypeOf((*pkg.BasicType)(nil)).Elem(),
		"CharType":        reflect.TypeOf((*pkg.CharType)(nil)).Elem(),
		"UcharType":       reflect.TypeOf((*pkg.UcharType)(nil)).Elem(),
		"IntType":         reflect.TypeOf((*pkg.IntType)(nil)).Elem(),
		"UintType":        reflect.TypeOf((*pkg.UintType)(nil)).Elem(),
		"FloatType":       reflect.TypeOf((*pkg.FloatType)(nil)).Elem(),
		"ComplexType":     reflect.TypeOf((*pkg.ComplexType)(nil)).Elem(),
		"BoolType":        reflect.TypeOf((*pkg.BoolType)(nil)).Elem(),
		"AddrType":        reflect.TypeOf((*pkg.AddrType)(nil)).Elem(),
		"UnspecifiedType": reflect.TypeOf((*pkg.UnspecifiedType)(nil)).Elem(),
		"QualType":        reflect.TypeOf((*pkg.QualType)(nil)).Elem(),
		"ArrayType":       reflect.TypeOf((*pkg.ArrayType)(nil)).Elem(),
		"VoidType":        reflect.TypeOf((*pkg.VoidType)(nil)).Elem(),
		"PtrType":         reflect.TypeOf((*pkg.PtrType)(nil)).Elem(),
		"StructType":      reflect.TypeOf((*pkg.StructType)(nil)).Elem(),
		"StructField":     reflect.TypeOf((*pkg.StructField)(nil)).Elem(),
		"EnumType":        reflect.TypeOf((*pkg.EnumType)(nil)).Elem(),
		"EnumValue":       reflect.TypeOf((*pkg.EnumValue)(nil)).Elem(),
		"FuncType":        reflect.TypeOf((*pkg.FuncType)(nil)).Elem(),
		"DotDotDotType":   reflect.TypeOf((*pkg.DotDotDotType)(nil)).Elem(),
		"TypedefType":     reflect.TypeOf((*pkg.TypedefType)(nil)).Elem(),
		"UnsupportedType": reflect.TypeOf((*pkg.UnsupportedType)(nil)).Elem(),
	})
}
