package set

// A Set is an unordered collection of unique items of type T.
//
// Three types in this package implement Set:
//   - [ArraySet] uses an array as its backing storage, and is optimized for small sets
//   - [HashSet] uses a hashmap (with empty values) as its backing storage, optimized for large sets
//   - [DynamicSet] starts out as an ArraySet, but transforms itself to a HashSet once it reaches a
//     size threshold
type Set[T comparable] interface {
	ComparableSet[T]
	Add(item T)
	AddMultiple(items ...T)
	AddFromSlice(items []T)
	MergeWith(otherSet ComparableSet[T])
	Remove(item T)
	Clear()
}

// A ComparableSet is the value type for a Set, with only the methods that will not mutate the set.
type ComparableSet[T comparable] interface {
	Contains(item T) bool
	Size() int
	IsEmpty() bool
	Equals(otherSet ComparableSet[T]) bool
	IsSubsetOf(otherSet ComparableSet[T]) bool
	IsSupersetOf(otherSet ComparableSet[T]) bool
	Union(otherSet ComparableSet[T]) Set[T]
	Intersection(otherSet ComparableSet[T]) Set[T]
	ToSlice() []T
	ToMap() map[T]struct{}
	ToArraySet() ArraySet[T]
	ToHashSet() HashSet[T]
	ToDynamicSet() DynamicSet[T]
	Copy() Set[T]
	String() string
	Iterate(loopBody func(item T) bool) bool
}
