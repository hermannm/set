package set

// A Set is an unordered collection of unique elements of type E.
//
// Three types in this package implement Set:
//   - [ArraySet] uses an array as its backing storage, optimized for small sets
//   - [HashSet] uses a hashmap (with empty values) as its backing storage, optimized for large sets
//   - [DynamicSet] starts out as an ArraySet, but transforms itself to a HashSet once it reaches a
//     size threshold
type Set[E comparable] interface {
	ComparableSet[E]

	// Add adds the given element to the set.
	// If the element is already present in the set, Add is a no-op.
	Add(element E)

	// AddMultiple adds the given elements to the set. Duplicate elements are added only once, and
	// elements already present in the set are not added.
	AddMultiple(elements ...E)

	// AddFromSlice adds the elements from the given slice to the set. Duplicate elements are added
	// only once, and elements already present in the set are not added.
	AddFromSlice(elements []E)

	// AddFromSet adds elements from the given other set to the set.
	AddFromSet(otherSet ComparableSet[E])

	// Remove removes the given element from the set.
	// If the element is not present in the set, Remove is a no-op.
	Remove(element E)

	// Clear removes all elements from the set. When possible, it will retain the same capacity as
	// before.
	Clear()
}

// A ComparableSet is the value type for a Set, with only the methods that will not mutate the set.
type ComparableSet[E comparable] interface {
	// Contains checks if given element is present in the set.
	Contains(element E) bool

	// Size returns the number of elements in the set.
	Size() int

	// IsEmpty checks if there are 0 elements in the set.
	IsEmpty() bool

	// Equals checks if the set contains exactly the same elements as the other given set.
	Equals(otherSet ComparableSet[E]) bool

	// IsSubsetOf checks if all of the elements in the set exist in the other given set.
	IsSubsetOf(otherSet ComparableSet[E]) bool

	// IsSupersetOf checks if the set contains all of the elements in the other given set.
	IsSupersetOf(otherSet ComparableSet[E]) bool

	// Union creates a new set that contains all the elements of the receiver set and the other
	// given set. The underlying type of the returned set will be the same as the receiver.
	Union(otherSet ComparableSet[E]) Set[E]

	// Intersection creates a new set with only the elements that exist in both the receiver set and
	// the other given set. The underlying type of the returned set will be the same as the
	// receiver.
	Intersection(otherSet ComparableSet[E]) Set[E]

	// ToSlice creates a slice with all the elements in the set.
	//
	// Since sets are unordered, the order of elements in the slice is non-deterministic, and may
	// vary even when called multiple times on the same set.
	ToSlice() []E

	// ToMap creates a map with all the set's elements as keys.
	ToMap() map[E]struct{}

	// Copy creates a new set with all the same elements as the original set, and the same
	// underlying type.
	Copy() Set[E]

	// String returns a string representation of the set, implementing [fmt.Stringer].
	//
	// Since sets are unordered, the order of elements in the string may differ each time it is
	// called.
	String() string

	// Iterate loops over every element in the set, and calls the given function on it.
	// It stops iteration if the function returns false.
	//
	// Since sets are unordered, iteration order is non-deterministic.
	//
	// The boolean return from Iterate is there to satisfy the future interface for
	// [range-over-func] in Go, and is always false.
	//
	// [range-over-func]: https://github.com/golang/go/issues/61405
	Iterate(loopBody func(element E) bool) bool
}
