package throttlelock

import "sync"

type ThrottleLock struct {
	turnCondition,
	awaitAllCondition *sync.Cond
	waitGroup *sync.WaitGroup
	runCapacity,
	currentRunning,
	totalExpectedWaiters,
	currentWaiters uint64
	awaitAllIsSleeping bool
}

func NewThrottleLock(runCapacity, totalExpectedWaiters uint64) *ThrottleLock {
	tl := new(ThrottleLock)
	tl.turnCondition = sync.NewCond(&sync.Mutex{})
	tl.awaitAllCondition = sync.NewCond(&sync.Mutex{})
	tl.waitGroup = &sync.WaitGroup{}
	tl.runCapacity = runCapacity
	tl.currentRunning = uint64(0)
	tl.totalExpectedWaiters = totalExpectedWaiters
	tl.currentWaiters = uint64(0)
	tl.awaitAllIsSleeping = false
	return tl
}

func (tl *ThrottleLock) AwaitAll() {
	tl.awaitAllCondition.L.Lock()
	if tl.currentWaiters == tl.totalExpectedWaiters {
		tl.awaitAllCondition.L.Unlock()
	} else {
		tl.awaitAllIsSleeping = true
		tl.awaitAllCondition.Wait()
	}
	tl.waitGroup.Wait()
}

func (tl *ThrottleLock) WaitForTurn() {
	tl.waitGroup.Add(1)
	tl.awaitAllCondition.L.Lock()
	tl.currentWaiters++
	if tl.currentWaiters == tl.totalExpectedWaiters && tl.awaitAllIsSleeping {
		tl.awaitAllCondition.Signal()
	}
	tl.awaitAllCondition.L.Unlock()
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
