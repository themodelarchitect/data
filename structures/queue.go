package structures

type Queue[T any] struct {
	items *Array[T]
}

func (q *Queue[T]) Enqueue(item T) {
	q.items.Push(item)
}

func (q *Queue[T]) Dequeue() T {
	if q.items.Length() < 1 {
		var val T
		return val // returns nil if the queue is empty
	}
	item := q.items.Pop()
	return item
}
