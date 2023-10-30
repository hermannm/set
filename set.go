package set

// A Set is an unordered collection of unique items of type T.
//
// Three types in this package implement Set:
//   - [ArraySet] uses an array as its backing storage, optimized for small sets
//   - [HashSet] uses a hashmap (with empty values) as its backing storage, optimized for large sets
//   - [DynamicSet] starts out as an ArraySet, but transforms itself to a HashSet once it reaches a
//     size threshold
type Set[T comparable] interface {
	ComparableSet[T]

	// Add adds the given item to the set.
	// If the item is already present in the set, Add is a no-op.
	Add(item T)

	// AddMultiple adds the given items to the set.
	// Duplicate items are added only once, and items already present in the set are not added.
	AddMultiple(items ...T)

	// AddFromSlice adds the items from the given slice to the set.
	// Duplicate items are added only once, and items already present in the set are not added.
	AddFromSlice(items []T)

	// MergeWith adds items from the given other set to the set.
	MergeWith(otherSet ComparableSet[T])

	// Remove removes the given item from the set.
	// If the item is not present in the set, Remove is a no-op.
	Remove(item T)

	// Clear removes all items from the set, leaving an empty set with the same capacity as before.
	Clear()
}

// A ComparableSet is the value type for a Set, with only the methods that will not mutate the set.
type ComparableSet[T comparable] interface {
	// Contains checks if given item is present in the set.
	Contains(item T) bool

	// Size returns the number of items in the set.
	Size() int

	// IsEmpty checks if there are 0 items in the set.
	IsEmpty() bool

	// Equals check if the set contains exactly the same items as the other given set.
	Equals(otherSet ComparableSet[T]) bool

	// IsSubsetOf checks if all of the items in the set exist in the other given set.
	IsSubsetOf(otherSet ComparableSet[T]) bool

	// IsSupersetOf checks if hte set contains all of the items in the other given set.
	IsSupersetOf(otherSet ComparableSet[T]) bool

	// Union creates a new set that contains all the items of receiver set and the other given set.
	// The underlying type of the returned set will be the same as the receiver.
	Union(otherSet ComparableSet[T]) Set[T]

	// Intersection creates a new set with only the items that exist in both the receiver set and
	// the other given set.
	// The underlying type of the returned set will be the same as the receiver.
	Intersection(otherSet ComparableSet[T]) Set[T]

	// ToSlice creates a slice with all the items in the set.
	//
	// Since sets are unordered, the order of items in the slice is non-deterministic, and may vary
	// even when called multiple times on the same set.
	ToSlice() []T

	// ToMap creates a map with all the set's items as keys.
	ToMap() map[T]struct{}

	// ToArraySet creates an [ArraySet] from the items in the set.
	// If the set is already an ArraySet, this is equivalent to calling CopyArraySet on it.
	ToArraySet() ArraySet[T]

	// ToHashSet creates a [HashSet] from the items in the set.
	// If the set is already a HashSet, this is equivalent to calling CopyHashSet on it.
	ToHashSet() HashSet[T]

	// ToDynamicSet creates a [DynamicSet] from the items in the set.
	// If the set is already a DynamicSet, this is equivalent to calling CopyDynamicSet on it.
	ToDynamicSet() DynamicSet[T]

	// Copy creates a new set with all the same items as the set it is called on, and the same
	// underlying type.
	Copy() Set[T]

	// String implements [fmt.Stringer] to customize the print format of the Set.
	//
	// Since sets are unordered, the order of items in the string may differ each time it is called.
	String() string

	// Iterate takes a function, which it then calls on every item in the set.
	// It stops iteration if the function returns false.
	//
	// The boolean return from Iterate is there to satisfy the future interface for
	// [range-over-func] in Go, and is always false.
	//
	// [range-over-func]: https://github.com/golang/go/issues/61405
	Iterate(loopBody func(item T) bool) bool
}
