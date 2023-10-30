package set

import (
	"fmt"
	"strings"
)

type HashSet[E comparable] struct {
	elements map[E]struct{}
}

var _ Set[int] = (*HashSet[int])(nil)
var _ ComparableSet[int] = HashSet[int]{}

func NewHashSet[E comparable]() HashSet[E] {
	return HashSet[E]{elements: make(map[E]struct{})}
}

func HashSetWithCapacity[E comparable](initialCapacity int) HashSet[E] {
	return HashSet[E]{elements: make(map[E]struct{}, initialCapacity)}
}

func HashSetOf[E comparable](elements ...E) HashSet[E] {
	return HashSetFromSlice(elements)
}

func HashSetFromSlice[E comparable](elements []E) HashSet[E] {
	set := HashSet[E]{elements: make(map[E]struct{}, len(elements))}

	for _, element := range elements {
		set.elements[element] = struct{}{}
	}

	return set
}

func (set *HashSet[E]) Add(element E) {
	if set.elements == nil {
		set.elements = make(map[E]struct{})
	}

	set.elements[element] = struct{}{}
}

func (set *HashSet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

func (set *HashSet[E]) AddFromSlice(elements []E) {
	if set.elements == nil {
		set.elements = make(map[E]struct{}, len(elements))
	}

	for _, element := range elements {
		set.elements[element] = struct{}{}
	}
}

func (set *HashSet[E]) MergeWith(otherSet ComparableSet[E]) {
	if set.elements == nil {
		set.elements = make(map[E]struct{}, otherSet.Size())
	}

	otherSet.Iterate(func(element E) bool {
		set.Add(element)
		return true
	})
}

func (set HashSet[E]) Remove(element E) {
	delete(set.elements, element)
}

func (set HashSet[E]) Clear() {
	for element := range set.elements {
		delete(set.elements, element)
	}
}

func (set HashSet[E]) Contains(element E) bool {
	if set.elements == nil {
		return false
	}

	_, contains := set.elements[element]
	return contains
}

func (set HashSet[E]) Size() int {
	return len(set.elements)
}

func (set HashSet[E]) IsEmpty() bool {
	return len(set.elements) == 0
}

func (set HashSet[E]) Equals(otherSet ComparableSet[E]) bool {
	return set.Size() == otherSet.Size() && set.IsSubsetOf(otherSet)
}

func (set HashSet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	for element := range set.elements {
		if !otherSet.Contains(element) {
			return false
		}
	}

	return true
}

func (set HashSet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	return otherSet.IsSubsetOf(set)
}

func (set HashSet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionHashSet(otherSet)
	return &union
}

func (set HashSet[E]) UnionHashSet(otherSet ComparableSet[E]) HashSet[E] {
	union := HashSetWithCapacity[E](set.Size() + otherSet.Size())

	for element := range set.elements {
		union.Add(element)
	}

	otherSet.Iterate(func(element E) bool {
		union.Add(element)
		return true
	})

	return union
}

func (set HashSet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionHashSet(otherSet)
	return &intersection
}

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

func (set HashSet[E]) ToSlice() []E {
	slice := make([]E, len(set.elements))

	i := 0
	for element := range set.elements {
		slice[i] = element
		i++
	}

	return slice
}

func (set HashSet[E]) ToMap() map[E]struct{} {
	m := make(map[E]struct{}, len(set.elements))

	for element := range set.elements {
		m[element] = struct{}{}
	}

	return m
}

func (set HashSet[E]) ToArraySet() ArraySet[E] {
	arraySet := ArraySet[E]{elements: make([]E, len(set.elements))}

	i := 0
	for element := range set.elements {
		arraySet.elements[i] = element
		i++
	}

	return arraySet
}

func (set HashSet[E]) ToHashSet() HashSet[E] {
	return set.CopyHashSet()
}

func (set HashSet[E]) ToDynamicSet() DynamicSet[E] {
	dynamicSet := DynamicSet[E]{sizeThreshold: DefaultDynamicSetSizeThreshold}

	if len(set.elements) >= dynamicSet.sizeThreshold {
		dynamicSet.hash = set.CopyHashSet()
		return dynamicSet
	} else {
		dynamicSet.array = set.ToArraySet()
		return dynamicSet
	}
}

func (set HashSet[E]) Copy() Set[E] {
	newSet := set.CopyHashSet()
	return &newSet
}

func (set HashSet[E]) CopyHashSet() HashSet[E] {
	newSet := HashSet[E]{elements: make(map[E]struct{}, len(set.elements))}

	for element := range set.elements {
		newSet.elements[element] = struct{}{}
	}

	return newSet
}

func (set HashSet[E]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("HashSet{")

	i := 0
	for element := range set.elements {
		fmt.Fprint(&stringBuilder, element)

		if i < len(set.elements)-1 {
			stringBuilder.WriteString(", ")
		}

		i++
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

func (set HashSet[E]) Iterate(yield func(E) bool) bool {
	for element := range set.elements {
		if !yield(element) {
			return false
		}
	}

	return false
}
