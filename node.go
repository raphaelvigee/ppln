package ppln

import (
	"context"
	"fmt"
	"github.com/dlsniper/debugger"
	"slices"
	"sync"
	"sync/atomic"
)

type Node interface {
	Inputs() int
	Outputs() int

	Do(values []any, emit func(i int, v any))
	Machinery() *NodeMachinery
}

type Edge struct {
	From Node
	To   Node
}

type Value struct {
	Value any
	Index int
	ID    uint64
}

type NodeValueMapper interface {
	ValueMapper(v Value) []Value
}

type MapperNode[T any] struct {
	NodeOut1[T]
	NodeHas1Out
	Node

	Func func(v Value) []Value
}

var _ Node = (*MapperNode[any])(nil)

func (t MapperNode[T]) ValueMapper(v Value) []Value {
	return t.Func(v)
}

type MappableNode[T any] interface {
	Node
	NodeOut1[T]
	NodeHas1Out
}

func NewMapperNode[T any](node Node, f func(v Value) []Value) MappableNode[T] {
	if m, ok := node.(NodeValueMapper); ok {
		newf := f
		f = func(v Value) []Value {
			vs := m.ValueMapper(v)

			if len(vs) == 0 {
				return nil
			}

			var outvs []Value
			for _, v := range vs {
				outvs = append(outvs, newf(v)...)
			}

			return vs
		}
	}

	return &MapperNode[T]{
		Node: node,
		Func: f,
	}
}

func TakeN[T any](node Node, n int) MappableNode[T] {
	return NewMapperNode[T](node, func(v Value) []Value {
		if v.Index != n {
			return nil
		}

		v.Index = 0

		return []Value{v}
	})
}

func Filter[T any](node Node, f func(v T) bool) MappableNode[T] {
	return NewMapperNode[T](node, func(v Value) []Value {
		if f(CastInput[T](v.Value)) {
			return []Value{v}
		}

		return nil
	})
}

func FilterNode[T any](f func(v T) bool) StreamNode1x1[T, T] {
	return NewFuncStreamNode1x1[T, T](func(v T, emit1 func(T)) {
		if f(v) {
			emit1(v)
		}
	})
}

type Machinery interface {
	Machinery() *NodeMachinery
}

//goland:noinspection GoVetLostCancel
func Pipeline(to Machinery, from ...interface {
	Machinery
	NodeHas1Out
}) {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	for i, from := range from {
		wg.Add(1)

		go func() {
			debugger.SetLabels(func() []string {
				return []string{"where", fmt.Sprintf("Pipeline %p (%v) -> %p", from, i, to)}
			})

			defer cancel()

			ch := from.Machinery().Listen(ctx)

			wg.Done()

			for v := range ch {
				if m, ok := from.(NodeValueMapper); ok {
					vs := m.ValueMapper(v)

					for _, v := range vs {
						to.Machinery().Incoming(i, v)
					}
				} else {
					to.Machinery().Incoming(i, v)
				}

			}
		}()
	}

	wg.Wait()
}

func NewNodeMachinery(n Node) *NodeMachinery {
	return &NodeMachinery{
		n:            n,
		stagedValues: map[uint64][]Value{},
	}
}

type Listener func(i int, v Value)

type listenerContainer struct {
	ID   uint64
	Func Listener
}

type NodeMachinery struct {
	n Node

	listenerc atomic.Uint64
	listeners []listenerContainer
	listenerm sync.Mutex

	valueMiddleware Listener

	stagedValuesm sync.Mutex
	stagedValues  map[uint64][]Value
}

var sourceId atomic.Uint64

func (m *NodeMachinery) NewSourceRun(values ...any) {
	id := sourceId.Add(1)

	args := make([]Value, len(values))
	for i, v := range values {
		args[i] = Value{
			Value: v,
			Index: i,
			ID:    id,
		}
	}

	m.Run(args...)
}

func (m *NodeMachinery) Run(values ...Value) {
	args := make([]any, len(values))
	for i, v := range values {
		args[i] = v.Value
	}

	var id uint64
	if len(values) > 0 {
		id = values[0].ID // TODO compose ids together from input
	} else {
		id = sourceId.Add(1)
	}

	m.n.Do(args, func(i int, v any) {
		m.listenerm.Lock()
		defer m.listenerm.Unlock()

		for _, l := range m.listeners {
			l.Func(i, Value{
				Value: v,
				Index: i,
				ID:    id,
			})
		}
	})
}

func (m *NodeMachinery) Listen(ctx context.Context) chan Value {
	id := m.listenerc.Add(1)
	if id == 0 {
		panic("listener id wrapped around")
	}

	ch := make(chan Value)

	m.listenerm.Lock()
	m.listeners = append(m.listeners, listenerContainer{
		ID: id,
		Func: func(i int, v Value) {
			ch <- v
		},
	})
	m.listenerm.Unlock()

	go func() {
		debugger.SetLabels(func() []string {
			return []string{"where", fmt.Sprintf("Listen wait close: %p", m)}
		})

		<-ctx.Done()
		close(ch)

		m.listenerm.Lock()
		m.listeners = slices.DeleteFunc(m.listeners, func(l listenerContainer) bool {
			return l.ID == id
		})
		m.listenerm.Unlock()
	}()

	return ch
}

func (m *NodeMachinery) Incoming(i int, v Value) {
	values, do := m.processIncoming(i, v)
	if !do {
		return
	}

	m.Run(values...)
}

func (m *NodeMachinery) processIncoming(i int, v Value) ([]Value, bool) {
	if m.n.Inputs() == 1 {
		return []Value{v}, true
	}

	m.stagedValuesm.Lock()
	defer m.stagedValuesm.Unlock()

	values, ok := m.stagedValues[v.ID]
	if !ok {
		values = make([]Value, m.n.Inputs())
	}

	values[i] = v

	allReady := true
	for _, v := range values {
		if v.ID == 0 {
			allReady = false
			break
		}
	}

	if allReady {
		delete(m.stagedValues, v.ID)

		return values, true
	}

	m.stagedValues[v.ID] = values

	return nil, false
}

func CastInput[T any](v any) T {
	if v == nil {
		var zero T

		return zero
	}

	v2, ok := v.(T)
	if !ok {
		_ = v.(T)
	}

	return v2
}
