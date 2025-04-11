package models

// BatchStore is a efficient store for batches of data
// It is a slice of pointers to the data
type BatchStore[T any] struct {
	store []T
	counter int
	size int
}

// NewBatchStore creates a new batch store
func NewBatchStore[T any](size int) BatchStore[T] {
	return BatchStore[T]{
		store: make([]T, 0, size),
		counter: 0,
		size: size,
	}
}

// Reset resets the batch store
func (b *BatchStore[T]) GetAndReset() ([]T, error) {
	data := b.store
	b.store = make([]T, 0, b.size)
	b.counter = 0
	return data, nil
}

// Push pushes the batch store to the server
func (b *BatchStore[T]) Push(data T) (*[]T, error) {
	b.store = append(b.store, data)
	b.counter++
	if b.counter >= b.size {
		data, err := b.GetAndReset()
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return nil, nil
}
