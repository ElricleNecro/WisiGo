package wisur

import rg "ReadGadget"

type Node struct {
	Parent *Node
	Fils   [8]*Node
	Center [3]float32
	Size   float32
	Part   []rg.Particule
	level  int64
}

func New(part []rg.Particule, Center [3]float32, side float32) *Node {
	var res Node

	res.Parent = nil
	for i, _ := range res.Fils {
		res.Fils[i] = nil
	}

	res.Size = side
	res.Center = Center
	res.level = 0

	return &res
}

func (e *Node) Create(NbMin int32) {
	for i, v range e.Fils {
		if i & 2 {
			tmp = New()
		}
	}
}
