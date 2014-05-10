package tb

import (
	"errors"
	"sync"
	"time"
)

type TokenBucket struct {
	max     int
	cur     int
	stop    chan struct{}
	stopped chan struct{}
	lock    sync.Mutex
	running bool
}

// Get a token from the bucket, if available. Returns true
// if a token was subtracted from the bucket, or false
func (t *TokenBucket) Get() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.cur > 0 {
		t.cur--
		return true
	}
	return false
}

// See if a token is available without taking a token from
// the bucket. This does not gaurantee that the next time
// you use Get that you will have a token if other goroutines
// are accessing the bucket concurrently
func (t *TokenBucket) Peek() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.cur > 0 {
		return true
	}
	return false
}

// Set the maximum number of tokens the bucket can have
func (t *TokenBucket) Max(max int) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.max = max
}

// Set the number of tokens in the bucket. This is useful for
// initialization in combination with Max() before Start() is
// called
func (t *TokenBucket) Set(num int) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.cur = num
}

// Begin adding one token to the bucket every interval. Tokens
// will not be added past the last setting set by a call to Max()
func (t *TokenBucket) Start(d time.Duration) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if true == t.running {
		return errors.New("Bucket already running")
	}
	t.stop = make(chan struct{})
	t.stopped = make(chan struct{})
	go func(t *TokenBucket, d time.Duration) {
		ticker := time.NewTicker(d)
		for {
			select {
			case <-ticker.C:
				t.lock.Lock()
				if t.max > t.cur {
					t.cur++
				}
				t.lock.Unlock()
			case <-t.stop:
				close(t.stopped)
				return
			}
		}
	}(t, d)
	t.running = true
	return nil
}

// Stop adding tokens to the bucket. The bucket can be resumed
// with a call to Start()
func (t *TokenBucket) Stop() error {
	if false == t.running {
		return errors.New("Bucket already stopped")
	}
	close(t.stop)
	<-t.stopped
	t.running = false
	return nil
}
