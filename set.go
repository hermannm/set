package set

type Set[T comparable] interface {
	ComparableSet[T]
	Add(T)
	AddMultiple(...T)
	AddFromSlice([]T)
	MergeWith(ComparableSet[T])
	Remove(T)
	Clear()
}

type ComparableSet[T comparable] interface {
	Contains(T) bool
	Size() int
	IsEmpty() bool
	Equals(ComparableSet[T]) bool
	IsSubsetOf(ComparableSet[T]) bool
	IsSupersetOf(ComparableSet[T]) bool
	Union(ComparableSet[T]) Set[T]
	Intersection(ComparableSet[T]) Set[T]
	ToSlice() []T
	ToMap() map[T]struct{}
	ToArraySet() ArraySet[T]
	ToHashSet() HashSet[T]
	ToDynamicSet() DynamicSet[T]
	Copy() Set[T]
	String() string
	Iterate(loopBody func(item T) bool) bool
}
