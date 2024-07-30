package cdlistener

import (
	"time"
)

var (
	DefaultDuration              = time.Minute * 3
	DefaultRepollInterval        = time.Second * 2
	DefaultMaxSameRepollsToStall = 7
)

func DefaultResultChan[T comparable]() chan Result[T] {
	return make(chan Result[T], 1)
}

// @deprecated, From() method is preferred
func New[T comparable](duration, repollInterval time.Duration, fnPoll func() T, fnInterrupt func() InterruptCode) *Cd[T] {
	return &Cd[T]{
		duration:              duration,
		repollInterval:        repollInterval,
		C:                     DefaultResultChan[T](),
		host:                  BasicHoster[T]{FnPoll: fnPoll, FnInterrupt: func(pi PollInfo[T]) InterruptCode { return fnInterrupt() }},
		maxSameRepollsToStall: DefaultMaxSameRepollsToStall,
	}
}

func From[T comparable](host Host[T]) *Cd[T] {
	return &Cd[T]{
		host:                  host,
		duration:              DefaultDuration,
		repollInterval:        DefaultRepollInterval,
		maxSameRepollsToStall: DefaultMaxSameRepollsToStall,
		C:                     DefaultResultChan[T](),
	}
}

// fnPoll:                fnPoll,
// fnPrematureInterrupt:  fnInterrupt,

type Cd[T comparable] struct {
	C                     chan Result[T]
	maxSameRepollsToStall int
	lastResult            Result[T]
	startAt               time.Time
	duration              time.Duration
	repollInterval        time.Duration
	// fnPoll                func() T
	// fnPrematureInterrupt  func() InterruptCode
	host      Host[T]
	stopper   _Stopper
	finalized bool
}

func (cd *Cd[T]) Stop() {
	if !cd.finalized {
		cd.stopper.Flag = true
		cd.stopper.Ch = make(chan any, 1)
		<-cd.stopper.Ch
		cd.finalized = true
	}
}

func (cd *Cd[T]) finalize(r Result[T]) {
	cd.finalized = true
	cd.lastResult = r
	defer func() { recover() }()
	cd.C <- r
}

func (cd Cd[T]) Expired() bool {
	return time.Now().After(cd.startAt.Add(cd.duration))
}

func (cd *Cd[T]) Start() *Cd[T] {
	cd.startAt = time.Now()
	go func() {
		r := cd.poll()
		cd.finalize(r)
	}()
	return cd
}

func (cd *Cd[T]) poll() Result[T] {
	var iteration int
	history := NewHistory[T](cd.maxSameRepollsToStall)
	for {
		iteration++
		if cd.stopper.Flag {
			cd.finalize(Result[T]{StoppedByUser: true})
		}

		if cd.finalized {
			return cd.lastResult
		}

		// last := cd.fnPoll()
		last := cd.host.CdValue()
		history.Push(last)

		// if code := cd.fnPrematureInterrupt(); code != 0
		if code := cd.host.CdInterrupt(PollInfo[T]{
			Iteration: iteration,
			LastValue: last,
			Logs:      &history,
		}); code != 0 {
			return Result[T]{InterruptCode: code}
		}

		if history.StreakAll() {
			return Result[T]{CdStopped: true}
		}

		time.Sleep(cd.repollInterval)
	}
}

func (cd *Cd[T]) MaxSameRepollsToStall(value int) *Cd[T] {
	cd.maxSameRepollsToStall = value
	return cd
}

func (cd *Cd[T]) Duration(value time.Duration) *Cd[T] {
	cd.duration = value
	return cd
}

func (cd *Cd[T]) RepollInterval(value time.Duration) *Cd[T] {
	cd.repollInterval = value
	return cd
}
