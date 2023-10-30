package set

import (
	"fmt"
	"strings"
)

type DynamicSet[T comparable] struct {
	resizeCutoff int
	array        ArraySet[T]
	hash         HashSet[T]
}

var _ Set[int] = (*DynamicSet[int])(nil)
var _ ComparableSet[int] = DynamicSet[int]{}

const DefaultDynamicSetResizeCutoff = 20

func NewDynamicSet[T comparable]() DynamicSet[T] {
	return DynamicSet[T]{
		resizeCutoff: DefaultDynamicSetResizeCutoff,
		array:        ArraySet[T]{items: nil},
		hash:         HashSet[T]{items: nil},
	}
}

func DynamicSetWithCapacity[T comparable](capacity int) DynamicSet[T] {
	set := DynamicSet[T]{resizeCutoff: DefaultDynamicSetResizeCutoff}

	if capacity < set.resizeCutoff {
		set.array = ArraySet[T]{items: make([]T, 0, capacity)}
	} else {
		set.hash = HashSetWithCapacity[T](capacity)
	}

	return set
}

func DynamicSetOf[T comparable](items ...T) DynamicSet[T] {
	return DynamicSetFromSlice(items)
}

func DynamicSetFromSlice[T comparable](items []T) DynamicSet[T] {
	set := DynamicSet[T]{resizeCutoff: DefaultDynamicSetResizeCutoff}

	if len(items) < set.resizeCutoff {
		set.array = ArraySet[T]{items: make([]T, 0, len(items))}

		for _, item := range items {
			if set.array.Contains(item) {
				continue
			}

			set.array.items = append(set.array.items, item)
		}
	} else {
		set.hash = HashSetFromSlice(items)
	}

	return set
}

func (set *DynamicSet[T]) SetResizeCutoff(resizeCutoff int) {
	set.resizeCutoff = resizeCutoff
}

func (set *DynamicSet[T]) Add(item T) {
	if set.isArraySet() {
		set.array.Add(item)

		if set.arraySetReachedCutoff() {
			set.transformToHashSet()
		}
	} else {
		set.hash.Add(item)
	}
}

func (set *DynamicSet[T]) AddMultiple(items ...T) {
	set.AddFromSlice(items)
}

func (set *DynamicSet[T]) AddFromSlice(items []T) {
	if set.isArraySet() {
		set.array.AddFromSlice(items)

		if set.arraySetReachedCutoff() {
			set.transformToHashSet()
		}
	} else {
		set.hash.AddFromSlice(items)
	}
}

func (set *DynamicSet[T]) MergeWith(other ComparableSet[T]) {
	if set.isArraySet() {
		set.array.MergeWith(other)

		if set.arraySetReachedCutoff() {
			set.transformToHashSet()
		}
	} else {
		set.hash.MergeWith(other)
	}
}

func (set *DynamicSet[T]) Remove(item T) {
	if set.isArraySet() {
		set.array.Remove(item)
	} else {
		set.hash.Remove(item)

		if set.hashSetReachedCutoff() {
			set.transformToArraySet()
		}
	}
}

func (set *DynamicSet[T]) Clear() {
	set.hash.items = nil
}

func (set DynamicSet[T]) Contains(item T) bool {
	if set.isArraySet() {
		return set.array.Contains(item)
	} else {
		return set.hash.Contains(item)
	}
}

func (set DynamicSet[T]) Size() int {
	if set.isArraySet() {
		return set.array.Size()
	} else {
		return set.hash.Size()
	}
}

func (set DynamicSet[T]) IsEmpty() bool {
	if set.isArraySet() {
		return set.array.IsEmpty()
	} else {
		return set.hash.IsEmpty()
	}
}

func (set DynamicSet[T]) Equals(other ComparableSet[T]) bool {
	if set.isArraySet() {
		return set.array.Equals(other)
	} else {
		return set.hash.Equals(other)
	}
}

func (set DynamicSet[T]) IsSubsetOf(other ComparableSet[T]) bool {
	if set.isArraySet() {
		return set.array.IsSubsetOf(other)
	} else {
		return set.hash.IsSubsetOf(other)
	}
}

func (set DynamicSet[T]) IsSupersetOf(other ComparableSet[T]) bool {
	if set.isArraySet() {
		return set.array.IsSupersetOf(other)
	} else {
		return set.hash.IsSupersetOf(other)
	}
}

func (set DynamicSet[T]) Union(other ComparableSet[T]) Set[T] {
	union := set.UnionDynamicSet(other)
	return &union
}

func (set DynamicSet[T]) UnionDynamicSet(other ComparableSet[T]) DynamicSet[T] {
	union := DynamicSet[T]{resizeCutoff: set.resizeCutoff}

	if set.isArraySet() {
		union.array = set.array.UnionArraySet(other)

		if union.arraySetReachedCutoff() {
			union.transformToHashSet()
		}
	} else {
		union.hash = set.hash.UnionHashSet(other)
	}

	return union
}

func (set DynamicSet[T]) Intersection(other ComparableSet[T]) Set[T] {
	intersection := set.IntersectionDynamicSet(other)
	return &intersection
}

func (set DynamicSet[T]) IntersectionDynamicSet(other ComparableSet[T]) DynamicSet[T] {
	intersection := DynamicSet[T]{resizeCutoff: set.resizeCutoff}

	if set.isArraySet() {
		intersection.array = set.array.IntersectionArraySet(other)
	} else {
		intersection.hash = set.hash.IntersectionHashSet(other)

		if intersection.hashSetReachedCutoff() {
			intersection.transformToArraySet()
		}
	}

	return intersection
}

func (set DynamicSet[T]) ToSlice() []T {
	if set.isArraySet() {
		return set.array.ToSlice()
	} else {
		return set.hash.ToSlice()
	}
}

func (set DynamicSet[T]) ToMap() map[T]struct{} {
	if set.isArraySet() {
		return set.array.ToMap()
	} else {
		return set.hash.ToMap()
	}
}

func (set DynamicSet[T]) ToArraySet() ArraySet[T] {
	if set.isArraySet() {
		return set.array.CopyArraySet()
	} else {
		return set.hash.ToArraySet()
	}
}

func (set DynamicSet[T]) ToHashSet() HashSet[T] {
	if set.isArraySet() {
		return set.array.ToHashSet()
	} else {
		return set.hash.CopyHashSet()
	}
}

func (set DynamicSet[T]) ToDynamicSet() DynamicSet[T] {
	return set.CopyDynamicSet()
}

func (set DynamicSet[T]) Copy() Set[T] {
	newSet := set.CopyDynamicSet()
	return &newSet
}

func (set DynamicSet[T]) CopyDynamicSet() DynamicSet[T] {
	newSet := DynamicSet[T]{resizeCutoff: set.resizeCutoff}

	if set.isArraySet() {
		newSet.array = set.array.CopyArraySet()
	} else {
		newSet.hash = set.hash.CopyHashSet()
	}

	return newSet
}

func (set DynamicSet[T]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("DynamicSet{")

	if set.isArraySet() {
		for i, item := range set.array.items {
			fmt.Fprint(&stringBuilder, item)

			if i < len(set.array.items)-1 {
				stringBuilder.WriteString(", ")
			}
		}
	} else {
		i := 0
		for item := range set.hash.items {
			fmt.Fprint(&stringBuilder, item)

			if i < len(set.hash.items)-1 {
				stringBuilder.WriteString(", ")
			}

			i++
		}
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

func (set DynamicSet[T]) Iterate(yield func(item T) bool) bool {
	if set.isArraySet() {
		return set.array.Iterate(yield)
	} else {
		return set.hash.Iterate(yield)
	}
}

func (set DynamicSet[T]) isArraySet() bool {
	return set.hash.items == nil
}

func (set DynamicSet[T]) arraySetReachedCutoff() bool {
	return len(set.array.items) >= set.resizeCutoff
}

func (set DynamicSet[T]) hashSetReachedCutoff() bool {
	return len(set.hash.items) <= set.resizeCutoff/2
}

func (set *DynamicSet[T]) transformToHashSet() {
	set.hash = set.array.ToHashSet()
	set.array.items = nil
}

func (set *DynamicSet[T]) transformToArraySet() {
	set.array = set.hash.ToArraySet()
	set.hash.items = nil
}
