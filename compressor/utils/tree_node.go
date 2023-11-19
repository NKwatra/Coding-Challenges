package utils

type TreeNode[T any] struct {
	val    T
	left   *TreeNode[T]
	right  *TreeNode[T]
	weight int
}

func (t TreeNode[T]) CompareTo(node TreeNode[T]) int {
	diff := t.weight - node.weight
	switch {
	case diff == 0:
		return 0
	case diff > 0:
		return 1
	default:
		return -1
	}
}

func (t *TreeNode[T]) SetVal(val T) {
	t.val = val
}

func (t *TreeNode[T]) SetLeft(left *TreeNode[T]) {
	t.left = left
}

func (t *TreeNode[T]) SetRight(right *TreeNode[T]) {
	t.right = right
}

func (t *TreeNode[T]) Left() *TreeNode[T] {
	return t.left
}

func (t *TreeNode[T]) Right() *TreeNode[T] {
	return t.right
}

func (t *TreeNode[T]) Val() T {
	return t.val
}

func (t *TreeNode[T]) Weight() int {
	return t.weight
}

func NewTreeNode[T any](weight int) *TreeNode[T] {
	node := new(TreeNode[T])
	node.weight = weight
	return node
}
