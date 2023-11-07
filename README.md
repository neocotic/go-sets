# go-sets

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/neocotic/go-sets)
[![Build Status](https://img.shields.io/github/actions/workflow/status/neocotic/go-sets/ci.yml?style=for-the-badge)](https://github.com/neocotic/go-sets/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/neocotic/go-sets?style=for-the-badge)](https://github.com/neocotic/go-sets)
[![License](https://img.shields.io/github/license/neocotic/go-sets?style=for-the-badge)](https://github.com/neocotic/go-sets/blob/main/LICENSE.md)

Easy-to-use generic set collections for Go (golang).

Provides separate implementations for mutable and immutable sets whilst making it easy to create clones of varying
mutability. Immutable sets play well with concurrency out-of-the-box, however, a special implementation of a mutable set
is available for concurrent use without requiring additional locking or coordination.

| Set           | Elements | Mutable | Concurrency Safe |
|---------------|----------|---------|------------------|
| `Empty`       | 0        | No      | Yes              |
| `Hash`        | Infinite | No      | Yes              |
| `MutableHash` | Infinite | Yes     | No               |
| `Singleton`   | 1        | No      | Yes              |
| `SyncHash`    | Infinite | Yes     | Yes              |

## Installation

Install using [go install](https://go.dev/ref/mod#go-install):

``` sh
go install github.com/neocotic/go-sets
```

Then import the package into your own code:

``` go
import "github.com/neocotic/go-sets"
```

## Documentation

Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/neocotic/go-sets#section-documentation). It
contains an overview and reference.

### Example

``` go
set := sets.Hash(123, 456, 789, 456, 123)
set.Len() // => 3
set.IsEmpty() // => false
set.Contains(456) // => true
set.Contains(-456) // => false
set.Find(func(e int) bool { return e < 200 }) // => 123, true
set.Min(sets.Asc[int]) // => 123
set.Max(sets.Asc[int]) // => 789
set.Equal(sets.Hash(789, 456, 123)) // => true
set.Range(func(e int) bool { fmt.Printf("%v\n", e); return false })
set.Every(func(e int) bool { return e < 200 }) // => false
set.None(func(e int) bool { return e < 200 }) // => false
set.Some(func(e int) bool { return e < 200 }) // => true

filtered := set.Filter(func(e int) bool { return e < 200 })
filtered.Equal(sets.Singleton(123)) // => true

mset := set.Mutable()
mset.Put(0)
mset.Equal(sets.Hash(0, 123, 456, 789)) // => true

set.Intersection(mset).Equal(sets.Hash(123, 456, 789)) // => true
set.Diff(mset).IsEmpty() // => true
mset.Diff(set).Equal(sets.Singleton(0)) // => true
set.DiffSymmetric(mset).Equal(sets.Singleton(0)) // => true

cset := mset.Clone()
cset.Delete(789)
cset.Equal(sets.Hash(0, 123, 456)) // => true
cset.Clear()
cset.IsEmpty() // => true
mset.Equal(sets.Hash(0, 123, 456, 789)) // => true

mapped := sets.Map[int, string](set, func(e int) string { return strconv.FormatInt(int64(e), 10) })
mapped.Equal(sets.Hash("123", "456", "789")) // => true

sets.Reduce(set, func(acc uint, e int) uint { return acc + uint(e) }) // => 1368

sets.Min(set) // => 123
sets.Max(set) // => 789
sets.SortedSlice(set) // => [123 456 789]
sets.SortedJoinInt(set, ",") // => "123,456,789"
```

There's many more functions available to explore!

## Issues

If you have any problems or would like to see changes currently in development you can do so
[here](https://github.com/neocotic/go-sets/issues).

## Contributors

If you want to contribute, you're a legend! Information on how you can do so can be found in
[CONTRIBUTING.md](https://github.com/neocotic/go-sets/blob/main/CONTRIBUTING.md). We want your suggestions and pull
requests!

A list of contributors can be found in [AUTHORS.md](https://github.com/neocotic/go-sets/blob/main/AUTHORS.md).

## License

Copyright © 2023 neocotic

See [LICENSE.md](https://github.com/neocotic/go-sets/raw/main/LICENSE.md) for more information on our MIT license.
