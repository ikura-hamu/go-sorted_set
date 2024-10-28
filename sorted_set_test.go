package gosortedset_test

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	gosortedset "github.com/ikura-hamu/go-sorted_set"
)

//go:embed testdata/input.txt
var input string

//go:embed testdata/output_py.txt
var outputPy string

var out *strings.Builder

func printAll(t *testing.T, com string, ss *gosortedset.SortedSet[int]) {
	t.Helper()

	s := "all " + com
	for v := range ss.Values() {
		s += fmt.Sprintf(" %d", v)
	}
	fmt.Fprintln(out, s)
}

func commonPrint(t *testing.T, com string, v any) {
	t.Helper()

	fmt.Fprintln(out, com, v)
}

func commonPrint2None(t *testing.T, com string, v1 any, v2 bool) {
	t.Helper()

	if v2 {
		fmt.Fprintln(out, com, v1)
	} else {
		fmt.Fprintln(out, com, "None")
	}
}

func use(t *testing.T, ss *gosortedset.SortedSet[int], c string, x int) {

	switch c {
	case "add":
		ss.Add(x)
		printAll(t, c, ss)
	case "discard":
		ss.Discard(x)
		printAll(t, c, ss)
	case "get":
		v, err := ss.Get(x)
		if errors.Is(err, gosortedset.ErrIndexOutOfRange) {
			commonPrint(t, c, "None")
		} else if err != nil {
			t.Fatal(err)
		} else {
			commonPrint(t, c, v)
		}
	case "pop":
		v, err := ss.Pop(x)
		if errors.Is(err, gosortedset.ErrIndexOutOfRange) {
			commonPrint(t, c, "None")
		} else if err != nil {
			t.Fatal(err)
		} else {
			commonPrint(t, c, v)
		}
		printAll(t, c, ss)
	case "index":
		v := ss.CountLt(x)
		commonPrint(t, c, v)
	case "index_right":
		v := ss.CountLe(x)
		commonPrint(t, c, v)
	case "lt":
		v, ok := ss.Lt(x)
		commonPrint2None(t, c, v, ok)
	case "le":
		v, ok := ss.Le(x)
		commonPrint2None(t, c, v, ok)
	case "gt":
		v, ok := ss.Gt(x)
		commonPrint2None(t, c, v, ok)
	case "ge":
		v, ok := ss.Ge(x)
		commonPrint2None(t, c, v, ok)
	case "contains":
		v := ss.Contains(x)
		commonPrint(t, c, v)
	case "len":
		v := ss.Len()
		commonPrint(t, c, v)
	}
}

func TestXxx(t *testing.T) {
	lines := strings.Split(input, "\n")
	ss := gosortedset.New([]int{})

	out = &strings.Builder{}

	lines = append(lines, lines...)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		var c string
		xStr := "0"
		c = fields[0]
		if len(fields) > 1 {
			xStr = fields[1]
		}
		x, err := strconv.Atoi(xStr)
		if err != nil {
			t.Fatal(err)
		}

		use(t, ss, c, x)
	}

	outputGo := out.String()
	f, err := os.Create("testdata/output_go.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(outputGo)
	if err != nil {
		t.Fatal(err)
	}

	outputPyLines := strings.Split(outputPy, "\n")
	outputGoLines := strings.Split(outputGo, "\n")

	if len(outputPyLines) != len(outputGoLines) {
		t.Errorf("len(outputPyLines) = %d, len(outputGoLines) = %d", len(outputPyLines), len(outputGoLines))
	}
	for i := range outputPyLines {
		if outputPyLines[i] != outputGoLines[i] {
			t.Errorf("outputPyLines[%d] = %s, outputGoLines[%d] = %s", i, outputPyLines[i], i, outputGoLines[i])
		}
	}
}
