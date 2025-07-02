package queue

type PriorityQueue struct {
	data []int
}

func (pq *PriorityQueue) Parent(index int) int {
	return (index - 1) / 2
}

func (pq *PriorityQueue) LeftChild(index int) int {
	return 2*index + 1
}

func (pq *PriorityQueue) RightChild(index int) int {
	return 2*index + 2
}

func (pq *PriorityQueue) Peek() (int, bool) {
	if pq.IsEmpty() {
		return -1, false
	}
	return pq.data[0], true
}

func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.data) == 0
}

func (pq *PriorityQueue) Push(value int) {
	pq.data = append(pq.data, value)
	pq.HeapifyUp(len(pq.data) - 1)
}

func (pq *PriorityQueue) Pop() (int, bool) {
	if pq.IsEmpty() {
		return -1, false
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

func (pq *PriorityQueue) Delete(value_to_delete int) bool {
	if pq.IsEmpty() {
		return false
	}
	for idx, val := range pq.data {
		if val == value_to_delete {
			pq.data[idx] = pq.data[len(pq.data)-1]
			pq.data = pq.data[:len(pq.data)-1]
			if idx >= len(pq.data) {
				return true
			}
			parent := pq.Parent(idx)
			if parent >= 0 && pq.data[parent] > pq.data[idx] {
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
		if pq.data[parent] > pq.data[index] {
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

	if left < size && pq.data[left] < pq.data[smallest] {
		smallest = left
	}

	if right < size && pq.data[right] < pq.data[smallest] {
		smallest = right
	}

	if smallest != index {
		pq.data[index], pq.data[smallest] = pq.data[smallest], pq.data[index]
		pq.HeapifyDown(smallest)
	}

}
