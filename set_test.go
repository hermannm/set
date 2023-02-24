package set_test

import (
	"testing"

	"hermannm.dev/set"
)

func TestNew(t *testing.T) {
	newSet := set.New[string]()

	if newSet == nil {
		t.Error("expected set.New to give non-nil set")
	}

	if size := newSet.Size(); size != 0 {
		t.Errorf("expected set from set.New to have size 0, got %d", size)
	}
}

func TestWithCapacity(t *testing.T) {
	newSet := set.WithCapacity[string](3)

	if newSet == nil {
		t.Error("expected set.WithCapacity to give non-nil set")
	}

	if size := newSet.Size(); size != 0 {
		t.Errorf("expected set from set.WithCapacity to have size 0, got %d", size)
	}
}

func TestOf(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	if numbers == nil {
		t.Error("expected set.Of to give non-nil set")
	}

	if size := numbers.Size(); size != 3 {
		t.Errorf("expected set from set.Of(1, 2, 3) to have size 3, got %d", size)
	}
}

func TestFromSlice(t *testing.T) {
	numberSlice := []int{1, 2, 3}
	numberSet := set.FromSlice(numberSlice)

	if numberSet == nil {
		t.Error("expected set.FromSlice to give non-nil set")
	}

	setSize := numberSet.Size()
	sliceLength := len(numberSlice)
	if setSize != sliceLength {
		t.Errorf(
			"expected set from set.FromSlice to have size equal to slice length %d, got %d",
			sliceLength,
			setSize,
		)
	}
}

func TestFromSliceWithDuplicates(t *testing.T) {
	numbersWithDuplicates := []int{1, 1, 2, 2}
	numbersWithoutDuplicates := set.FromSlice(numbersWithDuplicates)

	if size := numbersWithoutDuplicates.Size(); size != 2 {
		t.Errorf("expected size 2 from set of 2 unique elements, got %d", size)
	}
}

func TestAdd(t *testing.T) {
	numbers := set.New[int]()

	numbers.Add(1)

	if size := numbers.Size(); size != 1 {
		t.Errorf("expected set size to be 1 after one Add(), got %d", size)
	}

	if !setContainsAll(numbers, 1) {
		t.Errorf("expected Set{1} after Add(1), got %v", numbers)
	}
}

func TestAddDuplicate(t *testing.T) {
	numbers := set.Of(1, 2, 3)
	size1 := numbers.Size()

	numbers.Add(3)
	size2 := numbers.Size()

	if size1 != size2 {
		t.Errorf(
			"expected adding of existing element to not change set size %d, but got %d",
			size1,
			size2,
		)
	}
}

func TestAddToNil(t *testing.T) {
	defer func() {
		err := recover()

		if err == nil {
			t.Error("expected adding to nil set to panic")
		}

		errMessage, ok := err.(string)
		if !ok {
			t.Errorf("expected Add to panic with string, got %v", err)
		}

		if expectedMessage := "called Add on nil Set"; errMessage != expectedMessage {
			t.Errorf(
				`expected Add to panic with message "%s", got "%s"`, expectedMessage, errMessage,
			)
		}
	}()

	var nilSet set.Set[string]
	nilSet.Add("test")
}

func TestRemove(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	numbers.Remove(3)

	if size := numbers.Size(); size != 2 {
		t.Errorf("expected size 2 after removing from 3-element set, got %d", size)
	}

	if !setContainsAll(numbers, 1, 2) {
		t.Errorf("expected Set{1, 2} after removing 3 from Set{1, 2, 3}, got %v", numbers)
	}
}

func TestRemoveNonExisting(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	numbers.Remove(4)

	if size := numbers.Size(); size != 3 {
		t.Errorf("expected unchanged size 3 after removing non-existing element, got %d", size)
	}

	if !setContainsAll(numbers, 1, 2, 3) {
		t.Errorf("expected unchanged Set{1, 2, 3} after removing 4, got %v", numbers)
	}
}

func TestClear(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	numbers.Clear()

	if size := numbers.Size(); size != 0 {
		t.Errorf("expected size 0 after clearing set, got %d", size)
	}
}

func TestSize(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	if size := numbers.Size(); size != 3 {
		t.Errorf("expected %v to have size 3, got %d", numbers, size)
	}
}

func TestIsEmpty(t *testing.T) {
	newSet := set.New[string]()

	if isEmpty := newSet.IsEmpty(); !isEmpty {
		t.Errorf("expected %v.IsEmpty() = true, got %v", newSet, isEmpty)
	}

	numbers := set.Of(1, 2, 3)

	if isEmpty := numbers.IsEmpty(); isEmpty {
		t.Errorf("expected %v.IsEmpty() = false, got %v", numbers, isEmpty)
	}
}

func TestContains(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	if contains := numbers.Contains(3); !contains {
		t.Errorf("expected %v.Contains(3) = true, got %v", numbers, contains)
	}

	if contains := numbers.Contains(4); contains {
		t.Errorf("expected %v.Contains(4) = false, got %v", numbers, contains)
	}
}

func TestEquals(t *testing.T) {
	numbers1 := set.Of(1, 2, 3)
	numbers2 := set.Of(1, 2, 3)

	if equal := numbers1.Equals(numbers2); !equal {
		t.Errorf("expected %v.Equals(%v) = true, got %v", numbers1, numbers2, equal)
	}

	numbers3 := set.Of(1, 2, 4)

	if equal := numbers1.Equals(numbers3); equal {
		t.Errorf("expected %v.Equals(%v) = false, got %v", numbers1, numbers3, equal)
	}
}

func TestIsSubsetOf(t *testing.T) {
	numbers1 := set.Of(1, 2, 3)
	numbers2 := set.Of(1, 2, 3, 4, 5, 6)

	if isSubset := numbers1.IsSubsetOf(numbers2); !isSubset {
		t.Errorf("expected %v.IsSubsetOf(%v) = true, got %v", numbers1, numbers2, isSubset)
	}

	if isSubset := numbers2.IsSubsetOf(numbers1); isSubset {
		t.Errorf("expected %v.IsSubsetOf(%v) = false, got %v", numbers2, numbers1, isSubset)
	}
}

func TestIsSupersetOf(t *testing.T) {
	numbers1 := set.Of("test1", "test2", "test3")
	numbers2 := set.Of("test1", "test2")

	if isSuperset := numbers1.IsSupersetOf(numbers2); !isSuperset {
		t.Errorf("expected %v.IsSupersetOf(%v) = true, got %v", numbers1, numbers2, isSuperset)
	}

	if isSuperset := numbers2.IsSupersetOf(numbers1); isSuperset {
		t.Errorf("expected %v.IsSupersetOf(%v) = false, got %v", numbers2, numbers1, isSuperset)
	}
}

func TestToSlice(t *testing.T) {
	numberSet := set.Of(1, 2, 3)
	numberSlice := numberSet.ToSlice()

	setSize := numberSet.Size()
	sliceLength := len(numberSlice)
	if setSize != sliceLength {
		t.Errorf(
			"expected %v.Size() = len(%v), got set size %d and slice length %d",
			numberSet,
			numberSlice,
			setSize,
			sliceLength,
		)
	}

	for setElement := range numberSet {
		containedInSlice := false

		for _, sliceElement := range numberSlice {
			if setElement == sliceElement {
				containedInSlice = true
				break
			}
		}

		if !containedInSlice {
			t.Errorf(
				"expected %v to contain all elements of %v, but did not contain %v",
				numberSlice,
				numberSet,
				setElement,
			)
		}
	}
}

func TestCopy(t *testing.T) {
	numbers := set.Of(1, 2, 3)
	numbersCopy := numbers.Copy()

	if !setContainsAll(numbersCopy, 1, 2, 3) {
		t.Errorf("expected copy %v to contain all elements of original %v", numbersCopy, numbers)
	}

	numbers.Add(4)

	if size := numbersCopy.Size(); size != 3 {
		t.Errorf("expected unchanged size 3 of copy after adding to original set, got %d", size)
	}
}

func TestString(t *testing.T) {
	numbers := set.Of(1, 2, 3)

	numbersString := numbers.String()
	expectedStrings := []string{
		"Set{1, 2, 3}",
		"Set{1, 3, 2}",
		"Set{2, 1, 3}",
		"Set{2, 3, 1}",
		"Set{3, 1, 2}",
		"Set{3, 2, 1}",
	}

	isExpectedString := false
	for _, expected := range expectedStrings {
		if numbersString == expected {
			isExpectedString = true
		}
	}

	if !isExpectedString {
		t.Errorf(
			"expected %v.String() to equal one of the strings %v, got %s",
			numbers,
			expectedStrings,
			numbersString,
		)
	}
}

func TestStringEmptySet(t *testing.T) {
	emptySet := set.New[string]()

	expected := "Set{}"
	actual := emptySet.String()
	if expected != actual {
		t.Errorf("expected %v.String() = %s, got %s", emptySet, expected, actual)
	}
}

func TestUnion(t *testing.T) {
	numbers1 := set.Of(1, 2, 3)
	numbers2 := set.Of(3, 4, 5)

	union := set.Union(numbers1, numbers2)

	if size := union.Size(); size != 5 {
		t.Errorf("expected union %v.Size() = 5, got %d", union, size)
	}

	if !setContainsAll(union, 1, 2, 3, 4, 5) {
		t.Errorf(
			"expected union %v to contain all elements of component sets %v and %v",
			union,
			numbers1,
			numbers2,
		)
	}
}

func TestIntersection(t *testing.T) {
	numbers1 := set.Of(1, 2, 3, 4)
	numbers2 := set.Of(2, 3, 4, 5)

	intersection := set.Intersection(numbers1, numbers2)

	if size := intersection.Size(); size != 3 {
		t.Errorf("expected intersection %v.Size() = 3, got %d", intersection, size)
	}

	if !setContainsAll(intersection, 2, 3, 4) {
		t.Errorf(
			"expected intersection %v to contain shared elements of component sets %v and %v",
			intersection,
			numbers1,
			numbers2,
		)
	}
}

// setContainsAll checks that the given set contains all of the given elements.
func setContainsAll[Element comparable](set set.Set[Element], elements ...Element) bool {
	for _, element := range elements {
		if !set.Contains(element) {
			return false
		}
	}

	return true
}
