package run

import (
	"strconv"
	"sync"
	"time"
)

// ClockScaleEnv names the variable that puts a process's waits on a scaled
// clock: RHIZA_TEST_CLOCK_SCALE=<n>, a positive integer. Testscripts set
// it; production leaves it unset and keeps RealClock.
const ClockScaleEnv = "RHIZA_TEST_CLOCK_SCALE"

// ScaledSleepFloor is the least real time a scaled Sleep waits, unless it
// asks for less: no poll on the scaled clock runs more than about 20 times
// a second of real time.
const ScaledSleepFloor = 50 * time.Millisecond

// TimeoutRealFloor caps the real time a timeout on the scaled clock keeps:
// it fires only once the smaller of the timeout and this much real time
// has passed since the wait began.
const TimeoutRealFloor = 3 * time.Second

// ClockFromEnv returns the clock ClockScaleEnv chooses in the environment
// getenv reads: a ScaledClock when it holds a positive integer, RealClock
// when it is unset, empty or anything else. A nil getenv is an empty
// environment.
func ClockFromEnv(getenv func(string) string) Clock {
	if getenv == nil {
		return RealClock{}
	}
	return clockForScale(getenv(ClockScaleEnv))
}

func clockForScale(v string) Clock {
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return RealClock{}
	}
	return NewScaledClock(n)
}

// ScaledClock is a Clock whose Sleep(d) waits max(d/n, min(d, 50 ms)) of
// real time, and whose Now is the wall clock plus the time its Sleeps
// skipped. Real time spent outside Sleep advances Now by itself only, so a
// timeout on real work is not cut by the scale. It is safe for concurrent
// use.
type ScaledClock struct {
	scale   int64
	mu      sync.Mutex
	skipped time.Duration
}

// NewScaledClock returns a ScaledClock with scale n (at least 1).
func NewScaledClock(n int) *ScaledClock {
	return &ScaledClock{scale: int64(max(n, 1))}
}

// Now returns the wall clock plus the time Sleeps have skipped.
func (c *ScaledClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return time.Now().Add(c.skipped)
}

// Sleep waits max(d/n, min(d, ScaledSleepFloor)) of real time and adds the
// rest of d to Now.
func (c *ScaledClock) Sleep(d time.Duration) {
	if d <= 0 {
		return
	}
	real := max(d/time.Duration(c.scale), min(d, ScaledSleepFloor))
	start := time.Now()
	time.Sleep(real)
	waited := time.Since(start)
	if skip := d - waited; skip > 0 {
		c.mu.Lock()
		c.skipped += skip
		c.mu.Unlock()
	}
}

// Wait measures one wait's timeout on a Clock. Begin it when the wait
// starts polling.
type Wait struct {
	c         Clock
	began     time.Time
	realBegan time.Time
}

// BeginWait starts measuring a wait on c.
func BeginWait(c Clock) Wait {
	return Wait{c: c, began: c.Now(), realBegan: time.Now()}
}

// Elapsed is the time since the wait began, on its clock.
func (w Wait) Elapsed() time.Duration {
	return w.c.Now().Sub(w.began)
}

// TimedOut reports whether timeout has passed. On RealClock or any clock a
// test injects it is exactly Elapsed() >= timeout. On a ScaledClock the
// timeout must also have passed min(timeout, TimeoutRealFloor) of real
// time, so a wait on another process's real work keeps that much of it.
func (w Wait) TimedOut(timeout time.Duration) bool {
	if w.Elapsed() < timeout {
		return false
	}
	if _, ok := w.c.(*ScaledClock); !ok {
		return true
	}
	return time.Since(w.realBegan) >= min(timeout, TimeoutRealFloor)
}
