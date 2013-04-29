package wisur

type Node struct {
	Parent *Node
	Fils   [8]*Node
}

func New() Node {
	var res Node

	return res
}
