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
	"testing"
)

func Test_Hash(t *testing.T) {
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
			set := Hash(tc.elements...)
			if exp, act := len(tc.elements), set.Len(); act != exp {
				t.Errorf("unexpected Set length; want %v, got %v", exp, act)
			}
			if set.IsMutable() {
				t.Error("unexpected Set mutability; want true, got false")
			}
		})
	}
}

func Test_HashFromJSON(t *testing.T) {
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
			set, err := HashFromJSON[int]([]byte(tc.json))
			if err != nil {
				t.Errorf("unexpected error; want nil, got %q", err)
			} else if set == nil {
				t.Error("unexpected nil Set")
			} else {
				if set.IsMutable() {
					t.Error("unexpected Set mutability; want true, got false")
				}

				opts := []cmp.Option{cmpopts.SortSlices(Asc[int])}
				if actualElements := set.Slice(); !cmp.Equal(tc.expectElements, actualElements, opts...) {
					t.Errorf("unexpected unmarshalled elements; got diff %v", cmp.Diff(tc.expectElements, actualElements, opts...))
				}
			}
		})
	}
}

func Test_HashFromSlice(t *testing.T) {
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
			set := HashFromSlice(tc.elements)
			if exp, act := len(tc.elements), set.Len(); act != exp {
				t.Errorf("unexpected Set length; want %v, got %v", exp, act)
			}
			if set.IsMutable() {
				t.Error("unexpected Set mutability; want true, got false")
			}
		})
	}
}

func Test_HashSet_Clone(t *testing.T) {
	set := Hash(123, 456, 789)
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
	if clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want false, got true")
	}
}

func Test_HashSet_Clone_Nil(t *testing.T) {
	var set *HashSet[int]
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
	if clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want false, got true")
	}
}

func Test_HashSet_Contains(t *testing.T) {
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
			set := Hash(123, 456, 789)
			result := set.Contains(tc.element)
			if result != tc.expect {
				t.Errorf("unexpected element contained within Set: %q; want %v, got %v", tc.element, tc.expect, result)
			}
		})
	}
}

func Test_HashSet_Contains_Nil(t *testing.T) {
	testCases := map[string]struct {
		element int
	}{
		"with non-matching zero value for element":     {0},
		"with non-matching non-zero value for element": {1},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *HashSet[int]
			if set.Contains(tc.element) {
				t.Errorf("unexpected element contained within Set: %q; want false, got true", tc.element)
			}
		})
	}
}

func Test_HashSet_Diff(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *HashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(-789, -456, -123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *HashSet": {
			expect: Hash(456, 789),
			other:  Hash(-123, 0, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *HashSet": {
			expect: Hash(789),
			other:  Hash(0, 123, 456),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *HashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    Hash(123, 456, 789),
		},
		"with empty Set on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash[int](),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    Hash[int](),
		},
		"with empty Set on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Hash[int](),
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
			if diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Diff_Nil(t *testing.T) {
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
			var set *HashSet[int]
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
			if diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_DiffSymmetric(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *HashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *HashSet": {
			expect: Hash(-789, -456, -123, 123, 456, 789),
			other:  Hash(-789, -456, -123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *HashSet": {
			expect: Hash(-123, 0, 456, 789),
			other:  Hash(-123, 0, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *HashSet": {
			expect: Hash(0, 789),
			other:  Hash(0, 123, 456),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *HashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    Hash(123, 456, 789),
		},
		"with empty Set on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash[int](),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set on empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    Hash[int](),
		},
		"with empty Set on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Hash[int](),
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
			if diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_DiffSymmetric_Nil(t *testing.T) {
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
			var set *HashSet[int]
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
			if diff.IsMutable() {
				t.Error("unexpected diff Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Equal(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
		set    *HashSet[int]
	}{
		"with nil *HashSet on non-empty *HashSet": {
			expect: false,
			other:  (*HashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *EmptySet on non-empty *HashSet": {
			expect: false,
			other:  (*EmptySet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *MutableHashSet on non-empty *HashSet": {
			expect: false,
			other:  (*MutableHashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *SingletonSet on non-empty *HashSet": {
			expect: false,
			other:  (*SingletonSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *SyncHashSet on non-empty *HashSet": {
			expect: false,
			other:  (*SyncHashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only same elements on non-empty *HashSet": {
			expect: true,
			other:  Hash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing same elements and others on non-empty *HashSet": {
			expect: false,
			other:  Hash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements on non-empty *HashSet": {
			expect: false,
			other:  Hash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements and others on non-empty *HashSet": {
			expect: false,
			other:  Hash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different elements on non-empty *HashSet": {
			expect: false,
			other:  Hash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *EmptySet on non-empty *HashSet": {
			expect: false,
			other:  Empty[int](),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only same elements on non-empty *HashSet": {
			expect: true,
			other:  MutableHash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing same elements and others on non-empty *HashSet": {
			expect: false,
			other:  MutableHash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements on non-empty *HashSet": {
			expect: false,
			other:  MutableHash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements and others on non-empty *HashSet": {
			expect: false,
			other:  MutableHash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different elements on non-empty *HashSet": {
			expect: false,
			other:  MutableHash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing same element on non-empty *HashSet": {
			expect: true,
			other:  Singleton(123),
			set:    Hash(123),
		},
		"with non-nil *SingletonSet containing same element but not others on non-empty *HashSet": {
			expect: false,
			other:  Singleton(123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing different element on non-empty *HashSet": {
			expect: false,
			other:  Singleton(12),
			set:    Hash(123),
		},
		"with non-nil *SyncHashSet containing only same elements on non-empty *HashSet": {
			expect: true,
			other:  SyncHash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing same elements and others on non-empty *HashSet": {
			expect: false,
			other:  SyncHash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements on non-empty *HashSet": {
			expect: false,
			other:  SyncHash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements and others on non-empty *HashSet": {
			expect: false,
			other:  SyncHash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different elements on non-empty *HashSet": {
			expect: false,
			other:  SyncHash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with nil *HashSet on empty *HashSet": {
			expect: true,
			other:  (*HashSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *EmptySet on empty *HashSet": {
			expect: true,
			other:  (*EmptySet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *MutableHashSet on empty *HashSet": {
			expect: true,
			other:  (*MutableHashSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *SingletonSet on empty *HashSet": {
			expect: true,
			other:  (*SingletonSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *SyncHashSet on empty *HashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
			set:    Hash[int](),
		},
		"with non-nil non-empty *HashSet on empty *HashSet": {
			expect: false,
			other:  Hash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *HashSet on empty *HashSet": {
			expect: true,
			other:  Hash[int](),
			set:    Hash[int](),
		},
		"with non-nil *EmptySet on empty *HashSet": {
			expect: true,
			other:  Empty[int](),
			set:    Hash[int](),
		},
		"with non-nil non-empty *MutableHashSet on empty *HashSet": {
			expect: false,
			other:  MutableHash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *MutableHashSet on empty *HashSet": {
			expect: true,
			other:  MutableHash[int](),
			set:    Hash[int](),
		},
		"with non-nil *SingletonSet on empty *HashSet": {
			expect: false,
			other:  Singleton(123),
			set:    Hash[int](),
		},
		"with non-nil non-empty *SyncHashSet on empty *HashSet": {
			expect: false,
			other:  SyncHash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *SyncHashSet on empty *HashSet": {
			expect: true,
			other:  SyncHash[int](),
			set:    Hash[int](),
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

func Test_HashSet_Equal_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
		"with nil *HashSet": {
			expect: true,
			other:  (*HashSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: true,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: true,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			expect: true,
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil empty *HashSet": {
			expect: true,
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: false,
			other:  Hash(0),
		},
		"with non-nil *EmptySet": {
			expect: true,
			other:  Empty[int](),
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
		"with non-nil empty *SyncHashSet": {
			expect: true,
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: false,
			other:  SyncHash(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *HashSet[int]
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_HashSet_Every(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *HashSet[int]
	}{
		"with always-matching predicate on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           Hash(123, 456, 789),
		},
		"with always-matching predicate on empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash[int](),
		},
		"with never-matching predicate on empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash[int](),
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

func Test_HashSet_Every_Nil(t *testing.T) {
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
			var set *HashSet[int]
			result := set.Every(tc.predicateFunc)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
		})
	}
}

func Test_HashSet_Filter(t *testing.T) {
	testCases := map[string]struct {
		expect     Set[int]
		filterFunc func(element int) bool
		set        *HashSet[int]
	}{
		"with always-matching predicate on non-empty *HashSet": {
			expect:     Hash(123, 456, 789),
			filterFunc: func(_ int) bool { return true },
			set:        Hash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *HashSet": {
			expect:     Hash[int](),
			filterFunc: func(_ int) bool { return false },
			set:        Hash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *HashSet": {
			expect:     Hash(123, 456, 789),
			filterFunc: func(element int) bool { return element > 0 },
			set:        Hash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *HashSet": {
			expect:     Hash(123),
			filterFunc: func(element int) bool { return element == 123 },
			set:        Hash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *HashSet": {
			expect:     Hash(123, 456, 789),
			filterFunc: func(element int) bool { return element > 0 },
			set:        Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *HashSet": {
			expect:     Hash[int](),
			filterFunc: func(element int) bool { return element < 0 },
			set:        Hash(123, 456, 789),
		},
		"with always-matching predicate on empty *HashSet": {
			expect:     Hash[int](),
			filterFunc: func(_ int) bool { return true },
			set:        Hash[int](),
		},
		"with never-matching predicate on empty *HashSet": {
			expect:     Hash[int](),
			filterFunc: func(_ int) bool { return false },
			set:        Hash[int](),
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
			if filtered.IsMutable() {
				t.Error("unexpected filtered Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Filter_Nil(t *testing.T) {
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
			var set *HashSet[int]
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
			if filtered.IsMutable() {
				t.Error("unexpected filtered Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Find(t *testing.T) {
	testCases := map[string]struct {
		expectElementIn Set[int]
		expectOK        bool
		searchFunc      func(element int) bool
		set             *HashSet[int]
	}{
		"with always-matching predicate on non-empty *HashSet": {
			expectElementIn: Hash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(_ int) bool { return true },
			set:             Hash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *HashSet": {
			expectElementIn: Hash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return false },
			set:             Hash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *HashSet": {
			expectElementIn: Hash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element > 0 },
			set:             Hash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *HashSet": {
			expectElementIn: Hash(123),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element == 123 },
			set:             Hash(123, 456, 789),
		},
		"with conditional predicate matching some elements on non-empty *HashSet": {
			expectElementIn: Hash(123, 456, 789),
			expectOK:        true,
			searchFunc:      func(element int) bool { return element > 0 },
			set:             Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *HashSet": {
			expectElementIn: Hash[int](),
			expectOK:        false,
			searchFunc:      func(element int) bool { return element < 0 },
			set:             Hash(123, 456, 789),
		},
		"with always-matching predicate on empty *HashSet": {
			expectElementIn: Hash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return true },
			set:             Hash[int](),
		},
		"with never-matching predicate on empty *HashSet": {
			expectElementIn: Hash[int](),
			expectOK:        false,
			searchFunc:      func(_ int) bool { return false },
			set:             Hash[int](),
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

func Test_HashSet_Find_Nil(t *testing.T) {
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
			var set *HashSet[int]
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

func Test_HashSet_Immutable(t *testing.T) {
	set := Hash(123, 456, 789)
	immutable := set.Immutable()
	if immutable == nil {
		t.Error("unexpected nil Set")
	}
	if immutable != set {
		t.Errorf("unexpected immutable Set; want %v, got %v", set, immutable)
	}
}

func Test_HashSet_Immutable_Nil(t *testing.T) {
	var set *HashSet[int]
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

func Test_HashSet_Intersection(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *HashSet[int]
	}{
		"with non-empty Set containing no intersections on non-empty *HashSet": {
			expect: Hash[int](),
			other:  Hash(-789, -456, -123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing single intersection on non-empty *HashSet": {
			expect: Hash(123),
			other:  Hash(-123, 0, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing multiple intersections on non-empty *HashSet": {
			expect: Hash(123, 456),
			other:  Hash(0, 123, 456),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set containing full intersection on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    Hash(123, 456, 789),
		},
		"with empty Set on non-empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Hash(123, 456, 789),
		},
		"with non-empty Set on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
			set:    Hash[int](),
		},
		"with empty Set on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Hash[int](),
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
			if intersection.IsMutable() {
				t.Error("unexpected intersection Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Intersection_Nil(t *testing.T) {
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
			var set *HashSet[int]
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
			if intersection.IsMutable() {
				t.Error("unexpected intersection Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_IsEmpty(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		set    *HashSet[int]
	}{
		"on non-empty *HashSet": {
			expect: false,
			set:    Hash(123, 456, 789),
		},
		"on empty *HashSet": {
			expect: true,
			set:    Hash[int](),
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

func Test_HashSet_IsEmpty_Nil(t *testing.T) {
	var set *HashSet[int]
	if !set.IsEmpty() {
		t.Error("unexpected result; want true, got false")
	}
}

func Test_HashSet_IsMutable(t *testing.T) {
	testHashSetIsMutable(t, Hash[int])
}

func Test_HashSet_IsMutable_Nil(t *testing.T) {
	testHashSetIsMutable(t, func(_ ...int) *HashSet[int] { return nil })
}

func testHashSetIsMutable(t *testing.T, setFunc func(elements ...int) *HashSet[int]) {
	set := setFunc(123, 456, 789)
	if set.IsMutable() {
		t.Error("unexpected result; want false, got true")
	}
}

func Test_HashSet_Join(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expect: []string{"123", "456", "789"},
			set:    Hash(123, 456, 789),
		},
		"on *HashSet containing single element": {
			expect: []string{"123"},
			set:    Hash(123),
		},
		"on *HashSet containing no elements": {
			expect: []string{},
			set:    Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, tc.set.Join(sep, getIntStringConverterWithDefaultOptions[int]()), sep, tc.expect)
		})
	}
}

func Test_HashSet_Join_Nil(t *testing.T) {
	var set *HashSet[int]
	sep := ","
	assertSetJoin(t, set.Join(sep, getIntStringConverterWithDefaultOptions[int]()), sep, []string{})
}

func Test_HashSet_Len(t *testing.T) {
	testCases := map[string]struct {
		expect int
		set    *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expect: 3,
			set:    Hash(123, 456, 789),
		},
		"on *HashSet containing single element": {
			expect: 1,
			set:    Hash(123),
		},
		"on *HashSet containing no elements": {
			expect: 0,
			set:    Hash[int](),
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

func Test_HashSet_Len_Nil(t *testing.T) {
	var set *HashSet[int]
	if l := set.Len(); l != 0 {
		t.Errorf("unexpected length; want 0, got %v", l)
	}
}

func Test_HashSet_Max(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expectElement: 789,
			expectOK:      true,
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *HashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           Hash(123),
		},
		"on *HashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           Hash[int](),
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

func Test_HashSet_Max_Nil(t *testing.T) {
	var set *HashSet[int]
	element, ok := set.Max(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_HashSet_Min(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expectElement: -789,
			expectOK:      true,
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *HashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           Hash(123),
		},
		"on *HashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           Hash[int](),
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

func Test_HashSet_Min_Nil(t *testing.T) {
	var set *HashSet[int]
	element, ok := set.Min(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_HashSet_Mutable(t *testing.T) {
	testCases := map[string]struct {
		set *HashSet[int]
	}{
		"on non-empty *HashSet": {
			set: Hash(123, 456, 789),
		},
		"on empty *HashSet": {
			set: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mutable := tc.set.Mutable()
			if internal.IsNil(mutable) {
				t.Error("unexpected nil MutableSet")
			}
			if !mutable.Equal(tc.set) {
				t.Errorf("unexpected MutableSet; want %v, got %v", tc.set, mutable)
			}
			if !mutable.IsMutable() {
				t.Error("unexpected MutableSet mutability; want true, got false")
			}
		})
	}
}

func Test_HashSet_Mutable_Nil(t *testing.T) {
	var set *HashSet[int]
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

func Test_HashSet_None(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *HashSet[int]
	}{
		"with always-matching predicate on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching some element on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           Hash(123, 456, 789),
		},
		"with always-matching predicate on empty *HashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash[int](),
		},
		"with never-matching predicate on empty *HashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash[int](),
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

func Test_HashSet_None_Nil(t *testing.T) {
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
			var set *HashSet[int]
			result := set.None(tc.predicateFunc)
			if !result {
				t.Errorf("unexpected match within Set; want true, got %v", result)
			}
		})
	}
}

func Test_HashSet_Range(t *testing.T) {
	testCases := map[string]struct {
		expectCallCount int
		iterFunc        func(element int) bool
		set             *HashSet[int]
	}{
		"with non-breaking iterator on non-empty *HashSet": {
			expectCallCount: 3,
			iterFunc:        func(_ int) bool { return false },
			set:             Hash(123, 456, 789),
		},
		"with breaking iterator on non-empty *HashSet": {
			expectCallCount: 3,
			iterFunc: func() func(element int) bool {
				var i int
				return func(_ int) bool {
					i++
					return i == 3
				}
			}(),
			set: Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with non-breaking iterator on empty *HashSet": {
			expectCallCount: 0,
			iterFunc:        func(_ int) bool { return false },
			set:             Hash[int](),
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

func Test_HashSet_Range_Nil(t *testing.T) {
	var funcCallCount int
	var set *HashSet[int]
	set.Range(func(_ int) bool {
		funcCallCount++
		return false
	})
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to iterator; want 0, got %v", funcCallCount)
	}
}

func Test_HashSet_Slice(t *testing.T) {
	testCases := map[string]struct {
		expect []int
		set    *HashSet[int]
	}{
		"on non-empty *HashSet": {
			expect: []int{123, 456, 789},
			set:    Hash(123, 456, 789),
		},
		"on empty *HashSet": {
			expect: []int{},
			set:    Hash[int](),
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

func Test_HashSet_Slice_Nil(t *testing.T) {
	var set *HashSet[int]
	elements := set.Slice()
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_HashSet_Some(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
		set           *HashSet[int]
	}{
		"with always-matching predicate on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash(123, 456, 789),
		},
		"with never-matching predicate on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching all elements on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching single element on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element == 123 },
			set:           Hash(123, 456, 789),
		},
		"with conditional predicate matching some element on non-empty *HashSet": {
			expect:        true,
			predicateFunc: func(element int) bool { return element > 0 },
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with conditional predicate matching no elements on non-empty *HashSet": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
			set:           Hash(123, 456, 789),
		},
		"with always-matching predicate on empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
			set:           Hash[int](),
		},
		"with never-matching predicate on empty *HashSet": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
			set:           Hash[int](),
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

func Test_HashSet_Some_Nil(t *testing.T) {
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
			var set *HashSet[int]
			result := set.Some(tc.predicateFunc)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
		})
	}
}

func Test_HashSet_SortedJoin(t *testing.T) {
	testCases := map[string]struct {
		expect string
		set    *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expect: "-789,-456,-123,0,123,456,789",
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"on *HashSet containing single element": {
			expect: "123",
			set:    Hash(123),
		},
		"on *HashSet containing no elements": {
			expect: "",
			set:    Hash[int](),
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

func Test_HashSet_SortedJoin_Nil(t *testing.T) {
	var set *HashSet[int]
	result := set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_HashSet_SortedSlice(t *testing.T) {
	testCases := map[string]struct {
		expect []int
		set    *HashSet[int]
	}{
		"on non-empty *HashSet": {
			expect: []int{123, 456, 789},
			set:    Hash(123, 456, 789),
		},
		"on empty *HashSet": {
			expect: []int{},
			set:    Hash[int](),
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

func Test_HashSet_SortedSlice_Nil(t *testing.T) {
	var set *HashSet[int]
	elements := set.SortedSlice(Asc[int])
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_HashSet_TryRange(t *testing.T) {
	testError := errors.New("test")
	testCases := map[string]struct {
		expectCallCount int
		expectError     error
		iterFunc        func(element int) error
		set             *HashSet[int]
	}{
		"with non-failing iterator on non-empty *HashSet": {
			expectCallCount: 3,
			iterFunc:        func(_ int) error { return nil },
			set:             Hash(123, 456, 789),
		},
		"with failing iterator on non-empty *HashSet": {
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
			set: Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with non-failing iterator on empty *HashSet": {
			expectCallCount: 0,
			iterFunc:        func(_ int) error { return nil },
			set:             Hash[int](),
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

func Test_HashSet_TryRange_Nil(t *testing.T) {
	var funcCallCount int
	var set *HashSet[int]
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

func Test_HashSet_Union(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *HashSet[int]
	}{
		"with nil Set on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  nil,
			set:    Hash(123, 456, 789),
		},
		"with nil *HashSet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  (*HashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *EmptySet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  (*EmptySet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *MutableHashSet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  (*MutableHashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *SingletonSet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  (*SingletonSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with nil *SyncHashSet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  (*SyncHashSet[int])(nil),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  Hash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing some same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  Hash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different elements on non-empty *HashSet": {
			expect: Hash(12, 34, 56, 123, 456, 789),
			other:  Hash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *EmptySet on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Empty[int](),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  MutableHash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  MutableHash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  MutableHash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing some same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  MutableHash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different elements on non-empty *HashSet": {
			expect: Hash(12, 34, 56, 123, 456, 789),
			other:  MutableHash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing same element on non-empty *HashSet": {
			expect: Hash(123),
			other:  Singleton(123),
			set:    Hash(123),
		},
		"with non-nil *SingletonSet containing same element but not others on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Singleton(123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SingletonSet containing different element on non-empty *HashSet": {
			expect: Hash(12, 123),
			other:  Singleton(12),
			set:    Hash(123),
		},
		"with non-nil *SyncHashSet containing only same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  SyncHash(789, 456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  SyncHash(789, 456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements on non-empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  SyncHash(456, 123),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing some same elements and others on non-empty *HashSet": {
			expect: Hash(0, 123, 456, 789),
			other:  SyncHash(456, 123, 0),
			set:    Hash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different elements on non-empty *HashSet": {
			expect: Hash(12, 34, 56, 123, 456, 789),
			other:  SyncHash(12, 34, 56),
			set:    Hash(123, 456, 789),
		},
		"with nil Set on empty *HashSet": {
			expect: Hash[int](),
			other:  nil,
			set:    Hash[int](),
		},
		"with nil *HashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  (*HashSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *EmptySet on empty *HashSet": {
			expect: Hash[int](),
			other:  (*EmptySet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *MutableHashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  (*MutableHashSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *SingletonSet on empty *HashSet": {
			expect: Hash[int](),
			other:  (*SingletonSet[int])(nil),
			set:    Hash[int](),
		},
		"with nil *SyncHashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  (*SyncHashSet[int])(nil),
			set:    Hash[int](),
		},
		"with non-nil non-empty *HashSet on empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *HashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Hash[int](),
		},
		"with non-nil *EmptySet on empty *HashSet": {
			expect: Hash[int](),
			other:  Empty[int](),
			set:    Hash[int](),
		},
		"with non-nil non-empty *MutableHashSet on empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  MutableHash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *MutableHashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  MutableHash[int](),
			set:    Hash[int](),
		},
		"with non-nil *SingletonSet on empty *HashSet": {
			expect: Hash(123),
			other:  Singleton(123),
			set:    Hash[int](),
		},
		"with non-nil non-empty *SyncHashSet on empty *HashSet": {
			expect: Hash(123, 456, 789),
			other:  SyncHash(123, 456, 789),
			set:    Hash[int](),
		},
		"with non-nil empty *SyncHashSet on empty *HashSet": {
			expect: Hash[int](),
			other:  SyncHash[int](),
			set:    Hash[int](),
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
			if union.IsMutable() {
				t.Error("unexpected union Set mutability; want false, got true")
			}
		})
	}
}

func Test_HashSet_Union_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with nil Set": {
			expect: nil,
			other:  nil,
		},
		"with nil *HashSet": {
			expect: nil,
			other:  (*HashSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: nil,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: nil,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			expect: nil,
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: nil,
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil empty *HashSet": {
			expect: Hash[int](),
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: Hash(0),
			other:  Hash(0),
		},
		"with non-nil *EmptySet": {
			expect: Hash[int](),
			other:  Empty[int](),
		},
		"with non-nil empty *MutableHashSet": {
			expect: Hash[int](),
			other:  MutableHash[int](),
		},
		"with non-nil non-empty *MutableHashSet": {
			expect: Hash(0),
			other:  MutableHash(0),
		},
		"with non-nil *SingletonSet": {
			expect: Hash(0),
			other:  Singleton(0),
		},
		"with non-nil empty *SyncHashSet": {
			expect: Hash[int](),
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: Hash(0),
			other:  SyncHash(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *HashSet[int]
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
				if union.IsMutable() {
					t.Error("unexpected union Set mutability; want false, got true")
				}
			}
		})
	}
}

func Test_HashSet_String(t *testing.T) {
	set := Hash(123, 456, 789)
	assertSetString(t, set.String(), []string{"123", "456", "789"})
}

func Test_HashSet_String_Nil(t *testing.T) {
	var set *HashSet[int]
	assertSetString(t, set.String(), []string{})
}

func Test_HashSet_MarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    *HashSet[int]
	}{
		"on *HashSet containing multiple elements": {
			expect: []string{"123", "456", "789"},
			set:    Hash(123, 456, 789),
		},
		"on *HashSet containing single element": {
			expect: []string{"123"},
			set:    Hash(123),
		},
		"on *HashSet containing no elements": {
			expect: []string{},
			set:    Hash[int](),
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

func Test_HashSet_MarshalJSON_Nil(t *testing.T) {
	var set *HashSet[int]
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("null"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_HashSet_UnmarshalJSON(t *testing.T) {
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
			set := &HashSet[int]{}
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
