package set

import (
	"fmt"
	"strings"
	"testing"
)

// Benchmark result: 535.3 ns/op
func (set Set[Element]) stringWithSuffixRemove() string {
	if set.Size() == 0 {
		return "Set{}"
	}

	setString := "Set{"
	for element := range set {
		setString += fmt.Sprintf("%v, ", element)
	}
	setString = setString[:len(setString)-2]
	setString += "}"

	return setString
}

// Benchmark result: 518.7 ns/op
func (set Set[Element]) stringWithIndexCheck() string {
	size := set.Size()

	setString := "Set{"

	index := 0
	for element := range set {
		var format string
		if index < size-1 {
			format = "%v, "
		} else {
			format = "%v"
		}

		setString += fmt.Sprintf(format, element)

		index++
	}

	setString += "}"

	return setString
}

// Benchmark result: 415.6 ns/op
func (set Set[Element]) stringWithBuilder() string {
	size := set.Size()

	var stringBuilder strings.Builder

	stringBuilder.WriteString("Set{")

	index := 0
	for element := range set {
		var format string
		if index < size-1 {
			format = "%v, "
		} else {
			format = "%v"
		}

		fmt.Fprintf(&stringBuilder, format, element)

		index++
	}

	stringBuilder.WriteString("}")

	return stringBuilder.String()
}

func BenchmarkStringWithSuffixRemove(b *testing.B) {
	numbers := Of(1, 2, 3, 4, 5)

	for n := 0; n < b.N; n++ {
		numbers.stringWithSuffixRemove()
	}
}

func BenchmarkStringWithIndexCheck(b *testing.B) {
	numbers := Of(1, 2, 3, 4, 5)

	for n := 0; n < b.N; n++ {
		numbers.stringWithIndexCheck()
	}
}

func BenchmarkStringWithBuilder(b *testing.B) {
	numbers := Of(1, 2, 3, 4, 5)

	for n := 0; n < b.N; n++ {
		numbers.stringWithBuilder()
	}
}
