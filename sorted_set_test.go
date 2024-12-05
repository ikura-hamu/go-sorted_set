package gosortedset_test

import (
	_ "embed"
	"slices"
	"testing"

	gosortedset "github.com/ikura-hamu/go-sorted_set"
)

//go:embed testdata/input.json
var inputJSONStr []byte

//go:embed testdata/output_py.json
var outputJSONStr string

type testOperation struct {
	Method string `json:"method"`
	Arg    any    `json:"arg"`
}

type output struct {
	Method   string    `json:"method"`
	Arg      any       `json:"arg"`
	Result   any       `json:"result"`
	Contents []float64 `json:"contents"`
}

func cast[T any](t *testing.T, v any) T {
	t.Helper()
	r, ok := v.(T)
	requireTrue(t, ok)
	return r
}

func castSlice[E any, S ~[]E](t *testing.T, v any) S {
	anyArgs, ok := v.([]any)
	requireTrue(t, ok)
	args := make(S, len(anyArgs))
	for i, arg := range anyArgs {
		args[i], ok = arg.(E)
		requireTrue(t, ok)
	}

	return args
}

// func TestSortedSet(t *testing.T) {
// 	var testCases map[string][]testOperation
// 	err := json.Unmarshal(inputJSONStr, &testCases)
// 	requireNoError(t, err)

// 	var expectedOutput map[string][]output
// 	err = json.Unmarshal([]byte(outputJSONStr), &expectedOutput)
// 	requireNoError(t, err)

// 	for name, testCase := range testCases {
// 		t.Run(name, func(t *testing.T) {
// 			ss := gosortedset.New([]float64{})
// 			outputList := make([]output, 0, len(testCase))
// 			for _, op := range testCase {
// 				var res any
// 				switch op.Method {
// 				case "init":
// 					ss = gosortedset.New(castSlice[float64, []float64](t, op.Arg))
// 					res = ss.Buckets()
// 				case "__contains__":
// 					res = ss.Contains(cast[float64](t, op.Arg))
// 				case "add":
// 					res = ss.Add(cast[float64](t, op.Arg))
// 				case "discard":
// 					res = ss.Discard(cast[float64](t, op.Arg))
// 				case "lt":
// 					lt, exist := ss.Lt(cast[float64](t, op.Arg))
// 					if !exist {
// 						res = nil
// 					} else {
// 						res = lt
// 					}
// 				case "le":
// 					le, exist := ss.Le(cast[float64](t, op.Arg))
// 					if !exist {
// 						res = nil
// 					} else {
// 						res = le
// 					}
// 				case "gt":
// 					gt, exist := ss.Gt(cast[float64](t, op.Arg))
// 					if !exist {
// 						res = nil
// 					} else {
// 						res = gt
// 					}
// 				case "ge":
// 					ge, exist := ss.Ge(cast[float64](t, op.Arg))
// 					if !exist {
// 						res = nil
// 					} else {
// 						res = ge
// 					}
// 				case "__getitem__":
// 					item, err := ss.GetItem(int(cast[float64](t, op.Arg)))
// 					if errors.Is(err, gosortedset.ErrIndexOutOfRange) {
// 						res = "index error"
// 					} else if err != nil {
// 						t.Fatalf("unexpected error: %v", err)
// 					} else {
// 						res = item
// 					}
// 				case "pop":
// 					item, err := ss.Pop(int(cast[float64](t, op.Arg)))
// 					if errors.Is(err, gosortedset.ErrIndexOutOfRange) {
// 						res = "index error"
// 					} else if err != nil {
// 						t.Fatalf("unexpected error: %v", err)
// 					} else {
// 						res = item
// 					}
// 				case "index":
// 					res = ss.CountLe(cast[float64](t, op.Arg))
// 				case "index_right":
// 					res = ss.CountLt(cast[float64](t, op.Arg))
// 				default:
// 					t.Fatalf("unknown method %s", op.Method)
// 				}

// 				output := output{
// 					Method:   op.Method,
// 					Arg:      op.Arg,
// 					Result:   res,
// 					Contents: slices.Collect(ss.Values()),
// 				}
// 				outputList = append(outputList, output)
// 			}

// 			expected, ok := expectedOutput[name]
// 			requireTrue(t, ok)

// 			for i, output := range outputList {
// 				expectedResult := expected[i]

// 				switch expectedR := expectedResult.Result.(type) {
// 				case bool:
// 					r, ok := output.Result.(bool)
// 					requireTrue(t, ok)
// 					if r != expectedR {
// 						t.Errorf("expected result %v, got %v", expectedR, r)
// 					}
// 				case []float64:
// 					r, ok := output.Result.([]float64)
// 					requireTrue(t, ok)
// 					if !slices.Equal(r, expectedR) {
// 						t.Errorf("expected result %v, got %v", expectedR, r)
// 					}
// 				}

// 				if !slices.Equal(output.Contents, expectedResult.Contents) {
// 					t.Errorf("expected contents %v, got %v", expectedResult.Contents, output.Contents)
// 				}
// 			}

// 		})
// 	}

// }

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
