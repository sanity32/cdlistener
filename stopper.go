package cdlistener

type _Stopper struct {
	Flag bool
	Ch   chan any
}
