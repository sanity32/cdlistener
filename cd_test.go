package cdlistener

import (
	"testing"
	"time"
)

func TestCd_Expired(t *testing.T) {
	i := 1
	var interruptCode int

	time.AfterFunc(time.Second*3, func() {
		interruptCode = -3
	})

	cd := New[int](
		time.Second*30,
		time.Second,
		func() int { return i },
		func() int { return interruptCode },
	).Start()

	r := <-cd.C
	t.Logf("%#v", r)
}
