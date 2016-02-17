package throttlelock

import "sync"

type ThrottleLock struct {
	turnCondition *sync.Cond
	waitGroup     *sync.WaitGroup
	runCapacity,
	currentRunning int
}

func NewThrottleLock(runCapacity, totalExpectedWaiters int) *ThrottleLock {
	tl := new(ThrottleLock)
	tl.turnCondition = sync.NewCond(&sync.Mutex{})
	tl.waitGroup = &sync.WaitGroup{}
	tl.runCapacity = runCapacity
	tl.currentRunning = 0
	tl.waitGroup.Add(totalExpectedWaiters)
	return tl
}

func (tl *ThrottleLock) AwaitAll() {
	tl.waitGroup.Wait()
}

func (tl *ThrottleLock) WaitForTurn() {
	tl.turnCondition.L.Lock()
	if tl.currentRunning == tl.runCapacity {
		tl.turnCondition.Wait()
	}
	tl.currentRunning++
	tl.turnCondition.L.Unlock()
}

func (tl *ThrottleLock) Done() {
	tl.turnCondition.L.Lock()
	tl.currentRunning--
	tl.turnCondition.L.Unlock()
	tl.turnCondition.Signal()
	tl.waitGroup.Done()
}
