package ds

import (
	"errors"
	"math"
)

type Node struct {
	// 최소 힙(Min Heap)을 사용
	// 최소 힙에서는 부모 노드의 값이 자식 노드의 값보다 항상 작거나 같아야 합니다.
	priority int
	// stable
	index int
}

type PriorityQueue []*Node

var seq = 0

// NewPriorityQueue
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

// Push add a new node
func (pq *PriorityQueue) Push(node *Node) error {
	if pq == nil || node == nil {
		return errors.New("nil reference")
	}
	if seq == math.MaxInt {
		return errors.New("reached max sequence")
	}
	node.index = seq
	seq++

	*pq = append(*pq, node)
	pq.up(len(*pq) - 1)
	return nil
}

// Pop remove the highest priority node and return it
func (pq *PriorityQueue) Pop() *Node {
	if pq == nil {
		return nil
	}
	n := pq.Len()
	if n == 0 {
		return nil
	}
	old := *pq
	node := old[0]
	old[0] = old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	pq.down(0)
	return node
}

// Len return length of q
func (pq *PriorityQueue) Len() int {
	if pq == nil {
		return 0
	}
	return len(*pq)
}

func (pq *PriorityQueue) up(index int) {
	for {
		parent := (index - 1) / 2 // 부모 노드의 인덱스
		if index == 0 || !pq.less(index, parent) {
			break
		}
		pq.swap(parent, index)
		index = parent
	}
}

func (pq *PriorityQueue) down(index int) {
	n := len(*pq)
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

func (pq *PriorityQueue) less(i, j int) bool {
	// 우선순위가 같으면 추가된 순서 (index) 비교
	if (*pq)[i].priority == (*pq)[j].priority {
		return (*pq)[i].index < (*pq)[j].index
	}
	return (*pq)[i].priority < (*pq)[j].priority
}

func (pq *PriorityQueue) swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}
