package models

// Mode specified by user. Can be either `redis` or `local`. If `redis` mode is specified, all tasks are pushed to the Redis data store. On the other hand, if `local` is specified, in-memory priority queue is used.
type Mode string

// Checks whether the mode is of valid type. Can be either `redis` or `local`
func (m Mode) IsValid() bool {
	return m == "redis" || m == "local"
}
