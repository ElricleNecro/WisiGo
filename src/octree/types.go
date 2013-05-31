package octree

type WriteStringer interface {
	WriteString(s string) (ret int, err error)
}
