package set_test

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

	"hermannm.dev/set"
)

const (
	setSize         = 20
	inputSize       = 100
	maxStringLength = 32
)

var (
	setInts     = createRandomIntSlice(setSize)
	intArraySet = set.ArraySetFromSlice(setInts)
	intHashSet  = set.HashSetFromSlice(setInts)
	inputInts   = createRandomIntSlice(inputSize)

	setStrings     = createRandomStringSlice(setSize)
	stringArraySet = set.ArraySetFromSlice(setStrings)
	stringHashSet  = set.HashSetFromSlice(setStrings)
	inputStrings   = createRandomStringSlice(inputSize)

	setStructs     = createRandomStructSlice(setSize)
	structArraySet = set.ArraySetFromSlice(setStructs)
	structHashSet  = set.HashSetFromSlice(setStructs)
	inputStructs   = createRandomStructSlice(inputSize)
)

// Global variables to avoid the compiler optimizing away our benchmarked function calls (see
// https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)
//
//goland:noinspection GoUnusedGlobalVariable
var globalContains = false

func BenchmarkIntArraySet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, i := range inputInts {
			contains = intArraySet.Contains(i)
		}
	}
	globalContains = contains
}

func BenchmarkIntHashSet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, i := range inputInts {
			contains = intHashSet.Contains(i)
		}
	}
	globalContains = contains
}

func BenchmarkStringArraySet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, s := range inputStrings {
			contains = stringArraySet.Contains(s)
		}
	}
	globalContains = contains
}

func BenchmarkStringHashSet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, s := range inputStrings {
			contains = stringHashSet.Contains(s)
		}
	}
	globalContains = contains
}

func BenchmarkStructArraySet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, s := range inputStructs {
			contains = structArraySet.Contains(s)
		}
	}
	globalContains = contains
}

func BenchmarkStructHashSet(b *testing.B) {
	var contains bool
	for range b.N {
		for _, s := range inputStructs {
			contains = structHashSet.Contains(s)
		}
	}
	globalContains = contains
}

func createRandomIntSlice(length int) []int {
	ints := make([]int, length*2)

	for i := range ints {
		ints[i] = i
	}

	for i := range ints {
		j := rand.Intn(i + 1)
		ints[i], ints[j] = ints[j], ints[i]
	}

	return ints[:length]
}

func createRandomStringSlice(length int) []string {
	strings := make([]string, length)

	for i := range strings {
		stringLength := rand.Intn(maxStringLength) + 1
		bytes := make([]byte, stringLength)
		if _, err := cryptorand.Read(bytes); err != nil {
			panic(fmt.Errorf("failed to create random string: %w", err))
		}

		strings[i] = string(bytes)
	}

	return strings
}

type testStruct struct {
	i int
	s string
}

func createRandomStructSlice(length int) []testStruct {
	structs := make([]testStruct, length)

	ints := createRandomIntSlice(length)
	strings := createRandomStringSlice(length)
	for i := range structs {
		structs[i] = testStruct{i: ints[i], s: strings[i]}
	}

	return structs
}
