package ds

import (
	"errors"
	"math"
)

// PriorityQueue is a Priority queue based on a min heap.
type PriorityQueue interface {
	Push(*PQNode) error
	Pop() *PQNode
	Len() int
}

type PQNode struct {
	// 최소 힙(Min Heap)을 사용
	// 최소 힙에서는 부모 노드의 값이 자식 노드의 값보다 항상 작거나 같아야 합니다.
	Priority int
	// stable for internal
	index int
}

type priorityQueue struct {
	q   []*PQNode
	seq int // sequence for index
}

// NewPriorityQueue returns a new PriorityQueue.
func NewPriorityQueue() PriorityQueue {
	return &priorityQueue{}
}

// Push add a new node
func (pq *priorityQueue) Push(node *PQNode) error {
	if pq == nil || node == nil {
		return errors.New("nil reference")
	}
	if pq.seq == math.MaxInt {
		return errors.New("reached max sequence")
	}
	node.index = pq.seq
	pq.seq++

	pq.q = append(pq.q, node)
	pq.up(len(pq.q) - 1)
	return nil
}

// Pop remove the highest priority node and return it
// returns nil if q is empty
func (pq *priorityQueue) Pop() *PQNode {
	if pq == nil {
		return nil
	}
	n := pq.Len()
	if n == 0 {
		return nil
	}

	node := pq.q[0]
	pq.q[0] = pq.q[n-1]
	//pq.q = append(pq.q[:0], pq.q[:n-1]...)
	pq.q = pq.q[:n-1]
	pq.down(0)
	return node
}

// Len return length of q
func (pq *priorityQueue) Len() int {
	if pq == nil || pq.q == nil {
		return 0
	}
	return len(pq.q)
}

func (pq *priorityQueue) up(index int) {
	for {
		parent := (index - 1) / 2 // 부모 노드의 인덱스
		if index == 0 || !pq.less(index, parent) {
			break
		}
		pq.swap(parent, index)
		index = parent
	}
}

func (pq *priorityQueue) down(index int) {
	n := pq.Len()
	for {
		left := 2*index + 1
		if left >= n || left < 0 { // left < 0: 오버플로우 방지
			break
		}
		smallest := left // 왼쪽 자식
		if right := left + 1; right < n && pq.less(right, left) {
			smallest = right // 오른쪽 자식이 더 작으면 선택
		}
		if !pq.less(smallest, index) {
			break
		}
		pq.swap(index, smallest)
		index = smallest
	}
}

func (pq *priorityQueue) less(i, j int) bool {
	// 우선순위가 같으면 추가된 순서 (index) 비교
	if pq.q[i].Priority == pq.q[j].Priority {
		return pq.q[i].index < pq.q[j].index
	}
	return pq.q[i].Priority < pq.q[j].Priority
}

func (pq *priorityQueue) swap(i, j int) {
	pq.q[i], pq.q[j] = pq.q[j], pq.q[i]
}
