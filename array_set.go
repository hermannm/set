package set

import (
	"fmt"
	"strings"
)

type ArraySet[T comparable] struct {
	items []T
}

var _ Set[int] = (*ArraySet[int])(nil)
var _ ComparableSet[int] = ArraySet[int]{}

func NewArraySet[T comparable]() ArraySet[T] {
	return ArraySet[T]{items: nil}
}

func ArraySetWithCapacity[T comparable](capacity int) ArraySet[T] {
	return ArraySet[T]{items: make([]T, 0, capacity)}
}

func ArraySetOf[T comparable](items ...T) ArraySet[T] {
	return ArraySetFromSlice(items)
}

func ArraySetFromSlice[T comparable](items []T) ArraySet[T] {
	set := ArraySet[T]{items: make([]T, 0, len(items))}

	for _, item := range items {
		if set.Contains(item) {
			continue
		}

		set.items = append(set.items, item)
	}

	return set
}

func (set *ArraySet[T]) Add(item T) {
	for _, alreadyAdded := range set.items {
		if item == alreadyAdded {
			return
		}
	}

	set.items = append(set.items, item)
}

func (set *ArraySet[T]) AddMultiple(items ...T) {
	set.AddFromSlice(items)
}

func (set *ArraySet[T]) AddFromSlice(items []T) {
	if set.items == nil {
		set.items = make([]T, 0, len(items))
	}

	for _, item := range items {
		set.Add(item)
	}
}

func (set *ArraySet[T]) MergeWith(otherSet ComparableSet[T]) {
	if set.items == nil {
		set.items = make([]T, 0, otherSet.Size())
	}

	otherSet.Iterate(func(item T) bool {
		set.Add(item)
		return true
	})
}

func (set *ArraySet[T]) Remove(item T) {
	for i, candidate := range set.items {
		if item == candidate {
			set.items = append(set.items[:i], set.items[i+1:]...)
			return
		}
	}
}

func (set *ArraySet[T]) Clear() {
	set.items = set.items[:0]
}

func (set ArraySet[T]) Contains(item T) bool {
	for _, candidate := range set.items {
		if item == candidate {
			return true
		}
	}

	return false
}

func (set ArraySet[T]) Size() int {
	return len(set.items)
}

func (set ArraySet[T]) IsEmpty() bool {
	return len(set.items) == 0
}

func (set ArraySet[T]) Equals(otherSet ComparableSet[T]) bool {
	return set.Size() == otherSet.Size() && set.IsSubsetOf(otherSet)
}

func (set ArraySet[T]) IsSubsetOf(otherSet ComparableSet[T]) bool {
	for _, item := range set.items {
		if !otherSet.Contains(item) {
			return false
		}
	}

	return true
}

func (set ArraySet[T]) IsSupersetOf(otherSet ComparableSet[T]) bool {
	return otherSet.IsSubsetOf(set)
}

func (set ArraySet[T]) Union(otherSet ComparableSet[T]) Set[T] {
	union := set.UnionArraySet(otherSet)
	return &union
}

func (set ArraySet[T]) UnionArraySet(otherSet ComparableSet[T]) ArraySet[T] {
	union := ArraySetWithCapacity[T](set.Size() + otherSet.Size())

	for _, item := range set.items {
		union.Add(item)
	}

	otherSet.Iterate(func(item T) bool {
		union.Add(item)
		return true
	})

	return union
}

func (set ArraySet[T]) Intersection(otherSet ComparableSet[T]) Set[T] {
	intersection := set.IntersectionArraySet(otherSet)
	return &intersection
}

func (set ArraySet[T]) IntersectionArraySet(otherSet ComparableSet[T]) ArraySet[T] {
	var capacity int
	if set.Size() < otherSet.Size() {
		capacity = set.Size()
	} else {
		capacity = otherSet.Size()
	}

	intersection := ArraySetWithCapacity[T](capacity)
	for _, item := range set.items {
		if otherSet.Contains(item) {
			intersection.Add(item)
		}
	}

	return intersection
}

func (set ArraySet[T]) ToSlice() []T {
	slice := make([]T, len(set.items))
	copy(slice, set.items)
	return slice
}

func (set ArraySet[T]) ToMap() map[T]struct{} {
	m := make(map[T]struct{}, len(set.items))

	for _, item := range set.items {
		m[item] = struct{}{}
	}

	return m
}

func (set ArraySet[T]) ToArraySet() ArraySet[T] {
	return set.CopyArraySet()
}

func (set ArraySet[T]) ToHashSet() HashSet[T] {
	hashSet := HashSet[T]{items: make(map[T]struct{}, len(set.items))}

	for _, item := range set.items {
		hashSet.items[item] = struct{}{}
	}

	return hashSet
}

func (set ArraySet[T]) ToDynamicSet() DynamicSet[T] {
	dynamicSet := DynamicSet[T]{sizeThreshold: DefaultDynamicSetSizeThreshold}

	if len(set.items) < dynamicSet.sizeThreshold {
		dynamicSet.array = set.CopyArraySet()
		return dynamicSet
	} else {
		dynamicSet.hash = set.ToHashSet()
		return dynamicSet
	}
}

func (set ArraySet[T]) Copy() Set[T] {
	newSet := set.CopyArraySet()
	return &newSet
}

func (set ArraySet[T]) CopyArraySet() ArraySet[T] {
	newSet := ArraySet[T]{items: make([]T, len(set.items), cap(set.items))}
	copy(newSet.items, set.items)
	return newSet
}

func (set ArraySet[T]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("ArraySet{")

	for i, item := range set.items {
		fmt.Fprint(&stringBuilder, item)

		if i < len(set.items)-1 {
			stringBuilder.WriteString(", ")
		}
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

func (set ArraySet[T]) Iterate(yield func(T) bool) bool {
	for _, item := range set.items {
		if !yield(item) {
			return false
		}
	}

	return false
}
