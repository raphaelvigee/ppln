package main

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/imports"
	"io"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

var ifaceInTpl = template.Must(template.New("").Funcs(funcs).Parse(`
type NodeIn{{.N}}[T any] interface {
	Node

	_input{{.N}}(T)
}
`))

var ifaceOutTpl = template.Must(template.New("").Funcs(funcs).Parse(`
type NodeOut{{.N}}[T any] interface {
	Node

	_output{{.N}}(T)
}
`))

var nodeTpl = template.Must(template.New("").Funcs(funcs).Parse(`
type Node{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] interface {
	Node

	{{range $i := loop .InCount}}
		 NodeIn{{$i}}[I{{$i}}]
	{{- end}}

	{{range $i := loop .OutCount}}
		 NodeOut{{$i}}[O{{$i}}]
	{{- end}}
}

type FuncNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] struct {
	Node{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]

	Func func({{.InputParametersDef}}) ({{.OutputReturnDef}})
}

func (f FuncNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Do(inputs []any) []any {
	{{- range $idx, $i := loop .InCount}}
		 i{{$i}} := inputs[{{$idx}}].(I{{$i}})
	{{- end}}

	{{if gt .OutCount 0}}
		{{.OutputVars}} := f.Func({{.InputVars}})
	
		return []any{ {{.OutputVars}} }
	{{- else}}
		f.Func({{.InputVars}})

		return nil
	{{- end}}
}
`))

var streamNodeTpl = template.Must(template.New("").Funcs(funcs).Parse(`
type StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] interface {
	StreamNode

	{{range $i := loop .InCount}}
		 NodeIn{{$i}}[I{{$i}}]
	{{- end}}

	{{range $i := loop .OutCount}}
		 NodeOut{{$i}}[O{{$i}}]
	{{- end}}
}

type FuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] struct {
	StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]

	Func func(
	{{range $i := loop .InCount -}}
		 I{{$i}},
	{{end -}}
	{{ range $idx, $i := loop .OutCount -}}
		 func (v O{{$i}}),
	{{ end -}}
	)
}

func (f FuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Do(inputs []any, emit func(i int, v any)) {
	{{- range $idx, $i := loop .InCount}}
		 i{{$i}} := inputs[{{$idx}}].(I{{$i}})
	{{- end}}

	f.Func(
	{{range $i := loop .InCount -}}
		 i{{$i}},
	{{end}}
	{{- range $idx, $i := loop .OutCount -}}
		 func (v O{{$i}}) {
			emit({{$idx}}, v)
		 },
	{{ end}}
	)
}
`))

var takeTpl = template.Must(template.New("").Funcs(funcs).Parse(`
func Take{{.N}}[T any](n NodeOut{{.N}}[T]) NodeOut1[T] {
	panic("TODO")
}

func Pipeline{{.N}}[{{g_generics "T" .N true}}]({{g_params "from" "NodeOut1[T#]" .N }}, to interface{ {{range $i := loop .N}} NodeIn{{$i}}[T{{$i}}]; {{end}} }) []Edge {
	return []Edge{
	{{- range $idx, $i := loop .N}}
		{
			From: from{{$i}},
			To:   to,
		},
	{{- end}}
	}
}
`))

func loop(n int, zero bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		start := 1
		if zero {
			start = 0
		}
		for i := start; i <= n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func genGenericsList(out string, prefix string, n int, withType bool) string {
	for i := range loop(n, false) {
		if out != "" {
			out += ", "
		}
		out += fmt.Sprintf("%v%v", prefix, i)
		if withType {
			out += " any"
		}
	}

	return out
}

func genParametersList(out string, varName, typePrefix string, n int) string {
	for i := range loop(n, false) {
		if out != "" {
			out += ", "
		}
		if varName == "" {
			out += fmt.Sprintf("%v%v", typePrefix, i)
		} else if typePrefix == "" {
			out += fmt.Sprintf("%v%v", varName, i)
		} else {
			out += fmt.Sprintf("%v%v %v%v", varName, i, typePrefix, i)
		}
	}

	return out
}

func genParametersList2(out string, varName, typeTemplate string, n int) string {
	for i := range loop(n, false) {
		if out != "" {
			out += ", "
		}

		t := strings.ReplaceAll(typeTemplate, "#", fmt.Sprint(i))

		if varName == "" {
			out += t
		} else {
			out += fmt.Sprintf("%v%v %v", varName, i, t)
		}
	}

	return out
}

var funcs = template.FuncMap{
	"loop": func(c int) []int {
		return slices.Collect(loop(c, false))
	},
	"g_generics": func(prefix string, n int, withType bool) string {
		return genGenericsList("", prefix, n, withType)
	},
	"g_params": func(varPrefix, typeTemplate string, n int) string {
		return genParametersList2("", varPrefix, typeTemplate, n)
	},
}

type Config struct {
	Package                 string
	InputCount, OutputCount int
}

func gen(w io.Writer, c Config) error {
	_, err := fmt.Fprintf(w, "package %v\n", c.Package)
	if err != nil {
		return err
	}

	ic := c.InputCount
	oc := c.OutputCount

	for i := range loop(ic, false) {
		err := ifaceInTpl.Execute(w, map[string]interface{}{
			"N": i,
		})
		if err != nil {
			return err
		}
	}

	for i := range loop(oc, false) {
		err := ifaceOutTpl.Execute(w, map[string]interface{}{
			"N": i,
		})
		if err != nil {
			return err
		}

		err = takeTpl.Execute(w, map[string]interface{}{
			"N": i,
		})
		if err != nil {
			return err
		}
	}

	for i := range loop(ic, true) {
		for o := range loop(oc, true) {
			if i == 0 && o == 0 {
				continue
			}

			genericsTypeDef := genGenericsList("", "I", i, true)
			genericsTypeDef = genGenericsList(genericsTypeDef, "O", o, true)

			genericsTypeRef := genGenericsList("", "I", i, false)
			genericsTypeRef = genGenericsList(genericsTypeRef, "O", o, false)

			{
				inputParametersDef := genParametersList("", "v", "I", i)
				outputReturnsDef := genParametersList("", "", "O", o)

				inputVars := genParametersList("", "i", "", i)
				outputVars := genParametersList("", "v", "", o)

				err := nodeTpl.Execute(w, map[string]interface{}{
					"InCount":            i,
					"OutCount":           o,
					"GenericsTypeDef":    genericsTypeDef,
					"GenericsTypeRef":    genericsTypeRef,
					"InputParametersDef": inputParametersDef,
					"InputVars":          inputVars,
					"OutputReturnDef":    outputReturnsDef,
					"OutputVars":         outputVars,
				})
				if err != nil {
					return err
				}
			}

			{
				err := streamNodeTpl.Execute(w, map[string]interface{}{
					"InCount":         i,
					"OutCount":        o,
					"GenericsTypeDef": genericsTypeDef,
					"GenericsTypeRef": genericsTypeRef,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func main() {
	out := os.Args[1]

	var buf bytes.Buffer

	err := gen(&buf, Config{
		Package:     "ppln",
		InputCount:  3,
		OutputCount: 3,
	})
	if err != nil {
		panic(err)
	}

	formatted, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println(err)
		formatted = buf.Bytes()
	}

	err = os.MkdirAll(filepath.Dir(out), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(out, formatted, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
