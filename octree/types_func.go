package octree

func (e ByDist) Len() int {
	return len(e)
}

func (e ByDist) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ByDist) Less(i, j int) bool {
	return e[i].Radius < e[j].Radius
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
