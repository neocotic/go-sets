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
	"sync"
)

// SyncHashSet is an implementation of MutableSet that contains a unique data set.
//
// While SyncHashSet is mutable it is safe for concurrent use by multiple goroutines without additional locking or
// coordination due to internal locking. If mutability is not required HashSet is a cheaper alternative.
type SyncHashSet[E comparable] struct {
	elements internal.Hash[E]
	mu       sync.RWMutex
}

var (
	_ MutableSet[any]  = (*SyncHashSet[any])(nil)
	_ fmt.Stringer     = (*SyncHashSet[any])(nil)
	_ json.Marshaler   = (*SyncHashSet[any])(nil)
	_ json.Unmarshaler = (*SyncHashSet[any])(nil)
)

// Clear removes all elements from the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.Clear is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) Clear() MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = make(internal.Hash[E])
	return s
}

// Clone returns a clone of the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.Clone returns nil.
func (s *SyncHashSet[E]) Clone() Set[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncHashSet[E]{elements: internal.Clone[E](s.elements)}
}

// Contains returns whether the SyncHashSet contains the element.
//
// If the SyncHashSet is nil, SyncHashSet.Contains returns false.
func (s *SyncHashSet[E]) Contains(element E) bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.elements[element]
	return ok
}

// Delete removes the element from the SyncHashSet as well as any additional elements specified.
//
// If the SyncHashSet is nil, SyncHashSet.Delete is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) Delete(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.Delete[E](s.elements, element, elements)
	return s
}

// DeleteAll removes all elements in the specified Set from the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.DeleteAll is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) DeleteAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.DeleteAll[E](s.elements, elements)
	return s
}

// DeleteSlice removes all elements in the specified slice from the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.DeleteSlice is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) DeleteSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.DeleteSlice[E](s.elements, elements)
	return s
}

// DeleteWhere removes all elements that match the predicate function from the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.DeleteWhere is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) DeleteWhere(predicate func(element E) bool) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.DeleteWhere[E](s.elements, predicate)
	return s
}

// Diff returns a new SyncHashSet struct containing only elements of the SyncHashSet that do not exist in another Set.
//
// If the SyncHashSet is nil, SyncHashSet.Diff returns nil.
func (s *SyncHashSet[E]) Diff(other Set[E]) Set[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncHashSet[E]{elements: internal.Diff[E](s.elements, other)}
}

// DiffSymmetric returns a new SyncHashSet struct containing elements that exist within the SyncHashSet or another Set,
// but not both.
//
// If the SyncHashSet is nil, SyncHashSet.DiffSymmetric returns nil.
func (s *SyncHashSet[E]) DiffSymmetric(other Set[E]) Set[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncHashSet[E]{elements: internal.DiffSymmetric[E](s.elements, other)}
}

// Equal returns whether the SyncHashSet contains the exact same elements as another Set.
//
// If the SyncHashSet is nil it is treated as having no elements and the same logic applies to the other Set. To
// clarify; this means that a nil Set is equal to a non-nil Set that contains no elements.
func (s *SyncHashSet[E]) Equal(other Set[E]) bool {
	if s == nil {
		return other == nil || other.IsEmpty()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if other == nil {
		return len(s.elements) == 0
	}
	return internal.ContainsOnly[E](s.elements, other.Slice())
}

// Every returns whether the SyncHashSet contains elements that all match the predicate function.
//
// If the SyncHashSet is nil, SyncHashSet.Every returns false.
func (s *SyncHashSet[E]) Every(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Every[E](s.elements, predicate)
}

// Filter returns a new SyncHashSet struct containing only elements of the SyncHashSet that match the filter function.
//
// If the SyncHashSet is nil, SyncHashSet.Filter returns nil.
func (s *SyncHashSet[E]) Filter(filter func(element E) bool) Set[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncHashSet[E]{elements: internal.Filter[E](s.elements, filter)}
}

// Find returns an element within the SyncHashSet that matches the search function as well as an indication of whether a
// match was found.
//
// Iteration order is not guaranteed to be consistent so results may vary.
//
// If the SyncHashSet is nil, SyncHashSet.Find returns the zero value for E and false.
func (s *SyncHashSet[E]) Find(search func(element E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Find[E](s.elements, search)
}

// Immutable returns an immutable clone of the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.Immutable returns nil.
func (s *SyncHashSet[E]) Immutable() Set[E] {
	if s == nil {
		var ns *HashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &HashSet[E]{internal.Clone[E](s.elements)}
}

// Intersection returns a new SyncHashSet struct containing only elements of the SyncHashSet that also exist in another
// Set.
//
// If the SyncHashSet is nil, SyncHashSet.Intersection returns nil.
func (s *SyncHashSet[E]) Intersection(other Set[E]) Set[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncHashSet[E]{elements: internal.Intersection[E](s.elements, other)}
}

// IsEmpty returns whether the SyncHashSet contains no elements.
//
// If the SyncHashSet is nil, SyncHashSet.IsEmpty returns true.
func (s *SyncHashSet[E]) IsEmpty() bool {
	if s == nil {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.elements) == 0
}

// IsMutable always returns true to conform with Set.IsMutable.
func (s *SyncHashSet[E]) IsMutable() bool {
	return true
}

// Join converts the elements within the SyncHashSet to strings which are then concatenated to create a single string,
// placing sep between the converted elements in the resulting string.
//
// The order of elements within the resulting string is not guaranteed to be consistent. SyncHashSet.SortedJoin should
// be used instead for such cases where consistent ordering is required.
//
// If the SyncHashSet is nil, SyncHashSet.Join returns an empty string.
func (s *SyncHashSet[E]) Join(sep string, convert func(element E) string) string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Join[E](s.elements, sep, convert)
}

// Len returns the number of elements within the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.Len returns zero.
func (s *SyncHashSet[E]) Len() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.elements)
}

// Max returns the maximum element within the SyncHashSet using the provided less function.
//
// If the SyncHashSet is nil, SyncHashSet.Max returns the zero value for E and false.
func (s *SyncHashSet[E]) Max(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Max[E](s.elements, less)
}

// Min returns the minimum element within the SyncHashSet using the provided less function.
//
// If the SyncHashSet is nil, SyncHashSet.Min returns the zero value for E and false.
func (s *SyncHashSet[E]) Min(less func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Min[E](s.elements, less)
}

// Mutable returns a reference to itself to conform with Set.Mutable.
//
// If the SyncHashSet is nil, SyncHashSet.Mutable returns nil.
func (s *SyncHashSet[E]) Mutable() MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	return s
}

// None returns whether the SyncHashSet contains no elements that match the predicate function.
//
// If the SyncHashSet is nil, SyncHashSet.None returns true.
func (s *SyncHashSet[E]) None(predicate func(element E) bool) bool {
	if s == nil {
		return true
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.None[E](s.elements, predicate)
}

// Put adds the element to the SyncHashSet as well as any additional elements specified. Nothing changes for elements
// that already exist within the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.Put is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) Put(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.Put[E](s.elements, element, elements)
	return s
}

// PutAll adds all elements in the specified Set to the SyncHashSet. Nothing changes for elements that already exist
// within the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.PutAll is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) PutAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.PutAll[E](s.elements, elements)
	return s
}

// PutSlice adds all elements in the specified slice to the SyncHashSet. Nothing changes for elements that already exist
// within the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.PutSlice is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) PutSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	internal.PutSlice[E](s.elements, elements)
	return s
}

// Range calls the iter function with each element within the SyncHashSet but will stop early whenever the iter function
// returns true.
//
// Iteration order is not guaranteed to be consistent.
//
// If the SyncHashSet is nil, SyncHashSet.Range is a no-op.
func (s *SyncHashSet[E]) Range(iter func(element E) bool) {
	if s == nil {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	internal.Range[E](s.elements, iter)
}

// Retain removes all elements from the SyncHashSet except the element(s) specified.
//
// If the SyncHashSet is nil, SyncHashSet.Retain is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) Retain(element E, elements ...E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = internal.Retaining[E](s.elements, element, elements)
	return s
}

// RetainAll removes all elements from the SyncHashSet except those in the specified Set.
//
// If the SyncHashSet is nil, SyncHashSet.RetainAll is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) RetainAll(elements Set[E]) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = internal.RetainingAll[E](s.elements, elements)
	return s
}

// RetainSlice removes all elements from the SyncHashSet except those in the specified slice.
//
// If the SyncHashSet is nil, SyncHashSet.RetainSlice is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) RetainSlice(elements []E) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = internal.RetainingSlice[E](s.elements, elements)
	return s
}

// RetainWhere removes all elements except those that match the predicate function from the SyncHashSet.
//
// If the SyncHashSet is nil, SyncHashSet.RetainWhere is a no-op.
//
// A reference to the SyncHashSet is returned for method chaining.
func (s *SyncHashSet[E]) RetainWhere(predicate func(element E) bool) MutableSet[E] {
	if s == nil {
		var ns *SyncHashSet[E]
		return ns
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.elements = internal.RetainingWhere[E](s.elements, predicate)
	return s
}

// Slice returns a slice containing all elements of the SyncHashSet.
//
// The order of elements within the resulting slice is not guaranteed to be consistent. SyncHashSet.SortedSlice should
// be used instead for such cases where consistent ordering is required.
//
// If the SyncHashSet is nil, SyncHashSet.Slice returns nil.
func (s *SyncHashSet[E]) Slice() []E {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Slice[E](s.elements)
}

// Some returns whether the SyncHashSet contains any element that matches the predicate function.
//
// If the SyncHashSet is nil, SyncHashSet.Some returns false.
func (s *SyncHashSet[E]) Some(predicate func(element E) bool) bool {
	if s == nil {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.Some[E](s.elements, predicate)
}

// SortedJoin sorts the elements within the SyncHashSet using the provided less function and then converts those
// elements into strings which are then joined using the specified separator to create the resulting string.
//
// If the SyncHashSet is nil, SyncHashSet.SortedJoin returns an empty string.
func (s *SyncHashSet[E]) SortedJoin(sep string, convert func(element E) string, less func(x, y E) bool) string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.SortedJoin[E](s.elements, sep, convert, less)
}

// SortedSlice returns a slice containing all elements of the SyncHashSet sorted using the provided less function.
//
// If the SyncHashSet is nil, SyncHashSet.SortedSlice returns nil.
func (s *SyncHashSet[E]) SortedSlice(less func(x, y E) bool) []E {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.SortedSlice[E](s.elements, less)
}

// TryRange calls the iter function with each element within the SyncHashSet but will stop early whenever the iter
// function returns an error.
//
// Iteration order is not guaranteed to be consistent.
//
// If the SyncHashSet is nil, SyncHashSet.TryRange is a no-op.
func (s *SyncHashSet[E]) TryRange(iter func(element E) error) error {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.TryRange[E](s.elements, iter)
}

// Union returns a new SyncHashSet containing a union of the SyncHashSet with another Set.
//
// If the SyncHashSet and the other Set are both nil, SyncHashSet.Union returns nil.
func (s *SyncHashSet[E]) Union(other Set[E]) Set[E] {
	if elements := internal.Union[E](s, other); elements != nil {
		return &SyncHashSet[E]{elements: elements}
	}
	var ns *SyncHashSet[E]
	return ns
}

func (s *SyncHashSet[E]) String() string {
	if s == nil {
		return internal.NilString
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.String[E](s.elements)
}

func (s *SyncHashSet[E]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return internal.MarshalJSONNil()
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return internal.MarshalJSON[E](s.elements)
}

func (s *SyncHashSet[E]) UnmarshalJSON(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if elements, err := internal.UnmarshalJSON[E](data); err != nil {
		return err
	} else {
		s.elements = elements
		return nil
	}
}

// SyncHash returns a SyncHashSet struct that implements MutableSet containing each unique element provided.
//
// While SyncHash returns a mutable struct it is safe for concurrent use by multiple goroutines without additional
// locking or coordination due to internal locking. If mutability is not required Hash provides a cheaper alternative.
func SyncHash[E comparable](elements ...E) *SyncHashSet[E] {
	return &SyncHashSet[E]{elements: internal.FromSlice[E](elements)}
}

// SyncHashFromJSON returns a SyncHashSet struct that implements MutableSet containing each unique element parsed from
// the JSON-encoded data provided.
//
// While SyncHashFromJSON returns a mutable struct it is safe for concurrent use by multiple goroutines without
// additional locking or coordination due to internal locking. If mutability is not required HashFromJSON provides a
// cheaper alternative.
func SyncHashFromJSON[E comparable](data []byte) (*SyncHashSet[E], error) {
	set := &SyncHashSet[E]{}
	if err := json.Unmarshal(data, set); err != nil {
		return nil, err
	}
	return set, nil
}

// SyncHashFromSlice returns a SyncHashSet struct that implements MutableSet containing each unique element from the
// slice provided.
//
// While SyncHashFromSlice returns a mutable struct it is safe for concurrent use by multiple goroutines without
// additional locking or coordination due to internal locking. If mutability is not required HashFromSlice provides a
// cheaper alternative.
func SyncHashFromSlice[E comparable](elements []E) *SyncHashSet[E] {
	return &SyncHashSet[E]{elements: internal.FromSlice[E](elements)}
}
