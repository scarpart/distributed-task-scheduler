package loadbalancer

import (
	"errors"
	"fmt"
)

type Heap []*RemoteServer

func (heap Heap) Len() int {
	return len(heap)
}

func (heap *Heap) Root() (*RemoteServer, error) {
	if (heap.Len() <= 0) { return nil, errors.New("The heap is empty.") }
	return (*heap)[0], nil
}

// Returns the root element of the min heap
func (heap *Heap) LeastConnections() *RemoteServer {
	fmt.Printf("least connections : %v\n", heap)
	leastConns := (*heap)[0]
	fmt.Printf("%v\n", heap)
	return leastConns	
}

func (heap *Heap) Poll() (*RemoteServer, error) {
	if (heap.Len() <= 0) { return nil, errors.New("The heap is empty.") }
	firstServer := (*heap)[0]
	(*heap)[0] = (*heap)[heap.Len()-1]
	heap.HeapifyDown(0)
	return firstServer, nil
}

func (heap *Heap) Add(server *RemoteServer) {
	*heap = append(*heap, server)
	heap.HeapifyUp(heap.Len() - 1)
	fmt.Printf("after add : %v\n", heap)
}

func (heap *Heap) HeapifyUp(index int) {
	for heap.HasParent(index) && heap.Parent(index) > (*heap)[0].Connections {
		heap.Swap(heap.GetParentIndex(index), index)
		index = heap.GetParentIndex(index)
	}
}

func (heap *Heap) HeapifyDown(index int) {
	for heap.HasLeftChild(index) {
		smallerChildIndex := heap.GetLeftChildIndex(index)
		if heap.HasRightChild(index) && heap.RightChild(index) < heap.LeftChild(index) {
			smallerChildIndex = heap.GetRightChildIndex(index)
		}
		if (*heap)[index].Connections < (*heap)[smallerChildIndex].Connections {
			return
		} else {
			heap.Swap(index, smallerChildIndex)	
		}
		index = smallerChildIndex
	}
}

func (heap *Heap) GetServerIndex(server *RemoteServer) (int, error) {
	for k, v := range *heap {
		if v == server {
			return k, nil
		}
	}
	return 0, errors.New("Could not find the server in the heap.") 
}

func (heap *Heap) Fix(server *RemoteServer) {
	index, err := heap.GetServerIndex(server)
	if err != nil {
		fmt.Println("err is not nil", err)
		return // handle this better, probably (since the heap is guaranteed to work, i think there is no need for an error)
	}

	fmt.Println("ongoing formatting heap")
	if heap.Parent(index) <= server.Connections {
		fmt.Println("heapifying down (correct)")
		heap.HeapifyDown(index)	
	} else {
		fmt.Println("heapifying up (incorrect)")
		heap.HeapifyUp(index)
	}
}

func (heap *Heap) Swap(i1 int, i2 int) {
	temp := (*heap)[i1]
	(*heap)[i1] = (*heap)[i2]
	(*heap)[i2] = temp
}

func (heap *Heap) GetRightChildIndex(index int) int { return 2 * index + 2 }
func (heap *Heap) GetLeftChildIndex(index int) int { return 2 * index + 1 }
func (heap *Heap) GetParentIndex(index int) int { return (index - 1) / 2 }

func (heap *Heap) HasLeftChild(index int) bool { return 2 * index + 1 < heap.Len() } 
func (heap *Heap) HasRightChild(index int) bool { return 2 * index + 2 < heap.Len() } 
func (heap *Heap) HasParent(index int) bool { return (index - 1) / 2 > 0}

func (heap *Heap) LeftChild(index int) int32 { 
	server := (*heap)[heap.GetLeftChildIndex(index)]
	return server.Connections
} 

func (heap *Heap) RightChild(index int) int32 { 
	server := (*heap)[heap.GetRightChildIndex(index)]
	return server.Connections
} 

func (heap *Heap) Parent(index int) int32 { 	
	server := (*heap)[heap.GetParentIndex(index)]
	return server.Connections
} 
