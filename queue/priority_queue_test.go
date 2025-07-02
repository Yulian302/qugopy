package queue

import (
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Push(5)
	pq.Push(2)
	pq.Push(3)
	pq.Push(1)

	if ok := pq.Delete(3); !ok {
		t.Error("priority queue could not delete value")
	}

	pq.Push(3)

	if val, ok := pq.Pop(); !ok || val != 1 {
		t.Errorf("priority queue pop value: %d, must be 1", val)
	}
	if val, ok := pq.Pop(); !ok || val != 2 {
		t.Errorf("priority queue pop value: %d, must be 2", val)
	}

	if val, ok := pq.Pop(); !ok || val != 3 {
		t.Errorf("priority queue pop value: %d, must be 3", val)
	}
	if val, ok := pq.Pop(); !ok || val != 5 {
		t.Errorf("priority queue pop value: %d, must be 5", val)
	}
}
