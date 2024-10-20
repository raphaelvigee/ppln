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
	is_node()
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

type TakeNode[T any] struct {
	NodeOut1[T]
	NodeHas1Out
	Node
	N int
}

func (t TakeNode[T]) ValueMapper(v Value) []Value {
	if v.Index != t.N {
		return nil
	}

	v.Index = 0

	return []Value{v}
}

var _ Node = (*TakeNode[any])(nil)

func TakeN[T any](node Node, n int) interface {
	Node
	NodeOut1[T]
	NodeHas1Out
} {
	return &TakeNode[T]{
		Node: node,
		N:    n,
	}
}

//goland:noinspection GoVetLostCancel
func Pipeline(to Node, from ...interface {
	Node
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

var source atomic.Uint64

func (m *NodeMachinery) NewSourceRun(values ...any) {
	id := source.Add(1)

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
		id = source.Add(1)
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
