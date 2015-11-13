package throttlelock

import "sync"

type ThrottleLock struct {
	C *sync.Cond
	W *sync.WaitGroup
	Cap,
	Cur int
}

func NewThrottleLock(capacity int) *ThrottleLock {
	tl := new(ThrottleLock)
	tl.C = sync.NewCond(&sync.Mutex{})
	tl.W = &sync.WaitGroup{}
	tl.Cap = capacity
	tl.Cur = 0
	return tl
}

func (tl *ThrottleLock) AwaitAll() {
	tl.W.Wait()
}

func (tl *ThrottleLock) WaitForTurn() {
	tl.W.Add(1)
	tl.C.L.Lock()
	if tl.Cur == tl.Cap {
		tl.C.Wait()
	}
	tl.Cur++
	tl.C.L.Unlock()
}

func (tl *ThrottleLock) Done() {
	tl.C.L.Lock()
	tl.Cur--
	tl.C.L.Unlock()
	tl.C.Signal()
	tl.W.Done()
}
