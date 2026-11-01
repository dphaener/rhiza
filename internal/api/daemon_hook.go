package api

// This file holds the daemon's notify hook: the owner's command, run once
// per new inbox item.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

// notifyTimeout is how long the notify command may run before it is killed.
const notifyTimeout = 10 * time.Second

// notifyTimeoutFor is notifyTimeout, or RHIZA_TEST_NOTIFY_TIMEOUT (a Go
// duration) in tests; like RHIZA_TEST_FREE_BYTES it applies only alongside
// RHIZA_TEST_NOW, so a real daemon never sees it.
func notifyTimeoutFor(ctx *Ctx) time.Duration {
	getenv := ctxGetenv(ctx)
	if v := getenv("RHIZA_TEST_NOTIFY_TIMEOUT"); v != "" && getenv("RHIZA_TEST_NOW") != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return notifyTimeout
}

// daemonHook runs command (through sh -c, in the repo root) once for each
// of items whose ID is not in hooked, and returns the IDs of the items that
// still wait and were hooked, now or before. A failing or timed-out command
// is logged with the item ID and not retried. An empty command runs
// nothing and keeps nothing.
func daemonHook(ctx *Ctx, command string, items []InboxItem, titles map[string]string, hooked []string) []string {
	if command == "" {
		return nil
	}
	var kept []string
	for _, it := range items {
		kept = append(kept, it.ID)
		if slices.Contains(hooked, it.ID) {
			continue
		}
		what := titles[it.Change]
		if what == "" {
			what = it.Change
		}
		if err := runNotifyCommand(ctx, command, it, inboxNotice(what, it)); err != nil && ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "rhiza: daemon: notify command for item %s: %v\n", it.ID, err)
		}
	}
	return kept
}

// runNotifyCommand runs command for item it, with notice as its title.
func runNotifyCommand(ctx *Ctx, command string, it InboxItem, notice string) error {
	in, err := json.Marshal(it)
	if err != nil {
		return err
	}
	return runHookCommand(ctx, command, in, []string{
		"RHIZA_INBOX_KIND=" + it.Kind,
		"RHIZA_INBOX_ID=" + it.ID,
		"RHIZA_INBOX_CHANGE=" + it.Change,
		"RHIZA_INBOX_TITLE=" + notice,
	})
}

// runHookCommand runs command through sh -c in the repo root with in and a
// newline on stdin and env added to the environment, killing it past the
// hook timeout.
func runHookCommand(ctx *Ctx, command string, in []byte, env []string) error {
	timeout := notifyTimeoutFor(ctx)
	c, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(c, "sh", "-c", command)
	cmd.Dir = ctx.Root
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = bytes.NewReader(append(in, '\n'))
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	// A child that outlives the shell keeps the pipes open: stop waiting.
	cmd.WaitDelay = time.Second
	err := cmd.Run()
	if c.Err() != nil {
		return fmt.Errorf("timed out after %s", timeout)
	}
	if err != nil {
		msg := strings.TrimSpace(out.String())
		if msg != "" {
			return fmt.Errorf("%v: %s", err, msg)
		}
		return err
	}
	return nil
}
