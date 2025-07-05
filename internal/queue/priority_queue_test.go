package queue

import (
	"os"
	"testing"

	"github.com/Yulian302/qugopy/models"
	"github.com/google/uuid"
)

var (
	q     *PriorityQueue
	tasks []*models.IntTask
)

func TestMain(m *testing.M) {
	q = &PriorityQueue{}
	tasks = []*models.IntTask{
		{ID: uuid.New().String(), Task: models.Task{Type: "python", Payload: "Task1", Priority: 5}},
		{ID: uuid.New().String(), Task: models.Task{Type: "go", Payload: "Task2", Priority: 2}},
		{ID: uuid.New().String(), Task: models.Task{Type: "python", Payload: "Task3", Priority: 3}},
		{ID: uuid.New().String(), Task: models.Task{Type: "go", Payload: "Task4", Priority: 1}},
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

	if val, ok := q.Pop(); !ok || val.Task.Priority != 1 {
		t.Errorf("priority queue pop value: %d, must be 1", val.Task.Priority)
	}
	if val, ok := q.Pop(); !ok || val.Task.Priority != 2 {
		t.Errorf("priority queue pop value: %d, must be 2", val.Task.Priority)
	}

	if val, ok := q.Pop(); !ok || val.Task.Priority != 3 {
		t.Errorf("priority queue pop value: %d, must be 3", val.Task.Priority)
	}
	if val, ok := q.Pop(); !ok || val.Task.Priority != 5 {
		t.Errorf("priority queue pop value: %d, must be 5", val.Task.Priority)
	}
}

func TestDelete(t *testing.T) {
	for _, task := range tasks {
		q.Push(*task)
	}

	if ok := q.Delete(1); !ok {
		t.Errorf("could not delete value 1")
	}

	intTask, ok := q.Peek()

	if !ok {
		t.Errorf("could not peek root element")
	}
	if intTask.Task.Priority != 2 {
		t.Errorf("wrong order after deleting root elem. Got %d, must be %d", intTask.Task.Priority, 2)
	}

	if ok := q.Delete(3); !ok {
		t.Errorf("could not delete value 3")
	}

	intTask, ok = q.Peek()

	if !ok {
		t.Errorf("could not peek root element")
	}
	if intTask.Task.Priority != 2 {
		t.Errorf("wrong order after deleting middle elem. Got %d, must be %d", intTask.Task.Priority, 2)
	}
	if q.data[len(q.data)-1].Task.Priority != 5 {
		t.Errorf("wrong order after deleting middle elem. Got %d, must be %d", q.data[len(q.data)-1].Task.Priority, 5)
	}

}
