package websocket

import (
	"sync"
	"sync/atomic"
)

type ClientMap[T any] struct {
	clients *sync.Map
}

func NewClientMap[T any]() *ClientMap[T] {
	return &ClientMap[T]{
		clients: &sync.Map{},
	}
}

func (cm *ClientMap[T]) Add(key string, value T) {
	cm.clients.Store(key, value)
}

func (cm *ClientMap[T]) Get(key string) (T, bool) {
	value, ok := cm.clients.Load(key)
	if !ok {
		var zero T
		return zero, false
	}
	return value.(T), true
}

func (cm *ClientMap[T]) Delete(key string) {
	cm.clients.Delete(key)
}

func (cm *ClientMap[T]) Exists(key string) bool {
	_, exists := cm.clients.Load(key)
	return exists
}

func (cm *ClientMap[T]) Range(f func(key string, value T) bool) {
	cm.clients.Range(func(k, v interface{}) bool {
		return f(k.(string), v.(T))
	})
}

func (cm *ClientMap[T]) Len() int {
	var length int64
	cm.clients.Range(func(k, v interface{}) bool {
		atomic.AddInt64(&length, 1)
		return true
	})
	return int(length)
}
