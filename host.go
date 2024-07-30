package cdlistener

type CdValuer[T comparable] interface {
	CdValue() T
}

type CdScanner[T comparable] interface {
	CdScan(PollInfo[T]) InterruptCode
}

type Host[T comparable] interface {
	CdScanner[T]
	CdValuer[T]
}

type HostEntity[T comparable] struct {
	FnPoll func() T
	FnScan func(PollInfo[T]) InterruptCode
}

func NewHostEntity[T comparable](fnPoll func() T, fnScan func(PollInfo[T]) InterruptCode) *HostEntity[T] {
	return &HostEntity[T]{
		FnPoll: fnPoll,
		FnScan: fnScan,
	}
}

func (b HostEntity[T]) CdValue() T {
	return b.FnPoll()
}

func (b HostEntity[T]) CdScan(pollInfo PollInfo[T]) InterruptCode {
	return b.FnScan(pollInfo)
}
