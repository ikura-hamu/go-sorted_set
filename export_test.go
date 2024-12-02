package gosortedset

func (ss *SortedSet[T]) Buckets() [][]T {
	return ss.buckets
}
