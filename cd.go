package cdlistener

import (
	"context"
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
func New[T comparable](duration, repollInterval time.Duration, fnPoll func() T, fnScan func() InterruptCode) *Cd[T] {
	return &Cd[T]{
		duration:              duration,
		repollInterval:        repollInterval,
		host:                  NewHostEntity(fnPoll, func(pi PollInfo[T]) InterruptCode { return fnScan() }),
		maxSameRepollsToStall: DefaultMaxSameRepollsToStall,
		// C:                     DefaultResultChan[T](),

	}
}

func From[T comparable](host Host[T]) *Cd[T] {
	return &Cd[T]{
		host:                  host,
		duration:              DefaultDuration,
		repollInterval:        DefaultRepollInterval,
		maxSameRepollsToStall: DefaultMaxSameRepollsToStall,
		// C:                     DefaultResultChan[T](),
	}
}

// fnPoll:                fnPoll,
// fnPrematureInterrupt:  fnInterrupt,

type Cd[T comparable] struct {
	C                     chan Result[T]
	ctx                   context.Context
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

func (cd *Cd[T]) Ctx() context.Context {
	if r := cd.ctx; r != nil {
		return r
	}
	return context.Background()
}

func (cd *Cd[T]) SetCtx(ctx context.Context) *Cd[T] {
	cd.ctx = ctx
	return cd
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

// Deprecated: use Watch() with context instead
func (cd *Cd[T]) Start() *Cd[T] {
	cd.C = DefaultResultChan[T]()
	cd.startAt = time.Now()
	go func() {
		r := cd.poll()
		cd.finalize(r)
	}()
	return cd
}

func (cd *Cd[T]) Watch(ctx context.Context) <-chan Result[T] {
	cd.SetCtx(ctx)
	return cd.Start().C
}

func (cd *Cd[T]) poll() Result[T] {
	var iteration int
	history := NewHistory[T](cd.maxSameRepollsToStall)
	go func() {
		<-cd.ctx.Done()
		cd.Stop()
	}()
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
		if code := cd.host.CdScan(PollInfo[T]{
			Iteration: iteration,
			LastValue: last,
			Logs:      &history,
		}); code != 0 {
			return Result[T]{InterruptCode: code}
		}

		if history.StreakAll() {
			return Result[T]{CdStopped: true}
		}

		select {
		case <-cd.Ctx().Done():
			return Result[T]{CdStopped: true}
		case <-time.NewTimer(cd.repollInterval).C:
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
