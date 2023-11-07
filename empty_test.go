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

func Test_Empty(t *testing.T) {
	set := Empty[int]()
	if !set.IsEmpty() {
		t.Error("unexpected Set emptiness; want true, got false")
	}
	if set.IsMutable() {
		t.Error("unexpected Set mutability; want true, got false")
	}
}

func Test_EmptyFromJSON(t *testing.T) {
	testCases := map[string]struct {
		expectError error
		json        string
	}{
		"with JSON string for empty array": {
			json: "[]",
		},
		"with JSON string for null": {
			json: "null",
		},
		"with JSON string for array containing zero value for element": {
			expectError: ErrJSONElementCount,
			json:        "[0]",
		},
		"with JSON string for array containing non-zero value for element": {
			expectError: ErrJSONElementCount,
			json:        "[1]",
		},
		"with JSON string for array containing null": {
			expectError: ErrJSONElementCount,
			json:        "[null]",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set, err := EmptyFromJSON[int]([]byte(tc.json))
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
				if !set.IsEmpty() {
					t.Error("unexpected Set emptiness; want true, got false")
				}
				if set.IsMutable() {
					t.Error("unexpected Set mutability; want true, got false")
				}
			}
		})
	}
}

func Test_EmptySet_Clone(t *testing.T) {
	set := Empty[int]()
	clone := set.Clone()
	if internal.IsNil(clone) {
		t.Error("unexpected nil Set")
	}
	if !clone.IsEmpty() {
		t.Error("unexpected cloned Set emptiness; want true, got false")
	}
	if clone.IsMutable() {
		t.Error("unexpected cloned Set mutability; want false, got true")
	}
}

func Test_EmptySet_Clone_Nil(t *testing.T) {
	var set *EmptySet[int]
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

func Test_EmptySet_Contains(t *testing.T) {
	testEmptySetContains(t, Empty[int])
}

func Test_EmptySet_Contains_Nil(t *testing.T) {
	testEmptySetContains(t, func() *EmptySet[int] { return nil })
}

func testEmptySetContains(t *testing.T, setFunc func() *EmptySet[int]) {
	testCases := map[string]struct {
		element int
	}{
		"with zero value for element":     {0},
		"with non-zero value for element": {1},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := setFunc()
			if set.Contains(tc.element) {
				t.Errorf("unexpected element contained within Set: %q; want false, got true", tc.element)
			}
		})
	}
}

func Test_EmptySet_Diff(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with non-empty Set": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
		},
		"with empty Set": {
			expect: Hash[int](),
			other:  Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Empty[int]()
			diff := set.Diff(tc.other)
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

func Test_EmptySet_Diff_Nil(t *testing.T) {
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
			var set *EmptySet[int]
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

func Test_EmptySet_DiffSymmetric(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with non-empty Set": {
			expect: Hash(123, 456, 789),
			other:  Hash(123, 456, 789),
		},
		"with empty Set": {
			expect: Hash[int](),
			other:  Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Empty[int]()
			diff := set.DiffSymmetric(tc.other)
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

func Test_EmptySet_DiffSymmetric_Nil(t *testing.T) {
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
			var set *EmptySet[int]
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

func Test_EmptySet_Equal(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
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
		"with nil *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
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
			set := Empty[int]()
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_EmptySet_Equal_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		other  Set[int]
	}{
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
		"with nil *SyncHashSet": {
			expect: true,
			other:  (*SyncHashSet[int])(nil),
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
			var set *EmptySet[int]
			result := set.Equal(tc.other)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_EmptySet_Every(t *testing.T) {
	testEmptySetEvery(t, Empty[int])
}

func Test_EmptySet_Every_Nil(t *testing.T) {
	testEmptySetEvery(t, func() *EmptySet[int] { return nil })
}

func testEmptySetEvery(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
	result := set.Every(func(_ int) bool {
		funcCallCount++
		return true
	})
	if result {
		t.Error("unexpected result; want false, got true")
	}
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
	}
}

func Test_EmptySet_Filter(t *testing.T) {
	var funcCallCount int
	set := Empty[int]()
	filtered := set.Filter(func(_ int) bool {
		funcCallCount++
		return true
	})
	if internal.IsNil(filtered) {
		t.Error("unexpected nil Set")
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
}

func Test_EmptySet_Filter_Nil(t *testing.T) {
	var (
		funcCallCount int
		set           *EmptySet[int]
	)
	filtered := set.Filter(func(_ int) bool {
		funcCallCount++
		return true
	})
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
}

func Test_EmptySet_Find(t *testing.T) {
	testEmptySetFind(t, Empty[int])
}

func Test_EmptySet_Find_Nil(t *testing.T) {
	testEmptySetFind(t, func() *EmptySet[int] { return nil })
}

func testEmptySetFind(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
	element, ok := set.Find(func(_ int) bool {
		funcCallCount++
		return true
	})
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to search; want 0, got %v", funcCallCount)
	}
}

func Test_EmptySet_Immutable(t *testing.T) {
	set := Empty[int]()
	immutable := set.Immutable()
	if immutable == nil {
		t.Error("unexpected nil Set")
	}
	if immutable != set {
		t.Errorf("unexpected immutable Set; want %v, got %v", set, immutable)
	}
}

func Test_EmptySet_Immutable_Nil(t *testing.T) {
	var set *EmptySet[int]
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

func Test_EmptySet_Intersection(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with non-empty Set": {
			expect: Hash[int](),
			other:  Hash(123, 456, 789),
		},
		"with empty Set": {
			expect: Hash[int](),
			other:  Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Empty[int]()
			intersection := set.Intersection(tc.other)
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

func Test_EmptySet_Intersection_Nil(t *testing.T) {
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
			var set *EmptySet[int]
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

func Test_EmptySet_IsEmpty(t *testing.T) {
	testEmptySetIsEmpty(t, Empty[int])
}

func Test_EmptySet_IsEmpty_Nil(t *testing.T) {
	testEmptySetIsEmpty(t, func() *EmptySet[int] { return nil })
}

func testEmptySetIsEmpty(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	if !set.IsEmpty() {
		t.Error("unexpected result; want true, got false")
	}
}

func Test_EmptySet_IsMutable(t *testing.T) {
	testEmptySetIsMutable(t, Empty[int])
}

func Test_EmptySet_IsMutable_Nil(t *testing.T) {
	testEmptySetIsMutable(t, func() *EmptySet[int] { return nil })
}

func testEmptySetIsMutable(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	if set.IsMutable() {
		t.Error("unexpected result; want false, got true")
	}
}

func Test_EmptySet_Join(t *testing.T) {
	testEmptySetJoin(t, Empty[int])
}

func Test_EmptySet_Join_Nil(t *testing.T) {
	testEmptySetJoin(t, func() *EmptySet[int] { return nil })
}

func testEmptySetJoin(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	result := set.Join(",", getIntStringConverterWithDefaultOptions[int]())
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_EmptySet_Len(t *testing.T) {
	testEmptySetLen(t, Empty[int])
}

func Test_EmptySet_Len_Nil(t *testing.T) {
	testEmptySetLen(t, func() *EmptySet[int] { return nil })
}

func testEmptySetLen(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	if l := set.Len(); l != 0 {
		t.Errorf("unexpected length; want 0, got %v", l)
	}
}

func Test_EmptySet_Max(t *testing.T) {
	testEmptySetMax(t, Empty[int])
}

func Test_EmptySet_Max_Nil(t *testing.T) {
	testEmptySetMax(t, func() *EmptySet[int] { return nil })
}

func testEmptySetMax(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	element, ok := set.Max(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_EmptySet_Min(t *testing.T) {
	testEmptySetMin(t, Empty[int])
}

func Test_EmptySet_Min_Nil(t *testing.T) {
	testEmptySetMin(t, func() *EmptySet[int] { return nil })
}

func testEmptySetMin(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	element, ok := set.Min(Asc[int])
	if ok {
		t.Error("unexpected bool result; want false, got true")
	}
	if element != 0 {
		t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
	}
}

func Test_EmptySet_Mutable(t *testing.T) {
	set := Empty[int]()
	mutable := set.Mutable()
	if internal.IsNil(mutable) {
		t.Error("unexpected nil MutableSet")
	}
	if !mutable.IsEmpty() {
		t.Error("unexpected MutableSet emptiness; want true, got false")
	}
	if !mutable.IsMutable() {
		t.Error("unexpected MutableSet mutability; want true, got false")
	}
}

func Test_EmptySet_Mutable_Nil(t *testing.T) {
	var set *EmptySet[int]
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

func Test_EmptySet_None(t *testing.T) {
	testEmptySetNone(t, Empty[int])
}

func Test_EmptySet_None_Nil(t *testing.T) {
	testEmptySetNone(t, func() *EmptySet[int] { return nil })
}

func testEmptySetNone(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
	result := set.None(func(_ int) bool {
		funcCallCount++
		return true
	})
	if !result {
		t.Error("unexpected result; want true, got false")
	}
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
	}
}

func Test_EmptySet_Range(t *testing.T) {
	testEmptySetRange(t, Empty[int])
}

func Test_EmptySet_Range_Nil(t *testing.T) {
	testEmptySetRange(t, func() *EmptySet[int] { return nil })
}

func testEmptySetRange(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
	set.Range(func(_ int) bool {
		funcCallCount++
		return false
	})
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to iterator; want 0, got %v", funcCallCount)
	}
}

func Test_EmptySet_Slice(t *testing.T) {
	set := Empty[int]()
	elements := set.Slice()
	if elements == nil {
		t.Error("unexpected nil slice")
	}
	if exp := []int{}; !cmp.Equal(exp, elements) {
		t.Errorf("unexpected slice; got diff %v", cmp.Diff(exp, elements))
	}
}

func Test_EmptySet_Slice_Nil(t *testing.T) {
	var set *EmptySet[int]
	elements := set.Slice()
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_EmptySet_Some(t *testing.T) {
	testEmptySetSome(t, Empty[int])
}

func Test_EmptySet_Some_Nil(t *testing.T) {
	testEmptySetSome(t, func() *EmptySet[int] { return nil })
}

func testEmptySetSome(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
	result := set.Some(func(_ int) bool {
		funcCallCount++
		return true
	})
	if result {
		t.Error("unexpected result; want false, got true")
	}
	if funcCallCount != 0 {
		t.Errorf("unexpected number of calls to predicate; want 0, got %v", funcCallCount)
	}
}

func Test_EmptySet_SortedJoin(t *testing.T) {
	testEmptySetSortedJoin(t, Empty[int])
}

func Test_EmptySet_SortedJoin_Nil(t *testing.T) {
	testEmptySetSortedJoin(t, func() *EmptySet[int] { return nil })
}

func testEmptySetSortedJoin(t *testing.T, setFunc func() *EmptySet[int]) {
	set := setFunc()
	result := set.SortedJoin(",", getIntStringConverterWithDefaultOptions[int](), Asc[int])
	if exp := ""; result != exp {
		t.Errorf("unexpected result; want %q, got %q", exp, result)
	}
}

func Test_EmptySet_SortedSlice(t *testing.T) {
	set := Empty[int]()
	elements := set.SortedSlice(Asc[int])
	if elements == nil {
		t.Error("unexpected nil slice")
	}
	if exp := []int{}; !cmp.Equal(exp, elements) {
		t.Errorf("unexpected slice; got diff %v", cmp.Diff(exp, elements))
	}
}

func Test_EmptySet_SortedSlice_Nil(t *testing.T) {
	var set *EmptySet[int]
	elements := set.SortedSlice(Asc[int])
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_EmptySet_TryRange(t *testing.T) {
	testEmptySetTryRange(t, Empty[int])
}

func Test_EmptySet_TryRange_Nil(t *testing.T) {
	testEmptySetTryRange(t, func() *EmptySet[int] { return nil })
}

func testEmptySetTryRange(t *testing.T, setFunc func() *EmptySet[int]) {
	var funcCallCount int
	set := setFunc()
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

func Test_EmptySet_Union(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with nil Set": {
			expect: Empty[int](),
			other:  nil,
		},
		"with nil *EmptySet": {
			expect: Empty[int](),
			other:  (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			expect: Empty[int](),
			other:  (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			expect: Empty[int](),
			other:  (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			expect: Empty[int](),
			other:  (*SingletonSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			expect: Empty[int](),
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil *EmptySet": {
			expect: Empty[int](),
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: Empty[int](),
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: Hash(0),
			other:  Hash(0),
		},
		"with non-nil empty *MutableHashSet": {
			expect: Empty[int](),
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
			expect: Empty[int](),
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: Hash(0),
			other:  SyncHash(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := Empty[int]()
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

func Test_EmptySet_Union_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		other  Set[int]
	}{
		"with nil Set": {
			expect: nil,
			other:  nil,
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
		"with nil *SyncHashSet": {
			expect: nil,
			other:  (*SyncHashSet[int])(nil),
		},
		"with non-nil *EmptySet": {
			expect: Empty[int](),
			other:  Empty[int](),
		},
		"with non-nil empty *HashSet": {
			expect: Empty[int](),
			other:  Hash[int](),
		},
		"with non-nil non-empty *HashSet": {
			expect: Hash(0),
			other:  Hash(0),
		},
		"with non-nil empty *MutableHashSet": {
			expect: Empty[int](),
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
			expect: Empty[int](),
			other:  SyncHash[int](),
		},
		"with non-nil non-empty *SyncHashSet": {
			expect: Hash(0),
			other:  SyncHash(0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var set *EmptySet[int]
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

func Test_EmptySet_String(t *testing.T) {
	set := Empty[int]()
	assertSetString(t, set.String(), []string{})
}

func Test_EmptySet_String_Nil(t *testing.T) {
	var set *EmptySet[int]
	assertSetString(t, set.String(), []string{})
}

func Test_EmptySet_MarshalJSON(t *testing.T) {
	set := Empty[int]()
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("[]"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_EmptySet_MarshalJSON_Nil(t *testing.T) {
	var set *EmptySet[int]
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected error; want nil, got %q", err)
	}
	if exp := []byte("null"); !cmp.Equal(exp, data) {
		t.Errorf("unexpected JSON data; got diff %v", cmp.Diff(exp, data))
	}
}

func Test_EmptySet_UnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		expectError error
		json        string
	}{
		"with JSON string for empty array": {
			json: "[]",
		},
		"with JSON string for null": {
			json: "null",
		},
		"with JSON string for array containing zero value for element": {
			expectError: ErrJSONElementCount,
			json:        "[0]",
		},
		"with JSON string for array containing non-zero value for element": {
			expectError: ErrJSONElementCount,
			json:        "[1]",
		},
		"with JSON string for array containing null": {
			expectError: ErrJSONElementCount,
			json:        "[null]",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			set := &EmptySet[int]{}
			err := json.Unmarshal([]byte(tc.json), set)
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
			} else if tc.expectError != nil {
				t.Errorf("unexpected error; want %q, got nil", tc.expectError)
			}
		})
	}
}
