# question-sessions-tile-run-window-verified

evidence · `01M68GZN0M03KTQQ0S0Q7JKGA4` · actor [claude-code-655fa71d](../actor/claude-code-655fa71d.md)

## result

PASS

## statement

At commit 2588c4b, go test ./internal/api -run 'TestStartArbiterSplits|TestStartAdvisorSplits' -count=1 passes; both tests seed the run with a session whose tab block is block:a1 (holdPane) and assert via wantSplitOf that the fake rex saw split -b block:a1 and no window new. go test ./cmd/rhiza -run TestScripts/arbiter-panes-one-window -count=1 passes: an autonomous plan run's fake claude asks two questions, two arbiter sessions answer them, and the fake rex log shows exactly one window new and four split -s session:s1 calls (two arbiters, reviewer, critic). make test passes for every package. make self-check flagged rendered files only for two actors whose session records were made after this commit (this verifier's own session and one other live session), so the commit is not the cause.

## Links

### Out

- verifies question-sessions-tile-run-window (retired)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
