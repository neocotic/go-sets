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

package internal

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Hash contains only unique elements.
type Hash[E comparable] map[E]struct{}

// NilString is a string representation of the elements within a nil Hash.
const NilString = "[]"

// Clone returns a clone of the Hash.
func Clone[E comparable](hash Hash[E]) Hash[E] {
	cloned := make(Hash[E])
	for element := range hash {
		cloned[element] = struct{}{}
	}
	return cloned
}

// ContainsOnly returns whether the Hash only contains the elements provided and no more or less.
func ContainsOnly[E comparable](hash Hash[E], elements []E) bool {
	if len(hash) != len(elements) {
		return false
	}
	for _, element := range elements {
		if _, ok := hash[element]; !ok {
			return false
		}
	}
	return true
}

// Delete removes the element from the Hash as well as any additional elements specified.
func Delete[E comparable](hash Hash[E], element E, elements []E) {
	delete(hash, element)
	for _, _element := range elements {
		delete(hash, _element)
	}
}

// DeleteAll removes all elements in the specified Collection from the Hash.
func DeleteAll[E comparable](hash Hash[E], elements Collection[E]) {
	if elements != nil {
		elements.Range(func(element E) bool {
			delete(hash, element)
			return false
		})
	}
}

// DeleteSlice removes all elements in the specified slice from the Hash.
func DeleteSlice[E comparable](hash Hash[E], elements []E) {
	for _, element := range elements {
		delete(hash, element)
	}
}

// DeleteWhere removes all elements that match the predicate function from the Hash.
func DeleteWhere[E comparable](hash Hash[E], predicate func(element E) bool) {
	for element := range hash {
		if predicate(element) {
			delete(hash, element)
		}
	}
}

// Diff returns a Hash containing only elements of the Hash that do not exist in the Collection provided.
func Diff[E comparable](hash Hash[E], elements Collection[E]) Hash[E] {
	if elements == nil {
		return Clone(hash)
	}
	diff := make(Hash[E])
	for element := range hash {
		if !elements.Contains(element) {
			diff[element] = struct{}{}
		}
	}
	return diff
}

// DiffAll returns a new Collection containing only elements of the specified Collection that do not exist in any other
// provided Collection.
//
// The Collection is inspected by the given flag function, allowing the tracking of its characteristics. The flags are
// then passed along with the Hash containing the differences to the specified factory function which is used to
// construct the Collection implementation that is returned by DiffAll.
func DiffAll[E comparable, C Collection[E]](
	factory func(hash Hash[E], flags CollectionFlag) C,
	flag func(col Collection[E]) CollectionFlag,
	col Collection[E],
	others []Collection[E],
) C {
	if IsNil(col) {
		return factory(nil, 0)
	}
	flags := flag(col)
	diff := make(Hash[E])
	var validOthers []Collection[E]
	for _, other := range others {
		if IsNotNil(other) {
			validOthers = append(validOthers, other)
		}
	}
	var iter func(element E) bool
	if len(validOthers) == 0 {
		iter = func(element E) bool {
			diff[element] = struct{}{}
			return false
		}
	} else {
		iter = func(element E) bool {
			var exists bool
			for _, other := range validOthers {
				if other.Contains(element) {
					exists = true
					break
				}
			}
			if !exists {
				diff[element] = struct{}{}
			}
			return false
		}
	}
	col.Range(iter)
	return factory(diff, flags)
}

// DiffSymmetric returns a Hash containing elements that exist within the Hash or the Collection provided, but not both.
func DiffSymmetric[E comparable](hash Hash[E], elements Collection[E]) Hash[E] {
	if elements == nil {
		return Clone(hash)
	}
	diff := make(Hash[E])
	for element := range hash {
		if !elements.Contains(element) {
			diff[element] = struct{}{}
		}
	}
	elements.Range(func(element E) bool {
		if _, ok := hash[element]; !ok {
			diff[element] = struct{}{}
		}
		return false
	})
	return diff
}

// DiffSymmetricAll returns a new Collection containing elements that exist within the specified Collection or any other
// Collection provided, but not in more than one.
//
// The Collection is inspected by the given flag function, allowing the tracking of its characteristics. The flags are
// then passed along with the Hash containing the differences to the specified factory function which is used to
// construct the Collection implementation that is returned by DiffSymmetricAll.
func DiffSymmetricAll[E comparable, C Collection[E]](
	factory func(hash Hash[E], flags CollectionFlag) C,
	flag func(col Collection[E]) CollectionFlag,
	col Collection[E],
	others []Collection[E],
) C {
	if IsNil(col) {
		return factory(nil, 0)
	}
	flags := flag(col)
	diff := make(Hash[E])
	tmp := make(Hash[E])
	col.Range(func(element E) bool {
		diff[element] = struct{}{}
		tmp[element] = struct{}{}
		return false
	})
	var validOthers []Collection[E]
	for _, other := range others {
		if IsNotNil(other) {
			validOthers = append(validOthers, other)
		}
	}
	for _, other := range validOthers {
		other.Range(func(element E) bool {
			if _, ok := tmp[element]; ok {
				delete(diff, element)
			} else {
				diff[element] = struct{}{}
				tmp[element] = struct{}{}
			}
			return false
		})
	}
	return factory(diff, flags)
}

// Every returns whether the Hash contains elements that all match the predicate function.
func Every[E comparable](hash Hash[E], predicate func(element E) bool) bool {
	if len(hash) == 0 {
		return false
	}
	for element := range hash {
		if !predicate(element) {
			return false
		}
	}
	return true
}

// Filter returns a Hash containing only elements of the Hash that match the filter function.
func Filter[E comparable](hash Hash[E], filter func(element E) bool) Hash[E] {
	filtered := make(Hash[E])
	for element := range hash {
		if filter(element) {
			filtered[element] = struct{}{}
		}
	}
	return filtered
}

// Find returns an element within the Hash that matches the search function as well as an indication of whether a match
// was found.
func Find[E comparable](hash Hash[E], search func(element E) bool) (E, bool) {
	for element := range hash {
		if search(element) {
			return element, true
		}
	}
	var zero E
	return zero, false
}

// FromSlice returns a Hash containing each unique element from the slice provided.
func FromSlice[E comparable](elements []E) Hash[E] {
	hash := make(Hash[E])
	for _, element := range elements {
		hash[element] = struct{}{}
	}
	return hash
}

// Intersection returns a Hash containing only elements of the Hash that also exist in the Collection provided.
func Intersection[E comparable](hash Hash[E], elements Collection[E]) Hash[E] {
	intersection := make(Hash[E])
	if elements != nil {
		elements.Range(func(element E) bool {
			if _, ok := hash[element]; ok {
				intersection[element] = struct{}{}
			}
			return false
		})
	}
	return intersection
}

// IntersectionAll returns a new Collection containing only elements of the specified Collection that also exist in any
// other provided Collection.
//
// The Collection is inspected by the given flag function, allowing the tracking of its characteristics. The flags are
// then passed along with the Hash containing the intersections to the specified factory function which is used to
// construct the Collection implementation that is returned by IntersectionAll.
func IntersectionAll[E comparable, C Collection[E]](
	factory func(hash Hash[E], flags CollectionFlag) C,
	flag func(col Collection[E]) CollectionFlag,
	col Collection[E],
	others []Collection[E],
) C {
	if IsNil(col) {
		return factory(nil, 0)
	}
	flags := flag(col)
	intersection := make(Hash[E])
	for _, other := range others {
		if IsNotNil(other) {
			other.Range(func(element E) bool {
				if col.Contains(element) {
					intersection[element] = struct{}{}
				}
				return false
			})
		}
	}
	return factory(intersection, flags)
}

// Join converts the elements within the Hash to strings which are then concatenated to create a single string, placing
// sep between the converted elements in the resulting string.
//
// The order of elements within the resulting string is not guaranteed to be consistent. SortedJoin should be used
// instead for such cases where consistent ordering is required.
func Join[E comparable](hash Hash[E], sep string, convert func(element E) string) string {
	var (
		i  int
		sb strings.Builder
	)
	for element := range hash {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(convert(element))
		i++
	}
	return sb.String()
}

// Map returns a Hash containing keys converted from the elements within the given Collection using the mapper function.
func Map[E comparable, T comparable](elements Collection[E], mapper func(element E) T) Hash[T] {
	mapped := make(Hash[T])
	var mappedElement T
	if elements != nil {
		elements.Range(func(element E) bool {
			mappedElement = mapper(element)
			mapped[mappedElement] = struct{}{}
			return false
		})
	}
	return mapped
}

// MarshalJSON returns the elements of the Hash serialized as a JSON array.
func MarshalJSON[E comparable](hash Hash[E]) ([]byte, error) {
	return json.Marshal(Slice(hash))
}

// MarshalJSONNil returns a serialization of a null JSON array used to represent a nil Hash.
func MarshalJSONNil() ([]byte, error) {
	return []byte("null"), nil
}

// Max returns the maximum element within the Hash using the provided less function.
func Max[E comparable](hash Hash[E], less func(x, y E) bool) (E, bool) {
	max, ok := TakeOne(hash)
	if !ok {
		return max, false
	}
	for element := range hash {
		if less(max, element) {
			max = element
		}
	}
	return max, true
}

// Min returns the minimum element within the Hash using the provided less function.
func Min[E comparable](hash Hash[E], less func(x, y E) bool) (E, bool) {
	min, ok := TakeOne(hash)
	if !ok {
		return min, false
	}
	for element := range hash {
		if less(element, min) {
			min = element
		}
	}
	return min, true
}

// None returns whether the Hash contains no elements that match the predicate function.
func None[E comparable](hash Hash[E], predicate func(element E) bool) bool {
	for element := range hash {
		if predicate(element) {
			return false
		}
	}
	return true
}

// Put adds the element to the Hash as well as any additional elements specified. Nothing changes for elements that
// already exist within the Hash.
func Put[E comparable](hash Hash[E], element E, elements []E) {
	hash[element] = struct{}{}
	for _, _element := range elements {
		hash[_element] = struct{}{}
	}
}

// PutAll adds all elements in the specified Collection to the Hash. Nothing changes for elements that already exist
// within the Hash.
func PutAll[E comparable](hash Hash[E], elements Collection[E]) {
	if elements != nil {
		elements.Range(func(element E) bool {
			hash[element] = struct{}{}
			return false
		})
	}
}

// PutSlice adds all elements in the specified slice to the Hash. Nothing changes for elements that already exist within
// the Hash.
func PutSlice[E comparable](hash Hash[E], elements []E) {
	for _, element := range elements {
		hash[element] = struct{}{}
	}
}

// Range calls the iter function with each element within the Hash but will stop early whenever the iter function
// returns true.
//
// Iteration order is not guaranteed to be consistent.
func Range[E comparable](hash Hash[E], iter func(element E) bool) {
	for element := range hash {
		if iter(element) {
			break
		}
	}
}

// Retaining returns a Hash containing only the specified element(s) if they exist in the Hash.
func Retaining[E comparable](hash Hash[E], element E, elements []E) Hash[E] {
	retained := make(Hash[E])
	if _, ok := hash[element]; ok {
		retained[element] = struct{}{}
	}
	for _, _element := range elements {
		if _, ok := hash[_element]; ok {
			retained[_element] = struct{}{}
		}
	}
	return retained
}

// RetainingAll returns a Hash containing only elements from the specified Collection if they exist in the Hash.
func RetainingAll[E comparable](hash Hash[E], elements Collection[E]) Hash[E] {
	retained := make(Hash[E])
	if elements != nil {
		elements.Range(func(element E) bool {
			if _, ok := hash[element]; ok {
				retained[element] = struct{}{}
			}
			return false
		})
	}
	return retained
}

// RetainingSlice returns a Hash containing only elements the specified slice if they exist in the Hash.
func RetainingSlice[E comparable](hash Hash[E], elements []E) Hash[E] {
	retained := make(Hash[E])
	for _, element := range elements {
		if _, ok := hash[element]; ok {
			retained[element] = struct{}{}
		}
	}
	return retained
}

// RetainingWhere returns a Hash containing only elements from Hash that match the predicate function.
func RetainingWhere[E comparable](hash Hash[E], predicate func(element E) bool) Hash[E] {
	retained := make(Hash[E])
	for element := range hash {
		if predicate(element) {
			retained[element] = struct{}{}
		}
	}
	return retained
}

// Singleton returns a Hash containing only the element provided.
func Singleton[E comparable](element E) Hash[E] {
	return Hash[E]{element: {}}
}

// Slice returns a slice containing all elements of the Hash.
//
// The order of elements within the resulting slice is not guaranteed to be consistent. SortedSlice should be used
// instead for such cases where consistent ordering is required.
func Slice[E comparable](hash Hash[E]) []E {
	var i int
	elements := make([]E, len(hash))
	for element := range hash {
		elements[i] = element
		i++
	}
	return elements
}

// Some returns whether the Hash contains any element that matches the predicate function.
func Some[E comparable](hash Hash[E], predicate func(element E) bool) bool {
	for element := range hash {
		if predicate(element) {
			return true
		}
	}
	return false
}

// SortedJoin sorts the elements within the Hash using the provided less function and then converts those elements into
// strings which are then joined using the specified separator to create the resulting string.
func SortedJoin[E comparable](hash Hash[E], sep string, convert func(element E) string, less func(x, y E) bool) string {
	elements := SortedSlice(hash, less)
	converted := make([]string, len(elements))
	for i, element := range elements {
		converted[i] = convert(element)
	}
	return strings.Join(converted, sep)
}

// SortedSlice returns a slice containing all elements of the Hash sorted using the provided less function.
func SortedSlice[E comparable](hash Hash[E], less func(x, y E) bool) []E {
	elements := Slice(hash)
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	return elements
}

// String returns a string representation of the elements within the Hash.
func String[E comparable](hash Hash[E]) string {
	return fmt.Sprintf("%v", Slice(hash))
}

// TakeOne returns any element within the Hash as well as an indication of whether the Hash contains any elements.
func TakeOne[E comparable](hash Hash[E]) (element E, ok bool) {
	for element = range hash {
		ok = true
		break
	}
	return
}

// TryMap returns a Hash containing keys converted from elements within the given Collection using the mapper function,
// which may return an error should an element fail to be mapped.
func TryMap[E comparable, T comparable](
	elements Collection[E],
	mapper func(element E) (T, error),
) (Hash[T], error) {
	mapped := make(Hash[T])
	var (
		err           error
		mappedElement T
	)
	if elements != nil {
		elements.Range(func(element E) bool {
			if mappedElement, err = mapper(element); err != nil {
				return true
			} else {
				mapped[mappedElement] = struct{}{}
				return false
			}
		})
	}
	if err != nil {
		return nil, err
	}
	return mapped, nil
}

// TryRange calls the iter function with each element within the Hash but will stop early whenever the iter function
// returns an error.
//
// Iteration order is not guaranteed to be consistent.
func TryRange[E comparable](hash Hash[E], iter func(element E) error) error {
	for element := range hash {
		if err := iter(element); err != nil {
			return err
		}
	}
	return nil
}

// Union returns a Hash containing all elements from across each provided Collection.
//
// The Hash is only allocated when at least one Collection is not nil.
func Union[E comparable](col, other Collection[E]) Hash[E] {
	var hash Hash[E]
	if IsNotNil(col) {
		hash = make(Hash[E])
		col.Range(func(element E) bool {
			hash[element] = struct{}{}
			return false
		})
	}
	if IsNotNil(other) {
		if hash == nil {
			hash = make(Hash[E])
		}
		other.Range(func(element E) bool {
			hash[element] = struct{}{}
			return false
		})
	}
	return hash
}

// UnionAll returns a new Collection containing all elements from across each Collection provided.
//
// Each non-nil Collection is inspected by the given flag function, allowing the tracking of characteristics of each
// Collection. The final flags are then passed along with the Hash containing all elements to the specified factory
// function which is used to construct the Collection implementation that is returned by UnionAll.
func UnionAll[E comparable, C Collection[E]](
	factory func(hash Hash[E], flags CollectionFlag) C,
	flag func(col Collection[E]) CollectionFlag,
	col Collection[E],
	others []Collection[E],
) C {
	var (
		flags CollectionFlag
		hash  Hash[E]
	)
	if IsNotNil(col) {
		flags |= flag(col)
		hash = make(Hash[E])
		col.Range(func(element E) bool {
			hash[element] = struct{}{}
			return false
		})
	}
	for _, other := range others {
		if IsNotNil(other) {
			flags |= flag(other)
			if hash == nil {
				hash = make(Hash[E])
			}
			other.Range(func(element E) bool {
				hash[element] = struct{}{}
				return false
			})
		}
	}
	return factory(hash, flags)
}

// UnmarshalJSON deserializes the given JSON data as a JSON array and returns a Hash containing each unique element.
func UnmarshalJSON[E comparable](data []byte) (Hash[E], error) {
	var elements []E
	if err := json.Unmarshal(data, &elements); err != nil {
		return nil, err
	}
	return FromSlice(elements), nil
}
