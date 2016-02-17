package throttlelock

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottleLock(t *testing.T) {
	tl := NewThrottleLock(10, 50)
	v := int64(0)
	for i := 0; i < 50; i++ {
		go func(t1 int) {
			tl.WaitForTurn()
			defer tl.Done()
			atomic.AddInt64(&v, int64(1))
		}(i)
	}
	tl.AwaitAll()
	if v != 50 {
		t.Errorf("Not all go routines ran. Value is %d", v)
	}
}

func TestThrottleLockWithDelayedWait(t *testing.T) {
	tl := NewThrottleLock(10, 50)
	v := int64(0)
	for i := 0; i < 50; i++ {
		go func(t1 int) {
			if t1 == 0 {
				// two seconds is more than enough time for 50
				// iterations to complete
				time.Sleep(time.Second * 2)
			}
			tl.WaitForTurn()
			defer tl.Done()
			atomic.AddInt64(&v, int64(1))
		}(i)
	}
	tl.AwaitAll()
	if v != 50 {
		t.Errorf("Not all go routines ran. Value is %d", v)
	}
}
