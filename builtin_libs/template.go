package builtin_libs

import (
	. "github.com/zgg-lang/zgg-go/runtime"
)

func libTemplate(c *Context) ValueObject {
	rv := NewObject()
	// rv.SetMember("renderHtml", NewNativeFunction("template.renderHtml", func(c *Context, this Value, args []Value) Value {
	// 	var (
	// 		templateSrc ValueStr
	// 		params      ValueObject
	// 	)
	// 	EnsureFuncParams(c, "template.renderHtml", args,
	// 		ArgRuleRequired{"templateSrc", TypeStr, &templateSrc},
	// 		ArgRuleOptionsl{"params", TypeObject, &params, NewObject()},
	// 	)
	// 	htemplate.Parse
	// 	return nil
	// }), c)
	return rv
}
