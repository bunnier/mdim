package types

type Set interface {
	// If the data exist in the Set return true, otherwise false.
	Exist(interface{}) bool

	// If the data exist in the Set return false, otherwise add to the Set then return true.
	Add(interface{}) bool

	// If the data do not exist in the Set return false, otherwise remove it then return true.
	Remove(interface{}) bool
}

func NewSet(cap int) Set {
	var set Set = &mapHashSet{
		innerMap: make(map[interface{}]bool, cap),
	}
	return set
}

// A implement of HashSet by map[interface{}]bool
type mapHashSet struct {
	innerMap map[interface{}]bool
}

// If the data exist in the Set return true, otherwise false.
func (set *mapHashSet) Exist(data interface{}) bool {
	_, exist := set.innerMap[data]
	return exist
}

// If the data exist in the Set return false, otherwise add to the Set then return true.
func (set *mapHashSet) Add(data interface{}) bool {
	if set.Exist(data) {
		return false
	}
	set.innerMap[data] = true
	return true
}

// If the data do not exist in the Set return false, otherwise remove it then return true.
func (set *mapHashSet) Remove(data interface{}) bool {
	if !set.Exist(data) {
		return false
	}
	delete(set.innerMap, data)
	return true
}
