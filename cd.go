package cdlistener

import (
	"time"
)

func New[T comparable](duration, repollInterval time.Duration, fnPoll func() T, fnInterrupt func() InterruptCode) *Cd[T] {
	return &Cd[T]{
		duration:              duration,
		repollInterval:        repollInterval,
		C:                     make(chan Result[T], 1),
		fnPoll:                fnPoll,
		fnPrematureInterrupt:  fnInterrupt,
		MaxSameRepollsToStall: 5,
	}
}

type Cd[T comparable] struct {
	C                     chan Result[T]
	MaxSameRepollsToStall int
	lastResult            Result[T]
	startAt               time.Time
	duration              time.Duration
	repollInterval        time.Duration
	fnPoll                func() T
	fnPrematureInterrupt  func() InterruptCode
	stopper               _Stopper
	finalized             bool
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
	history := NewHistory[T](cd.MaxSameRepollsToStall)
	for {
		iteration++
		if cd.stopper.Flag {
			cd.finalize(Result[T]{StoppedByUser: true})
		}

		if cd.finalized {
			return cd.lastResult
		}

		if code := cd.fnPrematureInterrupt(); code != 0 {
			return Result[T]{InterruptCode: code}
		}

		history.Push(cd.fnPoll())
		if history.StreakAll() {
			return Result[T]{CdStopped: true}
		}

		time.Sleep(cd.repollInterval)
	}
}
