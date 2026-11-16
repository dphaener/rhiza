package run

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

func TestScaledClockSleepSkipsTime(t *testing.T) {
	c := NewScaledClock(1000)
	before := c.Now()
	start := time.Now()
	c.Sleep(30 * time.Second)
	real := time.Since(start)
	if real >= 100*time.Millisecond || real < ScaledSleepFloor {
		t.Errorf("Sleep(30s) took %v of real time, want [50ms, 100ms)", real)
	}
	if adv := c.Now().Sub(before); adv < 30*time.Second {
		t.Errorf("Now advanced %v, want at least 30s", adv)
	}
}

func TestScaledClockShortSleepIsReal(t *testing.T) {
	c := NewScaledClock(1000)
	before := c.Now()
	start := time.Now()
	c.Sleep(50 * time.Millisecond)
	real := time.Since(start)
	adv := c.Now().Sub(before)
	if real < 50*time.Millisecond || real >= 100*time.Millisecond {
		t.Errorf("Sleep(50ms) took %v of real time, want about 50ms", real)
	}
	if adv < 50*time.Millisecond || adv >= 100*time.Millisecond {
		t.Errorf("Sleep(50ms) advanced Now %v, want about 50ms", adv)
	}
}

func TestScaledClockRealTimeOutsideSleep(t *testing.T) {
	c := NewScaledClock(1000)
	before := c.Now()
	time.Sleep(50 * time.Millisecond)
	adv := c.Now().Sub(before)
	if adv < 50*time.Millisecond || adv >= time.Second {
		t.Errorf("50ms of real time advanced Now %v, want about 50ms", adv)
	}
}

func TestScaledClockConcurrent(t *testing.T) {
	c := NewScaledClock(1000)
	before := c.Now()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 3 {
				c.Sleep(10 * time.Second)
				_ = c.Now()
			}
		}()
	}
	wg.Wait()
	// 24 sleeps of 10s each add their skipped time.
	if adv := c.Now().Sub(before); adv < 200*time.Second {
		t.Errorf("Now advanced %v, want at least 200s", adv)
	}
}

// scaleEnv is a getenv over the variables in vars.
func scaleEnv(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

func TestScaledClockFromEnv(t *testing.T) {
	for _, v := range []string{"", "0", "-5", "1.5", "abc"} {
		if got, ok := ClockFromEnv(scaleEnv(map[string]string{ClockScaleEnv: v})).(RealClock); !ok {
			t.Errorf("%s=%q gives %T, want RealClock", ClockScaleEnv, v, got)
		}
	}
	got := ClockFromEnv(scaleEnv(map[string]string{ClockScaleEnv: "1000"}))
	if c, ok := got.(*ScaledClock); !ok || c.scale != 1000 {
		t.Errorf("%s=1000 gives %#v, want a ScaledClock of scale 1000", ClockScaleEnv, got)
	}
}

// ClockFromEnv reads only the environment it is given: the process's
// RHIZA_TEST_CLOCK_SCALE does not reach it.
//
// Serial: it calls t.Setenv (RHIZA_TEST_CLOCK_SCALE), process-wide state.
func TestScaledClockFromEnvIgnoresProcess(t *testing.T) {
	t.Setenv(ClockScaleEnv, "1000")
	if _, ok := ClockFromEnv(scaleEnv(nil)).(RealClock); !ok {
		t.Errorf("unset in getenv, 1000 in the process: %T, want RealClock", ClockFromEnv(scaleEnv(nil)))
	}
	if _, ok := ClockFromEnv(nil).(RealClock); !ok {
		t.Errorf("nil getenv gives %T, want RealClock", ClockFromEnv(nil))
	}
	if err := os.Unsetenv(ClockScaleEnv); err != nil {
		t.Fatal(err)
	}
	if _, ok := ClockFromEnv(scaleEnv(map[string]string{ClockScaleEnv: "7"})).(*ScaledClock); !ok {
		t.Error("7 in getenv, unset in the process: want a ScaledClock")
	}
}

func TestTimedOutScaledKeepsRealFloor(t *testing.T) {
	c := NewScaledClock(1000)
	w := BeginWait(c)
	c.Sleep(30 * time.Second)
	if w.TimedOut(10 * time.Minute) {
		t.Fatal("10m timed out after 30s of virtual time")
	}
	for w.Elapsed() < 10*time.Minute {
		c.Sleep(30 * time.Second)
	}
	if real := time.Since(w.realBegan); real >= TimeoutRealFloor {
		t.Fatalf("reaching 10m of virtual time took %v of real time", real)
	}
	if w.TimedOut(10 * time.Minute) {
		t.Fatal("10m timed out before 3s of real time")
	}
	time.Sleep(TimeoutRealFloor - time.Since(w.realBegan))
	if !w.TimedOut(10 * time.Minute) {
		t.Fatal("10m did not time out after 10m virtual and 3s real")
	}
}

func TestTimedOutScaledShortTimeoutWaitsItsRealTime(t *testing.T) {
	c := NewScaledClock(1000)
	w := BeginWait(c)
	c.Sleep(30 * time.Second)
	if w.TimedOut(time.Second) {
		t.Fatal("1s timed out before 1s of real time")
	}
	time.Sleep(time.Second - time.Since(w.realBegan))
	if !w.TimedOut(time.Second) {
		t.Fatal("1s did not time out after 1s of real time")
	}
}

func TestTimedOutFakeClockIsExact(t *testing.T) {
	c := &fakeClock{now: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	start := time.Now()
	w := BeginWait(c)
	c.Sleep(10*time.Minute - time.Second)
	if w.TimedOut(10 * time.Minute) {
		t.Error("timed out at 9m59s")
	}
	c.Sleep(time.Second)
	if !w.TimedOut(10 * time.Minute) {
		t.Error("not timed out at 10m")
	}
	if real := time.Since(start); real >= 100*time.Millisecond {
		t.Errorf("fake clock waited %v of real time", real)
	}
}

func TestTimedOutRealClockIsExact(t *testing.T) {
	w := BeginWait(RealClock{})
	if !w.TimedOut(0) {
		t.Error("a 0 timeout has not passed on RealClock")
	}
	if w.TimedOut(time.Hour) {
		t.Error("1h passed at once on RealClock")
	}
}

// laterRunner holds every wait condition once its time has come, in real
// time.
type laterRunner struct {
	*fakeRunner
	at time.Time
}

func (l laterRunner) Holds(step workflow.Step, env []string) (bool, error) {
	l.polls[step.Name]++
	return !l.at.IsZero() && !time.Now().Before(l.at), nil
}

const scaledWaitWF = `
name = "wf"
[step.ci]
kind = "wait"
until = "background work"
interval = "30s"
timeout = "10m"
on.pass = "done"
on.timeout = "failed"
`

func TestWaitScaledClockHoldsAfterRealWork(t *testing.T) {
	e, r, _, _ := setup(t, scaledWaitWF)
	e.Clock = NewScaledClock(1000)
	e.Runner = laterRunner{r, time.Now().Add(time.Second)}
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	if st.Result != Succeeded || len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomePass {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
}

func TestWaitScaledClockTimesOutInRealSeconds(t *testing.T) {
	e, r, _, _ := setup(t, scaledWaitWF)
	e.Clock = NewScaledClock(1000)
	e.Runner = laterRunner{r, time.Time{}}
	start := time.Now()
	st, err := e.Run(runID)
	if err != nil {
		t.Fatal(err)
	}
	real := time.Since(start)
	if st.Result != Failed || len(st.Steps) != 1 || st.Steps[0].Outcome != workflow.OutcomeTimeout {
		t.Fatalf("result %q steps %v", st.Result, outcomes(st))
	}
	if real < TimeoutRealFloor || real >= 4*time.Second {
		t.Errorf("timed out after %v of real time, want [3s, 4s)", real)
	}
	// About 20 polls reach 10m of virtual time; the rest of the 3s floor
	// polls every 50ms, not in a busy loop.
	if n := r.polls["ci"]; n < 20 || n > 80 {
		t.Errorf("polls = %d, want 20 to 80", n)
	}
}
