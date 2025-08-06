package inmemoryQueue

import (
	"queue-broker/api/core"
	"sync"
	"time"
)

// Поскольку по заданию нежелательно использовать сторонние пакеты, то и соответственно вещи
// по типу "Кафка" использовать нельзя. Поэтому было разработано in-memory решение.
// При дальнейшей разработке можно реализовать что-то похожее, используя сторонние решения

type Queue struct {
	messages chan core.Message
}

func newQueue(maxSize int) *Queue {
	// По сути каналы - это естественные представители очереди.
	// Можно было бы построить решение чисто на использовании каналов
	// (но при расширении функционала пришлось бы менять функции, которые используют его +
	// у них была бы не единичная ответственность, поскольку требовался бы менеджмент как репозитория, так и входящих в него каналов).
	// Поэтому хотелось некоторый объект, к которому можно было обращаться и в дальнейшем расширять/изменять его поведение.
	return &Queue{
		messages: make(chan core.Message, maxSize),
	}
}

// put adds a message to the queue.
func (q *Queue) put(message core.Message) error {
	select {
	case q.messages <- message:
		return nil
	default:
		return core.ErrFullQueue
	}
}

// getWithTimeout retrieves or waits given timeout and removes the first message from the queue if success.
func (q *Queue) getWithTimeout(timeoutSec int) (core.Message, error) {
	select {
	case message := <-q.messages:
		return message, nil
	case <-time.After(time.Duration(timeoutSec) * time.Second):
		return core.Message{}, core.ErrQueueGetTimeout
	}
}

// get retrieves and removes the first message from the queue.
func (q *Queue) get() core.Message {
	select {
	case message := <-q.messages:
		return message
	}
}

type InMemoryQueueRepository struct {
	mu          sync.RWMutex
	queues      map[string]*Queue
	maxMessages int
	maxQueues   int
}

func NewInMemoryQueueRepository(maxMessages, maxQueues int) *InMemoryQueueRepository {
	return &InMemoryQueueRepository{
		queues:      make(map[string]*Queue),
		maxMessages: maxMessages,
		maxQueues:   maxQueues,
	}
}

// Put adds a message to the specified queue.
// Creates queue if it doesn't exist (if under queue limit).
func (r *InMemoryQueueRepository) Put(queueName string, message core.Message) error {
	// Здесь mutex используется для того, чтобы не создать несколько раз очередь с заданным именем.
	// Т.е. получение само по себе безопасно, но вот именно создавать несколько раз небезопасно, т.к. возможна такая ситуация:
	// Очереди нет, мы зашли. В одном вызове смотрим, что очереди нет, создаём и пишем сообщение.
	// В другом вызове, чуть позже (до создания в первом вызове) смотрим и видим, что очереди нет. Создаём и перезатираем
	// сообщение, которое записалось во время первого вызова.
	r.mu.Lock()

	queue, exists := r.queues[queueName]
	if !exists {
		if len(r.queues) >= r.maxQueues {
			r.mu.Unlock()
			return core.ErrQueuesLimit
		}

		r.queues[queueName] = newQueue(r.maxMessages)
		queue = r.queues[queueName]
	}

	r.mu.Unlock()

	return queue.put(message)
}

// Get retrieves a message from the specified queue
func (r *InMemoryQueueRepository) Get(queueName string, timeoutSec int) (core.Message, error) {
	queue, exists := r.queues[queueName]

	if !exists {
		return core.Message{}, core.ErrQueueNotFound
	}

	if timeoutSec > 0 {
		return queue.getWithTimeout(timeoutSec)
	}
	return queue.get(), nil
}
