package api

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// scriptedRex answers each rex call from replies in turn: an error, or the
// output of a success.
type scriptedRex struct {
	replies []string // "" with fail true is an error; else the output
	fail    []bool
	calls   int
}

func (s *scriptedRex) run(string, []string, string, []string) ([]byte, error) {
	i := s.calls
	s.calls++
	if s.fail[i] {
		return nil, errors.New(s.replies[i])
	}
	return []byte(s.replies[i]), nil
}

func TestSessionRunRexRetryRecovers(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	rex := &scriptedRex{
		replies: []string{"rex is not answering", `{"window_id":"","block_ids":[]}`, rexWindowJSON},
		fail:    []bool{true, false, false},
	}
	var waits []time.Duration
	f.ctx.Exec = rex.run
	f.ctx.Sleep = func(d time.Duration) { waits = append(waits, d) }
	resp := startSession(t, f.ctx, `{"permission_mode":"manual"}`)
	if rex.calls != 3 {
		t.Errorf("rex ran %d times, want 3", rex.calls)
	}
	if len(waits) != 2 || waits[1] <= waits[0] {
		t.Errorf("waits %v, want two growing", waits)
	}
	if len(resp.RexArgv) == 0 {
		t.Errorf("no rex argv in %+v", resp)
	}
}

func TestSessionRunRexRetryGivesUp(t *testing.T) {
	f := newSessionFixture(t, "readable-tabs")
	newChange(t, f.ctx, `{"intent":"Readable tabs"}`)
	rex := &scriptedRex{
		replies: []string{"e1", "e2", "e3", "rex is down"},
		fail:    []bool{true, true, true, true},
	}
	var waits []time.Duration
	f.ctx.Exec = rex.run
	f.ctx.Sleep = func(d time.Duration) { waits = append(waits, d) }
	_, err := callOp(t, f.ctx, "session.start", `{"permission_mode":"manual"}`)
	wantExit(t, err, 1)
	if !strings.Contains(AsError(err).Message, "rex is down") {
		t.Errorf("message %q lacks Rex's last error", AsError(err).Message)
	}
	if rex.calls != 4 || len(waits) != 3 {
		t.Errorf("calls %d, waits %v; want 4 and 3", rex.calls, waits)
	}
}
