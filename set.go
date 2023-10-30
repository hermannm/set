package set

type Set[T comparable] interface {
	ComparableSet[T]
	Add(item T)
	AddMultiple(items ...T)
	AddFromSlice(items []T)
	MergeWith(otherSet ComparableSet[T])
	Remove(item T)
	Clear()
}

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
