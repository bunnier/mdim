package types

type Set interface {
	// Exist If the data has already exist in the Set it return true, otherwise return false.
	Exist(interface{}) bool

	// Put Put data to the Set.
	Put(interface{})

	// Add If the data has already exist in the Set it return false, otherwise add the data to the Set then return true.
	Add(interface{}) bool

	// Remove If the data do not exist in the Set it will return false, otherwise remove the data then return true.
	Remove(interface{}) bool

	// Len Return lentgh of the set.
	Len() int

	// IsEmpty Return true if the set empty.
	IsEmpty() bool

	// ToSlice Turn all elements in the Set into a slice
	ToSlice() []interface{}

	// Merge Merge anotherSet to the Set.
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

// Exist If the data has already exist in the Set it return true, otherwise return false.
func (set *mapHashSet) Exist(data interface{}) bool {
	_, exist := set.innerMap[data]
	return exist
}

// Put Put data to the Set.
func (set *mapHashSet) Put(data interface{}) {
	set.innerMap[data] = true
}

// Add If the data has already exist in the Set it return false, otherwise add the data to the Set then return true.
func (set *mapHashSet) Add(data interface{}) bool {
	if set.Exist(data) {
		return false
	}
	set.Put(data)
	return true
}

// Remove If the data do not exist in the Set it will return false, otherwise remove the data then return true.
func (set *mapHashSet) Remove(data interface{}) bool {
	if !set.Exist(data) {
		return false
	}
	delete(set.innerMap, data)
	return true
}

// Len Return lentgh of the set.
func (set *mapHashSet) Len() int {
	return len(set.innerMap)
}

// IsEmpty Return true if the set empty.
func (set *mapHashSet) IsEmpty() bool {
	return len(set.innerMap) == 0
}

// ToSlice Turn all elements in the Set into a slice
func (set *mapHashSet) ToSlice() []interface{} {
	slice := make([]interface{}, 0, set.Len())
	for k := range set.innerMap {
		slice = append(slice, k)
	}
	return slice
}

// Merge Merge anotherSet to the Set.
func (set *mapHashSet) Merge(anotherSet Set) {
	for _, v := range anotherSet.ToSlice() {
		set.Put(v)
	}
}
