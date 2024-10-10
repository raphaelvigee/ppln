package ppln

import "sync"

func NewDAG(defs ...[]Edge) *DAG {
	d := &DAG{
		edges:  map[Node][]Node{},
		states: map[Node]*NodeState{},
	}

	for _, def := range defs {
		d.Add(def)
	}

	return d
}

type DAG struct {
	edges  map[Node][]Node
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
		d.edges[edge.From] = append(d.edges[edge.From], edge.To)

		d.state(edge.From)
		d.state(edge.To)
	}
}

func (d *DAG) emit(from, to Node, idx int, v any) {
	d.state(to).emit(idx, v)
}

func (d *DAG) state(n Node) *NodeState {
	d.m.Lock()
	defer d.m.Unlock()

	if s, ok := d.states[n]; ok {
		return s
	}

	s := &NodeState{
		Node:   n,
		Inputs: make(chan Value),
	}
	d.states[n] = s

	go d.waitForRun(s, n)

	return s
}

func (d *DAG) waitForRun(s *NodeState, n Node) {
	//for v := range s.Inputs {
	//
	//}
}

func (d *DAG) run(s *NodeState, n Node, inputs []any) {
	switch n := n.(type) {
	case BatchNode:
		d.runBatch(s, n, inputs)
	case StreamNode:
		d.runStream(s, n, inputs)
	default:
		panic("not supposed to happen")
	}
}

func (d *DAG) runBatch(s *NodeState, n BatchNode, inputs []any) {
	outputs := n.Do(inputs)
	for i, v := range outputs {
		d.emitToDeps(n, i, v)
	}
}

func (d *DAG) runStream(s *NodeState, n StreamNode, inputs []any) {
	n.Do(inputs, func(i int, v any) {
		d.emitToDeps(n, i, v)
	})
}

func (d *DAG) emitToDeps(from Node, i int, v any) {
	for _, to := range d.edges[from] {
		d.emit(from, to, i, v)
	}
}

//func (d *DAG) Run(n Node) *State {
//	s := &State{d: d}
//	go s.run(n)
//
//	return s
//}
//
//type State struct {
//	m           sync.Mutex
//	d           *DAG
//	ns          *NodeState
//	nodeToState map[Node]*NodeState
//}
