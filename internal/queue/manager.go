package queue

import "sync"

type LocalQueue struct {
	PQ   PriorityQueue
	Lock sync.Mutex
}

var DefaultLocalQueue = &LocalQueue{
	PQ: PriorityQueue{},
}
