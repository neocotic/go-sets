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

// EmptySet is an immutable implementation of Set that contains no data.
//
// As EmptySet is immutable it is safe for concurrent use by multiple goroutines without additional locking or
// coordination.
type EmptySet[E comparable] struct{}

var (
	_ Set[any]         = (*EmptySet[any])(nil)
	_ fmt.Stringer     = (*EmptySet[any])(nil)
	_ json.Marshaler   = (*EmptySet[any])(nil)
	_ json.Unmarshaler = (*EmptySet[any])(nil)
)

// Clone returns a clone of the EmptySet.
//
// If the EmptySet is nil, EmptySet.Clone returns nil.
func (s *EmptySet[E]) Clone() Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	return &EmptySet[E]{}
}

// Contains always returns false to conform with Set.Contains.
func (s *EmptySet[E]) Contains(_ E) bool {
	return false
}

// Diff returns a new EmptySet struct to conform with Set.Diff.
//
// If the EmptySet is nil, EmptySet.Diff returns nil.
func (s *EmptySet[E]) Diff(_ Set[E]) Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	return &EmptySet[E]{}
}

// DiffSymmetric returns an immutable clone of another Set to conform with Set.DiffSymmetric.
//
// If the EmptySet is nil, EmptySet.DiffSymmetric returns nil.
func (s *EmptySet[E]) DiffSymmetric(other Set[E]) Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	if internal.IsNil(other) {
		return &EmptySet[E]{}
	}
	return other.Immutable()
}

// Equal returns whether the other Set also contains no elements.
//
// If the EmptySet is nil it is treated as having no elements and the same logic applies to the other Set. To clarify;
// this means that a nil Set is equal to a non-nil Set that contains no elements.
func (s *EmptySet[E]) Equal(other Set[E]) bool {
	return other == nil || other.IsEmpty()
}

// Every always returns false to conform with Set.Every.
func (s *EmptySet[E]) Every(_ func(element E) bool) bool {
	return false
}

// Filter returns a new EmptySet struct to conform with Set.Filter.
//
// If the EmptySet is nil, EmptySet.Filter returns nil.
func (s *EmptySet[E]) Filter(_ func(element E) bool) Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	return &EmptySet[E]{}
}

// Find always returns the zero value for E and false to conform with Set.Find.
func (s *EmptySet[E]) Find(_ func(element E) bool) (E, bool) {
	var zero E
	return zero, false
}

// Immutable returns a reference to itself to conform with Set.Immutable.
//
// If the EmptySet is nil, EmptySet.Immutable returns nil.
func (s *EmptySet[E]) Immutable() Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	return s
}

// Intersection returns a new EmptySet struct to conform with Set.Intersection.
//
// If the EmptySet is nil, EmptySet.Intersection returns nil.
func (s *EmptySet[E]) Intersection(_ Set[E]) Set[E] {
	if s == nil {
		var ns *EmptySet[E]
		return ns
	}
	return &EmptySet[E]{}
}

// IsEmpty always returns true to conform with Set.IsEmpty.
func (s *EmptySet[E]) IsEmpty() bool {
	return true
}

// IsMutable always returns false to conform with Set.IsMutable.
func (s *EmptySet[E]) IsMutable() bool {
	return false
}

// Join always returns an empty string to conform with Set.Join.
func (s *EmptySet[E]) Join(_ string, _ func(element E) string) string {
	return ""
}

// Len always returns zero to conform with Set.Len.
func (s *EmptySet[E]) Len() int {
	return 0
}

// Max always returns the zero value for E and false to conform with Set.Max.
func (s *EmptySet[E]) Max(_ func(x, y E) bool) (E, bool) {
	var zero E
	return zero, false
}

// Min always returns the zero value for E and false to conform with Set.Min.
func (s *EmptySet[E]) Min(_ func(x, y E) bool) (E, bool) {
	var zero E
	return zero, false
}

// Mutable returns a mutable clone of the EmptySet.
//
// If the EmptySet is nil, EmptySet.Mutable returns nil.
func (s *EmptySet[E]) Mutable() MutableSet[E] {
	if s == nil {
		var ns *MutableHashSet[E]
		return ns
	}
	return MutableHash[E]()
}

// None always returns true to conform with Set.None.
func (s *EmptySet[E]) None(_ func(element E) bool) bool {
	return true
}

// Range does nothing to conform with Set.Range.
func (s *EmptySet[E]) Range(_ func(element E) bool) {}

// Slice returns an empty slice to conform with Set.Slice.
//
// If the EmptySet is nil, EmptySet.Slice returns nil.
func (s *EmptySet[E]) Slice() []E {
	if s == nil {
		return nil
	}
	return make([]E, 0)
}

// Some always returns false to conform with Set.Some.
func (s *EmptySet[E]) Some(_ func(element E) bool) bool {
	return false
}

// SortedJoin always returns an empty string to conform with Set.SortedJoin.
func (s *EmptySet[E]) SortedJoin(_ string, _ func(element E) string, _ func(x, y E) bool) string {
	return ""
}

// SortedSlice returns an empty slice to conform with Set.SortedSlice.
//
// If the EmptySet is nil, EmptySet.SortedSlice returns nil.
func (s *EmptySet[E]) SortedSlice(_ func(x, y E) bool) []E {
	return s.Slice()
}

// TryRange does nothing and returns nil to conform with Set.TryRange.
func (s *EmptySet[E]) TryRange(_ func(element E) error) error {
	return nil
}

// Union returns a new immutable Set containing a union of the EmptySet with another Set.
//
// If the EmptySet and the other Set are both nil, EmptySet.Union returns nil.
func (s *EmptySet[E]) Union(other Set[E]) Set[E] {
	if elements := internal.Union[E](s, other); elements != nil {
		if len(elements) == 0 {
			return &EmptySet[E]{}
		}
		return &HashSet[E]{elements}
	}
	var ns *EmptySet[E]
	return ns
}

func (s *EmptySet[E]) String() string {
	return fmt.Sprintf("%v", s.Slice())
}

func (s *EmptySet[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Slice())
}

func (s *EmptySet[E]) UnmarshalJSON(data []byte) error {
	var elements []E
	if err := json.Unmarshal(data, &elements); err != nil {
		return err
	}
	if l := len(elements); l != 0 {
		return fmtErrJSONElementCount(0, l)
	}
	return nil
}

// Empty returns an immutable EmptySet struct that implements Set containing no data.
//
// As Empty returns an immutable struct it is safe for concurrent use by multiple goroutines without additional locking
// or coordination.
func Empty[E comparable]() *EmptySet[E] {
	return &EmptySet[E]{}
}

// EmptyFromJSON returns an immutable EmptySet struct that implements Set containing no data parsed from the
// JSON-encoded data provided.
//
// As EmptyFromJSON returns an immutable struct it is safe for concurrent use by multiple goroutines without additional
// locking or coordination.
//
// As EmptySet cannot contain any data, this function simply provides consistency with other Set implementations while
// also offering validation of sorts. That is; it will return an error if the JSON data does not form an empty array.
func EmptyFromJSON[E comparable](data []byte) (*EmptySet[E], error) {
	set := &EmptySet[E]{}
	if err := json.Unmarshal(data, set); err != nil {
		return nil, err
	}
	return set, nil
}
