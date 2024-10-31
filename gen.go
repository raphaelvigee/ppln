package ppln

import "sync"

type NodeHas0In interface {
	_input_layout()
}

type NodeIn1[T any] interface {
	_input1(T)
}

type NodeHas1In interface {
	_input_layout(any)
}

type NodeIn2[T any] interface {
	_input2(T)
}

type NodeHas2In interface {
	_input_layout(any, any)
}

type NodeIn3[T any] interface {
	_input3(T)
}

type NodeHas3In interface {
	_input_layout(any, any, any)
}

type NodeHas0Out interface {
	_out_layout()
}

type NodeOut1[T any] interface {
	_output1(T)
}

type NodeHas1Out interface {
	_out_layout(any)
}

type NodeOut2[T any] interface {
	_output2(T)
}

type NodeHas2Out interface {
	_out_layout(any, any)
}

type NodeOut3[T any] interface {
	_output3(T)
}

type NodeHas3Out interface {
	_out_layout(any, any, any)
}

func Take1[T any](n interface {
	Node
	NodeOut1[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 1-1)
}

func Pipeline1[T1 any](
	from1 interface {
		Node
		NodeOut1[T1]
		NodeHas1Out
	},
	to interface {
		Node
		NodeHas1In
		NodeIn1[T1]
	},
) {
	Pipeline(
		to,
		from1,
	)
}

func Take2[T any](n interface {
	Node
	NodeOut2[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 2-1)
}

func Pipeline2[T1 any, T2 any](
	from1 interface {
		Node
		NodeOut1[T1]
		NodeHas1Out
	}, from2 interface {
		Node
		NodeOut1[T2]
		NodeHas1Out
	},
	to interface {
		Node
		NodeHas2In
		NodeIn1[T1]
		NodeIn2[T2]
	},
) {
	Pipeline(
		to,
		from1,
		from2,
	)
}

func Take3[T any](n interface {
	Node
	NodeOut3[T]
}) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return TakeN[T](n, 3-1)
}

func Pipeline3[T1 any, T2 any, T3 any](
	from1 interface {
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
	},
	to interface {
		Node
		NodeHas3In
		NodeIn1[T1]
		NodeIn2[T2]
		NodeIn3[T3]
	},
) {
	Pipeline(
		to,
		from1,
		from2,
		from3,
	)
}

type FuncNode0x1[O1 any] func() O1

func NewFuncNode0x1[O1 any](f FuncNode0x1[O1]) StreamNode0x1[O1] {
	return NewFuncStreamNode0x1(func(
		emit1 func(*LineageRef, O1),
	) {
		v1 := f()

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
	})
}

type StreamNode0x1[O1 any] interface {
	Node
	NodeHas0In
	NodeHas1Out

	NodeOut1[O1]

	Run()
}

type FuncStreamNode0x1[O1 any] func(
	emit1 func(*LineageRef, O1),
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

func (f *funcStreamNode0x1[O1]) Inputs() uint8 {
	return 0
}

func (f *funcStreamNode0x1[O1]) Outputs() uint8 {
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

func (f *funcStreamNode0x1[O1]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {

	f.Func(
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
	)
}

type FuncNode0x2[O1 any, O2 any] func() (O1, O2)

func NewFuncNode0x2[O1 any, O2 any](f FuncNode0x2[O1, O2]) StreamNode0x2[O1, O2] {
	return NewFuncStreamNode0x2(func(
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
	) {
		v1, v2 := f()

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
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

func (f *funcStreamNode0x2[O1, O2]) Inputs() uint8 {
	return 0
}

func (f *funcStreamNode0x2[O1, O2]) Outputs() uint8 {
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

func (f *funcStreamNode0x2[O1, O2]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {

	f.Func(
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
	)
}

type FuncNode0x3[O1 any, O2 any, O3 any] func() (O1, O2, O3)

func NewFuncNode0x3[O1 any, O2 any, O3 any](f FuncNode0x3[O1, O2, O3]) StreamNode0x3[O1, O2, O3] {
	return NewFuncStreamNode0x3(func(
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
		emit3 func(*LineageRef, O3),
	) {
		v1, v2, v3 := f()

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
		emit3(l, v3)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
	emit3 func(*LineageRef, O3),
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

func (f *funcStreamNode0x3[O1, O2, O3]) Inputs() uint8 {
	return 0
}

func (f *funcStreamNode0x3[O1, O2, O3]) Outputs() uint8 {
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

func (f *funcStreamNode0x3[O1, O2, O3]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {

	f.Func(
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
		func(l *LineageRef, v O3) {
			emit(2, l, v)
		},
	)
}

type FuncNode1x0[I1 any] func(v1 I1)

func NewFuncNode1x0[I1 any](f FuncNode1x0[I1]) StreamNode1x0[I1] {
	return NewFuncStreamNode1x0(func(
		i1 I1,
	) {
		f(i1)
	})
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

func (f *funcStreamNode1x0[I1]) Inputs() uint8 {
	return 1
}

func (f *funcStreamNode1x0[I1]) Outputs() uint8 {
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

func (f *funcStreamNode1x0[I1]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])

	f.Func(
		i1,
	)
}

type FuncNode1x1[I1 any, O1 any] func(v1 I1) O1

func NewFuncNode1x1[I1 any, O1 any](f FuncNode1x1[I1, O1]) StreamNode1x1[I1, O1] {
	return NewFuncStreamNode1x1(func(
		i1 I1,
		emit1 func(*LineageRef, O1),
	) {
		v1 := f(i1)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
	})
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
	emit1 func(*LineageRef, O1),
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

func (f *funcStreamNode1x1[I1, O1]) Inputs() uint8 {
	return 1
}

func (f *funcStreamNode1x1[I1, O1]) Outputs() uint8 {
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

func (f *funcStreamNode1x1[I1, O1]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])

	f.Func(
		i1,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
	)
}

type FuncNode1x2[I1 any, O1 any, O2 any] func(v1 I1) (O1, O2)

func NewFuncNode1x2[I1 any, O1 any, O2 any](f FuncNode1x2[I1, O1, O2]) StreamNode1x2[I1, O1, O2] {
	return NewFuncStreamNode1x2(func(
		i1 I1,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
	) {
		v1, v2 := f(i1)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
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

func (f *funcStreamNode1x2[I1, O1, O2]) Inputs() uint8 {
	return 1
}

func (f *funcStreamNode1x2[I1, O1, O2]) Outputs() uint8 {
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

func (f *funcStreamNode1x2[I1, O1, O2]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])

	f.Func(
		i1,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
	)
}

type FuncNode1x3[I1 any, O1 any, O2 any, O3 any] func(v1 I1) (O1, O2, O3)

func NewFuncNode1x3[I1 any, O1 any, O2 any, O3 any](f FuncNode1x3[I1, O1, O2, O3]) StreamNode1x3[I1, O1, O2, O3] {
	return NewFuncStreamNode1x3(func(
		i1 I1,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
		emit3 func(*LineageRef, O3),
	) {
		v1, v2, v3 := f(i1)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
		emit3(l, v3)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
	emit3 func(*LineageRef, O3),
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

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Inputs() uint8 {
	return 1
}

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Outputs() uint8 {
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

func (f *funcStreamNode1x3[I1, O1, O2, O3]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])

	f.Func(
		i1,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
		func(l *LineageRef, v O3) {
			emit(2, l, v)
		},
	)
}

type FuncNode2x0[I1 any, I2 any] func(v1 I1, v2 I2)

func NewFuncNode2x0[I1 any, I2 any](f FuncNode2x0[I1, I2]) StreamNode2x0[I1, I2] {
	return NewFuncStreamNode2x0(func(
		i1 I1,
		i2 I2,
	) {
		f(i1, i2)
	})
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

func (f *funcStreamNode2x0[I1, I2]) Inputs() uint8 {
	return 2
}

func (f *funcStreamNode2x0[I1, I2]) Outputs() uint8 {
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

func (f *funcStreamNode2x0[I1, I2]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])

	f.Func(
		i1,
		i2,
	)
}

type FuncNode2x1[I1 any, I2 any, O1 any] func(v1 I1, v2 I2) O1

func NewFuncNode2x1[I1 any, I2 any, O1 any](f FuncNode2x1[I1, I2, O1]) StreamNode2x1[I1, I2, O1] {
	return NewFuncStreamNode2x1(func(
		i1 I1,
		i2 I2,
		emit1 func(*LineageRef, O1),
	) {
		v1 := f(i1, i2)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
	})
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
	emit1 func(*LineageRef, O1),
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

func (f *funcStreamNode2x1[I1, I2, O1]) Inputs() uint8 {
	return 2
}

func (f *funcStreamNode2x1[I1, I2, O1]) Outputs() uint8 {
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

func (f *funcStreamNode2x1[I1, I2, O1]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])

	f.Func(
		i1,
		i2,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
	)
}

type FuncNode2x2[I1 any, I2 any, O1 any, O2 any] func(v1 I1, v2 I2) (O1, O2)

func NewFuncNode2x2[I1 any, I2 any, O1 any, O2 any](f FuncNode2x2[I1, I2, O1, O2]) StreamNode2x2[I1, I2, O1, O2] {
	return NewFuncStreamNode2x2(func(
		i1 I1,
		i2 I2,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
	) {
		v1, v2 := f(i1, i2)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
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

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Inputs() uint8 {
	return 2
}

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Outputs() uint8 {
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

func (f *funcStreamNode2x2[I1, I2, O1, O2]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])

	f.Func(
		i1,
		i2,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
	)
}

type FuncNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any] func(v1 I1, v2 I2) (O1, O2, O3)

func NewFuncNode2x3[I1 any, I2 any, O1 any, O2 any, O3 any](f FuncNode2x3[I1, I2, O1, O2, O3]) StreamNode2x3[I1, I2, O1, O2, O3] {
	return NewFuncStreamNode2x3(func(
		i1 I1,
		i2 I2,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
		emit3 func(*LineageRef, O3),
	) {
		v1, v2, v3 := f(i1, i2)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
		emit3(l, v3)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
	emit3 func(*LineageRef, O3),
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

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Inputs() uint8 {
	return 2
}

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Outputs() uint8 {
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

func (f *funcStreamNode2x3[I1, I2, O1, O2, O3]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])

	f.Func(
		i1,
		i2,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
		func(l *LineageRef, v O3) {
			emit(2, l, v)
		},
	)
}

type FuncNode3x0[I1 any, I2 any, I3 any] func(v1 I1, v2 I2, v3 I3)

func NewFuncNode3x0[I1 any, I2 any, I3 any](f FuncNode3x0[I1, I2, I3]) StreamNode3x0[I1, I2, I3] {
	return NewFuncStreamNode3x0(func(
		i1 I1,
		i2 I2,
		i3 I3,
	) {
		f(i1, i2, i3)
	})
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

func (f *funcStreamNode3x0[I1, I2, I3]) Inputs() uint8 {
	return 3
}

func (f *funcStreamNode3x0[I1, I2, I3]) Outputs() uint8 {
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

func (f *funcStreamNode3x0[I1, I2, I3]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])
	i3 := Cast[I3](inputs[2])

	f.Func(
		i1,
		i2,
		i3,
	)
}

type FuncNode3x1[I1 any, I2 any, I3 any, O1 any] func(v1 I1, v2 I2, v3 I3) O1

func NewFuncNode3x1[I1 any, I2 any, I3 any, O1 any](f FuncNode3x1[I1, I2, I3, O1]) StreamNode3x1[I1, I2, I3, O1] {
	return NewFuncStreamNode3x1(func(
		i1 I1,
		i2 I2,
		i3 I3,
		emit1 func(*LineageRef, O1),
	) {
		v1 := f(i1, i2, i3)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
	})
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
	emit1 func(*LineageRef, O1),
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

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Inputs() uint8 {
	return 3
}

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Outputs() uint8 {
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

func (f *funcStreamNode3x1[I1, I2, I3, O1]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])
	i3 := Cast[I3](inputs[2])

	f.Func(
		i1,
		i2,
		i3,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
	)
}

type FuncNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any] func(v1 I1, v2 I2, v3 I3) (O1, O2)

func NewFuncNode3x2[I1 any, I2 any, I3 any, O1 any, O2 any](f FuncNode3x2[I1, I2, I3, O1, O2]) StreamNode3x2[I1, I2, I3, O1, O2] {
	return NewFuncStreamNode3x2(func(
		i1 I1,
		i2 I2,
		i3 I3,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
	) {
		v1, v2 := f(i1, i2, i3)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
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

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Inputs() uint8 {
	return 3
}

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Outputs() uint8 {
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

func (f *funcStreamNode3x2[I1, I2, I3, O1, O2]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])
	i3 := Cast[I3](inputs[2])

	f.Func(
		i1,
		i2,
		i3,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
	)
}

type FuncNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any] func(v1 I1, v2 I2, v3 I3) (O1, O2, O3)

func NewFuncNode3x3[I1 any, I2 any, I3 any, O1 any, O2 any, O3 any](f FuncNode3x3[I1, I2, I3, O1, O2, O3]) StreamNode3x3[I1, I2, I3, O1, O2, O3] {
	return NewFuncStreamNode3x3(func(
		i1 I1,
		i2 I2,
		i3 I3,
		emit1 func(*LineageRef, O1),
		emit2 func(*LineageRef, O2),
		emit3 func(*LineageRef, O3),
	) {
		v1, v2, v3 := f(i1, i2, i3)

		l := SameLineageRef()
		defer l.Done()
		emit1(l, v1)
		emit2(l, v2)
		emit3(l, v3)
	})
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
	emit1 func(*LineageRef, O1),
	emit2 func(*LineageRef, O2),
	emit3 func(*LineageRef, O3),
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

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Inputs() uint8 {
	return 3
}

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Outputs() uint8 {
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

func (f *funcStreamNode3x3[I1, I2, I3, O1, O2, O3]) Do(inputs []any, emit func(i uint8, l *LineageRef, v any)) {
	i1 := Cast[I1](inputs[0])
	i2 := Cast[I2](inputs[1])
	i3 := Cast[I3](inputs[2])

	f.Func(
		i1,
		i2,
		i3,
		func(l *LineageRef, v O1) {
			emit(0, l, v)
		},
		func(l *LineageRef, v O2) {
			emit(1, l, v)
		},
		func(l *LineageRef, v O3) {
			emit(2, l, v)
		},
	)
}
