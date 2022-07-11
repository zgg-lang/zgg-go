package runtime

import "reflect"

type funcEnv struct {
	cur, funcRoot, root *contextFrame
}

type ValueFunc struct {
	*ValueBase
	Name       string
	Args       []string
	env        *funcEnv
	ExpandLast bool
	This       Value
	Body       IEval
	BelongType ValueType
}

func NewFunc(name string, args []string, expandLast bool, body IEval) *ValueFunc {
	f := &ValueFunc{
		ValueBase:  &ValueBase{},
		Name:       name,
		Args:       args,
		ExpandLast: expandLast,
		Body:       body,
	}
	return f
}

func (v *ValueFunc) CloneWithEnv(c *Context) *ValueFunc {
	newFunc := *v
	cur, fr, r := c.CloneFrames()
	newFunc.env = &funcEnv{cur: cur, funcRoot: fr, root: r}
	return &newFunc
}

func (v *ValueFunc) GetIndex(index int, c *Context) Value {
	return constUndefined
}

func (v *ValueFunc) GetMember(name string, c *Context) Value {
	switch name {
	case "__name__":
		if v.Name != "" {
			return NewStr(v.Name)
		}
		return NewStr("anoymous")
	case "__args__":
		{
			args := NewArray(len(v.Args))
			for _, arg := range v.Args {
				args.PushBack(NewStr(arg))
			}
			return args
		}
	}
	return getCallableMember(v, name, c)
}

func (v *ValueFunc) Type() ValueType {
	return TypeFunc
}

func (v *ValueFunc) IsTrue() bool {
	return true
}

func (v *ValueFunc) CompareTo(other Value, c *Context) CompareResult {
	if f, ok := other.(*ValueFunc); ok && f == v {
		return CompareResultEqual
	}
	return CompareResultNotEqual
}

func (v *ValueFunc) ToString(*Context) string {
	if v.Name == "" {
		return "<anonymous function>"
	}
	return "<func " + v.Name + ">"
}

func (v *ValueFunc) GoType() reflect.Type {
	return reflect.TypeOf(nil)
}

func (v *ValueFunc) ToGoValue() interface{} {
	return nil
}

func (v *ValueFunc) GetName() string {
	if v.Name == "" {
		return "<anonymous function>"
	}
	return v.Name
}

func (v *ValueFunc) GetArgNames(*Context) []string {
	return v.Args
}

func (v *ValueFunc) Invoke(c *Context, thisArg Value, args []Value) {
	if e := v.env; e != nil {
		cur := c.curFrame
		froot := c.funcRootFrame
		root := c.rootFrame
		c.curFrame = e.cur
		c.funcRootFrame = e.funcRoot
		c.rootFrame = e.root
		defer func() {
			c.curFrame = cur
			c.funcRootFrame = froot
			c.rootFrame = root
		}()
	}
	c.PushFuncStack(v.GetName())
	defer c.PopStack()
	if v.This != nil {
		thisArg = v.This
	}
	if thisObj, isObj := thisArg.(ValueObject); isObj {
		if v.BelongType != nil {
			super := thisObj.Super(v.BelongType)
			c.ForceSetLocalValue("super", super)
			if thisObj.this != nil {
				thisArg = *(thisObj.this)
			}
		}
	}
	if thisArg == nil {
		thisArg = constUndefined
	}
	c.ForceSetLocalValue("this", thisArg)
	argumentsValue := NewArray(len(args))
	for _, arg := range args {
		argumentsValue.PushBack(arg)
	}
	c.SetLocalValue("arguments", argumentsValue)
	n := len(v.Args)
	inputN := len(args)
	if v.ExpandLast {
		n--
	}
	for i := 0; i < n; i++ {
		var argVal Value = constUndefined
		if i < inputN {
			argVal = args[i]
		}
		c.ForceSetLocalValue(v.Args[i], argVal)
	}
	if v.ExpandLast {
		if inputN >= len(v.Args) {
			argVal := NewArray(inputN - n)
			for i := n; i < inputN; i++ {
				argVal.PushBack(args[i])
			}
			c.ForceSetLocalValue(v.Args[n], argVal)
		} else {
			c.ForceSetLocalValue(v.Args[n], NewArray(0))
		}
	}
	v.Body.Eval(c)
}

func NoArgs() []Value {
	return []Value{}
}

func Args(v ...Value) func() []Value {
	return func() []Value { return v }
}
