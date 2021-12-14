package runtime

type ValueCallable interface {
	Value
	GetName() string
	GetArgNames() []string
	Invoke(*Context, Value, []Value)
}

func getCallableMember(v ValueCallable, name string, c *Context) Value {
	if member, found := builtinCallableMembers[name]; found {
		return makeMember(v, member)
	}
	return getExtMember(v, name, c)
}

var builtinCallableMembers = map[string]Value{}
