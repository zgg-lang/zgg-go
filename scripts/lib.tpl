package main

import (
{{ if .Exported }}
        pkg "{{.ImportPath}}"
{{end}}
        "reflect"
)

func New() map[string]interface{} {
	return map[string]interface{} {
        // Functions
        {{ range .Funcs }} "{{.}}": reflect.ValueOf(pkg.{{.}}),
        {{ end }}
        // Consts

        {{ range .Consts }} "{{.}}": reflect.ValueOf({{ if index $.TypeMapping . }}{{index $.TypeMapping .}}({{end}}pkg.{{.}}{{ if index $.TypeMapping . }}){{end}}),
        {{ end }}
        // Variables
        
        {{ range .Vars }} "{{.}}": reflect.ValueOf(&pkg.{{.}}),
        {{ end }}
        // Non interfaces

        {{ range .NonInterfaces }} "{{.}}": reflect.TypeOf((*pkg.{{.}})(nil)).Elem(),
        {{ end }}
        }
}