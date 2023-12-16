package memory

import (
	"fmt"
	"github.com/jkittell/data/store"
	"github.com/jkittell/data/structures/array"
	"math/rand"
	"sync"
	"time"
)

type Repository[T store.Entity] struct {
	data map[int]T
	mu   *sync.Mutex
}

func New[T store.Entity]() *Repository[T] {
	data := make(map[int]T)
	return &Repository[T]{data, &sync.Mutex{}}
}

func (r *Repository[T]) Create(entity T) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	entity.SetID(id)

	r.data[id] = entity

	return id, nil
}

func (r *Repository[T]) Update(id int, entity T) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[id] = entity
	return nil
}

func (r *Repository[T]) Get(id int) (T, error) {
	if val, ok := r.data[id]; ok {
		return val, nil
	}
	return *new(T), fmt.Errorf("not found")
}

func (r *Repository[T]) GetAll() *array.Array[T] {
	result := array.New[T]()
	for _, v := range r.data {
		result.Push(v)
	}
	return result
}

func (r *Repository[T]) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("not found")
	}

	delete(r.data, id)
	return nil
}
