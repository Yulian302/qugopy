package queue

import (
	"os"
	"testing"

	"github.com/Yulian302/qugopy/models"
)

var (
	q     *PriorityQueue
	tasks []*models.IntTask
)

func TestMain(m *testing.M) {
	q = &PriorityQueue{}
	tasks = []*models.IntTask{
		{ID: 1, Task: models.Task{Type: "python", Payload: "Task1", Priority: 5}},
		{ID: 2, Task: models.Task{Type: "go", Payload: "Task2", Priority: 2}},
		{ID: 3, Task: models.Task{Type: "python", Payload: "Task3", Priority: 3}},
		{ID: 4, Task: models.Task{Type: "go", Payload: "Task4", Priority: 1}},
	}
	code := m.Run()
	os.Exit(code)
}

func TestPushPop(t *testing.T) {

	for _, task := range tasks {
		q.Push(*task)
	}
	size := len(q.data)
	if size != len(tasks) {
		t.Errorf("size of queue is wrong. Must be %d, got %d", len(tasks), size)
	}

	if val, ok := q.Pop(); !ok || val.Priority != 1 {
		t.Errorf("priority queue pop value: %d, must be 1", val.Priority)
	}
	if val, ok := q.Pop(); !ok || val.Priority != 2 {
		t.Errorf("priority queue pop value: %d, must be 2", val.Priority)
	}

	if val, ok := q.Pop(); !ok || val.Priority != 3 {
		t.Errorf("priority queue pop value: %d, must be 3", val.Priority)
	}
	if val, ok := q.Pop(); !ok || val.Priority != 5 {
		t.Errorf("priority queue pop value: %d, must be 5", val.Priority)
	}
}

func TestDelete(t *testing.T) {
	for _, task := range tasks {
		q.Push(*task)
	}

	if ok := q.Delete(1); !ok {
		t.Errorf("could not delete value 1")
	}

	task, ok := q.Peek()

	if !ok {
		t.Errorf("could not peek root element")
	}
	if task.Priority != 2 {
		t.Errorf("wrong order after deleting root elem. Got %d, must be %d", task.Priority, 2)
	}

	if ok := q.Delete(3); !ok {
		t.Errorf("could not delete value 3")
	}

	task, ok = q.Peek()

	if !ok {
		t.Errorf("could not peek root element")
	}
	if task.Priority != 2 {
		t.Errorf("wrong order after deleting middle elem. Got %d, must be %d", task.Priority, 2)
	}
	if q.data[len(q.data)-1].Priority != 5 {
		t.Errorf("wrong order after deleting middle elem. Got %d, must be %d", q.data[len(q.data)-1].Priority, 5)
	}

}
