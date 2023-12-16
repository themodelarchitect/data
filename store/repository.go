package store

import (
	"github.com/jkittell/data/structures/array"
)

type Repository[T Entity] interface {
	Create(entity T) (int, error)
	Update(id int, entity T) error
	Get(id int) (T, error)
	GetAll() *array.Array[T]
	Delete(id int) error
}

type Entity interface {
	GetID() int
	SetID(id int)
}
