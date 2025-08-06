package core

type QueueRepo interface {
	Put(queueName string, msg Message) error
	Get(queueName string, timeoutSec int) (Message, error)
}
