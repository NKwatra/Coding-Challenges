package utils

type Comparable[T any] interface {
	CompareTo(other T) int
}
