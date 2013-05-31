package ReadGadget

func (e WriteError) Error() string {
	return string(e)
}

func (e ByType) Len() int {
	return len(e)
}

func (e ByType) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ByType) Less(i, j int) bool {
	return e[i].Type < e[j].Type
}

//Swap particle a and b!
func Swap(a, b *Particule) {
	var tmp Particule

	tmp = *a
	*a = *b
	*b = tmp
}

func Equal(a, b Particule) bool {
	for i, _ := range a.Pos {
		if a.Vel[i] != b.Vel[i] || a.Pos[i] != b.Pos[i] {
			return false
		}
	}

	if a.Mass != b.Mass || a.Type != b.Type || a.Id != b.Id {
		return false
	}

	return true
}
