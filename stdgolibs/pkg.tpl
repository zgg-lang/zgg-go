package stdgolibs

import (
{{ if .Exported }}
        pkg "{{.ImportPath}}"
{{end}}
        "reflect"
)

func init() {
	registerValues("{{.ImportPath}}", map[string]reflect.Value {
        // Functions
        {{ range .Funcs }} "{{.}}": reflect.ValueOf(pkg.{{.}}),
        {{ end }}
        // Consts

        {{ range .Consts }} "{{.}}": reflect.ValueOf({{ if index $.TypeMapping . }}{{index $.TypeMapping .}}({{end}}pkg.{{.}}{{ if index $.TypeMapping . }}){{end}}),
        {{ end }}
        // Variables
        
        {{ range .Vars }} "{{.}}": reflect.ValueOf(&pkg.{{.}}),
        {{ end }}
	})
	registerTypes("{{.ImportPath}}", map[string]reflect.Type {
        // Non interfaces

        {{ range .NonInterfaces }} "{{.}}": reflect.TypeOf((*pkg.{{.}})(nil)).Elem(),
        {{ end }}
        })
}