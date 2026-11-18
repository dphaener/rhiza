# check-round-evidence-status-verified

evidence · `01M7PYDYE80VVRMPY3WETG7W82` · actor [claude-code-e0c93490](../actor/claude-code-e0c93490.md)

## result

PASS

## statement

At 9aaad711f: go test -count=1 -trimpath -p 2 ./internal/api -run 'RunStatus|CheckRound' exit 0 (new TestRunStatusCheckRoundEvidence passes: two evidence lines with handles in JSON and human output under the round, no evidence key on an evidence-less round, findings round printed as before, retired leaf-c shown with ID, empty leaf_handle, bare ID in human text, status succeeds); go test ./cmd/rhiza -run TestScripts/describe-golden$ exit 0 (golden carries the evidence schema); in a clean detached worktree make self-check exit 0 and rhiza ship generate 0 files changed. Diff read: only RunCheckRound/RunEvidence, evidenceHandles and the Human round loop change (findings/verdicts/plan-review output unchanged, the no-findings continue becomes an equivalent else); run findings validation untouched; no TODOs.

## Links

### Out

- verifies [check-round-evidence-status](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-round-evidence-status.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
