package container

import (
	"cmp"
	"slices"
)

type SortedSet[T cmp.Ordered] map[T]bool

// NewSortedSet 새로운 SortedSet 생성
func NewSortedSet[T cmp.Ordered]() SortedSet[T] {
	return make(map[T]bool)
}

// Insert 삽입 연산
func (s SortedSet[T]) Insert(element T) {
	s[element] = true
}

// Remove 삭제 연산
func (s SortedSet[T]) Remove(element T) {
	delete(s, element)
}

// Has 포함 여부 확인
func (s SortedSet[T]) Has(element T) bool {
	_, exists := s[element]
	return exists
}

// Len 크기 반환
func (s SortedSet[T]) Len() int {
	return len(s)
}

// Equals 서로 동일한지 판별하는 함수
func (s SortedSet[T]) Equals(other SortedSet[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for key := range s {
		if !other.Has(key) {
			return false
		}
	}
	return true
}

// sort.Interface
// Less(i, j int) bool
// Swap(i, j int)

// Iterate 순회
func (s SortedSet[T]) Iterate() []T {
	var keys []T
	for key := range s {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

// Intersection 교집합 연산
func (s SortedSet[T]) Intersection(other SortedSet[T]) SortedSet[T] {
	result := NewSortedSet[T]()
	for key := range s {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union 합집합 연산
func (s SortedSet[T]) Union(other SortedSet[T]) SortedSet[T] {
	result := NewSortedSet[T]()
	for key := range s {
		result.Insert(key)
	}
	for key := range other {
		result.Insert(key)
	}
	return result
}

// Difference 차집합 연산
func (s SortedSet[T]) Difference(other SortedSet[T]) SortedSet[T] {
	result := NewSortedSet[T]()
	for key := range s {
		if !other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// SymmetricDifference 대칭 차집합 연산, 두 Set 중 하나에만 포함된 요소
func (s SortedSet[T]) SymmetricDifference(other SortedSet[T]) SortedSet[T] {
	result := NewSortedSet[T]()
	for key := range s {
		if !other.Has(key) {
			result.Insert(key)
		}
	}
	for key := range other {
		if !s.Has(key) {
			result.Insert(key)
		}
	}
	return result
}
