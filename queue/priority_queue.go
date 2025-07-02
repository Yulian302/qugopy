// Package queue provides a priority queue implementation specifically for Tasks.
// It supports a min-heap (or max-heap using negative priorities) ordering. Tasks with highest priorities (lowest) come first. Example: 1 > 2 > 3 > 5 ...
//
// Example (Min-Heap):
//
// pq := &PriorityQueue{}
// 	pq.Push(IntTask{ID: 1, Type: "python", Payload: "Task1", Priority: 5})
// 	pq.Push(IntTask{ID: 2, Type: "go", Payload: "Task2", Priority: 2})
// 	pq.Push(IntTask{ID: 3, Type: "python", Payload: "Task3", Priority: 3})
// 	pq.Push(IntTask{ID: 4, Type: "go", Payload: "Task4", Priority: 1})
// 1("Task4") -> 2("Task2") -> 3("Task3") -> 5("Task1")

package queue

type PriorityQueue struct {
	data []IntTask
}

// Heap defines the standard operations for a heap data structure.
// Implementations should ensure O(log n) time complexity for Push/Pop.
type Heap[T any] interface {
	// Parent returns the index of the parent node for the given index.
	Parent(idx int) int

	// LeftChild returns the index of the left child for the given index.
	LeftChild(idx int) int

	// RightChild returns the index of the right child for the given index.
	RightChild(idx int) int

	// Peek returns the root element without removing it.
	// The boolean indicates whether the heap was non-empty.
	Peek() (T, bool)

	// IsEmpty returns true if the heap has no elements.
	IsEmpty() bool

	// Push adds a new element to the heap.
	Push(val T)

	// Pop removes and returns the root element.
	// The boolean indicates whether the heap was non-empty.
	Pop() (T, bool)

	// Delete removes the specified element if found.
	// Returns true if the element was removed.
	Delete(priority uint16) bool

	// HeapifyUp rebalances the heap upward from the given index.
	HeapifyUp(idx int)

	// HeapifyDown rebalances the heap downward from the given index.
	HeapifyDown(idx int)
}

var _ Heap[IntTask] = (*PriorityQueue)(nil)

func (pq *PriorityQueue) Parent(index int) int {
	return (index - 1) / 2
}

func (pq *PriorityQueue) LeftChild(index int) int {
	return 2*index + 1
}

func (pq *PriorityQueue) RightChild(index int) int {
	return 2*index + 2
}

func (pq *PriorityQueue) Peek() (IntTask, bool) {
	if pq.IsEmpty() {
		var zero IntTask
		return zero, false
	}
	return pq.data[0], true
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.data) == 0
}

func (pq *PriorityQueue) Push(value IntTask) {
	pq.data = append(pq.data, value)
	pq.HeapifyUp(len(pq.data) - 1)
}

func (pq *PriorityQueue) Pop() (IntTask, bool) {
	if pq.IsEmpty() {
		var zero IntTask
		return zero, false
	}
	root := pq.data[0]
	last := pq.data[len(pq.data)-1]
	pq.data = pq.data[:len(pq.data)-1]

	if len(pq.data) > 0 {
		pq.data[0] = last
		pq.HeapifyDown(0)
	}
	return root, true
}

func (pq *PriorityQueue) Delete(priority uint16) bool {
	if pq.IsEmpty() {
		return false
	}
	for idx, val := range pq.data {
		if val.Priority == uint16(priority) {
			pq.data[idx] = pq.data[len(pq.data)-1]
			pq.data = pq.data[:len(pq.data)-1]
			if idx >= len(pq.data) {
				return true
			}
			parent := pq.Parent(idx)
			if parent >= 0 && pq.data[parent].GT(&pq.data[idx]) {
				pq.HeapifyUp(idx)
			} else {
				pq.HeapifyDown(idx)
			}
			return true
		}
	}
	return false
}

func (pq *PriorityQueue) HeapifyUp(index int) {
	for index > 0 {
		parent := pq.Parent(index)
		if pq.data[parent].GT(&pq.data[index]) {
			pq.data[parent], pq.data[index] = pq.data[index], pq.data[parent]
			index = parent
		} else {
			break
		}
	}
}

func (pq *PriorityQueue) HeapifyDown(index int) {
	size := len(pq.data)

	smallest := index
	left := pq.LeftChild(index)
	right := pq.RightChild(index)

	if left < size && pq.data[left].LT(&pq.data[smallest]) {
		smallest = left
	}

	if right < size && pq.data[right].LT(&pq.data[smallest]) {
		smallest = right
	}

	if smallest != index {
		pq.data[index], pq.data[smallest] = pq.data[smallest], pq.data[index]
		pq.HeapifyDown(smallest)
	}

}
