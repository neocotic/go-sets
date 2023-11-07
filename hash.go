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
	"encoding/json"
	"fmt"
	"github.com/neocotic/go-sets/internal"
)

// HashSet is an immutable implementation of Set that contains a unique data set.
//
// As HashSet is immutable it is safe for concurrent use by multiple goroutines without additional locking or
// coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use. That said; HashSet only implements json.Unmarshaler for the purpose of being able to have
// a HashSet field value on a struct being unmarshalled. It's recommended to unmarshal JSON into a HashSet using
// HashFromJSON as JSON is typically only unmarshalled into a struct once.
type HashSet[E comparable] struct {
	elements internal.Hash[E]
}

var (
	_ Set[any]         = (*HashSet[any])(nil)
	_ fmt.Stringer     = (*HashSet[any])(nil)
	_ json.Marshaler   = (*HashSet[any])(nil)
	_ json.Unmarshaler = (*HashSet[any])(nil)
)

// Clone returns a clone of the HashSet.
//
// If the HashSet is nil, HashSet.Clone returns nil.
func (s *HashSet[E]) Clone() Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.Clone[E](s.elements)}
}

// Contains returns whether the HashSet contains the element.
//
// If the HashSet is nil, HashSet.Contains returns false.
func (s *HashSet[E]) Contains(element E) bool {
	if s == nil {
		return false
	}
	_, ok := s.elements[element]
	return ok
}

// Diff returns a new HashSet struct containing only elements of the HashSet that do not exist in another Set.
//
// If the HashSet is nil, HashSet.Diff returns nil.
func (s *HashSet[E]) Diff(other Set[E]) Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.Diff[E](s.elements, other)}
}

// DiffSymmetric returns a new HashSet struct containing elements that exist within the HashSet or another Set, but not
// both.
//
// If the HashSet is nil, HashSet.DiffSymmetric returns nil.
func (s *HashSet[E]) DiffSymmetric(other Set[E]) Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.DiffSymmetric[E](s.elements, other)}
}

// Equal returns whether the HashSet contains the exact same elements as another Set.
//
// If the HashSet is nil it is treated as having no elements and the same logic applies to the other Set. To clarify;
// this means that a nil Set is equal to a non-nil Set that contains no elements.
func (s *HashSet[E]) Equal(other Set[E]) bool {
	if s == nil {
		return other == nil || other.IsEmpty()
	} else if other == nil {
		return s.IsEmpty()
	}
	return internal.ContainsOnly[E](s.elements, other.Slice())
}

// Every returns whether the HashSet contains elements that all match the predicate function.
//
// If the HashSet is nil, HashSet.Every returns false.
func (s *HashSet[E]) Every(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	return internal.Every[E](s.elements, predicate)
}

// Filter returns a new HashSet struct containing only elements of the HashSet that match the filter function.
//
// If the HashSet is nil, HashSet.Filter returns nil.
func (s *HashSet[E]) Filter(filter func(element E) bool) Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.Filter[E](s.elements, filter)}
}

// Find returns an element within the HashSet that matches the search function as well as an indication of whether a
// match was found.
//
// Iteration order is not guaranteed to be consistent so results may vary.
//
// If the HashSet is nil, HashSet.Find returns the zero value for E and false.
func (s *HashSet[E]) Find(search func(element E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Find[E](s.elements, search)
}

// Immutable returns a reference to itself to conform with Set.Immutable.
//
// If the HashSet is nil, HashSet.Immutable returns nil.
func (s *HashSet[E]) Immutable() Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return s
}

// Intersection returns a new HashSet struct containing only elements of the HashSet that also exist in another Set.
//
// If the HashSet is nil, HashSet.Intersection returns nil.
func (s *HashSet[E]) Intersection(other Set[E]) Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.Intersection[E](s.elements, other)}
}

// IsEmpty returns whether the HashSet contains no elements.
//
// If the HashSet is nil, HashSet.IsEmpty returns true.
func (s *HashSet[E]) IsEmpty() bool {
	if s == nil {
		return true
	}
	return len(s.elements) == 0
}

// IsMutable always returns false to conform with Set.IsMutable.
func (s *HashSet[E]) IsMutable() bool {
	return false
}

// Join converts the elements within the HashSet to strings which are then concatenated to create a single string,
// placing sep between the converted elements in the resulting string.
//
// The order of elements within the resulting string is not guaranteed to be consistent. HashSet.SortedJoin should be
// used instead for such cases where consistent ordering is required.
//
// If the HashSet is nil, HashSet.Join returns an empty string.
func (s *HashSet[E]) Join(sep string, convert func(element E) string) string {
	if s == nil {
		return ""
	}
	return internal.Join[E](s.elements, sep, convert)
}

// Len returns the number of elements within the HashSet.
//
// If the HashSet is nil, HashSet.Len returns zero.
func (s *HashSet[E]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

// Max returns the maximum element within the HashSet using the provided less function.
//
// If the HashSet is nil, HashSet.Max returns the zero value for E and false.
func (s *HashSet[E]) Max(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Max[E](s.elements, less)
}

// Min returns the minimum element within the HashSet using the provided less function.
//
// If the HashSet is nil, HashSet.Min returns the zero value for E and false.
func (s *HashSet[E]) Min(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Min[E](s.elements, less)
}

// Mutable returns a mutable clone of the HashSet.
//
// If the HashSet is nil, HashSet.Mutable returns nil.
func (s *HashSet[E]) Mutable() MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.Clone[E](s.elements)}
}

// None returns whether the HashSet contains no elements that match the predicate function.
//
// If the HashSet is nil, HashSet.None returns true.
func (s *HashSet[E]) None(predicate func(element E) bool) bool {
	if s == nil {
		return true
	}
	return internal.None[E](s.elements, predicate)
}

// Range calls the iter function with each element within the HashSet but will stop early whenever the iter function
// returns true.
//
// Iteration order is not guaranteed to be consistent.
//
// If the HashSet is nil, HashSet.Range is a no-op.
func (s *HashSet[E]) Range(iter func(element E) bool) {
	if s != nil {
		internal.Range[E](s.elements, iter)
	}
}

// Slice returns a slice containing all elements of the HashSet.
//
// The order of elements within the resulting slice is not guaranteed to be consistent. HashSet.SortedSlice should be
// used instead for such cases where consistent ordering is required.
//
// If the HashSet is nil, HashSet.Slice returns nil.
func (s *HashSet[E]) Slice() []E {
	if s == nil {
		return nil
	}
	return internal.Slice[E](s.elements)
}

// Some returns whether the HashSet contains any element that matches the predicate function.
//
// If the HashSet is nil, HashSet.Some returns false.
func (s *HashSet[E]) Some(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	return internal.Some[E](s.elements, predicate)
}

// SortedJoin sorts the elements within the HashSet using the provided less function and then converts those elements
// into strings which are then joined using the specified separator to create the resulting string.
//
// If the HashSet is nil, HashSet.SortedJoin returns an empty string.
func (s *HashSet[E]) SortedJoin(sep string, convert func(element E) string, less func(x, y E) bool) string {
	if s == nil {
		return ""
	}
	return internal.SortedJoin[E](s.elements, sep, convert, less)
}

// SortedSlice returns a slice containing all elements of the HashSet sorted using the provided less function.
//
// If the HashSet is nil, HashSet.SortedSlice returns nil.
func (s *HashSet[E]) SortedSlice(less func(x, y E) bool) []E {
	if s == nil {
		return nil
	}
	return internal.SortedSlice[E](s.elements, less)
}

// TryRange calls the iter function with each element within the HashSet but will stop early whenever the iter function
// returns an error.
//
// Iteration order is not guaranteed to be consistent.
//
// If the HashSet is nil, HashSet.TryRange is a no-op.
func (s *HashSet[E]) TryRange(iter func(element E) error) error {
	if s == nil {
		return nil
	}
	return internal.TryRange[E](s.elements, iter)
}

// Union returns a new HashSet containing a union of the HashSet with another Set.
//
// If the HashSet and the other Set are both nil, HashSet.Union returns nil.
func (s *HashSet[E]) Union(other Set[E]) Set[E] {
	if elements := internal.Union[E](s, other); elements != nil {
		return &HashSet[E]{elements}
	}
	var ns *HashSet[E]
	return ns
}

func (s *HashSet[E]) String() string {
	if s == nil {
		return internal.NilString
	}
	return internal.String[E](s.elements)
}

func (s *HashSet[E]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return internal.MarshalJSONNil()
	}
	return internal.MarshalJSON[E](s.elements)
}

func (s *HashSet[E]) UnmarshalJSON(data []byte) error {
	if elements, err := internal.UnmarshalJSON[E](data); err != nil {
		return err
	} else {
		s.elements = elements
		return nil
	}
}

// Hash returns an immutable HashSet struct that implements Set containing each unique element provided.
//
// As Hash returns an immutable struct it is safe for concurrent use by multiple goroutines without additional locking
// or coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use.
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use. That said; HashSet only implements json.Unmarshaler for the purpose of being able to have
// a HashSet field value on a struct being unmarshalled. It's recommended to unmarshal JSON into a HashSet using
// HashFromJSON as JSON is typically only unmarshalled into a struct once.
func Hash[E comparable](elements ...E) *HashSet[E] {
	return &HashSet[E]{internal.FromSlice[E](elements)}
}

// HashFromJSON returns an immutable HashSet struct that implements Set containing each unique element parsed from the
// JSON-encoded data provided.
//
// As HashFromJSON returns an immutable struct it is safe for concurrent use by multiple goroutines without additional
// locking or coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use and, as JSON is typically only unmarshalled into a struct once, it's unlikely that this
// needs to be called on the returned HashSet again after calling this function.
func HashFromJSON[E comparable](data []byte) (*HashSet[E], error) {
	set := &HashSet[E]{}
	if err := json.Unmarshal(data, set); err != nil {
		return nil, err
	}
	return set, nil
}

// HashFromSlice returns an immutable HashSet struct that implements Set containing each unique element from the slice
// provided.
//
// As HashFromSlice returns an immutable struct it is safe for concurrent use by multiple goroutines without additional
// locking or coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use.
func HashFromSlice[E comparable](elements []E) *HashSet[E] {
	return &HashSet[E]{internal.FromSlice[E](elements)}
}
