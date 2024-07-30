package cdlistener

type Result[T any] struct {
	LastValue     T
	InterruptCode InterruptCode
	CdStopped     bool
	StoppedByUser bool
}
