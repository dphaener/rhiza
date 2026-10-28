# review-signoff-lock-verified

evidence · `01M606757QQESRSJRDY5WCASJV` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

At a201949 (includes 0783bed sign-off under the run's lock and f194355 tidy): review-signoff.txtar shows a review run stopped at signoff with signoff_waiting true until rhiza run approve, then succeeded with one signoffs entry {actor test, outcome approved, rule owner, step signoff, visit 1}; review-signoff-changes.txtar shows run request-changes with one finding on leaf thing, after check alpha used cycles 1-3, start fixer for thing (cycle 4, findings=yes, commit Leaf: thing), its verifier, alpha re-run, and the run waiting at signoff again, a third request refused exit 2 with the journal unchanged; review-signoff.txtar shows approve and request-changes exit 2 with no run, a run stopped at another step, and a finished run, journal lines unchanged; review-signoff-autonomous.txtar shows an --autonomous run pass signoff with a journal entry rule autonomous, outcome approved, no actor; describe-golden matches. internal/api TestSignoffUnderLockHeld, TestSignoffRecheckedUnderLock, TestSignoffAlreadySignedUnderLock, TestSignoffUnderLockJournalsOnce pass; git show 0783bed f194355 read: Engine.Continue takes the lock, re-reads the journal and runs the sign-off checks (signoffWaitingStep plus same step/visit) before appending, a refusal writes nothing. Hands-on with a scratch build and fake claude/rex: 12 trials of concurrent rhiza run approve + rhiza run request-changes on a review run waiting at sign-off (launch order alternated). Every trial journalled exactly one type=signoff entry for visit 1, and its outcome was that of the command that exited 0 (6 approve wins: run succeeded; 6 request-changes wins: fixer for beta ran, run waits at signoff visit 2). The loser exited 3 (running in another process) or 2 (finished, or running, from the pre-lock check) and wrote no sign-off.

## tests

commit a201949; run by claude-code-0e5a9d92-verifier on 2026-10-28; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT for each of: go test ./cmd/rhiza -run 'TestScripts/review-signoff' -v -count=1; go test ./internal/run -run TestSignoff -v -count=1; go test ./internal/workflow -run TestSignoffStep -v -count=1; go test ./internal/api -run 'TestSignoff' -v -count=1; go test ./cmd/rhiza -run TestScripts/describe-golden -count=1; go test ./... -count=1 (all ok); go build -o $TMPDIR/rhiza-vs2/rhiza ./cmd/rhiza and a 12-trial concurrent approve/request-changes script, output in .shipit/evidence/reverify-signoff.txt

## Links

### Out

- verifies [review-signoff](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-signoff.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
