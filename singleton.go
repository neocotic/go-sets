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

// SingletonSet is an immutable implementation of Set that contains a single datum.
//
// As SingletonSet is immutable it is safe for concurrent use by multiple goroutines without additional locking or
// coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use. That said; SingletonSet only implements json.Unmarshaler for the purpose of being able to
// have a SingletonSet field value on a struct being unmarshalled. It's recommended to unmarshal JSON into a
// SingletonSet using SingletonFromJSON as JSON is typically only unmarshalled into a struct once.
type SingletonSet[E comparable] struct {
	element E
}

var (
	_ Set[any]         = (*SingletonSet[any])(nil)
	_ fmt.Stringer     = (*SingletonSet[any])(nil)
	_ json.Marshaler   = (*SingletonSet[any])(nil)
	_ json.Unmarshaler = (*SingletonSet[any])(nil)
)

// Clone returns a clone of the SingletonSet.
//
// If the SingletonSet is nil, SingletonSet.Clone returns nil.
func (s *SingletonSet[E]) Clone() Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	return &SingletonSet[E]{s.element}
}

// Contains returns whether the SingletonSet contains the element.
//
// If the SingletonSet is nil, SingletonSet.Contains returns false.
func (s *SingletonSet[E]) Contains(element E) bool {
	return s != nil && s.element == element
}

// Diff returns a new SingletonSet struct containing the element of the SingletonSet if it does not exist in another
// Set; otherwise an EmptySet.
//
// If the SingletonSet is nil, SingletonSet.Diff returns nil.
func (s *SingletonSet[E]) Diff(other Set[E]) Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	diff := internal.Diff[E](internal.Singleton(s.element), other)
	if element, ok := internal.TakeOne(diff); ok {
		return &SingletonSet[E]{element}
	}
	return &EmptySet[E]{}
}

// DiffSymmetric returns a new HashSet struct containing elements that exist within the SingletonSet or another Set, but
// not both.
//
// If the SingletonSet is nil, SingletonSet.DiffSymmetric returns nil.
func (s *SingletonSet[E]) DiffSymmetric(other Set[E]) Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	return &HashSet[E]{internal.DiffSymmetric[E](internal.Singleton(s.element), other)}
}

// Equal returns whether the other Set also contains the same element.
//
// If the SingletonSet is nil it is treated as having no elements and the same logic applies to the other Set. To
// clarify; this means that a nil Set is equal to a non-nil Set that contains no elements.
func (s *SingletonSet[E]) Equal(other Set[E]) bool {
	if s == nil {
		return other == nil || other.IsEmpty()
	} else if other == nil {
		return false
	}
	return internal.ContainsOnly[E](internal.Singleton(s.element), other.Slice())
}

// Every returns whether the element within the SingletonSet matches the predicate function.
//
// If the SingletonSet is nil, SingletonSet.Every returns false.
func (s *SingletonSet[E]) Every(predicate func(element E) bool) bool {
	return s != nil && predicate(s.element)
}

// Filter returns a clone of the SingletonSet if its element matches the filter function; otherwise an EmptySet.
//
// If the SingletonSet is nil, SingletonSet.Filter returns nil.
func (s *SingletonSet[E]) Filter(filter func(element E) bool) Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	if filter(s.element) {
		return &SingletonSet[E]{s.element}
	}
	return &EmptySet[E]{}
}

// Find returns the element within the SingletonSet if it matches the search function as well as an indication of
// whether it was matched.
//
// If the SingletonSet is nil, SingletonSet.Find returns the zero value for E and false.
func (s *SingletonSet[E]) Find(search func(element E) bool) (E, bool) {
	if s != nil && search(s.element) {
		return s.element, true
	}
	var zero E
	return zero, false
}

// Immutable returns a reference to itself to conform with Set.Immutable.
//
// If the SingletonSet is nil, SingletonSet.Immutable returns nil.
func (s *SingletonSet[E]) Immutable() Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	return s
}

// Intersection returns a clone of the SingletonSet if its element exists in another Set; otherwise an EmptySet.
//
// If the SingletonSet is nil, SingletonSet.Intersection returns nil.
func (s *SingletonSet[E]) Intersection(other Set[E]) Set[E] {
	if s == nil {
		var ns *SingletonSet[E]
		return ns
	}
	intersection := internal.Intersection[E](internal.Singleton(s.element), other)
	if element, ok := internal.TakeOne(intersection); ok {
		return &SingletonSet[E]{element}
	}
	return &EmptySet[E]{}
}

// IsEmpty returns whether the SingletonSet is nil to conform with Set.IsEmpty.
func (s *SingletonSet[E]) IsEmpty() bool {
	return s == nil
}

// IsMutable always returns false to conform with Set.IsMutable.
func (s *SingletonSet[E]) IsMutable() bool {
	return false
}

// Join returns the element within the SingletonSet converted to a string to conform with Set.Join.
//
// If the SingletonSet is nil, SingletonSet.Join returns an empty string.
func (s *SingletonSet[E]) Join(_ string, convert func(element E) string) string {
	if s == nil {
		return ""
	}
	return convert(s.element)
}

// Len returns one if the SingletonSet is not nil; otherwise zero.
func (s *SingletonSet[E]) Len() int {
	if s == nil {
		return 0
	}
	return 1
}

// Max returns the element within the SingletonSet to conform with Set.Max.
//
// If the SingletonSet is nil, SingletonSet.Max returns the zero value for E and false.
func (s *SingletonSet[E]) Max(_ func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return s.element, true
}

// Min returns the element within the SingletonSet to conform with Set.Min.
//
// If the SingletonSet is nil, SingletonSet.Min returns the zero value for E and false.
func (s *SingletonSet[E]) Min(_ func(x, y E) bool) (E, bool) {
	if s == nil {
		var zero E
		return zero, false
	}
	return s.element, true
}

// Mutable returns a mutable clone of the SingletonSet.
//
// If the SingletonSet is nil, SingletonSet.Mutable returns nil.
func (s *SingletonSet[E]) Mutable() MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return MutableHash[E](s.element)
}

// None returns whether the element within the SingletonSet does not match the predicate function.
//
// If the SingletonSet is nil, SingletonSet.None returns true.
func (s *SingletonSet[E]) None(predicate func(element E) bool) bool {
	return s == nil || !predicate(s.element)
}

// Range calls the iter function with the element within the SingletonSet.
//
// If the SingletonSet is nil, SingletonSet.Range is a no-op.
func (s *SingletonSet[E]) Range(iter func(element E) bool) {
	if s == nil {
		return
	}
	iter(s.element)
}

// Slice returns a slice containing the element within the SingletonSet.
//
// If the SingletonSet is nil, SingletonSet.Slice returns nil.
func (s *SingletonSet[E]) Slice() []E {
	if s == nil {
		return nil
	}
	return []E{s.element}
}

// Some returns whether the element within the SingletonSet matches the predicate function.
//
// If the SingletonSet is nil, SingletonSet.Some returns false.
func (s *SingletonSet[E]) Some(predicate func(element E) bool) bool {
	return s != nil && predicate(s.element)
}

// SortedJoin returns the element within the SingletonSet converted to a string to conform with Set.SortedJoin.
//
// If the SingletonSet is nil, SingletonSet.SortedJoin returns an empty string.
func (s *SingletonSet[E]) SortedJoin(sep string, convert func(element E) string, _ func(x, y E) bool) string {
	return s.Join(sep, convert)
}

// SortedSlice returns a slice containing the element within the SingletonSet to conform with Set.SortedSlice.
//
// If the SingletonSet is nil, SingletonSet.SortedSlice returns nil.
func (s *SingletonSet[E]) SortedSlice(_ func(x, y E) bool) []E {
	return s.Slice()
}

// TryRange calls the iter function with the element within the SingletonSet, which may return an error.
//
// If the SingletonSet is nil, SingletonSet.TryRange is a no-op.
func (s *SingletonSet[E]) TryRange(iter func(element E) error) error {
	if s == nil {
		return nil
	}
	return iter(s.element)
}

// Union returns a new immutable Set containing a union of the SingletonSet with another Set.
//
// If the SingletonSet and the other Set are both nil, SingletonSet.Union returns nil.
func (s *SingletonSet[E]) Union(other Set[E]) Set[E] {
	if elements := internal.Union[E](s, other); elements != nil {
		if len(elements) == 1 {
			if element, ok := internal.TakeOne(elements); ok {
				return &SingletonSet[E]{element}
			}
		}
		return &HashSet[E]{elements}
	}
	var ns *SingletonSet[E]
	return ns
}

func (s *SingletonSet[E]) String() string {
	return fmt.Sprintf("%v", s.Slice())
}

func (s *SingletonSet[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

func (s *SingletonSet[E]) UnmarshalJSON(data []byte) error {
	var elements []E
	if err := json.Unmarshal(data, &elements); err != nil {
		return err
	}
	if l := len(elements); l != 1 {
		return fmtErrJSONElementCount(1, l)
	}
	s.element = elements[0]
	return nil
}

// Singleton returns an immutable SingletonSet struct that implements Set containing a single datum.
//
// As Singleton returns an immutable struct it is safe for concurrent use by multiple goroutines without additional
// locking or coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use. That said; SingletonSet only implements json.Unmarshaler for the purpose of being able to
// have a SingletonSet field value on a struct being unmarshalled. It's recommended to unmarshal JSON into a
// SingletonSet using SingletonFromJSON as JSON is typically only unmarshalled into a struct once.
func Singleton[E comparable](element E) *SingletonSet[E] {
	return &SingletonSet[E]{element}
}

// SingletonFromJSON returns an immutable SingletonSet struct that implements Set containing a single datum parsed from
// the JSON-encoded data provided.
//
// As SingletonFromJSON returns an immutable struct it is safe for concurrent use by multiple goroutines without
// additional locking or coordination.
//
// The exception to its immutability is when passed to json.Unmarshal, however, this has been implemented in a way that
// is safe for concurrent use and, as JSON is typically only unmarshalled into a struct once, it's unlikely that this
// needs to be called on the returned SingletonSet again after calling this function.
func SingletonFromJSON[E comparable](data []byte) (*SingletonSet[E], error) {
	set := &SingletonSet[E]{}
	if err := json.Unmarshal(data, set); err != nil {
		return nil, err
	}
	return set, nil
}
