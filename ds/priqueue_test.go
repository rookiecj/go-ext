package ds

import (
	"testing"
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

func TestPriorityQueue_Push(t *testing.T) {
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
				pq = NewPriorityQueue()
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

func TestPriorityQueue_PushMany(t *testing.T) {
	q := NewPriorityQueue()
	limit := 10240
	for idx := 0; idx < limit; idx++ {
		q.Push(&PQNode{
			Priority: idx,
		})
	}

	var gotNodes []*PQNode
	for q.Len() > 0 {
		gotNodes = append(gotNodes, q.Pop())
	}

	if len(gotNodes) != limit {
		t.Errorf("PushMany length want=%d, got=%d", limit, len(gotNodes))
	}

	for idx := 0; idx < limit; idx++ {
		if gotNodes[idx].Priority != idx {
			t.Errorf("PushMany Priority want=%d, got=%d", idx, gotNodes[idx].Priority)
		}
	}
}
