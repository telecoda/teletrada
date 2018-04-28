package domain

import (
	"sync"
	"time"
)

var ServerTime = realTime

var fakeTimeLock sync.RWMutex
var fakeTime time.Time = time.Now().UTC() // default to a reasonable value

var realTime = func() time.Time { return time.Now().UTC() }

// fakeServerTime - returns fakeTime to allow tests to be deterministic
func fakeServerTime() time.Time {
	fakeTimeLock.RLock()
	defer fakeTimeLock.RUnlock()
	return fakeTime
}

// UseRealTime - reset ServerTime func to use the real time
func UseRealTime() {
	ServerTime = realTime
}

// UseFakeTime - return value of fakeTime variable
func UseFakeTime() {
	ServerTime = fakeServerTime
}

func TickFakeTime(d time.Duration) {
	fakeTimeLock.Lock()
	fakeTime = fakeTime.Add(d)
	fakeTimeLock.Unlock()
}

func InitFakeTime() {
	fakeTimeLock.Lock()
	fakeTime = time.Now().UTC()
	fakeTimeLock.Unlock()
}

func SetFakeTime(t time.Time) {
	fakeTimeLock.Lock()
	fakeTime = t
	fakeTimeLock.Unlock()
}
