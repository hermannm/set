# set

A Go package that provides generic Set data structures (collections of unique elements). It
implements a [`HashSet`](https://pkg.go.dev/hermannm.dev/set#HashSet), an
[`ArraySet`](https://pkg.go.dev/hermannm.dev/set#ArraySet) and a
[`DynamicSet`](https://pkg.go.dev/hermannm.dev/set#DynamicSet), with a common interface between
them.

Run `go get hermannm.dev/set` to add it to your project!

**Docs:** [pkg.go.dev/hermannm.dev/set](https://pkg.go.dev/hermannm.dev/set)

**Contents:**

- [Usage](#usage)
- [Developer's guide](#developers-guide)

## Usage

<!-- @formatter:off -->
```go
import (
	"fmt"

	"hermannm.dev/set"
)

func main() {
	numbers := set.HashSetOf(1, 2, 3)

	fmt.Println(numbers.Contains(3)) // true
	fmt.Println(numbers.Contains(4)) // false

	numbers.Add(4)
	fmt.Println(numbers.Contains(4)) // true

	otherNumbers := set.ArraySetOf(1, 2)
	fmt.Println(otherNumbers.IsSubsetOf(numbers))   // true
	fmt.Println(otherNumbers.IsSupersetOf(numbers)) // false

	overlappingNumbers := set.DynamicSetOf(3, 4, 5)

	union := numbers.Union(overlappingNumbers)
	fmt.Println(union.Size()) // 5

	intersection := numbers.Intersection(overlappingNumbers)
	fmt.Println(intersection.Size()) // 2

	numbers.Clear()
	fmt.Println(numbers.IsEmpty()) // true
}
```
<!-- @formatter:on -->

See the [docs](https://pkg.go.dev/hermannm.dev/set) for more details.

## Developer's guide

When publishing a new release:

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then create a GitHub release with the pushed tag's changelog entry
