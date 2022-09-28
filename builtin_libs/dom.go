package builtin_libs

import (
	. "github.com/zgg-lang/zgg-go/runtime"

	"github.com/PuerkitoBio/goquery"
)

func domLoadFromUrl(c *Context, this Value, args []Value) Value {
	var url ValueStr
	EnsureFuncParams(c, "fromUrl", args, ArgRuleRequired("url", TypeStr, &url))
	doc, err := goquery.NewDocument(url.Value())
	if err != nil {
		c.RaiseRuntimeError("dom.fromUrl: load from url fail: %s", err)
	}
	return NewGoValue(doc)
}

func libDom(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("__call__", NewNativeFunction("__call__", domLoadFromUrl), c)
	lib.SetMember("fromUrl", NewNativeFunction("fromUrl", domLoadFromUrl), c)
	return lib
}
