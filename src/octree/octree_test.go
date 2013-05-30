package octree

import (
	rg "ReadGadget"
	"math/rand"
	"os"
	"testing"
)

func TestIn(t *testing.T) {
	var part rg.Particule
	part.Pos[0] = 0.5
	part.Pos[1] = 0.9
	part.Pos[2] = 0.0

	center := [...]float32{0., 0., 0.}
	root := New(nil, center, 2.0)

	if root.In(part) == false {
		t.Error("Particle : ", part.Pos, " not in cube of center : ", root.Center, " and size : ", root.Size)
	}

	part.Pos[1] = -3.5
	if root.In(part) == true {
		t.Error("Particle : ", part.Pos, " in cube of center : ", root.Center, " and size : ", root.Size)
	}
}

func TestSwap(t *testing.T) {
	var part1, part2, opart1, opart2 rg.Particule

	part1.Pos[0] = 5.0
	part1.Pos[1] = 5.0
	part1.Pos[2] = 5.0

	part2.Pos[0] = -340.0
	part2.Pos[1] = -340.0
	part2.Pos[2] = -340.0

	opart1 = part1
	opart2 = part2

	if !rg.Equal(part1, opart1) {
		t.Error("Problem with the Equal test function for part1!")
	}
	if !rg.Equal(part2, opart2) {
		t.Error("Problem with the Equal test function for part2!")
	}

	rg.Swap(&part1, &part2)

	if !rg.Equal(part1, opart2) {
		t.Error("Problem with the Swap function for part1!")
	}
	if !rg.Equal(part2, opart1) {
		t.Error("Problem with the Swap function for part2!")
	}
}

func TestCreate(t *testing.T) {
	rdm := rand.New(rand.NewSource(int64(33)))
	nb_part := 10
	part := make([]rg.Particule, nb_part)

	for i := 0; i < nb_part; i++ {
		for j, _ := range part[i].Pos {
			part[i].Pos[j] = 2.0*rdm.Float32() - 1.0
		}
	}

	center := [...]float32{0., 0., 0.}
	root := New(part, center, 2.0)
	t.Log("Testing Creation")
	root.Create(1)

	file, err := os.Create("test/part.dat")
	defer file.Close()
	if err != nil {
		t.Error(err)
	}
	if err = root.savePart(file); err != nil {
		t.Error(err)
	}
	file.Close()
	file, err = os.Create("test/tree.dat")
	if err != nil {
		t.Error(err)
	}
	if err = root.SaveNode(file); err != nil {
		t.Error(err)
	}
}
