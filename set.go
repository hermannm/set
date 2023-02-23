// Package set provides a generic Set data structure, an unordered collection of unique elements.
package set

import "fmt"

// A Set is an unordered collection of unique elements, whose type is given by the Element type
// parameter.
//
// The zero value of a Set is nil, and calling Add on that will panic. Therefore, it should always
// be instantiated with one of the constructors in this package: New, WithCapacity, Of or FromSlice.
//
// Since Set is implemented as a type alias for a map, for-range loops can be used to iterate over a
// Set. Because sets are unordered, the iteration order is non-deterministic. The example below may
// print 1, 2, 3 in different order each time it runs.
//
//	uniqueNumbers := set.Of(1, 2, 3)
//	for number := range uniqueNumbers {
//	    fmt.Println(number)
//	}
type Set[Element comparable] map[Element]struct{}

// New creates an empty set for elements of the given type.
func New[Element comparable]() Set[Element] {
	return make(Set[Element])
}

// WithCapacity creates an empty set for elements of the given type.
// Allocates initial space for the given number of elements.
func WithCapacity[Element comparable](initialCapacity int) Set[Element] {
	return make(Set[Element], initialCapacity)
}

// Of creates a set of the given elements.
func Of[Element comparable](elements ...Element) Set[Element] {
	set := WithCapacity[Element](len(elements))
	set.Add(elements...)
	return set
}

// FromSlice creates a set of the unique elements from the given slice.
func FromSlice[Element comparable](slice []Element) Set[Element] {
	return Of(slice...)
}

// Add adds the given elements to the set.
// Since the set holds unique values, adding an already present element does not change the set.
//
// If the set is nil, Add will panic. To avoid this, ensure the set is instantiated with one of the
// constructors in this package: New, WithCapacity, Of or FromSlice.
func (set Set[Element]) Add(elements ...Element) {
	// Calling Add on a nil Set will give the panic message "assignment to entry in nil map".
	// Since users of the library may not know that a Set is implemented as a map, the panic is
	// intercepted here to give a more descriptive message.
	defer func() {
		if err := recover(); err != nil {
			panic("called Add on nil Set")
		}
	}()

	for _, element := range elements {
		set[element] = struct{}{}
	}
}

// Remove deletes the given element from the set.
// If the element is not present in the set, Remove is a no-op.
func (set Set[Element]) Remove(element Element) {
	delete(set, element)
}

// Clear removes all elements from the set, leaving an empty set with the same capacity as before.
func (set Set[Element]) Clear() {
	for element := range set {
		set.Remove(element)
	}
}

// Size returns the number of elements in the set.
func (set Set[Element]) Size() int {
	return len(set)
}

// IsEmpty checks if there are 0 elements in the set.
func (set Set[Element]) IsEmpty() bool {
	return set.Size() == 0
}

// Contains checks if the given element exists in the set.
func (set Set[Element]) Contains(element Element) bool {
	_, contains := set[element]
	return contains
}

// Equals checks if the set contains exactly the same elements as the other given set.
func (set1 Set[Element]) Equals(set2 Set[Element]) bool {
	// If both sets are the same size, and one is a subset of the other, then they are equal.
	return set1.Size() == set2.Size() && set1.IsSubsetOf(set2)
}

// IsSubsetOf checks if all of the elements in the set exist in the other given set.
func (set1 Set[Element]) IsSubsetOf(set2 Set[Element]) bool {
	for element := range set1 {
		if !set2.Contains(element) {
			return false
		}
	}

	return true
}

// IsSupersetOf checks if the set contains all of the elements in the other given set.
func (set1 Set[Element]) IsSupersetOf(set2 Set[Element]) bool {
	return set2.IsSubsetOf(set1)
}

// ToSlice returns a slice with all the elements in the set.
//
// Since sets are unordered, the order of elements in the slice is non-deterministic, and may vary
// even when called multiple times on the same set.
func (set Set[Element]) ToSlice() []Element {
	slice := make([]Element, 0, set.Size())

	for element := range set {
		slice = append(slice, element)
	}

	return slice
}

// Copy creates a new set with all the same elements as the set it is called on.
func (set Set[Element]) Copy() Set[Element] {
	copy := WithCapacity[Element](set.Size())

	for element := range set {
		copy.Add(element)
	}

	return copy
}

// String implements fmt.Stringer to customize the print format of Set.
//
// Since sets are unordered, the order of elements in the string may differ each time it is called.
func (set Set[Element]) String() string {
	if set.Size() == 0 {
		return "Set{}"
	}

	setString := "Set{"
	for element := range set {
		setString += fmt.Sprintf("%v, ", element)
	}
	setString = setString[:len(setString)-2] // Removes ", " after the last element.
	setString += "}"

	return setString
}

// Union creates a new set that contains all the elements of both of the given sets.
func Union[Element comparable](set1 Set[Element], set2 Set[Element]) Set[Element] {
	union := WithCapacity[Element](set1.Size() + set2.Size())

	for element := range set1 {
		union.Add(element)
	}
	for element := range set2 {
		union.Add(element)
	}

	return union
}

// Intersection creates a new set with only the elements that exist in both of the given sets.
func Intersection[Element comparable](set1 Set[Element], set2 Set[Element]) Set[Element] {
	size1 := set1.Size()
	size2 := set2.Size()

	var initialCapacity int
	if size1 < size2 {
		initialCapacity = size1
	} else {
		initialCapacity = size2
	}

	intersection := WithCapacity[Element](initialCapacity)

	for element := range set1 {
		if set2.Contains(element) {
			intersection.Add(element)
		}
	}

	return intersection
}
