package cdlistener

import "fmt"

const MIN_HISTORY_LENGTH = 5

func NewHistory[T comparable](n int) History[T] {
	return make(History[T], n)
}

type History[T comparable] []T

func (h History[T]) Zero() (t T) {
	return
}

func (h *History[T]) Get(n int) T {
	return (*h)[n]
}

func (h History[T]) Set(n int, v T) {
	h[n] = v
}

func (h History[T]) Last() T {
	if len(h) == 0 {
		return h.Zero()
	}
	return h[0]
}

func (h *History[T]) init() {
	if len(*h) < MIN_HISTORY_LENGTH {
		fmt.Println("min history", MIN_HISTORY_LENGTH)
		r := make(History[T], MIN_HISTORY_LENGTH)
		copy(r, *h)
		*h = r
	}
}

func (h History[T]) Push(v T) History[T] {
	h.init()
	for i := len(h) - 1; i > 0; i-- {
		h[i] = h[i-1]
	}
	h[0] = v
	return h
}

func (h History[T]) IsLastZero() bool {
	h.init()
	return h.Last() == h.Zero()
}

func (h History[T]) Streak(n int) bool {
	l := len(h)
	if n >= l {
		n = l - 1
	}
	for i := 1; i < n; i++ {
		if h[i] != h[i-1] {
			return false
		}
	}
	return true
}

func (h History[T]) StreakAll() bool {
	return h.Streak(len(h) - 1)
}
