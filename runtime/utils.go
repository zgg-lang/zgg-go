package runtime

import (
	"math"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

func Nullish(v Value) bool {
	switch v.(type) {
	case ValueNil:
		return true
	case ValueUndefined:
		return true
	}
	return false
}

type ArgRule interface {
	name() string
	expectedTypeName() string
	allowed(*Context, []Value, int) bool
}

func argTypeMatched(c *Context, arg Value, typ ValueType) bool {
	if typ == TypeAny {
		return true
	}
	if typ == TypeCallable {
		return c.IsCallable(arg)
	}
	if _, is := arg.(ValueInt); is && typ == TypeFloat {
		return true
	}
	return arg.Type().IsSubOf(typ)
}

func storeValue(storeTo interface{}, v Value) {
	if _, ok := v.(ValueBoundMethod); ok {
		if _, ok := storeTo.(*ValueBoundMethod); !ok {
			v = Unbound(v)
		}
	}
	if toFloat, is := storeTo.(*ValueFloat); is {
		if vi, isInt := v.(ValueInt); isInt {
			*toFloat = NewFloat(float64(vi.Value()))
			return
		}
	}
	rv := reflect.ValueOf(v)
	reflect.ValueOf(storeTo).Elem().Set(rv)
}

type argRuleRequired struct {
	ArgName      string
	ExpectedType ValueType
	StoreTo      interface{}
}

func ArgRuleRequired(name string, expectedType ValueType, storeTo interface{}) ArgRule {
	return argRuleRequired{ArgName: name, ExpectedType: expectedType, StoreTo: storeTo}
}

func (r argRuleRequired) name() string {
	return r.ArgName
}

func (r argRuleRequired) expectedTypeName() string {
	return r.ExpectedType.GetName()
}

func (r argRuleRequired) allowed(c *Context, args []Value, i int) bool {
	if i >= len(args) {
		return false
	}
	arg := args[i]
	if !argTypeMatched(c, arg, r.ExpectedType) {
		return false
	}
	if argVal, ok := arg.(ValueBoundMethod); ok {
		arg = argVal.Value
	}
	storeValue(r.StoreTo, arg)
	return true
}

type argRuleOptional struct {
	ArgName      string
	ExpectedType ValueType
	StoreTo      interface{}
	DefaultValue Value
}

func ArgRuleOptional(name string, expectedType ValueType, storeTo interface{}, defaultValue Value) ArgRule {
	return argRuleOptional{ArgName: name, ExpectedType: expectedType, StoreTo: storeTo, DefaultValue: defaultValue}
}

func (r argRuleOptional) name() string {
	return r.ArgName + "?"
}

func (r argRuleOptional) expectedTypeName() string {
	return r.ExpectedType.GetName()
}

func (r argRuleOptional) allowed(c *Context, args []Value, i int) bool {
	defaultValue := reflect.ValueOf(r.DefaultValue)
	if i >= len(args) {
		if defaultValue.IsValid() {
			storeValue(r.StoreTo, r.DefaultValue)
		}
		return true
	}
	arg := args[i]
	if !argTypeMatched(c, arg, r.ExpectedType) {
		if _, ok := arg.(ValueUndefined); ok {
			if defaultValue.IsValid() {
				storeValue(r.StoreTo, r.DefaultValue)
			}
			return true
		}
		return false
	}
	storeValue(r.StoreTo, arg)
	return true
}

type argRuleOneOf struct {
	ArgName       string
	ExpectedTypes []ValueType
	StoreTos      []interface{}
	Selected      *int
	DefaultStore  interface{}
	DefaultValue  Value
}

func ArgRuleOneOf(
	argName string,
	expectedTypes []ValueType,
	storeTos []interface{},
	selected *int,
	defaultStore interface{},
	defaultValue Value,
) ArgRule {
	return argRuleOneOf{
		ArgName:       argName,
		ExpectedTypes: expectedTypes,
		StoreTos:      storeTos,
		Selected:      selected,
		DefaultStore:  defaultStore,
		DefaultValue:  defaultValue,
	}
}

func (r argRuleOneOf) name() string {
	return r.ArgName
}

func (r argRuleOneOf) expectedTypeName() string {
	var b strings.Builder
	b.WriteString("(")
	for i, e := range r.ExpectedTypes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.GetName())
	}
	b.WriteString(")")
	return b.String()
}

func (r argRuleOneOf) allowed(c *Context, args []Value, i int) bool {
	if len(r.ExpectedTypes) < 1 || len(r.ExpectedTypes) != len(r.StoreTos) {
		panic("ArgRuleOneOf: len(ExpectedTypes) must > 0 && len(ExpectedTypes) must == len(StoreTos)")
	}
	if i >= len(args) {
		if r.DefaultStore != nil {
			*r.Selected = -1
			storeValue(r.DefaultStore, r.DefaultValue)
			return true
		}
		return false
	}
	arg := args[i]
	for j, expType := range r.ExpectedTypes {
		if argTypeMatched(c, arg, expType) {
			storeValue(r.StoreTos[j], arg)
			*r.Selected = j
			return true
		}
	}
	return false
}

func EnsureFuncParams(c *Context, funcName string, args []Value, rules ...ArgRule) {
	minArgs := 0
	for i := 0; i < len(rules); i++ {
		if _, ok := rules[i].(argRuleRequired); ok {
			minArgs++
		} else {
			break
		}
	}
	if len(args) < minArgs {
		funcName += "("
		for i, r := range rules {
			if i > 0 {
				funcName += ", "
			}
			funcName += r.name()
		}
		if minArgs == len(rules) {
			c.RaiseRuntimeError("%s) requires %d argument(s), but got %d", funcName, len(rules), len(args))
		} else {
			c.RaiseRuntimeError("%s) requires at least %d argument(s), but got %d", funcName, minArgs, len(args))
		}
		return
	}
	for i, rule := range rules {
		if !rule.allowed(c, args, i) {
			funcName += "("
			for i, r := range rules {
				if i > 0 {
					funcName += ", "
				}
				funcName += r.name()
			}
			c.RaiseRuntimeError("%s) arg %s should be a(n) %s", funcName, rule.name(), rule.expectedTypeName())
			return
		}
	}
}

type oneOfHelper struct {
	argName         string
	types           []ValueType
	callbacks       []func(Value)
	defaultCallback func()
}

func NewOneOfHelper(name string) *oneOfHelper {
	return &oneOfHelper{argName: name}
}

func (h *oneOfHelper) On(t ValueType, f func(Value)) *oneOfHelper {
	h.types = append(h.types, t)
	h.callbacks = append(h.callbacks, f)
	return h
}

func (h *oneOfHelper) Default(f func()) *oneOfHelper {
	h.defaultCallback = f
	return h
}

func (h *oneOfHelper) name() string {
	return h.argName
}

func (h *oneOfHelper) expectedTypeName() string {
	var b strings.Builder
	b.WriteString("(")
	for i, e := range h.types {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.GetName())
	}
	b.WriteString(")")
	return b.String()
}

func (h *oneOfHelper) allowed(c *Context, args []Value, i int) bool {
	if len(h.types) < 1 || len(h.types) != len(h.callbacks) {
		panic("ArgRuleOneOf: len(ExpectedTypes) must > 0 && len(ExpectedTypes) must == len(StoreTos)")
	}
	if i >= len(args) {
		if h.defaultCallback != nil {
			h.defaultCallback()
			return true
		}
		return false
	}
	arg := args[i]
	for j, expType := range h.types {
		if argTypeMatched(c, arg, expType) {
			h.callbacks[j](arg)
			return true
		}
	}
	return false
}

var (
	pathIdField  = regexp.MustCompile(`^\.([a-zA-Z_\$][a-zA-Z\d_\$]*)`)
	pathStrField = regexp.MustCompile(`^\.'(([^\\]|\\[uU][0-9a-fA-F]{4}|\\[xX][0-9a-fA-F]{2}|\\[^xXuU])*)'`)
	escapedChar  = regexp.MustCompile(`\\[uU][0-9a-fA-F]{4}|\\[xX][0-9a-fA-F]{2}|\\[^xXuU]`)
	pathIndex    = regexp.MustCompile(`^\[(\-?\d+)\]`)
)

func GetValueByPath(c *Context, v Value, path string) Value {
	res, err := jsonPathLookup(c, v, path)
	if err != nil {
		c.RaiseRuntimeError("find value by path %s error %+v", path, err)
	}
	switch rv := res.(type) {
	case Value:
		return rv
	case []any:
		retArr := NewArray(len(rv))
		for _, v := range rv {
			if vv, is := v.(Value); is {
				retArr.PushBack(vv)
			} else {
				retArr.PushBack(NewGoValue(v))
			}
		}
		return retArr
	default:
		return NewGoValue(res)
	}
}

var (
	commonMembers         map[string]*ValueBuiltinFunction
	commonMembersInitOnce sync.Once
)

func getMemberByType(c *Context, v Value, name string) Value {
	t := v.Type()
	if member := t.findMember(name); member != nil {
		return makeMember(v, member, c)
	}
	if getAttr, ok := c.GetCallable(t.findMember("__getAttr__")); ok {
		c.Invoke(getAttr, v, func() []Value { return []Value{NewStr(name)} })
		if _, isUndefiend := c.RetVal.(ValueUndefined); !isUndefiend {
			return makeMember(v, c.RetVal, c)
		}
	}
	commonMembersInitOnce.Do(func() {
		commonMembers = map[string]*ValueBuiltinFunction{
			"must": NewNativeFunction("must", func(c *Context, this Value, args []Value) Value {
				thisArg := Args(this)
				for _, a := range args {
					if callable, is := c.GetCallable(a); is {
						c.Invoke(callable, nil, thisArg)
						if !c.RetVal.IsTrue() {
							c.RaiseRuntimeError("assert 'must' failed")
						}
					} else if !c.ValuesEqual(this, a) {
						c.RaiseRuntimeError("assert 'must' failed")
					}
				}
				return this
			}),
			"notNil": NewNativeFunction("notNil", func(c *Context, this Value, args []Value) Value {
				switch this.(type) {
				case ValueNil:
				case ValueUndefined:
				default:
					return this
				}
				c.RaiseRuntimeError("assert not nil/undefined failed")
				return nil
			}),
		}
	})
	if f, found := commonMembers[name]; found {
		return makeMember(v, f, c)
	}
	return getExtMember(v, name, c)
}

type iteratorInfo struct {
	nextFn  func() Value
	closeFn func()
	closed  bool
}

var (
	iteratorType ValueType
	iteratorInit sync.Once
)

func MakeIterator(c *Context, nextFn func() Value, closeFn func()) ValueObject {
	iteratorInit.Do(func() {
		endRet := NewArrayByValues(constUndefined, NewBool(false))
		iteratorType = NewClassBuilder("iterator").
			Constructor(func(c *Context, this ValueObject, args []Value) {
				this.Reserved = &iteratorInfo{
					nextFn:  args[0].ToGoValue(c).(func() Value),
					closeFn: args[1].ToGoValue(c).(func()),
				}
			}).
			Method("__call__", func(c *Context, this ValueObject, args []Value) Value {
				info := this.Reserved.(*iteratorInfo)
				if info.closed {
					return endRet
				}
				v := info.nextFn()
				if v != nil {
					return NewArrayByValues(v, NewBool(true))
				} else {
					if info.closeFn != nil && !info.closed {
						info.closeFn()
					}
					info.closed = true
					return endRet
				}
			}).
			Method("close", func(c *Context, this ValueObject, args []Value) Value {
				info := this.Reserved.(*iteratorInfo)
				if info.closeFn != nil && !info.closed {
					info.closeFn()
				}
				info.closed = true
				return constUndefined
			}).
			Build()
	})
	rv := NewObject()
	rv.SetMember("__iter__", NewNativeFunction("__iter__", func(c *Context, _ Value, _ []Value) Value {
		return NewObjectAndInit(iteratorType, c, NewGoValue(nextFn), NewGoValue(closeFn))
	}), nil)
	return rv
}

func fixSliceRange(begin, end, size int64) (int64, int64) {
	if end == math.MaxInt64 {
		end = size
	}
	if begin < 0 {
		begin += size
	}
	if begin < 0 {
		begin = 0
	} else if begin > size {
		begin = size
	}
	if end < 0 {
		end += size
	}
	if end < 0 {
		end = 0
	} else if end > size {
		end = size
	}
	return begin, end
}
