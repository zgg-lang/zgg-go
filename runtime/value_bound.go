package runtime

import "reflect"

type ValueBoundMethod struct {
	*ValueBase
	Value ValueCallable
}

func NewBoundMethod(owner Value, value ValueCallable) ValueBoundMethod {
	rv := ValueBoundMethod{ValueBase: new(ValueBase), Value: Unbound(value).(ValueCallable)}
	rv.SetOwner(owner)
	return rv
}

func (v ValueBoundMethod) IsTrue() bool {
	return v.Value.IsTrue()
}

func (v ValueBoundMethod) CompareTo(other Value, c *Context) CompareResult {
	return v.Value.CompareTo(other, c)
}

func (v ValueBoundMethod) GetMember(name string, c *Context) Value {
	return v.Value.GetMember(name, c)
}

func (v ValueBoundMethod) GetIndex(index int, c *Context) Value {
	return v.Value.GetIndex(index, c)
}

func (v ValueBoundMethod) Type() ValueType {
	return v.Value.Type()
}

func (v ValueBoundMethod) GoType() reflect.Type {
	return v.Value.GoType()
}

func (v ValueBoundMethod) ToGoValue(c *Context) interface{} {
	return v.Value.ToGoValue(c)
}

func (v ValueBoundMethod) ToString(c *Context) string {
	return v.Value.ToString(c)
}

func (v ValueBoundMethod) GetName() string {
	return v.Value.GetName()
}

func (v ValueBoundMethod) GetArgNames(c *Context) []string {
	return v.Value.GetArgNames(c)
}

func (v ValueBoundMethod) Invoke(c *Context, thisVal Value, args []Value) {
	v.Value.Invoke(c, v.Owner, args)
}

func Unbound(v Value) Value {
	for v != nil {
		if bv, ok := v.(ValueBoundMethod); ok {
			v = bv.Value
		} else {
			break
		}
	}
	return v
}
