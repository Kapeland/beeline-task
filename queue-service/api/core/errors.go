package core

import "errors"

var (
	ErrQueueNotFound   = errors.New("queue not found")
	ErrQueuesLimit     = errors.New("queues limit exceeded")
	ErrQueueGetTimeout = errors.New("queue get timeout")
	ErrFullQueue       = errors.New("queue is full")
)
