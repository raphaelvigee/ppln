package ppln

type NodeIn1[T any] interface {
	_input1(T)

	SetDep1(NodeOut1[T])
}

type NodeIn2[T any] interface {
	_input2(T)

	SetDep2(NodeOut1[T])
}

type NodeIn3[T any] interface {
	_input3(T)

	SetDep3(NodeOut1[T])
}

type NodeOut1[T any] interface {
	_output1(T)
}

func Take1[T any](n NodeIn1[T]) NodeOut1[T] {
	panic("TODO")
}

type NodeOut2[T any] interface {
	_output2(T)
}

func Take2[T any](n NodeIn2[T]) NodeOut1[T] {
	panic("TODO")
}

type NodeOut3[T any] interface {
	_output3(T)
}

func Take3[T any](n NodeIn3[T]) NodeOut1[T] {
	panic("TODO")
}

type Node0x1[O1 any] interface {
	Node

	NodeOut1[O1]
}

type FuncNode0x1[O1 any] struct {
	Node0x1[O1]

	outputs []any

	Func func() O1
}

func (f FuncNode0x1[O1]) Do() {

	v1 := f.Func()

	f.outputs = []any{
		v1,
	}
}

func (f FuncNode0x1[O1]) OutputsValues() []any {
	return f.outputs
}

type Node0x2[O1 any, O2 any] interface {
	Node

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode0x2[O1 any, O2 any] struct {
	Node0x2[O1, O2]

	outputs []any

	Func func() (O1, O2)
}

func (f FuncNode0x2[O1, O2]) Do() {

	v1, v2 := f.Func()

	f.outputs = []any{
		v1,
		v2,
	}
}

func (f FuncNode0x2[O1, O2]) OutputsValues() []any {
	return f.outputs
}

type Node0x3[O1 any, O2 any, O3 any] interface {
	Node

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncNode0x3[O1 any, O2 any, O3 any] struct {
	Node0x3[O1, O2, O3]

	outputs []any

	Func func() (O1, O2, O3)
}

func (f FuncNode0x3[O1, O2, O3]) Do() {

	v1, v2, v3 := f.Func()

	f.outputs = []any{
		v1,
		v2,
		v3,
	}
}

func (f FuncNode0x3[O1, O2, O3]) OutputsValues() []any {
	return f.outputs
}

type Node1x0[I1 any] interface {
	Node

	NodeIn1[I1]
}

type FuncNode1x0[I1 any] struct {
	Node1x0[I1]

	outputs []any

	Func func(v1 I1)
}

func (f FuncNode1x0[I1]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)

	f.Func(
		i1,
	)
}

func (f FuncNode1x0[I1]) OutputsValues() []any {
	return f.outputs
}

type Node1x1[I1 any, O1 any] interface {
	Node

	NodeIn1[I1]

	NodeOut1[O1]
}

type FuncNode1x1[I1 any, O1 any] struct {
	Node1x1[I1, O1]

	outputs []any

	Func func(v1 I1) O1
}

func (f FuncNode1x1[I1, O1]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)

	v1 := f.Func(
		i1,
	)

	f.outputs = []any{
		v1,
	}
}

func (f FuncNode1x1[I1, O1]) OutputsValues() []any {
	return f.outputs
}

type Node1x2[I1 any, O1 any, O2 any] interface {
	Node

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode1x2[I1 any, O1 any, O2 any] struct {
	Node1x2[I1, O1, O2]

	outputs []any

	Func func(v1 I1) (O1, O2)
}

func (f FuncNode1x2[I1, O1, O2]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)

	v1, v2 := f.Func(
		i1,
	)

	f.outputs = []any{
		v1,
		v2,
	}
}

func (f FuncNode1x2[I1, O1, O2]) OutputsValues() []any {
	return f.outputs
}

type Node1x3[I1 any, O1 any, O2 any, O3 any] interface {
	Node

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncNode1x3[I1 any, O1 any, O2 any, O3 any] struct {
	Node1x3[I1, O1, O2, O3]

	outputs []any

	Func func(v1 I1) (O1, O2, O3)
}

func (f FuncNode1x3[I1, O1, O2, O3]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)

	v1, v2, v3 := f.Func(
		i1,
	)

	f.outputs = []any{
		v1,
		v2,
		v3,
	}
}

func (f FuncNode1x3[I1, O1, O2, O3]) OutputsValues() []any {
	return f.outputs
}

type Node2x0[I1 any, I2 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
}

type FuncNode2x0[I1 any, I2 any] struct {
	Node2x0[I1, I2]

	outputs []any

	Func func(v1 I1, v2 I2)
}

func (f FuncNode2x0[I1, I2]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)

	f.Func(
		i1,
		i2,
	)
}

func (f FuncNode2x0[I1, I2]) OutputsValues() []any {
	return f.outputs
}

type Node2x1[I1 any, I2 any, O1 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
}

type FuncNode2x1[I1 any, I2 any, O1 any] struct {
	Node2x1[I1, I2, O1]

	outputs []any

	Func func(v1 I1, v2 I2) O1
}

func (f FuncNode2x1[I1, I2, O1]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)

	v1 := f.Func(
		i1,
		i2,
	)

	f.outputs = []any{
		v1,
	}
}

func (f FuncNode2x1[I1, I2, O1]) OutputsValues() []any {
	return f.outputs
}

type Node2x2[I1 any, I2 any, O1 any, O2 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode2x2[I1 any, I2 any, O1 any, O2 any] struct {
	Node2x2[I1, I2, O1, O2]

	outputs []any

	Func func(v1 I1, v2 I2) (O1, O2)
}

func (f FuncNode2x2[I1, I2, O1, O2]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)

	v1, v2 := f.Func(
		i1,
		i2,
	)

	f.outputs = []any{
		v1,
		v2,
	}
}

func (f FuncNode2x2[I1, I2, O1, O2]) OutputsValues() []any {
	return f.outputs
}

type Node2x3[I1 any, I2 any, O1 any, O2 any, O3 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] struct {
	Node2x3[I1, I2, O1, O2, O3]

	outputs []any

	Func func(v1 I1, v2 I2) (O1, O2, O3)
}

func (f FuncNode2x3[I1, I2, O1, O2, O3]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)

	v1, v2, v3 := f.Func(
		i1,
		i2,
	)

	f.outputs = []any{
		v1,
		v2,
		v3,
	}
}

func (f FuncNode2x3[I1, I2, O1, O2, O3]) OutputsValues() []any {
	return f.outputs
}

type Node3x0[I1 any, I2 any, I3 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]
}

type FuncNode3x0[I1 any, I2 any, I3 any] struct {
	Node3x0[I1, I2, I3]

	outputs []any

	Func func(v1 I1, v2 I2, v3 I3)
}

func (f FuncNode3x0[I1, I2, I3]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)
	i3 := f.Inputs()[2].OutputsValues()[0].(I3)

	f.Func(
		i1,
		i2,
		i3,
	)
}

func (f FuncNode3x0[I1, I2, I3]) OutputsValues() []any {
	return f.outputs
}

type Node3x1[I1 any, I2 any, I3 any, O1 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
}

type FuncNode3x1[I1 any, I2 any, I3 any, O1 any] struct {
	Node3x1[I1, I2, I3, O1]

	outputs []any

	Func func(v1 I1, v2 I2, v3 I3) O1
}

func (f FuncNode3x1[I1, I2, I3, O1]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)
	i3 := f.Inputs()[2].OutputsValues()[0].(I3)

	v1 := f.Func(
		i1,
		i2,
		i3,
	)

	f.outputs = []any{
		v1,
	}
}

func (f FuncNode3x1[I1, I2, I3, O1]) OutputsValues() []any {
	return f.outputs
}

type Node3x2[I1 any, I2 any, I3 any, O1 any, O2 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] struct {
	Node3x2[I1, I2, I3, O1, O2]

	outputs []any

	Func func(v1 I1, v2 I2, v3 I3) (O1, O2)
}

func (f FuncNode3x2[I1, I2, I3, O1, O2]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)
	i3 := f.Inputs()[2].OutputsValues()[0].(I3)

	v1, v2 := f.Func(
		i1,
		i2,
		i3,
	)

	f.outputs = []any{
		v1,
		v2,
	}
}

func (f FuncNode3x2[I1, I2, I3, O1, O2]) OutputsValues() []any {
	return f.outputs
}

type Node3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] struct {
	Node3x3[I1, I2, I3, O1, O2, O3]

	outputs []any

	Func func(v1 I1, v2 I2, v3 I3) (O1, O2, O3)
}

func (f FuncNode3x3[I1, I2, I3, O1, O2, O3]) Do() {
	i1 := f.Inputs()[0].OutputsValues()[0].(I1)
	i2 := f.Inputs()[1].OutputsValues()[0].(I2)
	i3 := f.Inputs()[2].OutputsValues()[0].(I3)

	v1, v2, v3 := f.Func(
		i1,
		i2,
		i3,
	)

	f.outputs = []any{
		v1,
		v2,
		v3,
	}
}

func (f FuncNode3x3[I1, I2, I3, O1, O2, O3]) OutputsValues() []any {
	return f.outputs
}
