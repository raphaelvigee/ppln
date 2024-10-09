package ppln

type Node interface {
	is_node()
}

type BatchNode interface {
	Node

	Do(values []any) []any
}

type StreamNode interface {
	Node

	Do(values []any, emit func(i int, v any))
}

type Edge struct {
	From Node
	To   Node
}

type Value struct {
	Value   any
	Index   int
	Lineage struct{}
}

type NodeState struct {
	Node   Node
	Inputs []chan Value
}

func DAG(defs ...[]Edge) {
	
}
