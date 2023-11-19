package utils

import "fmt"

type MinPQueue[T Comparable[T]] struct {
	data []T
}

func left(index int) int {
	return 2*index + 1
}

func right(index int) int {
	return 2*index + 2
}

func parent(index int) int {
	return (index - 1) / 2
}

func heapify[U Comparable[U]](q *MinPQueue[U], root int) {
	smallest := root
	leftChild := left(root)
	rightChild := right(root)

	if leftChild < len(q.data) && q.data[leftChild].CompareTo(q.data[root]) == -1 {
		smallest = leftChild
	}
	if rightChild < len(q.data) && q.data[rightChild].CompareTo(q.data[smallest]) == -1 {
		smallest = rightChild
	}

	if smallest != root {
		temp := q.data[smallest]
		q.data[smallest] = q.data[root]
		q.data[root] = temp
		heapify[U](q, smallest)
	}
}

func (q *MinPQueue[T]) Add(item T) {
	q.data = append(q.data, item)
	index := len(q.data) - 1
	for index > 0 {
		parentIndex := parent(index)
		if q.data[index].CompareTo(q.data[parentIndex]) == -1 {
			temp := q.data[index]
			q.data[index] = q.data[parentIndex]
			q.data[parentIndex] = temp
			index = parentIndex
		} else {
			break
		}
	}
}

func (q *MinPQueue[T]) Peek() (T, error) {
	if len(q.data) == 0 {
		return *new(T), fmt.Errorf("cannot peek an empty queue")
	}
	return q.data[0], nil
}

func (q *MinPQueue[T]) Poll() (T, error) {
	if len(q.data) == 0 {
		return *new(T), fmt.Errorf("cannot poll an empty queue")
	}
	element := q.data[0]
	q.data[0] = q.data[len(q.data)-1]
	q.data = q.data[:len(q.data)-1]
	heapify[T](q, 0)
	return element, nil
}

func (q *MinPQueue[T]) Size() int {
	return len(q.data)
}

func NewMinPQueue[T Comparable[T]]() *MinPQueue[T] {
	q := new(MinPQueue[T])
	q.data = make([]T, 0)
	return q
}
