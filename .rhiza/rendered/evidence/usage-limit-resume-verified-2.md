# usage-limit-resume-verified-2

evidence · `01M65KZYQHBBFGFWVE5JM9EA4B` · actor [claude-code-bee34efe](../actor/claude-code-bee34efe.md)

## result

PASS

## statement

Verifier at 4ba849e (attempt 2): go test -count=1 ./cmd/rhiza -run TestScripts/session-run-usage-limit passes. Its fake claude writes 'Claude AI usage limit reached|<past reset>' on its first launch and finishes on the second. session run reports 'ended: the agent exited with code 0', and claude.log holds exactly one 'start --session-id <id>' and one 'resume --resume <id>' for the same stored session ID, so no second session starts. go test -race -count=20 of TestSessionRunStopsAgentOpenAtLimit, TestSessionRunFailsWhenPausesRunOut and TestSessionRunFreshSessionWhenTranscriptGone passes. These tests cover the review findings: an agent left open at its limit is killed by the transcript watcher and resumed, and running out of the 24 pauses fails the session with a StateError exit instead of passing. Decision usage-limit-detect-and-wait-answer answers the blocking question. Diff read: watchLimit polls the transcript past its pre-launch length and kills p, and the loop now records a failure at the cap. The first make check failed only in the test target right after the timing-budget output, a known timing-contention flake; a rerun passed clean (exit 0), and go test ./... passed.

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
