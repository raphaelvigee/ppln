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
{{- if gt .N 0}}
type NodeIn{{.N}}[T any] interface {
	_input{{.N}}(T)
}
{{- end}}

type NodeHas{{.N}}In interface {
	_input_layout({{g_params "" "any" .N }})
}
`))

var ifaceOutTpl = template.Must(template.New("").Funcs(funcs).Parse(`
{{- if gt .N 0}}
type NodeOut{{.N}}[T any] interface {
	_output{{.N}}(T)
}
{{- end}}

type NodeHas{{.N}}Out interface {
	_out_layout({{g_params "" "any" .N }})
}
`))

var nodeTpl = template.Must(template.New("").Funcs(funcs).Parse(` 
type FuncNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] func({{.InputParametersDef}}) ({{.OutputReturnDef}})

func NewFuncNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}](f FuncNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}] {
	return NewFuncStreamNode{{.InCount}}x{{.OutCount}}(func(
		{{range $i := loop .InCount -}}
			 i{{$i}} I{{$i}},
		{{end -}}
		{{ range $idx, $i := loop .OutCount -}}
			 emit{{$i}} func (O{{$i}}),
		{{ end -}}
	) {
		{{- if gt .OutCount 0}}
			{{.OutputVars}} := f({{.InputVars}})
	
			{{- range $idx, $i := loop .OutCount }}
				emit{{$i}}(v{{$i}})
			{{ end -}}
		{{- else}}
			f({{.InputVars}})
		{{- end -}}
	})
}
`))

var streamNodeTpl = template.Must(template.New("").Funcs(funcs).Parse(`
type StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] interface {
	Node
	NodeHas{{.InCount}}In
	NodeHas{{.OutCount}}Out

	{{range $i := loop .InCount}}
		 NodeIn{{$i}}[I{{$i}}]
	{{- end}}

	{{range $i := loop .OutCount}}
		 NodeOut{{$i}}[O{{$i}}]
	{{- end}}

	Run({{.InputParametersDef}})
}

type FuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] func(
{{range $i := loop .InCount -}}
	 _ I{{$i}},
{{end -}}
{{ range $idx, $i := loop .OutCount -}}
	 emit{{$i}} func (O{{$i}}),
{{ end -}}
)

func NewFuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}](f FuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}] {
	return &funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]{Func: f}
}

type funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeDef}}] struct {
	StreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]

	Func FuncStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]

	machineryOnce sync.Once
	machinery *NodeMachinery
}

func (f *funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Inputs() int {
	return {{.InCount}}
}

func (f *funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Outputs() int {
	return {{.OutCount}}
}

func (f *funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Run({{.InputParametersDef}}) {
	f.Machinery().NewSourceRun(
		{{- range $idx, $i := loop .InCount}}
		 v{{$i}},
		{{- end}}
	)
}

func (f *funcStreamNode{{.InCount}}x{{.OutCount}}[{{.GenericsTypeRef}}]) Do(inputs []any, emit func(i int, v any)) {
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
func Take{{.N}}[T any](n interface{Node; NodeOut{{.N}}[T]}) interface{Node; NodeOut1[T]; NodeHas1Out } {
	return TakeN[T](n, {{.N}}-1)
}

func Pipeline{{.N}}[{{g_generics "T" .N true}}](
	{{g_params "from" "interface{Node; NodeOut1[T#]; NodeHas1Out}" .N }},
	to interface{ Node; NodeHas{{.N}}In; {{range $i := loop .N}} NodeIn{{$i}}[T{{$i}}]; {{end}} },
) {
	Pipeline(
		to,
	{{- range $idx, $i := loop .N}}
		from{{$i}},
	{{- end}}
	)
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

	for i := range loop(ic, true) {
		err := ifaceInTpl.Execute(w, map[string]interface{}{
			"N": i,
		})
		if err != nil {
			return err
		}
	}

	for i := range loop(oc, true) {
		err := ifaceOutTpl.Execute(w, map[string]interface{}{
			"N": i,
		})
		if err != nil {
			return err
		}
	}

	for i := range loop(oc, false) {
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
				inputParametersDef := genParametersList("", "v", "I", i)

				err := streamNodeTpl.Execute(w, map[string]interface{}{
					"InCount":            i,
					"OutCount":           o,
					"InputParametersDef": inputParametersDef,
					"GenericsTypeDef":    genericsTypeDef,
					"GenericsTypeRef":    genericsTypeRef,
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
