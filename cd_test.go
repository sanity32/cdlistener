package cdlistener

import (
	"testing"
	"time"
)

func TestCd_Expired(t *testing.T) {
	i := 1
	var interruptCode int

	cd := New[int](
		time.Second*30,
		time.Second,
		func() int { return i },
		func() int { return interruptCode },
	).Start()

	time.AfterFunc(time.Second*1, func() {
		// interruptCode = -3
		t.Log("Stopping ...")
		cd.Stop()
		t.Log("Stop ok!")
	})

	r := <-cd.C
	t.Logf("%#v", r)
}
