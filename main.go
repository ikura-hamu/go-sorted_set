//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	gosortedset "github.com/ikura-hamu/go-sorted_set"
)

var s *bufio.Scanner
var w *bufio.Writer

func input() string {
	s.Scan()
	return s.Text()
}

func inputT[T any]() T {
	sr := strings.NewReader(input())
	var val T
	_, err := fmt.Fscanf(sr, "%v", &val)
	if err != nil {
		panic(err)
	}
	return val
}

func inputSlice[T any](len int) []T {
	sr := strings.NewReader(input())
	res := make([]T, len)
	for i := 0; i < len; i++ {
		fmt.Fscan(sr, &res[i])
	}
	return res
}

func print(a ...interface{}) {
	fmt.Fprintln(w, a...)
}

func main() {
	s = bufio.NewScanner(os.Stdin)
	w = bufio.NewWriter(os.Stdout)
	defer w.Flush()

	LQ := inputSlice[int](2)
	L, Q := LQ[0], LQ[1]
	ss := gosortedset.New([]int{})
	for range Q {
		q := inputSlice[int](2)
		c, x := q[0], q[1]
		switch c {
		case 1:
			ss.Add(x)
		case 2:
			l, lok := ss.Lt(x)
			if !lok {
				l = 0
			}
			r, rok := ss.Gt(x)
			if !rok {
				r = L
			}
			print(r - l)
		}
	}
}
