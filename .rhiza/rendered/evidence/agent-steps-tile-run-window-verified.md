# agent-steps-tile-run-window-verified

evidence · `01M68F4EHBW8HXBAXGNJPC5AHN` · actor [claude-code-1d55cef9](../actor/claude-code-1d55cef9.md)

## result

PASS

## statement

At 6f62e29 I read the diff: StartAgent, StartArbiter and StartAdvisor in internal/api/run.go now set req.split and req.anchor = r.paneAnchor(""), and session start splits whenever the anchor has a block, else opens a window. I ran go test ./internal/api -run 'TestStartAgent|TestStartCheck|TestStartFix' -count=1 -v: all pass. The new TestStartAgent in internal/api/run_checks_panes_test.go uses paneRunner, calls StartAgent for a plan interview step after session S-first with block:a1, and wantSplitOf asserts one rex split -b block:a1 and no window new. TestStartAgentFirstSession asserts exactly one window new and no split. It has no session of its own run (the store holds only another run's session, not an empty store), which reaches the same paneAnchor path that returns an empty Tab. Both use the fake rex exec of sessionFixture, so no real Rex runs. make check exits 0 with a clean tree.

## Links

### Out

- verifies [agent-steps-tile-run-window](../capability/rhiza/runtime/rex-sessions/agent-steps-tile-run-window.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
