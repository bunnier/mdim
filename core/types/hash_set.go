package types

type Set interface {
	// If the data exist in the Set return true, otherwise false.
	Exist(interface{}) bool

	// If the data exist in the Set return false, otherwise add to the Set then return true.
	Add(interface{}) bool

	// Put data to the Set.
	Put(interface{})

	// If the data do not exist in the Set return false, otherwise remove it then return true.
	Remove(interface{}) bool

	// Return len of the set.
	Len() int

	// Return if the set empty.
	IsEmpty() bool

	// Merge anotherSet to current.
	Merge(anotherSet Set)

	// Turn all elements into a slice
	ToSlice() []interface{}
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

// Put data to the Set.
func (set *mapHashSet) Put(data interface{}) {
	set.innerMap[data] = true
}

// If the data exist in the Set return false, otherwise add to the Set then return true.
func (set *mapHashSet) Add(data interface{}) bool {
	if set.Exist(data) {
		return false
	}
	set.Put(data)
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

// Return len of the set.
func (set *mapHashSet) Len() int {
	return len(set.innerMap)
}

// Return if the set empty.
func (set *mapHashSet) IsEmpty() bool {
	return len(set.innerMap) == 0
}

// Turn all elements into a slice
func (set *mapHashSet) ToSlice() []interface{} {
	slice := make([]interface{}, 0, set.Len())
	for k := range set.innerMap {
		slice = append(slice, k)
	}
	return slice
}

// Merge anotherSet to current.
func (set *mapHashSet) Merge(anotherSet Set) {
	for _, v := range anotherSet.ToSlice() {
		set.Put(v)
	}
}
