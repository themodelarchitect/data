package database

import (
	"fmt"
	"sync"
)

type InMemoryKV[K int, V any] struct {
	data map[K]V
	mu   *sync.Mutex
}

func NewInMemoryKV[K int, V any]() *InMemoryKV[K, V] {
	data := make(map[K]V)
	return &InMemoryKV[K, V]{data: data, mu: &sync.Mutex{}}
}

func (kv *InMemoryKV[K, V]) Set(key K, val V) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = val
	return nil
}

func (kv *InMemoryKV[K, V]) Get(key K) (V, error) {
	if val, ok := kv.data[key]; ok {
		return val, nil
	}
	return *new(V), fmt.Errorf("not found")
}

func (kv *InMemoryKV[K, V]) Delete(key K) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, key)
	return nil
}
