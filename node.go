package ppln

import (
	"context"
	"fmt"
	"github.com/dlsniper/debugger"
	"maps"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type LineageRef struct {
	lineage LineageID
	m       sync.Mutex
}

func NewLineageRef() *LineageRef {
	return &LineageRef{}
}

func (l *LineageRef) get(id LineageID) LineageID {
	l.m.Lock()
	defer l.m.Unlock()

	if len(l.lineage) > 0 {
		return l.lineage
	}

	l.lineage = append(id, sourceId.Add(1))

	return l.lineage
}

type Node interface {
	Inputs() uint8
	Outputs() uint8

	Do(values []any, emit func(i uint8, l *LineageRef, v any))
	Machinery() *NodeMachinery
}

type Edge struct {
	From Node
	To   Node
}

type LineageID []uint64

func (self LineageID) Correlates(other LineageID) bool {
	return self.correlates(other) || other.correlates(self)
}

func (self LineageID) correlates(other LineageID) bool {
	if len(self) > len(other) {
		return false
	}

	return slices.Equal(self, other[:len(self)])
}

type Value struct {
	Value   any
	Index   uint8
	Lineage LineageID
}

func (v Value) IsSet() bool {
	return len(v.Lineage) > 0
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

func TakeN[T any](node Node, n uint8) MappableNode[T] {
	return NewMapperNode[T](node, func(v Value) []Value {
		if v.Index != n {
			return nil
		}

		v.Index = 0

		return []Value{v}
	})
}

func Filter[T any](f func(v T) bool) StreamNode1x1[T, T] {
	return NewFuncStreamNode1x1[T, T](func(v T, emit1 func(*LineageRef, T)) {
		if f(v) {
			emit1(nil, v)
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
		n: n,
		stagedValues: valueStore{
			mem:    map[string]Value{},
			staged: nil,
		},
		id: sourceId.Add(1),
	}
}

type Listener func(i uint8, v Value)

type listenerContainer struct {
	ID   uint64
	Func Listener
}

type NodeMachinery struct {
	n  Node
	id uint64

	listenerc atomic.Uint64
	listeners []listenerContainer
	listenerm sync.RWMutex

	valueMiddleware Listener

	stagedValuesm sync.Mutex
	stagedValues  valueStore
}

var sourceId atomic.Uint64

func (m *NodeMachinery) NewSourceRun(values ...any) {
	id := LineageID{sourceId.Add(1)}

	args := make([]Value, len(values))
	for i, v := range values {
		args[i] = Value{
			Value:   v,
			Index:   uint8(i),
			Lineage: id,
		}
	}

	m.Run(args...)
}

func (m *NodeMachinery) Run(inputs ...Value) {
	args := make([]any, len(inputs))
	for i, v := range inputs {
		args[i] = v.Value
	}

	var id LineageID
	if m.n.Outputs() > 0 { // if we gonna need an id at some point
		switch len(inputs) {
		case 1:
			id = inputs[0].Lineage
			if len(id) == 0 {
				panic("empty lineage")
			}
		default:
			idm := map[uint64]struct{}{}
			for _, value := range inputs {
				for _, id := range value.Lineage {
					idm[id] = struct{}{}
				}
			}
			if len(idm) == 0 {
				idm[sourceId.Add(1)] = struct{}{}
			}

			id = slices.Collect(maps.Keys(idm))
			slices.Sort(id)

			if len(id) == 0 {
				panic("empty lineage")
			}
		}
	}

	m.n.Do(args, func(i uint8, l *LineageRef, v any) {
		if i >= m.n.Outputs() {
			panic(fmt.Sprintf("trying to output %v, but node has %v outputs", i, m.n.Outputs()))
		}

		if len(id) == 0 {
			panic("empty lineage")
		}

		emitId := id
		if l != nil {
			emitId = l.get(id)
		}

		m.listenerm.RLock()
		defer m.listenerm.RUnlock()

		for _, l := range m.listeners {
			l.Func(i, Value{
				Value:   v,
				Index:   i,
				Lineage: emitId,
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
		Func: func(i uint8, v Value) {
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

type valueStore struct {
	mem    map[string]Value
	staged [][]Value
}

func (s *valueStore) strid(id LineageID) string {
	ids := make([]string, len(id))
	for i, v := range id {
		ids[i] = fmt.Sprint(v)
	}

	return strings.Join(ids, ",")
}

func (s *valueStore) set(incoming Value, cap uint8) []Value {
	s.mem[s.strid(incoming.Lineage)] = incoming

	for stagedIdx, inputs := range s.staged {
		valueAtIndex := inputs[incoming.Index]

		if !valueAtIndex.IsSet() {
			allSet := true
			compatible := true
			for _, value := range inputs {
				if !value.IsSet() {
					if value.Index != incoming.Index {
						allSet = false
					}
				} else {
					if !incoming.Lineage.Correlates(value.Lineage) {
						compatible = false
						break
					}
				}
			}

			if compatible {
				if allSet {
					inputs[incoming.Index] = incoming

					s.staged = slices.Delete(s.staged, stagedIdx, stagedIdx+1)

					return inputs
				} else {
					inputs[incoming.Index] = incoming

					s.staged[stagedIdx] = inputs

					return nil
				}
			}
		}
	}

	values := make([]Value, cap)
	for i, value := range values {
		value.Index = uint8(i)
		values[i] = value
	}
	values[incoming.Index] = incoming

	for _, value := range s.mem {
		if values[value.Index].IsSet() {
			continue
		}

		if value.Lineage.Correlates(incoming.Lineage) {
			values[value.Index] = value
		}
	}

	allSet := true
	for _, value := range values {
		if !value.IsSet() {
			allSet = false
			break
		}
	}

	if allSet {
		return values
	}

	s.staged = append(s.staged, values)

	return nil
}

func (m *NodeMachinery) processIncoming(i int, v Value) ([]Value, bool) {
	v.Index = uint8(i)

	if m.n.Inputs() == 1 {
		return []Value{v}, true
	}

	m.stagedValuesm.Lock()
	defer m.stagedValuesm.Unlock()

	values := m.stagedValues.set(v, m.n.Inputs())

	return values, len(values) > 0
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
