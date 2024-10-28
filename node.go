package ppln

import (
	"context"
	"fmt"
	"github.com/dlsniper/debugger"
	"maps"
	"slices"
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

var pipeId atomic.Uint64

func Pipeline(to Machinery, from ...interface {
	Machinery
	NodeHas1Out
}) {
	ctx := context.Background()

	pipeId := pipeId.Add(1)

	var startedWg sync.WaitGroup

	for i, from := range from {
		startedWg.Add(1)

		go func() {
			debugger.SetLabels(func() []string {
				return []string{"where", fmt.Sprintf("Pipeline %p (%v) -> %p", from, i, to)}
			})

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			ch := from.Machinery().Listen(ctx)

			startedWg.Done()

			for v := range ch {
				if _, ok := v.Value.(SourceDone); ok {
					continue
				}

				if m, ok := from.(NodeValueMapper); ok {
					vs := m.ValueMapper(v)

					for _, v := range vs {
						v.Index = uint8(i)
						to.Machinery().Incoming(pipeId, v)
					}
				} else {
					v.Index = uint8(i)
					to.Machinery().Incoming(pipeId, v)
				}
			}
		}()
	}

	startedWg.Wait()
}

func NewNodeMachinery(n Node) *NodeMachinery {
	return &NodeMachinery{
		n: n,
		stagedValues: valueStore{
			mem: NewIDTrie(int(n.Inputs())),
		},
		id: sourceId.Add(1),
	}
}

type Listener func(v Value)

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

	stagedValues valueStore
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

func (m *NodeMachinery) emit(v Value) {
	m.listenerm.RLock()
	defer m.listenerm.RUnlock()

	for _, l := range m.listeners {
		l.Func(v)
	}
}

type SourceDone struct{}

func (m *NodeMachinery) Run(inputs ...Value) {
	var id LineageID
	switch len(inputs) {
	case 1:
		id = inputs[0].Lineage
		if len(id) == 0 {
			panic("empty lineage")
		}
	default:
		var pipeId uint64
		idm := map[uint64]struct{}{}
		for _, value := range inputs {
			for i, id := range value.Lineage {
				if i == 0 {
					if pipeId == 0 {
						pipeId = id
					} else {
						if pipeId != id {
							panic(fmt.Sprintf("mismatch pipeid: %v %v", pipeId, id))
						}
					}

					continue
				}

				idm[id] = struct{}{}
			}
		}

		id = make(LineageID, 1, len(idm)+1)
		id[0] = pipeId
		for idk := range maps.Keys(idm) {
			id = append(id, idk)
		}
		slices.Sort(id[1:])

		if len(id) == 1 {
			id = append(id, sourceId.Add(1))
		}
	}

	args := make([]any, len(inputs))
	for i, v := range inputs {
		args[i] = v.Value
	}

	defer func() {
		m.emit(Value{
			Value:   SourceDone{},
			Lineage: id,
		})
	}()

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

		m.emit(Value{
			Value:   v,
			Index:   i,
			Lineage: emitId,
		})
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
		Func: func(v Value) {
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

func (m *NodeMachinery) Incoming(pipeId uint64, v Value) {
	if m.n.Inputs() == 1 {
		m.Run(v)
		return
	}

	valuesCh := make(chan []Value)

	go m.stagedValues.set(pipeId, v, int(m.n.Inputs()), valuesCh)

	for values := range valuesCh {
		m.Run(values...)
	}
}

type valueStore struct {
	mem    *IDTrie
	staged [][]Value
	m      sync.Mutex
}

func (s *valueStore) set(pipeId uint64, incoming Value, cap int, ch chan []Value) {
	defer close(ch)

	s.m.Lock()
	defer s.m.Unlock()

	s.mem.Insert(pipeId, incoming)

	dbg := func(s string, args ...any) {
		//fmt.Printf(s+"\n", args...)
	}

	dbg("SET %v", incoming)

	staged := false
	for stagedIdx := 0; stagedIdx < len(s.staged); stagedIdx++ {
		inputs := s.staged[stagedIdx]

		dbg("CHECK %v", inputs)

		valueAtIndex := inputs[incoming.Index]

		if valueAtIndex.IsSet() {
			dbg("SET")

			continue
		}

		allSet := true
		compatible := true
		for _, value := range inputs {
			dbg("CORR %v", value)

			if !value.IsSet() {
				dbg("NOT SET")
				if value.Index != incoming.Index {
					allSet = false
				}
			} else {
				if !incoming.Lineage.Correlates(value.Lineage) {
					compatible = false
					break
				}
				dbg("CORRELATES")
			}
		}

		if compatible {
			dbg("COMPATIBLE")

			staged = true

			inputs[incoming.Index] = incoming

			if allSet {
				dbg("ALLSET")

				s.staged = slices.Delete(s.staged, stagedIdx, stagedIdx+1)
				stagedIdx--

				ch <- inputs
			} else {
				dbg("STAGE")

				s.staged[stagedIdx] = inputs
			}
		}
	}

	if staged {
		dbg("HAS STAGED")

		return
	}

	dbg("STAGING")

	values := make([]Value, cap)
	for i := range values {
		values[i] = Value{
			Index: uint8(i),
		}
	}
	values[incoming.Index] = incoming

	if t, ok := s.mem.GetTrie(pipeId, nil /*LineageID{incoming.Lineage[0]}*/); ok {
		dbg("HAS TRIE")

		t.Walk(func(value Value) {
			dbg("CHECK %v", value)

			if values[value.Index].IsSet() {
				return
			}

			if value.Lineage.Correlates(incoming.Lineage) {
				dbg("ADD %v", value)
				values[value.Index] = value
			}
		})

		allSet := true
		for _, value := range values {
			if !value.IsSet() {
				allSet = false
				break
			}
		}

		if allSet {
			dbg("ALL SET")

			ch <- values

			return
		}
	}

	s.staged = append(s.staged, values)
}

func Cast[T any](v any) T {
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
