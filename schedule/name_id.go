package schedule

type NameId struct {
	id    int
	idMap map[string]int
}

func NewNameId() *NameId {
	n := &NameId{}
	n.idMap = make(map[string]int)

	return n
}

func (n *NameId) GetId(name string) int {
	id, has := n.idMap[name]
	if has {
		return id
	}

	n.id++
	n.idMap[name] = n.id

	return n.id
}
