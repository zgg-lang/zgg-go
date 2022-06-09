package image

import (
	. "github.com/zgg-lang/zgg-go/runtime"
)

var ImageClass = NewClassBuilder("Image").
	Constructor(func(c *Context, this ValueObject, args []Value) {

	}).
	Build()
