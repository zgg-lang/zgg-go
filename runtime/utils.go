package runtime

import (
	"reflect"
	"regexp"
	"strings"
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
	if rv, is := res.(Value); is {
		return rv
	} else {
		return NewGoValue(res)
	}
}

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
	return getExtMember(v, name, c)
}
