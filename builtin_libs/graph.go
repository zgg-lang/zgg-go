package builtin_libs

import (
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	graphImageType ValueType
)

func libGraph(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("readFile", NewNativeFunction("readFile", func(c *Context, this Value, args []Value) Value {
		return nil
	}), c)
	return lib
}

func init() {
	graphImageType = NewClassBuilder("Image").
		Build()
}
