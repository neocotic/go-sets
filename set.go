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

type (
	// Set represents a data set which contains only unique elements.
	Set[E comparable] interface {
		// Clone returns a clone of the Set.
		//
		// The returned struct implementation of Set will always match that of the Set being cloned.
		//
		// If the Set is nil, Set.Clone returns nil.
		Clone() Set[E]
		// Contains returns whether the Set contains the element.
		//
		// If the Set is nil, Set.Contains returns false.
		Contains(element E) bool
		// Diff returns a new Set struct containing only elements of the Set that do not exist in another Set.
		//
		// The returned struct implementation of Set should match that of the Set, where possible, but must never differ
		// in mutability.
		//
		// If the Set is nil, Set.Diff returns nil.
		Diff(other Set[E]) Set[E]
		// DiffSymmetric returns a new Set struct containing elements that exist within the Set or another Set, but not
		// both.
		//
		// The returned struct implementation of Set should match that of the Set, where possible, but must never differ
		// in mutability.
		//
		// If the Set is nil, Set.DiffSymmetric returns nil.
		DiffSymmetric(other Set[E]) Set[E]
		// Equal returns whether the Set contains the exact same elements as another Set.
		//
		// If the Set is nil it is treated as having no elements and the same logic applies to the other Set. To
		// clarify; this means that a nil Set is equal to a non-nil Set that contains no elements.
		Equal(other Set[E]) bool
		// Every returns whether the Set contains elements that all match the predicate function.
		//
		// If the Set is nil, Set.Every returns false.
		Every(predicate func(element E) bool) bool
		// Filter returns a new Set struct containing only elements of the Set that match the filter function.
		//
		// The returned struct implementation of Set should match that of the Set being filtered, where possible, but
		// must never differ in mutability.
		//
		// If the Set is nil, Set.Filter returns nil.
		Filter(filter func(element E) bool) Set[E]
		// Find returns an element within the Set that matches the search function as well as an indication of whether a
		// match was found.
		//
		// Iteration order is not guaranteed to be consistent so results may vary.
		//
		// If the Set is nil, Set.Find returns the zero value for E and false.
		Find(search func(element E) bool) (E, bool)
		// Immutable returns an immutable version of the Set.
		//
		// The Set is returned if it is already immutable, otherwise an immutable clone is returned.
		//
		// If the Set is nil, Set.Immutable returns nil.
		Immutable() Set[E]
		// Intersection returns a new Set struct containing only elements of the Set that also exist in another Set.
		//
		// The returned struct implementation of Set should match that of the Set, where possible, but must never differ
		// in mutability.
		//
		// If the Set is nil, Set.Intersection returns nil.
		Intersection(other Set[E]) Set[E]
		// IsEmpty returns whether the Set contains no elements.
		//
		// If the Set is nil, Set.IsEmpty returns true.
		IsEmpty() bool
		// IsMutable returns whether the Set is mutable.
		IsMutable() bool
		// Join converts the elements within the Set to strings which are then concatenated to create a single string,
		// placing sep between the converted elements in the resulting string.
		//
		// The order of elements within the resulting string is not guaranteed to be consistent. Set.SortedJoin or
		// should be used instead for such cases where consistent ordering is required.
		//
		// If the Set is nil, Set.Join returns an empty string.
		Join(sep string, convert func(element E) string) string
		// Len returns the number of elements within the Set.
		//
		// If the Set is nil, Set.Len returns zero.
		Len() int
		// Max returns the maximum element within the Set using the provided less function.
		//
		// If the Set is nil, Set.Max returns the zero value for E and false.
		Max(less func(x, y E) bool) (E, bool)
		// Min returns the minimum element within the Set using the provided less function.
		//
		// If the Set is nil, Set.Min returns the zero value for E and false.
		Min(less func(x, y E) bool) (E, bool)
		// Mutable returns a mutable version of the Set.
		//
		// The Set is returned if it is already mutable, otherwise a mutable clone is returned.
		//
		// If the Set is nil, Set.Mutable returns nil.
		Mutable() MutableSet[E]
		// None returns whether the Set contains no elements that match the predicate function.
		//
		// If the Set is nil, Set.None returns true.
		None(predicate func(element E) bool) bool
		// Range calls the iter function with each element within the Set but will stop early whenever the iter function
		// returns true.
		//
		// Iteration order is not guaranteed to be consistent.
		//
		// If the Set is nil, Set.Range is a no-op.
		Range(iter func(element E) bool)
		// Slice returns a slice containing all elements of the Set.
		//
		// The order of elements within the resulting slice is not guaranteed to be consistent. Set.SortedSlice should
		// be used instead for such cases where consistent ordering is required.
		//
		// If the Set is nil, Set.Slice returns nil.
		Slice() []E
		// Some returns whether the Set contains any element that matches the predicate function.
		//
		// If the Set is nil, Set.Some returns false.
		Some(predicate func(element E) bool) bool
		// SortedJoin sorts the elements within the Set using the provided less function and then converts those
		// elements into strings which are then joined using the specified separator to create the resulting string.
		//
		// If the Set is nil, Set.SortedJoin returns an empty string.
		SortedJoin(sep string, convert func(element E) string, less func(x, y E) bool) string
		// SortedSlice returns a slice containing all elements of the Set sorted using the provided less function.
		//
		// If the Set is nil, Set.SortedSlice returns nil.
		SortedSlice(less func(x, y E) bool) []E
		// TryRange calls the iter function with each element within the Set but will stop early whenever the iter
		// function returns an error.
		//
		// Iteration order is not guaranteed to be consistent.
		//
		// If the Set is nil, Set.TryRange is a no-op.
		TryRange(iter func(element E) error) error
		// Union returns a new Set containing a union of the Set with another Set.
		//
		// The returned struct implementation of Set should match that of the Set, where possible, but must never
		// differ in mutability.
		//
		// If the Set and the other Set are both nil, Set.Union returns nil.
		Union(other Set[E]) Set[E]
	}

	// MutableSet represents a mutable Set.
	MutableSet[E comparable] interface {
		// Clear removes all elements from the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.Clear is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		Clear() MutableSet[E]
		// Delete removes the element from the MutableSet as well as any additional elements specified.
		//
		// If the MutableSet is nil, MutableSet.Delete is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		Delete(element E, elements ...E) MutableSet[E]
		// DeleteAll removes all elements in the specified Set from the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.DeleteAll is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		DeleteAll(elements Set[E]) MutableSet[E]
		// DeleteSlice removes all elements in the specified slice from the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.DeleteSlice is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		DeleteSlice(elements []E) MutableSet[E]
		// DeleteWhere removes all elements that match the predicate function from the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.DeleteWhere is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		DeleteWhere(predicate func(element E) bool) MutableSet[E]
		// Put adds the element to the MutableSet as well as any additional elements specified. Nothing changes for
		// elements that already exist within the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.Put is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		Put(element E, elements ...E) MutableSet[E]
		// PutAll adds all elements in the specified Set to the MutableSet. Nothing changes for elements that already
		// exist within the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.PutAll is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		PutAll(elements Set[E]) MutableSet[E]
		// PutSlice adds all elements in the specified slice to the MutableSet. Nothing changes for elements that
		// already exist within the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.PutSlice is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		PutSlice(elements []E) MutableSet[E]
		// Retain removes all elements from the MutableSet except the element(s) specified.
		//
		// If the MutableSet is nil, MutableSet.Retain is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		Retain(element E, elements ...E) MutableSet[E]
		// RetainAll removes all elements from the MutableSet except those in the specified Set.
		//
		// If the MutableSet is nil, MutableSet.RetainAll is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		RetainAll(elements Set[E]) MutableSet[E]
		// RetainSlice removes all elements from the MutableSet except those in the specified slice.
		//
		// If the MutableSet is nil, MutableSet.RetainSlice is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		RetainSlice(elements []E) MutableSet[E]
		// RetainWhere removes all elements except those that match the predicate function from the MutableSet.
		//
		// If the MutableSet is nil, MutableSet.RetainWhere is a no-op.
		//
		// A reference to the MutableSet is returned for method chaining.
		RetainWhere(predicate func(element E) bool) MutableSet[E]
		Set[E]
	}
)
