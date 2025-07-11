package internal

type RingBuffer struct {
	buffer   []string
	capacity int
	start    int
	end      int
	size     int
	current  int
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		buffer:   make([]string, capacity),
		capacity: capacity,
		current:  -1,
	}
}

func (r *RingBuffer) Add(value string) {
	if value == "" {
		return
	}

	r.buffer[r.end] = value
	r.end = (r.end + 1) % r.capacity
	if r.size < r.capacity {
		r.size++
	} else {
		r.start = (r.start + 1) % r.capacity
	}
	r.current = -1
}

func (r *RingBuffer) Next() (string, bool) {
	if r.size == 0 || r.current == -1 {
		return "", false
	}
	next := (r.current + 1) % r.capacity

	if next == r.end {
		r.current = -1
		return "", false
	}

	r.current = next
	return r.buffer[r.current], true
}

func (r *RingBuffer) Prev() (string, bool) {
	if r.size == 0 {
		return "", false
	}
	if r.current == -1 {
		r.current = (r.end - 1 + r.capacity) % r.capacity
	} else if r.current != r.start {
		r.current = (r.current - 1 + r.capacity) % r.capacity
	} else {
		return r.buffer[r.current], true
	}
	return r.buffer[r.current], true
}

func (r *RingBuffer) Get(i int) string {
	if i < 0 || i >= r.size {
		return ""
	}
	return r.buffer[(r.start+i)%r.capacity]
}
