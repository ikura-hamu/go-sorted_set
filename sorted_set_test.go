package gosortedset_test

import (
	"errors"
	"slices"
	"testing"

	gosortedset "github.com/ikura-hamu/go-sorted_set"
)

func assertEqualSlice[E comparable, S ~[]E](t *testing.T, expected S, actual S) {
	t.Helper()
	if !slices.Equal(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func assertEqualBuckets[E comparable, B ~[][]E](t *testing.T, expected, actual B) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

	for i, b := range actual {
		if !slices.Equal(b, expected[i]) {
			t.Errorf("buckets[%d]: expected %v, got %v", i, expected[i], actual[i])
		}
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial         []int
		expected        []int
		expectedBuckets [][]int
	}{
		"ok": {
			initial:         []int{1, 2, 3, 4, 5},
			expected:        []int{1, 2, 3, 4, 5},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
		"empty": {
			initial:         []int{},
			expected:        []int{},
			expectedBuckets: [][]int{},
		},
		"not sorted": {
			initial:         []int{5, 4, 3, 2, 1},
			expected:        []int{1, 2, 3, 4, 5},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
		"multiple buckets": {
			initial:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			expected:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15, 16, 17}},
		},
		"duplicate": {
			initial:         []int{1, 2, 3, 4, 5, 5, 5, 5, 5},
			expected:        []int{1, 2, 3, 4, 5},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			assertEqualSlice(t, testCase.expected, slices.Collect(ss.Values()))
			assertEqualBuckets(t, testCase.expectedBuckets, ss.Buckets())
		})
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  []int
		expected []int
	}{
		"ok": {
			initial:  []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		"multiple buckets": {
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			actual := make([]int, 0, len(testCase.initial))
			idx := make([]int, 0, len(testCase.initial))
			expectedIdx := make([]int, 0, len(testCase.initial))
			for i := range len(testCase.initial) {
				expectedIdx = append(expectedIdx, i)
			}
			ss.All()(func(i int, v int) bool {
				idx = append(idx, i)
				actual = append(actual, v)
				return true
			})

			assertEqualSlice(t, expectedIdx, idx)
			assertEqualSlice(t, testCase.expected, actual)
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial         []int
		operation       func(ss *gosortedset.SortedSet[int])
		expected        []int
		expectedBuckets [][]int
	}{
		"ok": {
			initial: []int{1, 2, 4, 5},
			operation: func(ss *gosortedset.SortedSet[int]) {
				ss.Add(3)
			},
			expected:        []int{1, 2, 3, 4, 5},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
		"add to empty": {
			initial: []int{},
			operation: func(ss *gosortedset.SortedSet[int]) {
				ss.Add(1)
			},
			expected:        []int{1},
			expectedBuckets: [][]int{{1}},
		},
		"add and split": {
			initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			operation: func(ss *gosortedset.SortedSet[int]) {
				for i := 17; i <= 25; i++ {
					ss.Add(i)
				}
			},
			expected:        []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, {13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}},
		},
		"contains same item": {
			initial: []int{1, 2, 3, 4, 5},
			operation: func(ss *gosortedset.SortedSet[int]) {
				ss.Add(3)
			},
			expected:        []int{1, 2, 3, 4, 5},
			expectedBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			testCase.operation(ss)
			assertEqualSlice(t, testCase.expected, slices.Collect(ss.Values()))

			assertEqualBuckets(t, testCase.expectedBuckets, ss.Buckets())
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  []int
		arg      int
		expected bool
	}{
		"contains": {
			initial:  []int{1, 2, 3, 4, 5},
			arg:      3,
			expected: true,
		},
		"not contains": {
			initial:  []int{1, 2, 3, 4, 5},
			arg:      6,
			expected: false,
		},
		"empty": {
			initial:  []int{},
			arg:      1,
			expected: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			if ss.Contains(testCase.arg) != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, ss.Contains(testCase.arg))
			}
		})
	}
}

func TestDiscard(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		preOperation  func(ss *gosortedset.SortedSet[int])
		arg           int
		expected      bool
		result        []int
		resultBuckets [][]int
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           3,
			expected:      true,
			result:        []int{1, 2, 4, 5},
			resultBuckets: [][]int{{1, 2, 4, 5}},
		},
		"not contains": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           6,
			expected:      false,
			result:        []int{1, 2, 3, 4, 5},
			resultBuckets: [][]int{{1, 2, 3, 4, 5}},
		},
		"empty": {
			initial:       []int{},
			arg:           1,
			expected:      false,
			result:        []int{},
			resultBuckets: [][]int{},
		},
		"set empty": {
			initial:       []int{1},
			arg:           1,
			expected:      true,
			result:        []int{},
			resultBuckets: [][]int{},
		},
		"bucket empty": {
			initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			preOperation: func(ss *gosortedset.SortedSet[int]) {
				for i := 1; i <= 7; i++ {
					ss.Discard(i)
				}
			},
			arg:           8,
			expected:      true,
			result:        []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
			resultBuckets: [][]int{{9, 10, 11, 12, 13, 14, 15, 16, 17}},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			if testCase.preOperation != nil {
				testCase.preOperation(ss)
			}

			if ss.Discard(testCase.arg) != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, ss.Discard(testCase.arg))
			}

			assertEqualSlice(t, testCase.result, slices.Collect(ss.Values()))
			assertEqualBuckets(t, testCase.resultBuckets, ss.Buckets())
		})
	}
}

func TestLt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		arg           int
		expectedValue int
		expectedExist bool
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           3,
			expectedValue: 2,
			expectedExist: true,
		},
		"not contains": {
			initial:       []int{1, 2, 4, 5},
			arg:           3,
			expectedValue: 2,
			expectedExist: true,
		},
		"return false": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           1,
			expectedValue: 0,
			expectedExist: false,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           11,
			expectedValue: 10,
			expectedExist: true,
		},
		"empty": {
			initial:       []int{},
			arg:           1,
			expectedValue: 0,
			expectedExist: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			value, ok := ss.Lt(testCase.arg)
			if value != testCase.expectedValue || ok != testCase.expectedExist {
				t.Errorf("expected %v, %v, got %v, %v", testCase.expectedValue, testCase.expectedExist, value, ok)
			}
		})
	}
}

func TestLe(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		arg           int
		expectedValue int
		expectedExist bool
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           3,
			expectedValue: 3,
			expectedExist: true,
		},
		"not contains": {
			initial:       []int{1, 2, 4, 5},
			arg:           3,
			expectedValue: 2,
			expectedExist: true,
		},
		"return false": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           0,
			expectedValue: 0,
			expectedExist: false,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           9,
			expectedValue: 9,
			expectedExist: true,
		},
		"empty": {
			initial:       []int{},
			arg:           1,
			expectedValue: 0,
			expectedExist: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			value, ok := ss.Le(testCase.arg)
			if value != testCase.expectedValue || ok != testCase.expectedExist {
				t.Errorf("expected %v, %v, got %v, %v", testCase.expectedValue, testCase.expectedExist, value, ok)
			}
		})
	}
}

func TestGt(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		arg           int
		expectedValue int
		expectedExist bool
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           3,
			expectedValue: 4,
			expectedExist: true,
		},
		"not contains": {
			initial:       []int{1, 2, 4, 5},
			arg:           3,
			expectedValue: 4,
			expectedExist: true,
		},
		"return false": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           5,
			expectedValue: 0,
			expectedExist: false,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           9,
			expectedValue: 10,
			expectedExist: true,
		},
		"empty": {
			initial:       []int{},
			arg:           1,
			expectedValue: 0,
			expectedExist: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			value, ok := ss.Gt(testCase.arg)
			if value != testCase.expectedValue || ok != testCase.expectedExist {
				t.Errorf("expected %v, %v, got %v, %v", testCase.expectedValue, testCase.expectedExist, value, ok)
			}
		})
	}
}

func TestGe(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		arg           int
		expectedValue int
		expectedExist bool
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           3,
			expectedValue: 3,
			expectedExist: true,
		},
		"not contains": {
			initial:       []int{1, 2, 4, 5},
			arg:           3,
			expectedValue: 4,
			expectedExist: true,
		},
		"return false": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           6,
			expectedValue: 0,
			expectedExist: false,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           9,
			expectedValue: 9,
			expectedExist: true,
		},
		"empty": {
			initial:       []int{},
			arg:           1,
			expectedValue: 0,
			expectedExist: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			value, ok := ss.Ge(testCase.arg)
			if value != testCase.expectedValue || ok != testCase.expectedExist {
				t.Errorf("expected %v, %v, got %v, %v", testCase.expectedValue, testCase.expectedExist, value, ok)
			}
		})
	}

}

func TestGetItem(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial       []int
		arg           int
		expectedValue int
		expectedError error
	}{
		"ok": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           2,
			expectedValue: 3,
		},
		"index out of range": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           5,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
		"negative index": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           -1,
			expectedValue: 5,
		},
		"negative index out of range": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           -6,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           15,
			expectedValue: 16,
		},
		"multiple buckets negative index": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           -2,
			expectedValue: 16,
		},
		"empty": {
			initial:       []int{},
			arg:           0,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			value, err := ss.GetItem(testCase.arg)
			if testCase.expectedError != nil {
				if !errors.Is(err, testCase.expectedError) {
					t.Errorf("expected error %v, got %v", testCase.expectedError, err)
				}
				return
			}

			if value != testCase.expectedValue {
				t.Errorf("expected %v, got %v", testCase.expectedValue, value)
			}
		})
	}
}

func TestPop(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial         []int
		preOperation    func(ss *gosortedset.SortedSet[int])
		arg             int
		expectedValue   int
		expectedError   error
		expected        []int
		expectedBuckets [][]int
	}{
		"ok": {
			initial:         []int{1, 2, 3, 4, 5},
			arg:             2,
			expectedValue:   3,
			expected:        []int{1, 2, 4, 5},
			expectedBuckets: [][]int{{1, 2, 4, 5}},
		},
		"index out of range": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           5,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
		"negative index": {
			initial:         []int{1, 2, 3, 4, 5},
			arg:             -1,
			expectedValue:   5,
			expected:        []int{1, 2, 3, 4},
			expectedBuckets: [][]int{{1, 2, 3, 4}},
		},
		"negative index out of range": {
			initial:       []int{1, 2, 3, 4, 5},
			arg:           -6,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
		"multiple buckets": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           15,
			expectedValue: 16,
			expected:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17},
			expectedBuckets: [][]int{
				{1, 2, 3, 4, 5, 6, 7, 8},
				{9, 10, 11, 12, 13, 14, 15, 17},
			},
		},
		"multiple buckets negative index": {
			initial:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:           -2,
			expectedValue: 16,
			expected:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17},
			expectedBuckets: [][]int{
				{1, 2, 3, 4, 5, 6, 7, 8},
				{9, 10, 11, 12, 13, 14, 15, 17},
			},
		},
		"empty": {
			initial:       []int{},
			arg:           0,
			expectedError: gosortedset.ErrIndexOutOfRange,
		},
		"bucket become empty": {
			initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			preOperation: func(ss *gosortedset.SortedSet[int]) {
				for range 7 {
					_, _ = ss.Pop(0)
				}
			},
			arg:             0,
			expectedValue:   8,
			expected:        []int{9, 10, 11, 12, 13, 14, 15, 16, 17},
			expectedBuckets: [][]int{{9, 10, 11, 12, 13, 14, 15, 16, 17}},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			if testCase.preOperation != nil {
				testCase.preOperation(ss)
			}

			value, err := ss.Pop(testCase.arg)
			if testCase.expectedError != nil {
				if !errors.Is(err, testCase.expectedError) {
					t.Errorf("expected error %v, got %v", testCase.expectedError, err)
				}
				return
			}

			if value != testCase.expectedValue {
				t.Errorf("expected %v, got %v", testCase.expectedValue, value)
			}

			assertEqualSlice(t, testCase.expected, slices.Collect(ss.Values()))
			assertEqualBuckets(t, testCase.expectedBuckets, ss.Buckets())
		})
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  []int
		arg      int
		expected int
	}{
		"ok": {
			initial:  []int{1, 2, 3, 4, 5},
			arg:      3,
			expected: 2,
		},
		"not contains": {
			initial:  []int{1, 2, 4, 5},
			arg:      3,
			expected: 2,
		},
		"empty": {
			initial:  []int{},
			arg:      1,
			expected: 0,
		},
		"multiple buckets": {
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:      15,
			expected: 14,
		},
		"multiple buckets not contains": {
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17, 18},
			arg:      15,
			expected: 14,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			if ss.Index(testCase.arg) != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, ss.Index(testCase.arg))
			}
		})
	}
}

func TestIndexRight(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initial  []int
		arg      int
		expected int
	}{
		"ok": {
			initial:  []int{1, 2, 3, 4, 5},
			arg:      3,
			expected: 3,
		},
		"not contains": {
			initial:  []int{1, 2, 4, 5},
			arg:      3,
			expected: 2,
		},
		"empty": {
			initial:  []int{},
			arg:      1,
			expected: 0,
		},
		"multiple buckets": {
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			arg:      15,
			expected: 15,
		},
		"multiple buckets not contains": {
			initial:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16, 17, 18},
			arg:      15,
			expected: 14,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ss := gosortedset.New(testCase.initial)
			if ss.IndexRight(testCase.arg) != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, ss.Index(testCase.arg))
			}
		})
	}
}
