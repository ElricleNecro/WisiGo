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
