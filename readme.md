ThrottleLock
------------
Throttle lock is a limited use case library that obfuscates the complexity of scheduling A LOT (we're talking millions) of go routines at once, but limiting how many of them can run at one time. Unit tests, scraping tasks, and lots of io work are examples of good fits for this library.

Channels are not good for this kind of thing as they timeout if they wait too long. Using waitgroups and conditionals is really the only path forward as they put go routines to sleep (which is necessary if your go routine isn't going to run hours in the future).

Throttle lock is working for my own use cases, but be forwarned! It contains a lot of synchronization logic in a small amount of code, it could do anything from stealing your lunch money to race locking a production server. Be careful!!!!!

Obligatory example:
```go
	// run one million routines, only 100 at a time.
	tl := NewThrottleLock(100, 1000000)
	for i := 0; i < 1000000; i++ {
		go func(t1 int) {
			tl.WaitForTurn()
			defer tl.Done()
			// do something
		}(i)
	}
	tl.AwaitAll()
```