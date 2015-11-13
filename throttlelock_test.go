package throttlelock

import (
	"fmt"
	"testing"
)

func TestThrottleLock(t *testing.T) {
	tl := NewThrottleLock(10)
	for i := 0; i < 50; i++ {
		go func(t int) {
			tl.WaitForTurn()
			defer tl.Done()
			fmt.Println("I'm job #", t)
		}(i)
	}
	tl.AwaitAll()
}
