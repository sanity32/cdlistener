package cdlistener

type Result[T any] struct {
	LastValue     T
	InterruptCode InterruptCode
	CdStopped     bool
	StoppedByUser bool
}

func (res Result[T]) Success() bool {
	return res.InterruptCode == SUCCESS
}
