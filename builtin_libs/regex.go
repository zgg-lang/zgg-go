package builtin_libs

import (
	"regexp"
	"unicode/utf8"

	. "github.com/zgg-lang/zgg-go/runtime"
)

func regexStrRunes(s string, begin, end int) int64 {
	return int64(utf8.RuneCountInString(s[begin:end]))
}

func regexMakeMatchGroupArray(c *Context, s string, begin, end int) Value {
	if begin < 0 || end < begin {
		return NewArrayByValues(NewStr(""), NewInt(-1), NewInt(-1))
	}
	oBegin := regexStrRunes(s, 0, begin)
	oEnd := oBegin + regexStrRunes(s, begin, end)
	return NewArrayByValues(
		NewStr(s[begin:end]),
		NewInt(oBegin),
		NewInt(oEnd),
	)
}

func regexMakeMatchGroupObject(c *Context, s string, begin, end int) Value {
	rv := NewObject()
	if begin < 0 || end < begin {
		rv.SetMember("text", NewStr(""), c)
		rv.SetMember("begin", NewInt(-1), c)
		rv.SetMember("end", NewInt(-1), c)
	} else {
		oBegin := regexStrRunes(s, 0, begin)
		oEnd := oBegin + regexStrRunes(s, begin, end)
		rv.SetMember("text", NewStr(s[begin:end]), c)
		rv.SetMember("begin", NewInt(oBegin), c)
		rv.SetMember("end", NewInt(oEnd), c)
	}
	return rv
}

func libRegexInner(c *Context, libName string, regexMakeMatchGroup func(*Context, string, int, int) Value) ValueObject {
	lib := NewObject()
	cls := buildRegexpClass(regexMakeMatchGroup)
	lib.SetMember("Regexp", cls, nil)
	lib.SetMember("findAll", NewNativeFunction("findAll", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			text    ValueStr
		)
		EnsureFuncParams(c, libName+".findAll", args,
			ArgRuleRequired("pattern", TypeStr, &pattern),
			ArgRuleRequired("text", TypeStr, &text),
		)
		re := NewObjectAndInit(cls, c, pattern)
		return c.InvokeMethod(re, "findAll", Args(args[1:]...))
	}, "pattern", "text"), nil)
	lib.SetMember("find", NewNativeFunction("find", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			text    ValueStr
			offset  ValueInt
		)
		EnsureFuncParams(c, libName+".find", args,
			ArgRuleRequired("pattern", TypeStr, &pattern),
			ArgRuleRequired("text", TypeStr, &text),
			ArgRuleOptional("offset", TypeInt, &offset, NewInt(0)),
		)
		re := NewObjectAndInit(cls, c, pattern)
		return c.InvokeMethod(re, "find", Args(args[1:]...))
	}, "pattern", "text", "offset"), nil)
	lib.SetMember("split", NewNativeFunction("split", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			text    ValueStr
		)
		EnsureFuncParams(c, libName+".split", args,
			ArgRuleRequired("pattern", TypeStr, &pattern),
			ArgRuleRequired("text", TypeStr, &text),
		)
		re := NewObjectAndInit(cls, c, pattern)
		return c.InvokeMethod(re, "split", Args(args[1:]...))
	}, "pattern", "text"), nil)
	lib.SetMember("replaceAll", NewNativeFunction("replaceAll", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			src     ValueStr
			repl    Value
		)
		EnsureFuncParams(c, libName+".replaceAll", args,
			ArgRuleRequired("pattern", TypeStr, &pattern),
			ArgRuleRequired("src", TypeStr, &src),
			ArgRuleRequired("repl", TypeAny, &repl),
		)
		re := NewObjectAndInit(cls, c, pattern)
		return c.InvokeMethod(re, "replaceAll", Args(args[1:]...))
	}, "pattern", "src", "repl"), nil)
	return lib
}

func libRegex(c *Context) ValueObject {
	return libRegexInner(c, "regex", regexMakeMatchGroupArray)
}

func libRegex2(c *Context) ValueObject {
	return libRegexInner(c, "regex2", regexMakeMatchGroupObject)
}

func buildRegexpClass(regexMakeMatchGroup func(*Context, string, int, int) Value) ValueType {
	return NewClassBuilder("Regexp").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var p ValueStr
			EnsureFuncParams(c, "Regexp.__init__", args, ArgRuleRequired("pattern", TypeStr, &p))
			re, err := regexp.Compile(p.Value())
			if err != nil {
				c.RaiseRuntimeError("compile regexp %s error %+v", p.Value(), err)
			}
			this.Reserved = re
		}).
		Method("find", func(c *Context, this ValueObject, args []Value) Value {
			var (
				text   ValueStr
				offset ValueInt
			)
			EnsureFuncParams(c, "Regexp.find", args,
				ArgRuleRequired("text", TypeStr, &text),
				ArgRuleOptional("offset", TypeInt, &offset, NewInt(0)),
			)
			re := this.Reserved.(*regexp.Regexp)
			n := offset.AsInt()
			runes := text.Runes()
			if n >= len(runes) || n < -len(runes) {
				return NewArray(0)
			}
			if n < 0 {
				n += len(runes)
			}
			s := text.Value()
			if n > 0 {
				nInBytes := 0
				for i := 0; i < n; i++ {
					nInBytes += utf8.RuneLen(runes[i])
				}
				n = nInBytes
			}
			rv := re.FindStringSubmatchIndex(s[n:])
			if rv == nil {
				return NewArray(0)
			}
			rvItem := NewArray(len(rv) / 2)
			for i := 0; i < len(rv); i += 2 {
				begin, end := rv[i]+n, rv[i+1]+n
				rvItem.PushBack(regexMakeMatchGroup(c, s, begin, end))
			}
			return rvItem
		}, "text", "offset").
		Method("findAll", func(c *Context, this ValueObject, args []Value) Value {
			var text ValueStr
			EnsureFuncParams(c, "Regexp.findAll", args,
				ArgRuleRequired("text", TypeStr, &text),
			)
			re := this.Reserved.(*regexp.Regexp)
			s := text.Value()
			rv := re.FindAllStringSubmatchIndex(s, -1)
			if rv == nil {
				return NewArray(0)
			}
			rvArr := NewArray(len(rv))
			for _, rs := range rv {
				rvItem := NewArray(len(rs) / 2)
				for i := 0; i < len(rs); i += 2 {
					begin, end := rs[i], rs[i+1]
					rvItem.PushBack(regexMakeMatchGroup(c, s, begin, end))
				}
				rvArr.PushBack(rvItem)
			}
			return rvArr
		}, "text").
		Method("split", func(c *Context, this ValueObject, args []Value) Value {
			var (
				text ValueStr
				n    ValueInt
			)
			EnsureFuncParams(c, "Regexp.split", args,
				ArgRuleRequired("text", TypeStr, &text),
				ArgRuleOptional("n", TypeInt, &n, NewInt(-1)),
			)
			var (
				re   = this.Reserved.(*regexp.Regexp)
				subs = re.Split(text.Value(), n.AsInt())
				rv   = NewArray(len(subs))
			)
			for _, sub := range subs {
				rv.PushBack(NewStr(sub))
			}
			return rv
		}, "text").
		Method("replaceAll", func(c *Context, this ValueObject, args []Value) Value {
			var (
				src       ValueStr
				replStr   ValueStr
				replFunc  ValueCallable
				replWhich int
			)
			EnsureFuncParams(c, "Regexp.replaceAll", args,
				ArgRuleRequired("src", TypeStr, &src),
				ArgRuleOneOf("repl",
					[]ValueType{TypeStr, TypeCallable},
					[]any{&replStr, &replFunc},
					&replWhich, nil, nil,
				),
			)
			p := this.Reserved.(*regexp.Regexp)
			if replWhich == 1 {
				return NewStr(p.ReplaceAllStringFunc(src.Value(), func(s string) string {
					c.Invoke(replFunc, nil, Args(NewStr(s)))
					return c.RetVal.ToString(c)
				}))
			} else {
				return NewStr(p.ReplaceAllString(src.Value(), replStr.Value()))
			}
		}, "src", "repl").
		Build()
}
