package utils

import "testing"

type MyInt struct {
	val int
}

func (m MyInt) CompareTo(o MyInt) int {
	diff := m.val - o.val
	switch {
	case diff == 0:
		return 0
	case diff > 0:
		return 1
	default:
		return -1
	}
}

func TestPeekOnEmpty(t *testing.T) {
	q := MinPQueue[MyInt]{make([]MyInt, 0)}
	val, err := q.Peek()
	if err == nil {
		t.Errorf("expected error with peek on empty queue, received %v", val)
	}
}

func TestPollOnEmpty(t *testing.T) {
	q := MinPQueue[MyInt]{make([]MyInt, 0)}
	_, err := q.Poll()
	if err == nil {
		t.Errorf("expected error with poll on empty queue")
	}
}

func TestMinQueue(t *testing.T) {
	q := MinPQueue[MyInt]{make([]MyInt, 0)}
	q.Add(MyInt{32})
	q.Add(MyInt{42})
	q.Add(MyInt{120})
	q.Add(MyInt{7})
	q.Add(MyInt{42})
	q.Add(MyInt{24})
	q.Add(MyInt{37})
	q.Add(MyInt{2})
	value, err := q.Peek()
	if err != nil || value.val != 2 {
		t.Errorf("expected err=nil, received=%v and expected value=2, received=%d", err, value.val)
	}
	order := [8]int{2, 7, 24, 32, 37, 42, 42, 120}
	for i := 0; i < 8; i++ {
		value, err = q.Poll()
		if err != nil || value.val != order[i] {
			t.Errorf("expected err=nil, received=%v and expected value=%d, received=%d", err, order[i], value.val)
		}
	}
}
