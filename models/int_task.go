package models

// IntTask (internal task) represents a task with priority-based ordering capabilities.
// It's designed for use in priority queues where tasks are ordered by: priority
type IntTask struct {
	// User defined task (omits internal properties)
	Task

	// ID uniquely identifies the task. Used for equality comparisons.
	ID int `json:"id"`
}

// GT (Greater Than) compares task priorities.
// Returns true if t1 has higher priority than t2.
//
// Example:
//
//	task1 := &IntTask{Priority: 5}
//	task2 := &IntTask{Priority: 3}
//	task1.GT(task2) // true
func (t1 *IntTask) GT(t2 *IntTask) bool {
	return t1.Priority > t2.Priority
}

// GTE (Greater Than or Equal) compares task priorities.
// Returns true if t1 has equal or higher priority than t2.
func (t1 *IntTask) GTE(t2 *IntTask) bool {
	return t1.Priority >= t2.Priority
}

// LT (Less Than) compares task priorities.
// Returns true if t1 has lower priority than t2.
func (t1 *IntTask) LT(t2 *IntTask) bool {
	return t1.Priority < t2.Priority
}

// LTE (Less Than or Equal) compares task priorities.
// Returns true if t1 has equal or lower priority than t2.
func (t1 *IntTask) LTE(t2 *IntTask) bool {
	return t1.Priority <= t2.Priority
}

// EQ (Equal) checks task identity.
// Uses ID field rather than Priority for equality comparison.
// Returns true if tasks have the same ID.
//
// Note: This differs from priority-based comparisons as it's
// used for task identification, not ordering.
func (t1 *IntTask) EQ(t2 *IntTask) bool {
	return t1.ID == t2.ID // Fixed: Changed <= to == for proper equality
}
