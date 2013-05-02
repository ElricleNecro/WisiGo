//-tags: main
package octree

import rg "WisiGo/ReadGadget"
import "fmt"

var (
	UseGoRoutine = false //Use of Go Routine for the tree creation? If you set this to true, don't forget to set runtime.GOMAXPROCS to the correct value.
)

//These Node structure will contain every needed informations about Tree Node.
type Node struct {
	Parent, Fils, Frere *Node          //These is the Main Node to Now : the parent, the first son and the next brother.
	Center              [3]float32     //Center of the cube.
	Size                float32        //Side size of the cube.
	Part                []rg.Particule //List of the particle contain in the node.
	level               int64          //Level of the tree.
}

// This function is used to create a tree using a list of the particle in the node, the center of the cube and his size. Her essential case of use will be for create the root of the tree.
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

//This function is used to create all sons of a node.
func (e *Node) Create(NbMin int32) {
	//If we have less particle than NbMin, we don't need do go down in the tree.
	if int32(len(e.Part)) <= NbMin {
		return
	}

	var t1 *Node = e.Fils
	var NbUse int32 = 0
	center := [3]float32{0., 0., 0.}
	side := e.Size / 2.

	for i := 0; i < 8; i++ {
		//We place the center of son :
		if (i & 1) == 0 {
			center[0] = e.Center[0] + side
		} else {
			center[0] = e.Center[0] - side
		}
		if (i & 2) == 0 {
			center[1] = e.Center[1] + side
		} else {
			center[1] = e.Center[1] - side
		}
		if (i & 4) == 0 {
			center[2] = e.Center[2] + side
		} else {
			center[2] = e.Center[2] - side
		}

		//Creation of the Node :
		t1 = New(e.Part[NbUse:], center, side)
		// We set the level of the Node :
		t1.level = e.level + 1

		//We are calculating the particle contain in the Node :
		NbUtil := t1.setPart(NbMin)
		fmt.Println(t1.level, NbUtil)

		//If there is particle in the node, we are actualizing the Particle slice, launching calculation of the son and preparing the next brother.
		if NbUtil != 0 {
			t1.Part = e.Part[NbUse:(NbUse + NbUtil)]
			if UseGoRoutine {
				go t1.Create(NbMin)
			} else {
				t1.Create(NbMin)
			}
			NbUse += NbUtil
			t1 = t1.Frere
		}
	}
}

func (e *Node) setPart(NbMin int32) (NbUse int32) {
	NbUse = 0
	for i, _ := range e.Part {
		if e.In(e.Part[i]) {
			rg.Swap(&e.Part[NbUse], &e.Part[i])
			NbUse += 1
		}
	}

	return
}

//Particle a is she in the node or not?
func (e Node) In(a rg.Particule) bool {
	for i, v := range a.Pos {
		if (e.Center[i]-e.Size/2.) >= v || (e.Center[i]+e.Size/2.) < v {
			return false
		}
	}

	return true
}

func (e Node) savePart(out WriteStringer) error {
	var err error
	if e.Fils != nil {
		err = e.Fils.savePart(out)
		if err != nil {
			return err
		}
	} else {
		for _, v := range e.Part {
			_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d %d\n", v.Pos[0], v.Pos[1], v.Pos[2], v.Id, e.level, len(e.Part)))
			if err != nil {
				return err
			}
		}
	}

	if e.Frere != nil {
		err = e.Frere.savePart(out)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e Node) SaveNode(out WriteStringer) (err error) {
	if e.Fils != nil {
		err = e.Fils.SaveNode(out)
		if err != nil {
			return
		}
	} else {
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]+e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]-e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]-e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]+e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]-e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]-e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]-e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]-e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]-e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]-e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]+e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]+e.Size/2, e.Center[1]+e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]-e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString(fmt.Sprintf("%g %g %g %d %d\n", e.Center[0]-e.Size/2, e.Center[1]+e.Size/2, e.Center[2]+e.Size/2, e.level, len(e.Part)))
		_, err = out.WriteString("\n")

		_, err = out.WriteString("\n")

		if err != nil {
			return
		}
	}

	if e.Frere != nil {
		err = e.Frere.SaveNode(out)
		if err != nil {
			return
		}
	}

	err = nil
	return
}
