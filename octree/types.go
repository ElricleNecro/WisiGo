package octree

import rg "github.com/ElricleNecro/WisiGo/ReadGadget"

//These Node structure will contain every needed informations about Tree Node.
type Node struct {
	Parent, Fils, Frere *Node          //These is the Main Node to Now : the parent, the first son and the next brother.
	Center              [3]float32     //Center of the cube.
	Size                float32        //Side size of the cube.
	Part                []rg.Particule //List of the particle contain in the node.
	level               int64          //Level of the tree.
}

type WriteStringer interface {
	WriteString(s string) (ret int, err error)
}

type Search struct {
	Radius float64
	Part   rg.Particule
}

type ByDist []Search
