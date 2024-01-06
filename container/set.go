package container

type Set[T comparable] map[T]bool

// NewSet 새로운 Set 생성
func NewSet[T comparable]() Set[T] {
	return make(map[T]bool)
}

// Insert 존재하지 않으면 추가, 존재하면 업데이트
func (s Set[T]) Insert(element T) {
	s[element] = true
}

// Remove 삭제 연산
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Has 포함 여부 확인
func (s Set[T]) Has(element T) bool {
	_, exists := s[element]
	return exists
}

// Len 크기 반환
func (s Set[T]) Len() int {
	return len(s)
}

// Iterate 순회
func (s Set[T]) Iterate() []T {
	var keys []T
	for key := range s {
		keys = append(keys, key)
	}
	return keys
}

// Equals 서로 동일한지 판별하는 함수
func (s Set[T]) Equals(other Set[T]) bool {
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

func (s Set[T]) Copy() Set[T] {
	result := NewSet[T]()
	for ele := range s {
		result.Insert(ele)
	}
	return result
}

// Intersection 교집합 연산
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := NewSet[T]()
	for key := range s {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union 합집합 연산
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := NewSet[T]()
	for key := range s {
		result.Insert(key)
	}
	for key := range other {
		result.Insert(key)
	}
	return result
}

// Difference 차집합 연산
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := NewSet[T]()
	for key := range s {
		if !other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// SymmetricDifference 대칭 차집합 연산, 두 Set 중 하나에만 포함된 요소
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	result := NewSet[T]()
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
