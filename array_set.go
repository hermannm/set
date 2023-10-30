package set

import (
	"fmt"
	"strings"
)

type ArraySet[E comparable] struct {
	elements []E
}

var _ Set[int] = (*ArraySet[int])(nil)
var _ ComparableSet[int] = ArraySet[int]{}

func NewArraySet[E comparable]() ArraySet[E] {
	return ArraySet[E]{elements: nil}
}

func ArraySetWithCapacity[E comparable](capacity int) ArraySet[E] {
	return ArraySet[E]{elements: make([]E, 0, capacity)}
}

func ArraySetOf[E comparable](elements ...E) ArraySet[E] {
	return ArraySetFromSlice(elements)
}

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

func (set *ArraySet[E]) Add(element E) {
	for _, alreadyAdded := range set.elements {
		if element == alreadyAdded {
			return
		}
	}

	set.elements = append(set.elements, element)
}

func (set *ArraySet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

func (set *ArraySet[E]) AddFromSlice(elements []E) {
	if set.elements == nil {
		set.elements = make([]E, 0, len(elements))
	}

	for _, element := range elements {
		set.Add(element)
	}
}

func (set *ArraySet[E]) MergeWith(otherSet ComparableSet[E]) {
	if set.elements == nil {
		set.elements = make([]E, 0, otherSet.Size())
	}

	otherSet.Iterate(func(element E) bool {
		set.Add(element)
		return true
	})
}

func (set *ArraySet[E]) Remove(element E) {
	for i, candidate := range set.elements {
		if element == candidate {
			set.elements = append(set.elements[:i], set.elements[i+1:]...)
			return
		}
	}
}

func (set *ArraySet[E]) Clear() {
	set.elements = set.elements[:0]
}

func (set ArraySet[E]) Contains(element E) bool {
	for _, candidate := range set.elements {
		if element == candidate {
			return true
		}
	}

	return false
}

func (set ArraySet[E]) Size() int {
	return len(set.elements)
}

func (set ArraySet[E]) IsEmpty() bool {
	return len(set.elements) == 0
}

func (set ArraySet[E]) Equals(otherSet ComparableSet[E]) bool {
	return set.Size() == otherSet.Size() && set.IsSubsetOf(otherSet)
}

func (set ArraySet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	for _, element := range set.elements {
		if !otherSet.Contains(element) {
			return false
		}
	}

	return true
}

func (set ArraySet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	return otherSet.IsSubsetOf(set)
}

func (set ArraySet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionArraySet(otherSet)
	return &union
}

func (set ArraySet[E]) UnionArraySet(otherSet ComparableSet[E]) ArraySet[E] {
	union := ArraySetWithCapacity[E](set.Size() + otherSet.Size())

	for _, element := range set.elements {
		union.Add(element)
	}

	otherSet.Iterate(func(element E) bool {
		union.Add(element)
		return true
	})

	return union
}

func (set ArraySet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionArraySet(otherSet)
	return &intersection
}

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

func (set ArraySet[E]) ToSlice() []E {
	slice := make([]E, len(set.elements))
	copy(slice, set.elements)
	return slice
}

func (set ArraySet[E]) ToMap() map[E]struct{} {
	m := make(map[E]struct{}, len(set.elements))

	for _, element := range set.elements {
		m[element] = struct{}{}
	}

	return m
}

func (set ArraySet[E]) ToArraySet() ArraySet[E] {
	return set.CopyArraySet()
}

func (set ArraySet[E]) ToHashSet() HashSet[E] {
	hashSet := HashSet[E]{elements: make(map[E]struct{}, len(set.elements))}

	for _, element := range set.elements {
		hashSet.elements[element] = struct{}{}
	}

	return hashSet
}

func (set ArraySet[E]) ToDynamicSet() DynamicSet[E] {
	dynamicSet := DynamicSet[E]{sizeThreshold: DefaultDynamicSetSizeThreshold}

	if len(set.elements) < dynamicSet.sizeThreshold {
		dynamicSet.array = set.CopyArraySet()
		return dynamicSet
	} else {
		dynamicSet.hash = set.ToHashSet()
		return dynamicSet
	}
}

func (set ArraySet[E]) Copy() Set[E] {
	newSet := set.CopyArraySet()
	return &newSet
}

func (set ArraySet[E]) CopyArraySet() ArraySet[E] {
	newSet := ArraySet[E]{elements: make([]E, len(set.elements), cap(set.elements))}
	copy(newSet.elements, set.elements)
	return newSet
}

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

func (set ArraySet[E]) Iterate(yield func(E) bool) bool {
	for _, element := range set.elements {
		if !yield(element) {
			return false
		}
	}

	return false
}
