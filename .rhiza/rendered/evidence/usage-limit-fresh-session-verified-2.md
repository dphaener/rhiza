# usage-limit-fresh-session-verified-2

evidence · `01M65M74HXENCPYV7X14FP4DJ9` · actor [claude-code-80bfa814](../actor/claude-code-80bfa814.md)

## result

PASS

## statement

At 4a8f5f3, ran go test ./internal/api -run TestSessionRunFreshSessionWhenTranscriptGone -count=1 -v: PASS. The test deletes the stored session's transcript during the usage-limit pause. session.run then launches without resume, using the session's first prompt, warns 'started a fresh session'. It also records the restart: the session record's new Fresh field holds one timestamp, and the session's run journal holds one fresh-session entry with the session ID, the step (build) and that detail. Read the diff of 2af1b61 and 4a8f5f3: recordFreshSession in internal/api/session_run.go stamps Record.Fresh and appends run.TypeFreshSession, which is documented in internal/run/journal.go. Both commits carry Leaf: usage-limit-fresh-session. make fmt-check, vet and test pass (the first test run hit one failure that did not come back on a rerun). make self-check fails only on stale .rhiza/rendered/actor files for actors recorded during this run, not on the leaf's code.

## Links

### Out

- verifies [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
