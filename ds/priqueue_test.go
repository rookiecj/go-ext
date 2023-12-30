package ds

import (
	"testing"
	"unsafe"
)

type myStruct struct {
	Node  Node
	value string
}

var nodeOffset = unsafe.Offsetof(myStruct{}.Node)

func equalValue(left *Node, right *Node) bool {
	if left == nil || right == nil {
		return false
	}
	leftStruct := (*myStruct)(unsafe.Add(unsafe.Pointer(left), -nodeOffset))
	rightStruct := (*myStruct)(unsafe.Add(unsafe.Pointer(left), -nodeOffset))

	return leftStruct.value == rightStruct.value
}

func TestPriorityQueue_Push(t *testing.T) {
	type args struct {
		items []*myStruct
	}
	tests := []struct {
		name string
		pq   *PriorityQueue
		args args
		want args
	}{
		{
			name: "priority - 1,2,3",
			pq: func() *PriorityQueue {
				pq := NewPriorityQueue()
				return pq
			}(),
			args: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
					},
					{
						Node: Node{
							priority: 2,
						},
					},
					{
						Node: Node{
							priority: 3,
						},
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
					},
					{
						Node: Node{
							priority: 2,
						},
					},
					{
						Node: Node{
							priority: 3,
						},
					},
				},
			},
		},

		{
			name: "priority - 1,1,1",
			pq: func() *PriorityQueue {
				pq := NewPriorityQueue()
				return pq
			}(),
			args: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
						value: "1",
					},
					{
						Node: Node{
							priority: 1,
						},
						value: "2",
					},
					{
						Node: Node{
							priority: 1,
						},
						value: "3",
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
						value: "1",
					},
					{
						Node: Node{
							priority: 1,
						},
						value: "2",
					},
					{
						Node: Node{
							priority: 1,
						},
						value: "3",
					},
				},
			},
		},

		{
			name: "priority - 1,2,3,1,2",
			pq: func() *PriorityQueue {
				pq := NewPriorityQueue()
				return pq
			}(),
			args: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
						value: "1",
					},
					{
						Node: Node{
							priority: 2,
						},
						value: "2",
					},
					{
						Node: Node{
							priority: 3,
						},
						value: "3",
					},
					{
						Node: Node{
							priority: 1,
						},
						value: "4",
					},
					{
						Node: Node{
							priority: 2,
						},
						value: "5",
					},
				},
			},
			want: args{
				items: []*myStruct{
					{
						Node: Node{
							priority: 1,
						},
					},
					{
						Node: Node{
							priority: 1,
						},
					},
					{
						Node: Node{
							priority: 2,
						},
					},
					{
						Node: Node{
							priority: 2,
						},
					},
					{
						Node: Node{
							priority: 3,
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
				pq = NewPriorityQueue()
			}
			for _, item := range tt.args.items {
				pq.Push(&item.Node)
			}

			wantLen := pq.Len()

			var gotNodes []*Node
			for pq.Len() > 0 {
				gotNodes = append(gotNodes, pq.Pop())
			}
			gotLen := len(gotNodes)
			if gotLen != wantLen {
				t.Errorf("PriorityQueue size want %d got %d", wantLen, gotLen)
			}

			for idx, got := range gotNodes {
				//if !reflect.DeepEqual(&tt.want.items[idx].Node, got) {
				//	t.Errorf("PriorityQueue priority diff want %v, got %v", tt.want.items[idx].Node, got)
				//}
				want := &tt.want.items[idx].Node
				if want.priority != got.priority || !equalValue(want, got) {
					t.Errorf("PriorityQueue priority diff want %v, got %v", tt.want.items[idx].Node, got)
				}
			}
		})
	}
}
