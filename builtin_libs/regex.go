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
	lib.SetMember("findAll", NewNativeFunction("findAll", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			text    ValueStr
		)
		EnsureFuncParams(c, libName+".findAll", args,
			ArgRuleRequired{"pattern", TypeStr, &pattern},
			ArgRuleRequired{"text", TypeStr, &text},
		)
		re, err := regexp.Compile(pattern.Value())
		if err != nil {
			c.RaiseRuntimeError(libName + ".findAll: pattern error " + err.Error())
			return nil
		}
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
	}, "pattern", "text"), nil)
	lib.SetMember("find", NewNativeFunction("find", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			text    ValueStr
			offset  ValueInt
		)
		EnsureFuncParams(c, libName+".find", args,
			ArgRuleRequired{"pattern", TypeStr, &pattern},
			ArgRuleRequired{"text", TypeStr, &text},
			ArgRuleOptional{"offset", TypeInt, &offset, NewInt(0)},
		)
		re, err := regexp.Compile(pattern.Value())
		if err != nil {
			c.RaiseRuntimeError(libName + ".find: pattern error " + err.Error())
			return nil
		}
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
	}, "pattern", "text", "offset"), nil)
	lib.SetMember("replaceAll", NewNativeFunction("replaceAll", func(c *Context, this Value, args []Value) Value {
		var (
			pattern ValueStr
			src     ValueStr
			repl    ValueStr
		)
		EnsureFuncParams(c, libName+".replaceAll", args,
			ArgRuleRequired{"pattern", TypeStr, &pattern},
			ArgRuleRequired{"src", TypeStr, &src},
			ArgRuleRequired{"repl", TypeStr, &repl},
		)
		p, err := regexp.Compile(pattern.Value())
		if err != nil {
			c.RaiseRuntimeError("invalid regexp %s", pattern.Value())
		}
		return NewStr(p.ReplaceAllString(src.Value(), repl.Value()))
	}, "pattern", "src", "repl"), nil)
	return lib
}

func libRegex(c *Context) ValueObject {
	return libRegexInner(c, "regex", regexMakeMatchGroupArray)
}

func libRegex2(c *Context) ValueObject {
	return libRegexInner(c, "regex2", regexMakeMatchGroupObject)
}
