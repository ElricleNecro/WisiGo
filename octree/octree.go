//-tags: main
package octree

import rg "github.com/ElricleNecro/WisiGo/ReadGadget"
import "fmt"
import m "math"

var (
	UseGoRoutine = false //Use of Go Routine for the tree creation?
	// If you set this to true, don't forget to set runtime.GOMAXPROCS
	// to the correct value.
)

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

	var t1 **Node
	t1 = &e.Fils
	var NbUse int32 = 0
	center := [3]float32{0., 0., 0.}
	side := e.Size / 2.
	side2 := side / 2.

	for i := 0; i < 8; i++ {
		//We place the center of son :
		if (i & 1) == 0 {
			center[0] = e.Center[0] + side2
		} else {
			center[0] = e.Center[0] - side2
		}
		if (i & 2) == 0 {
			center[1] = e.Center[1] + side2
		} else {
			center[1] = e.Center[1] - side2
		}
		if (i & 4) == 0 {
			center[2] = e.Center[2] + side2
		} else {
			center[2] = e.Center[2] - side2
		}

		//Creation of the Node :
		*t1 = New(e.Part[NbUse:], center, side)
		// We set the level of the Node :
		(*t1).level = e.level + 1

		//We are calculating the particle contain in the Node :
		NbUtil := t1.setPart()

		println((*t1).level, NbUtil, NbUse)
		print((*t1).level, " ")
		for z := 0; z < 3; z++ {
			print("] ", (*t1).Center[z]-(*t1).Size/2., "; ", (*t1).Center[z]+(*t1).Size/2., "] ")
		}
		println("")

		//If there is particle in the node, we are actualizing the Particle slice, launching calculation of the son and preparing the next brother.
		if NbUtil != 0 {
			(*t1).Part = e.Part[NbUse:(NbUse + NbUtil)]

			fmt.Println("Going down!")
			(*t1).Create(NbMin)
			fmt.Println("Going up!")

			NbUse += NbUtil
			t1 = &(*t1).Frere
		}
	}
}

type Search struct {
	Radius float64
	Part   rg.Particule
}

func Insert(e []Search, to Search) {
	var tmp, bak Search
	var i int

	// We search where to place the particule :
	for i, _ = range e {
		if e[i].Radius > to.Radius {
			tmp = e[i]
			e[i] = to
			i++
			break
		}
	}

	// We move all particule from where we place 'to' to the end of the array
	// (we don't keep those who get out from the array) :
	for ; i < len(e); i++ {
		bak = e[i]
		e[i] = tmp
		tmp = bak
	}
}

func (e *Node) Nearest(part rg.Particule, NbVois int) []Search {
	var res []Search = make([]Search, NbVois)

	for i, _ := range res {
		res[i].Radius = 2. * (float64)(e.Size)
	}

	e.SearchNeighbor(part, res[:])

	return res
}

func (e *Node) Dist(part rg.Particule) float64 {
	dx := m.Max(0., m.Abs((float64)(e.Center[0]-part.Pos[0]))-(float64)(e.Size/2.0))
	dy := m.Max(0., m.Abs((float64)(e.Center[1]-part.Pos[1]))-(float64)(e.Size/2.0))
	dz := m.Max(0., m.Abs((float64)(e.Center[2]-part.Pos[2]))-(float64)(e.Size/2.0))
	return m.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (e *Node) fill_neighboorhood(part rg.Particule, searchy []Search) {
	for _, v := range e.Part {
		if (float64)(v.Dist(part)) < searchy[len(searchy)-1].Radius {
			Insert(searchy[:], Search{(float64)(v.Dist(part)), v})
		}
	}
}

func (e *Node) SearchNeighbor(part rg.Particule, searchy []Search) {
	if e.Dist(part) > searchy[len(searchy)-1].Radius {
		return
	}
	if e.Fils != nil {
		e.Fils.SearchNeighbor(part, searchy[:])
	} else {
		e.fill_neighboorhood(part, searchy[:])
	}
}

func (e *Node) setPart() (NbUse int32) {
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
