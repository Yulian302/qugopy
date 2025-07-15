package queue

import "sync"

type LocalQueue struct {
	PQ   PriorityQueue
	Lock sync.Mutex
}

var (
	PythonLocalQueue = &LocalQueue{
		PQ: PriorityQueue{},
	}
	GoLocalQueue = &LocalQueue{
		PQ: PriorityQueue{},
	}
)
