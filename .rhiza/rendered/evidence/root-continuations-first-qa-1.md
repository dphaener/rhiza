# root-continuations-first-qa-1

evidence · `01M6XYC0C3BACJ4730ARSMZKNY` · actor [claude-code-c038643d](../actor/claude-code-c038643d.md)

## result

PASS

## statement

QA at 728608f3 with the built binary (go build ./cmd/rhiza into the session scratchpad), hands-on in a scratch root outside ~/code with two registered projects (a, b), fake rex, RHIZA_* unset, hermetic HOME. Project a (first in turn order) had two accepted milestone change sets ready; project b had a failed flaky-workflow run queued with rhiza run resume --retry. rhiza daemon --max-runs 1 --once at the root printed only 'Resumed run <b's run>', no 'Started a run of' for a. Then in b: one live running run (journal naming a live sleep pid) and two interrupted runs queued with rhiza run resume. With --max-runs 2 the one free slot went to b's oldest continuation (Resumed run QUEUEDONE01), not a's change sets; with --max-runs 4 the picks were QUEUEDTWO02 plus a's One and Two. go test -trimpath ./internal/daemon/... ok; TestScripts daemon-continue-slot, daemon-round-robin, describe-golden PASS.

## Links

### Out

- verifies [root-continuations-first](../capability/rhiza/orchestration/daemon/root-continuations-first.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
