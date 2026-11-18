# retry-brief-carries-failure-verified

evidence · `01M7QRH1KSRCXBYQRQHGAGDWHW` · actor [claude-code-91f56e46](../actor/claude-code-91f56e46.md)

## result

PASS

## statement

At c172b4b6c re-ran the Checks: block under the capped scope: go test -count=1 -trimpath -p 2 ./internal/run/... ./internal/guide exit 0; go test -count=1 -trimpath -p 2 ./internal/api -run 'Brief|Retry|Resume|Handback' exit 0; a scratch build's rhiza ship generate (0 changed) and rhiza render exit 0. With -v, TestBriefPreviousAttempt (section after the change set with failed step, outcome, failing command, output tail, findings as check, severity: text; a brief with no previous attempt equals the pre-change todayBrief byte for byte), TestJournalPreviousAttempt (never-retried run has none; only the latest round before the retry), TestJournalPreviousAttemptCheckRounds and TestRunRetryBriefPreviousAttempt (api: run fails at a review-check command step after a findings round, run.resume retry restarts the agent step whose visit-2 brief.md holds the refusal output and 'plan review, major: ...'; the visit-1 session started before the retry has no section, before and after the retry) all PASS. Diff read: only the leaf's files plus actor records; prepareKeyed (leaf/check/spawned sessions) passes nil, handback.md untouched, nothing writes .rhiza/log, retry restart unchanged; no TODOs.

## Links

### Out

- verifies [retry-brief-carries-failure](../capability/rhiza/workflows/workflow-runs/resume-failed-run/retry-brief-carries-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
