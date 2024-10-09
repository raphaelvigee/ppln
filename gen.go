package ppln

type NodeIn1[T any] interface {
	Node

	_input1(T)
}

type NodeIn2[T any] interface {
	Node

	_input2(T)
}

type NodeIn3[T any] interface {
	Node

	_input3(T)
}

type NodeOut1[T any] interface {
	Node

	_output1(T)
}

func Take1[T any](n NodeOut1[T]) NodeOut1[T] {
	panic("TODO")
}

func Pipeline1[T1 any](from1 NodeOut1[T1], to interface{ NodeIn1[T1] }) []Edge {
	return []Edge{
		{
			From: from1,
			To:   to,
		},
	}
}

type NodeOut2[T any] interface {
	Node

	_output2(T)
}

func Take2[T any](n NodeOut2[T]) NodeOut1[T] {
	panic("TODO")
}

func Pipeline2[T1 any, T2 any](from1 NodeOut1[T1], from2 NodeOut1[T2], to interface {
	NodeIn1[T1]
	NodeIn2[T2]
}) []Edge {
	return []Edge{
		{
			From: from1,
			To:   to,
		},
		{
			From: from2,
			To:   to,
		},
	}
}

type NodeOut3[T any] interface {
	Node

	_output3(T)
}

func Take3[T any](n NodeOut3[T]) NodeOut1[T] {
	panic("TODO")
}

func Pipeline3[T1 any, T2 any, T3 any](from1 NodeOut1[T1], from2 NodeOut1[T2], from3 NodeOut1[T3], to interface {
	NodeIn1[T1]
	NodeIn2[T2]
	NodeIn3[T3]
}) []Edge {
	return []Edge{
		{
			From: from1,
			To:   to,
		},
		{
			From: from2,
			To:   to,
		},
		{
			From: from3,
			To:   to,
		},
	}
}

type Node0x1[O1 any] interface {
	Node

	NodeOut1[O1]
}

type FuncNode0x1[O1 any] struct {
	Node0x1[O1]

	Func func() O1
}

func (f FuncNode0x1[O1]) Do(inputs []any) []any {

	v1 := f.Func()

	return []any{v1}
}

type StreamNode0x1[O1 any] interface {
	StreamNode

	NodeOut1[O1]
}

type FuncStreamNode0x1[O1 any] struct {
	StreamNode0x1[O1]

	Func func(
		func(v O1),
	)
}

func (f FuncStreamNode0x1[O1]) Do(inputs []any, emit func(i int, v any)) {

	f.Func(
		func(v O1) {
			emit(0, v)
		},
	)
}

type Node0x2[O1 any, O2 any] interface {
	Node

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode0x2[O1 any, O2 any] struct {
	Node0x2[O1, O2]

	Func func() (O1, O2)
}

func (f FuncNode0x2[O1, O2]) Do(inputs []any) []any {

	v1, v2 := f.Func()

	return []any{v1, v2}
}

type StreamNode0x2[O1 any, O2 any] interface {
	StreamNode

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncStreamNode0x2[O1 any, O2 any] struct {
	StreamNode0x2[O1, O2]

	Func func(
		func(v O1),
		func(v O2),
	)
}

func (f FuncStreamNode0x2[O1, O2]) Do(inputs []any, emit func(i int, v any)) {

	f.Func(
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
	)
}

type Node0x3[O1 any, O2 any, O3 any] interface {
	Node

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncNode0x3[O1 any, O2 any, O3 any] struct {
	Node0x3[O1, O2, O3]

	Func func() (O1, O2, O3)
}

func (f FuncNode0x3[O1, O2, O3]) Do(inputs []any) []any {

	v1, v2, v3 := f.Func()

	return []any{v1, v2, v3}
}

type StreamNode0x3[O1 any, O2 any, O3 any] interface {
	StreamNode

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncStreamNode0x3[O1 any, O2 any, O3 any] struct {
	StreamNode0x3[O1, O2, O3]

	Func func(
		func(v O1),
		func(v O2),
		func(v O3),
	)
}

func (f FuncStreamNode0x3[O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {

	f.Func(
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
		func(v O3) {
			emit(2, v)
		},
	)
}

type Node1x0[I1 any] interface {
	Node

	NodeIn1[I1]
}

type FuncNode1x0[I1 any] struct {
	Node1x0[I1]

	Func func(v1 I1)
}

func (f FuncNode1x0[I1]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)

	f.Func(i1)

	return nil
}

type StreamNode1x0[I1 any] interface {
	StreamNode

	NodeIn1[I1]
}

type FuncStreamNode1x0[I1 any] struct {
	StreamNode1x0[I1]

	Func func(
		I1,
	)
}

func (f FuncStreamNode1x0[I1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)

	f.Func(
		i1,
	)
}

type Node1x1[I1 any, O1 any] interface {
	Node

	NodeIn1[I1]

	NodeOut1[O1]
}

type FuncNode1x1[I1 any, O1 any] struct {
	Node1x1[I1, O1]

	Func func(v1 I1) O1
}

func (f FuncNode1x1[I1, O1]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)

	v1 := f.Func(i1)

	return []any{v1}
}

type StreamNode1x1[I1 any, O1 any] interface {
	StreamNode

	NodeIn1[I1]

	NodeOut1[O1]
}

type FuncStreamNode1x1[I1 any, O1 any] struct {
	StreamNode1x1[I1, O1]

	Func func(
		I1,
		func(v O1),
	)
}

func (f FuncStreamNode1x1[I1, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)

	f.Func(
		i1,
		func(v O1) {
			emit(0, v)
		},
	)
}

type Node1x2[I1 any, O1 any, O2 any] interface {
	Node

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncNode1x2[I1 any, O1 any, O2 any] struct {
	Node1x2[I1, O1, O2]

	Func func(v1 I1) (O1, O2)
}

func (f FuncNode1x2[I1, O1, O2]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)

	v1, v2 := f.Func(i1)

	return []any{v1, v2}
}

type StreamNode1x2[I1 any, O1 any, O2 any] interface {
	StreamNode

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncStreamNode1x2[I1 any, O1 any, O2 any] struct {
	StreamNode1x2[I1, O1, O2]

	Func func(
		I1,
		func(v O1),
		func(v O2),
	)
}

func (f FuncStreamNode1x2[I1, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)

	f.Func(
		i1,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
	)
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

	Func func(v1 I1) (O1, O2, O3)
}

func (f FuncNode1x3[I1, O1, O2, O3]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)

	v1, v2, v3 := f.Func(i1)

	return []any{v1, v2, v3}
}

type StreamNode1x3[I1 any, O1 any, O2 any, O3 any] interface {
	StreamNode

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncStreamNode1x3[I1 any, O1 any, O2 any, O3 any] struct {
	StreamNode1x3[I1, O1, O2, O3]

	Func func(
		I1,
		func(v O1),
		func(v O2),
		func(v O3),
	)
}

func (f FuncStreamNode1x3[I1, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)

	f.Func(
		i1,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
		func(v O3) {
			emit(2, v)
		},
	)
}

type Node2x0[I1 any, I2 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
}

type FuncNode2x0[I1 any, I2 any] struct {
	Node2x0[I1, I2]

	Func func(v1 I1, v2 I2)
}

func (f FuncNode2x0[I1, I2]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(i1, i2)

	return nil
}

type StreamNode2x0[I1 any, I2 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]
}

type FuncStreamNode2x0[I1 any, I2 any] struct {
	StreamNode2x0[I1, I2]

	Func func(
		I1,
		I2,
	)
}

func (f FuncStreamNode2x0[I1, I2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(
		i1,
		i2,
	)
}

type Node2x1[I1 any, I2 any, O1 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
}

type FuncNode2x1[I1 any, I2 any, O1 any] struct {
	Node2x1[I1, I2, O1]

	Func func(v1 I1, v2 I2) O1
}

func (f FuncNode2x1[I1, I2, O1]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	v1 := f.Func(i1, i2)

	return []any{v1}
}

type StreamNode2x1[I1 any, I2 any, O1 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
}

type FuncStreamNode2x1[I1 any, I2 any, O1 any] struct {
	StreamNode2x1[I1, I2, O1]

	Func func(
		I1,
		I2,
		func(v O1),
	)
}

func (f FuncStreamNode2x1[I1, I2, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(
		i1,
		i2,
		func(v O1) {
			emit(0, v)
		},
	)
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

	Func func(v1 I1, v2 I2) (O1, O2)
}

func (f FuncNode2x2[I1, I2, O1, O2]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	v1, v2 := f.Func(i1, i2)

	return []any{v1, v2}
}

type StreamNode2x2[I1 any, I2 any, O1 any, O2 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncStreamNode2x2[I1 any, I2 any, O1 any, O2 any] struct {
	StreamNode2x2[I1, I2, O1, O2]

	Func func(
		I1,
		I2,
		func(v O1),
		func(v O2),
	)
}

func (f FuncStreamNode2x2[I1, I2, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(
		i1,
		i2,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
	)
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

	Func func(v1 I1, v2 I2) (O1, O2, O3)
}

func (f FuncNode2x3[I1, I2, O1, O2, O3]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	v1, v2, v3 := f.Func(i1, i2)

	return []any{v1, v2, v3}
}

type StreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncStreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] struct {
	StreamNode2x3[I1, I2, O1, O2, O3]

	Func func(
		I1,
		I2,
		func(v O1),
		func(v O2),
		func(v O3),
	)
}

func (f FuncStreamNode2x3[I1, I2, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(
		i1,
		i2,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
		func(v O3) {
			emit(2, v)
		},
	)
}

type Node3x0[I1 any, I2 any, I3 any] interface {
	Node

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]
}

type FuncNode3x0[I1 any, I2 any, I3 any] struct {
	Node3x0[I1, I2, I3]

	Func func(v1 I1, v2 I2, v3 I3)
}

func (f FuncNode3x0[I1, I2, I3]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	f.Func(i1, i2, i3)

	return nil
}

type StreamNode3x0[I1 any, I2 any, I3 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]
}

type FuncStreamNode3x0[I1 any, I2 any, I3 any] struct {
	StreamNode3x0[I1, I2, I3]

	Func func(
		I1,
		I2,
		I3,
	)
}

func (f FuncStreamNode3x0[I1, I2, I3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	f.Func(
		i1,
		i2,
		i3,
	)
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

	Func func(v1 I1, v2 I2, v3 I3) O1
}

func (f FuncNode3x1[I1, I2, I3, O1]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	v1 := f.Func(i1, i2, i3)

	return []any{v1}
}

type StreamNode3x1[I1 any, I2 any, I3 any, O1 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
}

type FuncStreamNode3x1[I1 any, I2 any, I3 any, O1 any] struct {
	StreamNode3x1[I1, I2, I3, O1]

	Func func(
		I1,
		I2,
		I3,
		func(v O1),
	)
}

func (f FuncStreamNode3x1[I1, I2, I3, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	f.Func(
		i1,
		i2,
		i3,
		func(v O1) {
			emit(0, v)
		},
	)
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

	Func func(v1 I1, v2 I2, v3 I3) (O1, O2)
}

func (f FuncNode3x2[I1, I2, I3, O1, O2]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	v1, v2 := f.Func(i1, i2, i3)

	return []any{v1, v2}
}

type StreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
}

type FuncStreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] struct {
	StreamNode3x2[I1, I2, I3, O1, O2]

	Func func(
		I1,
		I2,
		I3,
		func(v O1),
		func(v O2),
	)
}

func (f FuncStreamNode3x2[I1, I2, I3, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	f.Func(
		i1,
		i2,
		i3,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
	)
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

	Func func(v1 I1, v2 I2, v3 I3) (O1, O2, O3)
}

func (f FuncNode3x3[I1, I2, I3, O1, O2, O3]) Do(inputs []any) []any {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	v1, v2, v3 := f.Func(i1, i2, i3)

	return []any{v1, v2, v3}
}

type StreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] interface {
	StreamNode

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]
}

type FuncStreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] struct {
	StreamNode3x3[I1, I2, I3, O1, O2, O3]

	Func func(
		I1,
		I2,
		I3,
		func(v O1),
		func(v O2),
		func(v O3),
	)
}

func (f FuncStreamNode3x3[I1, I2, I3, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)

	f.Func(
		i1,
		i2,
		i3,
		func(v O1) {
			emit(0, v)
		},
		func(v O2) {
			emit(1, v)
		},
		func(v O3) {
			emit(2, v)
		},
	)
}
