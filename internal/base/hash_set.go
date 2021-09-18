package base

type Set interface {
	// Exist determine the data whether have already existed in the Set or not, if they have, return true, otherwise return false.
	Exist(interface{}) bool

	// Put data to the Set.
	Put(interface{})

	// Add data to the Set, if the data have already exist in the Set, this method return false, otherwise return true.
	Add(interface{}) bool

	// Remove data from the Set, if the data do not exist in the Set, this method will return false, otherwise will return true.
	Remove(interface{}) bool

	// Len return lentgh of the set.
	Len() int

	// IsEmpty return true if the set empty.
	IsEmpty() bool

	// ToSlice turn all elements in the Set into a slice
	ToSlice() []interface{}

	// Merge anotherSet to the Set.
	Merge(anotherSet Set)
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

// Exist determine the data whether have already existed in the Set or not, if they have, return true, otherwise return false.
func (set *mapHashSet) Exist(data interface{}) bool {
	_, exist := set.innerMap[data]
	return exist
}

// Put data to the Set.
func (set *mapHashSet) Put(data interface{}) {
	set.innerMap[data] = true
}

// Add data to the Set, if the data have already exist in the Set, this method return false, otherwise return true.
func (set *mapHashSet) Add(data interface{}) bool {
	if set.Exist(data) {
		return false
	}
	set.Put(data)
	return true
}

// Remove data from the Set, if the data do not exist in the Set, this method will return false, otherwise will return true.
func (set *mapHashSet) Remove(data interface{}) bool {
	if !set.Exist(data) {
		return false
	}
	delete(set.innerMap, data)
	return true
}

// Len return lentgh of the set.
func (set *mapHashSet) Len() int {
	return len(set.innerMap)
}

// IsEmpty return true if the set empty.
func (set *mapHashSet) IsEmpty() bool {
	return len(set.innerMap) == 0
}

// ToSlice turn all elements in the Set into a slice
func (set *mapHashSet) ToSlice() []interface{} {
	slice := make([]interface{}, 0, set.Len())
	for k := range set.innerMap {
		slice = append(slice, k)
	}
	return slice
}

// Merge anotherSet to the Set.
func (set *mapHashSet) Merge(anotherSet Set) {
	for _, v := range anotherSet.ToSlice() {
		set.Put(v)
	}
}
