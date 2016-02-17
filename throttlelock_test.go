package throttlelock

import (
	"testing"
	"time"
)

func TestThrottleLock(t *testing.T) {
	tl := NewThrottleLock(10, 50)
	for i := 0; i < 50; i++ {
		go func(t1 int) {
			tl.WaitForTurn()
			defer tl.Done()
		}(i)
	}
	tl.AwaitAll()
}

func TestThrottleLockWithDelayedWait(t *testing.T) {
	tl := NewThrottleLock(10, 50)
	for i := 0; i < 50; i++ {
		go func(t1 int) {
			if t1 == 0 {
				time.Sleep(time.Second * 10)
			}
			tl.WaitForTurn()
			defer tl.Done()
		}(i)
	}
	tl.AwaitAll()
}
