package cdlistener

type Result[T any] struct {
	Value         T
	InterruptCode InterruptCode
	CdStopped     bool
	StoppedByUser bool
}
