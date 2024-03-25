package cdlistener

type Result[T any] struct {
	Value         T
	InterruptCode int
	CdStopped     bool
}
