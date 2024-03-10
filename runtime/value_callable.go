package runtime

type ValueCallable interface {
	Value
	GetName() string
	GetArgNames(*Context) []string
	Invoke(*Context, Value, []Value)
}

func getCallableMember(v ValueCallable, name string, c *Context) Value {
	if member, found := builtinCallableMembers[name]; found {
		return makeMember(v, member, c)
	}
	return getMemberByType(c, v, name)
}

var builtinCallableMembers = map[string]Value{}
