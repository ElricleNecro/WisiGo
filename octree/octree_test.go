package octree

import (
	rg "WisiGo/ReadGadget"
	"math/rand"
	"testing"
)

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
	root.Create(1)
}
