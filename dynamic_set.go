package set

import (
	"fmt"
	"strings"
)

// A DynamicSet is a collection of unique elements of type E. It starts out as an [ArraySet],
// optimized for small sets. But as elements are added to it and it reaches a certain size
// threshold, it transforms itself to a [HashSet], optimized for large sets.
//
// The size threshold defaults to [DefaultDynamicSetSizeThreshold], but can be customized with
// [DynamicSet.SetSizeThreshold].
//
// The zero value for a DynamicSet is ready to use. It must not be copied after first use.
//
// DynamicSet implements [Set] when passed by pointer, and [ComparableSet] when passed by value.
type DynamicSet[E comparable] struct {
	sizeThreshold int
	array         ArraySet[E]
	hash          HashSet[E]
}

var _ Set[int] = (*DynamicSet[int])(nil)
var _ ComparableSet[int] = DynamicSet[int]{}

// DefaultDynamicSetSizeThreshold is the default size at which a DynamicSet will transform from an
// ArraySet to a HashSet. From the benchmarks in benchmark_test.go, it appears that 20 elements is
// around where HashSet.Contains performs better than ArraySet.Contains, though this varies by the
// element type of the set.
const DefaultDynamicSetSizeThreshold = 20

// NewDynamicSet creates a new [DynamicSet] for elements of type E.
// It must not be copied after first use.
func NewDynamicSet[E comparable]() DynamicSet[E] {
	return DynamicSet[E]{
		sizeThreshold: DefaultDynamicSetSizeThreshold,
		array:         ArraySet[E]{elements: nil},
		hash:          HashSet[E]{elements: nil},
	}
}

// DynamicSetWithCapacity creates a new [DynamicSet], with at least the given initial capacity.
// It must not be copied after first use.
func DynamicSetWithCapacity[E comparable](capacity int) DynamicSet[E] {
	return DynamicSet[E]{
		sizeThreshold: DefaultDynamicSetSizeThreshold,
		array:         ArraySetWithCapacity[E](capacity),
	}
}

// DynamicSetOf creates a new [DynamicSet] from the given elements.
// It must not be copied after first use.
// Duplicate elements are added only once.
func DynamicSetOf[E comparable](elements ...E) DynamicSet[E] {
	return DynamicSetFromSlice(elements)
}

// DynamicSetFromSlice creates a new [DynamicSet] from the elements in the given slice.
// It must not be copied after first use.
// Duplicate elements in the slice are added only once.
func DynamicSetFromSlice[E comparable](elements []E) DynamicSet[E] {
	set := DynamicSet[E]{
		sizeThreshold: DefaultDynamicSetSizeThreshold,
		array:         ArraySet[E]{elements: make([]E, 0, len(elements))},
	}

	set.array.AddFromSlice(elements)

	if set.arraySetReachedThreshold() {
		set.transformToHashSet()
	}

	return set
}

// SetSizeThreshold sets the size at which the DynamicSet will transform from an ArraySet to a
// HashSet.
//
// If the set is an ArraySet above the given size threshold, it transforms to a HashSet immediately.
// If the set is a HashSet below the given size threshold, it transforms to an ArraySet.
func (set *DynamicSet[E]) SetSizeThreshold(sizeThreshold int) {
	set.sizeThreshold = sizeThreshold

	if set.IsArraySet() {
		if len(set.array.elements) >= sizeThreshold {
			set.transformToHashSet()
		}
	} else {
		if len(set.hash.elements) < sizeThreshold {
			set.transformToArraySet()
		}
	}
}

// Add adds the given element to the set.
// If the element is already present in the set, Add is a no-op.
//
// If the DynamicSet is an ArraySet, it transforms to a HashSet if adding the element brings it
// above the set's size threshold.
func (set *DynamicSet[E]) Add(element E) {
	if set.IsArraySet() {
		set.array.Add(element)

		if set.arraySetReachedThreshold() {
			set.transformToHashSet()
		}
	} else {
		set.hash.Add(element)
	}
}

// AddMultiple adds the given elements to the set. Duplicate elements are added only once, and
// elements already present in the set are not added.
//
// If the DynamicSet is an ArraySet, it transforms to a HashSet if adding the elements brings it
// above the set's size threshold.
func (set *DynamicSet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

// AddFromSlice adds the elements from the given slice to the set. Duplicate elements are added
// only once, and elements already present in the set are not added.
//
// If the DynamicSet is an ArraySet, it transforms to a HashSet if adding the elements brings it
// above the set's size threshold.
func (set *DynamicSet[E]) AddFromSlice(elements []E) {
	if set.IsArraySet() {
		set.array.AddFromSlice(elements)

		if set.arraySetReachedThreshold() {
			set.transformToHashSet()
		}
	} else {
		set.hash.AddFromSlice(elements)
	}
}

// MergeWith adds elements from the given other set to the set.
//
// If the DynamicSet is an ArraySet, it transforms to a HashSet if adding the elements brings it
// above the set's size threshold.
func (set *DynamicSet[E]) MergeWith(otherSet ComparableSet[E]) {
	if set.IsArraySet() {
		set.array.MergeWith(otherSet)

		if set.arraySetReachedThreshold() {
			set.transformToHashSet()
		}
	} else {
		set.hash.MergeWith(otherSet)
	}
}

// Remove removes the given element from the set.
// If the element is not present in the set, Remove is a no-op.
//
// If the DynamicSet is a HashSet, it transforms to an ArraySet if adding the elements brings it
// below half the set's size threshold.
func (set *DynamicSet[E]) Remove(element E) {
	if set.IsArraySet() {
		set.array.Remove(element)
	} else {
		set.hash.Remove(element)

		if set.hashSetReachedThreshold() {
			set.transformToArraySet()
		}
	}
}

// Clear removes all elements from the set.
func (set *DynamicSet[E]) Clear() {
	if set.IsArraySet() {
		set.array.Clear()
	} else {
		set.hash.elements = nil
	}
}

// Contains checks if given element is present in the set.
func (set DynamicSet[E]) Contains(element E) bool {
	if set.IsArraySet() {
		return set.array.Contains(element)
	} else {
		return set.hash.Contains(element)
	}
}

// Size returns the number of elements in the set.
func (set DynamicSet[E]) Size() int {
	if set.IsArraySet() {
		return set.array.Size()
	} else {
		return set.hash.Size()
	}
}

// IsEmpty checks if there are 0 elements in the set.
func (set DynamicSet[E]) IsEmpty() bool {
	if set.IsArraySet() {
		return set.array.IsEmpty()
	} else {
		return set.hash.IsEmpty()
	}
}

// Equals checks if the set contains exactly the same elements as the other given set.
func (set DynamicSet[E]) Equals(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.Equals(otherSet)
	} else {
		return set.hash.Equals(otherSet)
	}
}

// IsSubsetOf checks if all of the elements in the set exist in the other given set.
func (set DynamicSet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.IsSubsetOf(otherSet)
	} else {
		return set.hash.IsSubsetOf(otherSet)
	}
}

// IsSupersetOf checks if the set contains all of the elements in the other given set.
func (set DynamicSet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.IsSupersetOf(otherSet)
	} else {
		return set.hash.IsSupersetOf(otherSet)
	}
}

// Union creates a new set that contains all the elements of the receiver set and the other given
// set. The underlying type of the returned set is a *DynamicSet - to get a value type, use
// [DynamicSet.UnionDynamicSet] instead.
func (set DynamicSet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionDynamicSet(otherSet)
	return &union
}

// UnionDynamicSet creates a new DynamicSet that contains all the elements of the receiver set and
// the other given set.
func (set DynamicSet[E]) UnionDynamicSet(otherSet ComparableSet[E]) DynamicSet[E] {
	union := DynamicSet[E]{sizeThreshold: set.sizeThreshold}

	if set.IsArraySet() {
		union.array = set.array.UnionArraySet(otherSet)

		if union.arraySetReachedThreshold() {
			union.transformToHashSet()
		}
	} else {
		union.hash = set.hash.UnionHashSet(otherSet)
	}

	return union
}

// Intersection creates a new set with only the elements that exist in both the receiver set and the
// other given set. The underlying type of the returned set is a *DynamicSet - to get a value type,
// use [DynamicSet.IntersectionDynamicSet] instead.
func (set DynamicSet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionDynamicSet(otherSet)
	return &intersection
}

// IntersectionDynamicSet creates a new DynamicSet with only the elements that exist in both the
// receiver set and the other given set.
func (set DynamicSet[E]) IntersectionDynamicSet(otherSet ComparableSet[E]) DynamicSet[E] {
	intersection := DynamicSet[E]{sizeThreshold: set.sizeThreshold}

	if set.IsArraySet() {
		intersection.array = set.array.IntersectionArraySet(otherSet)
	} else {
		intersection.hash = set.hash.IntersectionHashSet(otherSet)

		if intersection.hashSetReachedThreshold() {
			intersection.transformToArraySet()
		}
	}

	return intersection
}

// ToSlice creates a slice with all the elements in the set.
//
// Since sets are unordered, the order of elements in the slice is non-deterministic, and may
// vary even when called multiple times on the same set.
func (set DynamicSet[E]) ToSlice() []E {
	if set.IsArraySet() {
		return set.array.ToSlice()
	} else {
		return set.hash.ToSlice()
	}
}

// ToMap creates a map with all the set's elements as keys.
func (set DynamicSet[E]) ToMap() map[E]struct{} {
	if set.IsArraySet() {
		return set.array.ToMap()
	} else {
		return set.hash.ToMap()
	}
}

// ToArraySet creates an [ArraySet] from the elements in the set.
func (set DynamicSet[E]) ToArraySet() ArraySet[E] {
	if set.IsArraySet() {
		return set.array.CopyArraySet()
	} else {
		return set.hash.ToArraySet()
	}
}

// ToHashSet creates a [HashSet] from the elements in the set.
func (set DynamicSet[E]) ToHashSet() HashSet[E] {
	if set.IsArraySet() {
		return set.array.ToHashSet()
	} else {
		return set.hash.CopyHashSet()
	}
}

// ToHashSet is equivalent to calling [DynamicSet.CopyDynamicSet]. It is implemented to satisfy the
// [Set] interface.
func (set DynamicSet[E]) ToDynamicSet() DynamicSet[E] {
	return set.CopyDynamicSet()
}

// Copy creates a new set with all the same elements and capacity as the original set.
// The underlying type of the returned set is a *DynamicSet - to get a value type, use
// [DynamicSet.CopyDynamicSet] instead.
func (set DynamicSet[E]) Copy() Set[E] {
	newSet := set.CopyDynamicSet()
	return &newSet
}

// CopyDynamicSet creates a new DynamicSet with all the same elements and capacity as the original
// set.
func (set DynamicSet[E]) CopyDynamicSet() DynamicSet[E] {
	newSet := DynamicSet[E]{sizeThreshold: set.sizeThreshold}

	if set.IsArraySet() {
		newSet.array = set.array.CopyArraySet()
	} else {
		newSet.hash = set.hash.CopyHashSet()
	}

	return newSet
}

// String returns a string representation of the set, implementing [fmt.Stringer].
//
// Since sets are unordered, the order of elements in the string may differ each time it is
// called.
//
// A DynamicSet of elements 1, 2 and 3 will be printed as: DynamicSet{1, 2, 3} (though the order may
// vary).
func (set DynamicSet[E]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("DynamicSet{")

	if set.IsArraySet() {
		for i, element := range set.array.elements {
			fmt.Fprint(&stringBuilder, element)

			if i < len(set.array.elements)-1 {
				stringBuilder.WriteString(", ")
			}
		}
	} else {
		i := 0
		for element := range set.hash.elements {
			fmt.Fprint(&stringBuilder, element)

			if i < len(set.hash.elements)-1 {
				stringBuilder.WriteString(", ")
			}

			i++
		}
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

// Iterate loops over every element in the set, and calls the given function on it.
// It stops iteration if the function returns false.
//
// Since sets are unordered, iteration order is non-deterministic.
//
// The boolean return from Iterate is there to satisfy the future interface for [range-over-func] in
// Go, and is always false.
//
// [range-over-func]: https://github.com/golang/go/issues/61405
func (set DynamicSet[E]) Iterate(yield func(element E) bool) bool {
	if set.IsArraySet() {
		return set.array.Iterate(yield)
	} else {
		return set.hash.Iterate(yield)
	}
}

// IsArraySet checks if the DynamicSet is an ArraySet internally, i.e. that it is yet to transform
// to a HashSet due to being below its size threshold.
func (set DynamicSet[E]) IsArraySet() bool {
	return set.hash.elements == nil
}

// IsHashSet checks if the DynamicSet is a HashSet internally, i.e. that is has transformed after
// reaching its size threshold.
func (set DynamicSet[E]) IsHashSet() bool {
	return set.hash.elements != nil
}

func (set DynamicSet[E]) arraySetReachedThreshold() bool {
	return len(set.array.elements) >= set.sizeThreshold
}

func (set DynamicSet[E]) hashSetReachedThreshold() bool {
	return len(set.hash.elements) <= set.sizeThreshold/2
}

func (set *DynamicSet[E]) transformToHashSet() {
	set.hash = set.array.ToHashSet()
	set.array.elements = nil
}

func (set *DynamicSet[E]) transformToArraySet() {
	set.array = set.hash.ToArraySet()
	set.hash.elements = nil
}
