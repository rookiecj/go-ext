package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/ds"
	"unsafe"
)

type myStruct struct {
	Node  ds.PQNode
	value string
}

var nodeOffset = unsafe.Offsetof(myStruct{}.Node)

func getItem(node *ds.PQNode) *myStruct {
	s := (*myStruct)(unsafe.Add(unsafe.Pointer(node), -nodeOffset))
	return s
}

func main() {
	items := []*myStruct{
		{
			Node: ds.PQNode{
				Priority: 1,
			},
			value: "pri 1 - first",
		},
		{
			Node: ds.PQNode{
				Priority: 2,
			},
			value: "pri 2 - second",
		},
		{
			Node: ds.PQNode{
				Priority: 1,
			},
			value: "pri 1 - third",
		},
	}

	pq := ds.NewPriorityQueue()
	for _, item := range items {
		pq.Push(&item.Node)
	}

	fmt.Println("total:", pq.Len())

	for pq.Len() > 0 {
		node := pq.Pop()
		item := getItem(node)
		fmt.Printf("pri: %d, value=%s\n", item.Node.Priority, item.value)
	}
}
