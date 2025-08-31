package set

import (
	"fmt"
	"strings"
)

// A HashSet is an unordered collection of unique elements of type E.
// It uses a hashmap (with empty values) as its backing storage, optimized for large sets (around 20
// elements or larger - see benchmark_test.go for benchmarks).
//
// The zero value for a HashSet is ready to use. It must not be copied after first use.
//
// HashSet implements [Set] when passed by pointer, and [ComparableSet] when passed by value.
type HashSet[E comparable] struct {
	elements map[E]struct{}
}

// NewHashSet creates a new [HashSet] for elements of type E.
// It must not be copied after first use.
func NewHashSet[E comparable]() HashSet[E] {
	return HashSet[E]{elements: make(map[E]struct{})}
}

// HashSetWithCapacity creates a new [HashSet], with at least the given initial capacity.
// It must not be copied after first use.
func HashSetWithCapacity[E comparable](capacity int) HashSet[E] {
	return HashSet[E]{elements: make(map[E]struct{}, capacity)}
}

// HashSetOf creates a new [HashSet] from the given elements.
// It must not be copied after first use.
// Duplicate elements are added only once.
func HashSetOf[E comparable](elements ...E) HashSet[E] {
	return HashSetFromSlice(elements)
}

// HashSetFromSlice creates a new [HashSet] from the elements in the given slice.
// It must not be copied after first use.
// Duplicate elements in the slice are added only once.
func HashSetFromSlice[E comparable](elements []E) HashSet[E] {
	set := HashSet[E]{elements: make(map[E]struct{}, len(elements))}

	for _, element := range elements {
		set.elements[element] = struct{}{}
	}

	return set
}

// Add adds the given element to the set.
// If the element is already present in the set, Add is a no-op.
//
// If the hash set was not previously initialized through one of the constructors in this package,
// it will be initialized here.
func (set *HashSet[E]) Add(element E) {
	if set.elements == nil {
		set.elements = make(map[E]struct{})
	}

	set.elements[element] = struct{}{}
}

// AddMultiple adds the given elements to the set. Duplicate elements are added only once, and
// elements already present in the set are not added.
//
// If the hash set was not previously initialized through one of the constructors in this package,
// it will be initialized here.
func (set *HashSet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

// AddFromSlice adds the elements from the given slice to the set. Duplicate elements are added only
// once, and elements already present in the set are not added.
//
// If the hash set was not previously initialized through one of the constructors in this package,
// it will be initialized here.
func (set *HashSet[E]) AddFromSlice(elements []E) {
	if set.elements == nil {
		set.elements = make(map[E]struct{}, len(elements))
	}

	for _, element := range elements {
		set.elements[element] = struct{}{}
	}
}

// AddFromSet adds elements from the given other set to the set.
//
// If the hash set was not previously initialized through one of the constructors in this package,
// it will be initialized here.
func (set *HashSet[E]) AddFromSet(otherSet ComparableSet[E]) {
	if set.elements == nil {
		set.elements = make(map[E]struct{}, otherSet.Size())
	}

	otherSet.All()(
		func(element E) bool {
			set.Add(element)
			return true
		},
	)
}

// Remove removes the given element from the set.
// If the element is not present in the set, Remove is a no-op.
func (set HashSet[E]) Remove(element E) {
	delete(set.elements, element)
}

// Clear removes all elements from the set, leaving an empty set with the same capacity as before.
func (set HashSet[E]) Clear() {
	for element := range set.elements {
		delete(set.elements, element)
	}
}

// Contains checks if given element is present in the set.
func (set HashSet[E]) Contains(element E) bool {
	if set.elements == nil {
		return false
	}

	_, contains := set.elements[element]
	return contains
}

// Size returns the number of elements in the set.
func (set HashSet[E]) Size() int {
	return len(set.elements)
}

// IsEmpty checks if there are 0 elements in the set.
func (set HashSet[E]) IsEmpty() bool {
	return len(set.elements) == 0
}

// Equals checks if the set contains exactly the same elements as the other given set.
func (set HashSet[E]) Equals(otherSet ComparableSet[E]) bool {
	return set.Size() == otherSet.Size() && set.IsSubsetOf(otherSet)
}

// IsSubsetOf checks if all of the elements in the set exist in the other given set.
func (set HashSet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	for element := range set.elements {
		if !otherSet.Contains(element) {
			return false
		}
	}

	return true
}

// IsSupersetOf checks if the set contains all of the elements in the other given set.
func (set HashSet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	return otherSet.IsSubsetOf(set)
}

// Union creates a new set that contains all the elements of the receiver set and the other given
// set. The underlying type of the returned set is a *HashSet - to get a value type, use
// [HashSet.UnionHashSet] instead.
func (set HashSet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionHashSet(otherSet)
	return &union
}

// UnionHashSet creates a new HashSet that contains all the elements of the receiver set and the
// other given set.
func (set HashSet[E]) UnionHashSet(otherSet ComparableSet[E]) HashSet[E] {
	union := HashSetWithCapacity[E](set.Size() + otherSet.Size())

	for element := range set.elements {
		union.Add(element)
	}

	otherSet.All()(
		func(element E) bool {
			union.Add(element)
			return true
		},
	)

	return union
}

// Intersection creates a new set with only the elements that exist in both the receiver set and the
// other given set. The underlying type of the returned set is a *HashSet - to get a value type,
// use [HashSet.IntersectionHashSet] instead.
func (set HashSet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionHashSet(otherSet)
	return &intersection
}

// IntersectionHashSet creates a new HashSet with only the elements that exist in both the receiver
// set and the other given set.
func (set HashSet[E]) IntersectionHashSet(otherSet ComparableSet[E]) HashSet[E] {
	var capacity int
	if set.Size() < otherSet.Size() {
		capacity = set.Size()
	} else {
		capacity = otherSet.Size()
	}

	intersection := HashSetWithCapacity[E](capacity)
	for element := range set.elements {
		if otherSet.Contains(element) {
			intersection.Add(element)
		}
	}

	return intersection
}

// ToSlice creates a slice with all the elements in the set.
//
// Since sets are unordered, the order of elements in the slice is non-deterministic, and may vary
// even when called multiple times on the same set.
func (set HashSet[E]) ToSlice() []E {
	slice := make([]E, len(set.elements))

	i := 0
	for element := range set.elements {
		slice[i] = element
		i++
	}

	return slice
}

// ToMap returns a map with all the set's elements as keys.
//
// Mutating the map will also mutate the set, since it uses the same backing storage. To avoid this,
// call CopyHashSet first.
func (set HashSet[E]) ToMap() map[E]struct{} {
	return set.elements
}

// Copy creates a new set with all the same elements and capacity as the original set.
// The underlying type of the returned set is a *HashSet - to get a value type, use
// [HashSet.CopyHashSet] instead.
func (set HashSet[E]) Copy() Set[E] {
	newSet := set.CopyHashSet()
	return &newSet
}

// CopyHashSet creates a new HashSet with all the same elements and capacity as the original set.
func (set HashSet[E]) CopyHashSet() HashSet[E] {
	newSet := HashSet[E]{elements: make(map[E]struct{}, len(set.elements))}

	for element := range set.elements {
		newSet.elements[element] = struct{}{}
	}

	return newSet
}

// String returns a string representation of the set, implementing [fmt.Stringer].
//
// Since sets are unordered, the order of elements in the string may differ each time it is called.
//
// A HashSet of elements 1, 2 and 3 will be printed as: HashSet{1, 2, 3} (though the order may
// vary).
func (set HashSet[E]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("HashSet{")

	i := 0
	for element := range set.elements {
		_, _ = fmt.Fprint(&stringBuilder, element)

		if i < len(set.elements)-1 {
			stringBuilder.WriteString(", ")
		}

		i++
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

// All returns an [Iterator] function, which when called will loop over the elements in the set and
// call the given yield function on each element. If yield returns false, iteration stops.
//
// Since sets are unordered, iteration order is non-deterministic.
func (set HashSet[E]) All() Iterator[E] {
	return func(yield func(element E) bool) {
		for element := range set.elements {
			if !yield(element) {
				break
			}
		}
	}
}
