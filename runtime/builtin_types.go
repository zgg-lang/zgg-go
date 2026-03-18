package runtime

import "sync"

const (
	builtinTypeUndefined = iota
	builtinTypeNil
	builtinTypeInt
	builtinTypeFloat
	builtinTypeBigNum
	builtinTypeBool
	builtinTypeStr
	builtinTypeBytes
	builtinTypeFunc
	builtinTypeArray
	builtinTypeObject
	builtinTypeMap
	builtinTypeType
	builtinTypeGoValue
	builtinTypeGoType
	builtinTypeCallable
	builtinTypeAny
)

var (
	TypeUndefined = NewType(builtinTypeUndefined, "Undefined")
	TypeNil       = NewType(builtinTypeNil, "Nil")
	TypeInt       = NewType(builtinTypeInt, "Int")
	TypeBigNum    = NewType(builtinTypeBigNum, "BigNum")
	TypeFloat     = NewType(builtinTypeFloat, "Float")
	TypeBool      = NewType(builtinTypeBool, "Bool")
	TypeStr       = NewType(builtinTypeStr, "Str")
	TypeBytes     = NewType(builtinTypeBytes, "Bytes")
	TypeFunc      = NewType(builtinTypeFunc, "Func")
	TypeArray     = NewType(builtinTypeArray, "Array")
	TypeObject    = NewType(builtinTypeObject, "Object")
	TypeMap       = NewTypeWithCreator(builtinTypeMap, "Map", NewMapWithPairs)
	TypeType      = NewType(builtinTypeType, "Type")
	TypeGoValue   = NewType(builtinTypeGoValue, "GoValue")
	TypeGoType    = NewType(builtinTypeGoType, "GoType")
	// Any Types
	TypeCallable = NewType(builtinTypeCallable, "AnyCallable")
	TypeAny      = NewType(builtinTypeAny, "Any")
)

var builtinTypes = map[string]ValueType{}

func init() {
	types := []ValueType{
		TypeUndefined,
		TypeNil,
		TypeInt,
		TypeBigNum,
		TypeFloat,
		TypeBool,
		TypeStr,
		TypeBytes,
		TypeFunc,
		TypeArray,
		TypeObject,
		TypeMap,
		TypeType,
		TypeGoValue,
		TypeGoType,
		// TypeCallable is not available in zgg code
	}
	for _, t := range types {
		builtinTypes[t.Name] = t
	}
}

var (
	typeArrayOfMap = make(map[int]ValueType)
	typeArrayOfRW  sync.RWMutex
)

func TypeArrayOf(t ValueType) ValueType {
	if t == nil {
		return TypeArray
	}
	typeArrayOfRW.RLock()
	t2 := typeArrayOfMap[t.TypeId]
	typeArrayOfRW.RUnlock()
	if t2 != nil {
		return t2
	}
	typeArrayOfRW.Lock()
	defer typeArrayOfRW.Unlock()
	if t2 = typeArrayOfMap[t.TypeId]; t2 != nil {
		return t2
	}
	t2 = NewType(NextTypeId(), "ArrayOf:"+t.Name)
	t2.Bases = []ValueType{t}
	typeArrayOfMap[t.TypeId] = t2
	return t2
}
