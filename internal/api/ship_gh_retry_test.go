package api

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// failingCreates wraps ctx.Exec so the first n gh pr create calls fail
// with msg and the rest go to the fake gh. It returns the count of pr
// create calls seen and the waits slept.
func failingCreates(ctx *Ctx, n int, msg string) (*int, *[]time.Duration) {
	inner := ctx.Exec
	creates := 0
	var waits []time.Duration
	ctx.Sleep = func(d time.Duration) { waits = append(waits, d) }
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name == "gh" && len(args) > 1 && args[0] == "pr" && args[1] == "create" {
			creates++
			if creates <= n {
				return nil, errors.New(msg)
			}
		}
		return inner(name, args, dir, env)
	}
	return &creates, &waits
}

func TestShipGHRetryRecovers(t *testing.T) {
	for _, msg := range []string{
		"gh pr create: HTTP 502: Bad Gateway",
		"GraphQL: Something went wrong while executing your query",
	} {
		ctx, _, _, _ := openPRFixture(t, true)
		creates, waits := failingCreates(ctx, 2, msg)
		if _, err := callOp(t, ctx, "ship.open-pr", `{}`); err != nil {
			t.Fatalf("%q: %v", msg, err)
		}
		if *creates != 3 {
			t.Errorf("%q: %d pr create calls, want 3", msg, *creates)
		}
		if len(*waits) != 2 || (*waits)[1] <= (*waits)[0] {
			t.Errorf("%q: waits %v, want two growing", msg, *waits)
		}
	}
}

func TestShipGHRetryGivesUp(t *testing.T) {
	ctx, _, _, _ := openPRFixture(t, true)
	creates, waits := failingCreates(ctx, 10, "gh pr create: HTTP 503: unavailable")
	_, err := callOp(t, ctx, "ship.open-pr", `{}`)
	wantExit(t, err, 1)
	if !strings.Contains(AsError(err).Message, "HTTP 503") {
		t.Errorf("message %q lacks the last error", AsError(err).Message)
	}
	if *creates != 4 || len(*waits) != 3 {
		t.Errorf("creates %d, waits %v; want 4 and 3", *creates, *waits)
	}
}

func TestShipGHRetrySkipsOtherErrors(t *testing.T) {
	for _, msg := range []string{"gh pr create: HTTP 422: Validation Failed", `no pull requests found for branch "x"`} {
		ctx, _, _, _ := openPRFixture(t, true)
		creates, waits := failingCreates(ctx, 10, msg)
		_, err := callOp(t, ctx, "ship.open-pr", `{}`)
		wantExit(t, err, 1)
		if *creates != 1 || len(*waits) != 0 {
			t.Errorf("%q: creates %d, waits %v; want 1 and none", msg, *creates, *waits)
		}
	}
}
