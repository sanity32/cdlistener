package cdlistener

type InterruptCode int

const (
	SUCCESS InterruptCode = 1
	FAILURE InterruptCode = -1
	NONE    InterruptCode = 0
)
