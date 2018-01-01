package conduit

import (
	tmpl "text/template"
)

var template = tmpl.Must(tmpl.New("template").Parse(`
{{define "package"}}
package {{.}}
{{end}}

{{define "imports"}}
{{range .}}
import {{.}}
{{end}}
{{end}}

{{define "stage"}}
// generated from {{.Name}}(in {{.InputType}}) (out {{.OutputType}})
func {{.Name}}Stage(inc <-chan {{.InputType}}, cancel <-chan struct{}) <-chan {{.OutputType}} {
	ouc := make(chan {{.OutputType}})
	go func() {
		defer close(ouc)
		for in := range inc {
			ouv := {{.Name}}(in)
			select {
			case <-cancel:
				return
			case ouc <- ouv:
			}
		}
	}()
	return ouc
}
{{end}}

{{define "source"}}
// generated from {{.Name}}() (out {{.OutputType}})
func {{.Name}}Source(cancel <-chan struct{}) <-chan {{.OutputType}} {
	ouc := make(chan {{.OutputType}})
	go func() {
		defer close(ouc)
		for {
			ouv := {{.Name}}()
			select {
			case <-cancel:
				return
			case ouc <- ouv:
			}
		}
	}()
	return ouc
}
{{end}}

{{define "sink"}}
// generated from {{.Name}}(in {{.InputType}}) ()
func {{.Name}}Sink(inc <- chan {{.InputType}}, cancel <-chan struct{}) {
	go func() {
		for {
			select {
			case <-cancel:
				return
			default:
			}
			in, ok := <- inc
			if ok == false {
				return
			}
			{{.Name}}(in)
		}
	}()
}
{{end}}`))
