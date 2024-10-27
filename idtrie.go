package ppln

const numBuckets = 32

type IDTrieNode struct {
	Parent   *IDTrieNode
	Values   []Value
	Children [numBuckets]*IDTrieNode
}

func (n *IDTrieNode) Walk(f func(Value)) {
	for _, child := range n.Children {
		if child == nil {
			continue
		}

		for _, value := range child.Values {
			if value.IsSet() {
				f(value)
			}
		}

		child.Walk(f)
	}
}

func NewIDTrieNode(parent *IDTrieNode, numValues int) *IDTrieNode {
	node := &IDTrieNode{
		Parent: parent,
		Values: make([]Value, numValues),
	}

	return node
}

type IDTrie struct {
	NumValues int
	RootNode  *IDTrieNode
}

func NewIDTrie(numValues int) *IDTrie {
	return NewIDTrieAt(NewIDTrieNode(nil, numValues), numValues)
}

func NewIDTrieAt(node *IDTrieNode, numValues int) *IDTrie {
	return &IDTrie{
		NumValues: numValues,
		RootNode:  node,
	}
}

func (t *IDTrie) findNode(pipeId uint64, id LineageID, create bool) *IDTrieNode {
	ogid := id

	id = make(LineageID, len(ogid)+1)
	id[0] = pipeId
	copy(id[1:], ogid)

	current := t.RootNode
	for len(id) > 0 {
		idv := id[0]

		var bucketIdx int
		if idv < numBuckets {
			bucketIdx = int(idv)
			id = id[1:]
		} else {
			bucketIdx = int(idv % numBuckets)
			id[0] = idv / numBuckets
		}

		if current.Children[bucketIdx] == nil {
			if create {
				current.Children[bucketIdx] = NewIDTrieNode(current, t.NumValues)
			} else {
				return nil
			}
		}
		current = current.Children[bucketIdx]
	}

	return current
}

func (t *IDTrie) Insert(pipeId uint64, v Value) {
	node := t.findNode(pipeId, v.Lineage, true)

	node.Values[v.Index] = v
}

//func (t *IDTrie) Remove(pipeId uint64, id LineageID, idx int) {
//	node := t.findNode(pipeId, id, true)
//
//	if node == nil {
//		return
//	}
//
//	node.Values[idx] = Value{
//		Index: uint8(idx),
//	}
//
//	allEmpty := true
//	for _, value := range node.Values {
//		if value.IsSet() {
//			allEmpty = false
//		}
//	}
//
//	if !allEmpty {
//		return
//	}
//
//	for node.Parent != nil {
//
//	}
//}

func (t *IDTrie) Get(pipeId uint64, id LineageID, idx int) (Value, bool) {
	node := t.findNode(pipeId, id, false)

	if node == nil {
		return Value{}, false
	}

	v := node.Values[idx]

	return v, v.IsSet()
}

func (t *IDTrie) GetTrie(pipeId uint64, id LineageID) (*IDTrie, bool) {
	node := t.findNode(pipeId, id, false)

	if node == nil {
		return nil, false
	}

	return NewIDTrieAt(node, t.NumValues), true
}

func (t *IDTrie) Walk(f func(Value)) {
	t.RootNode.Walk(f)
}
