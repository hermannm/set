# set

A small Go package that provides a generic Set data structure, an unordered collection of unique
elements.

Run `go get hermannm.dev/set` to add it to your project!

The package documentation can be read in the source code, or at
[pkg.go.dev/hermannm.dev/set](https://pkg.go.dev/hermannm.dev/set).

## Usage

```go
import (
	"fmt"

	"hermannm.dev/set"
)

func main() {
	numbers := set.Of(1, 2, 3)

	fmt.Println(numbers.Contains(3)) // true
	fmt.Println(numbers.Contains(4)) // false

	numbers.Add(4)
	fmt.Println(numbers.Contains(4)) // true

	otherNumbers := set.Of(1, 2)
	fmt.Println(otherNumbers.IsSubsetOf(numbers))   // true
	fmt.Println(otherNumbers.IsSupersetOf(numbers)) // false

	overlappingNumbers := set.Of(3, 4, 5)

	union := set.Union(numbers, overlappingNumbers)
	fmt.Println(union.Size()) // 5

	intersection := set.Intersection(numbers, overlappingNumbers)
	fmt.Println(intersection.Size()) // 2

	numbers.Clear()
	fmt.Println(numbers.IsEmpty()) // true
}
```

Refer to the [documentation](https://pkg.go.dev/hermannm.dev/set) for more details.
