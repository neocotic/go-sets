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
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/neocotic/go-sets/internal"
	"sync"
	"testing"
)

const DefaultTestConcurrency = 10

func Test_SyncHash(t *testing.T) {
	testCases := map[string]struct {
		elements []int
	}{
		"with multiple elements": {
			elements: []int{123, 456, 789},
		},
		"with single element": {
			elements: []int{123},
		},
		"with no elements": {
			elements: []int{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := SyncHash(tc.elements...)
			if exp, act := len(tc.elements), set.Len(); act != exp {
				t.Errorf("unexpected Set length; want %v, got %v", exp, act)
			}
			if !set.IsMutable() {
				t.Error("unexpected Set mutability; want false, got true")
			}
		})
	}
}

func Test_SyncHashFromJSON(t *testing.T) {
	testCases := map[string]struct {
		expectElements []int
		json           string
	}{
		"with JSON string for array containing multiple elements": {
			expectElements: []int{123, 456, 789},
			json:           "[123,456,789]",
		},
		"with JSON string for array containing single element": {
			expectElements: []int{123},
			json:           "[123]",
		},
		"with JSON string for array containing duplicated elements": {
			expectElements: []int{123, 456, 789},
			json:           "[123,456,789,456,123]",
		},
		"with JSON string for array containing null element": {
			expectElements: []int{0},
			json:           "[null]",
		},
		"with JSON string for empty array": {
			expectElements: []int{},
			json:           "[]",
		},
		"with JSON string for null": {
			expectElements: []int{},
			json:           "null",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set, err := SyncHashFromJSON[int]([]byte(tc.json))
			if err != nil {
				t.Errorf("unexpected error; want nil, got %q", err)
			} else if set == nil {
				t.Error("unexpected nil Set")
			} else {
				if !set.IsMutable() {
					t.Error("unexpected Set mutability; want false, got true")
				}

				opts := []cmp.Option{cmpopts.SortSlices(Asc[int])}
				if actualElements := set.Slice(); !cmp.Equal(tc.expectElements, actualElements, opts...) {
					t.Errorf("unexpected unmarshalled elements; got diff %v", cmp.Diff(tc.expectElements, actualElements, opts...))
				}
			}
		})
	}
}

func Test_SyncHashFromSlice(t *testing.T) {
	testCases := map[string]struct {
		elements []int
	}{
		"with slice containing multiple elements": {
			elements: []int{123, 456, 789},
		},
		"with slice containing single element": {
			elements: []int{123},
		},
		"with slice containing no elements": {
			elements: []int{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := SyncHashFromSlice(tc.elements)
			if exp, act := len(tc.elements), set.Len(); act != exp {
				t.Errorf("unexpected Set length; want %v, got %v", exp, act)
			}
			if !set.IsMutable() {
				t.Error("unexpected Set mutability; want false, got true")
			}
		})
	}
}

func Test_SyncHashSet_Clear(t *testing.T) {
	testCases := map[string]struct {
		set *SyncHashSet[int]
	}{
		"on non-empty *SyncHashSet": {
			set: SyncHash(123, 456, 789),
		},
		"on empty *SyncHashSet": {
			set: SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.Clear()

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Clear_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Clear()
	})
}

func Test_SyncHashSet_Clear_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	ret := set.Clear()

	if internal.IsNotNil(ret) {
		t.Errorf("unexpected MutableSet; want nil, got %v", ret)
	}
	if !set.IsEmpty() {
		t.Error("unexpected MutableSet emptiness; want true, got false")
	}
}

func Test_SyncHashSet_Clone(t *testing.T) {
	set := SyncHash(123, 456, 789)
	clone := set.Clone()
	if internal.IsNil(clone) {
		t.Error("unexpected nil Set")
	}
	if l := clone.Len(); l != 3 {
		t.Errorf("unexpected cloned Set length; want 3, got %v", l)
	}
	if !clone.Equal(set) {
		t.Errorf("unexpected cloned Set; want %v, got %v", set, clone)
	}
	if !clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want true, got false")
	}
}

func Test_SyncHashSet_Clone_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Clone()
	})
}

func Test_SyncHashSet_Clone_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	clone := set.Clone()
	if clone == nil {
		t.Error("unexpected nil Set")
	}
	if internal.IsNotNil(clone) {
		t.Errorf("unexpected cloned Set; want nil, got %#v", clone)
	}
	if !clone.IsEmpty() {
		t.Error("unexpected cloned Set emptiness; want true, got false")
	}
	if !clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want true, got false")
	}
}

func Test_SyncHashSet_Contains(t *testing.T) {
	testCases := map[string]struct {
		element int
		expect  bool
	}{
		"with matching element": {
			element: 123,
			expect:  true,
		},
		"with non-matching zero value for element": {
			element: 0,
			expect:  false,
		},
		"with non-matching non-zero value for element": {
			element: 1,
			expect:  false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := SyncHash(123, 456, 789)
			result := set.Contains(tc.element)
			if result != tc.expect {
				t.Errorf("unexpected element contained within Set: %q; want %v, got %v", tc.element, tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Contains_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Contains(123)
	})
}

func Test_SyncHashSet_Contains_Nil(t *testing.T) {
	testCases := map[string]struct {
		element int
	}{
		"with non-matching zero value for element":     {0},
		"with non-matching non-zero value for element": {1},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			if set.Contains(tc.element) {
				t.Errorf("unexpected element contained within Set: %q; want false, got true", tc.element)
			}
		})
	}
}

func Test_SyncHashSet_Delete(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with multiple elements that do not exist on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, -789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that all exist on non-empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that some exist on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, 789},
			expect:   Hash(123, 456),
			set:      SyncHash(123, 456, 789),
		},
		"with single element that does not exist on non-empty *SyncHashSet": {
			element: -123,
			expect:  Hash(123, 456, 789),
			set:     SyncHash(123, 456, 789),
		},
		"with single element that exists on non-empty *SyncHashSet": {
			element: 123,
			expect:  Hash(456, 789),
			set:     SyncHash(123, 456, 789),
		},
		"with multiple elements on empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with single element on empty *SyncHashSet": {
			element: 123,
			expect:  Hash[int](),
			set:     SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.Delete(tc.element, tc.elements...)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_Delete_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Delete(123)
	})
}

func Test_SyncHashSet_Delete_Nil(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
	}{
		"with multiple elements": {
			element:  123,
			elements: []int{456, 789},
		},
		"with single element": {
			element: 123,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			ret := set.Delete(tc.element, tc.elements...)

			if internal.IsNotNil(ret) {
				t.Errorf("unexpected MutableSet; want nil, got %v", ret)
			}
			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_DeleteAll(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with Set containing multiple elements that do not exist on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, -789),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, 789),
			expect:   Hash(123, 456),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing single element that does not exist on non-empty *SyncHashSet": {
			elements: Hash(-123),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing single element that exists on non-empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash(456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing no elements on non-empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements on empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with Set containing single element on empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with Set containing no elements on empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.DeleteAll(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_DeleteAll_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.DeleteAll(Singleton(123))
	})
}

func Test_SyncHashSet_DeleteAll_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
	}{
		"with Set containing multiple elements": {
			elements: Hash(123, 456, 789),
		},
		"with Set containing single element": {
			elements: Hash(123),
		},
		"with Set containing no elements": {
			elements: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.DeleteAll(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_DeleteSlice(t *testing.T) {
	testCases := map[string]struct {
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with slice containing multiple elements that do not exist on non-empty *SyncHashSet": {
			elements: []int{-123, -456, -789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: []int{-123, -456, 789},
			expect:   Hash(123, 456),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that does not exist on non-empty *SyncHashSet": {
			elements: []int{-123},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that exists on non-empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash(456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing no elements on non-empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements on empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing single element on empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing no elements on empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.DeleteSlice(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_DeleteSlice_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.DeleteSlice([]int{123})
	})
}

func Test_SyncHashSet_DeleteSlice_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements []int
	}{
		"with slice containing multiple elements": {
			elements: []int{123, 456, 789},
		},
		"with slice containing single element": {
			elements: []int{123},
		},
		"with slice containing no elements": {
			elements: []int{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.DeleteSlice(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_DeleteWhere(t *testing.T) {
	testCases := map[string]struct {
		expect        Set[int]
		predicateFunc func(element int) bool
		set           *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:        Hash(123, 456, 789),
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:        Hash(456, 789),
			predicateFunc: func(element int) bool { return element == 123 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *SyncHashSet": {
			expect:        Hash(-789, -456, -123, 0),
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:        Hash(123, 456, 789),
			predicateFunc: func(element int) bool { return element < 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.DeleteWhere(tc.predicateFunc)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_DeleteWhere_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.DeleteWhere(func(element int) bool { return true })
	})
}

func Test_SyncHashSet_DeleteWhere_Nil(t *testing.T) {
	testCases := map[string]struct {
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			predicateFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.DeleteWhere(tc.predicateFunc)
			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Diff(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SyncHashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *SyncHashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(-789, -456, -123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *SyncHashSet": {
			expect: Hash(456, 789),
			other:  Hash(-123, 0, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *SyncHashSet": {
			expect: Hash(789),
			other:  Hash(0, 123, 456),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    SyncHash(123, 456, 789),
		},
		"with empty Set on non-empty *SyncHashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash[int](),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set on empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with empty Set on empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := tc.set.Diff(tc.other)
			if internal.IsNil(diff) {
				t.Error("unexpected nil Set")
			}
			if !diff.Equal(tc.expect) {
				t.Errorf("unexpected diff Set; want %v, got %v", tc.expect, diff)
			}
			if !diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Diff_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Diff(Singleton(123))
	})
}

func Test_SyncHashSet_Diff_Nil(t *testing.T) {
	testCases := map[string]struct {
		other Set[int]
	}{
		"with non-empty Set": {
			other: Hash(123, 456, 789),
		},
		"with empty Set": {
			other: Hash[int](),
		},
		"with nil Set": {
			other: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			diff := set.Diff(tc.other)
			if diff == nil {
				t.Error("unexpected nil Set")
			}
			if internal.IsNotNil(diff) {
				t.Errorf("unexpected diff Set; want nil, got %#v", diff)
			}
			if !diff.IsEmpty() {
				t.Error("unexpected diff Set emptiness; want true, got false")
			}
			if !diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_DiffSymmetric(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SyncHashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *SyncHashSet": {
			expect: Hash(-789, -456, -123, 123, 456, 789),
			other:  Hash(-789, -456, -123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *SyncHashSet": {
			expect: Hash(-123, 0, 456, 789),
			other:  Hash(-123, 0, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *SyncHashSet": {
			expect: Hash(0, 789),
			other:  Hash(0, 123, 456),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    SyncHash(123, 456, 789),
		},
		"with empty Set on non-empty *SyncHashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash[int](),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set on empty *SyncHashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with empty Set on empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := tc.set.DiffSymmetric(tc.other)
			if internal.IsNil(diff) {
				t.Error("unexpected nil Set")
			}
			if !diff.Equal(tc.expect) {
				t.Errorf("unexpected diff Set; want %v, got %v", tc.expect, diff)
			}
			if !diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_DiffSymmetric_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.DiffSymmetric(Singleton(123))
	})
}

func Test_SyncHashSet_DiffSymmetric_Nil(t *testing.T) {
	testCases := map[string]struct {
		other Set[int]
	}{
		"with non-empty Set": {
			other: Hash(123, 456, 789),
		},
		"with empty Set": {
			other: Hash[int](),
		},
		"with nil Set": {
			other: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			diff := set.DiffSymmetric(tc.other)
			if diff == nil {
				t.Error("unexpected nil Set")
			}
			if internal.IsNotNil(diff) {
				t.Errorf("unexpected diff Set; want nil, got %#v", diff)
			}
			if !diff.IsEmpty() {
				t.Error("unexpected diff Set emptiness; want true, got false")
			}
			if !diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Equal(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
		set    *SyncHashSet[int]
	}{
		"with nil *SyncHashSet on non-empty *SyncHashSet": {
			expect: false,
			other:  (*SyncHashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *EmptySet on non-empty *SyncHashSet": {
			expect: false,
			other:  (*EmptySet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *HashSet on non-empty *SyncHashSet": {
			expect: false,
			other:  (*HashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *MutableHashSet on non-empty *SyncHashSet": {
			expect: false,
			other:  (*MutableHashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *SingletonSet on non-empty *SyncHashSet": {
			expect: false,
			other:  (*SingletonSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only same elements on non-empty *SyncHashSet": {
			expect: true,
			other:  SyncHash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements on non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different elements on non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *EmptySet on non-empty *SyncHashSet": {
			expect: false,
			other:  Empty[int](),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing only same elements on non-empty *SyncHashSet": {
			expect: true,
			other:  Hash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  Hash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements on non-empty *SyncHashSet": {
			expect: false,
			other:  Hash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  Hash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different elements on non-empty *SyncHashSet": {
			expect: false,
			other:  Hash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only same elements on non-empty *SyncHashSet": {
			expect: true,
			other:  MutableHash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  MutableHash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements on non-empty *SyncHashSet": {
			expect: false,
			other:  MutableHash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: false,
			other:  MutableHash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different elements on non-empty *SyncHashSet": {
			expect: false,
			other:  MutableHash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing same element on non-empty *SyncHashSet": {
			expect: true,
			other:  Singleton(123),
			set:    SyncHash(123),
		},
		"with non-nil *SingletonSet containing same element but not others on non-empty *SyncHashSet": {
			expect: false,
			other:  Singleton(123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing different element on non-empty *SyncHashSet": {
			expect: false,
			other:  Singleton(12),
			set:    SyncHash(123),
		},
		"with nil *SyncHashSet on empty *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *EmptySet on empty *SyncHashSet": {
			expect: true,
			other:  (*EmptySet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *HashSet on empty *SyncHashSet": {
			expect: true,
			other:  (*HashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *MutableHashSet on empty *SyncHashSet": {
			expect: true,
			other:  (*MutableHashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *SingletonSet on empty *SyncHashSet": {
			expect: true,
			other:  (*SingletonSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet on empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *SyncHashSet on empty *SyncHashSet": {
			expect: true,
			other:  SyncHash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil *EmptySet on empty *SyncHashSet": {
			expect: true,
			other:  Empty[int](),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *HashSet on empty *SyncHashSet": {
			expect: false,
			other:  Hash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *HashSet on empty *SyncHashSet": {
			expect: true,
			other:  Hash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *MutableHashSet on empty *SyncHashSet": {
			expect: false,
			other:  MutableHash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *MutableHashSet on empty *SyncHashSet": {
			expect: true,
			other:  MutableHash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil *SingletonSet on empty *SyncHashSet": {
			expect: false,
			other:  Singleton(123),
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Equal_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Equal(Hash(123, 456, 789))
	})
}

func Test_SyncHashSet_Equal_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
		"with nil *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: true,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			expect: true,
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: true,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			expect: true,
			other:  (*SingletonSet[int])(nil),
		},
		"with non-nil empty *SyncHashSet": {
			expect: true,
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(0),
		},
		"with non-nil *EmptySet": {
			expect: true,
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: true,
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: false,
			other:  Hash(0),
		},
		"with non-nil empty *MutableHashSet": {
			expect: true,
			other:  MutableHash[int](),
		},
		"with non-nil non-empty *MutableHashSet": {
			expect: false,
			other:  MutableHash(0),
		},
		"with non-nil *SingletonSet": {
			expect: false,
			other:  Singleton(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Every(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.Every(tc.predicateFunc)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Every_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Every(func(_ int) bool { return true })
	})
}

func Test_SyncHashSet_Every_Nil(t *testing.T) {
	testCases := map[string]struct {
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			predicateFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			result := set.Every(tc.predicateFunc)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
		})
	}
}

func Test_SyncHashSet_Filter(t *testing.T) {
	testCases := map[string]struct {
		expect     Set[int]
		filterFunc func(element int) bool
		set        *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:     SyncHash(123, 456, 789),
			filterFunc: func(_ int) bool { return true },
			set:        SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:     SyncHash[int](),
			filterFunc: func(_ int) bool { return false },
			set:        SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:     SyncHash(123, 456, 789),
			filterFunc: func(element int) bool { return element > 0 },
			set:        SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:     SyncHash(123),
			filterFunc: func(element int) bool { return element == 123 },
			set:        SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *SyncHashSet": {
			expect:     SyncHash(123, 456, 789),
			filterFunc: func(element int) bool { return element > 0 },
			set:        SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:     SyncHash[int](),
			filterFunc: func(element int) bool { return element < 0 },
			set:        SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:     SyncHash[int](),
			filterFunc: func(_ int) bool { return true },
			set:        SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:     SyncHash[int](),
			filterFunc: func(_ int) bool { return false },
			set:        SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			filtered := tc.set.Filter(tc.filterFunc)
			if internal.IsNil(filtered) {
				t.Error("unexpected nil Set")
			}
			if !filtered.Equal(tc.expect) {
				t.Errorf("unexpected filtered Set; want %v, got %v", tc.expect, filtered)
			}
			if !filtered.IsMutable() {
				t.Error("unexpected filtered Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Filter_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Filter(func(_ int) bool { return true })
	})
}

func Test_SyncHashSet_Filter_Nil(t *testing.T) {
	testCases := map[string]struct {
		filterFunc func(element int) bool
	}{
		"with always-matching predicate": {
			filterFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			filterFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			filtered := set.Filter(tc.filterFunc)
			if filtered == nil {
				t.Error("unexpected nil Set")
			}
			if internal.IsNotNil(filtered) {
				t.Errorf("unexpected filtered Set; want nil, got %#v", filtered)
			}
			if !filtered.IsEmpty() {
				t.Error("unexpected filtered Set emptiness; want true, got false")
			}
			if !filtered.IsMutable() {
				t.Error("unexpected filtered Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Find(t *testing.T) {
	testCases := map[string]struct {
		expectElementIn Set[int]
		expectOK        bool
		searchFunc      func(element int) bool
		set             *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expectElementIn: SyncHash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(_ int) bool { return true },
			set:             SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expectElementIn: SyncHash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return false },
			set:             SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expectElementIn: SyncHash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element > 0 },
			set:             SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expectElementIn: SyncHash(123),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element == 123 },
			set:             SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *SyncHashSet": {
			expectElementIn: SyncHash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element > 0 },
			set:             SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expectElementIn: SyncHash[int](),
			expectOK:        false,
			searchFunc:      func(element int) bool { return element < 0 },
			set:             SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expectElementIn: SyncHash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return true },
			set:             SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expectElementIn: SyncHash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return false },
			set:             SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := tc.set.Find(tc.searchFunc)
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if tc.expectElementIn.IsEmpty() {
				if element != 0 {
					t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
				}
			} else if !tc.expectElementIn.Contains(element) {
				t.Errorf("unexpected element result; want one of %v, got %v", tc.expectElementIn, element)
			}
		})
	}
}

func Test_SyncHashSet_Find_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_, _ = set.Find(func(_ int) bool { return false })
	})
}

func Test_SyncHashSet_Find_Nil(t *testing.T) {
	testCases := map[string]struct {
		searchFunc func(element int) bool
	}{
		"with always-matching predicate": {
			searchFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			searchFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			element, ok := set.Find(tc.searchFunc)
			if ok {
				t.Error("unexpected bool result; want false, got true")
			}
			if element != 0 {
				t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
			}
		})
	}
}

func Test_SyncHashSet_Immutable(t *testing.T) {
	testCases := map[string]struct {
		set *SyncHashSet[int]
	}{
		"on non-empty *SyncHashSet": {
			set: SyncHash(123, 456, 789),
		},
		"on empty *SyncHashSet": {
			set: SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mutable := tc.set.Immutable()
			if internal.IsNil(mutable) {
				t.Error("unexpected nil Set")
			}
			if !mutable.Equal(tc.set) {
				t.Errorf("unexpected Set; want %v, got %v", tc.set, mutable)
			}
			if mutable.IsMutable() {
				t.Error("unexpected Set mutability; want false, got true")
			}
		})
	}
}

func Test_SyncHashSet_Immutable_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Immutable()
	})
}

func Test_SyncHashSet_Immutable_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	immutable := set.Immutable()
	if immutable == nil {
		t.Error("unexpected nil Set")
	}
	if internal.IsNotNil(immutable) {
		t.Errorf("unexpected immutable Set; want nil, got %#v", immutable)
	}
	if !immutable.IsEmpty() {
		t.Error("unexpected immutable Set emptiness; want true, got false")
	}
	if immutable.IsMutable() {
		t.Error("unexpected immutable Set mutability; want false, got true")
	}
}

func TestSyncHashSet_Intersection(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SyncHashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash(-789, -456, -123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *SyncHashSet": {
			expect: Hash(123),
			other:  Hash(-123, 0, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *SyncHashSet": {
			expect: Hash(123, 456),
			other:  Hash(0, 123, 456),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *SyncHashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    SyncHash(123, 456, 789),
		},
		"with empty Set on non-empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty Set on empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with empty Set on empty *SyncHashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			intersection := tc.set.Intersection(tc.other)
			if internal.IsNil(intersection) {
				t.Error("unexpected nil Set")
			}
			if !intersection.Equal(tc.expect) {
				t.Errorf("unexpected intersection Set; want %v, got %v", tc.expect, intersection)
			}
			if !intersection.IsMutable() {
				t.Error("unexpected intersection Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Intersection_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Intersection(Singleton(123))
	})
}

func Test_SyncHashSet_Intersection_Nil(t *testing.T) {
	testCases := map[string]struct {
		other Set[int]
	}{
		"with non-empty Set": {
			other: Hash(123, 456, 789),
		},
		"with empty Set": {
			other: Hash[int](),
		},
		"with nil Set": {
			other: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			intersection := set.Intersection(tc.other)
			if intersection == nil {
				t.Error("unexpected nil Set")
			}
			if internal.IsNotNil(intersection) {
				t.Errorf("unexpected intersection Set; want nil, got %#v", intersection)
			}
			if !intersection.IsEmpty() {
				t.Error("unexpected intersection Set emptiness; want true, got false")
			}
			if !intersection.IsMutable() {
				t.Error("unexpected intersection Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_IsEmpty(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		set    *SyncHashSet[int]
	}{
		"on non-empty *SyncHashSet": {
			expect: false,
			set:    SyncHash(123, 456, 789),
		},
		"on empty *SyncHashSet": {
			expect: true,
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.IsEmpty()
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_IsEmpty_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.IsEmpty()
	})
}

func Test_SyncHashSet_IsEmpty_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	if !set.IsEmpty() {
		t.Error("unexpected result; want true, got false")
	}
}

func Test_SyncHashSet_IsMutable(t *testing.T) {
	testSyncHashSetIsMutable(t, SyncHash[int])
}

func Test_SyncHashSet_IsMutable_Nil(t *testing.T) {
	testSyncHashSetIsMutable(t, func(_ ...int) *SyncHashSet[int] { return nil })
}

func testSyncHashSetIsMutable(t *testing.T, setFunc func(elements ...int) *SyncHashSet[int]) {
	set := setFunc(123, 456, 789)
	if !set.IsMutable() {
		t.Error("unexpected result; want true, got false")
	}
}

func Test_SyncHashSet_Join(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expect: []string{"123", "456", "789"},
			set:    SyncHash(123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expect: []string{"123"},
			set:    SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expect: []string{},
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, tc.set.Join(sep, getIntStringConverterWithDefaultOptions[int]()), sep, tc.expect)
		})
	}
}

func Test_SyncHashSet_Join_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Join(",", getIntStringConverterWithDefaultOptions[int]())
	})
}

func Test_SyncHashSet_Join_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	sep := ","
	assertSetJoin(t, set.Join(sep, getIntStringConverterWithDefaultOptions[int]()), sep, []string{})
}

func Test_SyncHashSet_Len(t *testing.T) {
	testCases := map[string]struct {
		expect int
		set    *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expect: 3,
			set:    SyncHash(123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expect: 1,
			set:    SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expect: 0,
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.Len()
			if result != tc.expect {
				t.Errorf("unexpected length; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Len_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Len()
	})
}

func Test_SyncHashSet_Len_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	if l := set.Len(); l != 0 {
		t.Errorf("unexpected length; want 0, got %v", l)
	}
}

func Test_SyncHashSet_Max(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expectElement: 789,
			expectOK:      true,
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := tc.set.Max(Asc[int])
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if element != tc.expectElement {
				t.Errorf("unexpected element result; want %v, got %v", tc.expectElement, element)
			}
		})
	}
}

func Test_SyncHashSet_Max_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_, _ = set.Max(Asc[int])
	})
}

func Test_SyncHashSet_Max_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	element, ok := set.Max(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_SyncHashSet_Min(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expectElement: -789,
			expectOK:      true,
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := tc.set.Min(Asc[int])
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if element != tc.expectElement {
				t.Errorf("unexpected element result; want %v, got %v", tc.expectElement, element)
			}
		})
	}
}

func Test_SyncHashSet_Min_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_, _ = set.Min(Asc[int])
	})
}

func Test_SyncHashSet_Min_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	element, ok := set.Min(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_SyncHashSet_Mutable(t *testing.T) {
	set := SyncHash(123, 456, 789)
	mutable := set.Mutable()
	if mutable == nil {
		t.Error("unexpected nil MutableSet")
	}
	if mutable != set {
		t.Errorf("unexpected MutableSet; want %v, got %v", set, mutable)
	}
}

func Test_SyncHashSet_Mutable_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Mutable()
	})
}

func Test_SyncHashSet_Mutable_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	mutable := set.Mutable()
	if mutable == nil {
		t.Error("unexpected nil MutableSet")
	}
	if internal.IsNotNil(mutable) {
		t.Errorf("unexpected MutableSet; want nil, got %#v", mutable)
	}
	if !mutable.IsEmpty() {
		t.Error("unexpected MutableSet emptiness; want true, got false")
	}
	if !mutable.IsMutable() {
		t.Error("unexpected MutableSet mutability; want true, got false")
	}
}

func Test_SyncHashSet_None(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some element on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.None(tc.predicateFunc)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_None_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.None(func(_ int) bool { return false })
	})
}

func Test_SyncHashSet_None_Nil(t *testing.T) {
	testCases := map[string]struct {
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			predicateFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			result := set.None(tc.predicateFunc)
			if !result {
				t.Errorf("unexpected match within Set; want true, got %v", result)
			}
		})
	}
}

func Test_SyncHashSet_Put(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with multiple elements on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, -789},
			expect:   Hash(-123, -456, -789, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that all exist on non-empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that some exist on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, 789},
			expect:   Hash(-456, -123, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with single element on non-empty *SyncHashSet": {
			element: -123,
			expect:  Hash(-123, 123, 456, 789),
			set:     SyncHash(123, 456, 789),
		},
		"with single element that exists on non-empty *SyncHashSet": {
			element: 123,
			expect:  Hash(123, 456, 789),
			set:     SyncHash(123, 456, 789),
		},
		"with multiple elements on empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash[int](),
		},
		"with single element on empty *SyncHashSet": {
			element: 123,
			expect:  Hash(123),
			set:     SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.Put(tc.element, tc.elements...)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected Set; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_Put_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Put(0)
	})
}

func Test_SyncHashSet_Put_Nil(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
	}{
		"with multiple elements": {
			element:  123,
			elements: []int{456, 789},
		},
		"with single element": {
			element: 123,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.Put(tc.element, tc.elements...)

			if !set.IsEmpty() {
				t.Error("unexpected Set emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_PutAll(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with Set containing multiple elements on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, -789),
			expect:   Hash(-123, -456, -789, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, 789),
			expect:   Hash(-456, -123, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing single element on non-empty *SyncHashSet": {
			elements: Hash(-123),
			expect:   Hash(-123, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing single element that exists on non-empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing no elements on non-empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with Set containing multiple elements on empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash(123, 456, 789),
			set:      SyncHash[int](),
		},
		"with Set containing single element on empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash(123),
			set:      SyncHash[int](),
		},
		"with Set containing no elements on empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.PutAll(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected Set; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_PutAll_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.PutAll(Singleton(0))
	})
}

func Test_SyncHashSet_PutAll_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
	}{
		"with Set containing multiple elements": {
			elements: Hash(123, 456, 789),
		},
		"with Set containing single element": {
			elements: Hash(123),
		},
		"with Set containing no elements": {
			elements: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.PutAll(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected Set emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_PutSlice(t *testing.T) {
	testCases := map[string]struct {
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with slice containing multiple elements on non-empty *SyncHashSet": {
			elements: []int{-123, -456, -789},
			expect:   Hash(-123, -456, -789, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: []int{-123, -456, 789},
			expect:   Hash(-456, -123, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element on non-empty *SyncHashSet": {
			elements: []int{-123},
			expect:   Hash(-123, 123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that exists on non-empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing no elements on non-empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements on empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash[int](),
		},
		"with slice containing single element on empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash(123),
			set:      SyncHash[int](),
		},
		"with slice containing no elements on empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.PutSlice(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected Set; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_PutSlice_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.PutSlice([]int{0})
	})
}

func Test_SyncHashSet_PutSlice_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements []int
	}{
		"with slice containing multiple elements": {
			elements: []int{123, 456, 789},
		},
		"with slice containing single element": {
			elements: []int{123},
		},
		"with slice containing no elements": {
			elements: []int{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.PutSlice(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected Set emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Range(t *testing.T) {
	testCases := map[string]struct {
		expectCallCount int
		iterFunc        func(element int) bool
		set             *SyncHashSet[int]
	}{
		"with non-breaking iterator on non-empty *SyncHashSet": {
			expectCallCount: 3,
			iterFunc:        func(_ int) bool { return false },
			set:             SyncHash(123, 456, 789),
		},
		"with breaking iterator on non-empty *SyncHashSet": {
			expectCallCount: 3,
			iterFunc: func() func(element int) bool {
				var i int
				return func(_ int) bool {
					i++
					return i == 3
				}
			}(),
			set: SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with non-breaking iterator on empty *SyncHashSet": {
			expectCallCount: 0,
			iterFunc:        func(_ int) bool { return false },
			set:             SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			tc.set.Range(func(element int) bool {
				funcCallCount++
				return tc.iterFunc(element)
			})
			if funcCallCount != tc.expectCallCount {
				t.Errorf("unexpected number of calls to iterator; want %v, got %v", tc.expectCallCount, funcCallCount)
			}
		})
	}
}

func Test_SyncHashSet_Range_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		set.Range(func(_ int) bool { return false })
	})
}

func Test_SyncHashSet_Range_Nil(t *testing.T) {
	var funcCallCount int
	var set *SyncHashSet[int]
	set.Range(func(_ int) bool {
		funcCallCount++
		return false
	})
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to iterator; want 0, got %v", funcCallCount)
	}
}

func Test_SyncHashSet_Retain(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with multiple elements that do not exist on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, -789},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that all exist on non-empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with multiple elements that some exist on non-empty *SyncHashSet": {
			element:  -123,
			elements: []int{-456, 789},
			expect:   Hash(789),
			set:      SyncHash(123, 456, 789),
		},
		"with single element that does not exist on non-empty *SyncHashSet": {
			element: -123,
			expect:  Hash[int](),
			set:     SyncHash(123, 456, 789),
		},
		"with single element that exists on non-empty *SyncHashSet": {
			element: 123,
			expect:  Hash(123),
			set:     SyncHash(123, 456, 789),
		},
		"with multiple elements on empty *SyncHashSet": {
			element:  123,
			elements: []int{456, 789},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with single element on empty *SyncHashSet": {
			element: 123,
			expect:  Hash[int](),
			set:     SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.Retain(tc.element, tc.elements...)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_Retain_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Retain(123)
	})
}

func Test_SyncHashSet_Retain_Nil(t *testing.T) {
	testCases := map[string]struct {
		element  int
		elements []int
	}{
		"with multiple elements": {
			element:  123,
			elements: []int{456, 789},
		},
		"with single element": {
			element: 123,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.Retain(tc.element, tc.elements...)

			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_RetainAll(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with slice containing multiple elements that do not exist on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, -789),
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: Hash(-123, -456, 789),
			expect:   Hash(789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that does not exist on non-empty *SyncHashSet": {
			elements: Hash(-123),
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that exists on non-empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash(123),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing no elements on non-empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements on empty *SyncHashSet": {
			elements: Hash(123, 456, 789),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing single element on empty *SyncHashSet": {
			elements: Hash(123),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing no elements on empty *SyncHashSet": {
			elements: Hash[int](),
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.RetainAll(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_RetainAll_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.RetainAll(Singleton(123))
	})
}

func Test_SyncHashSet_RetainAll_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements Set[int]
	}{
		"with slice containing multiple elements": {
			elements: Hash(123, 456, 789),
		},
		"with slice containing single element": {
			elements: Hash(123),
		},
		"with slice containing no elements": {
			elements: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.RetainAll(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_RetainSlice(t *testing.T) {
	testCases := map[string]struct {
		elements []int
		expect   Set[int]
		set      *SyncHashSet[int]
	}{
		"with slice containing multiple elements that do not exist on non-empty *SyncHashSet": {
			elements: []int{-123, -456, -789},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that all exist on non-empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash(123, 456, 789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements that some exist on non-empty *SyncHashSet": {
			elements: []int{-123, -456, 789},
			expect:   Hash(789),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that does not exist on non-empty *SyncHashSet": {
			elements: []int{-123},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing single element that exists on non-empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash(123),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing no elements on non-empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash[int](),
			set:      SyncHash(123, 456, 789),
		},
		"with slice containing multiple elements on empty *SyncHashSet": {
			elements: []int{123, 456, 789},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing single element on empty *SyncHashSet": {
			elements: []int{123},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
		"with slice containing no elements on empty *SyncHashSet": {
			elements: []int{},
			expect:   Hash[int](),
			set:      SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.RetainSlice(tc.elements)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_RetainSlice_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.RetainSlice([]int{0})
	})
}

func Test_SyncHashSet_RetainSlice_Nil(t *testing.T) {
	testCases := map[string]struct {
		elements []int
	}{
		"with slice containing multiple elements": {
			elements: []int{123, 456, 789},
		},
		"with slice containing single element": {
			elements: []int{123},
		},
		"with slice containing no elements": {
			elements: []int{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.RetainSlice(tc.elements)

			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_RetainWhere(t *testing.T) {
	testCases := map[string]struct {
		expect        Set[int]
		predicateFunc func(element int) bool
		set           *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:        Hash(123, 456, 789),
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:        Hash(123, 456, 789),
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:        Hash(123),
			predicateFunc: func(element int) bool { return element == 123 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *SyncHashSet": {
			expect:        Hash(123, 456, 789),
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(element int) bool { return element < 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:        Hash[int](),
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ret := tc.set.RetainWhere(tc.predicateFunc)

			if internal.IsNil(ret) {
				t.Error("unexpected nil MutableSet")
			}
			if tc.set != ret {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, ret)
			}
			if !tc.expect.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.expect, tc.set)
			}
		})
	}
}

func Test_SyncHashSet_RetainWhere_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.RetainWhere(func(_ int) bool { return true })
	})
}

func Test_SyncHashSet_RetainWhere_Nil(t *testing.T) {
	testCases := map[string]struct {
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			predicateFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			set.RetainWhere(tc.predicateFunc)
			if !set.IsEmpty() {
				t.Error("unexpected MutableSet emptiness; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Slice(t *testing.T) {
	testCases := map[string]struct {
		expect []int
		set    *SyncHashSet[int]
	}{
		"on non-empty *SyncHashSet": {
			expect: []int{123, 456, 789},
			set:    SyncHash(123, 456, 789),
		},
		"on empty *SyncHashSet": {
			expect: []int{},
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			elements := tc.set.Slice()
			if elements == nil {
				t.Error("unexpected nil slice")
			}
			opts := []cmp.Option{cmpopts.SortSlices(Asc[int])}
			if !cmp.Equal(tc.expect, elements, opts...) {
				t.Errorf("unexpected slice; got diff %v", cmp.Diff(tc.expect, elements, opts...))
			}
		})
	}
}

func Test_SyncHashSet_Slice_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Slice()
	})
}

func Test_SyncHashSet_Slice_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	elements := set.Slice()
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_SyncHashSet_Some(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *SyncHashSet[int]
	}{
		"with always-matching predicate on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           SyncHash(123, 456, 789),
		},
		"with conditional predicate matching some element on non-empty *SyncHashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           SyncHash(123, 456, 789),
		},
		"with always-matching predicate on empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           SyncHash[int](),
		},
		"with never-matching predicate on empty *SyncHashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.Some(tc.predicateFunc)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_Some_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Some(func(_ int) bool { return false })
	})
}

func Test_SyncHashSet_Some_Nil(t *testing.T) {
	testCases := map[string]struct {
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			predicateFunc: func(_ int) bool { return false },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			result := set.Some(tc.predicateFunc)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
		})
	}
}

func Test_SyncHashSet_SortedJoin(t *testing.T) {
	testCases := map[string]struct {
		expect string
		set    *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expect: "-789,-456,-123,0,123,456,789",
			set:    SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expect: "123",
			set:    SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expect: "",
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SyncHashSet_SortedJoin_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	})
}

func Test_SyncHashSet_SortedJoin_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	result := set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_SyncHashSet_SortedSlice(t *testing.T) {
	testCases := map[string]struct {
		expect []int
		set    *SyncHashSet[int]
	}{
		"on non-empty *SyncHashSet": {
			expect: []int{123, 456, 789},
			set:    SyncHash(123, 456, 789),
		},
		"on empty *SyncHashSet": {
			expect: []int{},
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			elements := tc.set.SortedSlice(Asc[int])
			if elements == nil {
				t.Error("unexpected nil slice")
			}
			if !cmp.Equal(tc.expect, elements) {
				t.Errorf("unexpected slice; got diff %v", cmp.Diff(tc.expect, elements))
			}
		})
	}
}

func Test_SyncHashSet_SortedSlice_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.SortedSlice(Asc[int])
	})
}

func Test_SyncHashSet_SortedSlice_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	elements := set.SortedSlice(Asc[int])
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_SyncHashSet_TryRange(t *testing.T) {
	testError := errors.New("test")
	testCases := map[string]struct {
		expectCallCount int
		expectError     error
		iterFunc        func(element int) error
		set             *SyncHashSet[int]
	}{
		"with non-failing iterator on non-empty *SyncHashSet": {
			expectCallCount: 3,
			iterFunc:        func(_ int) error { return nil },
			set:             SyncHash(123, 456, 789),
		},
		"with failing iterator on non-empty *SyncHashSet": {
			expectCallCount: 3,
			expectError:     testError,
			iterFunc: func() func(element int) error {
				var i int
				return func(_ int) error {
					i++
					if i == 3 {
						return testError
					}
					return nil
				}
			}(),
			set: SyncHash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with non-failing iterator on empty *SyncHashSet": {
			expectCallCount: 0,
			iterFunc:        func(_ int) error { return nil },
			set:             SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			err := tc.set.TryRange(func(element int) error {
				funcCallCount++
				return tc.iterFunc(element)
			})
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
			} else if tc.expectError != nil {
				t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
			}
			if funcCallCount != tc.expectCallCount {
				t.Errorf("unexpected number of calls to iterator; want %v, got %v", tc.expectCallCount, funcCallCount)
			}
		})
	}
}

func Test_SyncHashSet_TryRange_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.TryRange(func(_ int) error { return nil })
	})
}

func Test_SyncHashSet_TryRange_Nil(t *testing.T) {
	var funcCallCount int
	var set *SyncHashSet[int]
	err := set.TryRange(func(_ int) error {
		funcCallCount++
		return errors.New("test")
	})
	if err != nil {
		t.Errorf("unexpected error; want nil, got %q", err)
	}
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to iterator; want 0, got %v", funcCallCount)
	}
}

func Test_SyncHashSet_Union(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SyncHashSet[int]
	}{
		"with nil Set on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  nil,
			set:    SyncHash(123, 456, 789),
		},
		"with nil *SyncHashSet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  (*SyncHashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *EmptySet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  (*EmptySet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *HashSet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  (*HashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *MutableHashSet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  (*MutableHashSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with nil *SingletonSet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  (*SingletonSet[int])(nil),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  SyncHash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  SyncHash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  SyncHash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  SyncHash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different elements on non-empty *SyncHashSet": {
			expect: SyncHash(12, 34, 56, 123, 456, 789),
			other:  SyncHash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *EmptySet on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  Empty[int](),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing only same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  Hash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  Hash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  Hash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  Hash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different elements on non-empty *SyncHashSet": {
			expect: SyncHash(12, 34, 56, 123, 456, 789),
			other:  Hash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  MutableHash(789, 456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  MutableHash(789, 456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  MutableHash(456, 123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements and others on non-empty *SyncHashSet": {
			expect: SyncHash(0, 123, 456, 789),
			other:  MutableHash(456, 123, 0),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different elements on non-empty *SyncHashSet": {
			expect: SyncHash(12, 34, 56, 123, 456, 789),
			other:  MutableHash(12, 34, 56),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing same element on non-empty *SyncHashSet": {
			expect: SyncHash(123),
			other:  Singleton(123),
			set:    SyncHash(123),
		},
		"with non-nil *SingletonSet containing same element but not others on non-empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  Singleton(123),
			set:    SyncHash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing different element on non-empty *SyncHashSet": {
			expect: SyncHash(12, 123),
			other:  Singleton(12),
			set:    SyncHash(123),
		},
		"with nil Set on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  nil,
			set:    SyncHash[int](),
		},
		"with nil *SyncHashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  (*SyncHashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *EmptySet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  (*EmptySet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *HashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  (*HashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *MutableHashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  (*MutableHashSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with nil *SingletonSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  (*SingletonSet[int])(nil),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet on empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  SyncHash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *SyncHashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  SyncHash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil *EmptySet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  Empty[int](),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *HashSet on empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *HashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  Hash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil non-empty *MutableHashSet on empty *SyncHashSet": {
			expect: SyncHash(123, 456, 789),
			other:  MutableHash(123, 456, 789),
			set:    SyncHash[int](),
		},
		"with non-nil empty *MutableHashSet on empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  MutableHash[int](),
			set:    SyncHash[int](),
		},
		"with non-nil *SingletonSet on empty *SyncHashSet": {
			expect: SyncHash(123),
			other:  Singleton(123),
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			union := tc.set.Union(tc.other)
			if internal.IsNil(union) {
				t.Error("unexpected nil Set")
			}
			if !union.Equal(tc.expect) {
				t.Errorf("unexpected union Set; want %v, got %v", tc.expect, union)
			}
			if !union.IsMutable() {
				t.Error("unexpected union Set mutability; want true, got false")
			}
		})
	}
}

func Test_SyncHashSet_Union_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.Union(Singleton(0))
	})
}

func Test_SyncHashSet_Union_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with nil Set": {
			expect: nil,
			other:  nil,
		},
		"with nil *SyncHashSet": {
			expect: nil,
			other:  (*SyncHashSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: nil,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			expect: nil,
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: nil,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			expect: nil,
			other:  (*SingletonSet[int])(nil),
		},
		"with non-nil empty *SyncHashSet": {
			expect: SyncHash[int](),
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: SyncHash(0),
			other:  SyncHash(0),
		},
		"with non-nil *EmptySet": {
			expect: SyncHash[int](),
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: SyncHash[int](),
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: SyncHash(0),
			other:  Hash(0),
		},
		"with non-nil empty *MutableHashSet": {
			expect: SyncHash[int](),
			other:  MutableHash[int](),
		},
		"with non-nil non-empty *MutableHashSet": {
			expect: SyncHash(0),
			other:  MutableHash(0),
		},
		"with non-nil *SingletonSet": {
			expect: SyncHash(0),
			other:  Singleton(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SyncHashSet[int]
			union := set.Union(tc.other)
			if tc.expect == nil {
				if internal.IsNotNil(union) {
					t.Errorf("unexpected Set; want nil, got %v", union)
				}
			} else {
				if internal.IsNil(union) {
					t.Errorf("unexpected Set; want %v, got nil", tc.expect)
				}
				if !union.Equal(tc.expect) {
					t.Errorf("unexpected union Set; want %v, got %v", tc.expect, union)
				}
				if !union.IsMutable() {
					t.Error("unexpected union Set mutability; want true, got false")
				}
			}
		})
	}
}

func Test_SyncHashSet_String(t *testing.T) {
	set := SyncHash(123, 456, 789)
	assertSetString(t, set.String(), []string{"123", "456", "789"})
}

func Test_SyncHashSet_String_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.String()
	})
}

func Test_SyncHashSet_String_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	assertSetString(t, set.String(), []string{})
}

func Test_SyncHashSet_MarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    *SyncHashSet[int]
	}{
		"on *SyncHashSet containing multiple elements": {
			expect: []string{"123", "456", "789"},
			set:    SyncHash(123, 456, 789),
		},
		"on *SyncHashSet containing single element": {
			expect: []string{"123"},
			set:    SyncHash(123),
		},
		"on *SyncHashSet containing no elements": {
			expect: []string{},
			set:    SyncHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(tc.set)
			if err != nil {
				t.Fatalf("unexpected error; want nil, got %q", err)
			}
			assertSetJSON(t, string(data), tc.expect)
		})
	}
}

func Test_SyncHashSet_MarshalJSON_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_, _ = set.MarshalJSON()
	})
}

func Test_SyncHashSet_MarshalJSON_Nil(t *testing.T) {
	var set *SyncHashSet[int]
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("null"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_SyncHashSet_UnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		expectElements []int
		json           string
	}{
		"with JSON string for array containing multiple elements": {
			expectElements: []int{123, 456, 789},
			json:           "[123,456,789]",
		},
		"with JSON string for array containing single element": {
			expectElements: []int{123},
			json:           "[123]",
		},
		"with JSON string for array containing duplicated elements": {
			expectElements: []int{123, 456, 789},
			json:           "[123,456,789,456,123]",
		},
		"with JSON string for array containing null element": {
			expectElements: []int{0},
			json:           "[null]",
		},
		"with JSON string for empty array": {
			expectElements: []int{},
			json:           "[]",
		},
		"with JSON string for null": {
			expectElements: []int{},
			json:           "null",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := &SyncHashSet[int]{}
			err := json.Unmarshal([]byte(tc.json), set)
			if err != nil {
				t.Errorf("unexpected error; want nil, got %q", err)
			}
			opts := []cmp.Option{cmpopts.SortSlices(Asc[int])}
			if actualElements := set.Slice(); !cmp.Equal(tc.expectElements, actualElements, opts...) {
				t.Errorf("unexpected unmarshalled elements; got diff %v", cmp.Diff(tc.expectElements, actualElements, opts...))
			}
		})
	}
}

func Test_SyncHashSet_UnmarshalJSON_Concurrent(t *testing.T) {
	testConcurrently(func(set *SyncHashSet[int], _ int) {
		_ = set.UnmarshalJSON([]byte(`[123, 456, 789]`))
	})
}

func testConcurrently(fn func(set *SyncHashSet[int], i int), n ...int) {
	_n := DefaultTestConcurrency
	if len(n) > 0 {
		_n = n[0]
	}
	var wg sync.WaitGroup
	wg.Add(_n)
	set := SyncHash(123, 456, 789)
	for i := 0; i < _n; i++ {
		go func(i int) {
			defer wg.Done()
			fn(set, i)
		}(i)
	}
	wg.Wait()
}
