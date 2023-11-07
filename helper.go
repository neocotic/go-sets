// Copyright (C) 2023 neocotic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sets

import (
	"github.com/neocotic/go-sets/internal"
	"golang.org/x/exp/constraints"
	"strconv"
	"strings"
)

const (
	// collectionFlagMutable flags an internal.Collection as mutable.
	collectionFlagMutable internal.CollectionFlag = 1 << iota
	// collectionFlagSync flags an internal.Collection as synchronized.
	collectionFlagSync
)

// Asc is a convenient generic less function sorts in ascending order.
func Asc[E constraints.Ordered](x, y E) bool {
	return x < y
}

// Desc is a convenient generic less function sorts in descending order.
func Desc[E constraints.Ordered](x, y E) bool {
	return x > y
}

// Diff returns a new Set struct containing only elements of the Set that do not exist in any other provided Set.
//
// Unlike Set.Diff, the return struct implementation of Set is determined by important characteristics of the Set
// provided. That is; if the Set is mutable, then the returned struct implementation of Set will also be mutable.
// Otherwise, it will be immutable. Likewise for whether any Set is synchronized.
//
// If the Set is nil, Diff returns nil.
func Diff[E comparable](set Set[E], others ...Set[E]) Set[E] {
	return internal.DiffAll[E, Set[E]](createSet[E], flagSet[E], set, asCollections(others))
}

// DiffSymmetric returns a new Set struct containing elements that exist within the Set or any other Set, but not in
// more than one.
//
// Unlike Set.DiffSymmetric, the return struct implementation of Set is determined by important characteristics of the
// Set provided. That is; if the Set is mutable, then the returned struct implementation of Set will also be mutable.
// Otherwise, it will be immutable. Likewise for whether any Set is synchronized.
//
// If the Set is nil, DiffSymmetric returns nil.
func DiffSymmetric[E comparable](set Set[E], others ...Set[E]) Set[E] {
	return internal.DiffSymmetricAll[E, Set[E]](createSet[E], flagSet[E], set, asCollections(others))
}

// Equal is a convenient shorthand for Set.Equal where the Set can be compared against one or more other Set.
//
// If the Set is nil it is treated as having no elements and the same logic applies to the others. To clarify; this
// means that a nil Set is equal to a non-nil Set that contains no elements.
func Equal[E comparable](set Set[E], others ...Set[E]) bool {
	if set == nil {
		var empty *EmptySet[E]
		return equalAll[E](empty, others)
	}
	return equalAll(set, others)
}

// Group returns a map containing the elements within the Set grouped using the grouper function.
//
// The mapped struct implementations of Set are always immutable.
//
// If the Set is nil, Group returns nil.
func Group[E comparable, G comparable](set Set[E], grouper func(element E) G) map[G]Set[E] {
	if internal.IsNil(set) {
		return nil
	}
	groups := make(map[G]Set[E])
	set.Range(func(element E) bool {
		group := grouper(element)
		groups[group] = Singleton(element).Union(groups[group])
		return false
	})
	return groups
}

// Intersection returns a new Set struct containing only elements of the Set that also exist in any other provided Set.
//
// Unlike Set.Intersection, the return struct implementation of Set is determined by important characteristics of the
// Set provided. That is; if the Set is mutable, then the returned struct implementation of Set will also be mutable.
// Otherwise, it will be immutable. Likewise for whether any Set is synchronized.
//
// If the Set is nil, Intersection returns nil.
func Intersection[E comparable](set Set[E], others ...Set[E]) Set[E] {
	return internal.IntersectionAll[E, Set[E]](createSet[E], flagSet[E], set, asCollections(others))
}

// JoinBool is a convenient shorthand for Set.Join where the generic type is a bool, replacing the need for a convert
// function to be provided for casting each element to a string with strconv.FormatBool.
//
// If the Set is nil, JoinBool returns an empty string.
func JoinBool[E ~bool](set Set[E], sep string) string {
	if set == nil {
		return ""
	}
	return set.Join(sep, func(element E) string {
		return strconv.FormatBool(bool(element))
	})
}

// JoinComplex64 is a convenient shorthand for Set.Join where the generic type is a complex64, replacing the need for a
// convert function to be provided for casting each element to a string with strconv.FormatComplex which can be
// controlled by passing options.
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseComplex will return the complex64 element exactly.
//
// If the Set is nil, JoinComplex64 returns an empty string.
func JoinComplex64[E ~complex64](set Set[E], sep string, opts ...JoinComplexOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinComplexOptions(opts)
	return set.Join(sep, getComplexStringConverter[E](64, o))
}

// JoinComplex128 is a convenient shorthand for Set.Join where the generic type is a complex128, replacing the need for
// a convert function to be provided for casting each element to a string with strconv.FormatComplex which can be
// controlled by passing options.
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseComplex will return the complex128 element exactly.
//
// If the Set is nil, JoinComplex128 returns an empty string.
func JoinComplex128[E ~complex128](set Set[E], sep string, opts ...JoinComplexOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinComplexOptions(opts)
	return set.Join(sep, getComplexStringConverter[E](128, o))
}

// JoinFloat32 is a convenient shorthand for Set.Join where the generic type is a float32, replacing the need for a
// convert function to be provided for casting each element to a string with strconv.FormatFloat which can be controlled
// by passing options (excluding sorting options).
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseFloat will return the float32 element exactly.
//
// If the Set is nil, JoinFloat32 returns an empty string.
func JoinFloat32[E ~float32](set Set[E], sep string, opts ...JoinFloatOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinFloatOptions(opts)
	return set.Join(sep, getFloatStringConverter[E](32, o))
}

// JoinFloat64 is a convenient shorthand for Set.Join where the generic type is a float64, replacing the need for a
// convert function to be provided for casting each element to a string with strconv.FormatFloat which can be controlled
// by passing options (excluding sorting options).
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseFloat will return the float64 element exactly.
//
// If the Set is nil, JoinFloat64 returns an empty string.
func JoinFloat64[E ~float64](set Set[E], sep string, opts ...JoinFloatOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinFloatOptions(opts)
	return set.Join(sep, getFloatStringConverter[E](64, o))
}

// JoinInt is a convenient shorthand for Set.Join where the generic type is a signed integer, replacing the need for a
// convert function to be provided for casting each element to a string with strconv.FormatInt which can be controlled
// by passing options (excluding sorting options).
//
// By default, the elements are formatted using base-10.
//
// If the Set is nil, JoinInt returns an empty string.
func JoinInt[E constraints.Signed](set Set[E], sep string, opts ...JoinIntOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinIntOptions(opts)
	return set.Join(sep, getIntStringConverter[E](o))
}

// JoinRune is a convenient shorthand for Set.Join where the generic type is a rune, removing the need for a convert
// function to be provided for casting each element to a string (excluding sorting options).
//
// If the Set is nil, JoinRune returns an empty string.
func JoinRune[E ~rune](set Set[E], sep string) string {
	if set == nil {
		return ""
	}
	var (
		i  int
		sb strings.Builder
	)
	set.Range(func(element E) bool {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteRune(rune(element))
		i++
		return false
	})
	return sb.String()
}

// JoinString is a convenient shorthand for Set.Join where the generic type is a string, removing the need for a convert
// function to be provided for casting each element to a string.
//
// If the Set is nil, JoinString returns an empty string.
func JoinString[E ~string](set Set[E], sep string) string {
	if set == nil {
		return ""
	}
	var (
		i  int
		sb strings.Builder
	)
	set.Range(func(element E) bool {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(string(element))
		i++
		return false
	})
	return sb.String()
}

// JoinUint is a convenient shorthand for Set.Join where the generic type is an unsigned integer, replacing the need for
// a convert function to be provided for casting each element to a string with strconv.FormatUint which can be
// controlled by passing options (excluding sorting options).
//
// By default, the elements are formatted using base-10.
//
// If the Set is nil, JoinUint returns an empty string.
func JoinUint[E constraints.Unsigned](set Set[E], sep string, opts ...JoinUintOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinUintOptions(opts)
	return set.Join(sep, getUintStringConverter[E](o))
}

// Map returns a new Set struct containing values converted from elements within the Set using the mapper function.
//
// The returned struct implementation of Set should match that of the Set being mapped, where possible, but must never
// differ in mutability.
//
// If the Set is nil, Map returns nil.
func Map[E comparable, T comparable](set Set[E], mapper func(element E) T) Set[T] {
	if set == nil {
		return nil
	}
	switch v := set.(type) {
	case *EmptySet[E]:
		var mapped *EmptySet[T]
		if v != nil {
			mapped = &EmptySet[T]{}
		}
		return mapped
	case *HashSet[E]:
		var mapped *HashSet[T]
		if v != nil {
			mapped = &HashSet[T]{internal.Map[E, T](set, mapper)}
		}
		return mapped
	case *MutableHashSet[E]:
		var mapped *MutableHashSet[T]
		if v != nil {
			mapped = &MutableHashSet[T]{internal.Map[E, T](set, mapper)}
		}
		return mapped
	case *SingletonSet[E]:
		var mapped *SingletonSet[T]
		if v != nil {
			mapped = &SingletonSet[T]{mapper(v.element)}
		}
		return mapped
	case *SyncHashSet[E]:
		var mapped *SyncHashSet[T]
		if v != nil {
			mapped = &SyncHashSet[T]{elements: internal.Map[E, T](set, mapper)}
		}
		return mapped
	default:
		if set.IsMutable() {
			var mapped *MutableHashSet[T]
			if internal.IsNotNil(set) {
				mapped = &MutableHashSet[T]{internal.Map[E, T](set, mapper)}
			}
			return mapped
		}
		var mapped *HashSet[T]
		if internal.IsNotNil(set) {
			mapped = &HashSet[T]{internal.Map[E, T](set, mapper)}
		}
		return mapped
	}
}

// Max is a convenient shorthand for Set.Max where the generic type is ordered, removing the need for a less function to
// be provided to control sorting.
//
// If the Set is nil, Max returns the zero value for E and false.
func Max[E constraints.Ordered](set Set[E]) (E, bool) {
	if set == nil {
		var zero E
		return zero, false
	}
	return set.Max(Asc[E])
}

// Min is a convenient shorthand for Set.Min where the generic type is ordered, removing the need for a less function to
// be provided to control sorting.
//
// If the Set is nil, Min returns the zero value for E and false.
func Min[E constraints.Ordered](set Set[E]) (E, bool) {
	if set == nil {
		var zero E
		return zero, false
	}
	return set.Min(Asc[E])
}

// Reduce returns the final result of running the reducer function across all elements within the Set as a single value.
//
// Optionally, an initial value can be specified. Otherwise, the zero value of R is used.
//
// If the Set is nil, Reduce returns initial value or the zero value of R if not specified.
func Reduce[E comparable, R any](set Set[E], reducer func(acc R, element E) R, initValue ...R) R {
	var acc R
	if len(initValue) > 0 {
		acc = initValue[0]
	}
	if set != nil {
		set.Range(func(element E) bool {
			acc = reducer(acc, element)
			return false
		})
	}
	return acc
}

// SortedJoinFloat32 is a convenient shorthand for Set.Join where the generic type is a float32, removing the need for a
// less function to be provided for sorting elements and replacing the need for a convert function to be provided for
// casting each element to a string with strconv.FormatFloat which can be controlled by passing options.
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseFloat will return the float32 element exactly.
//
// If the Set is nil, SortedJoinFloat32 returns an empty string.
func SortedJoinFloat32[E ~float32](set Set[E], sep string, opts ...JoinFloatOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinFloatOptions(opts)
	return set.SortedJoin(sep, getFloatStringConverter[E](32, o), func(x, y E) bool {
		return o.less(float64(x), float64(y))
	})
}

// SortedJoinFloat64 is a convenient shorthand for Set.Join where the generic type is a float64, removing the need for a
// less function to be provided for sorting elements and replacing the need for a convert function to be provided for
// casting each element to a string with strconv.FormatFloat which can be controlled by passing options.
//
// By default, the elements are formatted the 'f' (-ddd.dddd, no exponent) format with the smallest number of digits
// necessary such that strconv.ParseFloat will return the float64 element exactly.
//
// If the Set is nil, SortedJoinFloat64 returns an empty string.
func SortedJoinFloat64[E ~float64](set Set[E], sep string, opts ...JoinFloatOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinFloatOptions(opts)
	return set.SortedJoin(sep, getFloatStringConverter[E](64, o), func(x, y E) bool {
		return o.less(float64(x), float64(y))
	})
}

// SortedJoinInt is a convenient shorthand for Set.Join where the generic type is a signed integer, removing the need
// for a less function to be provided for sorting elements and replacing the need for a convert function to be provided
// for casting each element to a string with strconv.FormatInt which can be controlled by passing options.
//
// By default, the elements are formatted using base-10.
//
// If the Set is nil, SortedJoinInt returns an empty string.
func SortedJoinInt[E constraints.Signed](set Set[E], sep string, opts ...JoinIntOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinIntOptions(opts)
	return set.SortedJoin(sep, getIntStringConverter[E](o), func(x, y E) bool {
		return o.less(int64(x), int64(y))
	})
}

// SortedJoinRune is a convenient shorthand for Set.SortedJoin where the generic type is a rune, removing the need for
// less and convert functions to be provided for sorting elements and then casting them into a string which can be
// controlled by passing options.
//
// If the Set is nil, SortedJoinRune returns an empty string.
func SortedJoinRune[E ~rune](set Set[E], sep string, opts ...SortedJoinRuneOption) string {
	if set == nil {
		return ""
	}
	o := applySortedJoinRuneOptions(opts)
	elements := set.SortedSlice(func(x, y E) bool {
		return o.less(rune(x), rune(y))
	})
	var sb strings.Builder
	for i, element := range elements {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteRune(rune(element))
	}
	return sb.String()
}

// SortedJoinString is a convenient shorthand for Set.SortedJoin where the generic type is a string, removing the need
// for less and convert functions to be provided for sorting elements and then casting them into a string which can be
// controlled by passing options.
//
// If the Set is nil, SortedJoinString returns an empty string.
func SortedJoinString[E ~string](set Set[E], sep string, opts ...SortedJoinStringOption) string {
	if set == nil {
		return ""
	}
	o := applySortedJoinStringOptions(opts)
	elements := set.SortedSlice(func(x, y E) bool {
		return o.less(string(x), string(y))
	})
	var sb strings.Builder
	for i, element := range elements {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(string(element))
	}
	return sb.String()
}

// SortedJoinUint is a convenient shorthand for Set.Join where the generic type is an unsigned integer, removing the
// need for a less function to be provided for sorting elements and replacing the need for a convert function to be
// provided for casting each element to a string with strconv.FormatUint which can be controlled by passing options.
//
// By default, the elements are formatted using base-10.
//
// If the Set is nil, SortedJoinUint returns an empty string.
func SortedJoinUint[E constraints.Unsigned](set Set[E], sep string, opts ...JoinUintOption) string {
	if set == nil {
		return ""
	}
	o := applyJoinUintOptions(opts)
	return set.SortedJoin(sep, getUintStringConverter[E](o), func(x, y E) bool {
		return o.less(uint64(x), uint64(y))
	})
}

// SortedSlice is a convenient shorthand for Set.SortedSlice where the generic type is ordered, removing the need for a
// less function to be provided to control sorting. However, a less function can still be passed optionally for more
// granular control over sorting.
//
// If the Set is nil, SortedSlice returns nil.
func SortedSlice[E constraints.Ordered](set Set[E], less ...func(x, y E) bool) []E {
	if set == nil {
		return nil
	}
	_less := unwrapLess(less)
	return set.SortedSlice(_less)
}

// TryMap returns a new Set struct containing values converted from elements within the Set using the mapper function,
// which may return an error should an element fail to be mapped.
//
// The returned struct implementation of Set should match that of the Set being mapped, where possible, but must never
// differ in mutability.
//
// If the Set is nil, TryMap returns nil.
func TryMap[E comparable, T comparable](set Set[E], mapper func(element E) (T, error)) (Set[T], error) {
	if set == nil {
		return nil, nil
	}
	switch v := set.(type) {
	case *EmptySet[E]:
		var mapped *EmptySet[T]
		if v != nil {
			mapped = &EmptySet[T]{}
		}
		return mapped, nil
	case *HashSet[E]:
		var mapped *HashSet[T]
		if v == nil {
			return mapped, nil
		} else if elements, err := internal.TryMap[E, T](set, mapper); err != nil {
			return mapped, err
		} else {
			mapped = &HashSet[T]{elements}
			return mapped, nil
		}
	case *MutableHashSet[E]:
		var mapped *MutableHashSet[T]
		if v == nil {
			return mapped, nil
		} else if elements, err := internal.TryMap[E, T](set, mapper); err != nil {
			return mapped, err
		} else {
			mapped = &MutableHashSet[T]{elements}
			return mapped, nil
		}
	case *SingletonSet[E]:
		var mapped *SingletonSet[T]
		if v == nil {
			return mapped, nil
		} else if element, err := mapper(v.element); err != nil {
			return mapped, err
		} else {
			mapped = &SingletonSet[T]{element}
			return mapped, nil
		}
	case *SyncHashSet[E]:
		var mapped *SyncHashSet[T]
		if v == nil {
			return mapped, nil
		} else if elements, err := internal.TryMap[E, T](set, mapper); err != nil {
			return mapped, err
		} else {
			mapped = &SyncHashSet[T]{elements: elements}
			return mapped, nil
		}
	default:
		if set.IsMutable() {
			var mapped *MutableHashSet[T]
			if internal.IsNil(set) {
				return mapped, nil
			} else if elements, err := internal.TryMap[E, T](set, mapper); err != nil {
				return mapped, err
			} else {
				mapped = &MutableHashSet[T]{elements}
				return mapped, nil
			}
		}
		var mapped *HashSet[T]
		if internal.IsNil(set) {
			return mapped, nil
		} else if elements, err := internal.TryMap[E, T](set, mapper); err != nil {
			return mapped, err
		} else {
			mapped = &HashSet[T]{elements}
			return mapped, nil
		}
	}
}

// TryReduce returns the final result of running the reducer function across all elements within the Set as a single
// value, which may return an error should an element fail to be reduced.
//
// Optionally, an initial value can be specified. Otherwise, the zero value of T is used.
//
// If the Set is nil, TryReduce returns initial value or the zero value of T if not specified.
func TryReduce[E comparable, T any](set Set[E], reducer func(acc T, element E) (T, error), initValue ...T) (T, error) {
	var acc T
	if len(initValue) > 0 {
		acc = initValue[0]
	}
	var err error
	if set != nil {
		set.Range(func(element E) bool {
			acc, err = reducer(acc, element)
			return err != nil
		})
	}
	return acc, err
}

// Union returns a new Set containing a union of each Set.
//
// Unlike Set.Union, the return struct implementation of Set is determined by important characteristics of each Set
// provided. That is; if any Set is mutable, then the returned struct implementation of Set will also be mutable.
// Otherwise, it will be immutable. Likewise for whether any Set is synchronized.
//
// If each given Set is nil, Union returns nil.
func Union[E comparable](set Set[E], others ...Set[E]) Set[E] {
	return internal.UnionAll[E, Set[E]](createSet[E], flagSet[E], set, asCollections(others))
}

type (
	// JoinComplexOption allows control over the conversion of complex64/complex128 elements into strings when calling
	// JoinComplex64 or JoinComplex128 respectively.
	JoinComplexOption func(opts *joinComplexOptions)

	// joinComplexOptions contains information used to control over the conversion of complex64/complex128 elements into
	// strings using strconv.FormatComplex as well as how complex64/complex128 elements are sorted.
	joinComplexOptions struct {
		format    byte
		precision int
	}
)

// WithComplexFormat controls the format in which the complex64/complex128 element is formatted into a string.
//
// By default, the 'f' (-ddd.dddd, no exponent) format is used.
func WithComplexFormat(format byte) JoinComplexOption {
	return func(opts *joinComplexOptions) {
		opts.format = format
	}
}

// WithComplexPrecision controls the precision to which the complex64/complex128 element is formatted into a string.
//
// By default, the smallest number of digits necessary such that strconv.ParseComplex will return the
// complex64/complex128 element exactly.
func WithComplexPrecision(precision int) JoinComplexOption {
	return func(opts *joinComplexOptions) {
		opts.precision = precision
	}
}

type (
	// JoinFloatOption allows control over the conversion of float32/float64 elements into strings when calling
	// JoinFloat64 or SortedJoinFloat64 or the 32-bit equivalents. Sorting is also controllable for the latter
	// functions.
	JoinFloatOption func(opts *joinFloatOptions)

	// joinFloatOptions contains information used to control over the conversion of float32/float64 elements into
	// strings using strconv.FormatFloat as well as how float32/float64 elements are sorted.
	joinFloatOptions struct {
		format    byte
		less      func(x, y float64) bool
		precision int
	}
)

// WithFloatFormat controls the format in which the float32/float64 element is formatted into a string.
//
// By default, the 'f' (-ddd.dddd, no exponent) format is used.
func WithFloatFormat(format byte) JoinFloatOption {
	return func(opts *joinFloatOptions) {
		opts.format = format
	}
}

// WithFloatPrecision controls the precision to which the float32/float64 element is formatted into a string.
//
// By default, the smallest number of digits necessary such that strconv.ParseFloat will return the float32/float64
// element exactly.
func WithFloatPrecision(precision int) JoinFloatOption {
	return func(opts *joinFloatOptions) {
		opts.precision = precision
	}
}

// WithFloatSorting controls the sorting of float32/float64 elements.
//
// By default, float32/float64 elements are sorted in ascending order.
func WithFloatSorting(less func(x, y float64) bool) JoinFloatOption {
	return func(opts *joinFloatOptions) {
		opts.less = less
	}
}

// WithFloatSortingAsc controls the sorting of float32/float64 elements to use ascending ordering.
//
// This is the default ordering for float32/float64 elements.
func WithFloatSortingAsc() JoinFloatOption {
	return func(opts *joinFloatOptions) {
		opts.less = Asc[float64]
	}
}

// WithFloatSortingDesc controls the sorting of float32/float64 elements to use descending ordering.
//
// By default, float32/float64 elements are sorted in ascending order.
func WithFloatSortingDesc() JoinFloatOption {
	return func(opts *joinFloatOptions) {
		opts.less = Desc[float64]
	}
}

type (
	// JoinIntOption allows control over the conversion of signed integer elements into strings when calling JoinInt or
	// SortedJoinInt. Sorting is also controllable for the latter functions.
	JoinIntOption func(opts *joinIntOptions)

	// joinIntOptions contains information used to control the conversion of signed integer elements into strings using
	// strconv.FormatInt as well as how signed integer elements are sorted.
	joinIntOptions struct {
		base int
		less func(x, y int64) bool
	}
)

// WithIntBase controls the base in which the signed integer element is formatted into a string.
//
// By default, base-10 is used.
func WithIntBase(base int) JoinIntOption {
	return func(opts *joinIntOptions) {
		opts.base = base
	}
}

// WithIntSorting controls the sorting of signed integer elements.
//
// By default, signed integer elements are sorted in ascending order.
func WithIntSorting(less func(x, y int64) bool) JoinIntOption {
	return func(opts *joinIntOptions) {
		opts.less = less
	}
}

// WithIntSortingAsc controls the sorting of signed integer elements to use ascending ordering.
//
// This is the default ordering for signed integer elements.
func WithIntSortingAsc() JoinIntOption {
	return func(opts *joinIntOptions) {
		opts.less = Asc[int64]
	}
}

// WithIntSortingDesc controls the sorting of signed integer elements to use descending ordering.
//
// By default, signed integer elements are sorted in ascending order.
func WithIntSortingDesc() JoinIntOption {
	return func(opts *joinIntOptions) {
		opts.less = Desc[int64]
	}
}

type (
	// JoinUintOption allows control over the conversion of unsigned integer elements into strings when calling
	// JoinUint or SortedJoinUint. Sorting is also controllable for the latter functions.
	JoinUintOption func(opts *joinUintOptions)

	// joinUintOptions contains information used to control over the conversion of unsigned integer elements into
	// strings using strconv.FormatUint as well as how unsigned integer elements are sorted.
	joinUintOptions struct {
		base int
		less func(x, y uint64) bool
	}
)

// WithUintBase controls the base in which the unsigned integer element is formatted into a string.
//
// By default, base-10 is used.
func WithUintBase(base int) JoinUintOption {
	return func(opts *joinUintOptions) {
		opts.base = base
	}
}

// WithUintSorting controls the sorting of unsigned integer elements.
//
// By default, unsigned integer elements are sorted in ascending order.
func WithUintSorting(less func(x, y uint64) bool) JoinUintOption {
	return func(opts *joinUintOptions) {
		opts.less = less
	}
}

// WithUintSortingAsc controls the sorting of unsigned integer elements to use ascending ordering.
//
// This is the default ordering for unsigned integer elements.
func WithUintSortingAsc() JoinUintOption {
	return func(opts *joinUintOptions) {
		opts.less = Asc[uint64]
	}
}

// WithUintSortingDesc controls the sorting of unsigned integer elements to use descending ordering.
//
// By default, unsigned integer elements are sorted in ascending order.
func WithUintSortingDesc() JoinUintOption {
	return func(opts *joinUintOptions) {
		opts.less = Desc[uint64]
	}
}

type (
	// SortedJoinRuneOption allows control over the sorting of rune elements when calling SortedJoinRune.
	SortedJoinRuneOption func(opts *sortedJoinRuneOptions)

	// sortedJoinRuneOptions contains information used to control the sorting of rune elements when calling
	// SortedJoinByte.
	sortedJoinRuneOptions struct {
		less func(x, y rune) bool
	}
)

// WithRuneSorting controls the sorting of rune elements.
//
// By default, rune elements are sorted in ascending order.
func WithRuneSorting(less func(x, y rune) bool) SortedJoinRuneOption {
	return func(opts *sortedJoinRuneOptions) {
		opts.less = less
	}
}

// WithRuneSortingAsc controls the sorting of rune elements to use ascending ordering.
//
// This is the default ordering for rune elements.
func WithRuneSortingAsc() SortedJoinRuneOption {
	return func(opts *sortedJoinRuneOptions) {
		opts.less = Asc[rune]
	}
}

// WithRuneSortingDesc controls the sorting of rune elements to use descending ordering.
//
// By default, rune elements are sorted in ascending order.
func WithRuneSortingDesc() SortedJoinRuneOption {
	return func(opts *sortedJoinRuneOptions) {
		opts.less = Desc[rune]
	}
}

type (
	// SortedJoinStringOption allows control over the sorting of string elements when calling SortedJoinString.
	SortedJoinStringOption func(opts *sortedJoinStringOptions)

	// sortedJoinStringOptions contains information used to control the sorting of string elements when calling
	// SortedJoinString.
	sortedJoinStringOptions struct {
		less func(x, y string) bool
	}
)

// WithStringSorting controls the sorting of string elements.
//
// By default, string elements are sorted in ascending order.
func WithStringSorting(less func(x, y string) bool) SortedJoinStringOption {
	return func(opts *sortedJoinStringOptions) {
		opts.less = less
	}
}

// WithStringSortingAsc controls the sorting of string elements to use ascending ordering.
//
// This is the default ordering for string elements.
func WithStringSortingAsc() SortedJoinStringOption {
	return func(opts *sortedJoinStringOptions) {
		opts.less = Asc[string]
	}
}

// WithStringSortingDesc controls the sorting of string elements to use descending ordering.
//
// By default, string elements are sorted in ascending order.
func WithStringSortingDesc() SortedJoinStringOption {
	return func(opts *sortedJoinStringOptions) {
		opts.less = Desc[string]
	}
}

// applyJoinComplexOptions returns a new joinComplexOptions struct with the given options applied over their defaults.
func applyJoinComplexOptions(opts []JoinComplexOption) *joinComplexOptions {
	o := &joinComplexOptions{
		format:    'f',
		precision: -1,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// applyJoinFloatOptions returns a new joinFloatOptions struct with the given options applied over their defaults.
func applyJoinFloatOptions(opts []JoinFloatOption) *joinFloatOptions {
	o := &joinFloatOptions{
		format:    'f',
		less:      Asc[float64],
		precision: -1,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// applyJoinIntOptions returns a new joinIntOptions struct with the given options applied over their defaults.
func applyJoinIntOptions(opts []JoinIntOption) *joinIntOptions {
	o := &joinIntOptions{
		base: 10,
		less: Asc[int64],
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// applyJoinUintOptions returns a new joinUintOptions struct with the given options applied over their defaults.
func applyJoinUintOptions(opts []JoinUintOption) *joinUintOptions {
	o := &joinUintOptions{
		base: 10,
		less: Asc[uint64],
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// applySortedJoinRuneOptions returns a new sortedJoinRuneOptions struct with the given options applied over their
// defaults.
func applySortedJoinRuneOptions(opts []SortedJoinRuneOption) *sortedJoinRuneOptions {
	o := &sortedJoinRuneOptions{
		less: Asc[rune],
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// applySortedJoinStringOptions returns a new sortedJoinStringOptions struct with the given options applied over their
// defaults.
func applySortedJoinStringOptions(opts []SortedJoinStringOption) *sortedJoinStringOptions {
	o := &sortedJoinStringOptions{
		less: Asc[string],
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// asCollections returns a clone of the given slice of Set interfaces as a slice of internal.Collection interfaces.
func asCollections[E comparable](sets []Set[E]) []internal.Collection[E] {
	cols := make([]internal.Collection[E], len(sets))
	for i, set := range sets {
		cols[i] = set
	}
	return cols
}

// createSet returns a new Set struct for the given internal.Hash based on the flags provided.
//
// If hash is nil, createSet returns a nil reference to an EmptySet.
func createSet[E comparable](hash internal.Hash[E], flags internal.CollectionFlag) Set[E] {
	if hash == nil {
		var ns *EmptySet[E]
		return ns
	} else if flags&collectionFlagSync != 0 {
		return &SyncHashSet[E]{elements: hash}
	} else if flags&collectionFlagMutable != 0 {
		return &MutableHashSet[E]{hash}
	}
	return &HashSet[E]{hash}
}

// equalAll is a convenient shorthand for calling Set.Equal on multiple others.
func equalAll[E comparable](set Set[E], others []Set[E]) bool {
	for _, other := range others {
		if !set.Equal(other) {
			return false
		}
	}
	return true
}

// flagSet returns characteristic flags for the given internal.Collection.
func flagSet[E comparable](col internal.Collection[E]) internal.CollectionFlag {
	if _, ok := col.(*SyncHashSet[E]); ok {
		return collectionFlagMutable | collectionFlagSync
	}
	if col.(Set[E]).IsMutable() {
		return collectionFlagMutable
	}
	return 0
}

// getComplexStringConverter returns a function that can be used to convert a complex64/complex128 element into a string
// using strconv.FormatComplex while allowing options to be passed to control the formatting.
//
// By default, the element will be formatted using the 'f' (-ddd.dddd, no exponent) format and with the smallest number
// of digits necessary such that strconv.ParseComplex will return the complex64/complex128 element exactly.
func getComplexStringConverter[E constraints.Complex](bitSize int, opts *joinComplexOptions) func(element E) string {
	return func(element E) string {
		return strconv.FormatComplex(complex128(element), opts.format, opts.precision, bitSize)
	}
}

// getFloatStringConverter returns a function that can be used to convert a float32/float64 element into a string using
// strconv.FormatFloat while allowing options to be passed to control the formatting.
//
// By default, the element will be formatted using the 'f' (-ddd.dddd, no exponent) format and with the smallest number
// of digits necessary such that strconv.ParseFloat will return the float32/float64 element exactly.
func getFloatStringConverter[E constraints.Float](bitSize int, opts *joinFloatOptions) func(element E) string {
	return func(element E) string {
		return strconv.FormatFloat(float64(element), opts.format, opts.precision, bitSize)
	}
}

// getIntStringConverter returns a function that can be used to convert a signed integer element into a string using
// strconv.FormatInt while allowing options to be passed to control the formatting.
//
// By default, the element will be formatted using base-10.
func getIntStringConverter[E constraints.Signed](opts *joinIntOptions) func(element E) string {
	return func(element E) string {
		return strconv.FormatInt(int64(element), opts.base)
	}
}

// getUintStringConverter returns a function that can be used to convert an unsigned integer element into a string using
// strconv.FormatUint while allowing options to be passed to control the formatting.
//
// By default, the element will be formatted using base-10.
func getUintStringConverter[E constraints.Unsigned](opts *joinUintOptions) func(element E) string {
	return func(element E) string {
		return strconv.FormatUint(uint64(element), opts.base)
	}
}

// unwrapLess is a convenient function for unwrapping an optional less function while supporting the accepted default of
// ascending order.
func unwrapLess[E constraints.Ordered](less []func(x, y E) bool) func(x, y E) bool {
	if len(less) > 0 {
		return less[0]
	}
	return Asc[E]
}
