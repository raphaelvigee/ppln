//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strings"
)

type Node interface {
	Do()

	OutputsValues() []any

	Inputs() []Node
	Outputs() []Node
}

type Node1x1[I, O any] interface {
	NodeIn1[I]
	NodeOut1[O]

	Node
}

type FuncNode1x1[I, O any] struct {
	Node1x1[I, O]

	outputs []any

	Func func(i I) O
}

func (f FuncNode1x1[I, O]) Do() {
	i := f.Inputs()[0].OutputsValues()[0].(I)

	v := f.Func(i)
	f.outputs = []any{v}
}

func (f FuncNode1x1[I, O]) OutputsValues() []any {
	return f.outputs
}

type NodeIn1[T any] interface {
	_input1(T)

	SetDep1(NodeOut1[T])
}

type NodeIn2[T any] interface {
	_input2(T)

	SetDep2(NodeOut1[T])
}

type NodeOut1[T any] interface {
	_output1(T)

	GetOut1() T
}
type NodeOut2[T any] interface {
	_output2(T)

	GetOut2() T
}

type Node2x1[I1, I2, O any] interface {
	NodeIn1[I1]
	NodeIn2[I2]
	NodeOut1[O]

	Node
}

type FuncNode2x1[I1, I2, O any] struct {
	Node2x1[I1, I2, O]

	outputs []any

	Func func(i1 I1, i2 I2) O
}

func (f FuncNode2x1[I1, I2, O]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)

	v := f.Func(i1, i2)
	f.outputs = []any{v}
}

func (f FuncNode2x1[I1, I2, O]) OutputsValues() []any {
	return f.outputs
}

func Link1x1[T any](n1 NodeOut1[T], n2 NodeIn1[T]) {
	n2.SetDep1(n1)
}

func Link2x1[T1, T2 any](i1 NodeOut1[T1], i2 NodeOut1[T2], n2 interface {
	NodeIn1[T1]
	NodeIn2[T2]
}) {

}

func Run(n any) {

}

func main() {

	{
		c1 := FuncNode1x1[string, int]{Func: func(s string) int {
			return len(strings.Split(s, "\n"))
		}}
		c2 := FuncNode1x1[string, int]{Func: func(s string) int {
			return len(strings.Split(s, "\n"))
		}}
		c3 := FuncNode2x1[int, int, struct{}]{Func: func(v1, v2 int) struct{} {
			fmt.Println(v1 + v2)

			return struct{}{}
		}}
		c3.SetDep1(c1)
		c3.SetDep2(c2)
		// c1 => out 2 things
		// c2 => out 2 things
		// c3 => input c1.1 c2.2
	}
	{
		c1 := FuncNode1x1[string, int]{Func: func(s string) int {
			return len(strings.Split(s, "\n"))
		}}
		c2 := FuncNode1x1[int, struct{}]{Func: func(i int) struct{} {
			fmt.Println(i)

			return struct{}{}
		}}

		Link1x1[int](c1, c2)
	}

	{
		c1 := FuncNode1x1[string, int]{Func: func(s string) int {
			return len(strings.Split(s, "\n"))
		}}
		c2 := FuncNode1x1[string, int]{Func: func(s string) int {
			return len(strings.Split(s, "\n"))
		}}
		n2 := FuncNode2x1[int, int, struct{}]{Func: func(i1, i2 int) struct{} {
			fmt.Println(i1, i2)

			return struct{}{}
		}}

		Link2x1[int, int](c1, c2, n2)
	}
}
