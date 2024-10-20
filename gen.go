package ppln

import "sync"

type NodeIn0[T any] interface {
	Node

	_input0(T)
}

type NodeHas0In interface {
	Node

	_input_layout()
}

type NodeIn1[T any] interface {
	Node

	_input1(T)
}

type NodeHas1In interface {
	Node

	_input_layout(any)
}

type NodeIn2[T any] interface {
	Node

	_input2(T)
}

type NodeHas2In interface {
	Node

	_input_layout(any, any)
}

type NodeIn3[T any] interface {
	Node

	_input3(T)
}

type NodeHas3In interface {
	Node

	_input_layout(any, any, any)
}

type NodeOut1[T any] interface {
	_output1(T)
}
type NodeHas1Out interface {
	_out_layout(any)
}

func Take1[T any](n interface {
	Node
	NodeOut1[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 1)
}

func Pipeline1[T1 any](from1 interface {
	Node
	NodeOut1[T1]
	NodeHas1Out
}, to interface {
	Node
	NodeHas1In
	NodeIn1[T1]
}) {
	Pipeline(
		to,
		from1,
	)
}

type NodeOut2[T any] interface {
	_output2(T)
}
type NodeHas2Out interface {
	_out_layout(any, any)
}

func Take2[T any](n interface {
	Node
	NodeOut2[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 2)
}

func Pipeline2[T1 any, T2 any](from1 interface {
	Node
	NodeOut1[T1]
	NodeHas1Out
}, from2 interface {
	Node
	NodeOut1[T2]
	NodeHas1Out
}, to interface {
	Node
	NodeHas2In
	NodeIn1[T1]
	NodeIn2[T2]
}) {
	Pipeline(
		to,
		from1,
		from2,
	)
}

type NodeOut3[T any] interface {
	_output3(T)
}
type NodeHas3Out interface {
	_out_layout(any, any, any)
}

func Take3[T any](n interface {
	Node
	NodeOut3[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 3)
}

func Pipeline3[T1 any, T2 any, T3 any](from1 interface {
	Node
	NodeOut1[T1]
	NodeHas1Out
}, from2 interface {
	Node
	NodeOut1[T2]
	NodeHas1Out
}, from3 interface {
	Node
	NodeOut1[T3]
	NodeHas1Out
}, to interface {
	Node
	NodeHas3In
	NodeIn1[T1]
	NodeIn2[T2]
	NodeIn3[T3]
}) {
	Pipeline(
		to,
		from1,
		from2,
		from3,
	)
}

type Node0x1[O1 any] interface {
	Node
	NodeHas0In
	NodeHas1Out

	NodeOut1[O1]

	Run()
}

type FuncNode0x1[O1 any] func() O1

func NewFuncNode0x1[O1 any](f FuncNode0x1[O1]) Node0x1[O1] {
	return &funcNode0x1[O1]{
		Func: f,
	}
}

type funcNode0x1[O1 any] struct {
	Node0x1[O1]

	Func FuncNode0x1[O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode0x1[O1]) Inputs() int {
	return 0
}

func (f *funcNode0x1[O1]) Outputs() int {
	return 1
}

func (f *funcNode0x1[O1]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcNode0x1[O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode0x1[O1]) Do(inputs []any, emit func(i int, v any)) {
	v1 := f.Func()
	emit(0, v1)
}

type StreamNode0x1[O1 any] interface {
	Node
	NodeHas0In
	NodeHas1Out

	NodeOut1[O1]

	Run()
}

type FuncStreamNode0x1[O1 any] func(
	emit1 func(O1),
)

func NewFuncStreamNode0x1[O1 any](f FuncStreamNode0x1[O1]) StreamNode0x1[O1] {
	return &funcStreamNode0x1[O1]{Func: f}
}

type funcStreamNode0x1[O1 any] struct {
	StreamNode0x1[O1]

	Func FuncStreamNode0x1[O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode0x1[O1]) Inputs() int {
	return 0
}

func (f *funcStreamNode0x1[O1]) Outputs() int {
	return 1
}

func (f *funcStreamNode0x1[O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode0x1[O1]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcStreamNode0x1[O1]) Do(inputs []any, emit func(i int, v any)) {

	f.Func(
		func(v O1) {
			emit(0, v)
		},
	)
}

type Node0x2[O1 any, O2 any] interface {
	Node
	NodeHas0In
	NodeHas2Out

	NodeOut1[O1]
	NodeOut2[O2]

	Run()
}

type FuncNode0x2[O1 any, O2 any] func() (O1, O2)

func NewFuncNode0x2[O1 any, O2 any](f FuncNode0x2[O1, O2]) Node0x2[O1, O2] {
	return &funcNode0x2[O1, O2]{
		Func: f,
	}
}

type funcNode0x2[O1 any, O2 any] struct {
	Node0x2[O1, O2]

	Func FuncNode0x2[O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode0x2[O1, O2]) Inputs() int {
	return 0
}

func (f *funcNode0x2[O1, O2]) Outputs() int {
	return 2
}

func (f *funcNode0x2[O1, O2]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcNode0x2[O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode0x2[O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	v1, v2 := f.Func()
	emit(0, v1)

	emit(1, v2)
}

type StreamNode0x2[O1 any, O2 any] interface {
	Node
	NodeHas0In
	NodeHas2Out

	NodeOut1[O1]
	NodeOut2[O2]

	Run()
}

type FuncStreamNode0x2[O1 any, O2 any] func(
	emit1 func(O1),
	emit2 func(O2),
)

func NewFuncStreamNode0x2[O1 any, O2 any](f FuncStreamNode0x2[O1, O2]) StreamNode0x2[O1, O2] {
	return &funcStreamNode0x2[O1, O2]{Func: f}
}

type funcStreamNode0x2[O1 any, O2 any] struct {
	StreamNode0x2[O1, O2]

	Func FuncStreamNode0x2[O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode0x2[O1, O2]) Inputs() int {
	return 0
}

func (f *funcStreamNode0x2[O1, O2]) Outputs() int {
	return 2
}

func (f *funcStreamNode0x2[O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode0x2[O1, O2]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcStreamNode0x2[O1, O2]) Do(inputs []any, emit func(i int, v any)) {

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
	NodeHas0In
	NodeHas3Out

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run()
}

type FuncNode0x3[O1 any, O2 any, O3 any] func() (O1, O2, O3)

func NewFuncNode0x3[O1 any, O2 any, O3 any](f FuncNode0x3[O1, O2, O3]) Node0x3[O1, O2, O3] {
	return &funcNode0x3[O1, O2, O3]{
		Func: f,
	}
}

type funcNode0x3[O1 any, O2 any, O3 any] struct {
	Node0x3[O1, O2, O3]

	Func FuncNode0x3[O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode0x3[O1, O2, O3]) Inputs() int {
	return 0
}

func (f *funcNode0x3[O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcNode0x3[O1, O2, O3]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcNode0x3[O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode0x3[O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	v1, v2, v3 := f.Func()
	emit(0, v1)

	emit(1, v2)

	emit(2, v3)
}

type StreamNode0x3[O1 any, O2 any, O3 any] interface {
	Node
	NodeHas0In
	NodeHas3Out

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run()
}

type FuncStreamNode0x3[O1 any, O2 any, O3 any] func(
	emit1 func(O1),
	emit2 func(O2),
	emit3 func(O3),
)

func NewFuncStreamNode0x3[O1 any, O2 any, O3 any](f FuncStreamNode0x3[O1, O2, O3]) StreamNode0x3[O1, O2, O3] {
	return &funcStreamNode0x3[O1, O2, O3]{Func: f}
}

type funcStreamNode0x3[O1 any, O2 any, O3 any] struct {
	StreamNode0x3[O1, O2, O3]

	Func FuncStreamNode0x3[O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode0x3[O1, O2, O3]) Inputs() int {
	return 0
}

func (f *funcStreamNode0x3[O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcStreamNode0x3[O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode0x3[O1, O2, O3]) Run() {
	f.Machinery().NewSourceRun()
}

func (f *funcStreamNode0x3[O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {

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
	NodeHas1In
	NodeHas0Out

	NodeIn1[I1]

	Run(v1 I1)
}

type FuncNode1x0[I1 any] func(v1 I1)

func NewFuncNode1x0[I1 any](f FuncNode1x0[I1]) Node1x0[I1] {
	return &funcNode1x0[I1]{
		Func: f,
	}
}

type funcNode1x0[I1 any] struct {
	Node1x0[I1]

	Func FuncNode1x0[I1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode1x0[I1]) Inputs() int {
	return 1
}

func (f *funcNode1x0[I1]) Outputs() int {
	return 0
}

func (f *funcNode1x0[I1]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcNode1x0[I1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode1x0[I1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	f.Func(i1)
}

type StreamNode1x0[I1 any] interface {
	Node
	NodeHas1In
	NodeHas0Out

	NodeIn1[I1]

	Run(v1 I1)
}

type FuncStreamNode1x0[I1 any] func(
	_ I1,
)

func NewFuncStreamNode1x0[I1 any](f FuncStreamNode1x0[I1]) StreamNode1x0[I1] {
	return &funcStreamNode1x0[I1]{Func: f}
}

type funcStreamNode1x0[I1 any] struct {
	StreamNode1x0[I1]

	Func FuncStreamNode1x0[I1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode1x0[I1]) Inputs() int {
	return 1
}

func (f *funcStreamNode1x0[I1]) Outputs() int {
	return 0
}

func (f *funcStreamNode1x0[I1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode1x0[I1]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcStreamNode1x0[I1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)

	f.Func(
		i1,
	)
}

type Node1x1[I1 any, O1 any] interface {
	Node
	NodeHas1In
	NodeHas1Out

	NodeIn1[I1]

	NodeOut1[O1]

	Run(v1 I1)
}

type FuncNode1x1[I1 any, O1 any] func(v1 I1) O1

func NewFuncNode1x1[I1 any, O1 any](f FuncNode1x1[I1, O1]) Node1x1[I1, O1] {
	return &funcNode1x1[I1, O1]{
		Func: f,
	}
}

type funcNode1x1[I1 any, O1 any] struct {
	Node1x1[I1, O1]

	Func FuncNode1x1[I1, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode1x1[I1, O1]) Inputs() int {
	return 1
}

func (f *funcNode1x1[I1, O1]) Outputs() int {
	return 1
}

func (f *funcNode1x1[I1, O1]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcNode1x1[I1, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode1x1[I1, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	v1 := f.Func(i1)
	emit(0, v1)
}

type StreamNode1x1[I1 any, O1 any] interface {
	Node
	NodeHas1In
	NodeHas1Out

	NodeIn1[I1]

	NodeOut1[O1]

	Run(v1 I1)
}

type FuncStreamNode1x1[I1 any, O1 any] func(
	_ I1,
	emit1 func(O1),
)

func NewFuncStreamNode1x1[I1 any, O1 any](f FuncStreamNode1x1[I1, O1]) StreamNode1x1[I1, O1] {
	return &funcStreamNode1x1[I1, O1]{Func: f}
}

type funcStreamNode1x1[I1 any, O1 any] struct {
	StreamNode1x1[I1, O1]

	Func FuncStreamNode1x1[I1, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode1x1[I1, O1]) Inputs() int {
	return 1
}

func (f *funcStreamNode1x1[I1, O1]) Outputs() int {
	return 1
}

func (f *funcStreamNode1x1[I1, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode1x1[I1, O1]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcStreamNode1x1[I1, O1]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas1In
	NodeHas2Out

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1)
}

type FuncNode1x2[I1 any, O1 any, O2 any] func(v1 I1) (O1, O2)

func NewFuncNode1x2[I1 any, O1 any, O2 any](f FuncNode1x2[I1, O1, O2]) Node1x2[I1, O1, O2] {
	return &funcNode1x2[I1, O1, O2]{
		Func: f,
	}
}

type funcNode1x2[I1 any, O1 any, O2 any] struct {
	Node1x2[I1, O1, O2]

	Func FuncNode1x2[I1, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode1x2[I1, O1, O2]) Inputs() int {
	return 1
}

func (f *funcNode1x2[I1, O1, O2]) Outputs() int {
	return 2
}

func (f *funcNode1x2[I1, O1, O2]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcNode1x2[I1, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode1x2[I1, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	v1, v2 := f.Func(i1)
	emit(0, v1)

	emit(1, v2)
}

type StreamNode1x2[I1 any, O1 any, O2 any] interface {
	Node
	NodeHas1In
	NodeHas2Out

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1)
}

type FuncStreamNode1x2[I1 any, O1 any, O2 any] func(
	_ I1,
	emit1 func(O1),
	emit2 func(O2),
)

func NewFuncStreamNode1x2[I1 any, O1 any, O2 any](f FuncStreamNode1x2[I1, O1, O2]) StreamNode1x2[I1, O1, O2] {
	return &funcStreamNode1x2[I1, O1, O2]{Func: f}
}

type funcStreamNode1x2[I1 any, O1 any, O2 any] struct {
	StreamNode1x2[I1, O1, O2]

	Func FuncStreamNode1x2[I1, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode1x2[I1, O1, O2]) Inputs() int {
	return 1
}

func (f *funcStreamNode1x2[I1, O1, O2]) Outputs() int {
	return 2
}

func (f *funcStreamNode1x2[I1, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode1x2[I1, O1, O2]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcStreamNode1x2[I1, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas1In
	NodeHas3Out

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1)
}

type FuncNode1x3[I1 any, O1 any, O2 any, O3 any] func(v1 I1) (O1, O2, O3)

func NewFuncNode1x3[I1 any, O1 any, O2 any, O3 any](f FuncNode1x3[I1, O1, O2, O3]) Node1x3[I1, O1, O2, O3] {
	return &funcNode1x3[I1, O1, O2, O3]{
		Func: f,
	}
}

type funcNode1x3[I1 any, O1 any, O2 any, O3 any] struct {
	Node1x3[I1, O1, O2, O3]

	Func FuncNode1x3[I1, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode1x3[I1, O1, O2, O3]) Inputs() int {
	return 1
}

func (f *funcNode1x3[I1, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcNode1x3[I1, O1, O2, O3]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcNode1x3[I1, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode1x3[I1, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	v1, v2, v3 := f.Func(i1)
	emit(0, v1)

	emit(1, v2)

	emit(2, v3)
}

type StreamNode1x3[I1 any, O1 any, O2 any, O3 any] interface {
	Node
	NodeHas1In
	NodeHas3Out

	NodeIn1[I1]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1)
}

type FuncStreamNode1x3[I1 any, O1 any, O2 any, O3 any] func(
	_ I1,
	emit1 func(O1),
	emit2 func(O2),
	emit3 func(O3),
)

func NewFuncStreamNode1x3[I1 any, O1 any, O2 any, O3 any](f FuncStreamNode1x3[I1, O1, O2, O3]) StreamNode1x3[I1, O1, O2, O3] {
	return &funcStreamNode1x3[I1, O1, O2, O3]{Func: f}
}

type funcStreamNode1x3[I1 any, O1 any, O2 any, O3 any] struct {
	StreamNode1x3[I1, O1, O2, O3]

	Func FuncStreamNode1x3[I1, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Inputs() int {
	return 1
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Run(v1 I1) {
	f.Machinery().NewSourceRun(
		v1,
	)
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas2In
	NodeHas0Out

	NodeIn1[I1]
	NodeIn2[I2]

	Run(v1 I1, v2 I2)
}

type FuncNode2x0[I1 any, I2 any] func(v1 I1, v2 I2)

func NewFuncNode2x0[I1 any, I2 any](f FuncNode2x0[I1, I2]) Node2x0[I1, I2] {
	return &funcNode2x0[I1, I2]{
		Func: f,
	}
}

type funcNode2x0[I1 any, I2 any] struct {
	Node2x0[I1, I2]

	Func FuncNode2x0[I1, I2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode2x0[I1, I2]) Inputs() int {
	return 2
}

func (f *funcNode2x0[I1, I2]) Outputs() int {
	return 0
}

func (f *funcNode2x0[I1, I2]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcNode2x0[I1, I2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode2x0[I1, I2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	f.Func(i1, i2)
}

type StreamNode2x0[I1 any, I2 any] interface {
	Node
	NodeHas2In
	NodeHas0Out

	NodeIn1[I1]
	NodeIn2[I2]

	Run(v1 I1, v2 I2)
}

type FuncStreamNode2x0[I1 any, I2 any] func(
	_ I1,
	_ I2,
)

func NewFuncStreamNode2x0[I1 any, I2 any](f FuncStreamNode2x0[I1, I2]) StreamNode2x0[I1, I2] {
	return &funcStreamNode2x0[I1, I2]{Func: f}
}

type funcStreamNode2x0[I1 any, I2 any] struct {
	StreamNode2x0[I1, I2]

	Func FuncStreamNode2x0[I1, I2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode2x0[I1, I2]) Inputs() int {
	return 2
}

func (f *funcStreamNode2x0[I1, I2]) Outputs() int {
	return 0
}

func (f *funcStreamNode2x0[I1, I2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode2x0[I1, I2]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcStreamNode2x0[I1, I2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)

	f.Func(
		i1,
		i2,
	)
}

type Node2x1[I1 any, I2 any, O1 any] interface {
	Node
	NodeHas2In
	NodeHas1Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]

	Run(v1 I1, v2 I2)
}

type FuncNode2x1[I1 any, I2 any, O1 any] func(v1 I1, v2 I2) O1

func NewFuncNode2x1[I1 any, I2 any, O1 any](f FuncNode2x1[I1, I2, O1]) Node2x1[I1, I2, O1] {
	return &funcNode2x1[I1, I2, O1]{
		Func: f,
	}
}

type funcNode2x1[I1 any, I2 any, O1 any] struct {
	Node2x1[I1, I2, O1]

	Func FuncNode2x1[I1, I2, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode2x1[I1, I2, O1]) Inputs() int {
	return 2
}

func (f *funcNode2x1[I1, I2, O1]) Outputs() int {
	return 1
}

func (f *funcNode2x1[I1, I2, O1]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcNode2x1[I1, I2, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode2x1[I1, I2, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	v1 := f.Func(i1, i2)
	emit(0, v1)
}

type StreamNode2x1[I1 any, I2 any, O1 any] interface {
	Node
	NodeHas2In
	NodeHas1Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]

	Run(v1 I1, v2 I2)
}

type FuncStreamNode2x1[I1 any, I2 any, O1 any] func(
	_ I1,
	_ I2,
	emit1 func(O1),
)

func NewFuncStreamNode2x1[I1 any, I2 any, O1 any](f FuncStreamNode2x1[I1, I2, O1]) StreamNode2x1[I1, I2, O1] {
	return &funcStreamNode2x1[I1, I2, O1]{Func: f}
}

type funcStreamNode2x1[I1 any, I2 any, O1 any] struct {
	StreamNode2x1[I1, I2, O1]

	Func FuncStreamNode2x1[I1, I2, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode2x1[I1, I2, O1]) Inputs() int {
	return 2
}

func (f *funcStreamNode2x1[I1, I2, O1]) Outputs() int {
	return 1
}

func (f *funcStreamNode2x1[I1, I2, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode2x1[I1, I2, O1]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcStreamNode2x1[I1, I2, O1]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas2In
	NodeHas2Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1, v2 I2)
}

type FuncNode2x2[I1 any, I2 any, O1 any, O2 any] func(v1 I1, v2 I2) (O1, O2)

func NewFuncNode2x2[I1 any, I2 any, O1 any, O2 any](f FuncNode2x2[I1, I2, O1, O2]) Node2x2[I1, I2, O1, O2] {
	return &funcNode2x2[I1, I2, O1, O2]{
		Func: f,
	}
}

type funcNode2x2[I1 any, I2 any, O1 any, O2 any] struct {
	Node2x2[I1, I2, O1, O2]

	Func FuncNode2x2[I1, I2, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode2x2[I1, I2, O1, O2]) Inputs() int {
	return 2
}

func (f *funcNode2x2[I1, I2, O1, O2]) Outputs() int {
	return 2
}

func (f *funcNode2x2[I1, I2, O1, O2]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcNode2x2[I1, I2, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode2x2[I1, I2, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	v1, v2 := f.Func(i1, i2)
	emit(0, v1)

	emit(1, v2)
}

type StreamNode2x2[I1 any, I2 any, O1 any, O2 any] interface {
	Node
	NodeHas2In
	NodeHas2Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1, v2 I2)
}

type FuncStreamNode2x2[I1 any, I2 any, O1 any, O2 any] func(
	_ I1,
	_ I2,
	emit1 func(O1),
	emit2 func(O2),
)

func NewFuncStreamNode2x2[I1 any, I2 any, O1 any, O2 any](f FuncStreamNode2x2[I1, I2, O1, O2]) StreamNode2x2[I1, I2, O1, O2] {
	return &funcStreamNode2x2[I1, I2, O1, O2]{Func: f}
}

type funcStreamNode2x2[I1 any, I2 any, O1 any, O2 any] struct {
	StreamNode2x2[I1, I2, O1, O2]

	Func FuncStreamNode2x2[I1, I2, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Inputs() int {
	return 2
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Outputs() int {
	return 2
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas2In
	NodeHas3Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1, v2 I2)
}

type FuncNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] func(v1 I1, v2 I2) (O1, O2, O3)

func NewFuncNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any](f FuncNode2x3[I1, I2, O1, O2, O3]) Node2x3[I1, I2, O1, O2, O3] {
	return &funcNode2x3[I1, I2, O1, O2, O3]{
		Func: f,
	}
}

type funcNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] struct {
	Node2x3[I1, I2, O1, O2, O3]

	Func FuncNode2x3[I1, I2, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode2x3[I1, I2, O1, O2, O3]) Inputs() int {
	return 2
}

func (f *funcNode2x3[I1, I2, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcNode2x3[I1, I2, O1, O2, O3]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcNode2x3[I1, I2, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode2x3[I1, I2, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	v1, v2, v3 := f.Func(i1, i2)
	emit(0, v1)

	emit(1, v2)

	emit(2, v3)
}

type StreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] interface {
	Node
	NodeHas2In
	NodeHas3Out

	NodeIn1[I1]
	NodeIn2[I2]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1, v2 I2)
}

type FuncStreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] func(
	_ I1,
	_ I2,
	emit1 func(O1),
	emit2 func(O2),
	emit3 func(O3),
)

func NewFuncStreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any](f FuncStreamNode2x3[I1, I2, O1, O2, O3]) StreamNode2x3[I1, I2, O1, O2, O3] {
	return &funcStreamNode2x3[I1, I2, O1, O2, O3]{Func: f}
}

type funcStreamNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] struct {
	StreamNode2x3[I1, I2, O1, O2, O3]

	Func FuncStreamNode2x3[I1, I2, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Inputs() int {
	return 2
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Run(v1 I1, v2 I2) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
	)
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas3In
	NodeHas0Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncNode3x0[I1 any, I2 any, I3 any] func(v1 I1, v2 I2, v3 I3)

func NewFuncNode3x0[I1 any, I2 any, I3 any](f FuncNode3x0[I1, I2, I3]) Node3x0[I1, I2, I3] {
	return &funcNode3x0[I1, I2, I3]{
		Func: f,
	}
}

type funcNode3x0[I1 any, I2 any, I3 any] struct {
	Node3x0[I1, I2, I3]

	Func FuncNode3x0[I1, I2, I3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode3x0[I1, I2, I3]) Inputs() int {
	return 3
}

func (f *funcNode3x0[I1, I2, I3]) Outputs() int {
	return 0
}

func (f *funcNode3x0[I1, I2, I3]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcNode3x0[I1, I2, I3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode3x0[I1, I2, I3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)
	f.Func(i1, i2, i3)
}

type StreamNode3x0[I1 any, I2 any, I3 any] interface {
	Node
	NodeHas3In
	NodeHas0Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncStreamNode3x0[I1 any, I2 any, I3 any] func(
	_ I1,
	_ I2,
	_ I3,
)

func NewFuncStreamNode3x0[I1 any, I2 any, I3 any](f FuncStreamNode3x0[I1, I2, I3]) StreamNode3x0[I1, I2, I3] {
	return &funcStreamNode3x0[I1, I2, I3]{Func: f}
}

type funcStreamNode3x0[I1 any, I2 any, I3 any] struct {
	StreamNode3x0[I1, I2, I3]

	Func FuncStreamNode3x0[I1, I2, I3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode3x0[I1, I2, I3]) Inputs() int {
	return 3
}

func (f *funcStreamNode3x0[I1, I2, I3]) Outputs() int {
	return 0
}

func (f *funcStreamNode3x0[I1, I2, I3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode3x0[I1, I2, I3]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcStreamNode3x0[I1, I2, I3]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas3In
	NodeHas1Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncNode3x1[I1 any, I2 any, I3 any, O1 any] func(v1 I1, v2 I2, v3 I3) O1

func NewFuncNode3x1[I1 any, I2 any, I3 any, O1 any](f FuncNode3x1[I1, I2, I3, O1]) Node3x1[I1, I2, I3, O1] {
	return &funcNode3x1[I1, I2, I3, O1]{
		Func: f,
	}
}

type funcNode3x1[I1 any, I2 any, I3 any, O1 any] struct {
	Node3x1[I1, I2, I3, O1]

	Func FuncNode3x1[I1, I2, I3, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode3x1[I1, I2, I3, O1]) Inputs() int {
	return 3
}

func (f *funcNode3x1[I1, I2, I3, O1]) Outputs() int {
	return 1
}

func (f *funcNode3x1[I1, I2, I3, O1]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcNode3x1[I1, I2, I3, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode3x1[I1, I2, I3, O1]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)
	v1 := f.Func(i1, i2, i3)
	emit(0, v1)
}

type StreamNode3x1[I1 any, I2 any, I3 any, O1 any] interface {
	Node
	NodeHas3In
	NodeHas1Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncStreamNode3x1[I1 any, I2 any, I3 any, O1 any] func(
	_ I1,
	_ I2,
	_ I3,
	emit1 func(O1),
)

func NewFuncStreamNode3x1[I1 any, I2 any, I3 any, O1 any](f FuncStreamNode3x1[I1, I2, I3, O1]) StreamNode3x1[I1, I2, I3, O1] {
	return &funcStreamNode3x1[I1, I2, I3, O1]{Func: f}
}

type funcStreamNode3x1[I1 any, I2 any, I3 any, O1 any] struct {
	StreamNode3x1[I1, I2, I3, O1]

	Func FuncStreamNode3x1[I1, I2, I3, O1]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Inputs() int {
	return 3
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Outputs() int {
	return 1
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas3In
	NodeHas2Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] func(v1 I1, v2 I2, v3 I3) (O1, O2)

func NewFuncNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any](f FuncNode3x2[I1, I2, I3, O1, O2]) Node3x2[I1, I2, I3, O1, O2] {
	return &funcNode3x2[I1, I2, I3, O1, O2]{
		Func: f,
	}
}

type funcNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] struct {
	Node3x2[I1, I2, I3, O1, O2]

	Func FuncNode3x2[I1, I2, I3, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode3x2[I1, I2, I3, O1, O2]) Inputs() int {
	return 3
}

func (f *funcNode3x2[I1, I2, I3, O1, O2]) Outputs() int {
	return 2
}

func (f *funcNode3x2[I1, I2, I3, O1, O2]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcNode3x2[I1, I2, I3, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode3x2[I1, I2, I3, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)
	v1, v2 := f.Func(i1, i2, i3)
	emit(0, v1)

	emit(1, v2)
}

type StreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] interface {
	Node
	NodeHas3In
	NodeHas2Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncStreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] func(
	_ I1,
	_ I2,
	_ I3,
	emit1 func(O1),
	emit2 func(O2),
)

func NewFuncStreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any](f FuncStreamNode3x2[I1, I2, I3, O1, O2]) StreamNode3x2[I1, I2, I3, O1, O2] {
	return &funcStreamNode3x2[I1, I2, I3, O1, O2]{Func: f}
}

type funcStreamNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] struct {
	StreamNode3x2[I1, I2, I3, O1, O2]

	Func FuncStreamNode3x2[I1, I2, I3, O1, O2]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Inputs() int {
	return 3
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Outputs() int {
	return 2
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Do(inputs []any, emit func(i int, v any)) {
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
	NodeHas3In
	NodeHas3Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] func(v1 I1, v2 I2, v3 I3) (O1, O2, O3)

func NewFuncNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any](f FuncNode3x3[I1, I2, I3, O1, O2, O3]) Node3x3[I1, I2, I3, O1, O2, O3] {
	return &funcNode3x3[I1, I2, I3, O1, O2, O3]{
		Func: f,
	}
}

type funcNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] struct {
	Node3x3[I1, I2, I3, O1, O2, O3]

	Func FuncNode3x3[I1, I2, I3, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcNode3x3[I1, I2, I3, O1, O2, O3]) Inputs() int {
	return 3
}

func (f *funcNode3x3[I1, I2, I3, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcNode3x3[I1, I2, I3, O1, O2, O3]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcNode3x3[I1, I2, I3, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcNode3x3[I1, I2, I3, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
	i1 := inputs[0].(I1)
	i2 := inputs[1].(I2)
	i3 := inputs[2].(I3)
	v1, v2, v3 := f.Func(i1, i2, i3)
	emit(0, v1)

	emit(1, v2)

	emit(2, v3)
}

type StreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] interface {
	Node
	NodeHas3In
	NodeHas3Out

	NodeIn1[I1]
	NodeIn2[I2]
	NodeIn3[I3]

	NodeOut1[O1]
	NodeOut2[O2]
	NodeOut3[O3]

	Run(v1 I1, v2 I2, v3 I3)
}

type FuncStreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] func(
	_ I1,
	_ I2,
	_ I3,
	emit1 func(O1),
	emit2 func(O2),
	emit3 func(O3),
)

func NewFuncStreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any](f FuncStreamNode3x3[I1, I2, I3, O1, O2, O3]) StreamNode3x3[I1, I2, I3, O1, O2, O3] {
	return &funcStreamNode3x3[I1, I2, I3, O1, O2, O3]{Func: f}
}

type funcStreamNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] struct {
	StreamNode3x3[I1, I2, I3, O1, O2, O3]

	Func FuncStreamNode3x3[I1, I2, I3, O1, O2, O3]

	machineryOnce sync.Once
	machinery     *NodeMachinery
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Inputs() int {
	return 3
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Outputs() int {
	return 3
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Machinery() *NodeMachinery {
	f.machineryOnce.Do(func() {
		f.machinery = NewNodeMachinery(f)
	})

	return f.machinery
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Run(v1 I1, v2 I2, v3 I3) {
	f.Machinery().NewSourceRun(
		v1,
		v2,
		v3,
	)
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Do(inputs []any, emit func(i int, v any)) {
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
