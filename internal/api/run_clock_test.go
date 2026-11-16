package api

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// stepClock is an injected clock whose Sleep only advances Now: the poll
// count RealClock gives, with no real waiting.
type stepClock struct {
	now    time.Time
	sleeps int
}

func (c *stepClock) Now() time.Time { return c.now }

func (c *stepClock) Sleep(d time.Duration) { c.now = c.now.Add(d); c.sleeps++ }

// scaleGetenv is a getenv whose RHIZA_TEST_CLOCK_SCALE is v and which has
// nothing else.
func scaleGetenv(v string) func(string) string {
	return func(k string) string {
		if k == run.ClockScaleEnv {
			return v
		}
		return ""
	}
}

// envCtx is a Ctx made as the CLI makes one: its clock is the one
// RHIZA_TEST_CLOCK_SCALE=v chooses in its own environment.
func envCtx(v string) *Ctx {
	g := scaleGetenv(v)
	return &Ctx{Getenv: g, Clock: run.ClockFromEnv(g)}
}

// useEnvClock makes the clock RHIZA_TEST_CLOCK_SCALE=v chooses override
// every Ctx's (runClock), for the test.
func useEnvClock(t *testing.T, v string) run.Clock {
	t.Helper()
	c := run.ClockFromEnv(scaleGetenv(v))
	useClock(t, c)
	return c
}

// Two Ctxs in one process, one made with RHIZA_TEST_CLOCK_SCALE=1000 and
// one with it empty, go by a scaled and a real clock while the process's
// environment has neither; an unparsable scale is the real clock; a test
// that sets runClock overrides both.
//
// Serial: it calls t.Setenv (RHIZA_TEST_CLOCK_SCALE) and swaps the
// package's runClock (useClock), process-wide state.
func TestRunClockPerCtx(t *testing.T) {
	t.Setenv(run.ClockScaleEnv, "")
	if err := os.Unsetenv(run.ClockScaleEnv); err != nil {
		t.Fatal(err)
	}
	scaled, real := envCtx("1000"), envCtx("")
	sc, ok := clockOf(scaled).(*run.ScaledClock)
	if !ok {
		t.Fatalf("scale 1000: clock %T, want the scaled clock", clockOf(scaled))
	}
	if c := clockOf(real); c != (run.RealClock{}) {
		t.Errorf("scale empty: clock %T, want RealClock", c)
	}
	for _, v := range []string{"abc", "1.5", "0"} {
		if c := clockOf(envCtx(v)); c != (run.RealClock{}) {
			t.Errorf("scale %q: clock %T, want RealClock", v, c)
		}
	}
	if c := clockOf(&Ctx{}); c != (run.RealClock{}) {
		t.Errorf("no Clock: %T, want RealClock", c)
	}

	// Each Ctx's waits go by its own clock: the scaled one's spawn poll
	// costs little real time and moves its Now by the whole poll; the real
	// one's sleep waits its whole length.
	before, start := sc.Now(), time.Now()
	spawnSleep(scaled)
	if took, moved := time.Since(start), sc.Now().Sub(before); took > time.Second || moved < spawnPoll {
		t.Errorf("scaled Ctx: slept %s of real time, Now moved %s; want under 1s and at least %s", took, moved, spawnPoll)
	}
	start = time.Now()
	ctxSleep(real)(200 * time.Millisecond)
	if took := time.Since(start); took < 200*time.Millisecond {
		t.Errorf("real Ctx: slept %s of real time, want 200ms", took)
	}

	// runClock, when a test sets it, overrides both.
	t0 := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	step := &stepClock{now: t0}
	useClock(t, step)
	for name, ctx := range map[string]*Ctx{"scaled": scaled, "real": real} {
		if c := clockOf(ctx); c != run.Clock(step) {
			t.Errorf("%s Ctx with runClock set: clock %T, want the override", name, c)
		}
	}
	start = time.Now()
	ctxSleep(real)(time.Hour)
	spawnSleep(scaled)
	if took := time.Since(start); took > time.Second || step.sleeps != 2 || step.now.Sub(t0) != time.Hour+spawnPoll {
		t.Errorf("override: %d sleeps, Now moved %s, %s real; want 2 sleeps moving it %s, no real wait", step.sleeps, step.now.Sub(t0), took, time.Hour+spawnPoll)
	}
}

// On the scaled clock a 10 min wait for a session that reaches its prompt
// only after 1 s of real time still sends the message.
//
// Serial: it swaps the package's runClock (useEnvClock), state every run
// in the binary reads.
func TestTypePointerScaledClockWaitsForRealWork(t *testing.T) {
	f := newMessageFixture(t)
	f.setHook(t, session.StateBlocked, session.KindQuestion)
	useEnvClock(t, "1000")
	done := make(chan error, 1)
	time.AfterFunc(time.Second, func() {
		_, err := f.st.Update(f.sid, func(r *session.Record) {
			r.Status.Hook = &session.Report{State: session.StateIdle, At: time.Now().UTC()}
		})
		done <- err
	})
	start := time.Now()
	if _, err := typePointer(f.ctx, f.st, f.sid, "hello", 10*time.Minute); err != nil {
		t.Fatalf("typePointer: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if real := time.Since(start); real < time.Second {
		t.Errorf("sent after %s of real time, before the session reached its prompt", real)
	}
	if len(f.rex.calls) == 0 {
		t.Error("the pointer was not typed")
	}
}

// On the scaled clock a 1 s timeout on a session that never reaches its
// prompt fails only after 1 s of real time.
//
// Serial: it swaps the package's runClock (useEnvClock), state every run
// in the binary reads.
func TestTypePointerScaledClockTimeoutKeepsRealTime(t *testing.T) {
	f := newMessageFixture(t)
	f.setHook(t, session.StateBlocked, session.KindQuestion)
	useEnvClock(t, "1000")
	start := time.Now()
	_, err := typePointer(f.ctx, f.st, f.sid, "hello", time.Second)
	real := time.Since(start)
	if err == nil || !strings.Contains(err.Error(), "did not reach its prompt within 1s") {
		t.Fatalf("err = %v, want the timeout", err)
	}
	if real < time.Second || real > 3*time.Second {
		t.Errorf("timed out after %s of real time, want 1 s to 3 s", real)
	}
	f.noRex(t, "timed out")
}

// Ship deploy measures its deadline on the clock it sleeps on: on the
// scaled clock it polls a probe that never turns healthy as often as
// RealClock would, in far less real time than the timeout.
//
// Serial: it swaps the package's runClock (useEnvClock), state every run
// in the binary reads.
func TestShipDeployScaledClockKeepsPollCount(t *testing.T) {
	srv, n := fakeHealth(503)
	defer srv.Close()
	ctx := func() *Ctx {
		c := deployCtx(srv.URL, "10s", &bytes.Buffer{})
		c.Getenv = func(k string) string {
			if k == "RHIZA_DEPLOY_INTERVAL" {
				return "1s"
			}
			return ""
		}
		return c
	}
	useClock(t, &stepClock{now: time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)})
	r, err := runShipDeploy(ctx(), ShipDeployRequest{})
	if err != nil {
		t.Fatal(err)
	}
	want := r.(ShipDeployResponse).Polls
	if want != 10 {
		t.Fatalf("RealClock rule polled %d times, want 10 in 10s at 1s", want)
	}
	n.Store(0)

	useEnvClock(t, "1000")
	start := time.Now()
	r, err = runShipDeploy(ctx(), ShipDeployRequest{})
	real := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	resp := r.(ShipDeployResponse)
	if resp.Healthy || resp.Polls != want || int(n.Load()) != want {
		t.Errorf("scaled clock: %+v after %d requests, want unhealthy after %d polls", resp, n.Load(), want)
	}
	if real > 3*time.Second {
		t.Errorf("scaled clock took %s of real time for a 10s timeout", real)
	}
}

// With the variable set, the spawn poll's fallback sleeps on the scaled
// clock: a 2 s poll costs a fraction of its real time and moves Now by 2 s.
// With it unset, it is time.Sleep and waits the whole poll.
//
// Serial: it swaps the package's runClock (useEnvClock), state every run
// in the binary reads.
func TestSpawnSleepClock(t *testing.T) {
	c := useEnvClock(t, "1000")
	before, start := c.Now(), time.Now()
	spawnSleep(&Ctx{})
	if real, moved := time.Since(start), c.Now().Sub(before); real > time.Second || moved < spawnPoll {
		t.Errorf("scaled clock: slept %s of real time, Now moved %s; want under 1s and at least %s", real, moved, spawnPoll)
	}

	if c := useEnvClock(t, ""); c != (run.RealClock{}) {
		t.Fatalf("unset: runClock %T, want RealClock", c)
	}
	start = time.Now()
	spawnSleep(&Ctx{})
	if real := time.Since(start); real < spawnPoll {
		t.Errorf("unset: slept %s of real time, want %s", real, spawnPoll)
	}
}

// retryGH and retryRex with no ctx.Sleep sleep on runClock: with the
// variable set their backoff costs little real time and moves the scaled
// clock's Now by the whole wait; with it unset each wait is real.
//
// Serial: it swaps the package's runClock (useEnvClock), state every run
// in the binary reads.
func TestRetryDefaultSleepClock(t *testing.T) {
	transient := errors.New("HTTP 502: Bad Gateway")
	gh := func(fails int) func() ([]byte, error) {
		calls := 0
		return func() ([]byte, error) {
			calls++
			if calls <= fails {
				return nil, transient
			}
			return []byte("ok"), nil
		}
	}
	rex := func(fails int) func() (session.Tab, error) {
		calls := 0
		return func() (session.Tab, error) {
			calls++
			if calls <= fails {
				return session.Tab{}, errors.New("rex is not answering")
			}
			return session.Tab{Window: "w1"}, nil
		}
	}
	var ghTotal, rexTotal time.Duration
	for _, d := range ghWaits {
		ghTotal += d
	}
	for _, d := range rexTabWaits {
		rexTotal += d
	}

	c := useEnvClock(t, "1000")
	before, start := c.Now(), time.Now()
	if _, err := retryGH(&Ctx{}, gh(len(ghWaits)+1)); err == nil || !strings.Contains(err.Error(), "after 3 retries") {
		t.Errorf("retryGH: %v, want it to give up", err)
	}
	if real, moved := time.Since(start), c.Now().Sub(before); real > time.Second || moved < ghTotal {
		t.Errorf("retryGH on the scaled clock: %s real, Now moved %s; want under 1s and at least %s", real, moved, ghTotal)
	}
	before, start = c.Now(), time.Now()
	if _, err := retryRex(&Ctx{}, rex(len(rexTabWaits)+1)); err == nil || !strings.Contains(err.Error(), "after 3 retries") {
		t.Errorf("retryRex: %v, want it to give up", err)
	}
	if real, moved := time.Since(start), c.Now().Sub(before); real > time.Second || moved < rexTotal {
		t.Errorf("retryRex on the scaled clock: %s real, Now moved %s; want under 1s and at least %s", real, moved, rexTotal)
	}

	useEnvClock(t, "")
	start = time.Now()
	if _, err := retryGH(&Ctx{}, gh(1)); err != nil {
		t.Errorf("retryGH: %v", err)
	}
	if real := time.Since(start); real < ghWaits[0] {
		t.Errorf("retryGH unset: %s real, want its first wait %s", real, ghWaits[0])
	}
	start = time.Now()
	if _, err := retryRex(&Ctx{}, rex(1)); err != nil {
		t.Errorf("retryRex: %v", err)
	}
	if real := time.Since(start); real < rexTabWaits[0] {
		t.Errorf("retryRex unset: %s real, want its first wait %s", real, rexTabWaits[0])
	}
}
