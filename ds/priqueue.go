package ds

// Node은 우선순위 큐에서 사용할 구조체입니다.
type Node struct {
	// 우선순위
	// 최소 힙(Min Heap)을 사용
	// 최소 힙에서는 부모 노드의 값이 자식 노드의 값보다 항상 작거나 같아야 합니다.
	priority int
}

// PriorityQueue는 Node의 슬라이스입니다.
type PriorityQueue []*Node

// NewPriorityQueue는 새로운 PriorityQueue를 생성합니다.
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

// Push는 새로운 항목을 큐에 추가합니다.
func (pq *PriorityQueue) Push(node *Node) {
	*pq = append(*pq, node)
	pq.up(len(*pq) - 1)
}

// Pop은 가장 높은 우선순위의 항목을 제거하고 반환합니다.
func (pq *PriorityQueue) Pop() *Node {
	old := *pq
	n := len(old)
	node := old[0]
	old[0] = old[n-1]
	old[n-1] = nil
	*pq = old[:n-1]
	pq.down(0)
	return node
}

// Len은 큐의 길이를 반환합니다.
func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

// up은 힙 속성을 유지하면서 항목을 올립니다.
func (pq *PriorityQueue) up(index int) {
	for {
		parent := (index - 1) / 2
		if index == 0 || (*pq)[parent].priority <= (*pq)[index].priority {
			break
		}
		(*pq)[parent], (*pq)[index] = (*pq)[index], (*pq)[parent]
		index = parent
	}
}

// down은 힙 속성을 유지하면서 항목을 내립니다.
func (pq *PriorityQueue) down(index int) {
	n := len(*pq)
	for {
		left := 2*index + 1
		if left >= n {
			break
		}
		smallest := left
		if right := left + 1; right < n && (*pq)[right].priority < (*pq)[left].priority {
			smallest = right
		}
		if (*pq)[index].priority <= (*pq)[smallest].priority {
			break
		}
		(*pq)[index], (*pq)[smallest] = (*pq)[smallest], (*pq)[index]
		index = smallest
	}
}
