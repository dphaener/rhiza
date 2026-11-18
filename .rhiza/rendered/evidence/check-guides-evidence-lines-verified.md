# check-guides-evidence-lines-verified

evidence · `01M7PYS34K5YZ094WHE011P8M9` · actor [claude-code-eaeb5d13](../actor/claude-code-eaeb5d13.md)

## result

PASS

## statement

At 1bfa91e9d: go test -trimpath -p 2 -count=1 ./internal/guide/... exit 0 (incl. TestCaps); rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only untracked session actor lines in index.md). Diff read: check-plan/check-code/check-qa show {"findings":[...],"evidence":[{"leaf","line"}]}, say one evidence line per acceptance point of each leaf, and that an uncovered leaf is refused in a conductor run; new 'coverage' assertions in check_plan/code/qa_test.go cover each; check_qa_test still asserts the rhiza apply Evidence example and check-qa.md keeps it. check-compliance/quality/model and Go code outside internal/guide tests untouched; no TODOs.

## Links

### Out

- verifies [check-guides-evidence-lines](../capability/rhiza/workflows/built-in-workflows/built-in-review/check-guides-evidence-lines.md)

## History

- `01M7PT8436SJQWKMBE1V1M26W6` Fix Shallow checks pass unchallenged: a conductor run's check can finish in 1-2 minutes with 'no findings', no summary and no Evidence, and the run accepts it. Seen: the plan check of run 01M7K2AFC78VJSQK60DZJ51YGR (One read path, 5 leaves deleting the interim model) took 80 s with no summary (a rerun the orchestrator asked for took 9 min); the QA check of run 01M7KHWS2DRWVTKAKNSZ99TR7W (Leaf contract) took 2 min with no findings and no Evidence on any leaf, where earlier QA rounds took 13-22 min. Expected: rhiza run findings --check (or review-check) refuses an empty round that carries no per-leaf evidence (QA: Evidence per leaf; plan check: one line of evidence per acceptance point), so the conductor reruns the check.
