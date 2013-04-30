//-tags: main
package main

import rg "ReadGadget"

type Node struct {
	Parent, Fils, Frere *Node
	Center              [3]float32
	Size                float32
	Part                []rg.Particule
	level               int64
}

func New(part []rg.Particule, Center [3]float32, side float32) *Node {
	var res Node

	res.Parent = nil
	res.Frere = nil
	res.Fils = nil

	res.Size = side
	res.Center = Center
	res.level = 0

	res.Part = part

	return &res
}

func (e *Node) Create(NbMin int32) {
	if int32(len(e.Part)) <= NbMin {
		return
	}

	var t1 *Node = e.Fils
	var NbUse int32 = 0
	center := [3]float32{0., 0., 0.}

	for i := 0; i < 8; i++ {
		t1 = New(e.Part[NbUse:], center, side)
		NbUse += t1.setPart(NbMin)
	}
}

func (e *Node) setPart(NbMin int32) (NbUse int32) {
	NbUse = 0
	for i, _ := range e.Part {
		if e.In(e.Part[i]) {
			e.Part.Swap(e.Part[NbUse], e.Part[i])
			NbUse += 1
		}
	}
}
