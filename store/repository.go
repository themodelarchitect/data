package store

type Repository[T any] interface {
	Create(entity T) (int, error)
	Update(id int, entity T) error
	Get(id int) (T, error)
	GetAll() map[int]T
	Delete(id int) error
}
