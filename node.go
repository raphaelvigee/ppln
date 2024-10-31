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
	new     bool
	lineage LineageID
	m       sync.Mutex
	onDone  func()
}

func SameLineageRef() *LineageRef {
	return &LineageRef{}
}

func NewLineageRef() *LineageRef {
	return &LineageRef{new: true}
}

func (l *LineageRef) get(id LineageID, onDone func()) LineageID {
	l.m.Lock()
	defer l.m.Unlock()

	l.onDone = onDone

	if l.new {
		l.lineage = append(id, sourceId.Add(1))
	} else {
		if len(l.lineage) > 0 {
			return l.lineage
		}

		l.lineage = id
	}

	return l.lineage
}

func (l *LineageRef) Done() {
	f := l.onDone
	if f == nil {
		return
	}

	f()
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

func ArrayEach[E any, S ~[]E]() StreamNode1x2[S, int, E] {
	return NewFuncStreamNode1x2[S, int, E](func(v S, emitIndex func(*LineageRef, int), emitValue func(*LineageRef, E)) {
		for i, e := range v {
			l := NewLineageRef()
			emitIndex(l, i)
			emitValue(l, e)
			l.Done()
		}
	})
}

func MapEach[K comparable, V any, S ~map[K]V]() StreamNode1x2[S, K, V] {
	return NewFuncStreamNode1x2[S, K, V](func(v S, emitKey func(*LineageRef, K), emitValue func(*LineageRef, V)) {
		for i, e := range v {
			l := NewLineageRef()
			emitKey(l, i)
			emitValue(l, e)
			l.Done()
		}
	})
}

type Machinery interface {
	Machinery() *NodeMachinery
}

var pipeId atomic.Uint64

//goland:noinspection ALL
func Pipeline(to Machinery, from ...interface {
	Machinery
	NodeHas1Out
}) {
	ctx := context.Background() // TODO

	pipeId := pipeId.Add(1)

	var startedWg sync.WaitGroup

	for i, from := range from {
		startedWg.Add(1)

		go func() {
			debugger.SetLabels(func() []string {
				return []string{"where", fmt.Sprintf("Pipeline %p (%v) -> %p", from, i, to)}
			})

			ch := from.Machinery().Listen(ctx)

			startedWg.Done()

			for v := range ch {
				if _, ok := v.Value.(SourceDone); ok {
					v.Index = uint8(i)
					to.Machinery().Incoming(pipeId, v)
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

	EnableDebug bool
}

var sourceId atomic.Uint64

func (m *NodeMachinery) NewSourceRun(values ...any) {
	if len(values) == 0 {
		m.Run()
		return
	}

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
	case 0:
		id = LineageID{sourceId.Add(1)}
	case 1:
		id = inputs[0].Lineage
	default:
		idm := map[uint64]struct{}{}
		for _, value := range inputs {
			for _, id := range value.Lineage {
				idm[id] = struct{}{}
			}
		}

		id = slices.Collect(maps.Keys(idm))
		slices.Sort(id)
	}

	if len(id) == 0 {
		panic("empty lineage")
	}

	args := make([]any, len(inputs))
	for i, v := range inputs {
		args[i] = v.Value
	}

	m.n.Do(args, func(i uint8, l *LineageRef, v any) {
		if i >= m.n.Outputs() {
			panic(fmt.Sprintf("trying to output %v, but node has %v outputs", i, m.n.Outputs()))
		}

		emitId := id
		if l != nil {
			emitId = l.get(id, func() {
				m.SourceDone(l)
			})
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
	if _, ok := v.Value.(SourceDone); ok {
		m.stagedValues.cleanup(pipeId, v)
		return
	}

	if m.n.Inputs() == 1 {
		m.Run(v)
		return
	}

	valuesCh := make(chan []Value)

	m.stagedValues.debug = m.EnableDebug

	go m.stagedValues.set(pipeId, v, int(m.n.Inputs()), valuesCh)

	for values := range valuesCh {
		m.Run(values...)
	}
}

func (m *NodeMachinery) SourceDone(l *LineageRef) {
	return
	m.emit(Value{
		Value:   SourceDone{},
		Lineage: l.lineage,
	})
}

func (m *NodeMachinery) Debug() {
	fmt.Println("Staged:")
	for _, values := range m.stagedValues.staged {
		fmt.Println("  ", values)
	}
	fmt.Println("Mem:")
	m.stagedValues.mem.Walk(func(value Value) {
		fmt.Println("  ", value)
	})
}

type valueStore struct {
	mem    *IDTrie
	staged [][]Value
	m      sync.Mutex
	debug  bool
}

func (s *valueStore) cleanup(pipeId uint64, v Value) {
	s.m.Lock()
	defer s.m.Unlock()

	s.mem.Remove(pipeId, v.Lineage, int(v.Index))

	for stagedIdx := 0; stagedIdx < len(s.staged); stagedIdx++ {
		inputs := s.staged[stagedIdx]

		if inputs[v.Index].Lineage.Correlates(v.Lineage) {
			s.staged = slices.Delete(s.staged, stagedIdx, stagedIdx+1)
			stagedIdx--
		}
	}
}

func (s *valueStore) set(pipeId uint64, incoming Value, cap int, ch chan []Value) {
	defer close(ch)

	s.m.Lock()
	defer s.m.Unlock()

	s.mem.Insert(pipeId, incoming)

	dbg := func(msg string, args ...any) {
		if s.debug {
			fmt.Printf(msg+"\n", args...)
		}
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
					dbg("NOT CORRELATES")
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
