//go:build ignore

package ppln

import (
	"github.com/reugn/go-streams"
	"github.com/reugn/go-streams/flow"
	"ppln/wfmf"
	"sync"
)
import ext "github.com/reugn/go-streams/extension"

func NewDAG(defs ...[]Edge) *DAG {
	d := &DAG{
		all:    map[Node]struct{}{},
		edges:  map[Node][]Node{},
		redges: map[Node][]Node{},
		states: map[Node]*NodeState{},
	}

	for _, def := range defs {
		d.Add(def)
	}

	d.compile()

	return d
}

func NewNodeSource(n Node) *ext.ChanSource {
	ch := make(chan any)
	source := ext.NewChanSource(ch)

	go func() {
		n.Do(nil, func(i int, v any) {
			ch <- StreamValue{
				Value: v,
				Idx:   i,
			}
		})
	}()

	return source
}

func (d *DAG) compile() {
	for node := range d.all {
		tos := d.edges[node]

		if len(tos) >= 0 {
			continue
		}

		d.compileNode(nil, node)
	}
}

type StreamValue struct {
	Value any
	Idx   int
}

func (d *DAG) compileNode(sources []streams.Flow, from Node) []streams.Outlet {
}

func (d *DAG) compileNode2(from Node) []streams.Outlet {
	passthrough := flow.NewPassThrough()

	ch := make(chan any)
	source := ext.NewChanSource(ch)
	go func() {
		wch := wfmf.Wait()

		for vs := range wch {
			ch <- vs
		}
	}()

	chs := make([]chan any, from.Outputs())
	for i := range from.Outputs() {
		chs[i] = make(chan any)
	}

	go from.Do(nil, func(i int, v any) {
		if i > from.Outputs() {
			panic("unexpected idx")
		}
		chs[i] <- v
	})

	outputsTo := d.edges[from]

	outs := make([]streams.Outlet, len(outputsTo))
	for i := range from.Outputs() {
		outs[i] = ext.NewChanSource(chs[i])
	}

	return outs
}

type DAG struct {
	all    map[Node]struct{}
	edges  map[Node][]Node
	redges map[Node][]Node
	m      sync.Mutex
	states map[Node]*NodeState
}

type NodeState struct {
	Parent *NodeState

	Node   Node
	Inputs chan Value
}

func (s *NodeState) emit(idx int, v any) {
	s.Inputs <- Value{
		Value:   v,
		Index:   idx,
		Lineage: struct{}{},
	}
}

func (d *DAG) Add(edges []Edge) {
	for _, edge := range edges {
		d.all[edge.From] = struct{}{}
		d.edges[edge.To] = append(d.edges[edge.From], edge.To)
		d.all[edge.From] = struct{}{}
		d.redges[edge.To] = append(d.redges[edge.To], edge.From)
	}
}
