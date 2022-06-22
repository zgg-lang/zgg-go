package builtin_libs

import (
	htpl "html/template"
	"strings"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
// templateLoaderClass ValueType
)

func libTemplate(c *Context) ValueObject {
	rv := NewObject()
	rv.SetMember("renderHtml", NewNativeFunction("template.renderHtml", func(c *Context, this Value, args []Value) Value {
		var (
			templateSrc ValueStr
			params      Value
		)
		EnsureFuncParams(c, "template.renderHtml", args,
			ArgRuleRequired{"templateSrc", TypeStr, &templateSrc},
			ArgRuleOptional{"params", TypeAny, &params, NewObject()},
		)
		t, err := htpl.New("tpl").Parse(templateSrc.Value())
		if err != nil {
			c.RaiseRuntimeError("renderThml: parse template error %s", err)
		}
		var b strings.Builder
		if err := t.Execute(&b, params.ToGoValue()); err != nil {
			c.RaiseRuntimeError("renderThml: execute template error %s", err)
		}
		return NewStr(b.String())
	}), c)
	return rv
}

func init() {
	// 	templateLoaderClass = NewClassBuilder("Loader").
	// 		Constructor(func(c *Context, this ValueObject, args []Value) {
	// 			var root ValueStr
	// 			EnsureFuncParams(c, "template.Loader.__init__", args,
	// 				ArgRuleRequired{"root", TypeStr, &root},
	// 			)
	// 			this.SetMember
	// 		}).
	// 		Build()
}
