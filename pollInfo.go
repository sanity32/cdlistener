package cdlistener

type PollInfo[T comparable] struct {
	Iteration int
	LastValue T
	Logs      *History[T]
}
