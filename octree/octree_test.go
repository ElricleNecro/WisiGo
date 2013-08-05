package octree

import (
	rg "github.com/ElricleNecro/WisiGo/ReadGadget"
	"math/rand"
	"os"
	"sort"
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

func explore_node(a *Node) {
	if a.Fils != nil {
		explore_node(a.Fils)
	} else {
		println(a.level)
	}
	if a.Frere != nil {
		explore_node(a.Frere)
	}
}

func TestCreate(t *testing.T) {
	println("\033[31mIn TestCreate\033[00m")
	UseGoRoutine = false
	rdm := rand.New(rand.NewSource(int64(33)))
	nb_part := 10
	part := make([]rg.Particule, nb_part)

	for i := 0; i < nb_part; i++ {
		for j, _ := range part[i].Pos {
			part[i].Pos[j] = 2.0*rdm.Float32() - 1.0
			part[i].Id = int64(i + 1)
		}
	}

	center := [...]float32{0., 0., 0.}
	root := New(part, center, 2.0)
	t.Log("Testing Creation")
	root.Create(1)

	//fmt.Println(root.Part)

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

	print("\033[31m")
	explore_node(root)
	print("\033[00m")
}

func TestNeighbor(t *testing.T) {

	println("\033[31mIn TestNeighbor\033[00m")

	// rand object
	rdm := rand.New(rand.NewSource(int64(33)))

	// number of particle to generate
	nb_part := 5000

	// number of particles to search
	nb_near := 10

	// allocate the memory of the slice
	part := make([]rg.Particule, nb_part)

	// the same for particles found
	searchy := make([]Search, nb_near)
	search_brut := make([]Search, nb_near)

	// loop over particles and give them positions
	for i := 0; i < nb_part; i++ {
		for j, _ := range part[i].Pos {
			part[i].Pos[j] = 2.0*rdm.Float32() - 1.0
			part[i].Id = int64(i + 1)
		}
	}

	// create the tree
	center := [...]float32{0.0, 0.0, 0.}
	tree := New(part, center, 2.0)

	file, err := os.Create("test/neigh.dat")
	defer file.Close()

	if err != nil {
		t.Error(err)
	}
	if err = tree.savePart(file); err != nil {
		t.Error(err)
	}
	file.Close()

	file, err = os.Create("test/tree_neigh.dat")
	if err != nil {
		t.Error(err)
	}
	if err = tree.SaveNode(file); err != nil {
		t.Error(err)
	}
	// loop over particles and get their 10 closest neighbors
	for _, p := range part {

		// search neighbors with the tree
		searchy = tree.Nearest(p, nb_near)

		// sort elements by distance
		sort.Sort(ByDist(searchy))

		// the same with brute force
		search_brut = bruteForceNeighbors(part, p, nb_part, nb_near)

		// loop over selected points
		for i := 0; i < nb_near; i++ {

			// Check comparison between brute force and the tree
			if searchy[i].Part.Id != search_brut[i].Part.Id {
				t.Fatal("Search nearest neighbors failed because particles are different with brute force !")
			}

			// check that there is not the same particle as we search
			if searchy[i].Part.Id == p.Id {
				t.Fatal("The particle used for search is included in the neighbors !")
			}
			//println(p.Id, searchy[i].Part.Id, search_brut[i].Part.Id)
			//println(searchy[i].Radius, search_brut[i].Radius)
		}

	}

}

func bruteForceNeighbors(parts []rg.Particule, part0 rg.Particule, nb_part int, n_search int) (searchy []Search) {

	// allocate memory for all search
	all_search := make([]Search, nb_part-1)

	// compute the distance for all particules
	i := 0
	for _, p := range parts {

		if part0.Id != p.Id {
			// Compute the distance between the point and all points
			all_search[i].Radius = float64(part0.Dist(p))
			all_search[i].Part = p
			i += 1
		}

	}

	// need to order the distances and corresponding points
	// by the distance
	sort.Sort(ByDist(all_search))

	// Get the neighbors
	searchy = all_search[:n_search]

	return

}

func TestSetPart(t *testing.T) {
	println("\033[31mIn TestSetPart\033[00m")
	rdm := rand.New(rand.NewSource(int64(33)))
	nb_part := 10
	part := make([]rg.Particule, nb_part)

	for i := 0; i < nb_part; i++ {
		for j, _ := range part[i].Pos {
			part[i].Pos[j] = 2.0*rdm.Float32() - 1.0
			part[i].Id = int64(i + 1)
		}
	}

	center := [...]float32{0.5, 0.5, 0.}
	root := New(part, center, 0.5)
	t.Log("Testing Creation")
	root.Create(1)
	root.setPart()
}

func TestPotential(t *testing.T) {
	rdm := rand.New(rand.NewSource(int64(33)))
	nb_part := 10
	part := make([]rg.Particule, nb_part)
	var testy rg.Particule

	for i := 0; i < nb_part; i++ {
		for j, _ := range part[i].Pos {
			part[i].Pos[j] = 2.0*rdm.Float32() - 1.0
			part[i].Id = int64(i + 1)
			part[i].Mass = 1.0
		}
	}

	center := [...]float32{0., 0., 0.}
	root := New(part, center, 2.0)
	Ggrav = 1.0

	testy.Mass = 1.0
	testy.Id = 0
	for i, _ := range testy.Pos {
		testy.Pos[i] = 0.
	}

	pot := root.Potential(testy, 0.0)
	pot2 := bruteforcePotential(part, testy)

	if pot != pot2 {
		t.Fatal("Tree Code is different from brute force:", pot, " != ", pot2)
	}
}

func bruteforcePotential(part []rg.Particule, testy rg.Particule) float64 {
	var pot float64 = 0.0
	for _, v := range part {
		pot += Ggrav * (float64)(v.Mass/v.Dist(testy))
	}
	return pot
}
