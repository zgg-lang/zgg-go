package runtime

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

type ValueStr struct {
	*ValueBase
	v []rune
	s string
}

func NewStr(v string, args ...interface{}) ValueStr {
	if v == "" {
		return emptyStr
	}
	if len(args) > 0 {
		v = fmt.Sprintf(v, args...)
	}
	return ValueStr{ValueBase: &ValueBase{}, v: []rune(v), s: v}
}

func NewStrByRunes(v []rune) ValueStr {
	return ValueStr{ValueBase: &ValueBase{}, v: v, s: string(v)}
}

func (v ValueStr) GoType() reflect.Type {
	var vv string
	return reflect.TypeOf(vv)
}

func (v ValueStr) ToGoValue() interface{} {
	return v.Value()
}

func (v ValueStr) GetIndex(index int, c *Context) Value {
	chs := []rune(v.v)
	if index < 0 || index >= len(chs) {
		return constUndefined
	}
	return NewStrByRunes(chs[index : index+1])
}

func (v ValueStr) GetMember(name string, c *Context) Value {
	return getMemberByType(c, v, name)
}

func (ValueStr) Type() ValueType {
	return TypeStr
}

func (v ValueStr) Value() string {
	return v.s
}

func (v ValueStr) Runes() []rune {
	return v.v
}

func (v ValueStr) Len() int {
	return len(v.v)
}

func (v ValueStr) IsTrue() bool {
	return len(v.v) > 0
}

func (v ValueStr) CompareTo(other Value, c *Context) CompareResult {
	if v2, ok := other.(ValueStr); ok {
		if v.Value() == v2.Value() {
			return CompareResultEqual
		} else if v.Value() < v2.Value() {
			return CompareResultLess
		} else {
			return CompareResultGreater
		}
	}
	return CompareResultNotEqual
}

func (v ValueStr) ToString(*Context) string {
	return v.s
}

var builtinStrMethods = map[string]ValueCallable{
	"substr": &ValueBuiltinFunction{
		name: "str.substr",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisStr := thisArg.(ValueStr)
			charSeq := thisStr.v
			begin, end := 0, len(charSeq)
			switch len(args) {
			case 0:
			case 2:
				{
					endVal, ok := args[1].(ValueInt)
					if !ok {
						c.RaiseRuntimeError("str.substr arguments 1 must be an integer")
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
						c.RaiseRuntimeError("str.substr arguments 0 must be an integer")
						return nil
					}
					begin = int(beginVal.Value())
					if begin < 0 {
						begin += len(charSeq)
					}
				}
			default:
				c.RaiseRuntimeError("str.substr requires 0~2 arguments")
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
				return NewStr("")
			}
			sliceChars := make([]rune, end-begin)
			for i := begin; i < end; i++ {
				sliceChars[i-begin] = charSeq[i]
			}
			return NewStrByRunes(sliceChars)
		},
	},
	"upper": &ValueBuiltinFunction{
		name: "str.upper",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisStr := thisArg.(ValueStr)
			return NewStr(strings.ToUpper(thisStr.Value()))
		},
	},
	"lower": &ValueBuiltinFunction{
		name: "str.lower",
		body: func(c *Context, thisArg Value, args []Value) Value {
			thisStr := thisArg.(ValueStr)
			return NewStr(strings.ToLower(thisStr.Value()))
		},
	},
	"find": &ValueBuiltinFunction{
		name: "str.find",
		body: func(c *Context, thisArg Value, args []Value) Value {
			this := thisArg.(ValueStr).Value()
			var pattern string
			start := 0
			ok := true
			if pattern, ok = getArgStr(args, 0); !ok {
				c.RaiseRuntimeError("str.find argument 0 pattern error")
				return nil
			}
			if len(args) > 1 {
				if start, ok = getArgInt(args, 1); !ok {
					c.RaiseRuntimeError("str.find argument 1 startPos error")
					return nil
				}
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				c.RaiseRuntimeError("str.find argument 0 pattern error: %s", err)
				return nil
			}
			return NewStr(re.FindString(this[start:]))
		},
	},
	"hash": NewNativeFunction("str.hash", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		seed := uint64(131) // 31 131 1313 13131 131313 etc..
		hash := uint64(0)
		for i := 0; i < len(str); i++ {
			hash = (hash * seed) + uint64(str[i])
		}
		return NewInt(int64(hash & 0x7FFFFFFF))
	}),
	"split": NewNativeFunction("str.split", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		sp := " "
		limit := -1
		switch len(args) {
		case 2:
			limit = int(c.MustInt(args[1], "str.split argument limit"))
			fallthrough
		case 1:
			sp = c.MustStr(args[0], "str.split argument sp")
		default:
			c.RaiseRuntimeError("str.split usage: split(sp[, limit=-1])")
		}
		items := strings.SplitN(str, sp, limit)
		rv := NewArray(len(items))
		for _, item := range items {
			rv.PushBack(NewStr(item))
		}
		return rv
	}),
	"splitr": NewNativeFunction("str.splitr", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		var (
			sp    ValueStr
			limit ValueInt
		)
		EnsureFuncParams(c, "str.splitr", args,
			ArgRuleOptional("sp", TypeStr, &sp, NewStr("\\s+")),
			ArgRuleOptional("limit", TypeInt, &limit, NewInt(-1)),
		)
		if re, err := regexp.Compile(sp.Value()); err != nil {
			c.RaiseRuntimeError("str.splitr: parse regexp %s error %s", sp, err)
			return nil
		} else {
			items := re.Split(str, limit.AsInt())
			rv := NewArray(len(items))
			for _, item := range items {
				rv.PushBack(NewStr(item))
			}
			return rv
		}
	}),
	"lines": NewNativeFunction("str.lines", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		lines := strings.Split(str, "\n")
		rv := NewArray(len(lines))
		for _, l := range lines {
			rv.PushBack(NewStr(l))
		}
		return rv
	}),
	"indexOf": NewNativeFunction("str.indexOf", func(c *Context, thisArg Value, args []Value) Value {
		str := thisArg.(ValueStr).v
		sub := mustGetArgStr(c, "str.indexOf", args, 0)
		offset := 0
		if s, ok := getArgInt(args, 1); ok {
			str = str[s:]
			offset = s
		}
		ss := string(str)
		idx := strings.Index(ss, sub)
		if idx < 0 {
			return NewInt(int64(idx))
		}
		return NewInt(int64(offset + utf8.RuneCountInString(ss[:idx])))
	}),
	"lastIndexOf": NewNativeFunction("str.lastIndexOf", func(c *Context, thisArg Value, args []Value) Value {
		str := thisArg.(ValueStr).v
		sub := mustGetArgStr(c, "str.lastIndexOf", args, 0)
		ss := string(str)
		idx := strings.LastIndex(ss, sub)
		if idx < 0 {
			return NewInt(int64(idx))
		}
		return NewInt(int64(utf8.RuneCountInString(ss[:idx])))
	}),
	"replaceOne": NewNativeFunction("str.replaceOne", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		sub := mustGetArgStr(c, "str.replaceOne", args, 0)
		repl := mustGetArgStr(c, "str.replaceOne", args, 1)
		return NewStr(strings.Replace(str, sub, repl, 1))
	}),
	"replaceAll": NewNativeFunction("str.replaceAll", func(c *Context, thisArg Value, args []Value) Value {
		var sub, repl ValueStr
		EnsureFuncParams(c, "str.replaceAll", args,
			ArgRuleRequired("sub", TypeStr, &sub),
			ArgRuleRequired("repl", TypeStr, &repl),
		)
		str := c.MustStr(thisArg)
		return NewStr(strings.ReplaceAll(str, sub.Value(), repl.Value()))
	}, "sub", "repl"),
	"replace": NewNativeFunction("str.replace", func(c *Context, thisArg Value, args []Value) Value {
		var (
			pattern  ValueStr
			repl     ValueStr
			replFunc ValueCallable
			replType int
		)
		EnsureFuncParams(c, "str.replace", args,
			ArgRuleRequired("pattern", TypeStr, &pattern),
			ArgRuleOneOf("repl", []ValueType{TypeStr, TypeCallable}, []interface{}{&repl, &replFunc}, &replType, nil, nil),
		)
		p, err := regexp.Compile(pattern.Value())
		if err != nil {
			c.RaiseRuntimeError("str.replace: regexp %s is invalid: %s", pattern.Value(), err)
		}
		switch replType {
		case 1: // By Callable
			return NewStr(p.ReplaceAllStringFunc(c.MustStr(thisArg), func(r string) string {
				c.Invoke(replFunc, nil, Args(NewStr(r)))
				return c.RetVal.ToString(c)
			}))
		default:
			return NewStr(p.ReplaceAllString(c.MustStr(thisArg), repl.Value()))
		}
	}, "pattern", "repl"),
	"trim": NewNativeFunction("str.trim", func(c *Context, this Value, args []Value) Value {
		str := c.MustStr(this)
		return NewStr(strings.TrimSpace(str))
	}),
	"startsWith": NewNativeFunction("str.startsWith", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		prefix := mustGetArgStr(c, "str.startsWith", args, 0)
		return NewBool(strings.HasPrefix(str, prefix))
	}, "prefix"),
	"endsWith": NewNativeFunction("str.endsWith", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		prefix := mustGetArgStr(c, "str.endsWith", args, 0)
		return NewBool(strings.HasSuffix(str, prefix))
	}, "suffix"),
	"contains": NewNativeFunction("str.contains", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		sub := mustGetArgStr(c, "str.contains", args, 0)
		return NewBool(strings.Contains(str, sub))
	}, "sub"),
	"__next__": NewNativeFunction("__next__", func(c *Context, this Value, args []Value) Value {
		s := this.(ValueStr)
		if len(s.v) == 1 {
			return NewStrByRunes([]rune{s.v[0] + 1})
		} else {
			return constUndefined
		}
	}),
	"decodeHex": NewNativeFunction("str.decodeHex", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		bs, err := hex.DecodeString(str)
		if err != nil {
			c.RaiseRuntimeError("str.decodeHex decode failed: %s", err)
		}
		return NewBytes(bs)
	}),
	"decodeBase64": NewNativeFunction("str.decodeBase64", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		bs, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			c.RaiseRuntimeError("str.decodeBase64 decode failed: %s", err)
		}
		return NewBytes(bs)
	}),
	"fillParams": func() *ValueBuiltinFunction {
		re := regexp.MustCompile(`\{\w+\}`)
		return NewNativeFunction("str.fillParams", func(c *Context, thisArg Value, args []Value) Value {
			str := c.MustStr(thisArg)
			var params ValueObject
			EnsureFuncParams(c, "str.fillParams", args, ArgRuleRequired("params", TypeObject, &params))
			return NewStr(re.ReplaceAllStringFunc(str, func(s string) string {
				return params.GetMember(s[1:len(s)-1], c).ToString(c)
			}))
		})
	}(),
	"code": NewNativeFunction("str.code", func(c *Context, thisArg Value, args []Value) Value {
		this, isStr := Unbound(thisArg).(ValueStr)
		if !isStr {
			c.RaiseRuntimeError("this is not a string!")
		}
		var argIndex ValueInt
		EnsureFuncParams(c, "str.code", args, ArgRuleOptional("index", TypeInt, &argIndex, NewInt(0)))
		index := argIndex.AsInt()
		if index < 0 {
			index += this.Len()
		}
		if index < 0 || index >= this.Len() {
			c.RaiseRuntimeError("code index %d out of range [%d, %d)", index, -this.Len(), this.Len())
		}
		return NewInt(int64(this.v[index]))
	}),
	"codes": NewNativeFunction("str.codes", func(c *Context, thisArg Value, args []Value) Value {
		this, isStr := Unbound(thisArg).(ValueStr)
		if !isStr {
			c.RaiseRuntimeError("this is not a string!")
		}
		rv := NewArray(this.Len())
		for _, c := range this.v {
			rv.PushBack(NewInt(int64(c)))
		}
		return rv
	}),
	"width": NewNativeFunction("str.width", func(c *Context, thisArg Value, args []Value) Value {
		str := c.MustStr(thisArg)
		return NewInt(int64(runewidth.StringWidth(str)))
	}),
}

var (
	emptyStr = ValueStr{ValueBase: &ValueBase{}, v: []rune{}, s: ""}
)

func init() {
	addMembersAndStatics(TypeStr, builtinStrMethods)
	TypeStr.Statics.Store("fromCodes", NewNativeFunction("Str.fromCodes", func(c *Context, this Value, args []Value) Value {
		if len(args) < 1 {
			return NewStr("")
		}
		s := make([]rune, 0, len(args))
		for _, a := range args {
			code := c.MustInt(a)
			s = append(s, rune(code))
		}
		return NewStrByRunes(s)
	}))
}
