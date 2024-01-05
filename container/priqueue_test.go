package container

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

type myStruct struct {
	Node  PQNode
	value string
}

var nodeOffset = unsafe.Offsetof(myStruct{}.Node)

func equalValue(left *PQNode, right *PQNode) bool {
	if left == nil || right == nil {
		return false
	}
	leftStruct := (*myStruct)(unsafe.Add(unsafe.Pointer(left), -nodeOffset))
	rightStruct := (*myStruct)(unsafe.Add(unsafe.Pointer(left), -nodeOffset))

	return leftStruct.value == rightStruct.value
}

func TestPriorityQueue_PushPop(t *testing.T) {
	type args struct {
		items []*myStruct
	}
	tests := []struct {
		name string
		pq   PriorityQueue
		args args
		want args
	}{
		{
			name: "Priority - 1,2,3",
			args: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
					},
					{
						Node: PQNode{
							Priority: 2,
						},
					},
					{
						Node: PQNode{
							Priority: 3,
						},
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
					},
					{
						Node: PQNode{
							Priority: 2,
						},
					},
					{
						Node: PQNode{
							Priority: 3,
						},
					},
				},
			},
		},

		{
			name: "Priority - 1,1,1",
			args: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "1",
					},
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "2",
					},
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "3",
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "1",
					},
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "2",
					},
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "3",
					},
				},
			},
		},

		{
			name: "Priority - 1,2,3,1,2",
			args: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "1",
					},
					{
						Node: PQNode{
							Priority: 2,
						},
						value: "2",
					},
					{
						Node: PQNode{
							Priority: 3,
						},
						value: "3",
					},
					{
						Node: PQNode{
							Priority: 1,
						},
						value: "4",
					},
					{
						Node: PQNode{
							Priority: 2,
						},
						value: "5",
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: PQNode{
							Priority: 1,
						},
					},
					{
						Node: PQNode{
							Priority: 1,
						},
					},
					{
						Node: PQNode{
							Priority: 2,
						},
					},
					{
						Node: PQNode{
							Priority: 2,
						},
					},
					{
						Node: PQNode{
							Priority: 3,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := tt.pq
			if pq == nil {
				pq = NewPriorityQueue(0)
			}
			for _, item := range tt.args.items {
				pq.Push(&item.Node)
			}

			wantLen := pq.Len()

			var gotNodes []*PQNode
			for pq.Len() > 0 {
				gotNodes = append(gotNodes, pq.Pop())
			}
			gotLen := len(gotNodes)
			if gotLen != wantLen {
				t.Errorf("PriorityQueue size want %d got %d", wantLen, gotLen)
			}

			for idx, got := range gotNodes {
				//if !reflect.DeepEqual(&tt.want.items[idx].PQNode, got) {
				//	t.Errorf("PriorityQueue Priority diff want %v, got %v", tt.want.items[idx].PQNode, got)
				//}
				want := &tt.want.items[idx].Node
				if want.Priority != got.Priority || !equalValue(want, got) {
					t.Errorf("PriorityQueue Priority diff want %v, got %v", tt.want.items[idx].Node, got)
				}
			}
		})
	}
}

func TestPriorityQueue_PushPopMany(t *testing.T) {
	q := NewPriorityQueue(0)
	limit := 10240 // linear 0s, rand 2 millis
	//limit := 10000000 // linear 3s, rand 11s
	//limit := 100000000 // linear 28s, rand 206s
	//limit := 1000000000 // linear ?,

	t.Run(fmt.Sprintf("limit %d", limit), func(t *testing.T) {

		start := time.Now()
		for idx := 0; idx < limit; idx++ {
			q.Push(&PQNode{
				//Priority: idx, // linear
				Priority: rand.Intn(limit), // rand
			})
		}

		var gotNodes []*PQNode
		for q.Len() > 0 {
			gotNodes = append(gotNodes, q.Pop())
		}
		fmt.Printf("PushMany took %d millis for %d items\n", time.Now().Sub(start).Milliseconds(), limit)

		if len(gotNodes) != limit {
			t.Errorf("PushMany length want=%d, got=%d", limit, len(gotNodes))
			return
		}

		prevPriority := 0
		for idx := 0; idx < limit; idx++ {
			if gotNodes[idx].Priority < prevPriority {
				t.Errorf("PushMany Priority want=%d, got=%d", prevPriority, gotNodes[idx].Priority)
				break
			}
			prevPriority = gotNodes[idx].Priority
		}
	})
}

func TestPriorityQueueN_PushPopMany(t *testing.T) {
	//limit := 10240 //  2 millis
	limit := 10000000 //  10s
	//limit := 100000000 //  192s
	//limit := 1000000000 //

	t.Run(fmt.Sprintf("limit %d", limit), func(t *testing.T) {
		q := NewPriorityQueue(limit)

		start := time.Now()
		for idx := 0; idx < limit; idx++ {
			q.Push(&PQNode{
				//Priority: idx, // linear
				Priority: rand.Intn(limit), // rand
			})
		}

		var gotNodes []*PQNode
		for q.Len() > 0 {
			gotNodes = append(gotNodes, q.Pop())
		}
		fmt.Printf("PushPopManyN took %d millis for %d items\n", time.Now().Sub(start).Milliseconds(), limit)

		if len(gotNodes) != limit {
			t.Errorf("PushPopManyN length want=%d, got=%d", limit, len(gotNodes))
			return
		}

		prevPriority := 0
		for idx := 0; idx < limit; idx++ {
			if gotNodes[idx].Priority < prevPriority {
				t.Errorf("PushPopManyN Priority want=%d, got=%d", prevPriority, gotNodes[idx].Priority)
				break
			}
			prevPriority = gotNodes[idx].Priority
		}
	})
}

func TestPriorityQueueN_PushPopRand(t *testing.T) {
	//limit := 10240 // rand 2 millis
	limit := 10000000 // 7479 millis for 6667121 items
	//limit := 100000000 // 130791 millis for 66665244 items

	t.Run(fmt.Sprintf("limit %d", limit), func(t *testing.T) {
		q := NewPriorityQueue(limit)

		start := time.Now()
		for idx := 0; idx < limit; idx++ {
			pri := rand.Intn(limit)
			q.Push(&PQNode{
				//Priority: idx, // linear
				Priority: pri, // rand
			})
			if pri%3 == 0 {
				q.Pop()
			}
		}

		var gotNodes []*PQNode
		for q.Len() > 0 {
			gotNodes = append(gotNodes, q.Pop())
		}
		fmt.Printf("PushPopRand took %d millis for %d items\n", time.Now().Sub(start).Milliseconds(), len(gotNodes))

		prevPriority := 0
		for idx := 0; idx < len(gotNodes); idx++ {
			if gotNodes[idx].Priority < prevPriority {
				t.Errorf("PushPopRand Priority want=%d, got=%d", prevPriority, gotNodes[idx].Priority)
				break
			}
			prevPriority = gotNodes[idx].Priority
		}
	})
}
