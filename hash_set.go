package set

import (
	"fmt"
	"strings"
)

type HashSet[T comparable] struct {
	items map[T]struct{}
}

var _ Set[int] = (*HashSet[int])(nil)
var _ ComparableSet[int] = HashSet[int]{}

func NewHashSet[T comparable]() HashSet[T] {
	return HashSet[T]{items: make(map[T]struct{})}
}

func HashSetWithCapacity[T comparable](initialCapacity int) HashSet[T] {
	return HashSet[T]{items: make(map[T]struct{}, initialCapacity)}
}

func HashSetOf[T comparable](items ...T) HashSet[T] {
	return HashSetFromSlice(items)
}

func HashSetFromSlice[T comparable](items []T) HashSet[T] {
	set := HashSet[T]{items: make(map[T]struct{}, len(items))}

	for _, item := range items {
		set.items[item] = struct{}{}
	}

	return set
}

func (set *HashSet[T]) Add(item T) {
	if set.items == nil {
		set.items = make(map[T]struct{})
	}

	set.items[item] = struct{}{}
}

func (set *HashSet[T]) AddMultiple(items ...T) {
	set.AddFromSlice(items)
}

func (set *HashSet[T]) AddFromSlice(items []T) {
	if set.items == nil {
		set.items = make(map[T]struct{}, len(items))
	}

	for _, item := range items {
		set.items[item] = struct{}{}
	}
}

func (set *HashSet[T]) MergeWith(other ComparableSet[T]) {
	if set.items == nil {
		set.items = make(map[T]struct{}, other.Size())
	}

	other.Iterate(func(item T) bool {
		set.Add(item)
		return true
	})
}

func (set HashSet[T]) Remove(item T) {
	delete(set.items, item)
}

func (set HashSet[T]) Clear() {
	for item := range set.items {
		delete(set.items, item)
	}
}

func (set HashSet[T]) Contains(item T) bool {
	if set.items == nil {
		return false
	}

	_, contains := set.items[item]
	return contains
}

func (set HashSet[T]) Size() int {
	return len(set.items)
}

func (set HashSet[T]) IsEmpty() bool {
	return len(set.items) == 0
}

func (set HashSet[T]) Equals(other ComparableSet[T]) bool {
	return set.Size() == other.Size() && set.IsSubsetOf(other)
}

func (set HashSet[T]) IsSubsetOf(other ComparableSet[T]) bool {
	for item := range set.items {
		if !other.Contains(item) {
			return false
		}
	}

	return true
}

func (set HashSet[T]) IsSupersetOf(other ComparableSet[T]) bool {
	return other.IsSubsetOf(set)
}

func (set HashSet[T]) Union(other ComparableSet[T]) Set[T] {
	union := set.UnionHashSet(other)
	return &union
}

func (set HashSet[T]) UnionHashSet(other ComparableSet[T]) HashSet[T] {
	union := HashSetWithCapacity[T](set.Size() + other.Size())

	for item := range set.items {
		union.Add(item)
	}

	other.Iterate(func(item T) bool {
		union.Add(item)
		return true
	})

	return union
}

func (set HashSet[T]) Intersection(other ComparableSet[T]) Set[T] {
	intersection := set.IntersectionHashSet(other)
	return &intersection
}

func (set HashSet[T]) IntersectionHashSet(other ComparableSet[T]) HashSet[T] {
	var capacity int
	if set.Size() < other.Size() {
		capacity = set.Size()
	} else {
		capacity = other.Size()
	}

	intersection := HashSetWithCapacity[T](capacity)
	for item := range set.items {
		if other.Contains(item) {
			intersection.Add(item)
		}
	}

	return intersection
}

func (set HashSet[T]) ToSlice() []T {
	slice := make([]T, len(set.items))

	i := 0
	for item := range set.items {
		slice[i] = item
		i++
	}

	return slice
}

func (set HashSet[T]) ToMap() map[T]struct{} {
	m := make(map[T]struct{}, len(set.items))

	for item := range set.items {
		m[item] = struct{}{}
	}

	return m
}

func (set HashSet[T]) ToArraySet() ArraySet[T] {
	arraySet := ArraySet[T]{items: make([]T, len(set.items))}

	i := 0
	for item := range set.items {
		arraySet.items[i] = item
		i++
	}

	return arraySet
}

func (set HashSet[T]) ToHashSet() HashSet[T] {
	return set.CopyHashSet()
}

func (set HashSet[T]) ToDynamicSet() DynamicSet[T] {
	dynamicSet := DynamicSet[T]{resizeCutoff: DefaultDynamicSetResizeCutoff}

	if len(set.items) >= dynamicSet.resizeCutoff {
		dynamicSet.hash = set.CopyHashSet()
		return dynamicSet
	} else {
		dynamicSet.array = set.ToArraySet()
		return dynamicSet
	}
}

func (set HashSet[T]) Copy() Set[T] {
	newSet := set.CopyHashSet()
	return &newSet
}

func (set HashSet[T]) CopyHashSet() HashSet[T] {
	newSet := HashSet[T]{items: make(map[T]struct{}, len(set.items))}

	for item := range set.items {
		newSet.items[item] = struct{}{}
	}

	return newSet
}

func (set HashSet[T]) String() string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("HashSet{")

	i := 0
	for item := range set.items {
		fmt.Fprint(&stringBuilder, item)

		if i < len(set.items)-1 {
			stringBuilder.WriteString(", ")
		}

		i++
	}

	stringBuilder.WriteByte('}')
	return stringBuilder.String()
}

func (set HashSet[T]) Iterate(yield func(T) bool) bool {
	for item := range set.items {
		if !yield(item) {
			return false
		}
	}

	return false
}
