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

// MutableHashSet is an implementation of MutableSet that contains a unique data set.
//
// As MutableHash is mutable it is not safe for concurrent use by multiple goroutines. SyncHashSet should be used
// instead for such cases where mutability is required, otherwise HashSet for a simple immutable Set.
type MutableHashSet[E comparable] struct {
	elements internal.Hash[E]
}

var (
	_ MutableSet[any]  = (*MutableHashSet[any])(nil)
	_ fmt.Stringer     = (*MutableHashSet[any])(nil)
	_ json.Marshaler   = (*MutableHashSet[any])(nil)
	_ json.Unmarshaler = (*MutableHashSet[any])(nil)
)

// Clear removes all elements from the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.Clear is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) Clear() MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	s.elements = make(internal.Hash[E])
	return s
}

// Clone returns a clone of the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.Clone returns nil.
func (s *MutableHashSet[E]) Clone() Set[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.Clone[E](s.elements)}
}

// Contains returns whether the MutableHashSet contains the element.
//
// If the MutableHashSet is nil, MutableHashSet.Contains returns false.
func (s *MutableHashSet[E]) Contains(element E) bool {
	if s == nil {
		return false
	}
	_, ok := s.elements[element]
	return ok
}

// Delete removes the element from the MutableHashSet as well as any additional elements specified.
//
// If the MutableHashSet is nil, MutableHashSet.Delete is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) Delete(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.Delete[E](s.elements, element, elements)
	return s
}

// DeleteAll removes all elements in the specified Set from the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.DeleteAll is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) DeleteAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.DeleteAll[E](s.elements, elements)
	return s
}

// DeleteSlice removes all elements in the specified slice from the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.DeleteSlice is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) DeleteSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.DeleteSlice[E](s.elements, elements)
	return s
}

// DeleteWhere removes all elements that match the predicate function from the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.DeleteWhere is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) DeleteWhere(predicate func(element E) bool) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.DeleteWhere[E](s.elements, predicate)
	return s
}

// Diff returns a new MutableHashSet struct containing only elements of the MutableHashSet that do not exist in another
// Set.
//
// If the MutableHashSet is nil, MutableHashSet.Diff returns nil.
func (s *MutableHashSet[E]) Diff(other Set[E]) Set[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.Diff[E](s.elements, other)}
}

// DiffSymmetric returns a new MutableHashSet struct containing elements that exist within the MutableHashSet or another
// Set, but not both.
//
// If the MutableHashSet is nil, MutableHashSet.DiffSymmetric returns nil.
func (s *MutableHashSet[E]) DiffSymmetric(other Set[E]) Set[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.DiffSymmetric[E](s.elements, other)}
}

// Equal returns whether the MutableHashSet contains the exact same elements as another Set.
//
// If the MutableHashSet is nil it is treated as having no elements and the same logic applies to the other Set. To
// clarify; this means that a nil Set is equal to a non-nil Set that contains no elements.
func (s *MutableHashSet[E]) Equal(other Set[E]) bool {
	if s == nil {
		return other == nil || other.IsEmpty()
	} else if other == nil {
		return s.IsEmpty()
	}
	return internal.ContainsOnly[E](s.elements, other.Slice())
}

// Every returns whether the MutableHashSet contains elements that all match the predicate function.
//
// If the MutableHashSet is nil, MutableHashSet.Every returns false.
func (s *MutableHashSet[E]) Every(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	return internal.Every[E](s.elements, predicate)
}

// Filter returns a new MutableHashSet struct containing only elements of the MutableHashSet that match the filter
// function.
//
// If the MutableHashSet is nil, MutableHashSet.Filter returns nil.
func (s *MutableHashSet[E]) Filter(filter func(element E) bool) Set[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.Filter[E](s.elements, filter)}
}

// Find returns an element within the MutableHashSet that matches the search function as well as an indication of
// whether a match was found.
//
// Iteration order is not guaranteed to be consistent so results may vary.
//
// If the MutableHashSet is nil, MutableHashSet.Find returns the zero value for E and false.
func (s *MutableHashSet[E]) Find(search func(element E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Find[E](s.elements, search)
}

// Immutable returns an immutable clone of the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.Immutable returns nil.
func (s *MutableHashSet[E]) Immutable() Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	return &HashSet[E]{internal.Clone[E](s.elements)}
}

// Intersection returns a new MutableHashSet struct containing only elements of the MutableHashSet that also exist in
// another Set.
//
// If the MutableHashSet is nil, MutableHashSet.Intersection returns nil.
func (s *MutableHashSet[E]) Intersection(other Set[E]) Set[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return &MutableHashSet[E]{internal.Intersection[E](s.elements, other)}
}

// IsEmpty returns whether the MutableHashSet contains no elements.
//
// If the MutableHashSet is nil, MutableHashSet.IsEmpty returns true.
func (s *MutableHashSet[E]) IsEmpty() bool {
	if s == nil {
		return true
	}
	return len(s.elements) == 0
}

// IsMutable always returns true to conform with Set.IsMutable.
func (s *MutableHashSet[E]) IsMutable() bool {
	return true
}

// Join converts the elements within the MutableHashSet to strings which are then concatenated to create a single
// string, placing sep between the converted elements in the resulting string.
//
// The order of elements within the resulting string is not guaranteed to be consistent. MutableHashSet.SortedJoin
// should be used instead for such cases where consistent ordering is required.
//
// If the MutableHashSet is nil, MutableHashSet.Join returns an empty string.
func (s *MutableHashSet[E]) Join(sep string, convert func(element E) string) string {
	if s == nil {
		return ""
	}
	return internal.Join[E](s.elements, sep, convert)
}

// Len returns the number of elements within the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.Len returns zero.
func (s *MutableHashSet[E]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

// Max returns the maximum element within the MutableHashSet using the provided less function.
//
// If the MutableHashSet is nil, MutableHashSet.Max returns the zero value for E and false.
func (s *MutableHashSet[E]) Max(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Max[E](s.elements, less)
}

// Min returns the minimum element within the MutableHashSet using the provided less function.
//
// If the MutableHashSet is nil, MutableHashSet.Min returns the zero value for E and false.
func (s *MutableHashSet[E]) Min(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return internal.Min[E](s.elements, less)
}

// Mutable returns a reference to itself to conform with Set.Mutable.
//
// If the MutableHashSet is nil, MutableHashSet.Mutable returns nil.
func (s *MutableHashSet[E]) Mutable() MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return s
}

// None returns whether the MutableHashSet contains no elements that match the predicate function.
//
// If the MutableHashSet is nil, MutableHashSet.None returns true.
func (s *MutableHashSet[E]) None(predicate func(element E) bool) bool {
	if s == nil {
		return true
	}
	return internal.None[E](s.elements, predicate)
}

// Put adds the element to the MutableHashSet as well as any additional elements specified. Nothing changes for elements
// that already exist within the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.Put is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) Put(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.Put[E](s.elements, element, elements)
	return s
}

// PutAll adds all elements in the specified Set to the MutableHashSet. Nothing changes for elements that already exist
// within the SyncHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.PutAll is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) PutAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.PutAll[E](s.elements, elements)
	return s
}

// PutSlice adds all elements in the specified slice to the MutableHashSet. Nothing changes for elements that already
// exist within the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.PutSlice is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) PutSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	internal.PutSlice[E](s.elements, elements)
	return s
}

// Range calls the iter function with each element within the MutableHashSet but will stop early whenever the iter
// function returns true.
//
// Iteration order is not guaranteed to be consistent.
//
// If the MutableHashSet is nil, MutableHashSet.Range is a no-op.
func (s *MutableHashSet[E]) Range(iter func(element E) bool) {
	if s != nil {
		internal.Range[E](s.elements, iter)
	}
}

// Retain removes all elements from the MutableHashSet except the element(s) specified.
//
// If the MutableHashSet is nil, MutableHashSet.Retain is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) Retain(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	s.elements = internal.Retaining[E](s.elements, element, elements)
	return s
}

// RetainAll removes all elements from the MutableHashSet except those in the specified Set.
//
// If the MutableHashSet is nil, MutableHashSet.RetainAll is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) RetainAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	s.elements = internal.RetainingAll[E](s.elements, elements)
	return s
}

// RetainSlice removes all elements from the MutableHashSet except those in the specified slice.
//
// If the MutableHashSet is nil, MutableHashSet.RetainSlice is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) RetainSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	s.elements = internal.RetainingSlice[E](s.elements, elements)
	return s
}

// RetainWhere removes all elements except those that match the predicate function from the MutableHashSet.
//
// If the MutableHashSet is nil, MutableHashSet.RetainWhere is a no-op.
//
// A reference to the MutableHashSet is returned for method chaining.
func (s *MutableHashSet[E]) RetainWhere(predicate func(element E) bool) MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	s.elements = internal.RetainingWhere[E](s.elements, predicate)
	return s
}

// Slice returns a slice containing all elements of the MutableHashSet.
//
// The order of elements within the resulting slice is not guaranteed to be consistent. MutableHashSet.SortedSlice
// should be used instead for such cases where consistent ordering is required.
//
// If the MutableHashSet is nil, MutableHashSet.Slice returns nil.
func (s *MutableHashSet[E]) Slice() []E {
	if s == nil {
		return nil
	}
	return internal.Slice[E](s.elements)
}

// Some returns whether the MutableHashSet contains any element that matches the predicate function.
//
// If the MutableHashSet is nil, MutableHashSet.Some returns false.
func (s *MutableHashSet[E]) Some(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	return internal.Some[E](s.elements, predicate)
}

// SortedJoin sorts the elements within the MutableHashSet using the provided less function and then converts those
// elements into strings which are then joined using the specified separator to create the resulting string.
//
// If the MutableHashSet is nil, MutableHashSet.SortedJoin returns an empty string.
func (s *MutableHashSet[E]) SortedJoin(sep string, convert func(element E) string, less func(x, y E) bool) string {
	if s == nil {
		return ""
	}
	return internal.SortedJoin[E](s.elements, sep, convert, less)
}

// SortedSlice returns a slice containing all elements of the MutableHashSet sorted using the provided less function.
//
// If the MutableHashSet is nil, MutableHashSet.SortedSlice returns nil.
func (s *MutableHashSet[E]) SortedSlice(less func(x, y E) bool) []E {
	if s == nil {
		return nil
	}
	return internal.SortedSlice[E](s.elements, less)
}

// TryRange calls the iter function with each element within the MutableHashSet but will stop early whenever the iter
// function returns an error.
//
// Iteration order is not guaranteed to be consistent.
//
// If the MutableHashSet is nil, MutableHashSet.TryRange is a no-op.
func (s *MutableHashSet[E]) TryRange(iter func(element E) error) error {
	if s == nil {
		return nil
	}
	return internal.TryRange[E](s.elements, iter)
}

// Union returns a new MutableHashSet containing a union of the MutableHashSet with another Set.
//
// If the MutableHashSet and the other Set are both nil, MutableHashSet.Union returns nil.
func (s *MutableHashSet[E]) Union(other Set[E]) Set[E] {
	if elements := internal.Union[E](s, other); elements != nil {
		return &MutableHashSet[E]{elements}
	}
	var ns *MutableHashSet[E]
	return ns
}

func (s *MutableHashSet[E]) String() string {
	if s == nil {
		return internal.NilString
	}
	return internal.String[E](s.elements)
}

func (s *MutableHashSet[E]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return internal.MarshalJSONNil()
	}
	return internal.MarshalJSON[E](s.elements)
}

func (s *MutableHashSet[E]) UnmarshalJSON(data []byte) error {
	if elements, err := internal.UnmarshalJSON[E](data); err != nil {
		return err
	} else {
		s.elements = elements
		return nil
	}
}

// MutableHash returns a MutableHashSet struct that implements MutableSet containing each unique element provided.
//
// As MutableHash returns a mutable struct it is not safe for concurrent use by multiple goroutines. SyncHash should be
// used instead for such cases where mutability is required, otherwise Hash for a simple immutable Set.
func MutableHash[E comparable](elements ...E) *MutableHashSet[E] {
	return &MutableHashSet[E]{internal.FromSlice[E](elements)}
}

// MutableHashFromJSON returns a MutableHashSet struct that implements MutableSet containing each unique element parsed
// from the JSON-encoded data provided.
//
// As MutableHashFromJSON returns a mutable struct it is not safe for concurrent use by multiple goroutines.
// SyncHashFromJSON should be used instead for such cases where mutability is required, otherwise HashFromJSON for a
// simple immutable Set.
func MutableHashFromJSON[E comparable](data []byte) (*MutableHashSet[E], error) {
	set := &MutableHashSet[E]{}
	if err := json.Unmarshal(data, set); err != nil {
		return nil, err
	}
	return set, nil
}

// MutableHashFromSlice returns a MutableHashSet struct that implements MutableSet containing each unique element from
// the slice provided.
//
// As MutableHashFromSlice returns a mutable struct it is not safe for concurrent use by multiple goroutines.
// SyncHashFromSlice should be used instead for such cases where mutability is required, otherwise HashFromSlice for a
// simple immutable Set.
func MutableHashFromSlice[E comparable](elements []E) *MutableHashSet[E] {
	return &MutableHashSet[E]{internal.FromSlice[E](elements)}
}
