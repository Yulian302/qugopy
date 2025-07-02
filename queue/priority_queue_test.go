package queue

import (
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := &PriorityQueue{}
	pq.Push(IntTask{ID: 1, Type: "python", Payload: "Task1", Priority: 5})
	pq.Push(IntTask{ID: 2, Type: "go", Payload: "Task2", Priority: 2})
	pq.Push(IntTask{ID: 3, Type: "python", Payload: "Task3", Priority: 3})
	pq.Push(IntTask{ID: 4, Type: "go", Payload: "Task4", Priority: 1})

	if ok := pq.Delete(3); !ok {
		t.Error("priority queue could not delete value")
	}

	pq.Push(IntTask{ID: 5, Type: "python", Payload: "Task3", Priority: 3})

	if val, ok := pq.Pop(); !ok || val.Priority != 1 {
		t.Errorf("priority queue pop value: %d, must be 1", val.Priority)
	}
	if val, ok := pq.Pop(); !ok || val.Priority != 2 {
		t.Errorf("priority queue pop value: %d, must be 2", val.Priority)
	}

	if val, ok := pq.Pop(); !ok || val.Priority != 3 {
		t.Errorf("priority queue pop value: %d, must be 3", val.Priority)
	}
	if val, ok := pq.Pop(); !ok || val.Priority != 5 {
		t.Errorf("priority queue pop value: %d, must be 5", val.Priority)
	}
}
