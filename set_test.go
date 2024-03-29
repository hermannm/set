package set_test

import (
	"testing"

	"hermannm.dev/set"
)

func TestNew(t *testing.T) {
	for _, set := range []set.ComparableSet[int]{
		set.NewArraySet[int](),
		set.NewHashSet[int](),
		set.NewDynamicSet[int](),
	} {
		assertSize(t, set, 0)
	}
}

func TestWithCapacity(t *testing.T) {
	for _, set := range []set.ComparableSet[int]{
		set.ArraySetWithCapacity[int](5),
		set.HashSetWithCapacity[int](5),
		set.DynamicSetWithCapacity[int](5),
	} {
		assertSize(t, set, 0)
	}
}

func TestOf(t *testing.T) {
	for _, set := range []set.ComparableSet[int]{
		set.ArraySetOf(1, 2, 3),
		set.HashSetOf(1, 2, 3),
		set.DynamicSetOf(1, 2, 3),
	} {
		assertSize(t, set, 3)
		assertContains(t, set, 1, 2, 3)
	}
}

func TestFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, set := range []set.ComparableSet[int]{
		set.ArraySetFromSlice(slice),
		set.HashSetFromSlice(slice),
		set.DynamicSetFromSlice(slice),
	} {
		assertSize(t, set, len(slice))
		assertContains(t, set, slice...)
	}
}

func TestFromSliceWithDuplicates(t *testing.T) {
	slice := []int{1, 1, 2, 2}

	for _, set := range []set.ComparableSet[int]{
		set.ArraySetFromSlice(slice),
		set.HashSetFromSlice(slice),
		set.DynamicSetFromSlice(slice),
	} {
		assertSize(t, set, 2)
		assertContains(t, set, 1, 2)
	}
}

func TestAdd(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.Add(1)

		assertSize(t, set, 1)
		assertContains(t, set, 1)
	})
}

func TestAddDuplicate(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)
		set.Add(3)

		assertSize(t, set, 3)
	})
}

func TestAddMultiple(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		assertSize(t, set, 3)
		assertContains(t, set, 1, 2, 3)
	})
}

func TestAddFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 3}

	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddFromSlice(slice)

		assertSize(t, set, 3)
		assertContains(t, set, 1, 2, 3)
	})
}

func TestAddFromSet(t *testing.T) {
	otherSet := set.ArraySetOf(3, 4, 5)

	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		set.AddFromSet(otherSet)

		assertSize(t, set, 5)
		assertContains(t, set, 1, 2, 3, 4, 5)
	})
}

func TestRemove(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		set.Remove(3)

		assertSize(t, set, 2)
		assertContains(t, set, 1, 2)
	})
}

func TestRemoveNonExisting(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		set.Remove(4)

		assertSize(t, set, 3)
		assertContains(t, set, 1, 2, 3)
	})
}

func TestClear(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		set.Clear()

		assertSize(t, set, 0)
	})
}

func TestSize(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		if size := set.Size(); size != 3 {
			t.Errorf("expected %v to have size 3, got %d", set, size)
		}
	})
}

func TestIsEmpty(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		if !set.IsEmpty() {
			t.Errorf("expected %v.IsEmpty() == true", set)
		}

		set.Add(1)

		if set.IsEmpty() {
			t.Errorf("expected %v.IsEmpty() == false", set)
		}
	})
}

func TestContains(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		if !set.Contains(3) {
			t.Errorf("expected %v.Contains(3) == true", set)
		}

		if set.Contains(4) {
			t.Errorf("expected %v.Contains(4) == false", set)
		}
	})
}

func TestEquals(t *testing.T) {
	testAllSetTypes(func(set1 set.Set[int], setName string) {
		set1.AddMultiple(1, 2, 3)

		set2 := set.ArraySetOf(1, 2, 3)

		if !set1.Equals(set2) {
			t.Errorf("expected %v.Equals(%v) == true", set1, set2)
		}

		set3 := set.ArraySetOf(1, 2, 4)

		if set1.Equals(set3) {
			t.Errorf("expected %v.Equals(%v) == false", set1, set3)
		}
	})
}

func TestIsSubsetOf(t *testing.T) {
	testAllSetTypes(func(set1 set.Set[int], setName string) {
		set1.AddMultiple(1, 2, 3)
		set2 := set.HashSetOf(1, 2, 3, 4, 5, 6)

		if !set1.IsSubsetOf(set2) {
			t.Errorf("expected %v.IsSubsetOf(%v) == true", set1, set2)
		}

		if set2.IsSubsetOf(set1) {
			t.Errorf("expected %v.IsSubsetOf(%v) == false", set2, set1)
		}
	})
}

func TestIsSupersetOf(t *testing.T) {
	testAllSetTypes(func(set1 set.Set[int], setName string) {
		set1.AddMultiple(1, 2, 3, 4, 5, 6)
		set2 := set.ArraySetOf(1, 2, 3)

		if !set1.IsSupersetOf(set2) {
			t.Errorf("expected %v.IsSupersetOf(%v) == true", set1, set2)
		}

		if set2.IsSupersetOf(set1) {
			t.Errorf("expected %v.IsSupersetOf(%v) == false", set2, set1)
		}
	})
}

func TestUnion(t *testing.T) {
	testAllSetTypes(func(set1 set.Set[int], setName string) {
		set1.AddMultiple(1, 2, 3)
		set2 := set.ArraySetOf(3, 4, 5)

		union := set1.Union(set2)

		assertSize(t, union, 5)
		assertContains(t, union, 1, 2, 3, 4, 5)
	})
}

func TestIntersection(t *testing.T) {
	testAllSetTypes(func(set1 set.Set[int], setName string) {
		set1.AddMultiple(1, 2, 3, 4)
		set2 := set.HashSetOf(2, 3, 4, 5)

		intersection := set1.Intersection(set2)

		assertSize(t, intersection, 3)
		assertContains(t, intersection, 2, 3, 4)
	})
}

func TestToSlice(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)
		slice := set.ToSlice()

		if len(slice) != set.Size() {
			t.Errorf(
				"expected len(%v) == %v.Size(), but got %d and %d",
				slice,
				set,
				len(slice),
				set.Size(),
			)
		}

		set.All()(func(setElement int) bool {
			containedInSlice := false

			for _, sliceElement := range slice {
				if setElement == sliceElement {
					containedInSlice = true
					break
				}
			}

			if !containedInSlice {
				t.Errorf(
					"expected %v to contain all elements of %v, but did not contain %v",
					slice,
					set,
					setElement,
				)
			}

			return true
		})
	})
}

func TestToMap(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)
		m := set.ToMap()

		if len(m) != set.Size() {
			t.Errorf(
				"expected len(%v) == %v.Size(), but got %d and %d",
				m,
				set,
				len(m),
				set.Size(),
			)
		}

		set.All()(func(element int) bool {
			if _, containedInMap := m[element]; !containedInMap {
				t.Errorf(
					"expected %v to contain all elements of %v, but did not contain %v",
					m,
					set,
					element,
				)
			}
			return true
		})
	})
}

func TestCopy(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)
		setCopy := set.Copy()

		assertContains(t, setCopy, 1, 2, 3)

		set.Add(4)

		assertSize(t, setCopy, 3)
	})
}

func TestString(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)

		setString := set.String()
		expectedStrings := []string{
			setName + "{1, 2, 3}",
			setName + "{1, 3, 2}",
			setName + "{2, 1, 3}",
			setName + "{2, 3, 1}",
			setName + "{3, 1, 2}",
			setName + "{3, 2, 1}",
		}

		isExpectedString := false
		for _, expected := range expectedStrings {
			if setString == expected {
				isExpectedString = true
			}
		}

		if !isExpectedString {
			t.Errorf(
				"expected %v.String() to equal one of the strings %v, got %s",
				set,
				expectedStrings,
				setString,
			)
		}
	})
}

func TestStringEmptySet(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		expected := setName + "{}"
		actual := set.String()
		if expected != actual {
			t.Errorf("expected %v.String() == %s, got %s", set, expected, actual)
		}
	})
}

func TestIterator(t *testing.T) {
	testAllSetTypes(func(set set.Set[int], setName string) {
		set.AddMultiple(1, 2, 3)
		results := map[int]bool{}

		set.All()(func(element int) bool {
			results[element] = true
			return true
		})

		if len(results) != 3 {
			t.Errorf("expected iteration result map to have length 3, got %d", len(results))
		}

		if !results[1] || !results[2] || !results[3] {
			t.Errorf(
				"expected iteration result map to contain all elements of original set %v, but got %v",
				set,
				results,
			)
		}
	})
}

func TestDynamicSetTransformation(t *testing.T) {
	var set set.DynamicSet[int]
	if !set.IsArraySet() {
		t.Errorf("expected zero value of %v to start out as ArraySet", set)
	}

	ints := createRandomIntSlice(set.SizeThreshold())
	for _, i := range ints {
		set.Add(i)
	}
	if !set.IsHashSet() {
		t.Errorf(
			"expected %v to transform to HashSet after adding elements up to size threshold",
			set,
		)
	}

	for i := 0; i < set.SizeThreshold()/2; i++ {
		set.Remove(ints[i])
	}
	if !set.IsArraySet() {
		t.Errorf(
			"expected %v to transform back to ArraySet after removing elements down to half of size threshold",
			set,
		)
	}

	set.SetSizeThreshold(set.Size() - 1)
	if !set.IsHashSet() {
		t.Errorf(
			"expected %v to transform back to HashSet after setting size threshold to below current size",
			set,
		)
	}
}

func TestDynamicSetConstructors(t *testing.T) {
	for _, dynamicSet := range []struct {
		set.DynamicSet[int]
		name string
	}{
		{set.DynamicSet[int]{}, "zero value"},
		{set.NewDynamicSet[int](), "NewDynamicSet"},
		{set.DynamicSetWithCapacity[int](5), "DynamicSetWithCapacity"},
		{set.DynamicSetOf(1, 2, 3), "DynamicSetOf"},
		{set.DynamicSetFromSlice([]int{1, 2, 3}), "DynamicSetFromSlice"},
	} {
		if !dynamicSet.IsArraySet() {
			t.Errorf("expected DynamicSet from %s to start out as ArraySet", dynamicSet.name)
		}

		if dynamicSet.SizeThreshold() != set.DefaultDynamicSetSizeThreshold {
			t.Errorf(
				"expected DynamicSet from %s to have size threshold %d, got %d",
				dynamicSet.name,
				set.DefaultDynamicSetSizeThreshold,
				dynamicSet.SizeThreshold(),
			)
		}
	}
}

func testAllSetTypes(testFunc func(set set.Set[int], setName string)) {
	testFunc(&set.ArraySet[int]{}, "ArraySet")
	testFunc(&set.HashSet[int]{}, "HashSet")
	testFunc(&set.DynamicSet[int]{}, "DynamicSet")
}

func assertSize[E comparable, Set set.ComparableSet[E]](t *testing.T, set Set, expectedSize int) {
	t.Helper()

	if actualSize := set.Size(); actualSize != expectedSize {
		t.Errorf("expected %s.Size() == %d, got %d", set.String(), expectedSize, actualSize)
	}
}

func assertContains[E comparable, Set set.ComparableSet[E]](
	t *testing.T,
	set Set,
	expectedElements ...E,
) {
	t.Helper()

	for _, element := range expectedElements {
		if !set.Contains(element) {
			t.Errorf("expected %s to contain %v", set.String(), expectedElements)
			return
		}
	}
}
