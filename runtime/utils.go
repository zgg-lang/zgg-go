package runtime

import (
	"reflect"
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
	return arg.Type().IsSubOf(typ)
}

func storeValue(storeTo interface{}, v Value) {
	rv := reflect.ValueOf(v)
	if _, ok := v.(ValueBoundMethod); ok {
		if _, ok := storeTo.(*ValueBoundMethod); !ok {
			rv = reflect.ValueOf(Unbound(v))
		}
	}
	reflect.ValueOf(storeTo).Elem().Set(rv)
}

type ArgRuleRequired struct {
	ArgName      string
	ExpectedType ValueType
	StoreTo      interface{}
}

func (r ArgRuleRequired) name() string {
	return r.ArgName
}

func (r ArgRuleRequired) expectedTypeName() string {
	return r.ExpectedType.GetName()
}

func (r ArgRuleRequired) allowed(c *Context, args []Value, i int) bool {
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

type ArgRuleOptional struct {
	ArgName      string
	ExpectedType ValueType
	StoreTo      interface{}
	DefaultValue Value
}

func (r ArgRuleOptional) name() string {
	return r.ArgName + "?"
}

func (r ArgRuleOptional) expectedTypeName() string {
	return r.ExpectedType.GetName()
}

func (r ArgRuleOptional) allowed(c *Context, args []Value, i int) bool {
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

type ArgRuleOneOf struct {
	ArgName       string
	ExpectedTypes []ValueType
	StoreTos      []interface{}
	Selected      *int
	DefaultStore  interface{}
	DefaultValue  Value
}

func (r ArgRuleOneOf) name() string {
	return r.ArgName
}

func (r ArgRuleOneOf) expectedTypeName() string {
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

func (r ArgRuleOneOf) allowed(c *Context, args []Value, i int) bool {
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
		if _, ok := rules[i].(ArgRuleRequired); ok {
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
			c.OnRuntimeError("%s) requires %d argument(s), but got %d", funcName, len(rules), len(args))
		} else {
			c.OnRuntimeError("%s) requires at least %d argument(s), but got %d", funcName, minArgs, len(args))
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
			c.OnRuntimeError("%s) arg %s should be a(n) %s", funcName, rule.name(), rule.expectedTypeName())
			return
		}
	}
}
