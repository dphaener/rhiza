# usage-limit-fresh-session-verified

evidence · `01M65JWK3TTC4J4ST578XA5WS9` · actor [claude-code-ea09c699](../actor/claude-code-ea09c699.md)

## result

PASS

## statement

At commit 2af1b61 ran go test ./internal/api -run TestSessionRunFreshSessionWhenTranscriptGone -count=1 -v: PASS. A fake agent's transcript says 'limited', and the run's sleep during the usage-limit pause deletes that transcript. session.run then launched a second time with the same session ID, Resume=false and the stored first prompt 'do the step', so it started a fresh session instead of resuming. The response's warnings record 'started a fresh session for the step instead of resuming'. Read the diff of internal/api/session_run.go: after the pause, readTranscript decides between resume (transcript present) and fresh launch (transcript missing), and the fresh-session warning is appended to resp.Warnings. Without the change, the second launch would have Resume=true and the test would fail. make check passed (exit 0).

## Links

### Out

- verifies [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
