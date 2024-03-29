package set

import (
	"fmt"
	"strings"
)

// An ArraySet is a collection of unique elements of type E.
// It uses an array as its backing storage, optimized for small sets (up to around 20 elements - see
// benchmark_test.go for benchmarks).
//
// The zero value for an ArraySet is ready to use. It must not be copied after first use.
//
// ArraySet implements [Set] when passed by pointer, and [ComparableSet] when passed by value.
type ArraySet[E comparable] struct {
	elements []E
}

// NewArraySet creates a new [ArraySet] for elements of type E.
// It must not be copied after first use.
func NewArraySet[E comparable]() ArraySet[E] {
	return ArraySet[E]{elements: nil}
}

// ArraySetWithCapacity creates a new [ArraySet], with at least the given initial capacity.
// It must not be copied after first use.
func ArraySetWithCapacity[E comparable](capacity int) ArraySet[E] {
	return ArraySet[E]{elements: make([]E, 0, capacity)}
}

// ArraySetOf creates a new [ArraySet] from the given elements.
// It must not be copied after first use.
// Duplicate elements are added only once.
func ArraySetOf[E comparable](elements ...E) ArraySet[E] {
	return ArraySetFromSlice(elements)
}

// ArraySetFromSlice creates a new [ArraySet] from the elements in the given slice.
// It must not be copied after first use.
// Duplicate elements in the slice are added only once.
func ArraySetFromSlice[E comparable](elements []E) ArraySet[E] {
	set := ArraySet[E]{elements: make([]E, 0, len(elements))}

	for _, element := range elements {
		if set.Contains(element) {
			continue
		}

		set.elements = append(set.elements, element)
	}

	return set
}

// Add adds the given element to the set.
// If the element is already present in the set, Add is a no-op.
func (set *ArraySet[E]) Add(element E) {
	for _, alreadyAdded := range set.elements {
		if element == alreadyAdded {
			return
		}
	}

	set.elements = append(set.elements, element)
}

// AddMultiple adds the given elements to the set. Duplicate elements are added only once, and
// elements already present in the set are not added.
func (set *ArraySet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

// AddFromSlice adds the elements from the given slice to the set. Duplicate elements are added only
// once, and elements already present in the set are not added.
func (set *ArraySet[E]) AddFromSlice(elements []E) {
	if set.elements == nil {
		set.elements = make([]E, 0, len(elements))
	}

	for _, element := range elements {
		set.Add(element)
	}
}

// AddFromSet adds elements from the given other set to the set.
func (set *ArraySet[E]) AddFromSet(otherSet ComparableSet[E]) {
	if set.elements == nil {
		set.elements = make([]E, 0, otherSet.Size())
	}

	otherSet.All()(func(element E) bool {
		set.Add(element)
		return true
	})
}

// Remove removes the given element from the set.
// If the element is not present in the set, Remove is a no-op.
func (set *ArraySet[E]) Remove(element E) {
	for i, candidate := range set.elements {
		if element == candidate {
			set.elements = append(set.elements[:i], set.elements[i+1:]...)
			return
		}
	}
}

// Clear removes all elements from the set, leaving an empty set with the same capacity as before.
func (set *ArraySet[E]) Clear() {
	set.elements = set.elements[:0]
}

// Contains checks if given element is present in the set.
func (set ArraySet[E]) Contains(element E) bool {
	for _, candidate := range set.elements {
		if element == candidate {
			return true
		}
	}

	return false
}

// Size returns the number of elements in the set.
func (set ArraySet[E]) Size() int {
	return len(set.elements)
}

// IsEmpty checks if there are 0 elements in the set.
func (set ArraySet[E]) IsEmpty() bool {
	return len(set.elements) == 0
}

// Equals checks if the set contains exactly the same elements as the other given set.
func (set ArraySet[E]) Equals(otherSet ComparableSet[E]) bool {
	return set.Size() == otherSet.Size() && set.IsSubsetOf(otherSet)
}

// IsSubsetOf checks if all of the elements in the set exist in the other given set.
func (set ArraySet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	for _, element := range set.elements {
		if !otherSet.Contains(element) {
			return false
		}
	}

	return true
}

// IsSupersetOf checks if the set contains all of the elements in the other given set.
func (set ArraySet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	return otherSet.IsSubsetOf(set)
}

// Union creates a new set that contains all the elements of the receiver set and the other given
// set. The underlying type of the returned set is an *ArraySet - to get a value type, use
// [ArraySet.UnionArraySet] instead.
func (set ArraySet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionArraySet(otherSet)
	return &union
}

// UnionArraySet creates a new ArraySet that contains all the elements of the receiver set and the
// other given set.
func (set ArraySet[E]) UnionArraySet(otherSet ComparableSet[E]) ArraySet[E] {
	union := ArraySetWithCapacity[E](set.Size() + otherSet.Size())

	for _, element := range set.elements {
		union.Add(element)
	}

	otherSet.All()(func(element E) bool {
		union.Add(element)
		return true
	})

	return union
}

// Intersection creates a new set with only the elements that exist in both the receiver set and the
// other given set. The underlying type of the returned set is an *ArraySet - to get a value type,
// use [ArraySet.IntersectionArraySet] instead.
func (set ArraySet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionArraySet(otherSet)
	return &intersection
}

// IntersectionArraySet creates a new ArraySet with only the elements that exist in both the
// receiver set and the other given set.
func (set ArraySet[E]) IntersectionArraySet(otherSet ComparableSet[E]) ArraySet[E] {
	var capacity int
	if set.Size() < otherSet.Size() {
		capacity = set.Size()
	} else {
		capacity = otherSet.Size()
	}

	intersection := ArraySetWithCapacity[E](capacity)
	for _, element := range set.elements {
		if otherSet.Contains(element) {
			intersection.Add(element)
		}
	}

	return intersection
}

// ToSlice returns a slice with all the elements in the set.
//
// Mutating the slice may invalidate the set, since it uses the same backing storage. To avoid this,
// call CopyArraySet first.
func (set ArraySet[E]) ToSlice() []E {
	return set.elements
}

// ToMap creates a map with all the set's elements as keys.
func (set ArraySet[E]) ToMap() map[E]struct{} {
	m := make(map[E]struct{}, len(set.elements))

	for _, element := range set.elements {
		m[element] = struct{}{}
	}

	return m
}

// Copy creates a new set with all the same elements and capacity as the original set.
// The underlying type of the returned set is an *ArraySet - to get a value type, use
// [ArraySet.CopyArraySet] instead.
func (set ArraySet[E]) Copy() Set[E] {
	newSet := set.CopyArraySet()
	return &newSet
}

// CopyArraySet creates a new ArraySet with all the same elements and capacity as the original set.
func (set ArraySet[E]) CopyArraySet() ArraySet[E] {
	newSet := ArraySet[E]{elements: make([]E, len(set.elements), cap(set.elements))}
	copy(newSet.elements, set.elements)
	return newSet
}

// String returns a string representation of the set, implementing [fmt.Stringer].
//
// An ArraySet of elements 1, 2 and 3 will be printed as: ArraySet{1, 2, 3}
func (set ArraySet[E]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("ArraySet{")

	for i, element := range set.elements {
		fmt.Fprint(&stringBuilder, element)

		if i < len(set.elements)-1 {
			stringBuilder.WriteString(", ")
		}
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

// All returns an [Iterator] function, which when called will loop over the elements in the set and
// call the given yield function on each element. If yield returns false, iteration stops.
//
// Since sets are unordered, iteration order is non-deterministic.
func (set ArraySet[E]) All() Iterator[E] {
	return func(yield func(element E) bool) {
		for _, element := range set.elements {
			if !yield(element) {
				break
			}
		}
	}
}
