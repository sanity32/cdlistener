package cdlistener

type CdValuer[T comparable] interface {
	CdValue() T
}

type CdInterrupter[T comparable] interface {
	CdInterrupt(PollInfo[T]) InterruptCode
}

type Host[T comparable] interface {
	CdInterrupter[T]
	CdValuer[T]
}

type BasicHoster[T comparable] struct {
	FnPoll      func() T
	FnInterrupt func(PollInfo[T]) InterruptCode
}

func NewBasicHoster[T comparable](fnPoll func() T, fnInterrupt func(PollInfo[T]) InterruptCode) *BasicHoster[T] {
	return &BasicHoster[T]{
		FnPoll:      fnPoll,
		FnInterrupt: fnInterrupt,
	}
}

func (b BasicHoster[T]) CdValue() T {
	return b.FnPoll()
}

func (b BasicHoster[T]) CdInterrupt(pollInfo PollInfo[T]) InterruptCode {
	return b.FnInterrupt(pollInfo)
}
