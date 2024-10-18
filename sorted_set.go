package gosortedset

import (
	"cmp"
	"fmt"
	"iter"
	"math"
	"slices"
	"strings"
)

const (
	bucketRatio = 16
	splitRatio  = 24
)

type SortedSet[T cmp.Ordered] struct {
	a    [][]T
	size int
}

func New[T cmp.Ordered](a []T) *SortedSet[T] {
	s := &SortedSet[T]{}
	if !slices.IsSorted(a) {
		slices.Sort(a)
	}
	a = slices.Compact(a)
	n := len(a)

	s.size = n

	numBucket := int(math.Ceil(math.Sqrt(float64(n) / float64(bucketRatio))))

	s.a = make([][]T, numBucket)
	for i := 0; i < numBucket; i++ {
		s.a[i] = make([]T, 0, n/numBucket)
		s.a[i] = append(s.a[i], a[i*n/numBucket:(i+1)*n/numBucket]...)
	}

	return s
}

func (s *SortedSet[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for _, bucket := range s.a {
			for _, v := range bucket {
				if !yield(idx, v) {
					return
				}
				idx++
			}
		}
	}
}

func (s *SortedSet[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, bucket := range s.a {
			for _, v := range bucket {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (s *SortedSet[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := s.size - 1
		for i := range s.size {
			for j := range len(s.a[s.size-i-1]) {
				if !yield(idx, s.a[s.size-i-1][len(s.a[s.size-i-1])-j-1]) {
					return
				}
				idx--
			}
		}
	}
}

func (s *SortedSet[T]) Len() int {
	return s.size
}

func (s *SortedSet[T]) Equals(other *SortedSet[T]) bool {
	if s.size != other.size {
		return false
	}
	for i := range s.a {
		if !slices.Equal(s.a[i], other.a[i]) {
			return false
		}
	}
	return true
}

func (s *SortedSet[T]) String() string {
	sb := &strings.Builder{}
	_, _ = sb.WriteString("SortedSet{")

	for i := range s.a {
		for j := range s.a[i] {
			_, _ = sb.WriteString(fmt.Sprintf("%v, ", s.a[i][j]))
		}
	}

	_, _ = sb.WriteString("}")

	return sb.String()
}

func (s *SortedSet[T]) position(x T) ([]T, int, int) {
	var bucket int
	var a []T
	for bucket, a = range s.a {
		if x <= a[len(a)-1] {
			break
		}
	}
	i, _ := slices.BinarySearch(a, x)
	return a, bucket, i
}

func (s *SortedSet[T]) Contains(x T) bool {
	if s.size == 0 {
		return false
	}
	a, _, i := s.position(x)
	return i < len(a) && a[i] == x
}

func (s *SortedSet[T]) Add(x T) bool {
	if s.size == 0 {
		s.a = [][]T{{x}}
		s.size = 1
		return true
	}
	a, b, i := s.position(x)
	if i < len(a) && a[i] == x {
		return false
	}
	a = slices.Insert(a, i, x)
	s.size++
	if len(a) > len(s.a)*splitRatio {
		mid := len(a) >> 1
		s.a = slices.Insert(s.a, b+1, a[:mid])
		s.a[b] = a[mid:]
	}
	return true
}

func (s *SortedSet[T]) pop(a []T, b int, i int) T {
	ans := slices.Delete(a, i, i+1)[0]
	s.size--
	if len(a) == 0 {
		_ = slices.Delete(s.a, b, b+1)
	}
	return ans
}

func (s *SortedSet[T]) Discard(x T) bool {
	if s.size == 0 {
		return false
	}
	a, b, i := s.position(x)
	if i == len(a) || a[i] != x {
		return false
	}
	_ = s.pop(a, b, i)
	return true
}

func (s *SortedSet[T]) Lt(x T) (T, bool) {
	for i := range s.a {
		a := s.a[len(s.a)-i-1]
		if a[0] < x {
			j, _ := slices.BinarySearch(a, x)
			return a[j-1], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) Le(x T) (T, bool) {
	for i := range s.a {
		a := s.a[len(s.a)-i-1]
		if a[0] <= x {
			j, ok := slices.BinarySearch(a, x)
			if !ok {
				return a[j-1], true
			}
			return a[j], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) Gt(x T) (T, bool) {
	for _, a := range s.a {
		if a[len(a)-1] > x {
			j, _ := slices.BinarySearch(a, x)
			return a[j+1], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) Ge(x T) (T, bool) {
	for _, a := range s.a {
		if a[len(a)-1] >= x {
			j, ok := slices.BinarySearch(a, x)
			if !ok {
				return a[j], true
			}
			return a[j+1], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) Get(idx int) (T, error) {
	if idx < 0 {
		for i := range s.a {
			a := s.a[len(s.a)-i-1]
			idx += len(a)
			if idx >= 0 {
				return a[idx], nil
			}
		}
	} else {
		for _, a := range s.a {
			if idx < len(a) {
				return a[idx], nil
			}
			idx -= len(a)
		}
	}

	var v T
	return v, ErrIndexOutOfRange
}

func (s *SortedSet[T]) Pop(idx int) (T, error) {
	if idx < 0 {
		for i := range s.a {
			b := len(s.a) - i - 1
			a := s.a[b]
			idx += len(a)
			if idx >= 0 {
				return s.pop(a, -(b + 1), idx), nil
			}
		}
	} else {
		for b, a := range s.a {
			if idx < len(a) {
				return s.pop(a, b, idx), nil
			}
			idx -= len(a)
		}
	}
	var v T
	return v, ErrIndexOutOfRange
}

func (s *SortedSet[T]) CountLt(x T) int {
	ans := 0
	for _, a := range s.a {
		if a[len(a)-1] >= x {
			i, _ := slices.BinarySearch(a, x)
			return ans + i
		}
		ans += len(a)
	}
	return ans
}

func (s *SortedSet[T]) CountLe(x T) int {
	ans := 0
	for _, a := range s.a {
		if a[len(a)-1] >= x {
			i, ok := slices.BinarySearch(a, x)
			if !ok {
				return ans + i
			}
			return ans + i + 1
		}
		ans += len(a)
	}
	return ans
}
