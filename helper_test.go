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
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/neocotic/go-sets/internal"
	"golang.org/x/exp/constraints"
	"sort"
	"strings"
	"testing"
)

func Test_Asc(t *testing.T) {
	elements := []int{789, 456, 123, 0, -123, -456, -789}
	expect := []int{-789, -456, -123, 0, 123, 456, 789}

	sort.SliceStable(elements, func(i, j int) bool { return Asc(elements[i], elements[j]) })

	if !cmp.Equal(expect, elements) {
		t.Errorf("unexpected sorted slice; got diff %v", cmp.Diff(expect, elements))
	}
}

func Test_Desc(t *testing.T) {
	elements := []int{-789, -456, -123, 0, 123, 456, 789}
	expect := []int{789, 456, 123, 0, -123, -456, -789}

	sort.SliceStable(elements, func(i, j int) bool { return Desc(elements[i], elements[j]) })

	if !cmp.Equal(expect, elements) {
		t.Errorf("unexpected sorted slice; got diff %v", cmp.Diff(expect, elements))
	}
}

func Test_Diff(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		others []Set[int]
		set    Set[int]
	}{
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing no intersections": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing single intersection": {
			expect: Hash(456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: Hash(789),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing full intersection": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing no intersections": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing single intersection": {
			expect: Hash(456, 789),
			others: []Set[int]{
				Singleton(-789),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing multiple intersections": {
			expect: Hash(789),
			others: []Set[int]{
				Singleton(0),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing full intersection": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil and empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nil Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and non-empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nil Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash[int](),
		},
		"with non-empty *MutableHashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: MutableHash(789),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: MutableHash(123, 456, 789),
		},
		"with empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: MutableHash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: MutableHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := Diff(tc.set, tc.others...)
			if internal.IsNil(diff) {
				t.Error("unexpected nil Set")
			}
			if !diff.Equal(tc.expect) {
				t.Errorf("unexpected diff Set; want %v, got %v", tc.expect, diff)
			}
			if tc.expect.IsMutable() != diff.IsMutable() {
				t.Errorf("unexpected diff Set mutability; want %v, got %v", tc.expect.IsMutable(), diff.IsMutable())
			}
		})
	}
}

func Test_Diff_Nil(t *testing.T) {
	testCases := map[string]struct {
		others []Set[int]
		set    Set[int]
	}{
		"with nil Set": {
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: nil,
		},
		"with nil *HashSet": {
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := Diff(tc.set, tc.others...)
			if internal.IsNotNil(diff) {
				t.Errorf("unexpected Set; want nil, got %v", diff)
			}
		})
	}
}

func Test_DiffSymmetric(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		others []Set[int]
		set    Set[int]
	}{
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing no intersections": {
			expect: Hash(-789, -456, -123, 123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing single intersection": {
			expect: Hash(-789, -123, 0, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: Hash(0, 789),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing full intersection": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing no intersections": {
			expect: Hash(-789, -456, -123, 123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing single intersection": {
			expect: Hash(-789, -123, 0, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing multiple intersections": {
			expect: Hash(0, 789),
			others: []Set[int]{
				Singleton(0),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing full intersection": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil and empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nil Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nil Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash[int](),
		},
		"with non-empty *MutableHashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: MutableHash(0, 789),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: MutableHash(123, 456, 789),
		},
		"with empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: MutableHash(123, 456, 789),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: MutableHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := DiffSymmetric(tc.set, tc.others...)
			if internal.IsNil(diff) {
				t.Error("unexpected nil Set")
			}
			if !diff.Equal(tc.expect) {
				t.Errorf("unexpected diff Set; want %v, got %v", tc.expect, diff)
			}
			if tc.expect.IsMutable() != diff.IsMutable() {
				t.Errorf("unexpected diff Set mutability; want %v, got %v", tc.expect.IsMutable(), diff.IsMutable())
			}
		})
	}
}

func Test_DiffSymmetric_Nil(t *testing.T) {
	testCases := map[string]struct {
		others []Set[int]
		set    Set[int]
	}{
		"with nil Set": {
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: nil,
		},
		"with nil *HashSet": {
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := DiffSymmetric(tc.set, tc.others...)
			if internal.IsNotNil(diff) {
				t.Errorf("unexpected Set; want nil, got %v", diff)
			}
		})
	}
}

func Test_Equal(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		others []Set[int]
		set    Set[int]
	}{
		"with non-empty *HashSet and equal Sets": {
			expect: true,
			others: []Set[int]{
				Hash(123, 456, 789),
				MutableHash(789, 456, 123),
				SyncHash(123, 456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of equal and non-equal Sets": {
			expect: false,
			others: []Set[int]{
				Hash(123, 456, 789),
				MutableHash(-789, -456, -123),
				SyncHash(123, 456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil and empty Sets": {
			expect: false,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    Hash(123, 456, 789),
		},
		"with empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and mix of nil and empty Sets": {
			expect: true,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    Hash[int](),
		},
		"with non-empty *MutableHashSet and equal Sets": {
			expect: true,
			others: []Set[int]{
				Hash(123, 456, 789),
				MutableHash(789, 456, 123),
				SyncHash(123, 456, 789),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and mix of equal and non-equal Sets": {
			expect: false,
			others: []Set[int]{
				Hash(123, 456, 789),
				MutableHash(-789, -456, -123),
				SyncHash(123, 456, 789),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and mix of nil and empty Sets": {
			expect: false,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    MutableHash(123, 456, 789),
		},
		"with empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: MutableHash[int](),
		},
		"with empty *MutableHashSet and mix of nil and empty Sets": {
			expect: true,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: MutableHash[int](),
		},
		"with empty *MutableHashSet and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    MutableHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := Equal(tc.set, tc.others...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_Equal_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect bool
		others []Set[int]
		set    Set[int]
	}{
		"with nil Set and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: nil,
		},
		"with nil Set and mix of nil and empty Sets": {
			expect: true,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: nil,
		},
		"with nil Set and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    nil,
		},
		"with nil *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: false,
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: (*HashSet[int])(nil),
		},
		"with nil *HashSet and mix of nil and empty Sets": {
			expect: true,
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: (*HashSet[int])(nil),
		},
		"with nil *HashSet and nothing else": {
			expect: true,
			others: []Set[int]{},
			set:    (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := Equal(tc.set, tc.others...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_Group(t *testing.T) {
	testCases := map[string]struct {
		expect      map[string]Set[int]
		grouperFunc func(element int) string
		set         Set[int]
	}{
		"with non-empty *HashSet with multi-group grouper": {
			expect: map[string]Set[int]{
				"negative": Hash(-789, -456, -123),
				"positive": Hash(123, 456, 789),
			},
			grouperFunc: func(element int) string {
				if element < 0 {
					return "negative"
				}
				return "positive"
			},
			set: Hash(-789, -456, -123, 123, 456, 789),
		},
		"with non-empty *HashSet with single-group grouper": {
			expect: map[string]Set[int]{
				"positive": Hash(123, 456, 789),
			},
			grouperFunc: func(element int) string { return "positive" },
			set:         Hash(123, 456, 789),
		},
		"with empty *HashSet": {
			expect:      map[string]Set[int]{},
			grouperFunc: func(element int) string { return "" },
			set:         Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			groups := Group(tc.set, tc.grouperFunc)
			if groups == nil {
				t.Error("unexpected nil map")
			}
			opts := []cmp.Option{cmp.Transformer("Set", func(in Set[int]) []int {
				return in.SortedSlice(Asc[int])
			})}
			if !cmp.Equal(groups, tc.expect, opts...) {
				t.Errorf("unexpected map; got diff %v", cmp.Diff(tc.expect, groups, opts...))
			}
		})
	}
}

func Test_Group_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			groups := Group(tc.set, func(element int) string {
				return ""
			})
			if groups != nil {
				t.Errorf("unexpected map; want nil, got %v", groups)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to grouper; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_Intersection(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		others []Set[int]
		set    Set[int]
	}{
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing no intersections": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing single intersection": {
			expect: Hash(123),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: Hash(123, 456),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets containing full intersection": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing no intersections": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(-789),
				Hash(-456, -123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing single intersection": {
			expect: Hash(123),
			others: []Set[int]{
				Singleton(-789),
				Hash(-123, 0, 123),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing multiple intersections": {
			expect: Hash(123, 456),
			others: []Set[int]{
				Singleton(0),
				Hash(123, 456),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and non-empty Sets containing full intersection": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and mix of nil and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nil Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and non-empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Singleton(123),
				Hash(456, 789),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nil Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash[int](),
		},
		"with non-empty *MutableHashSet and mix of nil, empty, and non-empty Sets containing multiple intersections": {
			expect: MutableHash(123, 456),
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: MutableHash(123, 456, 789),
		},
		"with empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: MutableHash[int](),
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: MutableHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			intersection := Intersection(tc.set, tc.others...)
			if internal.IsNil(intersection) {
				t.Error("unexpected nil Set")
			}
			if !intersection.Equal(tc.expect) {
				t.Errorf("unexpected intersection Set; want %v, got %v", tc.expect, intersection)
			}
			if tc.expect.IsMutable() != intersection.IsMutable() {
				t.Errorf("unexpected intersection Set mutability; want %v, got %v", tc.expect.IsMutable(), intersection.IsMutable())
			}
		})
	}
}

func Test_Intersection_Nil(t *testing.T) {
	testCases := map[string]struct {
		others []Set[int]
		set    Set[int]
	}{
		"with nil Set": {
			others: []Set[int]{
				Singleton(0),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(123, 456),
			},
			set: nil,
		},
		"with nil *HashSet": {
			others: []Set[int]{
				Singleton(123),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 789),
			},
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			diff := Intersection(tc.set, tc.others...)
			if internal.IsNotNil(diff) {
				t.Errorf("unexpected Set; want nil, got %v", diff)
			}
		})
	}
}

func Test_JoinBool(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    Set[bool]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"false", "true"},
			set:    Hash(false, true),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"true"},
			set:    Hash(true),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[bool](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinBool(tc.set, sep), sep, tc.expect)
		})
	}
}

func Test_JoinBool_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[bool]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[bool])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinBool(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinComplex64(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinComplexOption
		set    Set[complex64]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"(0+0i)", "(123.321+321.123i)", "(456.654+654.456i)", "(789.987+987.789i)"},
			set: Hash(
				c64(0, 0),
				c64(123.321, 321.123),
				c64(456.654, 654.456),
				c64(789.987, 987.789),
			),
		},
		"with *HashSet containing multiple elements and WithComplexFormat option": {
			expect: []string{
				"(0e+00+0e+00i)",
				"(1.23321e+02+3.21123e+02i)",
				"(4.56654e+02+6.54456e+02i)",
				"(7.89987e+02+9.87789e+02i)",
			},
			opts: []JoinComplexOption{WithComplexFormat('e')},
			set: Hash(
				c64(0, 0),
				c64(123.321, 321.123),
				c64(456.654, 654.456),
				c64(789.987, 987.789),
			),
		},
		"with *HashSet containing multiple elements and WithComplexPrecision option": {
			expect: []string{"(0.0+0.0i)", "(123.3+321.1i)", "(456.7+654.5i)", "(790.0+987.8i)"},
			opts:   []JoinComplexOption{WithComplexPrecision(1)},
			set: Hash(
				c64(0, 0),
				c64(123.321, 321.123),
				c64(456.654, 654.456),
				c64(789.987, 987.789),
			),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"(123.321+321.123i)"},
			set:    Hash(c64(123.321, 321.123)),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[complex64](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinComplex64(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinComplex64_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[complex64]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[complex64])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinComplex64(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinComplex128(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinComplexOption
		set    Set[complex128]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"(0+0i)", "(123.321+321.123i)", "(456.654+654.456i)", "(789.987+987.789i)"},
			set: Hash(
				c128(0, 0),
				c128(123.321, 321.123),
				c128(456.654, 654.456),
				c128(789.987, 987.789),
			),
		},
		"with *HashSet containing multiple elements and WithComplexFormat option": {
			expect: []string{
				"(0e+00+0e+00i)",
				"(1.23321e+02+3.21123e+02i)",
				"(4.56654e+02+6.54456e+02i)",
				"(7.89987e+02+9.87789e+02i)",
			},
			opts: []JoinComplexOption{WithComplexFormat('e')},
			set: Hash(
				c128(0, 0),
				c128(123.321, 321.123),
				c128(456.654, 654.456),
				c128(789.987, 987.789),
			),
		},
		"with *HashSet containing multiple elements and WithComplexPrecision option": {
			expect: []string{"(0.0+0.0i)", "(123.3+321.1i)", "(456.7+654.5i)", "(790.0+987.8i)"},
			opts:   []JoinComplexOption{WithComplexPrecision(1)},
			set: Hash(
				c128(0, 0),
				c128(123.321, 321.123),
				c128(456.654, 654.456),
				c128(789.987, 987.789),
			),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"(123.321+321.123i)"},
			set:    Hash(c128(123.321, 321.123)),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[complex128](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinComplex128(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinComplex128_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[complex128]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[complex128])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinComplex128(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinFloat32(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinFloatOption
		set    Set[float32]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"-789.987", "-456.654", "-123.321", "0", "123.321", "456.654", "789.987"},
			set:    Hash[float32](-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatFormat option": {
			expect: []string{"0e+00", "1.23321e+02", "4.56654e+02", "7.89987e+02"},
			opts:   []JoinFloatOption{WithFloatFormat('e')},
			set:    Hash[float32](0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatPrecision option": {
			expect: []string{"0.0", "123.3", "456.7", "790.0"},
			opts:   []JoinFloatOption{WithFloatPrecision(1)},
			set:    Hash[float32](0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"123.321"},
			set:    Hash[float32](123.321),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[float32](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinFloat32(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinFloat32_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[float32]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[float32])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinFloat32(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinFloat64(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinFloatOption
		set    Set[float64]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"-789.987", "-456.654", "-123.321", "0", "123.321", "456.654", "789.987"},
			set:    Hash(-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatFormat option": {
			expect: []string{"0e+00", "1.23321e+02", "4.56654e+02", "7.89987e+02"},
			opts:   []JoinFloatOption{WithFloatFormat('e')},
			set:    Hash(0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatPrecision option": {
			expect: []string{"0.0", "123.3", "456.7", "790.0"},
			opts:   []JoinFloatOption{WithFloatPrecision(1)},
			set:    Hash(0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"123.321"},
			set:    Hash(123.321),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[float64](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinFloat64(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinFloat64_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[float64]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[float64])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinFloat64(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinInt(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinIntOption
		set    Set[int]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"-789", "-456", "-123", "0", "123", "456", "789"},
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithIntBase option": {
			expect: []string{"0", "1", "1010", "1100100"},
			opts:   []JoinIntOption{WithIntBase(2)},
			set:    Hash(0, 1, 10, 100),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"123"},
			set:    Hash(123),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinInt(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinInt_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinInt(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinRune(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    Set[rune]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"a", "b", "c"},
			set:    Hash('a', 'b', 'c'),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"a"},
			set:    Hash('a'),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[rune](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinRune(tc.set, sep), sep, tc.expect)
		})
	}
}

func Test_JoinRune_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[rune]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[rune])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinRune(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinString(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		set    Set[string]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"foo", "bar"},
			set:    Hash("foo", "bar"),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"foo"},
			set:    Hash("foo"),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[string](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinString(tc.set, sep), sep, tc.expect)
		})
	}
}

func Test_JoinString_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[string]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[string])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinString(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_JoinUint(t *testing.T) {
	testCases := map[string]struct {
		expect []string
		opts   []JoinUintOption
		set    Set[uint]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: []string{"0", "123", "456", "789"},
			set:    Hash[uint](0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithUintBase option": {
			expect: []string{"0", "1", "1010", "1100100"},
			opts:   []JoinUintOption{WithUintBase(2)},
			set:    Hash[uint](0, 1, 10, 100),
		},
		"with *HashSet containing single element and no options": {
			expect: []string{"123"},
			set:    Hash[uint](123),
		},
		"with *HashSet containing no elements and no options": {
			expect: []string{},
			set:    Hash[uint](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sep := ","
			assertSetJoin(t, JoinUint(tc.set, sep, tc.opts...), sep, tc.expect)
		})
	}
}

func Test_JoinUint_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[uint]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[uint])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := JoinUint(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_Map(t *testing.T) {
	testCases := map[string]struct {
		expect     Set[string]
		mapperFunc func(element int) string
		set        Set[int]
	}{
		"with *EmptySet": {
			expect: Empty[string](),
			set:    Empty[int](),
		},
		"with empty *HashSet": {
			expect: Hash[string](),
			set:    Hash[int](),
		},
		"with non-empty *HashSet": {
			expect: Hash("123", "456", "789"),
			set:    Hash(123, 456, 789),
		},
		"with empty *MutableHashSet": {
			expect: MutableHash[string](),
			set:    MutableHash[int](),
		},
		"with non-empty *MutableHashSet": {
			expect: MutableHash("123", "456", "789"),
			set:    MutableHash(123, 456, 789),
		},
		"with *SingletonSet": {
			expect: Singleton("123"),
			set:    Singleton(123),
		},
		"with empty *SyncHashSet": {
			expect: SyncHash[string](),
			set:    SyncHash[int](),
		},
		"with non-empty *SyncHashSet": {
			expect: SyncHash("123", "456", "789"),
			set:    SyncHash(123, 456, 789),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mapped := Map(tc.set, getIntStringConverterWithDefaultOptions[int]())
			if internal.IsNil(mapped) {
				t.Error("unexpected nil Set")
			}
			if !mapped.Equal(tc.expect) {
				t.Errorf("unexpected mapped Set; want %v, got %v", tc.expect, mapped)
			}
		})
	}
}

func Test_Map_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *EmptySet": {
			set: (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			set: (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			set: (*SingletonSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			set: (*SyncHashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			mapped := Map(tc.set, func(element int) string {
				funcCallCount++
				return ""
			})
			if internal.IsNotNil(mapped) {
				t.Errorf("unexpected mapped Set; want nil, got %v", mapped)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to mapper; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_Max(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           Set[int]
	}{
		"with *HashSet containing multiple elements": {
			expectElement: 789,
			expectOK:      true,
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           Hash(123),
		},
		"with *HashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := Max(tc.set)
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if element != tc.expectElement {
				t.Errorf("unexpected element result; want %v, got %v", tc.expectElement, element)
			}
		})
	}
}

func Test_Max_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := Max(tc.set)
			if ok {
				t.Error("unexpected bool result; want false, got true")
			}
			if element != 0 {
				t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
			}
		})
	}
}

func Test_Min(t *testing.T) {
	testCases := map[string]struct {
		expectElement int
		expectOK      bool
		set           Set[int]
	}{
		"with *HashSet containing multiple elements": {
			expectElement: -789,
			expectOK:      true,
			set:           Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing single element": {
			expectElement: 123,
			expectOK:      true,
			set:           Hash(123),
		},
		"with *HashSet containing no elements": {
			expectElement: 0,
			expectOK:      false,
			set:           Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := Min(tc.set)
			if ok != tc.expectOK {
				t.Errorf("unexpected bool result; want %v, got %v", tc.expectOK, ok)
			}
			if element != tc.expectElement {
				t.Errorf("unexpected element result; want %v, got %v", tc.expectElement, element)
			}
		})
	}
}

func Test_Min_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			element, ok := Min(tc.set)
			if ok {
				t.Error("unexpected bool result; want false, got true")
			}
			if element != 0 {
				t.Errorf("unexpected non-zero value for element result; want 0, got %v", element)
			}
		})
	}
}

func Test_Reduce(t *testing.T) {
	testCases := map[string]struct {
		expect      uint
		initValue   []uint
		reducerFunc func(acc uint, element int) uint
		set         Set[int]
	}{
		"with non-empty *HashSet and additive reducer and initial value": {
			expect:    123 + 456 + 789 + 100,
			initValue: []uint{100},
			reducerFunc: func(acc uint, element int) uint {
				return acc + uint(element)
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and additive reducer and no initial value": {
			expect: 123 + 456 + 789,
			reducerFunc: func(acc uint, element int) uint {
				return acc + uint(element)
			},
			set: Hash(123, 456, 789),
		},
		"with empty *HashSet and initial value": {
			expect:    100,
			initValue: []uint{100},
			reducerFunc: func(acc uint, element int) uint {
				return acc + uint(element)
			},
			set: Hash[int](),
		},
		"with empty *HashSet and no initial value": {
			expect: 0,
			reducerFunc: func(acc uint, element int) uint {
				return acc + uint(element)
			},
			set: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := Reduce(tc.set, tc.reducerFunc, tc.initValue...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
		})
	}
}

func Test_Reduce_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect    uint
		initValue []uint
		set       Set[int]
	}{
		"with nil Set and initial value": {
			expect:    100,
			initValue: []uint{100},
			set:       nil,
		},
		"with nil Set and no initial value": {
			expect: 0,
			set:    nil,
		},
		"with nil *HashSet and initial value": {
			expect:    100,
			initValue: []uint{100},
			set:       (*HashSet[int])(nil),
		},
		"with nil *HashSet and no initial value": {
			expect: 0,
			set:    (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			result := Reduce[int, uint](tc.set, func(acc uint, element int) uint {
				funcCallCount++
				return 123
			}, tc.initValue...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to reducer; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_SortedJoinFloat32(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []JoinFloatOption
		set    Set[float32]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "-789.987,-456.654,-123.321,0,123.321,456.654,789.987",
			set:    Hash[float32](-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatFormat option": {
			expect: "0e+00,1.23321e+02,4.56654e+02,7.89987e+02",
			opts:   []JoinFloatOption{WithFloatFormat('e')},
			set:    Hash[float32](0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatPrecision option": {
			expect: "0.0,123.3,456.7,790.0",
			opts:   []JoinFloatOption{WithFloatPrecision(1)},
			set:    Hash[float32](0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSorting option": {
			expect: "789.987,456.654,123.321,0,-123.321,-456.654,-789.987",
			opts:   []JoinFloatOption{WithFloatSorting(Desc[float64])},
			set:    Hash[float32](-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSortingAsc option": {
			expect: "-789.987,-456.654,-123.321,0,123.321,456.654,789.987",
			opts:   []JoinFloatOption{WithFloatSortingAsc()},
			set:    Hash[float32](-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSortingDesc option": {
			expect: "789.987,456.654,123.321,0,-123.321,-456.654,-789.987",
			opts:   []JoinFloatOption{WithFloatSortingDesc()},
			set:    Hash[float32](-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing single element and no options": {
			expect: "123.321",
			set:    Hash[float32](123.321),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[float32](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinFloat32(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinFloat32_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[float32]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[float32])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinFloat32(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedJoinFloat64(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []JoinFloatOption
		set    Set[float64]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "-789.987,-456.654,-123.321,0,123.321,456.654,789.987",
			set:    Hash(-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatFormat option": {
			expect: "0e+00,1.23321e+02,4.56654e+02,7.89987e+02",
			opts:   []JoinFloatOption{WithFloatFormat('e')},
			set:    Hash(0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatPrecision option": {
			expect: "0.0,123.3,456.7,790.0",
			opts:   []JoinFloatOption{WithFloatPrecision(1)},
			set:    Hash(0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSorting option": {
			expect: "789.987,456.654,123.321,0,-123.321,-456.654,-789.987",
			opts:   []JoinFloatOption{WithFloatSorting(Desc[float64])},
			set:    Hash(-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSortingAsc option": {
			expect: "-789.987,-456.654,-123.321,0,123.321,456.654,789.987",
			opts:   []JoinFloatOption{WithFloatSortingAsc()},
			set:    Hash(-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing multiple elements and WithFloatSortingDesc option": {
			expect: "789.987,456.654,123.321,0,-123.321,-456.654,-789.987",
			opts:   []JoinFloatOption{WithFloatSortingDesc()},
			set:    Hash(-789.987, -456.654, -123.321, 0.0, 123.321, 456.654, 789.987),
		},
		"with *HashSet containing single element and no options": {
			expect: "123.321",
			set:    Hash(123.321),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[float64](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinFloat64(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinFloat64_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[float64]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[float64])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinFloat64(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedJoinInt(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []JoinIntOption
		set    Set[int]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "-789,-456,-123,0,123,456,789",
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithIntBase option": {
			expect: "0,1,1010,1100100",
			opts:   []JoinIntOption{WithIntBase(2)},
			set:    Hash(0, 1, 10, 100),
		},
		"with *HashSet containing multiple elements and WithIntSorting option": {
			expect: "789,456,123,0,-123,-456,-789",
			opts:   []JoinIntOption{WithIntSorting(Desc[int64])},
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithIntSortingAsc option": {
			expect: "-789,-456,-123,0,123,456,789",
			opts:   []JoinIntOption{WithIntSortingAsc()},
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithIntSortingDesc option": {
			expect: "789,456,123,0,-123,-456,-789",
			opts:   []JoinIntOption{WithIntSortingDesc()},
			set:    Hash(-789, -456, -123, 0, 123, 456, 789),
		},
		"with *HashSet containing single element and no options": {
			expect: "123",
			set:    Hash(123),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinInt(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinInt_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinInt(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedJoinRune(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []SortedJoinRuneOption
		set    Set[rune]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "a,b,c",
			set:    Hash('a', 'b', 'c'),
		},
		"with *HashSet containing multiple elements and WithRuneSorting option": {
			expect: "c,b,a",
			opts:   []SortedJoinRuneOption{WithRuneSorting(Desc[rune])},
			set:    Hash('a', 'b', 'c'),
		},
		"with *HashSet containing multiple elements and WithRuneSortingAsc option": {
			expect: "a,b,c",
			opts:   []SortedJoinRuneOption{WithRuneSortingAsc()},
			set:    Hash('a', 'b', 'c'),
		},
		"with *HashSet containing multiple elements and WithRuneSortingDesc option": {
			expect: "c,b,a",
			opts:   []SortedJoinRuneOption{WithRuneSortingDesc()},
			set:    Hash('a', 'b', 'c'),
		},
		"with *HashSet containing single element and no options": {
			expect: "a",
			set:    Hash('a'),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[rune](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinRune(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinRune_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[rune]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[rune])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinRune(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedJoinString(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []SortedJoinStringOption
		set    Set[string]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "bar,foo",
			set:    Hash("foo", "bar"),
		},
		"with *HashSet containing multiple elements and WithStringSorting option": {
			expect: "foo,bar",
			opts:   []SortedJoinStringOption{WithStringSorting(Desc[string])},
			set:    Hash("foo", "bar"),
		},
		"with *HashSet containing multiple elements and WithStringSortingAsc option": {
			expect: "bar,foo",
			opts:   []SortedJoinStringOption{WithStringSortingAsc()},
			set:    Hash("foo", "bar"),
		},
		"with *HashSet containing multiple elements and WithStringSortingDesc option": {
			expect: "foo,bar",
			opts:   []SortedJoinStringOption{WithStringSortingDesc()},
			set:    Hash("foo", "bar"),
		},
		"with *HashSet containing single element and no options": {
			expect: "foo",
			set:    Hash("foo"),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[string](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinString(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinString_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[string]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[string])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinString(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedJoinUint(t *testing.T) {
	testCases := map[string]struct {
		expect string
		opts   []JoinUintOption
		set    Set[uint]
	}{
		"with *HashSet containing multiple elements and no options": {
			expect: "0,123,456,789",
			set:    Hash[uint](0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithUintBase option": {
			expect: "0,1,1010,1100100",
			opts:   []JoinUintOption{WithUintBase(2)},
			set:    Hash[uint](0, 1, 10, 100),
		},
		"with *HashSet containing multiple elements and WithUintSorting option": {
			expect: "789,456,123,0",
			opts:   []JoinUintOption{WithUintSorting(Desc[uint64])},
			set:    Hash[uint](0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithUintSortingAsc option": {
			expect: "0,123,456,789",
			opts:   []JoinUintOption{WithUintSortingAsc()},
			set:    Hash[uint](0, 123, 456, 789),
		},
		"with *HashSet containing multiple elements and WithUintSortingDesc option": {
			expect: "789,456,123,0",
			opts:   []JoinUintOption{WithUintSortingDesc()},
			set:    Hash[uint](0, 123, 456, 789),
		},
		"with *HashSet containing single element and no options": {
			expect: "123",
			set:    Hash[uint](123),
		},
		"with *HashSet containing no elements and no options": {
			expect: "",
			set:    Hash[uint](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinUint(tc.set, ",", tc.opts...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %q, got %q", tc.expect, result)
			}
		})
	}
}

func Test_SortedJoinUint_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[uint]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *HashSet": {
			set: (*HashSet[uint])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := SortedJoinUint(tc.set, ",")
			if exp := ""; result != exp {
				t.Errorf("unexpected result; want %q, got %q", exp, result)
			}
		})
	}
}

func Test_SortedSlice(t *testing.T) {
	testCases := map[string]struct {
		expect []int
		less   func(x, y int) bool
		set    Set[int]
	}{
		"with non-empty *HashSet and default (ascending) sorting": {
			expect: []int{123, 456, 789},
			set:    Hash(123, 456, 789),
		},
		"with non-empty *HashSet and custom (descending) sorting": {
			expect: []int{789, 456, 123},
			less:   Desc[int],
			set:    Hash(123, 456, 789),
		},
		"with empty *HashSet": {
			expect: []int{},
			set:    Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			elements := SortedSlice[int](tc.set, wrapLess(tc.less)...)
			if elements == nil {
				t.Error("unexpected nil slice")
			}
			if !cmp.Equal(tc.expect, elements) {
				t.Errorf("unexpected slice; got diff %v", cmp.Diff(tc.expect, elements))
			}
		})
	}
}

func Test_SortedSlice_Nil(t *testing.T) {
	var set *HashSet[int]
	elements := SortedSlice[int](set)
	if elements != nil {
		t.Errorf("unexpected slice; want nil, got %v", elements)
	}
}

func Test_TryMap(t *testing.T) {
	testErr := errors.New("test")
	testCases := map[string]struct {
		expect      Set[string]
		expectError error
		mapperFunc  func(element int) string
		set         Set[int]
	}{
		"with *EmptySet and passing mapper": {
			expect: Empty[string](),
			set:    Empty[int](),
		},
		"with empty *HashSet and passing mapper": {
			expect: Hash[string](),
			set:    Hash[int](),
		},
		"with non-empty *HashSet and passing mapper": {
			expect: Hash("123", "456", "789"),
			set:    Hash(123, 456, 789),
		},
		"with non-empty *HashSet and failing mapper": {
			expectError: testErr,
			set:         Hash(123, 456, 789),
		},
		"with empty *MutableHashSet and passing mapper": {
			expect: MutableHash[string](),
			set:    MutableHash[int](),
		},
		"with non-empty *MutableHashSet and passing mapper": {
			expect: MutableHash("123", "456", "789"),
			set:    MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and failing mapper": {
			expectError: testErr,
			set:         MutableHash(123, 456, 789),
		},
		"with *SingletonSet and passing mapper": {
			expect: Singleton("123"),
			set:    Singleton(123),
		},
		"with *SingletonSet and failing mapper": {
			expectError: testErr,
			set:         Singleton(123),
		},
		"with empty *SyncHashSet and passing mapper": {
			expect: SyncHash[string](),
			set:    SyncHash[int](),
		},
		"with non-empty *SyncHashSet and passing mapper": {
			expect: SyncHash("123", "456", "789"),
			set:    SyncHash(123, 456, 789),
		},
		"with non-empty *SyncHashSet and failing mapper": {
			expectError: testErr,
			set:         SyncHash(123, 456, 789),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mapper := getIntStringConverterWithDefaultOptions[int]()
			mapped, err := TryMap(tc.set, func(element int) (string, error) {
				return mapper(element), tc.expectError
			})
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
				if internal.IsNotNil(mapped) {
					t.Errorf("unexpected mapped Set; want nil, got %v", mapped)
				}
			} else {
				if tc.expectError != nil {
					t.Errorf("unexpected error; want %q, got nil", tc.expectError)
				}
				if internal.IsNil(mapped) {
					t.Error("unexpected nil Set")
				}
				if !mapped.Equal(tc.expect) {
					t.Errorf("unexpected mapped Set; want %v, got %v", tc.expect, mapped)
				}
			}
		})
	}
}

func Test_TryMap_Nil(t *testing.T) {
	testCases := map[string]struct {
		set Set[int]
	}{
		"with nil Set": {
			set: nil,
		},
		"with nil *EmptySet": {
			set: (*EmptySet[int])(nil),
		},
		"with nil *HashSet": {
			set: (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet": {
			set: (*MutableHashSet[int])(nil),
		},
		"with nil *SingletonSet": {
			set: (*SingletonSet[int])(nil),
		},
		"with nil *SyncHashSet": {
			set: (*SyncHashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			mapped, err := TryMap(tc.set, func(element int) (string, error) {
				funcCallCount++
				return "", nil
			})
			if internal.IsNotNil(mapped) {
				t.Errorf("unexpected mapped Set; want nil, got %v", mapped)
			}
			if err != nil {
				t.Errorf("unexpected error; want nil, got %q", err)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to mapper; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_TryReduce(t *testing.T) {
	testErr := errors.New("test")
	testCases := map[string]struct {
		expectResult uint
		expectError  error
		initValue    []uint
		reducerFunc  func(acc uint, element int) (uint, error)
		set          Set[int]
	}{
		"with non-empty *HashSet and additive reducer and initial value": {
			expectResult: 123 + 456 + 789 + 100,
			initValue:    []uint{100},
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc + uint(element), nil
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and additive reducer and no initial value": {
			expectResult: 123 + 456 + 789,
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc + uint(element), nil
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and failing reducer and initial value": {
			expectResult: 100,
			expectError:  testErr,
			initValue:    []uint{100},
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc, testErr
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and failing reducer and no initial value": {
			expectResult: 0,
			expectError:  testErr,
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc, testErr
			},
			set: Hash(123, 456, 789),
		},
		"with empty *HashSet and initial value": {
			expectResult: 100,
			initValue:    []uint{100},
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc + uint(element), nil
			},
			set: Hash[int](),
		},
		"with empty *HashSet and no initial value": {
			expectResult: 0,
			reducerFunc: func(acc uint, element int) (uint, error) {
				return acc + uint(element), nil
			},
			set: Hash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := TryReduce(tc.set, tc.reducerFunc, tc.initValue...)
			if err != nil {
				if tc.expectError == nil {
					t.Errorf("unexpected error; want nil, got %q", err)
				} else if !errors.Is(err, tc.expectError) {
					t.Errorf("unexpected error; want %q, got %q", tc.expectError, err)
				}
				if result != tc.expectResult {
					t.Errorf("unexpected result; want %v, got %v", tc.expectResult, result)
				}
			} else {
				if tc.expectError != nil {
					t.Errorf("unexpected error; want %q, got nil", tc.expectError)
				}
				if result != tc.expectResult {
					t.Errorf("unexpected result; want %v, got %v", tc.expectResult, result)
				}
			}
		})
	}
}

func Test_TryReduce_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect    uint
		initValue []uint
		set       Set[int]
	}{
		"with nil Set and initial value": {
			expect:    100,
			initValue: []uint{100},
			set:       nil,
		},
		"with nil Set and no initial value": {
			expect: 0,
			set:    nil,
		},
		"with nil *HashSet and initial value": {
			expect:    100,
			initValue: []uint{100},
			set:       (*HashSet[int])(nil),
		},
		"with nil *HashSet and no initial value": {
			expect: 0,
			set:    (*HashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var funcCallCount int
			result, err := TryReduce[int, uint](tc.set, func(acc uint, element int) (uint, error) {
				funcCallCount++
				return 123, nil
			}, tc.initValue...)
			if result != tc.expect {
				t.Errorf("unexpected result; want %v, got %v", tc.expect, result)
			}
			if err != nil {
				t.Errorf("unexpected error; want nil, got %q", err)
			}
			if funcCallCount != 0 {
				t.Errorf("unexpected number of calls to reducer; want 0, got %v", funcCallCount)
			}
		})
	}
}

func Test_Union(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		others []Set[int]
		set    Set[int]
	}{
		"with non-empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash(-789, -456, -123, 0, 123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123, 0),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nil Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash(123, 456, 789),
		},
		"with non-empty *HashSet and nothing else": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{},
			set:    Hash(123, 456, 789),
		},
		"with empty *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 123),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nil Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: Hash[int](),
		},
		"with empty *HashSet and nothing else": {
			expect: Hash[int](),
			others: []Set[int]{},
			set:    Hash[int](),
		},
		"with non-empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: MutableHash(-789, -456, -123, 0, 123, 456, 789),
			others: []Set[int]{
				Singleton(-789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(-456, -123, 0),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and empty Sets": {
			expect: MutableHash(123, 456, 789),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and nil Sets": {
			expect: MutableHash(123, 456, 789),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: MutableHash(123, 456, 789),
		},
		"with non-empty *MutableHashSet and nothing else": {
			expect: MutableHash(123, 456, 789),
			others: []Set[int]{},
			set:    MutableHash(123, 456, 789),
		},
		"with empty *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: MutableHash(123, 456, 789),
			others: []Set[int]{
				Singleton(789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 123),
			},
			set: MutableHash[int](),
		},
		"with empty *MutableHashSet and empty Sets": {
			expect: MutableHash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: MutableHash[int](),
		},
		"with empty *MutableHashSet and nil Sets": {
			expect: MutableHash[int](),
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: MutableHash[int](),
		},
		"with empty *MutableHashSet and nothing else": {
			expect: MutableHash[int](),
			others: []Set[int]{},
			set:    MutableHash[int](),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			union := Union(tc.set, tc.others...)
			if internal.IsNil(union) {
				t.Error("unexpected nil Set")
			}
			if !union.Equal(tc.expect) {
				t.Errorf("unexpected union Set; want %v, got %v", tc.expect, union)
			}
			if tc.expect.IsMutable() != union.IsMutable() {
				t.Errorf("unexpected union Set mutability; want %v, got %v", tc.expect.IsMutable(), union.IsMutable())
			}
		})
	}
}

func Test_Union_Nil(t *testing.T) {
	testCases := map[string]struct {
		expect Set[int]
		others []Set[int]
		set    Set[int]
	}{
		"with nil Set and mix of nil, empty, and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 123),
			},
			set: nil,
		},
		"with nil Set and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: nil,
		},
		"with nil Set and nil Sets": {
			expect: nil,
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: nil,
		},
		"with nil Set and nothing else": {
			expect: nil,
			others: []Set[int]{},
			set:    nil,
		},
		"with nil *HashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 123),
			},
			set: (*HashSet[int])(nil),
		},
		"with nil *HashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: (*HashSet[int])(nil),
		},
		"with nil *HashSet and nil Sets": {
			expect: nil,
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: (*HashSet[int])(nil),
		},
		"with nil *HashSet and nothing else": {
			expect: nil,
			others: []Set[int]{},
			set:    (*HashSet[int])(nil),
		},
		"with nil *MutableHashSet and mix of nil, empty, and non-empty Sets": {
			expect: Hash(123, 456, 789),
			others: []Set[int]{
				Singleton(789),
				nil,
				Empty[int](),
				(*HashSet[int])(nil),
				Hash(456, 123),
			},
			set: (*MutableHashSet[int])(nil),
		},
		"with nil *MutableHashSet and empty Sets": {
			expect: Hash[int](),
			others: []Set[int]{
				Empty[int](),
				Hash[int](),
			},
			set: (*MutableHashSet[int])(nil),
		},
		"with nil *MutableHashSet and nil Sets": {
			expect: nil,
			others: []Set[int]{
				nil,
				(*HashSet[int])(nil),
			},
			set: (*MutableHashSet[int])(nil),
		},
		"with nil *MutableHashSet and nothing else": {
			expect: nil,
			others: []Set[int]{},
			set:    (*MutableHashSet[int])(nil),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			union := Union(tc.set, tc.others...)
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
				if tc.expect.IsMutable() != union.IsMutable() {
					t.Errorf("unexpected union Set mutability; want %v, got %v", tc.expect.IsMutable(), union.IsMutable())
				}
			}
		})
	}
}

func assertSetJoin(t *testing.T, result, sep string, expect []string) {
	if len(result) == 0 {
		if len(expect) > 0 {
			t.Errorf("unexpected string elements contained within result; want %v, got []", expect)
		}
	} else if elements := strings.Split(result, sep); !HashFromSlice(elements).Equal(HashFromSlice(expect)) {
		t.Errorf("unexpected string elements contained within result; want %v, got %v", expect, elements)
	}
}

func assertSetJSON(t *testing.T, result string, expect []string) {
	parseable, ok := strings.CutPrefix(result, "[")
	if !ok {
		t.Fatalf("unexpected prefix; want %q, got %q", "[", result[0])
	}
	parseable, ok = strings.CutSuffix(parseable, "]")
	if !ok {
		t.Fatalf("unexpected suffix; want %q, got %q", "]", result[len(result)-1])
	}
	// Will not play nicely if Set generic type is capable of marshalling bytes containing commas
	assertSetJoin(t, parseable, ",", expect)
}

func assertSetString(t *testing.T, result string, expect []string) {
	parseable, ok := strings.CutPrefix(result, "[")
	if !ok {
		t.Fatalf("unexpected prefix; want %q, got %q", "[", result[0])
	}
	parseable, ok = strings.CutSuffix(parseable, "]")
	if !ok {
		t.Fatalf("unexpected suffix; want %q, got %q", "]", result[len(result)-1])
	}
	// Will not play nicely if Set generic type is capable of producing string representations containing spaces
	assertSetJoin(t, parseable, " ", expect)
}

func c64(real, imag float32) complex64 {
	return complex(real, imag)
}

func c128(real, imag float64) complex128 {
	return complex(real, imag)
}

func getIntStringConverterWithDefaultOptions[E constraints.Signed]() func(element E) string {
	return getIntStringConverter[E](applyJoinIntOptions(nil))
}

func wrapLess[E constraints.Ordered](less func(x, y E) bool) []func(x, y E) bool {
	if less == nil {
		return nil
	}
	return []func(x, y E) bool{less}
}
