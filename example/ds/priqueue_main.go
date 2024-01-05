package main

import (
	"fmt"
	"github.com/rookiecj/go-langext/container"
	"unsafe"
)

type myStruct struct {
	Node  container.PQNode
	value string
}

var nodeOffset = unsafe.Offsetof(myStruct{}.Node)

func getItem(node *container.PQNode) *myStruct {
	s := (*myStruct)(unsafe.Add(unsafe.Pointer(node), -nodeOffset))
	return s
}

func main() {
	items := []*myStruct{
		{
			Node: container.PQNode{
				Priority: 1,
			},
			value: "pri 1 - first",
		},
		{
			Node: container.PQNode{
				Priority: 2,
			},
			value: "pri 2 - second",
		},
		{
			Node: container.PQNode{
				Priority: 1,
			},
			value: "pri 1 - third",
		},
	}

	pq := container.NewPriorityQueue(0)
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
