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
	buckets [][]T
	size    int
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

	s.buckets = make([][]T, numBucket)
	for i := 0; i < numBucket; i++ {
		s.buckets[i] = make([]T, 0, n/numBucket)
		s.buckets[i] = append(s.buckets[i], a[i*n/numBucket:(i+1)*n/numBucket]...)
	}

	return s
}

func (s *SortedSet[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for _, bucket := range s.buckets {
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
		for _, bucket := range s.buckets {
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
			for j := range len(s.buckets[s.size-i-1]) {
				if !yield(idx, s.buckets[s.size-i-1][len(s.buckets[s.size-i-1])-j-1]) {
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
	for i := range s.buckets {
		if !slices.Equal(s.buckets[i], other.buckets[i]) {
			return false
		}
	}
	return true
}

func (s *SortedSet[T]) String() string {
	sb := &strings.Builder{}
	_, _ = sb.WriteString("SortedSet{")

	for i := range s.buckets {
		for j := range s.buckets[i] {
			_, _ = sb.WriteString(fmt.Sprintf("%v, ", s.buckets[i][j]))
		}
	}

	_, _ = sb.WriteString("}")

	return sb.String()
}

// return the bucket, index of the bucket and position in which x should be.
func (s *SortedSet[T]) position(x T) (*[]T, int, int) {
	var bucket int
	var a *[]T
	for bucket = range s.buckets {
		a = &s.buckets[bucket]
		if x <= (*a)[len(*a)-1] {
			break
		}
	}
	i, _ := slices.BinarySearch(*a, x)
	return a, bucket, i
}

func (s *SortedSet[T]) Contains(x T) bool {
	if s.size == 0 {
		return false
	}
	a, _, i := s.position(x)
	return i < len(*a) && (*a)[i] == x
}

func (s *SortedSet[T]) Add(x T) bool {
	if s.size == 0 {
		s.buckets = [][]T{{x}}
		s.size = 1
		return true
	}
	a, b, i := s.position(x)
	if i != len(*a) && (*a)[i] == x {
		return false
	}
	*a = slices.Insert(*a, i, x)
	s.buckets[b] = *a
	s.size++

	if len(*a) > len(s.buckets)*splitRatio {
		mid := len(*a) >> 1
		s.buckets = slices.Insert(s.buckets, b+1, (*a)[mid:])
		s.buckets[b] = (*a)[:mid]
	}
	return true
}

func (s *SortedSet[T]) pop(a *[]T, b int, i int) T {
	ans := (*a)[i]
	*a = slices.Delete(*a, i, i+1)[:len(*a)-1]
	s.size--
	if len(*a) == 0 {
		if b < 0 {
			b = b + len(s.buckets)
		}
		s.buckets = slices.Delete(s.buckets, b, b+1)
		if len(s.buckets) > 1 { // TODO: 上で代入してるから、ここは必要ではないが、どっちが速いか調べる必要がある
			s.buckets = s.buckets[:len(s.buckets)-1]
		}
	}
	return ans
}

func (s *SortedSet[T]) Discard(x T) bool {
	if s.size == 0 {
		return false
	}
	a, b, i := s.position(x)
	if i == len(*a) || (*a)[i] != x {
		return false
	}
	_ = s.pop(a, b, i)

	return true
}

func (s *SortedSet[T]) Lt(x T) (T, bool) {
	for i := range s.buckets {
		a := s.buckets[len(s.buckets)-i-1]
		if a[0] < x {
			j, _ := slices.BinarySearch(a, x)
			return a[j-1], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) Le(x T) (T, bool) {
	for i := range s.buckets {
		a := s.buckets[len(s.buckets)-i-1]
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
	for _, a := range s.buckets {
		if a[len(a)-1] > x {
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

func (s *SortedSet[T]) Ge(x T) (T, bool) {
	for _, a := range s.buckets {
		if a[len(a)-1] >= x {
			j, ok := slices.BinarySearch(a, x)
			if !ok {
				return a[j], true
			}
			return a[j], true
		}
	}
	var v T
	return v, false
}

func (s *SortedSet[T]) GetItem(idx int) (T, error) {
	if idx < 0 {
		for i := range s.buckets {
			a := s.buckets[len(s.buckets)-i-1]
			idx += len(a)
			if idx >= 0 {
				return a[idx], nil
			}
		}
	} else {
		for _, a := range s.buckets {
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
		for b := range s.buckets {
			a := &s.buckets[len(s.buckets)-b-1]
			idx += len(*a)
			if idx >= 0 {
				v := s.pop(a, -(b + 1), idx)
				return v, nil
			}
		}
	} else {
		for b := range s.buckets {
			a := &s.buckets[b]
			if idx < len(*a) {
				return s.pop(a, b, idx), nil
			}
			idx -= len(*a)
		}
	}
	var v T
	return v, ErrIndexOutOfRange
}

func (s *SortedSet[T]) Index(x T) int {
	ans := 0
	for _, a := range s.buckets {
		if a[len(a)-1] >= x {
			i, _ := slices.BinarySearch(a, x)
			return ans + i
		}
		ans += len(a)
	}
	return ans
}

func (s *SortedSet[T]) IndexRight(x T) int {
	ans := 0
	for _, a := range s.buckets {
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
