package set

import (
	"fmt"
	"strings"
)

type DynamicSet[E comparable] struct {
	sizeThreshold int
	array         ArraySet[E]
	hash          HashSet[E]
}

var _ Set[int] = (*DynamicSet[int])(nil)
var _ ComparableSet[int] = DynamicSet[int]{}

const DefaultDynamicSetSizeThreshold = 20

func NewDynamicSet[E comparable]() DynamicSet[E] {
	return DynamicSet[E]{
		sizeThreshold: DefaultDynamicSetSizeThreshold,
		array:         ArraySet[E]{elements: nil},
		hash:          HashSet[E]{elements: nil},
	}
}

func DynamicSetWithCapacity[E comparable](capacity int) DynamicSet[E] {
	set := DynamicSet[E]{sizeThreshold: DefaultDynamicSetSizeThreshold}

	if capacity < set.sizeThreshold {
		set.array = ArraySet[E]{elements: make([]E, 0, capacity)}
	} else {
		set.hash = HashSetWithCapacity[E](capacity)
	}

	return set
}

func DynamicSetOf[E comparable](elements ...E) DynamicSet[E] {
	return DynamicSetFromSlice(elements)
}

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

func (set *DynamicSet[E]) AddMultiple(elements ...E) {
	set.AddFromSlice(elements)
}

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

func (set *DynamicSet[E]) Clear() {
	set.hash.elements = nil
}

func (set DynamicSet[E]) Contains(element E) bool {
	if set.IsArraySet() {
		return set.array.Contains(element)
	} else {
		return set.hash.Contains(element)
	}
}

func (set DynamicSet[E]) Size() int {
	if set.IsArraySet() {
		return set.array.Size()
	} else {
		return set.hash.Size()
	}
}

func (set DynamicSet[E]) IsEmpty() bool {
	if set.IsArraySet() {
		return set.array.IsEmpty()
	} else {
		return set.hash.IsEmpty()
	}
}

func (set DynamicSet[E]) Equals(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.Equals(otherSet)
	} else {
		return set.hash.Equals(otherSet)
	}
}

func (set DynamicSet[E]) IsSubsetOf(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.IsSubsetOf(otherSet)
	} else {
		return set.hash.IsSubsetOf(otherSet)
	}
}

func (set DynamicSet[E]) IsSupersetOf(otherSet ComparableSet[E]) bool {
	if set.IsArraySet() {
		return set.array.IsSupersetOf(otherSet)
	} else {
		return set.hash.IsSupersetOf(otherSet)
	}
}

func (set DynamicSet[E]) Union(otherSet ComparableSet[E]) Set[E] {
	union := set.UnionDynamicSet(otherSet)
	return &union
}

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

func (set DynamicSet[E]) Intersection(otherSet ComparableSet[E]) Set[E] {
	intersection := set.IntersectionDynamicSet(otherSet)
	return &intersection
}

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

func (set DynamicSet[E]) ToSlice() []E {
	if set.IsArraySet() {
		return set.array.ToSlice()
	} else {
		return set.hash.ToSlice()
	}
}

func (set DynamicSet[E]) ToMap() map[E]struct{} {
	if set.IsArraySet() {
		return set.array.ToMap()
	} else {
		return set.hash.ToMap()
	}
}

func (set DynamicSet[E]) ToArraySet() ArraySet[E] {
	if set.IsArraySet() {
		return set.array.CopyArraySet()
	} else {
		return set.hash.ToArraySet()
	}
}

func (set DynamicSet[E]) ToHashSet() HashSet[E] {
	if set.IsArraySet() {
		return set.array.ToHashSet()
	} else {
		return set.hash.CopyHashSet()
	}
}

func (set DynamicSet[E]) ToDynamicSet() DynamicSet[E] {
	return set.CopyDynamicSet()
}

func (set DynamicSet[E]) Copy() Set[E] {
	newSet := set.CopyDynamicSet()
	return &newSet
}

func (set DynamicSet[E]) CopyDynamicSet() DynamicSet[E] {
	newSet := DynamicSet[E]{sizeThreshold: set.sizeThreshold}

	if set.IsArraySet() {
		newSet.array = set.array.CopyArraySet()
	} else {
		newSet.hash = set.hash.CopyHashSet()
	}

	return newSet
}

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

func (set DynamicSet[E]) Iterate(yield func(element E) bool) bool {
	if set.IsArraySet() {
		return set.array.Iterate(yield)
	} else {
		return set.hash.Iterate(yield)
	}
}

func (set DynamicSet[E]) IsArraySet() bool {
	return set.hash.elements == nil
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
