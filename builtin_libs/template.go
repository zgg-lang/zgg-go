package builtin_libs

import (
	htpl "html/template"
	"io"
	"strings"
	ttpl "text/template"

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
			writer      Value
		)
		EnsureFuncParams(c, "template.renderHtml", args,
			ArgRuleRequired("templateSrc", TypeStr, &templateSrc),
			ArgRuleOptional("params", TypeAny, &params, NewObject()),
			ArgRuleOptional("writer", TypeGoValue, &writer, nil),
		)
		t, err := htpl.New("tpl").Parse(templateSrc.Value())
		if err != nil {
			c.RaiseRuntimeError("renderHtml: parse template error %s", err)
		}
		if writer != nil {
			if w, ok := writer.ToGoValue().(io.Writer); ok {
				if err := t.Execute(w, params.ToGoValue()); err != nil {
					c.RaiseRuntimeError("renderHtml: execute template error %s", err)
				}
			} else {
				c.RaiseRuntimeError("invalid writer %+v", writer.ToGoValue())
			}
			return Undefined()
		} else {
			var b strings.Builder
			if err := t.Execute(&b, params.ToGoValue()); err != nil {
				c.RaiseRuntimeError("renderHtml: execute template error %s", err)
			}
			return NewStr(b.String())
		}
	}), c)
	rv.SetMember("renderText", NewNativeFunction("template.renderText", func(c *Context, this Value, args []Value) Value {
		var (
			templateSrc ValueStr
			params      Value
			writer      Value
		)
		EnsureFuncParams(c, "template.renderText", args,
			ArgRuleRequired("templateSrc", TypeStr, &templateSrc),
			ArgRuleOptional("params", TypeAny, &params, NewObject()),
			ArgRuleOptional("writer", TypeGoValue, &writer, nil),
		)
		t, err := ttpl.New("tpl").Parse(templateSrc.Value())
		if err != nil {
			c.RaiseRuntimeError("renderThml: parse template error %s", err)
		}
		if writer != nil {
			if w, ok := writer.ToGoValue().(io.Writer); ok {
				if err := t.Execute(w, params.ToGoValue()); err != nil {
					c.RaiseRuntimeError("renderText: execute template error %s", err)
				}
			} else {
				c.RaiseRuntimeError("invalid writer %+v", writer.ToGoValue())
			}
			return Undefined()
		} else {
			var b strings.Builder
			if err := t.Execute(&b, params.ToGoValue()); err != nil {
				c.RaiseRuntimeError("renderText: execute template error %s", err)
			}
			return NewStr(b.String())
		}
	}), c)
	return rv
}

func init() {
	// 	templateLoaderClass = NewClassBuilder("Loader").
	// 		Constructor(func(c *Context, this ValueObject, args []Value) {
	// 			var root ValueStr
	// 			EnsureFuncParams(c, "template.Loader.__init__", args,
	// 				ArgRuleRequired("root", TypeStr, &root),
	// 			)
	// 			this.SetMember
	// 		}).
	// 		Build()
}
