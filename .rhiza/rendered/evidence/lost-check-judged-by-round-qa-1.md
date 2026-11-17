# lost-check-judged-by-round-qa-1

evidence · `01M7MJC1V454WRR3E5JXD07WZY` · actor [claude-code-596d20c6](../actor/claude-code-596d20c6.md)

## result

PASS

## statement

QA at ed7678fe4 with a fresh go build -trimpath ./cmd/rhiza in scratch (outside PATH), in hermetic scratch repos with a three-check review workflow, a fake rex that runs the tab command and answers block inspect 'not found', and a fake claude that records its round or verdicts then kill -9s its rhiza session run wrapper (a crash: no exit recorded); each run killed, then rhiza run resume --foreground. (1) All three checks recorded empty rounds then died: resume judged each 'pass (lost: Rex block not found; findings round recorded)', run succeeded in cycle 1, no new session started. (2) alpha died before recording a round: alpha judged fail (lost), cycle 2 ran alpha alone; beta and gamma passed on their rounds. (3) alpha's round held a major and a blocker: resume started alpha's critic; the critic recorded dropped verdicts on both then died; second resume judged the critic 'pass (lost; verdicts recorded)' and alpha passed on its verdicts. (4) Same with the critic recording a verdict on the major only: critic judged fail, check detail 'cycle 1: critic failed: lost: Rex block not found', cycle 2 ran alpha. (5) No crash, alpha exits 1 after recording its round: judged 'fail (exit 1)', cycle 2 ran alpha. go test -p 2 -trimpath ./internal/run/... ok; cmd/rhiza scripts review-checks-step, run-resume-lost-session, run-resume-lost-fixer, review-findings exit 0.

## Links

### Out

- verifies [lost-check-judged-by-round](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step/lost-check-judged-by-round.md)

## History

- `01M7MGF8J9MA30XXQ8GTR4KXRW` Fix Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes
