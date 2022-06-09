package runtime

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
)

type ValueBytes struct {
	*ValueBase
	v []byte
}

func NewBytes(v []byte) ValueBytes {
	return ValueBytes{ValueBase: &ValueBase{}, v: v}
}

func (v ValueBytes) GoType() reflect.Type {
	return reflect.TypeOf(v.v)
}

func (v ValueBytes) ToGoValue() interface{} {
	return v.Value()
}

func (v ValueBytes) GetIndex(index int, c *Context) Value {
	if index < 0 || index >= len(v.v) {
		return constUndefined
	}
	return NewInt(int64(v.v[index]))
}

func (v ValueBytes) SetIndex(index int, value Value, c *Context) {
	if index < 0 {
		index += v.Len()
	}
	if index < 0 || index >= v.Len() {
		c.RaiseRuntimeError(fmt.Sprintf("set bytes item error: Out of bound length %d index %d", v.Len(), index))
		return
	}
	iv := c.MustInt(value)
	if iv < 0 || iv > 255 {
		c.RaiseRuntimeError("set bytes item error: value must between 0 and 255, not %d", iv)
		return
	}
	v.v[index] = byte(iv)
}

func (v ValueBytes) GetMember(name string, c *Context) Value {
	if member, found := builtinBytesMethods[name]; found {
		return makeMember(v, member)
	}
	return getExtMember(v, name, c)
}

func (ValueBytes) Type() ValueType {
	return TypeBytes
}

func (v ValueBytes) Value() []byte {
	return v.v
}

func (v ValueBytes) IsTrue() bool {
	return len(v.v) > 0
}

func (v ValueBytes) CompareTo(other Value, c *Context) CompareResult {
	if v2, ok := other.(ValueBytes); ok {
		bs1, bs2 := v.Value(), v2.Value()
		l1, l2 := len(bs1), len(bs2)
		minLen := l1
		if l2 < l1 {
			minLen = l2
		}
		for i := 0; i < minLen; i++ {
			if bs1[i] < bs2[i] {
				return CompareResultLess
			} else if bs1[i] > bs2[i] {
				return CompareResultGreater
			}
		}
		if l1 < l2 {
			return CompareResultLess
		} else if l1 > l2 {
			return CompareResultGreater
		} else {
			return CompareResultEqual
		}
	}
	return CompareResultNotEqual
}

func (v ValueBytes) Len() int {
	return len(v.v)
}

func (v ValueBytes) ToString(*Context) string {
	return string(v.v)
}

var builtinBytesMethods = map[string]ValueCallable{
	"len": &ValueBuiltinFunction{
		name: "bytes.len",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisBytes := thisArg.(ValueBytes)
			return NewInt(int64(len(thisBytes.Value())))
		},
	},
	"slice": &ValueBuiltinFunction{
		name: "bytes.slice",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisStr := thisArg.(ValueBytes)
			charSeq := thisStr.v
			begin, end := 0, len(charSeq)
			switch len(args) {
			case 0:
			case 2:
				{
					endVal, ok := args[1].(ValueInt)
					if !ok {
						c.RaiseRuntimeError("bytes.slice arguments 1 must be an integer")
						return nil
					}
					end = int(endVal.Value())
					if end < 0 {
						end += len(charSeq)
					}
				}
				fallthrough
			case 1:
				{
					beginVal, ok := args[0].(ValueInt)
					if !ok {
						c.RaiseRuntimeError("bytes.slice arguments 0 must be an integer")
						return nil
					}
					begin = int(beginVal.Value())
					if begin < 0 {
						begin += len(charSeq)
					}
				}
			default:
				c.RaiseRuntimeError("bytes.slice requires 0~2 arguments")
				return nil
			}
			if begin < 0 {
				begin = 0
			}
			if begin >= len(charSeq) {
				begin = len(charSeq)
			}
			if end < 0 {
				end = 0
			}
			if end >= len(charSeq) {
				end = len(charSeq)
			}
			if begin >= end {
				return NewBytes([]byte{})
			}
			sliceChars := make([]byte, end-begin)
			for i := begin; i < end; i++ {
				sliceChars[i-begin] = charSeq[i]
			}
			return NewBytes(sliceChars)
		},
	},
	"hex": NewNativeFunction("bytes.hex", func(c *Context, thisArg Value, args []Value) Value {
		thisBytes := thisArg.(ValueBytes)
		return NewStr(hex.EncodeToString(thisBytes.v))
	}),
	"base64": NewNativeFunction("bytes.base64", func(c *Context, thisArg Value, args []Value) Value {
		thisBytes := thisArg.(ValueBytes)
		return NewStr(base64.StdEncoding.EncodeToString(thisBytes.v))
	}),
}

func init() {
	addMembersAndStatics(TypeBytes, builtinBytesMethods)
}
