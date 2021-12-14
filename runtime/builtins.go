package runtime

var builtins = map[string]Value{}

func init() {
	for name, value := range builtinFunctions {
		builtins[name] = value
	}
	for name, value := range builtinTypes {
		builtins[name] = value
	}
}
