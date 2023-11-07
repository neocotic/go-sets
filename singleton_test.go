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
	"github.com/neocotic/go-sets/internal"
	"testing"
)

func Test_Singleton(t *testing.T) {
	set := Singleton(123)
	if l := set.Len(); l != 1 {
		t.Errorf("unexpected Set length; want 1, got %v", l)
	}
	if set.IsMutable() {
		t.Error("unexpected Set mutability; want true, got false")
	}
}

func Test_SingletonFromJSON(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectError   error
		json          string
	}{
		"with JSON string for array containing zero value for element": {
			expectElement: 0,
			json:          "[0]",
		},
		"with JSON string for array containing non-zero value for element": {
			expectElement: 123,
			json:          "[123]",
		},
		"with JSON string for array containing null element": {
			expectElement: 0,
			json:          "[null]",
		},
		"with JSON string for empty array": {
			expectError: ErrJSONElementCount,
			json:        "[]",
		},
		"with JSON string for array containing multiple elements": {
			expectError: ErrJSONElementCount,
			json:        "[123,456]",
		},
		"with JSON string for null": {
			expectError: ErrJSONElementCount,
			json:        "null",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set, err := SingletonFromJSON[int]([]byte(tc.json))
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
			} else if tc.expectError != nil {
				t.Errorf("unexpected error; want %q, got nil", tc.expectError)
			} else if set == nil {
				t.Error("unexpected nil Set")
			} else {
				if set.IsMutable() {
					t.Error("unexpected Set mutability; want true, got false")
				}
				if set.element != tc.expectElement {
					t.Errorf("unexpected unmarshalled element; want %v, got %v", tc.expectElement, set)
				}
			}
		})
	}
}

func Test_SingletonSet_Clone(t *testing.T) {
	set := Singleton(123)
	clone := set.Clone()
	if internal.IsNil(clone) {
		t.Error("unexpected nil Set")
	}
	if l := clone.Len(); l != 1 {
		t.Errorf("unexpected cloned Set length; want 1, got %v", l)
	}
	if !clone.Equal(set) {
		t.Errorf("unexpected cloned Set; want %v, got %v", set, clone)
	}
	if clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want false, got true")
	}
}

func Test_SingletonSet_Clone_Nil(t *testing.T) {
	var set *SingletonSet[int]
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

func Test_SingletonSet_Contains(t *testing.T) {
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
			set := Singleton(123)
			result := set.Contains(tc.element)
			if result != tc.expect {
				t.Errorf("unexpected element contained within Set: %q; want %v, got %v", tc.element, tc.expect, result)
			}
		})
	}
}

func Test_SingletonSet_Contains_Nil(t *testing.T) {
	testCases := map[string]struct {
		element int
	}{
		"with non-matching zero value for element":     {0},
		"with non-matching non-zero value for element": {1},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *SingletonSet[int]
			if set.Contains(tc.element) {
				t.Errorf("unexpected element contained within Set: %q; want false, got true", tc.element)
			}
		})
	}
}

func Test_SingletonSet_Diff(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SingletonSet[int]
	}{
		"with non-empty Set containing no intersection": {
			expect: Hash(123),
			other:  Hash(-789, -456, -123),
			set:    Singleton(123),
		},
		"with non-empty Set containing single intersection": {
			expect: Hash[int](),
			other:  Hash(-123, 0, 123),
			set:    Singleton(123),
		},
		"with empty Set": {
			expect: Hash(123),
			other:  Hash[int](),
			set:    Singleton(123),
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

func Test_SingletonSet_Diff_Nil(t *testing.T) {
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
			var set *SingletonSet[int]
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

func Test_SingletonSet_DiffSymmetric(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SingletonSet[int]
	}{
		"with non-empty Set containing no intersection": {
			expect: Hash(-789, -456, -123, 123),
			other:  Hash(-789, -456, -123),
			set:    Singleton(123),
		},
		"with non-empty Set containing single intersection": {
			expect: Hash(-123, 0),
			other:  Hash(-123, 0, 123),
			set:    Singleton(123),
		},
		"with empty Set": {
			expect: Hash(123),
			other:  Hash[int](),
			set:    Singleton(123),
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

func Test_SingletonSet_DiffSymmetric_Nil(t *testing.T) {
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
			var set *SingletonSet[int]
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

func Test_SingletonSet_Equal(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
		"with nil *SingletonSet": {
			expect: false,
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: false,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			expect: false,
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: false,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: false,
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil *SingletonSet containing same element": {
			expect: true,
			other:  Singleton(123),
		},
		"with non-nil *SingletonSet containing different element": {
			expect: false,
			other:  Singleton(456),
		},
		"with non-nil *EmptySet": {
			expect: false,
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: false,
			other:  Hash[int](),
		},
		"with non-nil *HashSet containing only same element": {
			expect: true,
			other:  Hash(123),
		},
		"with non-nil *HashSet containing same element as well as different elements": {
			expect: false,
			other:  Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different element": {
			expect: false,
			other:  Hash(456),
		},
		"with non-nil empty *MutableHashSet": {
			expect: false,
			other:  MutableHash[int](),
		},
		"with non-nil *MutableHashSet containing only same element": {
			expect: true,
			other:  MutableHash(123),
		},
		"with non-nil *MutableHashSet containing same element as well as different elements": {
			expect: false,
			other:  MutableHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different element": {
			expect: false,
			other:  MutableHash(456),
		},
		"with non-nil empty *SyncHashSet": {
			expect: false,
			other:  SyncHash[int](),
		},
		"with non-nil *SyncHashSet containing only same element": {
			expect: true,
			other:  SyncHash(123),
		},
		"with non-nil *SyncHashSet containing same element as well as different elements": {
			expect: false,
			other:  SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different element": {
			expect: false,
			other:  SyncHash(456),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Singleton(123)
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SingletonSet_Equal_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
		"with nil *SingletonSet": {
			expect: true,
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: true,
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashHashSet": {
			expect: true,
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: true,
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil *SingletonSet": {
			expect: false,
			other:  Singleton(0),
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
			var set *SingletonSet[int]
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_SingletonSet_Every(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
		},
		"with conditional matching predicate": {
			expect:        true,
			predicateFunc: func(element int) bool { return element == 123 },
		},
		"with conditional non-matching predicate": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			predicate := func(element int) bool {
				funcCalls = append(funcCalls, element)
				return tc.predicateFunc(element)
			}
			set := Singleton(123)
			result := set.Every(predicate)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to predicate; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to predicate; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_Every_Nil(t *testing.T) {
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
			var (
				funcCallCount int
				set           *SingletonSet[int]
			)
			predicate := func(element int) bool {
				funcCallCount++
				return tc.predicateFunc(element)
			}
			result := set.Every(predicate)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SingletonSet_Filter(t *testing.T) {
	testCases := map[string]struct {
		expect     Set[int]
		filterFunc func(element int) bool
	}{
		"with always-matching predicate": {
			expect:     Singleton(123),
			filterFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			expect:     Empty[int](),
			filterFunc: func(_ int) bool { return false },
		},
		"with conditional matching predicate": {
			expect:     Singleton(123),
			filterFunc: func(element int) bool { return element == 123 },
		},
		"with conditional non-matching predicate": {
			expect:     Empty[int](),
			filterFunc: func(element int) bool { return element < 0 },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			filter := func(element int) bool {
				funcCalls = append(funcCalls, element)
				return tc.filterFunc(element)
			}
			set := Singleton(123)
			filtered := set.Filter(filter)
			if internal.IsNil(filtered) {
				t.Error("unexpected nil Set")
			}
			if !filtered.Equal(tc.expect) {
				t.Errorf("unexpected filtered Set; want %v, got %v", tc.expect, filtered)
			}
			if filtered.IsMutable() {
				t.Error("unexpected filtered Set mutability; want false, got true")
			}
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to filter; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to filter; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_Filter_Nil(t *testing.T) {
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
			var funcCallCount int
			filter := func(element int) bool {
				funcCallCount++
				return tc.filterFunc(element)
			}
			var set *SingletonSet[int]
			filtered := set.Filter(filter)
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
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to filter; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SingletonSet_Find(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		searchFunc    func(element int) bool
	}{
		"with always-matching predicate": {
			expectElement: 123,
			expectOK:      true,
			searchFunc:    func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			expectElement: 0,
			expectOK:      false,
			searchFunc:    func(_ int) bool { return false },
		},
		"with conditional matching predicate": {
			expectElement: 123,
			expectOK:      true,
			searchFunc:    func(element int) bool { return element == 123 },
		},
		"with conditional non-matching predicate": {
			expectElement: 0,
			expectOK:      false,
			searchFunc:    func(element int) bool { return element < 0 },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			search := func(element int) bool {
				funcCalls = append(funcCalls, element)
				return tc.searchFunc(element)
			}
			set := Singleton(123)
			element, ok := set.Find(search)
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if element != tc.expectElement {
				t.Errorf("unexpected element result; want %v, got %v", tc.expectElement, element)
			}
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to search; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to search; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_Find_Nil(t *testing.T) {
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
			var funcCallCount int
			search := func(element int) bool {
				funcCallCount++
				return tc.searchFunc(element)
			}
			var set *SingletonSet[int]
			element, ok := set.Find(search)
			if ok {
				t.Error("unexpected bool result; want false, got true")
			}
			if element != 0 {
				t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to search; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SingletonSet_Immutable(t *testing.T) {
	set := Singleton(123)
	immutable := set.Immutable()
	if immutable == nil {
		t.Error("unexpected nil Set")
	}
	if immutable != set {
		t.Errorf("unexpected immutable Set; want %v, got %v", set, immutable)
	}
}

func Test_SingletonSet_Immutable_Nil(t *testing.T) {
	var set *SingletonSet[int]
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

func Test_SingletonSet_Intersection(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
		set    *SingletonSet[int]
	}{
		"with non-empty Set containing no intersection": {
			expect: Hash[int](),
			other:  Hash(-789, -456, -123),
			set:    Singleton(123),
		},
		"with non-empty Set containing intersection": {
			expect: Singleton(123),
			other:  Hash(-123, 0, 123),
			set:    Singleton(123),
		},
		"with empty Set": {
			expect: Hash[int](),
			other:  Hash[int](),
			set:    Singleton(123),
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

func Test_SingletonSet_Intersection_Nil(t *testing.T) {
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
			var set *SingletonSet[int]
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

func Test_SingletonSet_IsEmpty(t *testing.T) {
	set := Singleton(123)
	if set.IsEmpty() {
		t.Error("unexpected result; want false, got true")
	}
}

func Test_SingletonSet_IsEmpty_Nil(t *testing.T) {
	var set *SingletonSet[int]
	if !set.IsEmpty() {
		t.Error("unexpected result; want true, got false")
	}
}

func Test_SingletonSet_IsMutable(t *testing.T) {
	testSingletonSetIsMutable(t, Singleton[int])
}

func Test_SingletonSet_IsMutable_Nil(t *testing.T) {
	testSingletonSetIsMutable(t, func(_ int) *SingletonSet[int] { return nil })
}

func testSingletonSetIsMutable(t *testing.T, setFunc func(element int) *SingletonSet[int]) {
	set := setFunc(123)
	if set.IsMutable() {
		t.Error("unexpected result; want false, got true")
	}
}

func Test_SingletonSet_Join(t *testing.T) {
	set := Singleton(123)
	result := set.Join(",", getIntStringConverterWithDefaultOptions[int]())
	if exp := "123"; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_SingletonSet_Join_Nil(t *testing.T) {
	var set *SingletonSet[int]
	result := set.Join(",", getIntStringConverterWithDefaultOptions[int]())
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_SingletonSet_Len(t *testing.T) {
	set := Singleton(123)
	if l := set.Len(); l != 1 {
		t.Errorf("unexpected length; want 1, got %v", l)
	}
}

func Test_SingletonSet_Len_Nil(t *testing.T) {
	var set *SingletonSet[int]
	if l := set.Len(); l != 0 {
		t.Errorf("unexpected length; want 0, got %v", l)
	}
}

func Test_SingletonSet_Max(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *SingletonSet[int]
	}{
		"with zero value for element": {
			expectElement: 0,
			expectOK:      true,
			set:           Singleton(0),
		},
		"with non-zero value for element": {
			expectElement: 123,
			expectOK:      true,
			set:           Singleton(123),
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

func Test_SingletonSet_Max_Nil(t *testing.T) {
	var set *SingletonSet[int]
	element, ok := set.Max(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_SingletonSet_Min(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           *SingletonSet[int]
	}{
		"with zero value for element": {
			expectElement: 0,
			expectOK:      true,
			set:           Singleton(0),
		},
		"with non-zero value for element": {
			expectElement: 123,
			expectOK:      true,
			set:           Singleton(123),
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

func Test_SingletonSet_Min_Nil(t *testing.T) {
	var set *SingletonSet[int]
	element, ok := set.Min(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_SingletonSet_Mutable(t *testing.T) {
	set := Singleton(123)
	mutable := set.Mutable()
	if internal.IsNil(mutable) {
		t.Error("unexpected nil MutableSet")
	}
	if l := mutable.Len(); l != 1 {
		t.Errorf("unexpected MutableSet length; want 1, got %v", l)
	}
	if !mutable.Equal(set) {
		t.Errorf("unexpected MutableSet; want %v, got %v", set, mutable)
	}
	if !mutable.IsMutable() {
		t.Error("unexpected MutableSet mutability; want true, got false")
	}
}

func Test_SingletonSet_Mutable_Nil(t *testing.T) {
	var set *SingletonSet[int]
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

func Test_SingletonSet_None(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			expect:        false,
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			expect:        true,
			predicateFunc: func(_ int) bool { return false },
		},
		"with conditional matching predicate": {
			expect:        false,
			predicateFunc: func(element int) bool { return element == 123 },
		},
		"with conditional non-matching predicate": {
			expect:        true,
			predicateFunc: func(element int) bool { return element < 0 },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			predicate := func(element int) bool {
				funcCalls = append(funcCalls, element)
				return tc.predicateFunc(element)
			}
			set := Singleton(123)
			result := set.None(predicate)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to predicate; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to predicate; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_None_Nil(t *testing.T) {
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
			var (
				funcCallCount int
				set           *SingletonSet[int]
			)
			predicate := func(element int) bool {
				funcCallCount++
				return tc.predicateFunc(element)
			}
			result := set.None(predicate)
			if !result {
				t.Errorf("unexpected match within Set; want true, got %v", result)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SingletonSet_Range(t *testing.T) {
	var funcCalls []int
	set := Singleton(123)
	set.Range(func(element int) bool {
		funcCalls = append(funcCalls, element)
		return false
	})
	if l := len(funcCalls); l != 1 {
		t.Errorf("unexpected number of calls to iterator; want 1, got %v", l)
	}
	if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
		t.Errorf("unexpected calls to iterator; got diff %v", cmp.Diff(exp, funcCalls))
	}
}

func Test_SingletonSet_Range_Nil(t *testing.T) {
	var funcCallCount int
	var set *SingletonSet[int]
	set.Range(func(_ int) bool {
		funcCallCount++
		return false
	})
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to iterator; want 0, got %v", funcCallCount)
	}
}

func Test_SingletonSet_Slice(t *testing.T) {
	set := Singleton(123)
	elements := set.Slice()
	if elements == nil {
		t.Error("unexpected nil slice")
	}
	if exp := []int{123}; !cmp.Equal(exp, elements) {
		t.Errorf("unexpected slice; got diff %v", cmp.Diff(exp, elements))
	}
}

func Test_SingletonSet_Slice_Nil(t *testing.T) {
	var set *SingletonSet[int]
	elements := set.Slice()
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_SingletonSet_Some(t *testing.T) {
	testCases := map[string]struct {
		expect        bool
		predicateFunc func(element int) bool
	}{
		"with always-matching predicate": {
			expect:        true,
			predicateFunc: func(_ int) bool { return true },
		},
		"with never-matching predicate": {
			expect:        false,
			predicateFunc: func(_ int) bool { return false },
		},
		"with conditional matching predicate": {
			expect:        true,
			predicateFunc: func(element int) bool { return element == 123 },
		},
		"with conditional non-matching predicate": {
			expect:        false,
			predicateFunc: func(element int) bool { return element < 0 },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			predicate := func(element int) bool {
				funcCalls = append(funcCalls, element)
				return tc.predicateFunc(element)
			}
			set := Singleton(123)
			result := set.Some(predicate)
			if result != tc.expect {
				t.Errorf("unexpected match within Set; want %v, got %v", tc.expect, result)
			}
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to predicate; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to predicate; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_Some_Nil(t *testing.T) {
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
			var (
				funcCallCount int
				set           *SingletonSet[int]
			)
			predicate := func(element int) bool {
				funcCallCount++
				return tc.predicateFunc(element)
			}
			result := set.Some(predicate)
			if result {
				t.Errorf("unexpected match within Set; want false, got %v", result)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SingletonSet_SortedJoin(t *testing.T) {
	set := Singleton(123)
	result := set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	if exp := "123"; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_SingletonSet_SortedJoin_Nil(t *testing.T) {
	var set *SingletonSet[int]
	result := set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_SingletonSet_SortedSlice(t *testing.T) {
	set := Singleton(123)
	elements := set.SortedSlice(Asc[int])
	if elements == nil {
		t.Error("unexpected nil slice")
	}
	if exp := []int{123}; !cmp.Equal(exp, elements) {
		t.Errorf("unexpected slice; got diff %v", cmp.Diff(exp, elements))
	}
}

func Test_SingletonSet_SortedSlice_Nil(t *testing.T) {
	var set *SingletonSet[int]
	elements := set.SortedSlice(Asc[int])
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_SingletonSet_TryRange(t *testing.T) {
	testError := errors.New("test")
	testCases := map[string]struct {
		expectError error
		iterFunc    func(element int) error
	}{
		"with non-failing iterator": {
			expectError: nil,
			iterFunc:    func(_ int) error { return nil },
		},
		"with failing iterator": {
			expectError: testError,
			iterFunc:    func(_ int) error { return testError },
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCalls []int
			set := Singleton(123)
			err := set.TryRange(func(element int) error {
				funcCalls = append(funcCalls, element)
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
			if l := len(funcCalls); l != 1 {
				t.Errorf("unexpected number of calls to iterator; want 1, got %v", l)
			}
			if exp := []int{123}; !cmp.Equal(exp, funcCalls) {
				t.Errorf("unexpected calls to iterator; got diff %v", cmp.Diff(exp, funcCalls))
			}
		})
	}
}

func Test_SingletonSet_TryRange_Nil(t *testing.T) {
	var funcCallCount int
	var set *SingletonSet[int]
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

func Test_SingletonSet_Union(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with nil Set": {
			expect: Singleton(123),
			other:  nil,
		},
		"with nil *SingletonSet": {
			expect: Singleton(123),
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *EmptySet": {
			expect: Singleton(123),
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			expect: Singleton(123),
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: Singleton(123),
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: Singleton(123),
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil *SingletonSet containing same element": {
			expect: Singleton(123),
			other:  Singleton(123),
		},
		"with non-nil *SingletonSet containing different element": {
			expect: Hash(123, 456),
			other:  Singleton(456),
		},
		"with non-nil *EmptySet": {
			expect: Singleton(123),
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: Singleton(123),
			other:  Hash[int](),
		},
		"with non-nil *HashSet containing only same element": {
			expect: Singleton(123),
			other:  Hash(123),
		},
		"with non-nil *HashSet containing same element as well as different elements": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
		},
		"with non-nil *HashSet containing only different element": {
			expect: Hash(123, 456),
			other:  Hash(456),
		},
		"with non-nil empty *MutableHashSet": {
			expect: Singleton(123),
			other:  MutableHash[int](),
		},
		"with non-nil *MutableHashSet containing only same element": {
			expect: Singleton(123),
			other:  MutableHash(123),
		},
		"with non-nil *MutableHashSet containing same element as well as different elements": {
			expect: Hash(123, 456, 789),
			other:  MutableHash(123, 456, 789),
		},
		"with non-nil *MutableHashSet containing only different element": {
			expect: Hash(123, 456),
			other:  MutableHash(456),
		},
		"with non-nil empty *SyncHashSet": {
			expect: Singleton(123),
			other:  SyncHash[int](),
		},
		"with non-nil *SyncHashSet containing only same element": {
			expect: Singleton(123),
			other:  SyncHash(123),
		},
		"with non-nil *SyncHashSet containing same element as well as different elements": {
			expect: Hash(123, 456, 789),
			other:  SyncHash(123, 456, 789),
		},
		"with non-nil *SyncHashSet containing only different element": {
			expect: Hash(123, 456),
			other:  SyncHash(456),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Singleton(123)
			union := set.Union(tc.other)
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

func Test_SingletonSet_Union_Nil(t *testing.T) {
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
			expect: Singleton(0),
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
			expect: Singleton(0),
			other:  MutableHash(0),
		},
		"with non-nil *SingletonSet": {
			expect: Singleton(0),
			other:  Singleton(0),
		},
		"with non-nil empty *SyncHashSet": {
			expect: Hash[int](),
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: Singleton(0),
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

func Test_SingletonSet_String(t *testing.T) {
	set := Singleton(123)
	assertSetString(t, set.String(), []string{"123"})
}

func Test_SingletonSet_String_Nil(t *testing.T) {
	var set *SingletonSet[int]
	assertSetString(t, set.String(), []string{})
}

func Test_SingletonSet_MarshalJSON(t *testing.T) {
	set := Singleton(123)
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("[123]"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_SingletonSet_MarshalJSON_Nil(t *testing.T) {
	var set *SingletonSet[int]
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("null"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_SingletonSet_UnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectError   error
		json          string
	}{
		"with JSON string for array containing zero value for element": {
			expectElement: 0,
			json:          "[0]",
		},
		"with JSON string for array containing non-zero value for element": {
			expectElement: 123,
			json:          "[123]",
		},
		"with JSON string for array containing null element": {
			expectElement: 0,
			json:          "[null]",
		},
		"with JSON string for empty array": {
			expectError: ErrJSONElementCount,
			json:        "[]",
		},
		"with JSON string for array containing multiple elements": {
			expectError: ErrJSONElementCount,
			json:        "[123,456]",
		},
		"with JSON string for null": {
			expectError: ErrJSONElementCount,
			json:        "null",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := &SingletonSet[int]{}
			err := json.Unmarshal([]byte(tc.json), set)
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
				if set.element != 0 {
					t.Errorf("unexpected unmarshalled element; want 0, got %v", set)
				}
			} else {
				if tc.expectError != nil {
					t.Errorf("unexpected error; want %q, got nil", tc.expectError)
				}
				if set.element != tc.expectElement {
					t.Errorf("unexpected unmarshalled element; want %v, got %v", tc.expectElement, set)
				}
			}
		})
	}
}
