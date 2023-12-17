package store

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type InMemoryRepository[T any] struct {
	data map[int]T
	mu   *sync.Mutex
}

func NewInMemoryRepository[T any]() *InMemoryRepository[T] {
	data := make(map[int]T)
	return &InMemoryRepository[T]{data, &sync.Mutex{}}
}

func (r *InMemoryRepository[T]) Create(entity T) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	r.data[id] = entity
	return id, nil
}

func (r *InMemoryRepository[T]) Update(id int, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; ok {
		r.data[id] = entity
		return nil
	} else {
		return errors.New(fmt.Sprintf("%d not found", id))
	}
}

func (r *InMemoryRepository[T]) Get(id int) (T, error) {
	if val, ok := r.data[id]; ok {
		return val, nil
	}
	return *new(T), fmt.Errorf("not found")
}

func (r *InMemoryRepository[T]) GetAll() map[int]T {
	return r.data
}

func (r *InMemoryRepository[T]) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("%d not found", id)
	}

	delete(r.data, id)
	return nil
}
