package ppln

type Node interface {
	Do()

	OutputsValues() []any

	Inputs() []Node
	Outputs() []Node
}
