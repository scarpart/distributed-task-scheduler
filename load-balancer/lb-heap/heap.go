package lbheap

import (
	"errors"

	"github.com/scarpart/distributed-task-scheduler/load-balancer/server"
)

type Heap struct {
	Array []*server.RemoteServer
	Size  int
}

func NewHeap() Heap {
	return Heap{
		Array: make([]*server.RemoteServer, 5),
		Size: 0,
	}
}

func (heap *Heap) EnsureExtraCapacity() {
	heap.Array = append(heap.Array, make([]*server.RemoteServer, len(heap.Array))...)
}

func (heap *Heap) Poll() (*server.RemoteServer, error) {
	if (heap.Size <= 0) { return nil, errors.New("The heap is empty.") }
	firstServer := heap.Array[0]
	heap.Array[0] = heap.Array[heap.Size-1]
	heap.Size--
	heap.HeapifyDown()
	return firstServer, nil
}

func (heap *Heap) Add(server *server.RemoteServer) {
	heap.EnsureExtraCapacity()
	heap.Array[heap.Size] = server
	heap.Size++
	heap.HeapifyUp()
}

func (heap *Heap) HeapifyUp() {
	index := heap.Size - 1
	for heap.HasParent(index) && heap.Parent(index) > heap.Array[index].Value() {
		heap.Swap(heap.GetParentIndex(index), index)
		index = heap.GetParentIndex(index)
	}
}

func (heap *Heap) HeapifyDown() {
	index := 0
	for heap.HasLeftChild(index) {
		smallerChildIndex := heap.GetLeftChildIndex(index)
		if heap.HasRightChild(index) && heap.RightChild(index) < heap.LeftChild(index) {
			smallerChildIndex = heap.GetRightChildIndex(index)
		}
		if heap.Array[index].Value() < heap.Array[smallerChildIndex].Value() {
			return
		} else {
			heap.Swap(index, smallerChildIndex)	
		}
		index = smallerChildIndex
	}
}

func (heap *Heap) Swap(i1 int, i2 int) {
	temp := heap.Array[i1]
	heap.Array[i1] = heap.Array[i2]
	heap.Array[i2] = temp
}

func (heap *Heap) GetRightChildIndex(index int) int { return 2 * index + 2 }
func (heap *Heap) GetLeftChildIndex(index int) int { return 2 * index + 1 }
func (heap *Heap) GetParentIndex(index int) int { return (index / 2) - 1 }

func (heap *Heap) HasLeftChild(index int) bool { return 2 * index + 1 < heap.Size } 
func (heap *Heap) HasRightChild(index int) bool { return 2 * index + 2 < heap.Size } 
func (heap *Heap) HasParent(index int) bool { return index / 2 - 1 < 0}

func (heap *Heap) LeftChild(index int) float32 { 
	server := heap.Array[heap.GetLeftChildIndex(index)]
	return float32(server.Weight) * (server.MEM_Usage + server.CPU_Usage)
} 

func (heap *Heap) RightChild(index int) float32 { 
	server := heap.Array[heap.GetRightChildIndex(index)]
	return float32(server.Weight) * (server.MEM_Usage + server.CPU_Usage)
} 

func (heap *Heap) Parent(index int) float32 { 
	server := heap.Array[heap.GetParentIndex(index)]
	return float32(server.Weight) * (server.MEM_Usage + server.CPU_Usage)
} 
